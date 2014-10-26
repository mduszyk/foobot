package proto

import(
    "strings"
)

type Msg struct {
    Raw string
    Cmd string
    Args string
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
    }
    if len(chunks) > 1 {
        msg.Args = chunks[1]
    }
    return msg
}