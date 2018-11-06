package teaconfigs

import (
	"regexp"
)

//  API定义
type API struct {
	Path           string          `yaml:"path" json:"path"`                     // 访问路径
	Address        string          `yaml:"address" json:"address"`               // 实际地址
	Methods        []string        `yaml:"methods" json:"methods"`               // TODO
	Params         []*APIParam     `yaml:"params" json:"params"`                 // TODO
	Name           string          `yaml:"name" json:"name"`                     // TODO
	Description    string          `yaml:"description" json:"description"`       // TODO
	Mock           []string        `yaml:"mock" json:"mock"`                     // TODO
	Author         string          `yaml:"author" json:"author"`                 // TODO
	Company        string          `yaml:"company" json:"company"`               // TODO
	IsAsynchronous bool            `yaml:"isAsynchronous" json:"isAsynchronous"` // TODO
	Timeout        float64         `yaml:"timeout" json:"timeout"`               // TODO
	MaxSize        uint            `yaml:"maxSize" json:"maxSize"`               // TODO
	Headers        []*HeaderConfig `yaml:"headers" json:"headers"`               // TODO
	TodoThings     []string        `yaml:"todo" json:"todo"`                     // TODO
	DoneThings     []string        `yaml:"done" json:"done"`                     // TODO
	Response       []byte          `yaml:"response" json:"response"`             // TODO
	Roles          []string        `yaml:"roles" json:"roles"`                   // TODO
	IsDeprecated   bool            `yaml:"isDeprecated" json:"isDeprecated"`     // TODO
	On             bool            `yaml:"on" json:"on"`                         // TODO
	Versions       []string        `yaml:"versions" json:"versions"`             // TODO
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
