package groups

import "github.com/TeaWeb/code/teawaf/rules"

var sqlGroup = &rules.RuleGroup{
	On:       true,
	Name:     "SQL注入",
	Code:     "sql",
	RuleSets: []*rules.RuleSet{},
}
