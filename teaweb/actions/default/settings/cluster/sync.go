package cluster

import (
	"github.com/TeaWeb/code/teacluster"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type SyncAction actions.Action

// 同步
func (this *SyncAction) RunPost(params struct{}) {
	node := teaconfigs.SharedNodeConfig()
	if node == nil {
		this.Fail("节点配置不存在")
	}

	if !teacluster.ClusterManager.IsActive() {
		this.Fail("当前节点没有连接到集群服务器")
	}

	if node.IsMaster() {
		teacluster.ClusterManager.PushItems()
	} else {
		teacluster.ClusterManager.PullItems()
	}

	teacluster.ClusterManager.SetIsChanged(false)

	this.Success()
}
