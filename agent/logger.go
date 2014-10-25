// TODO make buffer to be cyclic buffer with fixed size
package agent

import "log"
import "bytes"
import "strings"

const TRACE = 0
const INFO = 1
const WARN = 2
const ERROR = 3

var buf bytes.Buffer
var level int

var levelMap = map[string]int{
    "trace": TRACE,
    "info": INFO,
    "warn": WARN,
    "ERROR": ERROR,
}

type Log interface {
    Printf(format string, v ...interface{})
}

type logNull struct{}
func (l *logNull) Printf(format string, v ...interface{}) {}

var disabled logNull
var trace = log.New(&buf, "TRACE: ", log.Ldate | log.Ltime | log.Lshortfile)
var info = log.New(&buf, "INFO : ", log.Ldate | log.Ltime | log.Lshortfile)
var warn = log.New(&buf, "WARN : ", log.Ldate | log.Ltime | log.Lshortfile)
var err = log.New(&buf, "ERROR: ", log.Ldate | log.Ltime | log.Lshortfile)

var LogTrace = Log(&disabled)
var LogInfo = Log(info)
var LogWarn = Log(warn)
var LogError = Log(err)

func LogLevelStr(l string) {
    LogLevel(levelMap[strings.ToLower(l)])
}

func LogLevel(l int) {
    level = l
    if (TRACE >= level) {
        LogTrace = trace
    } else {
        LogTrace = &disabled
    }
    if (INFO >= level) {
        LogInfo = info
    } else {
        LogInfo = &disabled
    }
    if (WARN >= level) {
        LogWarn = warn
    } else {
        LogWarn = &disabled
    }
    if (ERROR >= level) {
        LogError = err
    } else {
        LogError = &disabled
    }
}

func LogTail(n int) string {
    // TODO 
    return buf.String()
}

