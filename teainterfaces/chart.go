package teainterfaces

// 图表接口
type ChartInterface interface {
	Id() string
	SetId(id string)
	Type() string
	Name() string
	Detail() string
}

// 进度条相关
type ProgressBarInterface interface {
	Value() float64
	Color() Color
}

// 仪表盘相关
type GaugeChartInterface interface {
	Value() float64
	Label() string
	Min() float64
	Max() float64
	Unit() string
}

// 线图相关
type LineInterface interface {
	Name() string
	Values() []interface{}
	Color() Color
	Filled() bool
	ShowItems() bool
}

type LineChartInterface interface {
	Lines() []interface{}
	Labels() []string

	Max() float64
	XShowTick() bool // X轴是否显示刻度

	YTickCount() uint // Y轴刻度分隔数量
	YShowTick() bool  // Y轴是否显示刻度
}

// 饼图相关
type PieChartInterface interface {
	Values() []interface{}
	Labels() []string
}

// 表格相关
type RowInterface interface {
	Columns() []interface{}
}

type ColumnInterface interface {
	Text() string
	Width() float64 // 百分比，比如 30 表示 30%
}

type TableInterface interface {
	Rows() []interface{}
}
