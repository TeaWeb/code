package teaconfigs

import "errors"

// SSL配置
type SSLConfig struct {
	On             bool   `yaml:"on" json:"on"`
	Certificate    string `yaml:"certificate" json:"certificate"`
	CertificateKey string `yaml:"certificateKey" json:"certificateKey"`
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
	return nil
}
