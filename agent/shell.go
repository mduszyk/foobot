package agent

import(
    "strings"
	"fuzzywookie/foobot/log"
	"fuzzywookie/foobot/proto"
)


type ShellModule struct {
    shells map[string]*Shell
}

func NewShellModule() *ShellModule {
    return &ShellModule{
        shells: make(map[string]*Shell),
    }
}

func (m *ShellModule) list() string {
    rsp := ""
    for k, v := range m.shells {
        rsp += k + ": " + strings.Join(v.proc.Args, " ") + "\n"
    }

    return rsp
}

func (m *ShellModule) Handle(msg *proto.Msg) string {
    rsp := ""

    msg2 := proto.Parse(msg.Args)
    switch msg2.Cmd {
        case "":
            log.TRACE.Printf("Shell list, addr: %s")
            rsp = m.list()
        case ":list":
            log.TRACE.Printf("Shell list, addr: %s")
            rsp = m.list()
        case ":kill":
            log.TRACE.Printf("Shell kill, addr: %s, args: %s", msg2.Args)
        default:
            sh, ok := m.shells[msg.Addr]
            if !ok {
                sh = NewShell()
                sh.Start()
                m.shells[msg.Addr] = sh
            }
            rsp = sh.Insert(msg.Args)
    }

    return rsp
}

