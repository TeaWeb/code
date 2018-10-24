package apps

import (
	"fmt"
	"github.com/TeaWeb/code/teaapps"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/TeaWeb/code/teaweb/actions/default/apputils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AllAction actions.Action

func (this *AllAction) Run(params struct{}) {
	allApps := []maps.Map{}

	for _, plugin := range teaplugins.Plugins() {
		apps := plugin.Apps
		if len(apps) == 0 {
			continue
		}

		for _, app := range apps {
			memory := app.SumMemoryUsage()

			func(app *teaapps.App) {
				go app.Reload()
			}(app)

			isFavored := apputils.FavorAppContains(app.UniqueId())
			allApps = append(allApps, maps.Map{
				"plugin": plugin.Name,

				"id":               app.Id,
				"isFavored":        isFavored,
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
			})
		}
	}

	this.Data["apps"] = allApps

	this.Show()
}
