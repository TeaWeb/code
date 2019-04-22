package groups

import "github.com/TeaWeb/code/teawaf/rules"

var blackListGroup = &rules.RuleGroup{
	On:       true,
	Name:     "黑名单",
	Code:     "blackList",
	RuleSets: []*rules.RuleSet{},
}
