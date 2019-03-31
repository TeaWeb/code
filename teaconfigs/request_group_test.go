package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestRequestGroup_Match(t *testing.T) {
	a := assert.NewAssertion(t)

	formatter := func(source string) string {
		if source == "${remoteAddr}" {
			return "192.168.1.100"
		}
		if source == "${arg.id}" {
			return "20"
		}
		return ""
	}

	{
		group := NewRequestGroup()
		group.Validate()
		a.IsTrue(group.Match(formatter))
	}

	{
		group := NewRequestGroup()
		{
			ipRange := NewIPRangeConfig()
			ipRange.Type = IPRangeTypeRange
			ipRange.Param = "${remoteAddr}"
			ipRange.IPFrom = "192.168.1.1"
			ipRange.IPTo = "192.168.1.200"
			group.AddIPRange(ipRange)
		}
		err := group.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(group.Match(formatter))
	}

	{
		group := NewRequestGroup()
		{
			ipRange := NewIPRangeConfig()
			ipRange.Type = IPRangeTypeRange
			ipRange.Param = "${remoteAddr}"
			ipRange.IPFrom = "192.168.1.1"
			ipRange.IPTo = "192.168.1.100"
			group.AddIPRange(ipRange)
		}
		err := group.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(group.Match(formatter))
	}

	{
		group := NewRequestGroup()
		{
			ipRange := NewIPRangeConfig()
			ipRange.Type = IPRangeTypeRange
			ipRange.Param = "${remoteAddr}"
			ipRange.IPFrom = "192.168.1.1"
			ipRange.IPTo = "192.168.1.99"
			group.AddIPRange(ipRange)
		}
		err := group.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(group.Match(formatter))
	}

	{
		group := NewRequestGroup()
		{
			ipRange := NewIPRangeConfig()
			ipRange.Type = IPRangeTypeCIDR
			ipRange.Param = "${remoteAddr}"
			ipRange.CIDR = "192.168.1.1/24"
			group.AddIPRange(ipRange)
		}
		err := group.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(group.Match(formatter))
	}

	{
		group := NewRequestGroup()
		{
			cond := NewRequestCond()
			cond.Param = "${arg.id}"
			cond.Operator = RequestCondOperatorGt
			cond.Value = "19"
			group.AddCond(cond)
		}
		{
			ipRange := NewIPRangeConfig()
			ipRange.Type = IPRangeTypeCIDR
			ipRange.Param = "${remoteAddr}"
			ipRange.CIDR = "192.168.1.1/24"
			group.AddIPRange(ipRange)
		}
		err := group.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(group.Match(formatter))
	}
}
