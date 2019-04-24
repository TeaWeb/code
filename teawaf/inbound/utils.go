package inbound

import "github.com/TeaWeb/code/teawaf/rules"

var InternalGroups = []*rules.RuleGroup{
	sqlGroup,
	xssGroup,
	emptyBodyGroup,
	whiteListGroup,
	blackListGroup,
}
