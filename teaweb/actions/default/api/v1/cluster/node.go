package cluster

import (
	"github.com/TeaWeb/code/teacluster"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type NodeAction actions.Action

// 节点信息
func (this *NodeAction) RunGet(params struct{}) {
	config := teaconfigs.SharedNodeConfig()
	if config == nil {
		apiutils.Fail(this, "not a node yet")
		return
	}
	apiutils.Success(this, maps.Map{
		"isActive":  teacluster.SharedManager.IsActive(),
		"isChanged": teacluster.SharedManager.IsChanged(),
		"config":    config,
	})
}
