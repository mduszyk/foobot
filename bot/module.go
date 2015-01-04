package bot

import(
    "reflect"
    "github.com/mduszyk/foobot/proto"
    "github.com/mduszyk/foobot/log"
)

func CallModule(module interface{}, msg *proto.Msg) string {
    modValue := reflect.ValueOf(module)

    prefix := "CMD_"
    name := prefix + msg.Cmd

    method := modValue.MethodByName(name)
    if !method.IsValid() {
        name = prefix
        method = modValue.MethodByName(name)
    }

    rsp := ""

    if method.IsValid() {
        log.TRACE.Printf("Calling cmd method: %s.%s", modValue.Type(), name)
        in := []reflect.Value{reflect.ValueOf(msg)}
        rsp = method.Call(in)[0].Interface().(string)
    }

    return rsp

}
