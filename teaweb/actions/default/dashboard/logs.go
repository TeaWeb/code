package dashboard

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type LogsAction actions.Action

// 实时日志
func (this *LogsAction) Run(params struct{}) {
	ones, err := tealogs.NewQuery().
		Desc("_id").
		Limit(10).
		Action(tealogs.QueryActionFindAll).
		Execute()
	if err != nil {
		logs.Error(err)
		this.Data["logs"] = []*tealogs.AccessLog{}
	} else {
		this.Data["logs"] = ones
	}

	this.Data["qps"] = tealogs.SharedLogger().QPS()

	this.Success()
}
