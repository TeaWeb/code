package teawaf

import (
	"github.com/TeaWeb/code/teawaf/actions"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/iwind/TeaGo/assert"
	"net/http"
	"testing"
)

func TestWAF_MatchRequest(t *testing.T) {
	a := assert.NewAssertion(t)

	set := rules.NewRuleSet()
	set.Name = "Name_Age"
	set.Connector = rules.RuleConnectorAnd
	set.Rules = []*rules.Rule{
		{
			Param:    "${arg.name}",
			Operator: rules.RuleOperatorEqString,
			Value:    "lu",
		},
		{
			Param:    "${arg.age}",
			Operator: rules.RuleOperatorEq,
			Value:    "20",
		},
	}
	set.Action = actions.ActionBlock

	group := rules.NewRuleGroup()
	group.AddRuleSet(set)
	group.IsInbound = true

	waf := NewWAF()
	waf.AddRuleGroup(group)
	err := waf.Init()
	if err != nil {
		t.Fatal(err)
	}

	waf.OnAction(func(action actions.ActionString) (goNext bool) {
		return action != actions.ActionBlock
	})

	req, err := http.NewRequest(http.MethodGet, "http://teaos.cn/hello?name=lu&age=20", nil)
	if err != nil {
		t.Fatal(err)
	}
	goNext, set, err := waf.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	if set == nil {
		t.Log("not match")
		return
	}
	t.Log("goNext:", goNext, "set:", set.Name)
	a.IsFalse(goNext)
}
