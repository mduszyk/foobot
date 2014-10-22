package agent

import (
    "fmt"
)

type Agent struct {
    proto Proto

}

func NewAgent() *Agent {
	agent := &Agent{}
    return agent
}

func (agent *Agent) Recv(addr string, msg string) {
    fmt.Printf("Agent recv, addr: %s, msg: %s\n", addr, msg)
    agent.proto.Send(addr, "ACK: " + msg)
}

func (agent *Agent) AttachProto(proto Proto) {
    proto.Register(agent)
    agent.proto = proto
}
