package teaproxy

import (
	"github.com/TeaWeb/code/teawaf/actions"
	"github.com/iwind/TeaGo/logs"
	"net/http"
)

// call request waf
func (this *Request) callWAFRequest(writer *ResponseWriter) (blocked bool) {
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
			this.SetAttr("waf.ruleset", ruleSet.Name)
			this.SetAttr("waf.id", this.waf.Id)
		}
	}

	return !goNext
}

// call response waf
func (this *Request) callWAFResponse(resp *http.Response, writer *ResponseWriter) (blocked bool) {
	if this.waf == nil {
		return
	}

	goNext, ruleSet, err := this.waf.MatchResponse(this.raw, resp, writer)
	if err != nil {
		logs.Error(err)
		return
	}

	if ruleSet != nil {
		if ruleSet.Action != actions.ActionAllow {
			this.SetAttr("waf.action", ruleSet.Action)
			this.SetAttr("waf.ruleset", ruleSet.Name)
			this.SetAttr("waf.id", this.waf.Id)
		}
	}

	return !goNext
}
