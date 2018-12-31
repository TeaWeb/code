package locations

import (
	"github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type CacheAction actions.Action

// 缓存设置
func (this *CacheAction) Run(params struct {
	Server     string
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

	this.Data["selectedTab"] = "location"
	this.Data["selectedSubTab"] = "cache"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = server
	this.Data["server"] = maps.Map{
		"filename": server.Filename,
	}

	this.Show()
}
