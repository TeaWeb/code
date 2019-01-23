package agents

import (
	"fmt"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/TeaWeb/code/teaweb/actions/default/apputils"
	"github.com/TeaWeb/jsapps/probes"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"os/exec"
	"strings"
)

type AllAction actions.Action

// 所有App列表
func (this *AllAction) Run(params struct{}) {
	allApps := []maps.Map{}

	// 检查命令 ps, pgrep, lsof
	messages := []string{}
	for _, cmd := range []string{"ps", "pgrep", "lsof"} {
		_, err := exec.LookPath(cmd)
		if err != nil {
			messages = append(messages, "需要先安装命令\""+cmd+"\"，<a href=\"https://github.com/TeaWeb/build/blob/master/docs/apps/Install"+strings.ToUpper(cmd[0:1])+cmd[1:]+".md\" target=\"_blank\">如何安装？</a>")
		}
	}
	this.Data["messages"] = messages

	// 探针
	parser := probes.NewParser(Tea.Root + Tea.DS + "plugins" + Tea.DS + "jsapps.js")
	_, f, err := parser.LoadFunctions()
	if err != nil {
		this.Data["countProbes"] = 0
	} else {
		this.Data["countProbes"] = len(f.Keys())
	}

	// 刷新
	teaplugins.ReloadAllApps()

	// 所有Apps
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

	this.Show()
}
