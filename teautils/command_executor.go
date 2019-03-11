package teautils

import (
	"bytes"
	"errors"
	"os/exec"
)

// 命令执行器
type CommandExecutor struct {
	commands []*Command
}

// 获取新对象
func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

func (this *CommandExecutor) Add(command string, arg ...string) {
	this.commands = append(this.commands, &Command{
		Name: command,
		Args: arg,
	})
}

func (this *CommandExecutor) Run() (output string, err error) {
	if len(this.commands) == 0 {
		return "", errors.New("no commands no run")
	}
	var lastCmd *exec.Cmd = nil
	var data []byte = nil
	for _, command := range this.commands {
		cmd := exec.Command(command.Name, command.Args ...)
		buf := bytes.NewBuffer([]byte{})
		cmd.Stdout = buf
		if lastCmd != nil {
			cmd.Stdin = bytes.NewBuffer(data)
		}
		err = cmd.Start()
		if err != nil {
			return "", err
		}

		err = cmd.Wait()
		if err != nil {
			return "", err
		}
		data = buf.Bytes()

		lastCmd = cmd
	}

	return string(bytes.TrimSpace(data)), nil
}
