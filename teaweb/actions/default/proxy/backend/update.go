package backend

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction actions.Action

// 修改后端服务器
func (this *UpdateAction) Run(params struct {
	Server  string
	Backend string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["proxy"] = server
	this.Data["selectedTab"] = "backend"
	this.Data["filename"] = server.Filename

	backend := server.FindBackend(params.Backend)
	if backend == nil {
		this.Fail("找不到要修改的后端服务器")
	}

	backend.Validate()

	this.Data["backend"] = maps.Map{
		"id":          backend.Id,
		"address":     backend.Address,
		"code":        backend.Code,
		"weight":      backend.Weight,
		"failTimeout": int(backend.FailTimeoutDuration().Seconds()),
		"on":          backend.On,
		"maxConns":    backend.MaxConns,
		"maxFails":    backend.MaxFails,
		"isDown":      backend.IsDown,
		"isBackup":    backend.IsBackup,
	}

	this.Show()
}

// 提交
func (this *UpdateAction) RunPost(params struct {
	Server      string
	BackendId   string
	Address     string
	Weight      uint
	On          bool
	Code        string
	FailTimeout uint
	MaxFails    uint
	MaxConns    uint
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

	backend := server.FindBackend(params.BackendId)
	if backend == nil {
		this.Fail("找不到要修改的后端服务器")
	}

	backend.Address = params.Address
	backend.Weight = params.Weight
	backend.On = params.On
	backend.IsDown = false
	backend.Code = params.Code
	backend.FailTimeout = fmt.Sprintf("%d", params.FailTimeout) + "s"
	backend.MaxFails = params.MaxFails
	backend.MaxConns = params.MaxConns
	backend.IsBackup = params.IsBackup

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	global.NotifyChange()

	this.Success()
}
