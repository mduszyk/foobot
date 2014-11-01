package agent

import (
    "time"
    "reflect"
    "strconv"
	"fuzzywookie/foobot/log"
	"fuzzywookie/foobot/conf"
	"fuzzywookie/foobot/proto"
)

type Agent struct {
    protos map[string]proto.Proto
    modules map[string]proto.Interpreter
    proto proto.Proto
    workers map[string]chan *proto.Msg
    cmdbuf int
    wrktimout time.Duration
}

func NewAgent() *Agent {
    buf, _ := strconv.Atoi(conf.Get("bot.cmdbuf", "10"))
    timout, _ := strconv.Atoi(conf.Get("bot.wrktimout", "10"))
	agent := &Agent{
        proto: nil,
        protos: make(map[string]proto.Proto),
        modules: make(map[string]proto.Interpreter),
        workers: make(map[string]chan *proto.Msg),
        cmdbuf: buf,
        wrktimout: time.Duration(timout) * time.Second,
    }
    return agent
}

func (agent *Agent) AddProto(name string, proto proto.Proto) {
    proto.Register(agent)
    agent.protos[name] = proto
    if agent.proto == nil {
        agent.proto = proto
    }
    log.INFO.Printf("Added proto, name: %s, type: %s", name, reflect.TypeOf(proto))
}

func (agent *Agent) AddModule(cmd string, module proto.Interpreter) {
    agent.modules[cmd] = module
    log.INFO.Printf("Added module, cmd: %s, type: %s", cmd, reflect.TypeOf(module))
}

func (agent *Agent) Run() {
    log.INFO.Printf("Starting agent")
    // run default proto
    agent.proto.Run()
}

func (agent *Agent) StartProto(name string) {
    proto, ok := agent.protos[name]
    if !ok {
        log.ERROR.Printf("Proto not found, name: %s", name)
        return
    }
    go proto.Run()
}

func (agent *Agent) runWorker(input chan *proto.Msg, addr string) {
    for {
        select {
            case msg := <-input:
                module, ok := agent.modules[msg.Cmd]
                if ok {
                    log.TRACE.Printf("Agent, msg.Addr: %s, msg.Raw: %s",
                        msg.Addr, msg.Raw)
                    rsp := module.Handle(msg)
                    msg.Src.Send(msg.Addr, rsp)
                }
            case <-time.After(agent.wrktimout):
                log.TRACE.Printf("Idle worker exiting, addr %s, timeout: %s",
                    addr, agent.wrktimout)
                delete(agent.workers, addr)
                return;
        }
    }
}

func (agent *Agent) Handle(msg *proto.Msg) string {

    input, ok := agent.workers[msg.Addr]
    if !ok {
        log.TRACE.Printf("Starting worker, addr: %s", msg.Addr)
        input = make(chan *proto.Msg, agent.cmdbuf)
        agent.workers[msg.Addr] = input
        go agent.runWorker(input, msg.Addr)
    }

    rsp := ""

    select {
        case input <- msg:
            log.TRACE.Printf("Msg sent to worker, addr: %s", msg.Addr)
        default:
            log.ERROR.Printf("Cmd buffer full, addr: %s", msg.Addr)
    }

    return rsp
}
