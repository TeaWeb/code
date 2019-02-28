package agent

import (
	"encoding/json"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"net/http"
)

type ItemAction actions.Action

// 监控项
func (this *ItemAction) Run(params struct {
	AgentId string
	ItemId  string
}) {
	apiutils.ValidateUser(this)

	this.AddHeader("Content-Type", "application/json; charset=utf-8")

	// 获取数据
	query := teamongo.NewAgentValueQuery()
	query.Agent(params.AgentId)
	query.Item(params.ItemId)
	query.Desc("_id")
	v, err := query.Find()
	if err != nil {
		logs.Error(err)
		this.Error(err.Error(), http.StatusInternalServerError)
	} else if v == nil {
		this.Error("item value not found", http.StatusNotFound)
	} else {
		data, err := json.Marshal(v.Value)
		if err != nil {
			logs.Error(err)
			this.Error(err.Error(), http.StatusInternalServerError)
		} else {
			this.Write(data)
		}
	}
}
