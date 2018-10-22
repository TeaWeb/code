package teacharts

import (
	"github.com/TeaWeb/plugin/charts"
)

func NewLineChart() *LineChart {
	p := &LineChart{}
	p.Type = "line"
	p.Lines = []*Line{}
	return p
}

func NewLineChartFromInterface(chart *charts.LineChart) *LineChart {
	p := &LineChart{}
	p.Type = "line"
	p.Name = chart.Name
	p.Detail = chart.Detail

	p.Lines = []*Line{}

	for _, line := range chart.Lines {
		p.AddLine(&Line{
			Name:      line.Name,
			Values:    line.Values,
			Color:     line.Color,
			Filled:    line.Filled,
			ShowItems: line.ShowItems,
		})
	}

	p.Labels = chart.Labels
	p.Max = chart.Max
	p.XShowTick = chart.XShowTick
	p.YTickCount = chart.YTickCount
	p.YShowTick = chart.YShowTick
	return p
}

type Line struct {
	Name      string        `json:"name"`
	Values    []interface{} `json:"values"`
	Color     Color         `json:"color"`
	Filled    bool          `json:"filled"`
	ShowItems bool          `json:"showItems"`
}

type LineChart struct {
	Chart

	Lines  []*Line  `json:"lines"`
	Labels []string `json:"labels"`

	Max       float64 `json:"max"`
	XShowTick bool    `json:"xShowTick"` // X轴是否显示刻度

	YTickCount uint `json:"yTickCount"` // Y轴刻度分隔数量
	YShowTick  bool `json:"yShowTick"`  // Y轴是否显示刻度
}

func (this *LineChart) AddLine(line *Line) {
	this.Lines = append(this.Lines, line)
}

func (this *LineChart) ResetLines() {
	this.Lines = []*Line{}
}
