package log

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teamongo"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
}) {
	this.Data["teaMenu"] = "log"

	// 检查MongoDB连接
	this.Data["mongoError"] = ""
	err := teamongo.Test()
	if err != nil {
		this.Data["mongoError"] = "此功能需要连接MongoDB"
	}

	this.Show()
}
