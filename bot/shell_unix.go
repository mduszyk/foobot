package bot

/*
#include <stdlib.h>
#include <sys/select.h>

static void _FD_ZERO(void *set) {
    FD_ZERO((fd_set*)set);
}

static void _FD_SET(int sysfd, void *set) {
    FD_SET(sysfd, (fd_set*)set);
}

static int _FD_ISSET (int sysfd, void *set) {
    return FD_ISSET(sysfd, (fd_set*)set);
}
*/
import "C"

import(
    "os"
    "os/exec"
    "syscall"
    "unsafe"
    "strings"
    "github.com/kr/pty"
    "github.com/mduszyk/foobot/log"
    "github.com/mduszyk/foobot/conf"
)

func FD_ZERO(set *syscall.FdSet) {
    s := unsafe.Pointer(set)
    C._FD_ZERO(s)
}

func FD_SET(sysfd int, set *syscall.FdSet) {
    s := unsafe.Pointer(set)
    fd := C.int(sysfd)
    C._FD_SET(fd, s)
}

func FD_ISSET(sysfd int, set *syscall.FdSet) bool {
    s := unsafe.Pointer(set)
    fd := C.int(sysfd)
    return C._FD_ISSET(fd, s) != 0
}

const PS1 = "--FoObOt--"


type Shell struct {
    cmd *exec.Cmd
    pty *os.File
    ps1 string
}

func NewShell() *Shell {
    return &Shell{}
}

func (sh *Shell) Start() {
    shell := conf.Get("bot.shell")

    chunks := strings.Split(shell, " ")
    var args []string
    if len(chunks) > 1 {
        args = chunks[1:]
    } else {
        args = []string{}
    }

    cmd := exec.Command(chunks[0], args...)

    fdm, err := pty.Start(cmd)
    if err != nil {
        log.ERROR.Printf("Failed starting shell pty: %s", err)
        return
    }

    log.TRACE.Printf("Started shell: %s", shell)

    sh.cmd = cmd
    sh.pty = fdm

    sh.setupPrompt(PS1)
}

func (sh *Shell) Kill() {
    if sh.cmd != nil {
        sh.cmd.Process.Kill()
    }
}

func (sh *Shell) Interrupt() {
    if sh.cmd != nil {
        sh.cmd.Process.Signal(syscall.SIGINT)
    }
}

func (sh *Shell) setupPrompt(ps1 string) {
    // setup shell prompt
    pscmd := "export PS1=\"" + ps1 + "\"\n"
    sh.pty.Write([]byte(pscmd))
    readBetween(sh.pty, "", ps1 + "\"\r\n" + ps1)
    sh.ps1 = ps1
    log.TRACE.Printf("Prompt ready: %s", ps1)
}

func readBetween(f *os.File, token1 string, token2 string) string {
    buf := make([]byte, 256)
    var outBuf []byte
    start := 0
    end := 0
    first := true

    fd := int(f.Fd())
    /* log.TRACE.Printf("Read, fd: %d", fd) */

    rfds := &syscall.FdSet{}

    timeout := &syscall.Timeval{}
    timeout.Sec = 0
    timeout.Usec = 500000

    for {
        FD_ZERO(rfds)
        FD_SET(fd, rfds)

        n, err := syscall.Select(fd + 1, rfds, nil, nil, timeout)
        if err != nil {
            if err == syscall.EINTR {
                log.TRACE.Printf("Select interrupted")
                continue
            } else {
                log.ERROR.Printf("Select failed, error: %s", err)
                break
            }
        }
        if !FD_ISSET(fd, rfds) {
            log.TRACE.Printf("Select timeout")
            break
        }

        n, err = f.Read(buf)
        if err != nil {
            log.ERROR.Printf("Failed reading shell output, error: %s", err)
            break
        }

        end += n
        outBuf = append(outBuf, buf[:n]...)
        text := string(outBuf[:end])
        /* log.TRACE.Printf("text: %s, token2: %s", text, token2) */

        if first {
            start = strings.Index(text, token1) + len(token1)
            first = false
            /* log.TRACE.Printf("start: %d", start) */
        }
        if strings.Contains(text, token2) {
            end -= len(token2)
            break
        }
    }

    /* log.TRACE.Printf("out: %s", string(outBuf[:end])) */
    return string(outBuf[start:end])
}

func (sh *Shell) Insert(line string) string {
    if sh.cmd == nil {
        sh.Start()
    }

    sh.pty.Write([]byte(line + "\n"))

    return readBetween(sh.pty, "\n", sh.ps1)
}
