package app

import (
	"fmt"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type ReloadAction actions.Action

func (this *ReloadAction) Run(params struct {
	AppId string
}) {
	plugin, app := FindApp(params.AppId)
	if app != nil {
		app.Reload()

		memory := app.SumMemoryUsage()

		this.Data["app"] = maps.Map{
			"plugin": plugin.Name,

			"id":               app.Id,
			"name":             app.Name,
			"site":             app.Site,
			"docSite":          app.DocSite,
			"isRunning":        app.IsRunning,
			"cpuPercent":       app.SumCPUUsage().Percent * 100,
			"memoryPercent":    memory.Percent * 100,
			"memoryRSS":        fmt.Sprintf("%.2f", float64(memory.RSS)/1024/1024), // 单位：M
			"countProcesses":   len(app.Processes),
			"countConnections": app.CountAllConnections(),
			"countOpenFiles":   app.CountAllOpenFiles(),
			"countListens":     app.CountAllListens(),
		}
		this.Success()
	}

	this.Fail()
}
