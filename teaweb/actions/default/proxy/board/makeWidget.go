package board

import "github.com/iwind/TeaGo/actions"

type MakeWidgetAction actions.Action

// 制作Widget
func (this *MakeWidgetAction) Run(params struct{}) {
	this.Data["code"] = `var widget = {
	"name": "", // Widget名称
	"code": "", // Widget代号
	"author": "", // 作者
	"version": "", // 版本
};

widget.run = function () {
	var chart = new charts.HTMLChart();
	chart.options.name = ""; // Chart名称;
	chart.options.columns = 1; // Chart宽度
	chart.html = ""; // Chart HTML内容，只对HTMLChart有效
	chart.render();
}`

	this.Show()
}
