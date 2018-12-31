package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 路径规则列表
func (this *IndexAction) Run(params struct {
	Server string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["selectedTab"] = "location"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = server
	this.Data["server"] = maps.Map{
		"filename": server.Filename,
	}

	locations := []maps.Map{}
	for _, location := range server.Locations {
		location.Validate()
		locations = append(locations, maps.Map{
			"on":                location.On,
			"id":                location.Id,
			"type":              location.PatternType(),
			"pattern":           location.PatternString(),
			"patternTypeName":   teaconfigs.FindLocationPatternTypeName(location.PatternType()),
			"isCaseInsensitive": location.IsCaseInsensitive(),
			"isReverse":         location.IsReverse(),
			"rewrite":           location.Rewrite,
			"headers":           location.Headers,
			"fastcgi":           location.Fastcgi,
			"root":              location.Root,
			"cachePolicy":       location.CachePolicyObject(),
		})
	}

	this.Data["locations"] = locations

	this.Show()
}
