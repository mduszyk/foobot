package bot

import(
    "fmt"
    "strconv"
    "crypto/sha256"
    "github.com/mduszyk/foobot/log"
    "github.com/mduszyk/foobot/proto"
    "github.com/mduszyk/foobot/conf"
    "github.com/mduszyk/foobot/module"
)

const PASS_SHA256 = `114ac7740c0b09ce0c97dd44f04aa8fae156a4221dc7e03a48f64072adfd81b8`

type AuthModule struct {
    auths map[string]int
}

func NewAuthModule() *AuthModule {
    return &AuthModule{
        auths: make(map[string]int),
    }
}

func (a *AuthModule) Login(user string, pass string) bool {
    sum := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))

    /* log.TRACE.Printf("Auth login, sum: %s", sum) */
    if PASS_SHA256 == sum || conf.Get("bot.pass") == sum {
        a.auths[user] = 0
        log.TRACE.Printf("Auth login success, user: %s", user)
        return true
    }

    log.TRACE.Printf("Auth login failed, user: %s", user)
    return false
}

func (a *AuthModule) Verify(user string) bool {
    _, ok := a.auths[user]
    return ok
}

func (a *AuthModule) Add(user string) {
    _, ok := a.auths[user]
    if !ok {
        a.auths[user] = 1
    }
}

func (a *AuthModule) Rm(user string) {
    v, ok := a.auths[user]
    if ok && v == 1 {
        delete(a.auths, user)
    }
}

func (a *AuthModule) CMD_(msg *proto.Msg) string {
    return a.CMD_list(msg)
}

func (a *AuthModule) CMD_list(msg *proto.Msg) string {
    log.TRACE.Printf("Auth list")
    rsp := ""
    for k, v := range a.auths {
        rsp += k + ": " + strconv.Itoa(v) + "\n"
    }

    return rsp
}

func (a *AuthModule) CMD_add(msg *proto.Msg) string {
    a.Add(msg.Args)
    return a.CMD_list(msg)
}

func (a *AuthModule) CMD_rm(msg *proto.Msg) string {
    a.Rm(msg.Args)
    return a.CMD_list(msg)
}

func (a *AuthModule) Handle(msg *proto.Msg) string {
    if !a.Verify(msg.User) {
        a.Login(msg.User, msg.Raw)
        return a.CMD_list(msg)
    }

    return module.CallCmdMethod(a, msg)
}
