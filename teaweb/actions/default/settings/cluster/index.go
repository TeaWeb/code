package cluster

import (
	"github.com/TeaWeb/code/teacluster"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

func (this *IndexAction) RunGet(params struct{}) {
	this.Data["teaMenu"] = "cluster"

	node := teaconfigs.SharedNodeConfig()

	manager := teacluster.ClusterManager

	if node != nil && node.On {
		manager.BuildSum()
	}

	this.Data["node"] = node
	this.Data["isActive"] = manager.IsActive()
	this.Data["error"] = manager.Error()

	if node != nil && node.IsMaster() {
		this.Data["isChanged"] = manager.IsChanged()
	} else {
		this.Data["isChanged"] = false
	}

	this.Show()
}
