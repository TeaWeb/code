package board

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 定义路由
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Module("").
			Prefix("/proxy/board").
			GetPost("", new(IndexAction)).
			GetPost("/charts", new(ChartsAction)).
			GetPost("/make", new(MakeAction)).
			Post("/test", new(TestAction)).
			GetPost("/chart", new(ChartAction)).
			GetPost("/updateChart", new(UpdateChartAction)).
			Post("/deleteChart", new(DeleteChartAction)).
			Post("/useChart", new(UseChartAction)).
			Post("/cancelChart", new(CancelChartAction)).
			Post("/moveChart", new(MoveChartAction)).
			EndAll()

		// 检查图表更新
		logs.Println("[proxy]check widget changes")
		if proxyutils.CheckChartChanges() {
			err := proxyutils.ApplyChartChanges()
			if err != nil {
				logs.Error(err)
			}
		}
	})
}
