package teaconfigs

import (
	"bytes"
	"errors"
	"github.com/iwind/TeaGo/utils/string"
	"net"
)

// IP Range类型
type IPRangeType = int

const (
	IPRangeTypeRange = IPRangeType(1)
	IPRangeTypeCIDR  = IPRangeType(2)
)

// IP Range
type IPRangeConfig struct {
	Id string `yaml:"id" json:"id"`

	Type IPRangeType `yaml:"type" json:"type"`

	Param  string `yaml:"param" json:"param"`
	CIDR   string `yaml:"cidr" json:"cidr"`
	IPFrom string `yaml:"ipFrom" json:"ipFrom"`
	IPTo   string `yaml:"ipTo" json:"ipTo"`

	cidr   *net.IPNet
	ipFrom net.IP
	ipTo   net.IP
}

// 获取新对象
func NewIPRangeConfig() *IPRangeConfig {
	return &IPRangeConfig{
		Id: stringutil.Rand(16),
	}
}

// 校验
func (this *IPRangeConfig) Validate() error {
	if this.Type == IPRangeTypeCIDR {
		if len(this.CIDR) == 0 {
			return errors.New("cidr should not be empty")
		}

		_, cidr, err := net.ParseCIDR(this.CIDR)
		if err != nil {
			return err
		}
		this.cidr = cidr
	}

	if this.Type == IPRangeTypeRange {
		this.ipFrom = net.ParseIP(this.IPFrom)
		this.ipTo = net.ParseIP(this.IPTo)

		if this.ipFrom.To4() == nil && this.ipFrom.To16() == nil {
			return errors.New("from ip should in IPv4 or IPV6 format")
		}

		if this.ipTo.To4() == nil && this.ipTo.To16() == nil {
			return errors.New("to ip should in IPv4 or IPV6 format")
		}
	}

	return nil
}

// 是否包含某个IP
func (this *IPRangeConfig) Contains(ipString string) bool {
	ip := net.ParseIP(ipString)
	if ip.To4() == nil {
		return false
	}
	if this.Type == IPRangeTypeCIDR {
		if this.cidr == nil {
			return false
		}
		return this.cidr.Contains(ip)
	}
	if this.Type == IPRangeTypeRange {
		if this.ipFrom == nil || this.ipTo == nil {
			return false
		}
		return bytes.Compare(ip, this.ipFrom) >= 0 && bytes.Compare(ip, this.ipTo) <= 0
	}
	return false
}
