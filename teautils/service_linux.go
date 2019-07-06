// +build linux

package teautils

import (
	"errors"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"io/ioutil"
	"os/exec"
	"regexp"
)

var serviceFile = "/etc/init.d/teaweb"

// 安装服务
func (this *ServiceManager) Install(exePath string, args []string) error {
	scriptFile := Tea.Root + "/scripts/teaweb"
	if !files.NewFile(scriptFile).Exists() {
		return errors.New("'scripts/teaweb' file not exists")
	}

	data, err := ioutil.ReadFile(scriptFile)
	if err != nil {
		return err
	}

	data = regexp.MustCompile("INSTALL_DIR=.+").ReplaceAll(data, []byte("INSTALL_DIR="+Tea.Root))
	err = ioutil.WriteFile(serviceFile, data, 0777)
	if err != nil {
		return err
	}

	chkCmd, err := exec.LookPath("chkconfig")
	if err != nil {
		return err
	}

	err = exec.Command(chkCmd, "--add", "teaweb").Start()
	if err != nil {
		return err
	}

	return nil
}

// 启动服务
func (this *ServiceManager) Start() error {
	return exec.Command("service", "teaweb", "start").Start()
}

// 删除服务
func (this *ServiceManager) Uninstall() error {
	f := files.NewFile(serviceFile)
	if f.Exists() {
		return f.Delete()
	}
	return nil
}
