package apps

import (
	probes2 "github.com/TeaWeb/jsapps/probes"
	"github.com/TeaWeb/plugin/apps"
	"github.com/TeaWeb/plugin/apps/probes"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type AddProbeAction actions.Action

// 添加规则
func (this *AddProbeAction) Run(params struct{}) {
	this.Show()
}

func (this *AddProbeAction) RunPost(params struct {
	Name            string
	Developer       string
	Site            string
	DocSite         string
	CommandName     string
	CommandPatterns []string
	CommandVersion  string
	IsTesting       bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入App名称").
		Field("commandName", params.CommandName).
		Require("请输入启动的命令名")

	probe := probes.NewProcessProbe()
	probe.Name = params.Name
	probe.Developer = params.Developer
	probe.Site = params.Site
	probe.DocSite = params.DocSite
	probe.CommandName = params.CommandName
	probe.CommandPatterns = params.CommandPatterns
	probe.CommandVersion = params.CommandVersion

	if params.IsTesting { // 测试
		results, err := probe.Run()
		if err != nil {
			logs.Error(err)
			this.Fail("未找到匹配的App")
		}

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
	} else { // 保存
		if len(params.CommandPatterns) == 0 {
			params.CommandPatterns = []string{}
		}

		probe := probes.NewProcessProbe()
		probe.Name = params.Name
		probe.Developer = params.Developer
		probe.Site = params.Site
		probe.DocSite = params.DocSite
		probe.CommandName = params.CommandName
		probe.CommandPatterns = params.CommandPatterns
		probe.CommandVersion = params.CommandVersion

		parser := probes2.NewParser(Tea.ConfigFile("jsapps.js"))
		err := parser.AddProbe(probe)
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
		this.Next("/apps/probes", nil)
		this.Success("保存成功")
	}
}
