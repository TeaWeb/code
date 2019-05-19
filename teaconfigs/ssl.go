package teaconfigs

import (
	"crypto/tls"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/pkg/errors"
	"strings"
)

// TLS Version
type TLSVersion = string

var AllTlsVersions = []TLSVersion{"SSL 3.0", "TLS 1.0", "TLS 1.1", "TLS 1.2"}

// Cipher Suites
type TLSCipherSuite = string

var AllTLSCipherSuites = []TLSCipherSuite{
	"TLS_RSA_WITH_RC4_128_SHA",
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	"TLS_RSA_WITH_AES_128_CBC_SHA",
	"TLS_RSA_WITH_AES_256_CBC_SHA",
	"TLS_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_RSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
	"TLS_ECDHE_RSA_WITH_RC4_128_SHA",
	"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
}

// SSL配置
type SSLConfig struct {
	On bool `yaml:"on" json:"on"` // 是否开启

	Certificate    string `yaml:"certificate" json:"certificate"`       // 证书文件, deprecated in v0.1.4
	CertificateKey string `yaml:"certificateKey" json:"certificateKey"` // 密钥, deprecated in v0.1.4

	Certs []*SSLCertConfig `yaml:"certs" json:"certs"`

	Listen       []string         `yaml:"listen" json:"listen"`             // 网络地址
	MinVersion   TLSVersion       `yaml:"minVersion" json:"minVersion"`     // 支持的最小版本
	CipherSuites []TLSCipherSuite `yaml:"cipherSuites" json:"cipherSuites"` // 加密算法套件

	nameMapping map[string]*tls.Certificate // dnsName => cert

	minVersion   uint16
	cipherSuites []uint16
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

	if len(this.Certs) == 0 {
		if len(this.Certificate) > 0 && len(this.CertificateKey) > 0 {
			newCert := NewSSLCertConfig(this.Certificate, this.CertificateKey)
			newCert.Id = "old_version_cert"
			this.Certs = append(this.Certs, newCert)
		}
	}

	if len(this.Certs) == 0 {
		return errors.New("no certificates in https config")
	}

	for _, cert := range this.Certs {
		err := cert.Validate()
		if err != nil {
			return err
		}
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

	// min version
	switch this.MinVersion {
	case "SSL 3.0":
		this.minVersion = tls.VersionSSL30
	case "TLS 1.0":
		this.minVersion = tls.VersionTLS10
	case "TLS 1.1":
		this.minVersion = tls.VersionTLS11
	case "TLS 1.2":
		this.minVersion = tls.VersionTLS12
	default:
		this.minVersion = tls.VersionTLS10
	}

	// cipher suites
	suites := []uint16{}
	for _, suite := range this.CipherSuites {
		switch suite {
		case "TLS_RSA_WITH_RC4_128_SHA":
			suites = append(suites, tls.TLS_RSA_WITH_RC4_128_SHA)
		case "TLS_RSA_WITH_3DES_EDE_CBC_SHA":
			suites = append(suites, tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA)
		case "TLS_RSA_WITH_AES_128_CBC_SHA":
			suites = append(suites, tls.TLS_RSA_WITH_AES_128_CBC_SHA)
		case "TLS_RSA_WITH_AES_256_CBC_SHA":
			suites = append(suites, tls.TLS_RSA_WITH_AES_256_CBC_SHA)
		case "TLS_RSA_WITH_AES_128_CBC_SHA256":
			suites = append(suites, tls.TLS_RSA_WITH_AES_128_CBC_SHA256)
		case "TLS_RSA_WITH_AES_128_GCM_SHA256":
			suites = append(suites, tls.TLS_RSA_WITH_AES_128_GCM_SHA256)
		case "TLS_RSA_WITH_AES_256_GCM_SHA384":
			suites = append(suites, tls.TLS_RSA_WITH_AES_256_GCM_SHA384)
		case "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":
			suites = append(suites, tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA)
		case "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":
			suites = append(suites, tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA)
		case "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":
			suites = append(suites, tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA)
		case "TLS_ECDHE_RSA_WITH_RC4_128_SHA":
			suites = append(suites, tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA)
		case "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":
			suites = append(suites, tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA)
		case "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":
			suites = append(suites, tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA)
		case "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":
			suites = append(suites, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA)
		case "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256":
			suites = append(suites, tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256)
		case "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":
			suites = append(suites, tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256)
		case "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":
			suites = append(suites, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256)
		case "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":
			suites = append(suites, tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256)
		case "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":
			suites = append(suites, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384)
		case "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":
			suites = append(suites, tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384)
		case "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":
			suites = append(suites, tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305)
		case "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":
			suites = append(suites, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305)
		}
	}
	this.cipherSuites = suites

	return nil
}

// 取得最小版本
func (this *SSLConfig) TLSMinVersion() uint16 {
	return this.minVersion
}

// 套件
func (this *SSLConfig) TLSCipherSuites() []uint16 {
	return this.cipherSuites
}

// 校验是否匹配某个域名
func (this *SSLConfig) MatchDomain(domain string) (cert *tls.Certificate, ok bool) {
	for _, cert := range this.Certs {
		if cert.MatchDomain(domain) {
			return cert.CertObject(), true
		}
	}
	return nil, false
}

// 取得第一个证书
func (this *SSLConfig) FirstCert() *tls.Certificate {
	for _, cert := range this.Certs {
		return cert.CertObject()
	}
	return nil
}

// 是否包含某个证书或密钥路径
func (this *SSLConfig) ContainsFile(file string) bool {
	for _, cert := range this.Certs {
		if cert.CertFile == file || cert.KeyFile == file {
			return true
		}
	}
	return false
}

// 删除证书文件
func (this *SSLConfig) DeleteFiles() error {
	var resultErr error = nil

	if len(this.Certificate) > 0 {
		err := files.NewFile(Tea.ConfigFile(this.Certificate)).Delete()
		if err != nil {
			resultErr = err
		}
	}

	if len(this.CertificateKey) > 0 {
		err := files.NewFile(Tea.ConfigFile(this.CertificateKey)).Delete()
		if err != nil {
			resultErr = err
		}
	}

	for _, cert := range this.Certs {
		err := cert.DeleteFiles()
		if err != nil {
			resultErr = err
		}
	}

	return resultErr
}

// 查找单个证书配置
func (this *SSLConfig) FindCert(certId string) *SSLCertConfig {
	for _, cert := range this.Certs {
		if cert.Id == certId {
			return cert
		}
	}
	return nil
}
