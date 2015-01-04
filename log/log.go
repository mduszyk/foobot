package log

import(
    /* "fmt" */
    "os"
    "io"
    "log"
    "strconv"
    "strings"
    /* "sync/atomic" */
    "github.com/mduszyk/foobot/proto"
    "github.com/mduszyk/foobot/module"
)

const LEVEL_TRACE = 0
const LEVEL_INFO = 1
const LEVEL_WARN = 2
const LEVEL_ERROR = 3

var STR_TO_LEVEL = map[string]int{
    "trace": LEVEL_TRACE,
    "info": LEVEL_INFO,
    "warn": LEVEL_WARN,
    "error": LEVEL_ERROR,
}

var LEVEL_TO_STR = []string{"trace", "info", "warn", "error"}


type MutableWriter struct {
    writer io.Writer
}

func (mw *MutableWriter) Write(p []byte) (n int, err error) {
    return mw.writer.Write(p)
}

func (mw *MutableWriter) SetWriter(w io.Writer) {
    mw.writer = w
}

type CircularWriter struct {
    buf [][]byte
    index int
    maxLines int
    maxLineSize int
}

func NewCircularWriter(maxLines int, maxLineSize int) *CircularWriter {
    w := &CircularWriter{
        buf: make([][]byte, maxLines),
        index: 0,
        maxLines: maxLines,
        maxLineSize: maxLineSize,
    }
    for i := 0; i < maxLines; i++ {
        w.buf[i] = make([]byte, maxLineSize)
    }
    return w
}

// TODO sync it, use cas?
func (cw *CircularWriter) Write(p []byte) (n int, err error) {
    /* fmt.Printf("write, p: %s\n", p) */
    copy(cw.buf[cw.index], p)
    cw.index = (cw.index + 1) % cap(cw.buf)
    return len(p), nil
}

// TODO sync it, use cas?
func (cw *CircularWriter) Tail(n int) string {
    rsp := ""
    /* fmt.Printf("tail, n: %d\n", n) */
    index := cw.index
    for i := n; i > 0; i-- {
        j := index - 1 - i
        if j < 0 {
            j = len(cw.buf) + j
        }
        rsp += string(cw.buf[j]) + "\n"
        /* fmt.Printf("tail, buf: %s\n", cw.buf[j]) */
    }

    return rsp
}

type Logger struct {
    buf *CircularWriter
    writer *MutableWriter
    level int
}

var buf = NewCircularWriter(128, 256)
var instance = &Logger{
    buf: buf,
    writer: &MutableWriter{buf},
    level: LEVEL_INFO,
}

func GetLogModule() *Logger {
    return instance
}

type Log interface {
    Printf(format string, v ...interface{})
}

type nullLog struct{}
func (l *nullLog) Printf(format string, v ...interface{}) {}

var disabled nullLog
var trace = log.New(instance.writer, "TRACE: ", log.Ldate | log.Ltime | log.Lshortfile)
var info = log.New(instance.writer, "INFO : ", log.Ldate | log.Ltime | log.Lshortfile)
var warn = log.New(instance.writer, "WARN : ", log.Ldate | log.Ltime | log.Lshortfile)
var err = log.New(instance.writer, "ERROR: ", log.Ldate | log.Ltime | log.Lshortfile)

var TRACE = Log(&disabled)
var INFO = Log(info)
var WARN = Log(warn)
var ERROR = Log(err)

func EnableStderr() {
    instance.writer.SetWriter(io.MultiWriter(buf, os.Stderr))
}

func SetLevelStr(l string) {
    SetLevel(STR_TO_LEVEL[strings.ToLower(l)])
}

func SetLevel(l int) {
    instance.level = l
    if (LEVEL_TRACE >= instance.level) {
        TRACE = trace
    } else {
        TRACE = &disabled
    }
    if (LEVEL_INFO >= instance.level) {
        INFO = info
    } else {
        INFO = &disabled
    }
    if (LEVEL_WARN >= instance.level) {
        WARN = warn
    } else {
        WARN = &disabled
    }
    if (LEVEL_ERROR >= instance.level) {
        ERROR = err
    } else {
        ERROR = &disabled
    }
}

func (l *Logger) CMD_tail(msg *proto.Msg) string {
    n, err := strconv.Atoi(msg.Args)
    if err != nil {
        n = 5
    }
    return l.buf.Tail(n)
}

func (l *Logger) CMD_level(msg *proto.Msg) string {
    if len(msg.Args) > 0 {
        SetLevelStr(msg.Args)
    }
    return "log level " + msg.Args
}

func (l *Logger) CMD_info(msg *proto.Msg) string {
    rsp := "log.max_lines: " + strconv.Itoa(l.buf.maxLines) + "\n"
    rsp += "log.max_line_size: " + strconv.Itoa(l.buf.maxLineSize) + "\n"
    rsp += "log.level: " + LEVEL_TO_STR[l.level] + "\n"
    return rsp
}

func (l *Logger) Handle(msg *proto.Msg) string {
    return module.CallCmdMethod(l, msg)
}

