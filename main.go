package main

import (
	"fuzzywookie/foobot/agent"
	"fuzzywookie/foobot/proto"
	/* "fuzzywookie/foobot/test" */
)


func main() {
    irc := proto.NewIrcProto()
    a := agent.NewAgent()
    a.AttachProto(irc)
    irc.Run()
    /* test.RunIrcBot() */
}
