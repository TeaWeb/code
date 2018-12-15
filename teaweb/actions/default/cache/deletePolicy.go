package cache

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type DeletePolicyAction actions.Action

// 删除缓存策略
func (this *DeletePolicyAction) Run(params struct {
	Filename string
}) {
	if len(params.Filename) == 0 {
		this.Fail("请指定要删除的缓存策略")
	}

	policy := teaconfigs.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到要删除的缓存策略")
	}

	config, _ := teaconfigs.SharedCacheConfig()
	config.DeletePolicy(params.Filename)
	err := config.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	err = policy.Delete()
	if err != nil {
		logs.Error(err)
	}

	this.Success()
}
