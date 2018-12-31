package proxy

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"strings"
)

type UpdateAction actions.Action

func (this *UpdateAction) Run(params struct {
	Filename string
}) {
	if len(params.Filename) == 0 {
		this.Fail("配置文件读取失败")
	}

	reader, err := files.NewReader(Tea.ConfigFile(params.Filename))
	if err != nil {
		this.Fail("配置文件读取失败")
	}

	config := &teaconfigs.ServerConfig{}
	err = reader.ReadYAML(config)
	if err != nil {
		logs.Error(err)
		this.Fail("配置文件读取失败")
	}

	this.Data["filename"] = params.Filename
	this.Data["server"] = config

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	Auth           *helpers.UserMustAuth
	Filename       string
	Name           []string
	ListenAddress  []string
	ListenPort     []int
	BackendAddress []string
	BackendPort    []int
	Must           *actions.Must
}) {
	if len(params.Filename) == 0 {
		this.Fail("配置文件名错误")
	}
	configFile := files.NewFile(Tea.ConfigFile(params.Filename))
	if !configFile.IsFile() {
		this.Fail("找不到要修改的配置")
	}

	if len(params.Name) == 0 {
		this.Fail("域名不能为空")
	}

	for index, name := range params.Name {
		name = strings.TrimSpace(name)
		if len(name) == 0 {
			this.Fail("域名不能为空")
		}
		params.Name[index] = name
	}

	for index, address := range params.ListenAddress {
		address = strings.TrimSpace(address)
		if len(address) == 0 {
			this.Fail("访问地址不能为空")
		}
		params.ListenAddress[index] = address
	}

	for index, port := range params.ListenPort {
		if port <= 0 || port >= 65535 {
			this.Fail("访问地址端口错误")
		}
		params.ListenPort[index] = port
	}

	for index, address := range params.BackendAddress {
		address = strings.TrimSpace(address)
		if len(address) == 0 {
			this.Fail("后端地址不能为空")
		}
		params.BackendAddress[index] = address
	}

	for index, port := range params.BackendPort {
		if port <= 0 || port >= 65535 {
			this.Fail("后端地址端口错误")
		}
		params.BackendPort[index] = port
	}

	// 保存
	server := teaconfigs.NewServerConfig()
	server.AddName(params.Name ...)
	for index, address := range params.ListenAddress {
		if index > len(params.ListenPort)-1 {
			continue
		}

		server.AddListen(fmt.Sprintf("%s:%d", address, params.ListenPort[index]))
	}
	for index, address := range params.BackendAddress {
		if index > len(params.BackendPort)-1 {
			continue
		}

		backend := &teaconfigs.ServerBackendConfig{
			Address: fmt.Sprintf("%s:%d", address, params.BackendPort[index]),
		}
		server.AddBackend(backend)
	}

	err := server.WriteToFile(configFile.Path())
	if err != nil {
		this.Fail("配置文件写入失败")
	}

	// 重启
	proxyutils.NotifyChange()

	this.Next("/proxy", nil, "").Success("服务保存成功")
}
