package groups

import "github.com/TeaWeb/code/teawaf/rules"

var xssGroup = &rules.RuleGroup{
	On:       true,
	Name:     "XSS跨站脚本",
	Code:     "xss",
	RuleSets: []*rules.RuleSet{},
}
