package main

import (
	"fuzzywookie/foobot/agent"
	"fuzzywookie/foobot/proto"
	/* "fuzzywookie/foobot/test" */
)


func main() {
    irc := proto.NewIrcProto()
    a := agent.NewAgent()
    a.Attach(irc)
    a.Run()

    /* test.RunIrcBot() */
}
