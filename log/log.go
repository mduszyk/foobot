// TODO make buffer to be cyclic buffer with fixed size
package log

import(
    "os"
    "io"
    "log"
    "bytes"
    "strconv"
    "strings"
	"github.com/mduszyk/foobot/proto"
)

const LEVEL_TRACE = 0
const LEVEL_INFO = 1
const LEVEL_WARN = 2
const LEVEL_ERROR = 3

type MutableWriter struct {
    writer io.Writer
}

func (mw *MutableWriter) Write(p []byte) (n int, err error) {
    return mw.writer.Write(p)
}

func (mw *MutableWriter) SetWriter(w io.Writer) {
    mw.writer = w
}

type LogData struct {
    writer *MutableWriter
    level int
}

var buf bytes.Buffer

var data = LogData {
    writer: &MutableWriter{&buf},
    level: LEVEL_INFO,
}

var levelMap = map[string]int{
    "trace": LEVEL_TRACE,
    "info": LEVEL_INFO,
    "warn": LEVEL_WARN,
    "error": LEVEL_ERROR,
}

type Log interface {
    Printf(format string, v ...interface{})
}

type nullLog struct{}
func (l *nullLog) Printf(format string, v ...interface{}) {}

var disabled nullLog
var trace = log.New(data.writer, "TRACE: ", log.Ldate | log.Ltime | log.Lshortfile)
var info = log.New(data.writer, "INFO : ", log.Ldate | log.Ltime | log.Lshortfile)
var warn = log.New(data.writer, "WARN : ", log.Ldate | log.Ltime | log.Lshortfile)
var err = log.New(data.writer, "ERROR: ", log.Ldate | log.Ltime | log.Lshortfile)

var TRACE = Log(&disabled)
var INFO = Log(info)
var WARN = Log(warn)
var ERROR = Log(err)

func EnableStderr() {
    data.writer.SetWriter(io.MultiWriter(&buf, os.Stderr))
}

func SetLevelStr(l string) {
    SetLevel(levelMap[strings.ToLower(l)])
}

func SetLevel(l int) {
    data.level = l
    if (LEVEL_TRACE >= data.level) {
        TRACE = trace
    } else {
        TRACE = &disabled
    }
    if (LEVEL_INFO >= data.level) {
        INFO = info
    } else {
        INFO = &disabled
    }
    if (LEVEL_WARN >= data.level) {
        WARN = warn
    } else {
        WARN = &disabled
    }
    if (LEVEL_ERROR >= data.level) {
        ERROR = err
    } else {
        ERROR = &disabled
    }
}

func Tail(n int) string {
    // TODO 
    return buf.String()
}

func NewLogModule() *LogData {
    return &data
}

func (data *LogData) Handle(msg *proto.Msg) string {
    rsp := ""

    switch msg.Cmd {
        case "tail":
            n, err := strconv.Atoi(msg.Args)
            if err == nil {
                n = 1
            }
            rsp = Tail(n)
        case "level":
            SetLevelStr(msg.Args)
            rsp = "log level " + msg.Args
    }

    return rsp
}

