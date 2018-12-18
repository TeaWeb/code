package apps

import (
	"fmt"
	"github.com/TeaWeb/code/teacharts"
	"github.com/TeaWeb/code/teaplugins"
	_ "github.com/TeaWeb/code/teaweb/actions/default/apps/app"
	"github.com/TeaWeb/code/teaweb/actions/default/apputils"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantApp,
			}).
			Helper(new(Helper)).
			Prefix("/apps").
			Get("", new(IndexAction)).
			Get("/all", new(AllAction)).
			Post("/watch", new(WatchAction)).
			Get("/probes", new(ProbesAction)).
			Get("/square", new(SquareAction)).
			GetPost("/addProbe", new(AddProbeAction)).
			Post("/deleteProbe", new(DeleteProbeAction)).
			GetPost("/probeScript", new(ProbeScriptAction)).
			Post("/copyProbe", new(CopyProbeAction)).
			EndAll()
	})

	// 增加widget
	addWidgets()
}

func addWidgets() {
	p := teaplugins.NewPlugin()

	widget := teaplugins.NewWidget()
	widget.MoreURL = "/apps"
	widget.Dashboard = true
	widget.Group = teaplugins.WidgetGroupService
	widget.Name = "Apps"
	widget.OnReload(func() {
		widget.ResetCharts()

		// 判断是否有关注的App
		hasFollowedApps := false
		for _, p1 := range teaplugins.Plugins() {
			for _, a := range p1.Apps {
				if apputils.FavorAppContains(a.UniqueId()) {
					hasFollowedApps = true
					break
				}
			}
		}

		for _, p1 := range teaplugins.Plugins() {
			for _, a := range p1.Apps {
				if !hasFollowedApps || apputils.FavorAppContains(a.UniqueId()) {
					c := teacharts.NewTable()
					c.Name = a.Name
					c.AddRow("<strong>" + a.Name + "</strong>")
					if a.IsRunning {
						if len(a.Processes) > 0 {
							c.AddRow(fmt.Sprintf("<span title=\"PID\">%d</span>", a.Processes[0].Pid), "正在运行", "<i class=\"ui icon dot circle green small\"></i>")
						} else {
							c.AddRow(" ", "正在运行", "<i class=\"ui icon dot circle green small\"></i>")
						}
					} else {
						c.AddRow("", "已停止", "<i class=\"ui icon dot circle grey small\"></i>")
					}
					c.SetWidth(35, 45, 20)

					widget.AddChart(c)
				}
			}
		}
	})
	p.AddWidget(widget)

	teaplugins.Register(p)
}
