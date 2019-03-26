package agents

import (
	"fmt"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/TeaWeb/code/teaweb/actions/default/apputils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type WatchAction actions.Action

// 监控是否变化
func (this *WatchAction) Run(params struct {
	AppIds []string
}) {
	allAppIds := []string{}
	for _, plugin := range teaplugins.Plugins() {
		for _, app := range plugin.Apps {
			allAppIds = append(allAppIds, app.Id)
		}
	}

	isChanged := false
	for _, appId := range params.AppIds {
		if !lists.ContainsString(allAppIds, appId) {
			isChanged = true
			break
		}
	}
	if !isChanged {
		for _, appId := range allAppIds {
			if !lists.ContainsString(params.AppIds, appId) {
				isChanged = true
				break
			}
		}
	}

	this.Data["isChanged"] = isChanged

	// 所有Apps
	if isChanged {
		allApps := []maps.Map{}

		for _, plugin := range teaplugins.Plugins() {
			apps := plugin.Apps
			if len(apps) == 0 {
				continue
			}

			for _, app := range apps {
				memory := app.SumMemoryUsage()

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
	}

	this.Success()
}
