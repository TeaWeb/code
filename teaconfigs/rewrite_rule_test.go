package teaconfigs

import (
	"testing"
	"github.com/iwind/TeaGo/assert"
)

func TestRewriteRule(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rule := RewriteRule{
		Pattern: "/(hello)/(world)",
		Replace: "/${1}/${2}",
	}
	a.IsNil(rule.Validate())

	{
		a.IsFalse(rule.Apply("/hello/worl", func(source string) string {
			return source
		}))
		a.Log("proxy:", rule.TargetProxy())
		a.Log("url:", rule.TargetURL())
	}

	{
		a.IsTrue(rule.Apply("/hello/world", func(source string) string {
			return source
		}))
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
	a.IsTrue(rule.Apply("/hello/world", func(source string) string {
		return source
	}))
	a.IsTrue(rule.TargetProxy() == "lb001")
	a.IsTrue(rule.TargetURL() == "/hello/world")
}
