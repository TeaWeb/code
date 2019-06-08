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

	this.Data["accessLogs"] = proxyutils.FormatAccessLog(server.AccessLog)

	this.Show()
}

// 保存提交
func (this *UpdateAction) RunPost(params struct {
	HttpOn      bool
	ServerId    string
	Description string
	Name        []string
	Listen      []string
	Root        string
	Charset     string
	Index       []string
	MaxBodySize float64
	MaxBodyUnit string

	EnableStat    bool
	GzipLevel     uint8
	GzipMinLength float64
	GzipMinUnit   string
	CacheStatic   bool

	PageStatus []string
	PageURL    []string

	ShutdownPageOn bool
	ShutdownPage   string

	RedirectToHttps bool

	Must *actions.Must
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

	// 访问日志
	server.AccessLog = proxyutils.ParseAccessLogForm(this.Request)

	server.DisableStat = !params.EnableStat
	if params.GzipLevel <= 9 {
		server.GzipLevel = params.GzipLevel
	}
	server.GzipMinLength = strconv.FormatFloat(params.GzipMinLength, 'f', -1, 64) + params.GzipMinUnit
	server.CacheStatic = params.CacheStatic

	server.Pages = []*teaconfigs.PageConfig{}
	for index, status := range params.PageStatus {
		if index < len(params.PageURL) {
			page := teaconfigs.NewPageConfig()
			page.Status = []string{status}
			page.URL = params.PageURL[index]
			server.AddPage(page)
		}
	}

	server.ShutdownPageOn = params.ShutdownPageOn
	if server.ShutdownPageOn && len(params.ShutdownPage) == 0 {
		this.FailField("shutdownPage", "请输入临时关闭页面文件路径")
	}
	server.ShutdownPage = params.ShutdownPage

	server.RedirectToHttps = params.RedirectToHttps

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
