package agents

import (
	"github.com/TeaWeb/code/teaconfigs/forms"
	"github.com/TeaWeb/code/teaconfigs/widgets"
)

// 插件
type PluginSource struct {
	Source `yaml:",inline"`

	form         *forms.Form
	name         string
	code         string
	description  string
	variables    []*SourceVariable
	thresholds   []*Threshold
	charts       []*widgets.Chart
	presentation *forms.Presentation
	platforms    []string
}

// 获取新对象
func NewPluginSource() *PluginSource {
	return &PluginSource{}
}

// 名称
func (this *PluginSource) Name() string {
	return this.name
}

// 设置名称
func (this *PluginSource) SetName(name string) {
	this.name = name
}

// 代号
func (this *PluginSource) Code() string {
	return this.code
}

// 设置代号
func (this *PluginSource) SetCode(code string) {
	this.code = code
}

// 描述
func (this *PluginSource) Description() string {
	return this.description
}

// 设置描述
func (this *PluginSource) SetDescription(description string) {
	this.description = description
}

// 执行
func (this *PluginSource) Execute(params map[string]string) (value interface{}, err error) {
	return
}

// 表单信息
func (this *PluginSource) Form() *forms.Form {
	return this.form
}

// 设置表单信息
func (this *PluginSource) SetForm(form *forms.Form) {
	this.form = form
}

// 变量
func (this *PluginSource) Variables() []*SourceVariable {
	return this.variables
}

// 添加变量
func (this *PluginSource) AddVariable(variable *SourceVariable) {
	this.variables = append(this.variables, variable)
}

// 阈值
func (this *PluginSource) Thresholds() []*Threshold {
	return this.thresholds
}

// 添加阈值
func (this *PluginSource) AddThreshold(threshold *Threshold) {
	this.thresholds = append(this.thresholds, threshold)
}

// 图表
func (this *PluginSource) Charts() []*widgets.Chart {
	return this.charts
}

// 添加图表
func (this *PluginSource) AddChart(chart *widgets.Chart) {
	this.charts = append(this.charts, chart)
}

// 显示信息
func (this *PluginSource) Presentation() *forms.Presentation {
	return this.presentation
}

// 设置显示信息
func (this *PluginSource) SetPresentation(presentation *forms.Presentation) {
	this.presentation = presentation
}

// 平台限制
func (this *PluginSource) Platforms() []string {
	return this.platforms
}

// 设置平台限制
func (this *PluginSource) SetPlatforms(platforms []string) {
	this.platforms = platforms
}
