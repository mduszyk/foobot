package agent

import(
    "io"
    "fmt"
    "os/exec"
    "strings"
)

type Shell struct {
    proc *exec.Cmd
    stdin io.WriteCloser
    stdout io.ReadCloser
}

func NewShell() *Shell {
    proc := exec.Command("/bin/bash")
    in, err := proc.StdinPipe()
    if err != nil {
        fmt.Printf("Error connecting shell stdin: %s\n", err)
        return nil
    }
    out, err := proc.StdoutPipe()
    if err != nil {
        fmt.Printf("Error connecting shell stdout: %s\n", err)
        return nil
    }
    proc.Start()

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

