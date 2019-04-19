package teawaf

import (
	"github.com/TeaWeb/code/teawaf/actions"
	"github.com/iwind/TeaGo/assert"
	"net/http"
	"testing"
)

func TestWAF_MatchRequest(t *testing.T) {
	a := assert.NewAssertion(t)

	set := NewRuleSet()
	set.Name = "Name_Age"
	set.Connector = RuleConnectorAnd
	set.Rules = []*Rule{
		{
			Param:    "${arg.name}",
			Operator: RuleOperatorEqString,
			Value:    "lu",
		},
		{
			Param:    "${arg.age}",
			Operator: RuleOperatorEq,
			Value:    "20",
		},
	}
	set.Action = actions.ActionBlock

	group := NewRuleGroup()
	group.AddRuleSet(set)

	waf := NewWAF()
	waf.AddRuleGroup(group)
	waf.Init()

	req, err := http.NewRequest(http.MethodGet, "http://teaos.cn/hello?name=lu&age=20", nil)
	if err != nil {
		t.Fatal(err)
	}
	goNext, set, err := waf.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("goNext:", goNext, "set:", set.Name)
	a.IsFalse(goNext)
}
