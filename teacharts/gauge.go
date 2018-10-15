package teacharts

import "github.com/TeaWeb/code/teainterfaces"

func NewGaugeChart() *GaugeChart {
	p := &GaugeChart{}
	p.Type = "gauge"
	return p
}

func NewGaugeChartFromInterface(chart teainterfaces.GaugeChartInterface) *GaugeChart {
	p := &GaugeChart{}
	p.Type = "gauge"
	p.Name = chart.(teainterfaces.ChartInterface).Name()
	p.Detail = chart.(teainterfaces.ChartInterface).Detail()

	p.Value = chart.Value()
	p.Label = chart.Label()
	p.Min = chart.Min()
	p.Max = chart.Max()
	p.Unit = chart.Unit()
	return p
}

type GaugeChart struct {
	Chart

	Value float64 `json:"value"`
	Label string  `json:"label"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Unit  string  `json:"unit"`
}

func (this *GaugeChart) UniqueId() string {
	return this.Id
}

func (this *GaugeChart) SetUniqueId(id string) {
	this.Id = id
}
