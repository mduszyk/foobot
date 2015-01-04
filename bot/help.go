package bot

import(
    "github.com/mduszyk/foobot/proto"
    "github.com/mduszyk/foobot/module"
)

type HelpModule struct {
    bot *Bot
}

func NewHelpModule(bot *Bot) *HelpModule {
    return &HelpModule{bot}
}

func (m *HelpModule) CMD_(msg *proto.Msg) string {
    rsp := ""
    for k, v := range m.bot.GetModules() {
        rsp += k
        methods := proto.CmdMethods(v)
        if len(methods) > 0 {
            rsp += " (" + methods + ")"
        }
        rsp += "\n"
    }

    return rsp
}

func (m *HelpModule) Handle(msg *proto.Msg) string {
    return module.CallCmdMethod(m, msg)
}
