package app

import (
	"fmt"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	AppId string
}) {
	plugin, app := FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到要查看的服务")
	}

	this.Data["app"] = app
	this.Data["cpu"] = app.SumCPUUsage()

	memoryUsage := app.SumMemoryUsage()
	this.Data["memoryRSS"] = fmt.Sprintf("%.2f", float64(memoryUsage.RSS)/1024/1024)
	this.Data["memoryVMS"] = fmt.Sprintf("%.2f", float64(memoryUsage.VMS)/1024/1024)
	this.Data["memory"] = memoryUsage

	this.Data["plugin"] = plugin.Name

	this.Success()
}
