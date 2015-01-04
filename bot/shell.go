package bot

import(
    "strings"
    "github.com/mduszyk/foobot/log"
    "github.com/mduszyk/foobot/proto"
)


type ShellModule struct {
    shells map[string]*Shell
}

func NewShellModule() *ShellModule {
    return &ShellModule{
        shells: make(map[string]*Shell),
    }
}

func (m *ShellModule) CMD_(msg *proto.Msg) string {
    if len(msg.Raw) == 0 {
        return m.CMD__list(msg)
    }

    // run shell command
    sh, ok := m.shells[msg.Addr]
    if !ok {
        sh = NewShell()
        sh.Start()
        m.shells[msg.Addr] = sh
    }
    return sh.Insert(msg.Raw)
}

func (m *ShellModule) CMD__list(msg *proto.Msg) string {
    log.TRACE.Printf("Shell list, addr: %s", msg.Addr)
    rsp := ""
    for k, v := range m.shells {
        rsp += k + ": " + strings.Join(v.cmd.Args, " ") + "\n"
    }

    return rsp
}

func (m *ShellModule) CMD__kill(msg *proto.Msg) string {
    log.TRACE.Printf("Shell kill, addr: %s, args: %s", msg.Addr, msg.Args)
    sh, ok := m.shells[msg.Args]
    if ok {
        sh.Kill()
        delete(m.shells, msg.Args)
    }
    return ""
}

func (m *ShellModule) CMD__int(msg *proto.Msg) string {
    log.TRACE.Printf("Shell interrupt, addr: %s, args: %s", msg.Addr, msg.Args)
    sh, ok := m.shells[msg.Args]
    if ok {
        sh.Kill()
        delete(m.shells, msg.Args)
    }
    return ""
}

func (m *ShellModule) Handle(msg *proto.Msg) string {
    return CallCmdMethod(m, msg)
}

