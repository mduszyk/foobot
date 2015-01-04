package proto

import(
    "reflect"
    "strings"
)

const CMD_PREFIX = "CMD_"

func CallCmdMethod(module interface{}, msg *Msg) string {
    modValue := reflect.ValueOf(module)

    name := CMD_PREFIX + msg.Cmd

    method := modValue.MethodByName(name)
    if !method.IsValid() {
        name = CMD_PREFIX
        method = modValue.MethodByName(name)
    }

    rsp := ""

    if method.IsValid() {
        in := []reflect.Value{reflect.ValueOf(msg)}
        rsp = method.Call(in)[0].Interface().(string)
    }

    return rsp

}

func CmdMethods(module interface{}) string {
    modType := reflect.TypeOf(module)
    rsp := ""
    for i := 0; i < modType.NumMethod(); i++ {
        methodName := modType.Method(i).Name
        if strings.HasPrefix(methodName, CMD_PREFIX) {
            if len(rsp) > 0 {
                rsp += ", "
            }
            methodName = strings.Replace(methodName, CMD_PREFIX, "", -1)
            if methodName == "" {
                methodName = "\"\""
            }
            rsp += methodName
        }
    }

    return rsp
}

