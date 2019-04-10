package dashboard

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
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
		if err != topology.ErrServerSelectionTimeout {
			logs.Error(err)
		}
		this.Data["logs"] = []*tealogs.AccessLog{}
	} else {
		this.Data["logs"] = ones
	}

	this.Data["qps"] = teaproxy.QPS

	this.Success()
}
