package locations

import (
	"github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateCacheAction actions.Action

// 更新缓存设置
func (this *UpdateCacheAction) Run(params struct {
	Server        string
	LocationIndex int
	Policy        string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail("找不到要操作的代理服务")
	}

	location := server.LocationAtIndex(params.LocationIndex)
	if location == nil {
		this.Fail("找不到要操作的路径规则")
	}

	if len(params.Policy) > 0 {
		policy := shared.NewCachePolicyFromFile(params.Policy)
		if policy == nil {
			this.Fail("找不到要使用的缓存策略")
		}
		this.Data["policy"] = maps.Map{
			"name":     policy.Name,
			"typeName": teacache.TypeName(policy.Type),
		}
	} else {
		this.Data["policy"] = maps.Map{
			"name":     "",
			"typeName": "",
		}
	}

	location.CacheOn = true
	location.CachePolicy = params.Policy
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
