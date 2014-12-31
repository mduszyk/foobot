package bot

import(
    "github.com/mduszyk/foobot/proto"
)

type Run struct {
}

func NewRunModule() *Run {
    return &Run{}
}

func (r *Run) Handle(msg *proto.Msg) string {
    rsp := ""
    rsp = "run module invoked. WIP."

    return rsp
}
