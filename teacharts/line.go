package teacharts

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

func NewLineChart() *LineChart {
	p := &LineChart{}
	p.Type = "line"
	p.Lines = []*Line{}
	return p
}

func (this *LineChart) UniqueId() string {
	return this.Id
}

func (this *LineChart) SetUniqueId(id string) {
	this.Id = id
}

func (this *LineChart) AddLine(line *Line) {
	this.Lines = append(this.Lines, line)
}

func (this *LineChart) ResetLines() {
	this.Lines = []*Line{}
}
