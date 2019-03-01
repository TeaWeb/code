package agents

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type MoveGroupAction actions.Action

// 移动分组位置
func (this *MoveGroupAction) Run(params struct {
	FromIndex int
	ToIndex   int
}) {
	config := agents.SharedGroupConfig()
	config.Move(params.FromIndex, params.ToIndex)
	err := config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
