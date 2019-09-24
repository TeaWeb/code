package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"strconv"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) Run(params struct {
	ServerId    string
	LocationId  string
	From        string
	ShowSpecial bool
}) {
	_, location := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "detail")

	this.Data["from"] = params.From
	this.Data["showSpecial"] = params.ShowSpecial

	this.Data["patternTypes"] = teaconfigs.AllLocationPatternTypes()
	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets
	this.Data["accessLogIsInherited"] = len(location.AccessLog) == 0
	if len(location.AccessLog) == 0 {
		location.AccessLog = []*teaconfigs.AccessLogConfig{teaconfigs.NewAccessLogConfig()}
	}
	this.Data["accessLogs"] = proxyutils.FormatAccessLog(location.AccessLog)

	if len(location.Pages) == 0 {
		location.Pages = []*teaconfigs.PageConfig{}
	}

	this.Data["location"] = maps.Map{
		"id":                location.Id,
		"on":                location.On,
		"pattern":           location.PatternString(),
		"type":              location.PatternType(),
		"name":              location.Name,
		"isReverse":         location.IsReverse(),
		"isCaseInsensitive": location.IsCaseInsensitive(),
		"root":              location.Root,
		"index":             location.Index,
		"charset":           location.Charset,
		"maxBodySize":       location.MaxBodySize,
		"enableStat":        !location.DisableStat,
		"gzipLevel":         location.GzipLevel,
		"gzipMinLength":     location.GzipMinLength,
		"redirectToHttps":   location.RedirectToHttps,
		"conds":             location.Cond,

		// 菜单用
		"rewrite":     location.Rewrite,
		"headers":     location.Headers,
		"fastcgi":     location.Fastcgi,
		"cachePolicy": location.CachePolicy,
		"websocket":   location.Websocket,
		"backends":    location.Backends,
		"wafOn":       location.WAFOn,
		"wafId":       location.WafId,

		"shutdownPage":   location.ShutdownPage,
		"shutdownPageOn": location.ShutdownPageOn,
		"pages":          location.Pages,
	}

	// 运算符
	this.Data["operators"] = shared.AllRequestOperators()

	// 变量
	this.Data["variables"] = proxyutils.DefaultRequestVariables()

	this.Show()
}

// 保存修改
func (this *UpdateAction) RunPost(params struct {
	ServerId             string
	LocationId           string
	Pattern              string
	PatternType          int
	Name                 string
	Root                 string
	Charset              string
	Index                []string
	MaxBodySize          float64
	MaxBodyUnit          string
	AccessLogIsInherited bool
	EnableStat           bool
	GzipLevel            int8
	GzipMinLength        float64
	GzipMinUnit          string
	RedirectToHttps      bool
	On                   bool
	IsReverse            bool
	IsCaseInsensitive    bool

	PageStatus []string
	PageURL    []string

	ShutdownPageOn bool
	ShutdownPage   string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	location := server.FindLocation(params.LocationId)
	if location == nil {
		this.Fail("找不到要修改的Location")
	}

	// 校验正则
	if params.PatternType == teaconfigs.LocationPatternTypeRegexp {
		_, err := regexp.Compile(params.Pattern)
		if err != nil {
			this.Fail("正则表达式校验失败：" + err.Error())
		}
	}

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	location.Cond = conds

	location.SetPattern(params.Pattern, params.PatternType, params.IsCaseInsensitive, params.IsReverse)
	location.On = params.On
	location.Name = params.Name
	location.Root = params.Root
	location.Charset = params.Charset
	location.MaxBodySize = strconv.FormatFloat(params.MaxBodySize, 'f', -1, 64) + params.MaxBodyUnit
	if params.AccessLogIsInherited {
		location.AccessLog = []*teaconfigs.AccessLogConfig{}
	} else {
		location.AccessLog = proxyutils.ParseAccessLogForm(this.Request)
	}
	location.DisableStat = !params.EnableStat
	location.GzipLevel = params.GzipLevel
	location.GzipMinLength = strconv.FormatFloat(params.GzipMinLength, 'f', -1, 64) + params.GzipMinUnit
	location.RedirectToHttps = params.RedirectToHttps

	// 特殊页面
	location.Pages = []*teaconfigs.PageConfig{}
	for index, status := range params.PageStatus {
		if index < len(params.PageURL) {
			page := teaconfigs.NewPageConfig()
			page.Status = []string{status}
			page.URL = params.PageURL[index]
			location.AddPage(page)
		}
	}

	location.ShutdownPageOn = params.ShutdownPageOn
	if location.ShutdownPageOn && len(params.ShutdownPage) == 0 {
		this.FailField("shutdownPage", "请输入临时关闭页面文件路径")
	}
	location.ShutdownPage = params.ShutdownPage

	// 首页
	index := []string{}
	for _, i := range params.Index {
		if len(i) > 0 && !lists.ContainsString(index, i) {
			index = append(index, i)
		}
	}
	location.Index = index

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()
	this.Success()
}
