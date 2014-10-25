package agent

import (
    "fmt"
)

type Agent struct {
    proto Proto
    sh *Shell
}

func NewAgent() *Agent {
    sh := NewShell()
	agent := &Agent{
        sh: sh,
    }
    return agent
}

func (agent *Agent) Recv(addr string, msg *Msg) {
    var rsp string

    switch {
        default:
            rsp = "ECHO: " + msg.Raw
        case ":sh" == msg.Cmd:
            rsp = agent.sh.Insert(msg.Args)
    }

    fmt.Printf("Agent addr: %s, msg: %s, rsp: %s\n", addr, msg, rsp)
    agent.proto.Send(addr, rsp)
}

func (agent *Agent) Attach(proto Proto) {
    proto.Register(agent)
    agent.proto = proto
}

func (agent *Agent) Run() {
    agent.proto.Run()
}
