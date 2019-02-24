package locations

import (
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type DetailAction actions.Action

// 路径规则详情
func (this *DetailAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	server, location := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "detail")

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
		"maxBodySize":     location.MaxBodySize,
		"enableAccessLog": !location.DisableAccessLog,
		"gzipLevel":       location.GzipLevel,
		"gzipMinLength":   location.GzipMinLength,

		"fastcgi":     location.Fastcgi,
		"headers":     location.Headers,
		"cachePolicy": location.CachePolicy,
		"rewrite":     location.Rewrite,
		"websocket":   location.Websocket,
	}
	this.Data["server"] = server

	// 字符集
	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	this.Show()
}
