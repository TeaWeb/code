package groups

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type DetailAction actions.Action

// 详情
func (this *DetailAction) Run(params struct {
	GroupId string
}) {
	if len(params.GroupId) == 0 {
		this.Data["group"] = maps.Map{
			"id":   "",
			"name": "默认分组",
			"on":   true,
		}
	} else {
		group := agents.SharedGroupConfig().FindGroup(params.GroupId)
		if group == nil {
			this.Fail("找不到Group")
		}
		this.Data["group"] = group
	}

	this.Show()
}
