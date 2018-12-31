package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type HeadersAction actions.Action

// 自定义Http Header
func (this *HeadersAction) Run(params struct {
	Server     string // 必填
	LocationId string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	location := server.FindLocation(params.LocationId)
	if location == nil {
		this.Fail("找不到要修改的Location")
	}
	this.Data["location"] = maps.Map{
		"id":          location.Id,
		"pattern":     location.PatternString(),
		"fastcgi":     location.Fastcgi,
		"rewrite":     location.Rewrite,
		"headers":     location.Headers,
		"cachePolicy": location.CachePolicy,
	}

	this.Data["selectedTab"] = "location"
	this.Data["selectedSubTab"] = "headers"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = server

	this.Data["headerQuery"] = maps.Map{
		"server":     params.Server,
		"locationId": params.LocationId,
	}

	this.Show()
}
