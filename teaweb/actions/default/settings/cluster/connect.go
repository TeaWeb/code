package cluster

import (
	"github.com/TeaWeb/code/teacluster"
	"github.com/iwind/TeaGo/actions"
)

type ConnectAction actions.Action

func (this *ConnectAction) RunPost(params struct{}) {
	teacluster.ClusterManager.Restart()
	this.Success()
}
