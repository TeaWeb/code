package agents

import (
	"github.com/TeaWeb/jsapps/probes"
	"github.com/TeaWeb/plugin/apps"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

type ProbeScriptAction actions.Action

// 探针脚本
func (this *ProbeScriptAction) Run(params struct {
	FromList bool
	ProbeId  string
}) {
	if len(params.ProbeId) == 0 {
		this.Fail("请选择要查看的探针")
	}

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

	this.Data["isAdded"] = lists.Contains(localProbeIds, params.ProbeId)
	this.Data["fromList"] = params.FromList
	this.Data["probeId"] = params.ProbeId

	if strings.HasPrefix(params.ProbeId, "local_") {
		parser := probes.NewParser(Tea.ConfigFile("jsapps.js"))
		f, _ := parser.FindProbeFunction(params.ProbeId)
		this.Data["func"] = f
		this.Data["isLocal"] = true
	} else {
		parser := probes.NewParser(Tea.Root + Tea.DS + "plugins" + Tea.DS + "jsapps.js")
		f, _ := parser.FindProbeFunction(params.ProbeId)

		this.Data["func"] = f
		this.Data["isLocal"] = false
	}

	this.Show()
}

func (this *ProbeScriptAction) RunPost(params struct {
	ProbeId   string
	Script    string
	IsTesting bool
}) {
	if params.IsTesting {
		engine := probes.NewScriptEngine()
		err := engine.RunScript("(" + params.Script + ")()")
		if err != nil {
			this.Fail("未找到匹配的App：" + err.Error())
		}

		results := engine.Apps()
		this.Data["apps"] = lists.Map(results, func(k int, v interface{}) interface{} {
			app := v.(*apps.App)
			return maps.Map{
				"name":      app.Name,
				"developer": app.Developer,
				"site":      app.Site,
				"docSite":   app.DocSite,
				"version":   app.Version,
				"file":      app.Processes[0].File,
				"dir":       app.Processes[0].Dir,
				"processes": lists.Map(app.Processes, func(k int, v interface{}) interface{} {
					var process = v.(*apps.Process)
					return maps.Map{
						"pid":     process.Pid,
						"cmdline": process.Cmdline,
					}
				}),
			}
		})

		this.Success()
	} else {
		parser := probes.NewParser(Tea.ConfigFile("jsapps.js"))
		err := parser.ReplaceProbe(params.ProbeId, params.Script)
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}

		this.Success()
	}
}
