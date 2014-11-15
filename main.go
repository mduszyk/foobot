package main

import (
    "flag"
    "github.com/VividCortex/godaemon"
    "github.com/mduszyk/foobot/log"
    "github.com/mduszyk/foobot/conf"
    "github.com/mduszyk/foobot/bot"
    "github.com/mduszyk/foobot/protoimpl"
)

var verbose *bool = flag.Bool("v", false, "Prints logs to stderr on trace level")
var pass *string = flag.String("P", "", "Set custom bot pass")
var ircServer *string = flag.String("s", "example.com:6697", "irc server socket")
var ircPass *string = flag.String("p", "", "irc server password")

func main() {
    flag.Parse()

    if *verbose {
        log.EnableStderr()
        log.SetLevel(log.LEVEL_TRACE)
    } else {
        godaemon.MakeDaemon(&godaemon.DaemonAttr{})
    }

    conf.Init()
    conf.Set("irc.channel", "#bot")
    conf.Set("irc.ident", "foobot")
    conf.Set("irc.name", "foobot")
    conf.Set("irc.pass", *ircPass)
    conf.Set("irc.server", *ircServer)
    conf.Set("irc.version", "foobot 1.0")
    conf.Set("net.server.type", "tcp")
    conf.Set("net.server.socket", "localhost:6600")
    conf.Set("bot.cmdbuf", "10")
    conf.Set("bot.wrktimout", "120")
    conf.Set("bot.pass", *pass)

    ircProto := protoimpl.NewIrcProto()
    netServerProto := protoimpl.NewNetServerProto()

    a := bot.NewBot()

    a.AddModule(":conf", conf.NewConfModule())
    a.AddModule(":irc", ircProto)
    a.AddModule(":info", bot.NewInfoModule())
    a.AddModule(":log", log.NewLogModule())
    a.AddModule(":sh", bot.NewShellModule())
    a.AddModule(":auth", bot.NewAuthModule())

    a.AddProto("irc", ircProto)
    a.AddProto("net", netServerProto)

    a.StartProto("net")

    a.Run()
}
