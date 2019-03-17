package agentutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io"
	"net"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"
)

// 安装器
type Installer struct {
	Master       string
	Dir          string
	Host         string
	Port         int
	AuthUsername string
	AuthPassword string
	Timeout      time.Duration
	GroupId      string

	HostName string
	HostIP   string
	OS       string
	Arch     string

	Logs        []string
	IsInstalled bool
}

// 获取新对象
func NewInstaller() *Installer {
	return &Installer{}
}

// 安装Agent
func (this *Installer) Start() error {
	this.log("start")

	if len(this.Master) == 0 {
		return errors.New("'master' should not be empty")
	}

	if len(this.Dir) == 0 {
		return errors.New("'dir' should not be empty")
	}

	var hostKeyCallback ssh.HostKeyCallback = nil
	if lists.Contains([]string{"linux", "darwin"}, runtime.GOOS) {
		user1, err := user.Current()
		if err == nil {
			file := user1.HomeDir + "/.ssh/known_hosts"
			if files.NewFile(file).Exists() {
				callback, err := knownhosts.New(file)
				if err != nil {
					logs.Error(err)
				} else {
					hostKeyCallback = callback
				}
			}
		}
	}

	if hostKeyCallback == nil {
		hostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}
	}
	config := &ssh.ClientConfig{
		User: this.AuthUsername,
		Auth: []ssh.AuthMethod{
			ssh.Password(this.AuthPassword),
		},
		HostKeyCallback: hostKeyCallback,
		Timeout:         this.Timeout,
	}

	this.log("connecting")
	client, err := ssh.Dial("tcp", this.Host+":"+fmt.Sprintf("%d", this.Port), config)
	if err != nil {
		return err
	}
	defer client.Close()

	// hostname
	this.log("get hostname")
	hostName, _, err := this.runCmdOnSSH(client, "hostname")
	if err != nil {
		return err
	}
	this.HostName = string(bytes.TrimSpace(hostName))

	// os
	this.log("get os and arch")
	uname, _, err := this.runCmdOnSSH(client, "uname -a")
	if err != nil {
		return err
	}
	if strings.Index(string(uname), "Darwin") >= 0 {
		this.OS = "darwin"
	} else if strings.Index(string(uname), "Linux") >= 0 {
		this.OS = "linux"
	} else {
		return errors.New("installer only supports darwin and linux")
	}

	if strings.Index(string(uname), "x86_64") > 0 {
		this.Arch = "amd64"
	} else {
		this.Arch = "386"
	}

	// upload installer
	this.log("finding installer file")
	filename := "agentinstaller_" + this.OS + "_" + this.Arch
	installerFile := files.NewFile(Tea.Root + "/installers/" + filename)
	if !installerFile.Exists() {
		return errors.New("installer file '" + filename + "' not found")
	}

	this.log("sftp connecting")
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	this.log("create installer file on /tmp")
	writer, err := sftpClient.Create("/tmp/agentinstaller")
	if err != nil {
		return err
	}
	isInstallerWriterClosed := false
	defer func() {
		if !isInstallerWriterClosed {
			writer.Close()
		}

		// 删除
		this.runCmdOnSSH(client, "unlink /tmp/agentinstaller")
	}()

	this.log("open installer file")
	reader, err := os.OpenFile(installerFile.Path(), os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer reader.Close()

	this.log("copy installer file to host")
	n, err := io.Copy(writer, reader)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("copy installer failed")
	}
	isInstallerWriterClosed = true
	writer.Close() // 明确close一次，以便于下面的chmod和运行

	// chmod
	this.log("chmod")
	_, _, err = this.runCmdOnSSH(client, "chmod 777 /tmp/agentinstaller")
	if err != nil {
		return err
	}

	// run
	this.log("installing")

	agentList, err := agents.SharedAgentList()
	if err != nil {
		return err
	}

	// 创建主机信息
	agent := agents.NewAgentConfig()
	this.log("create new agent " + agent.Id)
	agent.Name = string(bytes.TrimSpace(hostName))
	agent.AutoUpdates = true
	agent.CheckDisconnections = true
	agent.AllowAll = true
	agent.On = true
	agent.Key = stringutil.Rand(32)
	if len(this.GroupId) > 0 {
		agent.AddGroup(this.GroupId)
	}
	err = agent.Save()
	if err != nil {
		return err
	}
	agentList.AddAgent(agent.Filename())
	err = agentList.Save()
	if err != nil {
		return err
	}

	newAgentCreated := false
	defer func() {
		if !newAgentCreated {
			this.log("delete agent " + agent.Id)
			err = agent.Delete()
			if err != nil {
				logs.Error(err)
			}

			agentList.RemoveAgent(agent.Filename())
			err = agentList.Save()
			if err != nil {
				logs.Error(err)
			}
		}
	}()

	output, stderr, err := this.runCmdOnSSH(client, "/tmp/agentinstaller -dir=\""+this.Dir+"\" -master=\""+this.Master+"\" -id=\""+agent.Id+"\" -key=\""+agent.Key+"\"")
	if err != nil {
		return errors.New(err.Error() + ":" + string(stderr))
	}

	outputString := strings.TrimSpace(string(output))
	if len(outputString) == 0 {
		return errors.New("start failed: no response:" + string(stderr))
	}

	m := maps.Map{}
	err = json.Unmarshal([]byte(outputString), &m)
	if err != nil {
		return err
	}

	errString := m.GetString("err")
	isInstalled := m.GetBool("isInstalled")
	ip := m.GetString("ip")
	this.HostIP = ip
	if isInstalled {
		this.IsInstalled = true
		if len(errString) > 0 {
			return errors.New(errString)
		} else {
			newAgentCreated = true

			// 保存IP
			agent.Host = ip
			err := agent.Save()
			if err != nil {
				logs.Error(err)
			}

			this.log("finished")
			return nil
		}
	}

	return errors.New("error response:" + errString)
}

// 通过SSH运行一个命令
func (this *Installer) runCmdOnSSH(client *ssh.Client, cmd string) (stdoutBytes []byte, stderrBytes []byte, err error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}
	defer session.Close()

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	session.Stdout = stdout
	session.Stderr = stderr
	err = session.Run(cmd)
	if err != nil {
		return stdout.Bytes(), stderr.Bytes(), err
	}
	return stdout.Bytes(), stderr.Bytes(), nil
}

// 记录日志
func (this *Installer) log(message string) {
	this.Logs = append(this.Logs, message)
}