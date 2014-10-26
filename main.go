package main

import (
	"flag"
	"fuzzywookie/foobot/log"
	"fuzzywookie/foobot/conf"
	"fuzzywookie/foobot/agent"
	"fuzzywookie/foobot/protoimpl"
)

var verbose *bool = flag.Bool("v", false, "Prints logs to stdout on trace level")

func main() {
	flag.Parse()

    if *verbose {
        log.EnableStdout()
        log.SetLevel(log.LEVEL_TRACE)
    }

    conf.Init()
    conf.Set("irc.ident", "foobot")
    conf.Set("irc.name", "foobot")
    conf.Set("irc.version", "foobot 1.0")
    conf.Set("irc.pass", "baltycka")
    conf.Set("irc.server", "cube.mdevel.net:6697")

    ircProto := protoimpl.NewIrcProto()

    a := agent.NewAgent()
    a.AddModule(":log", log.NewLogModule())
    a.AddModule(":conf", conf.NewConfModule())
    a.AddModule(":sh", agent.NewShellModule())
    a.AddModule(":irc", ircProto)
    a.AddProto("irc", ircProto)
    a.Run()

}
