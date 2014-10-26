package agent

import (
    "strconv"
    "strings"
    "reflect"
	"fuzzywookie/foobot/log"
	"fuzzywookie/foobot/conf"
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

    switch msg.Cmd {
        default:
            rsp = "ECHO: " + msg.Raw
        case ":sh":
            rsp = agent.sh.Insert(msg.Args)
        case ":conf":
            if len(msg.Args) == 0 {
                rsp = conf.Dump()
            } else {
                rsp = msg.Args + ": " + conf.Get(msg.Args)
            }
        case ":log":
            if strings.HasPrefix(msg.Args, "level") {
                chunks := strings.SplitN(msg.Args, " ", 2)
                log.SetLevelStr(chunks[1])
                rsp = "log level " + chunks[1]
            } else {
                n, err := strconv.Atoi(msg.Args)
                if err == nil {
                    n = 1
                }
                rsp = log.Tail(n)
            }
    }

    log.TRACE.Printf("Agent cmd, addr: %s, msg: %s", addr, msg.Raw)
    agent.proto.Send(addr, rsp)
}

func (agent *Agent) Attach(proto Proto) {
    proto.Register(agent)
    agent.proto = proto
    log.INFO.Printf("Attached proto: %s", reflect.TypeOf(proto))
}

func (agent *Agent) Run() {
    log.INFO.Printf("Starting agent")
    agent.proto.Run()
}
