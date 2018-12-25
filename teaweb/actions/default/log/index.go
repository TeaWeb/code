package log

import (
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	Server string
}) {
	// 检查MongoDB连接
	this.Data["mongoError"] = ""
	err := teamongo.Test()
	if err != nil {
		this.Data["mongoError"] = "此功能需要连接MongoDB"
	}

	this.Data["server"] = maps.Map{
		"filename": params.Server,
	}

	this.Show()
}
