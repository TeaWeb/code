package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestRewriteRule(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rule := RewriteRule{
		Pattern: "/(hello)/(world)",
		Replace: "/${1}/${2}",
	}
	a.IsNil(rule.Validate())

	{
		_, b := rule.Match("/hello/worl", func(source string) string {
			return source
		})
		a.IsFalse(b)
		a.Log("proxy:", rule.TargetProxy())
		a.Log("url:", rule.TargetURL())
	}

	{
		_, b := rule.Match("/hello/world", func(source string) string {
			return source
		})
		a.IsTrue(b)
		a.Log("proxy:", rule.TargetProxy())
		a.Log("url:", rule.TargetURL())
	}
}

func TestRewriteRuleProxy(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rule := RewriteRule{
		Pattern: "/(hello)/(world)",
		Replace: "proxy://lb001/${1}/${2}",
	}
	a.IsNil(rule.Validate())

	_, b := rule.Match("/hello/world", func(source string) string {
		return source
	})
	a.IsTrue(b)
	a.IsTrue(rule.TargetProxy() == "lb001")
	a.IsTrue(rule.TargetURL() == "/hello/world")
}
