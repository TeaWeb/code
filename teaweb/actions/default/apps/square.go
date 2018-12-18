package apps

import (
	"github.com/TeaWeb/jsapps/probes"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/types"
)

type SquareAction actions.Action

// 探针广场
func (this *SquareAction) Run(params struct{}) {
	var localProbeIds = []string{}
	{
		parser := probes.NewParser(Tea.ConfigFile("jsapps.js"))
		result, err := parser.Parse()
		if err == nil {
			for _, m := range result {
				localProbeIds = append(localProbeIds, types.String(m["id"]))
			}
		}
	}

	{
		parser := probes.NewParser(Tea.Root + Tea.DS + "plugins" + Tea.DS + "jsapps.js")
		result, err := parser.Parse()
		if err != nil {
			this.Data["error"] = err.Error()
			this.Data["probes"] = []map[string]interface{}{}
		} else {
			this.Data["error"] = ""

			for _, m := range result {
				m["isAdded"] = lists.Contains(localProbeIds, m["id"])
			}

			this.Data["probes"] = result
		}
	}

	this.Show()
}
