package teaconfigs

import (
	"errors"
	"strings"
)

// SSL配置
type SSLConfig struct {
	On             bool     `yaml:"on" json:"on"`                         // 是否开启
	Certificate    string   `yaml:"certificate" json:"certificate"`       // 证书文件
	CertificateKey string   `yaml:"certificateKey" json:"certificateKey"` // 密钥
	Listen         []string `yaml:"listen" json:"listen"`                 // 网络地址
}

// 获取新对象
func NewSSLConfig() *SSLConfig {
	return &SSLConfig{}
}

// 校验配置
func (this *SSLConfig) Validate() error {
	if !this.On {
		return nil
	}
	if len(this.Certificate) == 0 {
		return errors.New("'certificate' should not be empty")
	}
	if len(this.CertificateKey) == 0 {
		return errors.New("'certificateKey' should not be empty")
	}
	if this.Listen == nil {
		this.Listen = []string{}
	} else {
		for index, addr := range this.Listen {
			portIndex := strings.Index(addr, ":")
			if portIndex < 0 {
				this.Listen[index] = addr + ":443"
			}
		}
	}
	return nil
}
