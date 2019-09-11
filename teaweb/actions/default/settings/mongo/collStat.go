package mongo

import (
	"github.com/TeaWeb/code/teadb"
	"github.com/iwind/TeaGo/actions"
)

type CollStatAction actions.Action

// 集合统计
func (this *CollStatAction) Run(params struct {
	CollNames []string
}) {
	statMap, err := teadb.SharedDB().StatTables(params.CollNames)
	if err != nil {
		this.Data["result"] = map[string]interface{}{}
		this.Fail("获取统计信息失败：" + err.Error())
	} else {
		this.Data["result"] = statMap
	}

	this.Success()
}
