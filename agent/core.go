package agent

import (
    "fmt"
    "strings"
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

func (agent *Agent) Recv(addr string, msg string) {
    var rsp string
    chunks := strings.SplitN(msg, " ", 2)

    if chunks[0] == ":sh" {
        rsp = agent.sh.Insert(chunks[1])
    } else {
        rsp = "ACK: " + msg
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
