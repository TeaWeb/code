package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/types"
	"regexp"
	"strconv"
)

type AddAction actions.Action

// 添加路径规则
func (this *AddAction) Run(params struct {
	ServerId string
	From     string
	Must     *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = server
	this.Data["selectedTab"] = "location"
	this.Data["selectedSubTab"] = "detail"
	this.Data["from"] = params.From

	this.Data["patternTypes"] = teaconfigs.AllLocationPatternTypes()
	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	accessLog := teaconfigs.NewAccessLogConfig()
	this.Data["accessLogIsInherited"] = true
	this.Data["accessLogs"] = proxyutils.FormatAccessLog([]*teaconfigs.AccessLogConfig{accessLog})

	// 运算符
	this.Data["condOperators"] = shared.AllRequestOperators()

	// 变量
	this.Data["condVariables"] = proxyutils.DefaultRequestVariables()

	//gzip
	this.Data["gzip"] = &teaconfigs.GzipConfig{
		Level:     -1,
		MinLength: "",
		MimeTypes: teaconfigs.DefaultGzipMimeTypes,
	}

	this.Show()
}

// 保存提交
func (this *AddAction) RunPost(params struct {
	ServerId             string
	Name                 string
	Pattern              string
	PatternType          int
	Root                 string
	Charset              string
	Index                []string
	URLPrefix            string `alias:"urlPrefix"`
	MaxBodySize          float64
	MaxBodyUnit          string
	AccessLogIsInherited bool
	EnableStat           bool

	// gzip
	GzipLevel          int8
	GzipMinLength      float64
	GzipMinUnit        string
	GzipMimeTypeValues []string

	RedirectToHttps   bool
	On                bool
	IsReverse         bool
	IsCaseInsensitive bool

	// pages
	PageStatusList    []string
	PageURLList       []string
	PageNewStatusList []string

	ShutdownPageOn     bool
	ShutdownPageURL    string
	ShutdownPageStatus int
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	// 校验正则
	if params.PatternType == teaconfigs.LocationPatternTypeRegexp {
		_, err := regexp.Compile(params.Pattern)
		if err != nil {
			this.Fail("正则表达式校验失败：" + err.Error())
		}
	}

	location := teaconfigs.NewLocation()

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	location.Cond = conds

	location.SetPattern(params.Pattern, params.PatternType, params.IsCaseInsensitive, params.IsReverse)
	location.On = params.On
	location.CacheOn = true
	location.Name = params.Name
	location.Root = params.Root
	location.URLPrefix = params.URLPrefix
	location.Charset = params.Charset
	location.MaxBodySize = strconv.FormatFloat(params.MaxBodySize, 'f', -1, 64) + params.MaxBodyUnit
	if params.AccessLogIsInherited {
		location.AccessLog = []*teaconfigs.AccessLogConfig{}
	} else {
		location.AccessLog = proxyutils.ParseAccessLogForm(this.Request)
	}
	location.DisableStat = !params.EnableStat
	location.RedirectToHttps = params.RedirectToHttps

	// gzip
	// 这里gzipLevel包括0，因为要指定不压缩
	if params.GzipLevel >= 0 && params.GzipLevel <= 9 {
		minLength := strconv.FormatFloat(params.GzipMinLength, 'f', -1, 64) + params.GzipMinUnit
		gzip := &teaconfigs.GzipConfig{
			Level:     params.GzipLevel,
			MinLength: minLength,
			MimeTypes: params.GzipMimeTypeValues,
		}
		location.Gzip = gzip
	} else {
		location.Gzip = nil
	}

	// 特殊页面
	location.Pages = []*teaconfigs.PageConfig{}
	for index, status := range params.PageStatusList {
		page := teaconfigs.NewPageConfig()
		page.Status = []string{status}
		if index < len(params.PageURLList) {
			page.URL = params.PageURLList[index]
		}
		if index < len(params.PageNewStatusList) {
			page.NewStatus = types.Int(params.PageNewStatusList[index])
			if page.NewStatus < 0 {
				page.NewStatus = 0
			}
		}
		location.AddPage(page)
	}

	if location.Shutdown != nil {
		location.Shutdown.On = params.ShutdownPageOn
		location.Shutdown.URL = params.ShutdownPageURL
		location.Shutdown.Status = params.ShutdownPageStatus
	} else if params.ShutdownPageOn {
		location.Shutdown = teaconfigs.NewShutdownConfig()
		location.Shutdown.On = params.ShutdownPageOn
		location.Shutdown.URL = params.ShutdownPageURL
		location.Shutdown.Status = params.ShutdownPageStatus
	}
	if location.Shutdown != nil && location.Shutdown.On && len(location.Shutdown.URL) == 0 {
		this.FailField("shutdownPageURL", "请输入临时关闭页面文件路径")
	}

	// 首页
	index := []string{}
	for _, i := range params.Index {
		if len(i) > 0 && !lists.ContainsString(index, i) {
			index = append(index, i)
		}
	}
	location.Index = index
	server.AddLocation(location)

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()
	this.Success()
}
