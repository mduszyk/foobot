package bot

import(
    "fmt"
    "strconv"
    "crypto/sha256"
	"github.com/mduszyk/foobot/log"
	"github.com/mduszyk/foobot/proto"
	"github.com/mduszyk/foobot/conf"
)

const PASS_SHA256 = `114ac7740c0b09ce0c97dd44f04aa8fae156a4221dc7e03a48f64072adfd81b8`

type AuthModule struct {
    auths map[string]int
}

var instance *AuthModule = nil

func NewAuthModule() *AuthModule {
    if instance == nil {
        instance = &AuthModule{
            auths: make(map[string]int),
        }
    }
    return instance
}

func (a *AuthModule) list() string {
    rsp := ""
    for k, v := range a.auths {
        rsp += k + ": " + strconv.Itoa(v) + "\n"
    }

    return rsp
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

func (a *AuthModule) Handle(msg *proto.Msg) string {
    rsp := ""

    if !a.Verify(msg.User) {
        a.Login(msg.User, msg.Raw)
        return a.list()
    }

    switch msg.Cmd {
        case "":
            log.TRACE.Printf("Auth list")
            rsp = a.list()
        case "list":
            log.TRACE.Printf("Auth list")
            rsp = a.list()
        case "add":
            a.Add(msg.Args)
            rsp = a.list()
        case "rm":
            a.Rm(msg.Args)
            rsp = a.list()
    }

    return rsp
}
