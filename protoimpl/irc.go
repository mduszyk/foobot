package protoimpl

import (
    "strings"
    "crypto/tls"
	"fuzzywookie/foobot/log"
	"fuzzywookie/foobot/conf"
	"fuzzywookie/foobot/proto"
	irc "github.com/fluffle/goirc/client"
)

type IrcProto struct {
    conf *irc.Config
    conn *irc.Conn
    disconn chan bool
    terminate bool
}

func NewIrcProto() *IrcProto {
    nick := conf.Get("irc.nick")
    cfg := irc.NewConfig(nick, conf.Get("irc.ident", nick), conf.Get("irc.name", nick))
	cfg.Version = conf.Get("irc.version", nick)
	cfg.QuitMessage = conf.Get("irc.quitmsg", "bye")
    cfg.Pass = conf.Get("irc.pass", "")
    cfg.SSL = true
    cfg.SSLConfig = &tls.Config{}
    cfg.SSLConfig.InsecureSkipVerify = true
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
        log.INFO.Printf("Connected to irc server")
        conn.Join(conf.Get("irc.channel", "#foobot"))
    })

	c.HandleFunc("disconnected", func(conn *irc.Conn, line *irc.Line) {
        log.INFO.Printf("Disconnected from irc server")
        proto.disconn <- true
    })

    return proto
}

func (proto *IrcProto) Run() {
	for !proto.terminate {
		// connect to server
		if err := proto.conn.Connect(); err != nil {
            log.ERROR.Printf("Connection error: %s", err)
			return
		}

		// wait on disconnect channel
		<-proto.disconn
	}
}

func (proto *IrcProto) Send(addr string, text string) {
    for _, e := range strings.Split(text, "\n") {
        proto.conn.Privmsg(addr, e)
    }
}

func (p *IrcProto) Register(i proto.Interpreter) {
    handler := func(conn *irc.Conn, line *irc.Line) {
        log.TRACE.Printf("Got message, line: %s", line)
        msg := proto.Parse(line.Text())
        // pass message to agent
        rsp := i.Handle(msg)
        p.Send(line.Target(), rsp)
    }
	p.conn.HandleFunc("PRIVMSG", handler)
}

func (p *IrcProto) Handle(msg *proto.Msg) string {
    msg2 := proto.Parse(msg.Args)
    switch msg2.Cmd {
        case "msg":
            p.conn.Privmsg(msg2.Arg[0], strings.Join(msg2.Arg[1:], " "))
        case "join":
            p.conn.Join(msg2.Arg[0])
        case "part":
            p.conn.Part(msg2.Arg[0], "bye")
        case "kick":
            p.conn.Kick(msg2.Arg[0], msg2.Arg[1], "bye")
        case "nick":
            p.conn.Nick(msg2.Arg[0])
    }

    return ""
}
