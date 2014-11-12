package agent

import (
    "time"
    "reflect"
    "strconv"
	"github.com/mduszyk/foobot/log"
	"github.com/mduszyk/foobot/conf"
	"github.com/mduszyk/foobot/proto"
)

type Agent struct {
    protos map[string]proto.Proto
    modules map[string]proto.Interpreter
    proto proto.Proto
    auth *AuthModule
    authCmd string
    workers map[string]chan *proto.Msg
    cmdbuf int
    wrktimout time.Duration
}

func NewAgent() *Agent {
    buf, _ := strconv.Atoi(conf.Get("bot.cmdbuf", "10"))
    timout, _ := strconv.Atoi(conf.Get("bot.wrktimout", "10"))
	a := &Agent{
        proto: nil,
        auth: nil,
        authCmd: "",
        protos: make(map[string]proto.Proto),
        modules: make(map[string]proto.Interpreter),
        workers: make(map[string]chan *proto.Msg),
        cmdbuf: buf,
        wrktimout: time.Duration(timout) * time.Second,
    }
    return a
}

func (a *Agent) AddProto(name string, proto proto.Proto) {
    proto.Register(a)
    a.protos[name] = proto
    if a.proto == nil {
        a.proto = proto
    }
    log.INFO.Printf("Added proto, name: %s, type: %s", name, reflect.TypeOf(proto))
}

func (a *Agent) AddModule(cmd string, module proto.Interpreter) {
    a.modules[cmd] = module
    if auth, ok := module.(*AuthModule); ok {
        a.auth = auth
        a.authCmd = cmd
    }
    log.INFO.Printf("Added module, cmd: %s, type: %s", cmd, reflect.TypeOf(module))
}

func (a *Agent) Run() {
    log.INFO.Printf("Starting agent")
    // run default proto
    a.proto.Run()
}

func (a *Agent) StartProto(name string) {
    proto, ok := a.protos[name]
    if !ok {
        log.ERROR.Printf("Proto not found, name: %s", name)
        return
    }
    go proto.Run()
}

func (a *Agent) runWorker(input chan *proto.Msg, addr string) {
    for {
        select {
            case msg := <-input:
                module, ok := a.modules[msg.Cmd]
                if ok {
                    if a.authCmd != msg.Cmd {
                        // don't loging auth proto raw messages
                        log.TRACE.Printf("Agent, msg.Addr: %s, msg.Raw: %s",
                            msg.Addr, msg.Raw)
                    }
                    rsp := module.Handle(proto.Pop(msg))
                    msg.Proto.Send(msg.Addr, rsp)
                }
            case <-time.After(a.wrktimout):
                log.TRACE.Printf("Idle worker exiting, addr %s, timeout: %s",
                    addr, a.wrktimout)
                delete(a.workers, addr)
                return;
        }
    }
}

func (a *Agent) Handle(msg *proto.Msg) string {
    if a.auth != nil && a.authCmd != msg.Cmd && !a.auth.Verify(msg.User) {
        log.INFO.Printf("Forbidden, user: %s, cmd: %s", msg.User, msg.Raw)
        return ""
    }

    input, ok := a.workers[msg.Addr]
    if !ok {
        log.TRACE.Printf("Starting worker, addr: %s", msg.Addr)
        input = make(chan *proto.Msg, a.cmdbuf)
        a.workers[msg.Addr] = input
        go a.runWorker(input, msg.Addr)
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
