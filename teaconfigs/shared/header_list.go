package shared

import (
	"github.com/iwind/TeaGo/lists"
)

// HeaderList相关操作接口
type HeaderListInterface interface {
	// 校验
	ValidateHeaders() error

	// 取得所有的IgnoreHeader
	AllIgnoreHeaders() []string

	// 添加IgnoreHeader
	AddIgnoreHeader(name string)

	// 判断是否包含IgnoreHeader
	ContainsIgnoreHeader(name string) bool

	// 移除IgnoreHeader
	RemoveIgnoreHeader(name string)

	// 修改IgnoreHeader
	UpdateIgnoreHeader(oldName string, newName string)

	// 取得所有的Header
	AllHeaders() []*HeaderConfig

	// 添加Header
	AddHeader(header *HeaderConfig)

	// 判断是否包含Header
	ContainsHeader(name string) bool

	// 查找Header
	FindHeader(headerId string) *HeaderConfig

	// 移除Header
	RemoveHeader(headerId string)

	// 格式化Headers
	FormatHeaders(formatter func(source string) string) []*HeaderConfig
}

// HeaderList定义
type HeaderList struct {
	hasHeaders bool

	// 添加的Headers
	Headers []*HeaderConfig `yaml:"headers" json:"headers"`

	// 忽略的Headers
	IgnoreHeaders []string `yaml:"ignoreHeaders" json:"ignoreHeaders"`
}

// 校验
func (this *HeaderList) ValidateHeaders() error {
	this.hasHeaders = len(this.Headers) > 0

	for _, h := range this.Headers {
		err := h.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// 是否有Headers
func (this *HeaderList) HasHeaders() bool {
	return this.hasHeaders
}

// 取得所有的IgnoreHeader
func (this *HeaderList) AllIgnoreHeaders() []string {
	if this.IgnoreHeaders == nil {
		return []string{}
	}
	return this.IgnoreHeaders
}

// 添加IgnoreHeader
func (this *HeaderList) AddIgnoreHeader(name string) {
	if !lists.ContainsString(this.IgnoreHeaders, name) {
		this.IgnoreHeaders = append(this.IgnoreHeaders, name)
	}
}

// 判断是否包含IgnoreHeader
func (this *HeaderList) ContainsIgnoreHeader(name string) bool {
	if len(this.IgnoreHeaders) == 0 {
		return false
	}
	return lists.ContainsString(this.IgnoreHeaders, name)
}

// 修改IgnoreHeader
func (this *HeaderList) UpdateIgnoreHeader(oldName string, newName string) {
	result := []string{}
	for _, h := range this.IgnoreHeaders {
		if h == oldName {
			result = append(result, newName)
		} else {
			result = append(result, h)
		}
	}
	this.IgnoreHeaders = result
}

// 移除IgnoreHeader
func (this *HeaderList) RemoveIgnoreHeader(name string) {
	result := []string{}
	for _, n := range this.IgnoreHeaders {
		if n == name {
			continue
		}
		result = append(result, n)
	}
	this.IgnoreHeaders = result
}

// 取得所有的Header
func (this *HeaderList) AllHeaders() []*HeaderConfig {
	if this.Headers == nil {
		return []*HeaderConfig{}
	}
	return this.Headers
}

// 添加Header
func (this *HeaderList) AddHeader(header *HeaderConfig) {
	this.Headers = append(this.Headers, header)
}

// 判断是否包含Header
func (this *HeaderList) ContainsHeader(name string) bool {
	for _, h := range this.Headers {
		if h.Name == name {
			return true
		}
	}
	return false
}

// 查找Header
func (this *HeaderList) FindHeader(headerId string) *HeaderConfig {
	for _, h := range this.Headers {
		if h.Id == headerId {
			return h
		}
	}
	return nil
}

// 移除Header
func (this *HeaderList) RemoveHeader(headerId string) {
	result := []*HeaderConfig{}
	for _, h := range this.Headers {
		if h.Id == headerId {
			continue
		}
		result = append(result, h)
	}
	this.Headers = result
}

// 格式化Header
func (this *HeaderList) FormatHeaders(formatter func(source string) string) []*HeaderConfig {
	result := []*HeaderConfig{}
	for _, h := range this.Headers {
		if !h.On {
			continue
		}
		newHeader := h.Copy()
		if h.hasVariables {
			newHeader.Value = formatter(h.Value)
		}
		result = append(result, newHeader)
	}
	return result
}
