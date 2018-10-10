package teaconfigs

import (
	"testing"
	"github.com/iwind/TeaGo/assert"
)

func TestRewriteCond(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	{
		cond := RewriteCond{
			Test:    "/hello",
			Pattern: "abc",
		}
		a.IsNil(cond.Validate())
		a.IsFalse(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RewriteCond{
			Test:    "/hello",
			Pattern: "/\\w+",
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RewriteCond{
			Test:    "/article/123.html",
			Pattern: `^/article/\d+\.html$`,
		}
		a.IsNil(cond.Validate())
		a.IsTrue(cond.Match(func(format string) string {
			return format
		}))
	}

	{
		cond := RewriteCond{
			Test:    "/hello",
			Pattern: "[",
		}
		a.IsNotNil(cond.Validate())
		a.IsFalse(cond.Match(func(format string) string {
			return format
		}))
	}
}
