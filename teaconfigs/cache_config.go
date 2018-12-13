package teaconfigs

import "github.com/iwind/TeaGo/utils/string"

// 缓存配置
type CacheConfig struct {
	On       bool   `yaml:"on" json:"on"`             // 是否开启
	Key      string `yaml:"key" json:"key"`           // Key
	Capacity string `yaml:"capacity" json:"capacity"` // 最大内容容量
	Life     string `yaml:"life" json:"life"`         // 时间
	Status   []int  `yaml:"status" json:"status"`     // 缓存的状态码列表
	MaxSize  string `yaml:"maxSize" json:"maxSize"`   // 能够请求的最大尺寸

	maxSize float64

	Type    string                 `yaml:"type" json:"type"`       // 类型
	Options map[string]interface{} `yaml:"options" json:"options"` // 选项
}

// 获取新对象
func NewCacheConfig() *CacheConfig {
	return &CacheConfig{}
}

// 校验
func (this *CacheConfig) Validate() error {
	this.maxSize, _ = stringutil.ParseFileSize(this.MaxSize)
	return nil
}

func (this *CacheConfig) MaxDataSize() float64 {
	return this.maxSize
}
