package board

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/TeaWeb/code/teastats"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type MakeAction actions.Action

// 制作图表
func (this *MakeAction) Run(params struct {
	Server string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail("找不到要查看的代理服务")
	}

	this.Data["server"] = maps.Map{
		"id":       server.Id,
		"filename": params.Server,
	}

	this.Data["items"] = teastats.FindAllStatFilters()

	this.Show()
}

// 保存提交
func (this *MakeAction) RunPost(params struct {
	Server         string
	Name           string
	Description    string
	Columns        uint8
	Items          []string
	JavascriptCode string
	On             bool
	Must           *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入名称")

	_, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail("找不到要查看的代理服务")
	}

	chart := widgets.NewChart()
	chart.On = params.On
	chart.Name = params.Name
	chart.Description = params.Description
	chart.Columns = params.Columns
	chart.Requirements = params.Items
	chart.Type = "javascript"
	chart.Options = maps.Map{
		"code": params.JavascriptCode,
	}

	widget := widgets.NewWidget()
	widget.AddChart(chart)
	err = widget.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
