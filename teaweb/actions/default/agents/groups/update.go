package groups

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction actions.Action

// 分组ID
func (this *UpdateAction) Run(params struct {
	GroupId string
}) {
	if params.GroupId == "default" {
		this.Data["group"] = agents.LoadDefaultGroup()
	} else {
		group := agents.SharedGroupConfig().FindGroup(params.GroupId)
		if group == nil {
			this.Fail("找不到Group")
		}
		this.Data["group"] = group
	}

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	GroupId string
	Name    string
	Must    *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")

	if params.GroupId == "default" {
		group := agents.LoadDefaultGroup()
		group.Name = params.Name
		err := group.WriteToFile(agents.DefaultGroupFile())
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	} else {
		config := agents.SharedGroupConfig()
		group := config.FindGroup(params.GroupId)
		if group == nil {
			this.Fail("找不到Group")
		}
		group.Name = params.Name
		err := config.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	}

	this.Success()
}
