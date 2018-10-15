package teacharts

import "github.com/TeaWeb/code/teainterfaces"

func NewPieChart() *PieChart {
	p := &PieChart{
		Values: []interface{}{},
		Labels: []string{},
	}
	p.Type = "pie"
	return p
}

func NewPieChartFromInterface(chart teainterfaces.PieChartInterface) *PieChart {
	p := &PieChart{
		Values: chart.Values(),
		Labels: chart.Labels(),
	}
	p.Type = "pie"
	p.Name = chart.(teainterfaces.ChartInterface).Name()
	p.Detail = chart.(teainterfaces.ChartInterface).Detail()
	return p
}

type PieChart struct {
	Chart
	Values []interface{} `json:"values"`
	Labels []string      `json:"labels"`
}

func (this *PieChart) UniqueId() string {
	return this.Id
}

func (this *PieChart) SetUniqueId(id string) {
	this.Id = id
}
