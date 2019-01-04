package locations

import (
	"github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/locations/locationutils"
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
	_, location := locationutils.SetCommonInfo(this, params.Server, params.LocationId, "cache")

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
