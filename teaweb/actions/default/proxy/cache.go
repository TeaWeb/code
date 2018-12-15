package proxy

import (
	"github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type CacheAction actions.Action

func (this *CacheAction) Run(params struct {
	Filename string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["server"] = maps.Map{
		"filename": params.Filename,
	}

	this.Data["selectedTab"] = "cache"
	this.Data["filename"] = params.Filename
	this.Data["proxy"] = server

	// 缓存策略
	this.Data["cachePolicy"] = ""
	if len(server.CachePolicy) > 0 {
		policy := teaconfigs.NewCachePolicyFromFile(server.CachePolicy)
		if policy != nil {
			this.Data["cachePolicy"] = policy.Name + "（" + teacache.TypeName(policy.Type) + "）"
		}
	}
	this.Data["cachePolicyFile"] = server.CachePolicy

	cache, _ := teaconfigs.SharedCacheConfig()
	this.Data["cachePolicyList"] = lists.Map(cache.FindAllPolicies(), func(k int, v interface{}) interface{} {
		policy := v.(*teaconfigs.CachePolicy)
		return maps.Map{
			"filename": policy.Filename,
			"name":     policy.Name,
			"type":     teacache.TypeName(policy.Type),
		}
	})

	this.Show()
}
