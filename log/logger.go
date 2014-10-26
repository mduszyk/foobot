// TODO make buffer to be cyclic buffer with fixed size
package log

import "log"
import "bytes"
import "strings"

const LEVEL_TRACE = 0
const LEVEL_INFO = 1
const LEVEL_WARN = 2
const LEVEL_ERROR = 3

var buf bytes.Buffer
var level int

var levelMap = map[string]int{
    "trace": LEVEL_TRACE,
    "info": LEVEL_INFO,
    "warn": LEVEL_WARN,
    "ERROR": LEVEL_ERROR,
}

type Log interface {
    Printf(format string, v ...interface{})
}

type nullLog struct{}
func (l *nullLog) Printf(format string, v ...interface{}) {}

var disabled nullLog
var trace = log.New(&buf, "TRACE: ", log.Ldate | log.Ltime | log.Lshortfile)
var info = log.New(&buf, "INFO : ", log.Ldate | log.Ltime | log.Lshortfile)
var warn = log.New(&buf, "WARN : ", log.Ldate | log.Ltime | log.Lshortfile)
var err = log.New(&buf, "ERROR: ", log.Ldate | log.Ltime | log.Lshortfile)

var TRACE = Log(&disabled)
var INFO = Log(info)
var WARN = Log(warn)
var ERROR = Log(err)

func SetLevelStr(l string) {
    SetLevel(levelMap[strings.ToLower(l)])
}

func SetLevel(l int) {
    level = l
    if (LEVEL_TRACE >= level) {
        TRACE = trace
    } else {
        TRACE = &disabled
    }
    if (LEVEL_INFO >= level) {
        INFO = info
    } else {
        INFO = &disabled
    }
    if (LEVEL_WARN >= level) {
        WARN = warn
    } else {
        WARN = &disabled
    }
    if (LEVEL_ERROR >= level) {
        ERROR = err
    } else {
        ERROR = &disabled
    }
}

func Tail(n int) string {
    // TODO 
    return buf.String()
}

