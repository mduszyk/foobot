package main

import (
    "flag"
    "fuzzywookie/foobot/log"
    "fuzzywookie/foobot/conf"
    "fuzzywookie/foobot/agent"
    "fuzzywookie/foobot/protoimpl"
    "github.com/VividCortex/godaemon"
)

var verbose *bool = flag.Bool("v", false, "Prints logs to stdout on trace level")

func main() {
	flag.Parse()

    if *verbose {
        log.EnableStdout()
        log.SetLevel(log.LEVEL_TRACE)
    } else {
        godaemon.MakeDaemon(&godaemon.DaemonAttr{})
    }

    conf.Init()
    conf.Set("irc.channel", "#foobot")
    conf.Set("irc.ident", "foobot")
    conf.Set("irc.name", "foobot")
    conf.Set("irc.pass", "baltycka")
    conf.Set("irc.server", "cube.mdevel.net:6697")
    conf.Set("irc.version", "foobot 1.0")
    conf.Set("net.server.type", "tcp")
    conf.Set("net.server.socket", "localhost:6600")
    conf.Set("bot.cmdbuf", "10")
    conf.Set("bot.wrktimout", "120")

    ircProto := protoimpl.NewIrcProto()
    netServerProto := protoimpl.NewNetServerProto()

    a := agent.NewAgent()

    a.AddModule(":conf", conf.NewConfModule())
    a.AddModule(":irc", ircProto)
    a.AddModule(":info", agent.NewInfoModule())
    a.AddModule(":log", log.NewLogModule())
    a.AddModule(":sh", agent.NewShellModule())

    a.AddProto("irc", ircProto)
    a.AddProto("net", netServerProto)

    a.StartProto("net")

    a.Run()
}
