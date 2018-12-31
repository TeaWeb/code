package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type DetailAction actions.Action

// 路径规则详情
func (this *DetailAction) Run(params struct {
	Server     string
	LocationId string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["server"] = maps.Map{
		"filename": server.Filename,
	}

	location := server.FindLocation(params.LocationId)
	if location == nil {
		this.Fail("找不到要修改的路径配置")
	}

	if location.Index == nil {
		location.Index = []string{}
	}

	this.Data["selectedTab"] = "location"
	this.Data["selectedSubTab"] = "detail"
	this.Data["filename"] = params.Server

	this.Data["location"] = maps.Map{
		"on":              location.On,
		"id":              location.Id,
		"type":            location.PatternType(),
		"pattern":         location.PatternString(),
		"caseInsensitive": location.IsCaseInsensitive(),
		"reverse":         location.IsReverse(),
		"root":            location.Root,
		"charset":         location.Charset,
		"index":           location.Index,
		"fastcgi":         location.Fastcgi,
		"headers":         location.Headers,
		"cachePolicy":     location.CachePolicy,
		"rewrite":         location.Rewrite,
	}
	this.Data["proxy"] = server

	// 字符集
	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	this.Show()
}
