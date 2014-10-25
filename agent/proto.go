package agent

type Msg struct {
    Raw string
    Cmd string
    Args string
}

type Receiver interface {
    Recv(addr string, msg *Msg)
}

type Proto interface {
    Send(addr string, text string)
    Register(r Receiver)
    Run()
}

