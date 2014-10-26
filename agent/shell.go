package agent

import(
    "io"
    "os/exec"
    "strings"
	"fuzzywookie/foobot/log"
)

type Shell struct {
    proc *exec.Cmd
    stdin io.WriteCloser
    stdout io.ReadCloser
}

func NewShell() *Shell {
    shell := "/bin/bash"
    proc := exec.Command(shell)
    in, err := proc.StdinPipe()
    if err != nil {
        log.ERROR.Printf("Error connecting shell stdin: %s", err)
        return nil
    }
    out, err := proc.StdoutPipe()
    if err != nil {
        log.ERROR.Printf("Error connecting shell stdout: %s", err)
        return nil
    }
    proc.Start()
    log.TRACE.Printf("Started shell: %s", shell)

    return &Shell{
		proc: proc,
        stdin: in,
		stdout: out,
	}
}

// TODO improve this
func (sh *Shell) Insert(line string) string {
    sh.stdin.Write([]byte(line + "; echo -e '\\x63\\x68\\x65\\x63\\x6b'\n"))
    var rsp string
    buf := make([]byte, 256)
    var rspBuf []byte
    l := 0
    for {
        n, _ := sh.stdout.Read(buf)
        l += n
        rspBuf = append(rspBuf, buf[:n]...)
        rsp = string(rspBuf[:l])
        if strings.Contains(rsp, "check") {
            break
        }
    }
    l -= 6

    return string(rspBuf[:l])
}

