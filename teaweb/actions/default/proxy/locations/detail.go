package locations

import (
	"github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type DetailAction actions.Action

// 路径规则详情
func (this *DetailAction) Run(params struct {
	Server string
	Index  int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["server"] = maps.Map{
		"filename": proxy.Filename,
	}

	location := proxy.LocationAtIndex(params.Index)
	if location == nil {
		this.Fail("找不到要修改的路径配置")
	}

	if location.Index == nil {
		location.Index = []string{}
	}

	this.Data["selectedTab"] = "location"
	this.Data["filename"] = params.Server
	this.Data["locationIndex"] = params.Index

	this.Data["location"] = maps.Map{
		"on":              location.On,
		"type":            location.PatternType(),
		"pattern":         location.PatternString(),
		"caseInsensitive": location.IsCaseInsensitive(),
		"reverse":         location.IsReverse(),
		"root":            location.Root,
		"rewrite": lists.Map(location.Rewrite, func(k int, v interface{}) interface{} {
			r := v.(*teaconfigs.RewriteRule)
			return maps.Map{
				"id":             r.Id,
				"on":             r.On,
				"headers":        r.Headers,
				"ignoreHeaders":  r.IgnoreHeaders,
				"replace":        r.Replace,
				"pattern":        r.Pattern,
				"cond":           r.Cond,
				"redirectMethod": r.RedirectMethod(),
			}
		}),
		"fastcgi": location.FastcgiAtIndex(0),
		"charset": location.Charset,
		"index":   location.Index,
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

	// headers
	this.Data["locationIndex"] = params.Index
	this.Data["headers"] = location.Headers
	this.Data["ignoreHeaders"] = lists.NewList(location.IgnoreHeaders).Map(func(k int, v interface{}) interface{} {
		return map[string]interface{}{
			"name": v,
		}
	}).Slice

	// 缓存策略
	this.Data["cachePolicy"] = ""
	if len(location.CachePolicy) > 0 {
		policy := shared.NewCachePolicyFromFile(location.CachePolicy)
		if policy != nil {
			this.Data["cachePolicy"] = policy.Name + "（" + teacache.TypeName(policy.Type) + "）"
		}
	}
	this.Data["cachePolicyFile"] = location.CachePolicy

	cache, _ := teaconfigs.SharedCacheConfig()
	this.Data["cachePolicyList"] = lists.Map(cache.FindAllPolicies(), func(k int, v interface{}) interface{} {
		policy := v.(*shared.CachePolicy)
		return maps.Map{
			"filename": policy.Filename,
			"name":     policy.Name,
			"type":     teacache.TypeName(policy.Type),
		}
	})

	this.Show()
}
