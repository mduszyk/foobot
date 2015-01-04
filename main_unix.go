// +build linux 

package main

import (
    "github.com/VividCortex/godaemon"
)

func daemon() {
    godaemon.MakeDaemon(&godaemon.DaemonAttr{})
}
