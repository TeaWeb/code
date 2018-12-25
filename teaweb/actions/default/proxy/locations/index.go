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
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["selectedTab"] = "location"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = proxy

	this.Data["typeOptions"] = []maps.Map{
		{
			"name":  "匹配前缀",
			"value": teaconfigs.LocationPatternTypePrefix,
		},
		{
			"name":  "精准匹配",
			"value": teaconfigs.LocationPatternTypeExact,
		},
		{
			"name":  "正则表达式匹配",
			"value": teaconfigs.LocationPatternTypeRegexp,
		},
	}

	locations := []maps.Map{}
	for _, location := range proxy.Locations {
		location.Validate()
		locations = append(locations, maps.Map{
			"on":              location.On,
			"type":            location.PatternType(),
			"pattern":         location.PatternString(),
			"caseInsensitive": location.IsCaseInsensitive(),
			"reverse":         location.IsReverse(),
		})
	}

	this.Data["locations"] = locations

	this.Show()
}
