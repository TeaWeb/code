package board

import (
	"github.com/TeaWeb/code/teastats"
	"github.com/iwind/TeaGo/actions"
)

type RefreshDataAction actions.Action

// 刷新数据
func (this *RefreshDataAction) RunPost(params struct {
	ServerId string
}) {
	queue := teastats.FindServerQueue(params.ServerId)
	if queue != nil {
		queue.Commit()
	}

	this.Success()
}
