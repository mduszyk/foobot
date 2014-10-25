package proto

import (
    "fmt"
    "strings"
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

func (proto *IrcProto) Send(addr string, text string) {
    for _, e := range strings.Split(text, "\n") {
        proto.conn.Privmsg(addr, e)
    }
}

func (proto *IrcProto) Register(r agent.Receiver) {
    handler := func(conn *irc.Conn, line *irc.Line) {
        addr := line.Target()
        text := line.Text()
        chunks := strings.SplitN(text, " ", 2)
        msg := &agent.Msg{
            Raw: text,
            Cmd: chunks[0],
            Args: "",
        }
        if len(chunks) > 1 {
            msg.Args = chunks[1]
        }
        switch {
            default:
                // pass message to agent
                r.Recv(addr, msg)
            case strings.HasPrefix(msg.Cmd, ":irc"):
                // handle message here
                fmt.Printf("Got irc proto command: %s", msg.Raw)
        }
    }
	proto.conn.HandleFunc("PRIVMSG", handler)
}

