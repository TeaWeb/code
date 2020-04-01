package cache

import (
	"github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaweb/actions/default/actionutils"
)

type CleanPolicyAction struct {
	actionutils.ParentAction
}

// 清理
func (this *CleanPolicyAction) Run(params struct {
	Filename string
}) {
	this.SecondMenu("clean")

	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Data["result"] = "找不到Policy"
		this.Fail()
	}
	this.Data["policy"] = policy
	this.Show()
}

// 执行清理
func (this *CleanPolicyAction) RunPost(params struct {
	Filename string
}) {
	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Data["result"] = "找不到Policy"
		this.Fail()
	}

	manager := teacache.FindCachePolicyManager(params.Filename)
	if manager == nil {
		manager = teacache.NewManagerFromConfig(policy)
		defer func() {
			_ = manager.Close()
		}()
	}

	if manager == nil {
		this.Fail("找不到管理器")
	}

	err := manager.Clean()
	if err != nil {
		this.Data["result"] = err.Error()
		this.Fail()
	}

	this.Data["result"] = "清理完成"
	this.Success()
}
