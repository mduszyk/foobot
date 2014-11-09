package conf

import (
    "os"
    "path/filepath"
	"fuzzywookie/foobot/log"
	"fuzzywookie/foobot/proto"
)

type ConfData map[string]string
var data = make(ConfData)

func Set(key string, value string) {
    data[key] = value
}

func Get(key string, args ...string) string {
    if val, ok := data[key]; ok {
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
    Set("bot.shell", "/bin/bash")

    Set("irc.nick", hostname)
}

func Dump() string {
    rsp := ""
    for k, v := range data {
        rsp += k + ": " + v + "\n"
    }

    return rsp
}

func NewConfModule() *ConfData {
    return &data
}

func (data *ConfData) Handle(msg *proto.Msg) string {
    rsp := ""
    if len(msg.Raw) == 0 {
        rsp = Dump()
    } else {
        rsp = msg.Raw + ": " + Get(msg.Raw)
    }
    return rsp
}
