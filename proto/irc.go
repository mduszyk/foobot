package proto

import (
    "fmt"
    "crypto/tls"
	"fuzzywookie/foobot/agent"
	irc "github.com/fluffle/goirc/client"
)

type IrcProto struct {
    conf *irc.Config
    conn *irc.Conn
    disconn chan bool
    terminate bool
}

func NewIrcProto() *IrcProto {
    // TODO move params to agent.conf
    cfg := irc.NewConfig("bot1", "foobot", "foobot")
	cfg.Version = "foobot 1.0"
	cfg.QuitMessage = "bye"
    cfg.Pass = "baltycka"
    cfg.SSL = true
    cfg.SSLConfig = &tls.Config{}
    cfg.SSLConfig.InsecureSkipVerify = true
    cfg.Server = "cube.mdevel.net:6697"
    cfg.NewNick = func(n string) string { return n + "^" }

	// create new IRC connection
    c := irc.Client(cfg)
	c.EnableStateTracking()

	proto := &IrcProto{
		conf: cfg,
		conn: c,
        disconn: make(chan bool),
        terminate: false,
	}

	c.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) { conn.Join("#foobot") })

	c.HandleFunc("disconnected",
		func(conn *irc.Conn, line *irc.Line) { proto.disconn <- true })

    return proto
}

func (proto *IrcProto) Run() {
	for !proto.terminate {
		// connect to server
		if err := proto.conn.Connect(); err != nil {
			fmt.Printf("Connection error: %s\n", err)
			return
		}

		// wait on disconnect channel
		<-proto.disconn
	}
}

func (proto *IrcProto) Send(addr string, msg string) {
    proto.conn.Privmsg(addr, msg)
}

func (proto *IrcProto) Register(r agent.Receiver) {
    handler := func(conn *irc.Conn, line *irc.Line) {
        r.Recv(line.Target(), line.Text())
    }
	proto.conn.HandleFunc("PRIVMSG", handler)
}

