package teacluster

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

// cluster -> master|node
type SumAction struct {
	Action
}

func (this *SumAction) Name() string {
	return "sum"
}

func (this *SumAction) OnSuccess(success *SuccessAction) error {
	if success.Data == nil {
		return nil
	}

	sumMap := success.Data.Get("sum")
	if sumMap == nil || !types.IsMap(sumMap) {
		return nil
	}

	sumList := []string{}
	m := maps.NewMap(sumMap)
	for id, sum := range m {
		sumList = append(sumList, id+"|"+types.String(sum))
	}

	// write to local file
	file := files.NewFile(Tea.ConfigFile("cluster.sum"))
	file.WriteString(strings.Join(sumList, "\n"))

	// push or pull
	node := teaconfigs.SharedNodeConfig()
	if node == nil {
		return nil
	}

	if node.IsMaster() {
		ClusterManager.PushItems()
	} else {
		ClusterManager.PullItems()
	}

	return nil
}

func (this *SumAction) OnFail(fail *FailAction) error {
	// TODO retry later
	return nil
}

func (this *SumAction) TypeId() int8 {
	return 9
}
