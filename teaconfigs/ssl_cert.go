package teaconfigs

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/utils/string"
	"strings"
)

// SSL证书
type SSLCertConfig struct {
	Id          string `yaml:"id" json:"id"`
	On          bool   `yaml:"on" json:"on"`
	Description string `yaml:"description" json:"description"`
	CertFile    string `yaml:"certFile" json:"certFile"`
	KeyFile     string `yaml:"keyFile" json:"keyFile"`
	IsLocal     bool   `yaml:"isLocal" json:"isLocal"` // if is local file
	TaskId      string `yaml:"taskId" json:"taskId"`   // 生成证书任务ID

	dnsNames []string
	cert     *tls.Certificate
}

// 获取新的SSL证书
func NewSSLCertConfig(certFile string, keyFile string) *SSLCertConfig {
	return &SSLCertConfig{
		On:       true,
		Id:       stringutil.Rand(16),
		CertFile: certFile,
		KeyFile:  keyFile,
	}
}

// 校验
func (this *SSLCertConfig) Validate() error {
	if len(this.CertFile) == 0 {
		return errors.New("cert file should not be empty")
	}
	if len(this.KeyFile) == 0 {
		return errors.New("key file should not be empty")
	}
	cert, err := tls.LoadX509KeyPair(this.FullCertPath(), this.FullKeyPath())
	if err != nil {
		return errors.New("load certificate '" + this.CertFile + "', '" + this.KeyFile + "' failed:" + err.Error())
	}

	this.dnsNames = []string{}
	for _, data := range cert.Certificate {
		c, err := x509.ParseCertificate(data)
		if err != nil {
			continue
		}
		dnsNames := c.DNSNames
		if len(dnsNames) > 0 {
			for _, dnsName := range dnsNames {
				if !lists.ContainsString(this.dnsNames, dnsName) {
					this.dnsNames = append(this.dnsNames, dnsName)
				}
			}
		}
	}

	this.cert = &cert
	return nil
}

// 证书文件路径
func (this *SSLCertConfig) FullCertPath() string {
	if !strings.ContainsAny(this.CertFile, "/\\") {
		return Tea.ConfigFile(this.CertFile)
	}
	return this.CertFile
}

// 密钥文件路径
func (this *SSLCertConfig) FullKeyPath() string {
	if !strings.ContainsAny(this.KeyFile, "/\\") {
		return Tea.ConfigFile(this.KeyFile)
	}
	return this.KeyFile
}

// 校验是否匹配某个域名
func (this *SSLCertConfig) MatchDomain(domain string) bool {
	if len(this.dnsNames) == 0 {
		return false
	}
	return teautils.MatchDomains(this.dnsNames, domain)
}

// 获取证书对象
func (this *SSLCertConfig) CertObject() *tls.Certificate {
	return this.cert
}

// 删除文件
func (this *SSLCertConfig) DeleteFiles() error {
	if this.IsLocal {
		return nil
	}

	var resultErr error = nil
	if len(this.CertFile) > 0 && !strings.ContainsAny(this.CertFile, "/\\") {
		err := files.NewFile(this.FullCertPath()).Delete()
		if err != nil {
			resultErr = err
		}
	}

	if len(this.KeyFile) > 0 && !strings.ContainsAny(this.KeyFile, "/\\") {
		err := files.NewFile(this.FullKeyPath()).Delete()
		if err != nil {
			resultErr = err
		}
	}
	return resultErr
}
