package app

import (
	"fmt"
	"github.com/TeaWeb/code/teaweb/actions/default/apputils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type ReloadAction actions.Action

// 刷新App状态
func (this *ReloadAction) Run(params struct {
	AppId string
}) {
	plugin, app := apputils.FindApp(params.AppId)
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
			"cpuPercent":       fmt.Sprintf("%.1f", app.SumCPUUsage().Percent),
			"memoryPercent":    fmt.Sprintf("%.1f", memory.Percent),
			"memoryRSS":        fmt.Sprintf("%.2f", float64(memory.RSS)/1024/1024), // 单位：M
			"countProcesses":   len(app.Processes),
			"countConnections": app.CountAllConnections(),
			"countOpenFiles":   app.CountAllOpenFiles(),
			"countListens":     app.CountAllListens(),
			"isFavored":        apputils.FavorAppContains(app.UniqueId()),
		}
		this.Success()
	}

	this.Fail()
}
