package backend

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type UpdateAction actions.Action

// 修改后端服务器
func (this *UpdateAction) Run(params struct {
	ServerId   string
	LocationId string
	Websocket  bool
	Backend    string
	From       string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = server
	if len(params.LocationId) > 0 {
		this.Data["selectedTab"] = "location"
	} else {
		this.Data["selectedTab"] = "backend"
	}
	this.Data["locationId"] = params.LocationId
	this.Data["websocket"] = types.Int(params.Websocket)
	this.Data["from"] = params.From

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}
	backend := backendList.FindBackend(params.Backend)
	if backend == nil {
		this.Fail("找不到要修改的后端服务器")
	}

	backend.Validate()

	this.Data["backend"] = maps.Map{
		"id":          backend.Id,
		"address":     backend.Address,
		"scheme":      backend.Scheme,
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
	ServerId    string
	LocationId  string
	Websocket   bool
	BackendId   string
	Address     string
	Scheme      string
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

	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}

	backend := backendList.FindBackend(params.BackendId)
	if backend == nil {
		this.Fail("找不到要修改的后端服务器")
	}

	backend.Address = params.Address
	backend.Scheme = params.Scheme
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

	proxyutils.NotifyChange()

	this.Success()
}
