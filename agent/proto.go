package agent

type Receiver interface {
    Recv(addr string, msg string)
}

type Proto interface {
    Send(addr string, msg string)
    Register(r Receiver)
    Run()
}

