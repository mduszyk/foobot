package conf

import (
    "os"
    "path/filepath"
    "github.com/mduszyk/foobot/log"
    "github.com/mduszyk/foobot/proto"
)

type ConfData map[string]string

var instance = make(ConfData)

func Set(key string, value string) {
    instance[key] = value
}

func Get(key string, args ...string) string {
    if val, ok := instance[key]; ok {
        return val
    }
    return args[0]
}

func GetBinDir() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.ERROR.Printf("Failed getting binary dir, error: %s", err)
    }

    return dir
}

func Init() {
    hostname, err := os.Hostname()
    if err != nil {
        log.ERROR.Printf("Failed getting hostname, error: %s", err)
    }

    Set("bot.bindir", GetBinDir())
    Set("irc.nick", hostname)
}

func Dump() string {
    rsp := ""
    for k, v := range instance {
        rsp += k + ": " + v + "\n"
    }

    return rsp
}

func GetConfModule() *ConfData {
    return &instance
}

func (cd *ConfData) Handle(msg *proto.Msg) string {
    rsp := ""

    switch msg.Cmd {
        case "":
            rsp = Dump()
        case "get":
            rsp = msg.Args + ": " + Get(msg.Args)
        case "set":
            Set(msg.Arg[0], msg.Arg[1])
    }

    return rsp
}
