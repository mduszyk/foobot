package bot

import(
    "os"
    "runtime"
    "strconv"
    "strings"
    "github.com/mduszyk/foobot/log"
    "github.com/mduszyk/foobot/conf"
    "github.com/mduszyk/foobot/proto"
    "github.com/mduszyk/foobot/module"
)

type Info struct {
}

func NewInfoModule() *Info {
    return &Info{}
}

func (i *Info) CMD_(msg *proto.Msg) string {
    hostname, err := os.Hostname()
    if err != nil {
        log.ERROR.Printf("Failed getting hostname, error: %s", err)
    }
    wd, err := os.Getwd()
    if err != nil {
        log.ERROR.Printf("Failed getting wd, error: %s", err)
    }
    var mem runtime.MemStats
    runtime.ReadMemStats(&mem)
    info := ""
    info += "os.hostname: " + hostname + "\n"
    info += "os.pid: " + strconv.Itoa(os.Getpid()) + "\n"
    info += "os.gid: " + strconv.Itoa(os.Getgid()) + "\n"
    info += "os.uid: " + strconv.Itoa(os.Getuid()) + "\n"
    info += "os.page: " + strconv.Itoa(os.Getpagesize()) + "\n"
    info += "os.wd: " + wd + "\n"
    info += "bot.bindir: " + conf.GetBinDir() + "\n"
    info += "runtime.GOOS: " + runtime.GOOS + "\n"
    info += "runtime.GOARCH: " + runtime.GOARCH + "\n"
    info += "runtime.GOMAXPROCS: " + strconv.Itoa(runtime.GOMAXPROCS(0)) + "\n"
    info += "runtime.NumCPU: " + strconv.Itoa(runtime.NumCPU()) + "\n"
    info += "runtime.NumGoroutine: " + strconv.Itoa(runtime.NumGoroutine()) + "\n"
    info += "runtime.Version: " + runtime.Version() + "\n"
    info += "runtime.MemStats.Alloc: " + strconv.FormatUint(mem.Alloc, 10) + "\n"
    info += "runtime.MemStats.Sys: " + strconv.FormatUint(mem.Sys, 10) + "\n"

    return info
}

func (i *Info) CMD_env(msg *proto.Msg) string {
    var info string
    key := msg.Args
    if len(key) == 0 {
        info = strings.Join(os.Environ(), "\n")
    } else {
        info = key + "=" + os.Getenv(key)
    }
    return info
}

func (i *Info) CMD_help(msg *proto.Msg) string {
    text := "Info module returns information about environment\n" 
    text += ":info\n"
    text += ":info env\n"
    text += ":info env PATH\n"
    return text
}

func (i *Info) Handle(msg *proto.Msg) string {
    return module.CallCmdMethod(i, msg)
}
