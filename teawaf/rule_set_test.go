package teawaf

import (
	"net/http"
	"testing"
)

func TestRuleSet_MatchRequest(t *testing.T) {
	set := NewRuleSet()
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

	err := set.Init()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://teaos.cn/hello?name=lu&age=20", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set.MatchRequest(req))
}

func TestRuleSet_MatchRequest2(t *testing.T) {
	set := NewRuleSet()
	set.Connector = RuleConnectorOr

	set.Rules = []*Rule{
		{
			Param:    "${arg.name}",
			Operator: RuleOperatorEqString,
			Value:    "lu",
		},
		{
			Param:    "${arg.age}",
			Operator: RuleOperatorEq,
			Value:    "21",
		},
	}

	err := set.Init()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://teaos.cn/hello?name=lu&age=20", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set.MatchRequest(req))
}
