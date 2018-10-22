package teacharts

import (
	"github.com/TeaWeb/plugin/charts"
)

func NewPieChart() *PieChart {
	p := &PieChart{
		Values: []interface{}{},
		Labels: []string{},
	}
	p.Type = "pie"
	return p
}

func NewPieChartFromInterface(chart *charts.PieChart) *PieChart {
	p := &PieChart{
		Values: chart.Values,
		Labels: chart.Labels,
	}
	p.Type = "pie"
	p.Name = chart.Name
	p.Detail = chart.Detail
	return p
}

type PieChart struct {
	Chart
	Values []interface{} `json:"values"`
	Labels []string      `json:"labels"`
}
