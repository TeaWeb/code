package cluster

import (
	"github.com/TeaWeb/code/teacluster"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

func (this *IndexAction) RunGet(params struct{}) {
	this.Data["teaMenu"] = "cluster"

	this.Data["node"] = teaconfigs.SharedNodeConfig()
	this.Data["isActive"] = teacluster.ClusterManager.IsActive()
	this.Data["error"] = teacluster.ClusterManager.Error()

	this.Show()
}
