package groups

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type AddAction actions.Action

// 添加分组
func (this *AddAction) Run(params struct {
	From string
}) {
	this.Data["from"] = params.From

	this.Show()
}

// 提交保存
func (this *AddAction) RunPost(params struct {
	Name string
	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")

	group := agents.NewGroup(params.Name)
	config := agents.SharedGroupConfig()
	config.AddGroup(group)
	err := config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
