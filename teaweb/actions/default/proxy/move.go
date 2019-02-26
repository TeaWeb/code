package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type MoveAction actions.Action

// 移动代理服务位置
func (this *MoveAction) Run(params struct {
	FromIndex int
	ToIndex   int
}) {
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	serverList.MoveServer(params.FromIndex, params.ToIndex)
	err = serverList.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
