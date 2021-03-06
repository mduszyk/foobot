package proto

import(
    "strings"
)

type Msg struct {
    Raw string
    Cmd string
    Args string
    Arg []string
    Addr string
    User string
    Proto Proto
}

type Interpreter interface {
    Handle(msg *Msg) string
}

type Proto interface {
    Send(addr string, text string)
    Register(i Interpreter)
    Run()
}

func Parse(text string) *Msg {
    chunks := strings.SplitN(text, " ", 2)
    msg := &Msg{
        Raw: text,
        Cmd: chunks[0],
        Args: "",
        Arg: nil,
    }
    if len(chunks) > 1 {
        msg.Args = chunks[1]
        msg.Arg = strings.Split(chunks[1], " ")
    }
    return msg
}

func Pop(m *Msg) *Msg {
    msg := Parse(m.Args)
    msg.Raw = m.Args
    msg.Addr = m.Addr
    msg.User = m.User
    msg.Proto = m.Proto
    return msg
}
