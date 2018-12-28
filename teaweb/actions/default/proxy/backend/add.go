package backend

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
)

type AddAction actions.Action

// 添加服务器
func (this *AddAction) Run(params struct {
	Server string
	Backup bool
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["selectedTab"] = "backend"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = proxy

	this.Data["isBackup"] = params.Backup

	this.Show()
}

// 提交
func (this *AddAction) RunPost(params struct {
	Server      string
	Address     string
	Weight      uint
	On          bool
	Code        string
	FailTimeout uint
	MaxFails    uint
	MaxConns    uint
	SlowStart   uint
	IsBackup    bool
	Must        *actions.Must
}) {
	params.Must.
		Field("address", params.Address).
		Require("请输入后端服务器地址")

	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	backend := teaconfigs.NewServerBackendConfig()
	backend.Address = params.Address
	backend.Weight = params.Weight
	backend.On = params.On
	backend.IsDown = false
	backend.Code = params.Code
	backend.FailTimeout = fmt.Sprintf("%d", params.FailTimeout) + "s"
	backend.MaxFails = params.MaxFails
	backend.MaxConns = params.MaxConns
	backend.SlowStart = fmt.Sprintf("%d", params.SlowStart) + "s"
	backend.IsBackup = params.IsBackup

	server.Backends = append(server.Backends, backend)
	server.Save()

	global.NotifyChange()

	this.Success()
}
