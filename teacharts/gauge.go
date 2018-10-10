package teacharts

type GaugeChart struct {
	Chart

	Value float64 `json:"value"`
	Label string  `json:"label"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Unit  string  `json:"unit"`
}

func NewGaugeChart() *GaugeChart {
	p := &GaugeChart{}
	p.Type = "gauge"
	return p
}

func (this *GaugeChart) UniqueId() string {
	return this.Id
}

func (this *GaugeChart) SetUniqueId(id string) {
	this.Id = id
}
