package protoimpl

import (
    "time"
    "strings"
    "strconv"
    "crypto/tls"
    irc "github.com/fluffle/goirc/client"
    "github.com/mduszyk/foobot/log"
    "github.com/mduszyk/foobot/conf"
    "github.com/mduszyk/foobot/proto"
    "github.com/mduszyk/foobot/module"
)

type IrcProto struct {
    conf *irc.Config
    conn *irc.Conn
    disconn chan bool
    terminate bool
}

func NewIrcProto() *IrcProto {
    nick := conf.Get("irc.nick")
    ssl, err := strconv.ParseBool(conf.Get("irc.ssl", "false"))
    if err != nil {
        log.ERROR.Printf("Incorrect value, irc.ssl: %s",
            conf.Get("irc.ssl", "false"))
        ssl = false
    }
    noverify, err := strconv.ParseBool(conf.Get("irc.ssl.noverify", "false"))
    if err != nil {
        log.ERROR.Printf("Incorrect value, irc.ssl.noverify: %s",
            conf.Get("irc.ssl.noverify", "false"))
        noverify = false
    }

    cfg := irc.NewConfig(nick, conf.Get("irc.ident", nick),
        conf.Get("irc.name", nick))
    cfg.Version = conf.Get("irc.version", nick)
    cfg.QuitMessage = conf.Get("irc.quitmsg", "bye")
    cfg.Pass = conf.Get("irc.pass", "")
    cfg.SSL = ssl
    cfg.SSLConfig = &tls.Config{}
    cfg.SSLConfig.InsecureSkipVerify = noverify
    cfg.Server = conf.Get("irc.server")
    cfg.NewNick = func(n string) string { return n + "^" }
    cfg.Flood = true

    // create new IRC connection
    c := irc.Client(cfg)
    c.EnableStateTracking()

    proto := &IrcProto{
        conf: cfg,
        conn: c,
        disconn: make(chan bool),
        terminate: false,
    }

    c.HandleFunc("connected", func(conn *irc.Conn, line *irc.Line) {
        log.INFO.Printf("Connected, server: %s, msg: %s", cfg.Server, line)
        conn.Join(conf.Get("irc.channel", "#foobot"))
    })

    c.HandleFunc("disconnected", func(conn *irc.Conn, line *irc.Line) {
        log.INFO.Printf("Disconnected, server: %s, msg: %s", cfg.Server, line)
        proto.disconn <- true
    })

    return proto
}

func (p *IrcProto) Run() {
    log.INFO.Printf("Starting irc proto")
    for !p.terminate {
        // connect to server
        if err := p.conn.Connect(); err != nil {
            log.ERROR.Printf("Connection error: %s", err)
        } else {
            // wait on disconnect channel
            <-p.disconn
        }
        time.Sleep(time.Second)
    }
}

func (p *IrcProto) Send(addr string, text string) {
    for _, e := range strings.Split(text, "\n") {
        p.conn.Privmsg(addr, e)
    }
}

func (p *IrcProto) Register(i proto.Interpreter) {
    handler := func(conn *irc.Conn, line *irc.Line) {
        text := line.Text()
        // ignore irc chat messages
        if text[0] != ':' {
            return
        }
        addr := line.Target()
        /* log.TRACE.Printf("Got message, addr: %s, irc line: %s", addr, line) */
        msg := proto.Parse(text)
        msg.Addr = addr
        msg.User = line.Src
        msg.Proto = p
        // pass message to agent
        rsp := i.Handle(msg)
        if rsp != "" {
            p.Send(addr, rsp)
        }
    }
    p.conn.HandleFunc("PRIVMSG", handler)
}

func (p *IrcProto) CMD_msg(msg *proto.Msg) string {
    p.conn.Privmsg(msg.Arg[0], strings.Join(msg.Arg[1:], " "))
    return ""
}

func (p *IrcProto) CMD_join(msg *proto.Msg) string {
    p.conn.Join(msg.Arg[0])
    return ""
}

func (p *IrcProto) CMD_part(msg *proto.Msg) string {
    p.conn.Part(msg.Arg[0], "bye")
    return ""
}

func (p *IrcProto) CMD_kick(msg *proto.Msg) string {
    p.conn.Kick(msg.Arg[0], msg.Arg[1], "bye")
    return ""
}

func (p *IrcProto) CMD_nick(msg *proto.Msg) string {
    p.conn.Nick(msg.Arg[0])
    return ""
}

func (p *IrcProto) Handle(msg *proto.Msg) string {
    return module.CallCmdMethod(p, msg)
}
