package groups

import "github.com/TeaWeb/code/teawaf/rules"

var whiteListGroup = &rules.RuleGroup{
	On:       true,
	Name:     "白名单",
	Code:     "whiteList",
	RuleSets: []*rules.RuleSet{},
}
