package teaconfigs

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/Tea"
	"strings"
)

// SSL配置
type SSLConfig struct {
	On             bool     `yaml:"on" json:"on"`                         // 是否开启
	Certificate    string   `yaml:"certificate" json:"certificate"`       // 证书文件
	CertificateKey string   `yaml:"certificateKey" json:"certificateKey"` // 密钥
	Listen         []string `yaml:"listen" json:"listen"`                 // 网络地址

	cert     *tls.Certificate
	dnsNames []string
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

	cert, err := tls.LoadX509KeyPair(Tea.ConfigFile(this.Certificate), Tea.ConfigFile(this.CertificateKey))
	if err != nil {
		return errors.New("load certificate '" + this.Certificate + "', '" + this.CertificateKey + "' failed:" + err.Error())
	}

	for _, data := range cert.Certificate {
		c, err := x509.ParseCertificate(data)
		if err != nil {
			continue
		}
		dnsNames := c.DNSNames
		if len(dnsNames) > 0 {
			this.dnsNames = append(this.dnsNames, dnsNames...)
		}
	}

	this.cert = &cert

	return nil
}

// 取得Certificate对象
func (this *SSLConfig) CertificateObject() *tls.Certificate {
	return this.cert
}

// 校验是否匹配某个域名
func (this *SSLConfig) MatchDomain(domain string) bool {
	if len(this.dnsNames) == 0 {
		return false
	}
	return teautils.MatchDomains(this.dnsNames, domain)
}
