package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
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
	this.Data["operators"] = teaconfigs.AllRequestOperators()

	// 变量
	this.Data["variables"] = proxyutils.DefaultRequestVariables()

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

	CondParams []string
	CondOps    []string
	CondValues []string
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
	location.Cond = []*teaconfigs.RequestCond{}
	if len(params.CondParams) > 0 {
		for index, param := range params.CondParams {
			if index < len(params.CondOps) && index < len(params.CondValues) {
				cond := teaconfigs.NewRequestCond()
				cond.Param = param
				cond.Value = params.CondValues[index]
				cond.Operator = params.CondOps[index]
				err := cond.Validate()
				if err != nil {
					this.Fail("匹配条件\"" + cond.Param + " " + cond.Value + "\"校验失败：" + err.Error())
				}
				location.AddCond(cond)
			}
		}
	}

	location.SetPattern(params.Pattern, params.PatternType, params.IsCaseInsensitive, params.IsReverse)
	location.On = params.On
	location.CacheOn = true
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

	index := []string{}
	for _, i := range params.Index {
		if len(i) > 0 && !lists.ContainsString(index, i) {
			index = append(index, i)
		}
	}
	location.Index = index
	server.AddLocation(location)

	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()
	this.Success()
}
