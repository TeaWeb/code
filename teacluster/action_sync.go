package teacluster

import (
	"github.com/TeaWeb/code/teacluster/configs"
	"github.com/TeaWeb/code/teahooks"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

// cluster -> node
type SyncAction struct {
	Action

	ItemActions []*configs.ItemAction
}

func (this *SyncAction) Name() string {
	return "sync"
}

func (this *SyncAction) Execute() error {
	for _, itemAction := range this.ItemActions {
		logs.Println("[cluster]"+itemAction.Action, "'"+itemAction.ItemId+"'")
		switch itemAction.Action {
		case configs.ItemActionAdd:
			fallthrough
		case configs.ItemActionChange:
			file := files.NewFile(Tea.ConfigFile(itemAction.ItemId))
			dir := file.Parent()
			if !dir.Exists() {
				err := dir.MkdirAll()
				if err != nil {
					logs.Error(err)
					return err
				}
			}
			file.Write(itemAction.Item.Data)
		case configs.ItemActionRemove:
			file := files.NewFile(Tea.ConfigFile(itemAction.ItemId))
			if file.Exists() {
				file.Delete()
			}
		}
	}

	ClusterManager.BuildSum()

	// reload system
	teahooks.Call(teahooks.EventReload)

	return nil
}

func (this *SyncAction) TypeId() int8 {
	return 8
}
