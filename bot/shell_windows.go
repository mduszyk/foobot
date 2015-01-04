package bot

import (
    "os/exec"
)

type Shell struct {
    cmd *exec.Cmd
}

func NewShell() *Shell {
    return &Shell{}
}

func (sh *Shell) Start() {
}

func (sh *Shell) Kill() {
}

func (sh *Shell) Interrupt() {
}

func (sh *Shell) Insert(line string) string {
    return "win shell not supported"
}
