package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) Run(params struct {
	Server      string
	LocationId  string
	From        string
	ShowSpecial bool
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
	this.Data["showSpecial"] = params.ShowSpecial

	location := server.FindLocation(params.LocationId)
	if location == nil {
		this.Fail("找不到要修改的Location")
	}

	this.Data["patternTypes"] = teaconfigs.AllLocationPatternTypes()
	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	this.Data["location"] = maps.Map{
		"id":                location.Id,
		"on":                location.On,
		"pattern":           location.PatternString(),
		"type":              location.PatternType(),
		"isReverse":         location.IsReverse(),
		"isCaseInsensitive": location.IsCaseInsensitive(),
		"root":              location.Root,
		"index":             location.Index,
		"charset":           location.Charset,

		// 菜单用
		"rewrite":     location.Rewrite,
		"headers":     location.Headers,
		"fastcgi":     location.Fastcgi,
		"cachePolicy": location.CachePolicy,
	}

	this.Show()
}

// 保存修改
func (this *UpdateAction) RunPost(params struct {
	Server            string
	LocationId        string
	Pattern           string
	PatternType       int
	Root              string
	Charset           string
	Index             []string
	On                bool
	IsReverse         bool
	IsCaseInsensitive bool
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	location := server.FindLocation(params.LocationId)
	if location == nil {
		this.Fail("找不到要修改的Location")
	}
	location.SetPattern(params.Pattern, params.PatternType, params.IsCaseInsensitive, params.IsReverse)
	location.On = params.On
	location.Root = params.Root
	location.Charset = params.Charset

	index := []string{}
	for _, i := range params.Index {
		if len(i) > 0 && !lists.Contains(index, i) {
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
