package teaconfigs

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/Tea"
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
	On             bool             `yaml:"on" json:"on"`                         // 是否开启
	Certificate    string           `yaml:"certificate" json:"certificate"`       // 证书文件
	CertificateKey string           `yaml:"certificateKey" json:"certificateKey"` // 密钥
	Listen         []string         `yaml:"listen" json:"listen"`                 // 网络地址
	MinVersion     TLSVersion       `yaml:"minVersion" json:"minVersion"`         // 支持的最小版本
	CipherSuites   []TLSCipherSuite `yaml:"cipherSuites" json:"cipherSuites"`     // 加密算法套件

	cert     *tls.Certificate
	dnsNames []string

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

// 取得最小版本
func (this *SSLConfig) TLSMinVersion() uint16 {
	return this.minVersion
}

// 套件
func (this *SSLConfig) TLSCipherSuites() []uint16 {
	return this.cipherSuites
}
