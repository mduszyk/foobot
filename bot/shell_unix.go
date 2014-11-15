package bot

import(
    "io"
    "os"
    "os/exec"
    "strings"
    "github.com/kr/pty"
	"github.com/mduszyk/foobot/log"
	"github.com/mduszyk/foobot/conf"
)

const PS1 = "--FoObOt--"

type Shell struct {
    proc *exec.Cmd
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

    proc := exec.Command(chunks[0], args...)

    fdm, err := pty.Start(proc)
    if err != nil {
        log.ERROR.Printf("Failed starting shell pty: %s", err)
        return
    }

    log.TRACE.Printf("Started shell: %s", shell)

    sh.proc = proc
    sh.pty = fdm

    sh.setupPrompt(PS1)
}

func (sh *Shell) setupPrompt(ps1 string) {
    // setup shell prompt
    pscmd := "export PS1=\"" + ps1 + "\"\n"
    sh.pty.Write([]byte(pscmd))
    readBetween(sh.pty, "", ps1 + "\"\r\n" + ps1)
    sh.ps1 = ps1
    log.TRACE.Printf("Prompt ready: %s", ps1)
}

func readBetween(r io.Reader, token1 string, token2 string) string {
    buf := make([]byte, 256)
    var outBuf []byte
    start := -1
    end := 0

    for {
        n, _ := r.Read(buf)
        end += n
        outBuf = append(outBuf, buf[:n]...)
        text := string(outBuf[:end])
        /* log.TRACE.Printf("text: %s, token2: %s", text, token2) */

        if start < 0 {
            start = strings.Index(text, token1) + len(token1)
            /* log.TRACE.Printf("start: %d", start) */
        }
        if strings.Contains(text, token2) {
            break
        }
    }

    /* log.TRACE.Printf("out: %s", string(outBuf[:end])) */

    end -= len(token2)
    return string(outBuf[start:end])
}

func (sh *Shell) Insert(line string) string {
    if sh.proc == nil {
        sh.Start()
    }

    sh.pty.Write([]byte(line + "\n"))

    return readBetween(sh.pty, "\n", sh.ps1)
}

