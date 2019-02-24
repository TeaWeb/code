package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"strconv"
)

type UpdateAction actions.Action

// 修改代理服务信息
func (this *UpdateAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}
	this.Data["server"] = server
	this.Data["selectedTab"] = "basic"

	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	this.Show()
}

// 保存提交
func (this *UpdateAction) RunPost(params struct {
	HttpOn          bool
	ServerId        string
	Description     string
	Name            []string
	Listen          []string
	Root            string
	Charset         string
	Index           []string
	MaxBodySize     float64
	MaxBodyUnit     string
	EnableAccessLog bool
	GzipLevel       uint8
	GzipMinLength   float64
	GzipMinUnit     string
	Must            *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	params.Must.
		Field("description", params.Description).
		Require("代理服务名称不能为空")

	server.Http = params.HttpOn
	server.Description = params.Description
	server.Name = params.Name
	server.Listen = params.Listen
	server.Root = params.Root
	server.Charset = params.Charset
	server.Index = params.Index
	server.MaxBodySize = strconv.FormatFloat(params.MaxBodySize, 'f', -1, 64) + params.MaxBodyUnit
	server.DisableAccessLog = !params.EnableAccessLog
	if params.GzipLevel <= 9 {
		server.GzipLevel = params.GzipLevel
	}
	server.GzipMinLength = strconv.FormatFloat(params.GzipMinLength, 'f', -1, 64) + params.GzipMinUnit

	err := server.Validate()
	if err != nil {
		this.Fail("校验失败：" + err.Error())
	}

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 重启
	proxyutils.NotifyChange()

	this.Success()
}
