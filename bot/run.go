package bot

import(
    "os/exec"
    "strings"
    "github.com/mduszyk/foobot/proto"
    "github.com/mduszyk/foobot/log"
)

type Run struct {
}

func NewRunModule() *Run {
    return &Run{}
}

func runCommand(cmd string, args string) string {
    log.TRACE.Printf("Command to execute: %s", cmd)

    arg_list := strings.Split(args, " ");
    log.TRACE.Printf("Argument list: %s", arg_list)
    log.TRACE.Printf("Argument's lenght: %d", len(arg_list))

    out, err := exec.Command(cmd, arg_list...).Output()
    if err != nil {
        log.ERROR.Printf(`Failed executing command %s,
                                arg_list = %s, error: %s`, cmd, arg_list, err)
    } else {
        log.INFO.Printf("Command %s successfully executed.", cmd)
        log.TRACE.Printf("Command's output: %s", out)
    }
    return string(out[:])
}

func (r *Run) Handle(msg *proto.Msg) string {
    rsp := ""

    if len(msg.Cmd) == 0 {
        rsp = "Command is empty"
    } else {
        rsp = runCommand(msg.Cmd, msg.Args)
    }

    return rsp
}
