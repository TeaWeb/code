package backend

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/scheduling"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type SchedulingAction actions.Action

// 调度算法
func (this *SchedulingAction) Run(params struct {
	Server string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["proxy"] = server
	this.Data["filename"] = server.Filename
	this.Data["selectedTab"] = "backend"

	if server.Scheduling == nil {
		server.Scheduling = &teaconfigs.SchedulingConfig{
			Code:    "random",
			Options: maps.Map{},
		}
	}
	this.Data["scheduling"] = server.Scheduling
	this.Data["schedulingTypes"] = scheduling.AllSchedulingTypes()

	this.Show()
}

// 保存提交
func (this *SchedulingAction) RunPost(params struct {
	Server      string
	Type        string
	HashKey     string
	StickyType  string
	StickyParam string
	Must        *actions.Must
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	options := maps.Map{}
	if params.Type == "hash" {
		params.Must.
			Field("hashKey", params.HashKey).
			Require("请输入Key")

		options["key"] = params.HashKey
	} else if params.Type == "sticky" {
		params.Must.
			Field("stickyType", params.StickyType).
			Require("请选择参数类型").
			Field("stickyParam", params.StickyParam).
			Require("请输入参数名").
			Match("^[a-zA-Z0-9]+$", "参数名只能是英文字母和数字的组合").
			MaxCharacters(50, "参数名长度不能超过50位")

		options["type"] = params.StickyType
		options["param"] = params.StickyParam
	}

	if scheduling.FindSchedulingType(params.Type) == nil {
		this.Fail("不支持此种算法")
	}

	server.Scheduling = &teaconfigs.SchedulingConfig{
		Code:    params.Type,
		Options: options,
	}

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	if len(server.Backends) > 0 {
		proxyutils.NotifyChange()
	}

	this.Success()
}
