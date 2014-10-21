package test

import (
	"bufio"
	"flag"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"os"
	"strings"
    "crypto/tls"
)

var host *string = flag.String("host", "cube.mdevel.net:6697", "IRC server")
var channel *string = flag.String("channel", "#bot", "IRC channel")
var pass *string = flag.String("pass", "baltycka", "Password")

func RunIrcBot() {
	flag.Parse()

	// create new IRC connection
    cfg := irc.NewConfig("bot1", "foobot", "foobot")
	cfg.Version = "foobot 1.0"
	cfg.QuitMessage = "bye"
    cfg.Pass = *pass
    cfg.SSL = true
    cfg.SSLConfig = &tls.Config{}
    cfg.SSLConfig.InsecureSkipVerify = true
    cfg.Server = *host
    cfg.NewNick = func(n string) string { return n + "^" }

    c := irc.Client(cfg)
	c.EnableStateTracking()
	c.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) { conn.Join(*channel) })

	// Set up a handler to notify of disconnect events.
	quit := make(chan bool)
	c.HandleFunc("disconnected",
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	// set up a goroutine to read commands from stdin
	in := make(chan string, 4)
	reallyquit := false
	go func() {
		con := bufio.NewReader(os.Stdin)
		for {
			s, err := con.ReadString('\n')
			if err != nil {
				// wha?, maybe ctrl-D...
				close(in)
				break
			}
			// no point in sending empty lines down the channel
			if len(s) > 2 {
				in <- s[0 : len(s)-1]
			}
		}
	}()

	// set up a goroutine to do parsey things with the stuff from stdin
	go func() {
		for cmd := range in {
			fmt.Printf("cmd: %s\n", cmd)
			if cmd[0] == ':' {
				switch idx := strings.Index(cmd, " "); {
				case cmd[1] == 'd':
					fmt.Printf(c.String())
				case cmd[1] == 'f':
					if len(cmd) > 2 && cmd[2] == 'e' {
						// enable flooding
						c.Config().Flood = true
					} else if len(cmd) > 2 && cmd[2] == 'd' {
						// disable flooding
						c.Config().Flood = false
					}
					for i := 0; i < 20; i++ {
						c.Privmsg(*channel, "flood test!")
					}
				case idx == -1:
					continue
				case cmd[1] == 'q':
					reallyquit = true
					c.Quit(cmd[idx+1 : len(cmd)])
				case cmd[1] == 'j':
					c.Join(cmd[idx+1 : len(cmd)])
				case cmd[1] == 'p':
					c.Part(cmd[idx+1 : len(cmd)])
				}
			} else {
				c.Raw(cmd)
			}
		}
	}()

	for !reallyquit {
		// connect to server
		if err := c.ConnectTo(*host); err != nil {
			fmt.Printf("Connection error: %s\n", err)
			return
		}

		// wait on quit channel
		<-quit
	}
}
