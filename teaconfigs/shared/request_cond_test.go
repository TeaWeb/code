package shared

import (
	"bytes"
	"github.com/iwind/TeaGo/assert"
	"net"
	"testing"
)

func TestRequestCond_Compare1(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		cond := RequestCond{
			Param:    "/hello",
			Operator: RequestCondOperatorRegexp,
			Value:    "abc",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello",
			Operator: RequestCondOperatorRegexp,
			Value:    "/\\w+",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RequestCond{
			Param:    "/article/123.html",
			Operator: RequestCondOperatorRegexp,
			Value:    `^/article/\d+\.html$`,
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello",
			Operator: RequestCondOperatorRegexp,
			Value:    "[",
		}
		a.IsNotNil(cond.Validate())
		a.IsFalse(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello",
			Operator: RequestCondOperatorNotRegexp,
			Value:    "abc",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello",
			Operator: RequestCondOperatorNotRegexp,
			Value:    "/\\w+",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RequestCond{
			Param:    "123.123",
			Operator: RequestCondOperatorEqInt,
			Value:    "123",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "123",
			Operator: RequestCondOperatorEqInt,
			Value:    "123",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "abc",
			Operator: RequestCondOperatorEqInt,
			Value:    "abc",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "123",
			Operator: RequestCondOperatorEqFloat,
			Value:    "123",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "123.0",
			Operator: RequestCondOperatorEqFloat,
			Value:    "123",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "123.123",
			Operator: RequestCondOperatorEqFloat,
			Value:    "123.12",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "123",
			Operator: RequestCondOperatorGtFloat,
			Value:    "1",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "123",
			Operator: RequestCondOperatorGtFloat,
			Value:    "125",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "125",
			Operator: RequestCondOperatorGteFloat,
			Value:    "125",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "125",
			Operator: RequestCondOperatorLtFloat,
			Value:    "127",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "125",
			Operator: RequestCondOperatorLteFloat,
			Value:    "127",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "125",
			Operator: RequestCondOperatorEqString,
			Value:    "125",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "125",
			Operator: RequestCondOperatorNeqString,
			Value:    "125",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "125",
			Operator: RequestCondOperatorNeqString,
			Value:    "127",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello/world",
			Operator: RequestCondOperatorHasPrefix,
			Value:    "/hello",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello/world",
			Operator: RequestCondOperatorHasPrefix,
			Value:    "/hello2",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello/world",
			Operator: RequestCondOperatorHasSuffix,
			Value:    "world",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello/world",
			Operator: RequestCondOperatorHasSuffix,
			Value:    "world/",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello/world",
			Operator: RequestCondOperatorContainsString,
			Value:    "wo",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello/world",
			Operator: RequestCondOperatorContainsString,
			Value:    "wr",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello/world",
			Operator: RequestCondOperatorNotContainsString,
			Value:    "HELLO",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "/hello/world",
			Operator: RequestCondOperatorNotContainsString,
			Value:    "hello",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}
}

func TestRequestCond_IP(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		cond := RequestCond{
			Param:    "hello",
			Operator: RequestCondOperatorEqIP,
			Value:    "hello",
		}
		a.IsNotNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.1.100",
			Operator: RequestCondOperatorEqIP,
			Value:    "hello",
		}
		a.IsNotNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.1.100",
			Operator: RequestCondOperatorEqIP,
			Value:    "192.168.1.100",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.1.100",
			Operator: RequestCondOperatorGtIP,
			Value:    "192.168.1.90",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.1.100",
			Operator: RequestCondOperatorGteIP,
			Value:    "192.168.1.90",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.1.80",
			Operator: RequestCondOperatorLtIP,
			Value:    "192.168.1.90",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.0.100",
			Operator: RequestCondOperatorLteIP,
			Value:    "192.168.1.90",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.0.100",
			Operator: RequestCondOperatorIPInRange,
			Value:    "192.168.0.90,",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.0.100",
			Operator: RequestCondOperatorIPInRange,
			Value:    "192.168.0.90,192.168.1.100",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.0.100",
			Operator: RequestCondOperatorIPInRange,
			Value:    ",192.168.1.100",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.1.100",
			Operator: RequestCondOperatorIPInRange,
			Value:    "192.168.0.90,192.168.1.99",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.1.100",
			Operator: RequestCondOperatorIPInRange,
			Value:    "192.168.0.90/24",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.1.100",
			Operator: RequestCondOperatorIPInRange,
			Value:    "192.168.0.90/18",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "192.168.1.100",
			Operator: RequestCondOperatorIPInRange,
			Value:    "a/18",
		}
		a.IsNotNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}
}

func TestRequestCondIPCompare(t *testing.T) {
	{
		ip1 := net.ParseIP("192.168.3.100")
		ip2 := net.ParseIP("192.168.2.100")
		t.Log(bytes.Compare(ip1, ip2))
	}

	{
		ip1 := net.ParseIP("192.168.3.100")
		ip2 := net.ParseIP("a")
		t.Log(bytes.Compare(ip1, ip2))
	}

	{
		ip1 := net.ParseIP("b")
		ip2 := net.ParseIP("192.168.2.100")
		t.Log(bytes.Compare(ip1, ip2))
	}

	{
		ip1 := net.ParseIP("b")
		ip2 := net.ParseIP("a")
		t.Log(ip1 == nil)
		t.Log(bytes.Compare(ip1, ip2))
	}
}

func TestRequestCond_In(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		cond := RequestCond{
			Param:    "a",
			Operator: RequestCondOperatorIn,
			Value:    `a`,
		}
		a.IsNotNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "a",
			Operator: RequestCondOperatorIn,
			Value:    `["a", "b"]`,
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "c",
			Operator: RequestCondOperatorNotIn,
			Value:    `["a", "b"]`,
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RequestCond{
			Param:    "a",
			Operator: RequestCondOperatorNotIn,
			Value:    `["a", "b"]`,
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}
}
