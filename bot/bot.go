package bot

import (
    "time"
    "reflect"
    "strconv"
	"github.com/mduszyk/foobot/log"
	"github.com/mduszyk/foobot/conf"
	"github.com/mduszyk/foobot/proto"
)

type Bot struct {
    protos map[string]proto.Proto
    modules map[string]proto.Interpreter
    proto proto.Proto
    auth *AuthModule
    authCmd string
    workers map[string]chan *proto.Msg
    cmdbuf int
    wrktimout time.Duration
}

func NewBot() *Bot {
    buf, _ := strconv.Atoi(conf.Get("bot.cmdbuf", "10"))
    timout, _ := strconv.Atoi(conf.Get("bot.wrktimout", "10"))
	a := &Bot{
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

func (b *Bot) AddProto(name string, proto proto.Proto) {
    proto.Register(b)
    b.protos[name] = proto
    if b.proto == nil {
        b.proto = proto
    }
    log.INFO.Printf("Added proto, name: %s, type: %s", name, reflect.TypeOf(proto))
}

func (b *Bot) AddModule(cmd string, module proto.Interpreter) {
    b.modules[cmd] = module
    if auth, ok := module.(*AuthModule); ok {
        b.auth = auth
        b.authCmd = cmd
    }
    log.INFO.Printf("Added module, cmd: %s, type: %s", cmd, reflect.TypeOf(module))
}

func (b *Bot) Run() {
    log.INFO.Printf("Starting agent")
    // run default proto
    b.proto.Run()
}

func (b *Bot) StartProto(name string) {
    proto, ok := b.protos[name]
    if !ok {
        log.ERROR.Printf("Proto not found, name: %s", name)
        return
    }
    go proto.Run()
}

func (b *Bot) runWorker(input chan *proto.Msg, addr string) {
    for {
        select {
            case msg := <-input:
                module, ok := b.modules[msg.Cmd]
                if ok {
                    if b.authCmd != msg.Cmd {
                        // don't loging auth proto raw messages
                        log.TRACE.Printf("Bot, msg.Addr: %s, msg.Raw: %s",
                            msg.Addr, msg.Raw)
                    }
                    rsp := module.Handle(proto.Pop(msg))
                    msg.Proto.Send(msg.Addr, rsp)
                }
            case <-time.After(b.wrktimout):
                log.TRACE.Printf("Idle worker exiting, addr %s, timeout: %s",
                    addr, b.wrktimout)
                delete(b.workers, addr)
                return;
        }
    }
}

func (b *Bot) Handle(msg *proto.Msg) string {
    if b.auth != nil && b.authCmd != msg.Cmd && !b.auth.Verify(msg.User) {
        log.INFO.Printf("Forbidden, user: %s, cmd: %s", msg.User, msg.Raw)
        return ""
    }

    input, ok := b.workers[msg.Addr]
    if !ok {
        log.TRACE.Printf("Starting worker, addr: %s", msg.Addr)
        input = make(chan *proto.Msg, b.cmdbuf)
        b.workers[msg.Addr] = input
        go b.runWorker(input, msg.Addr)
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
