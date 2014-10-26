package conf

import (
    "os"
    "strconv"
    "strings"
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

func getBinDir() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.ERROR.Printf("Failed getting binary dir, error: %s", err)
    }

    return dir
}

func Init() {
    hostname, err := os.Hostname()
    if err != nil {
        log.ERROR.Printf("Failed getting hostname, error: %s")
    }
    wd, err := os.Getwd()
    if err != nil {
        log.ERROR.Printf("Failed getting wd, error: %s")
    }

    Set("os.wd", wd)
    Set("os.hostname", hostname)
    Set("os.pid", strconv.Itoa(os.Getpid()))
    Set("os.uid", strconv.Itoa(os.Getuid()))

    Set("bot.bindir", getBinDir())
    Set("bot.shell", "/bin/bash")

    Set("irc.nick", hostname)
}

func Dump() string {
    buf := make([]string, 64)
    for k, v := range data {
        buf = append(buf, k + ": " + v)
    }

    return strings.Join(buf, "\n")
}

func NewConfModule() *ConfData {
    return &data
}

func (data *ConfData) Handle(msg *proto.Msg) string {
    rsp := ""
    if len(msg.Args) == 0 {
        rsp = Dump()
    } else {
        rsp = msg.Args + ": " + Get(msg.Args)
    }
    return rsp
}