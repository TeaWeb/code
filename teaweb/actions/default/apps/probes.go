package apps

import (
	"github.com/TeaWeb/jsapps/probes"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
)

type ProbesAction actions.Action

// 探针列表
func (this *ProbesAction) Run(params struct{}) {
	parser := probes.NewParser(Tea.ConfigFile("jsapps.js"))
	result, err := parser.Parse()
	if err != nil {
		this.Data["error"] = err.Error()
		this.Data["probes"] = []map[string]interface{}{}
	} else {
		this.Data["error"] = ""
		this.Data["probes"] = result
	}

	this.Show()
}
