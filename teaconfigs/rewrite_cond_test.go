package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestRewriteCond(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		cond := RewriteCond{
			Param:    "/hello",
			Operator: RewriteOperatorRegexp,
			Value:    "abc",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RewriteCond{
			Param:    "/hello",
			Operator: RewriteOperatorRegexp,
			Value:    "/\\w+",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RewriteCond{
			Param:    "/article/123.html",
			Operator: RewriteOperatorRegexp,
			Value:    `^/article/\d+\.html$`,
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RewriteCond{
			Param:    "/hello",
			Operator: RewriteOperatorRegexp,
			Value:    "[",
		}
		a.IsNotNil(cond.Validate())
		a.IsFalse(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RewriteCond{
			Param:    "123",
			Operator: RewriteOperatorGt,
			Value:    "1",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RewriteCond{
			Param:    "123",
			Operator: RewriteOperatorGt,
			Value:    "125",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RewriteCond{
			Param:    "125",
			Operator: RewriteOperatorGte,
			Value:    "125",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RewriteCond{
			Param:    "125",
			Operator: RewriteOperatorLt,
			Value:    "127",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RewriteCond{
			Param:    "125",
			Operator: RewriteOperatorLte,
			Value:    "127",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RewriteCond{
			Param:    "125",
			Operator: RewriteOperatorEq,
			Value:    "125",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RewriteCond{
			Param:    "125",
			Operator: RewriteOperatorNot,
			Value:    "125",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(source string) string {
			return source
		}))
	}

	{
		cond := RewriteCond{
			Param:    "125",
			Operator: RewriteOperatorNot,
			Value:    "127",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(source string) string {
			return source
		}))
	}
}
