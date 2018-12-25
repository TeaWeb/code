package proxy

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
	Server string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["server"] = maps.Map{
		"filename": params.Server,
	}

	this.Data["selectedTab"] = "cache"
	this.Data["filename"] = params.Server
	this.Data["proxy"] = server

	// 缓存策略
	this.Data["cachePolicy"] = ""
	if len(server.CachePolicy) > 0 {
		policy := shared.NewCachePolicyFromFile(server.CachePolicy)
		if policy != nil {
			this.Data["cachePolicy"] = policy.Name + "（" + teacache.TypeName(policy.Type) + "）"
		}
	}
	this.Data["cachePolicyFile"] = server.CachePolicy

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
