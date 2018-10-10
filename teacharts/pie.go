package teacharts

type PieChart struct {
	Chart
	Values []interface{} `json:"values"`
	Labels []string      `json:"labels"`
}

func NewPieChart() *PieChart {
	p := &PieChart{
		Values: []interface{}{},
		Labels: []string{},
	}
	p.Type = "pie"
	return p
}

func (this *PieChart) UniqueId() string {
	return this.Id
}

func (this *PieChart) SetUniqueId(id string) {
	this.Id = id
}
