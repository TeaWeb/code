package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"strconv"
)

type AddAction actions.Action

// 添加路径规则
func (this *AddAction) Run(params struct {
	Server string
	From   string
	Must   *actions.Must
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["server"] = maps.Map{
		"filename": params.Server,
	}
	this.Data["filename"] = params.Server
	this.Data["proxy"] = server
	this.Data["selectedTab"] = "location"
	this.Data["selectedSubTab"] = "detail"
	this.Data["from"] = params.From

	this.Data["patternTypes"] = teaconfigs.AllLocationPatternTypes()
	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	this.Show()
}

// 保存提交
func (this *AddAction) RunPost(params struct {
	Server            string
	Pattern           string
	PatternType       int
	Root              string
	Charset           string
	Index             []string
	MaxBodySize       float64
	MaxBodyUnit       string
	EnableAccessLog   bool
	GzipLevel         int8
	GzipMinLength     float64
	GzipMinUnit       string
	On                bool
	IsReverse         bool
	IsCaseInsensitive bool
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	// 校验正则
	if params.PatternType == teaconfigs.LocationPatternTypeRegexp {
		_, err := regexp.Compile(params.Pattern)
		if err != nil {
			this.Fail("正则表达式校验失败：" + err.Error())
		}
	}

	location := teaconfigs.NewLocation()
	location.SetPattern(params.Pattern, params.PatternType, params.IsCaseInsensitive, params.IsReverse)
	location.On = params.On
	location.Root = params.Root
	location.Charset = params.Charset
	location.MaxBodySize = strconv.FormatFloat(params.MaxBodySize, 'f', -1, 64) + params.MaxBodyUnit
	location.DisableAccessLog = !params.EnableAccessLog
	location.GzipLevel = params.GzipLevel
	location.GzipMinLength = strconv.FormatFloat(params.GzipMinLength, 'f', -1, 64) + params.GzipMinUnit

	index := []string{}
	for _, i := range params.Index {
		if len(i) > 0 && !lists.Contains(index, i) {
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
