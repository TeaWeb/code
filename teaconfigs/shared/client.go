package shared

import (
	"github.com/iwind/TeaGo/utils/string"
	"strings"
)

// 客户端配置
type ClientConfig struct {
	Id          string `yaml:"id" json:"id"`                   // ID
	On          bool   `yaml:"on" json:"on"`                   // 是否开启
	IP          string `yaml:"ip" json:"ip"`                   // IP
	Description string `yaml:"description" json:"description"` // 描述

	hasWildcard bool
	pieces      []string
}

// 取得新配置对象
func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		Id: stringutil.Rand(16),
		On: true,
	}
}

// 校验
func (this *ClientConfig) Validate() error {
	this.hasWildcard = strings.Contains(this.IP, "*")
	if this.hasWildcard && len(this.IP) > 0 {
		this.pieces = strings.Split(this.IP, ".")
	}
	return nil
}

// 判断是否匹配某个IP
func (this *ClientConfig) Match(ip string) bool {
	if len(ip) == 0 {
		return false
	}
	if this.hasWildcard {
		pieces2 := strings.Split(ip, ".")
		if len(pieces2) != len(this.pieces) {
			return false
		}
		for index, piece2 := range pieces2 {
			if this.pieces[index] != "*" && this.pieces[index] != piece2 {
				return false
			}
		}
		return true
	}
	return this.IP == ip
}
