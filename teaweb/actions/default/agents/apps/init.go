package apps

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAgent,
			}).
			Prefix("/agents/apps").
			Helper(new(Helper)).
			Get("", new(IndexAction)).
			GetPost("/add", new(AddAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/delete", new(DeleteAction)).
			Get("/detail", new(DetailAction)).
			Get("/schedule", new(ScheduleAction)).
			Get("/boot", new(BootAction)).
			Get("/manual", new(ManualAction)).
			GetPost("/addTask", new(AddTaskAction)).
			Post("/deleteTask", new(DeleteTaskAction)).
			Get("/taskDetail", new(TaskDetailAction)).
			GetPost("/updateTask", new(UpdateTaskAction)).
			Get("/runTask", new(RunTaskAction)).
			GetPost("/taskLogs", new(TaskLogsAction)).
			GetPost("/monitor", new(MonitorAction)).
			GetPost("/addItem", new(AddItemAction)).
			Post("/deleteItem", new(DeleteItemAction)).
			Get("/itemDetail", new(ItemDetailAction)).
			GetPost("/updateItem", new(UpdateItemAction)).
			GetPost("/itemValues", new(ItemValuesAction)).
			GetPost("/itemCharts", new(ItemChartsAction)).
			GetPost("/addItemChart", new(AddItemChartAction)).
			Post("/deleteItemChart", new(DeleteItemChartAction)).
			GetPost("/updateItemChart", new(UpdateItemChartAction)).
			Get("/widget", new(WidgetAction)).
			GetPost("/addWidget", new(AddWidgetAction)).
			GetPost("/makeWidget", new(MakeWidgetAction)).
			Post("/testWidget", new(TestWidgetAction)).
			EndAll()
	})
}
