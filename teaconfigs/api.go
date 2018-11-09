package teaconfigs

import (
	"regexp"
)

//  API定义
type API struct {
	Path           string          `yaml:"path" json:"path"`                     // 访问路径
	Address        string          `yaml:"address" json:"address"`               // 实际地址
	Methods        []string        `yaml:"methods" json:"methods"`               // 方法
	Params         []*APIParam     `yaml:"params" json:"params"`                 // 参数
	Name           string          `yaml:"name" json:"name"`                     // 名称
	Description    string          `yaml:"description" json:"description"`       // 描述 TODO 需要支持markdown
	Mock           []string        `yaml:"mock" json:"mock"`                     // TODO
	Author         string          `yaml:"author" json:"author"`                 // 作者
	Company        string          `yaml:"company" json:"company"`               // 公司或团队
	IsAsynchronous bool            `yaml:"isAsynchronous" json:"isAsynchronous"` // TODO
	Timeout        float64         `yaml:"timeout" json:"timeout"`               // TODO
	MaxSize        uint            `yaml:"maxSize" json:"maxSize"`               // TODO
	Headers        []*HeaderConfig `yaml:"headers" json:"headers"`               // TODO
	TodoThings     []string        `yaml:"todo" json:"todo"`                     // 待做事宜
	DoneThings     []string        `yaml:"done" json:"done"`                     // 已完成事宜
	Response       []byte          `yaml:"response" json:"response"`             // TODO
	IsDeprecated   bool            `yaml:"isDeprecated" json:"isDeprecated"`     // 是否过期
	On             bool            `yaml:"on" json:"on"`                         // 是否开启
	Versions       []string        `yaml:"versions" json:"versions"`             // 版本信息
	ModifiedAt     int64           `yaml:"modifiedAt" json:"modifiedAt"`         // 最后修改时间
	Username       string          `yaml:"username" json:"username"`             // 最后修改用户名
	Groups         []string        `yaml:"groups" json:"groups"`                 // 分组

	pathReg    *regexp.Regexp // 匹配模式
	pathParams []string
}

// 获取新API对象
func NewAPI() *API {
	return &API{
		On: true,
	}
}

// 执行校验
func (this *API) Validate() error {
	this.pathParams = []string{}
	reg := regexp.MustCompile(`:\w+`)
	if reg.MatchString(this.Path) {
		newPath := reg.ReplaceAllStringFunc(this.Path, func(s string) string {
			param := s[1:]
			this.pathParams = append(this.pathParams, param)
			return "(.+)"
		})

		pathReg, err := regexp.Compile(newPath)
		if err != nil {
			return err
		}
		this.pathReg = pathReg
	}
	return nil
}

// 添加参数
func (this *API) AddParam(param *APIParam) {
	this.Params = append(this.Params, param)
}

// 格式化Header
func (this *API) FormatHeaders(formatter func(source string) string) []*HeaderConfig {
	result := []*HeaderConfig{}
	for _, header := range this.Headers {
		result = append(result, &HeaderConfig{
			Name:   header.Name,
			Value:  formatter(header.Value),
			Always: header.Always,
			Status: header.Status,
		})
	}
	return result
}

// 使用正则匹配路径
func (this *API) Match(path string) (params map[string]string, matched bool) {
	if this.pathReg == nil {
		return nil, false
	}
	if !this.pathReg.MatchString(path) {
		return nil, false
	}

	params = map[string]string{}
	matched = true
	matches := this.pathReg.FindStringSubmatch(path)
	for index, match := range matches {
		if index == 0 {
			continue
		}
		params[this.pathParams[index-1]] = match
	}
	return
}

// 是否允许某个请求方法
func (this *API) AllowMethod(method string) bool {
	for _, m := range this.Methods {
		if m == method {
			return true
		}
	}
	return false
}

// 删除某个分组
func (this *API) RemoveGroup(name string) {
	result := []string{}
	for _, g := range this.Groups {
		if g != name {
			result = append(result, g)
		}
	}
	this.Groups = result
}

// 修改API某个分组名
func (this *API) ChangeGroup(oldName string, newName string) {
	result := []string{}
	for _, g := range this.Groups {
		if g == oldName {
			result = append(result, newName)
		} else {
			result = append(result, g)
		}
	}
	this.Groups = result
}

// 删除某个版本
func (this *API) RemoveVersion(name string) {
	result := []string{}
	for _, g := range this.Versions {
		if g != name {
			result = append(result, g)
		}
	}
	this.Versions = result
}

// 修改API某个版本号
func (this *API) ChangeVersion(oldName string, newName string) {
	result := []string{}
	for _, g := range this.Versions {
		if g == oldName {
			result = append(result, newName)
		} else {
			result = append(result, g)
		}
	}
	this.Versions = result
}

// 开始监控
func (this *API) StartWatching() {
	SharedApiWatching.Add(this.Path)
}

// 结束监控
func (this *API) StopWatching() {
	SharedApiWatching.Remove(this.Path)
}

// 是否在监控
func (this *API) IsWatching() bool {
	return SharedApiWatching.Contains(this.Path)
}
