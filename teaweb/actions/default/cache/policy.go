package cache

import (
	"github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type PolicyAction actions.Action

// 缓存策略详情
func (this *PolicyAction) Run(params struct {
	Filename string
}) {
	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到Policy")
	}

	this.Data["policy"] = policy

	// 类型
	this.Data["typeName"] = teacache.TypeName(policy.Type)

	this.Show()
}
