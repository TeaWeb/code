package teaproxy

import (
	"github.com/TeaWeb/code/teawaf/actions"
	"github.com/iwind/TeaGo/logs"
)

// call waf
func (this *Request) callWAF(writer *ResponseWriter) (blocked bool) {
	if this.waf == nil {
		return
	}
	goNext, ruleSet, err := this.waf.MatchRequest(this.raw, writer)
	if err != nil {
		logs.Error(err)
		return
	}

	if ruleSet != nil {
		if ruleSet.Action != actions.ActionAllow {
			this.SetAttr("waf.action", ruleSet.Action)
		}
	}

	return !goNext
}
