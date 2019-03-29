package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestRequestCond(t *testing.T) {
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
			Param:    "123",
			Operator: RequestCondOperatorGt,
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
			Operator: RequestCondOperatorGt,
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
			Operator: RequestCondOperatorGte,
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
			Operator: RequestCondOperatorLt,
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
			Operator: RequestCondOperatorLte,
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
			Operator: RequestCondOperatorEq,
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
			Operator: RequestCondOperatorNot,
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
			Operator: RequestCondOperatorNot,
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
			Operator: RequestCondOperatorPrefix,
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
			Operator: RequestCondOperatorPrefix,
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
			Operator: RequestCondOperatorSuffix,
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
			Operator: RequestCondOperatorSuffix,
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
			Operator: RequestCondOperatorContains,
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
			Operator: RequestCondOperatorContains,
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
			Operator: RequestCondOperatorNotContains,
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
			Operator: RequestCondOperatorNotContains,
			Value:    "hello",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}
}
