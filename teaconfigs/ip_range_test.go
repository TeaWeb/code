package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestGeoConfig_Contains(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		geo := NewIPRangeConfig()
		geo.Type = IPRangeTypeRange
		geo.IPFrom = "192.168.1.100"
		geo.IPTo = "192.168.1.110"
		geo.Validate()
		a.IsTrue(geo.Contains("192.168.1.100"))
		a.IsTrue(geo.Contains("192.168.1.101"))
		a.IsTrue(geo.Contains("192.168.1.110"))
		a.IsFalse(geo.Contains("192.168.1.111"))
	}

	{
		geo := NewIPRangeConfig()
		geo.Type = IPRangeTypeCIDR
		geo.CIDR = "192.168.1.1/24"
		geo.Validate()
		a.IsTrue(geo.Contains("192.168.1.100"))
		a.IsFalse(geo.Contains("192.168.2.100"))
	}

	{
		geo := NewIPRangeConfig()
		geo.Type = IPRangeTypeCIDR
		geo.CIDR = "192.168.1.1/16"
		geo.Validate()
		a.IsTrue(geo.Contains("192.168.2.100"))
	}
}
