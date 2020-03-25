package settings

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

func (this *IndexAction) RunGet(params struct{}) {
	this.Data["setting"] = teaconfigs.LoadProxySetting()

	this.Show()
}
