package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type DetailAction actions.Action

func (this *DetailAction) Run(params struct {
	Filename string
	Index    int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location == nil {
		this.Fail("找不到要修改的路径配置")
	}

	if location.Index == nil {
		location.Index = []string{}
	}

	this.Data["filename"] = params.Filename
	this.Data["locationIndex"] = params.Index
	this.Data["location"] = maps.Map{
		"on":              location.On,
		"type":            location.PatternType(),
		"pattern":         location.PatternString(),
		"caseInsensitive": location.IsCaseInsensitive(),
		"reverse":         location.IsReverse(),
		"root":            location.Root,
		"rewrite":         location.Rewrite,
		"fastcgi":         location.FastcgiAtIndex(0),
		"charset":         location.Charset,
		"index":           location.Index,
	}
	this.Data["proxy"] = proxy

	// 已经有的代理服务
	proxyConfigs := teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir())
	proxies := []maps.Map{}
	for _, proxyConfig := range proxyConfigs {
		if proxyConfig.Id == proxy.Id {
			continue
		}

		name := proxyConfig.Description
		if !proxyConfig.On {
			name += "(未启用)"
		}
		proxies = append(proxies, maps.Map{
			"id":   proxyConfig.Id,
			"name": name,
		})
	}
	this.Data["proxies"] = proxies

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

	// 字符集
	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	this.Show()
}
