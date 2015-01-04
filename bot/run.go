package bot

import(
    "os/exec"
    "strings"
    "github.com/mduszyk/foobot/proto"
    "github.com/mduszyk/foobot/log"
    "github.com/mduszyk/foobot/module"
)

type Run struct {
}

func NewRunModule() *Run {
    return &Run{}
}

func runCommand(cmd string, args string) string {
    log.TRACE.Printf("Command to execute: %s", cmd)

    arg_list := strings.Fields(args);
    log.TRACE.Printf("Argument's lenght: %d", len(arg_list))
    if len(arg_list) != 0 {
        log.TRACE.Printf("Argument list %s", arg_list)
    }

    var out []byte
    var err error
    if len(arg_list) == 0 {
        log.TRACE.Printf("Executing command %s with no arguments", cmd)
        out, err = exec.Command(cmd).CombinedOutput()
    } else {
        log.TRACE.Printf("Executing command %s with arguments %s",
                         cmd, arg_list)
        out, err = exec.Command(cmd, arg_list...).CombinedOutput()
    }

    if err != nil {
        log.ERROR.Printf(`Failed executing command %s,
                         arg_list = %s, error: %s`, cmd, arg_list, err)
    } else if out == nil{
        log.TRACE.Printf("Output is nil")
    } else {
        log.INFO.Printf("Command %s successfully executed.", cmd)
        log.TRACE.Printf("Command's output: %s", out)
    }
    return string(out[:])
}

func (r *Run) CMD_(msg *proto.Msg) string {
    rsp := ""

    if len(msg.Cmd) == 0 {
        rsp = "Command is empty"
    } else {
        rsp = runCommand(msg.Cmd, msg.Args)
    }

    return rsp
}

func (r *Run) Handle(msg *proto.Msg) string {
    return module.CallCmdMethod(r, msg)
}

