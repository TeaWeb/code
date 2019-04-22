package groups

import "github.com/TeaWeb/code/teawaf/rules"

var emptyBodyGroup = &rules.RuleGroup{
	On:       false,
	Name:     "空请求体",
	Code:     "emptyBody",
	RuleSets: []*rules.RuleSet{},
}
