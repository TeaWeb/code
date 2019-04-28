package rules

import (
	"github.com/TeaWeb/code/teawaf/checkpoints"
	"github.com/TeaWeb/code/teawaf/requests"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"net/url"
	"testing"
)

func TestRule_Init_Single(t *testing.T) {
	rule := NewRule()
	rule.Param = "${arg.name}"
	rule.Operator = RuleOperatorEqString
	rule.Value = "lu"
	err := rule.Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rule.singleParam, rule.singleCheckpoint)
	rawReq, err := http.NewRequest(http.MethodGet, "http://teaos.cn/hello?name=lu&age=20", nil)
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewRequest(rawReq)
	t.Log(rule.MatchRequest(req))
}

func TestRule_Init_Composite(t *testing.T) {
	rule := NewRule()
	rule.Param = "${arg.name} ${arg.age}"
	rule.Operator = RuleOperatorContains
	rule.Value = "lu"
	err := rule.Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rule.singleParam, rule.singleCheckpoint)

	rawReq, err := http.NewRequest(http.MethodGet, "http://teaos.cn/hello?name=lu&age=20", nil)
	if err != nil {
		t.Fatal(err)
	}
	req := requests.NewRequest(rawReq)
	t.Log(rule.MatchRequest(req))
}

func TestRule_Test(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		rule := NewRule()
		rule.Operator = RuleOperatorGt
		rule.Value = "123"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("124"))
		a.IsFalse(rule.Test("123"))
		a.IsFalse(rule.Test("122"))
		a.IsFalse(rule.Test("abcdef"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorGte
		rule.Value = "123"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("124"))
		a.IsTrue(rule.Test("123"))
		a.IsFalse(rule.Test("122"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorLt
		rule.Value = "123"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("124"))
		a.IsFalse(rule.Test("123"))
		a.IsTrue(rule.Test("122"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorLte
		rule.Value = "123"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("124"))
		a.IsTrue(rule.Test("123"))
		a.IsTrue(rule.Test("122"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorEq
		rule.Value = "123"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("124"))
		a.IsTrue(rule.Test("123"))
		a.IsFalse(rule.Test("122"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorNeq
		rule.Value = "123"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("124"))
		a.IsFalse(rule.Test("123"))
		a.IsTrue(rule.Test("122"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorEqString
		rule.Value = "123"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("124"))
		a.IsTrue(rule.Test("123"))
		a.IsFalse(rule.Test("122"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorEqString
		rule.Value = "abc"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("ABC"))
		a.IsTrue(rule.Test("abc"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorEqString
		rule.IsCaseInsensitive = true
		rule.Value = "abc"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("ABC"))
		a.IsTrue(rule.Test("abc"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorNeqString
		rule.Value = "abc"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("124"))
		a.IsFalse(rule.Test("abc"))
		a.IsTrue(rule.Test("122"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorNeqString
		rule.IsCaseInsensitive = true
		rule.Value = "abc"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("ABC"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorMatch
		rule.Value = "^\\d+"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("123"))
		a.IsFalse(rule.Test("abc123"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorMatch
		rule.Value = "abc"
		rule.IsCaseInsensitive = true
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("ABC"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorNotMatch
		rule.Value = "\\d+"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("123"))
		a.IsTrue(rule.Test("abc"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorNotMatch
		rule.Value = "abc"
		rule.IsCaseInsensitive = true
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("ABC"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorMatch
		rule.Value = "^(?i)[a-z]+$"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("ABC"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorContains
		rule.Value = "Hello"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("Hello, World"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorContains
		rule.Value = "hello"
		rule.IsCaseInsensitive = true
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("Hello, World"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorContains
		rule.Value = "Hello"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test([]string{"Hello", "World"}))
		a.IsTrue(rule.Test(maps.Map{
			"a": "World", "b": "Hello",
		}))
		a.IsFalse(rule.Test(maps.Map{
			"a": "World", "b": "Hello2",
		}))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorNotContains
		rule.Value = "Hello"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("Hello, World"))
		a.IsTrue(rule.Test("World"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorNotContains
		rule.Value = "hello"
		rule.IsCaseInsensitive = true
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("Hello, World"))
		a.IsTrue(rule.Test("World"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorPrefix
		rule.Value = "Hello"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("Hello, World"))
		a.IsFalse(rule.Test("World, Hello"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorPrefix
		rule.Value = "hello"
		rule.IsCaseInsensitive = true
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(rule.Test("Hello, World"))
		a.IsFalse(rule.Test("World, Hello"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorSuffix
		rule.Value = "Hello"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("Hello, World"))
		a.IsTrue(rule.Test("World, Hello"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorSuffix
		rule.Value = "hello"
		rule.IsCaseInsensitive = true
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("Hello, World"))
		a.IsTrue(rule.Test("World, Hello"))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorHasKey
		rule.Value = "Hello"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("Hello, World"))
		a.IsTrue(rule.Test(maps.Map{
			"Hello": "World",
		}))
		a.IsFalse(rule.Test(maps.Map{
			"Hello1": "World",
		}))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorHasKey
		rule.Value = "hello"
		rule.IsCaseInsensitive = true
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("Hello, World"))
		a.IsTrue(rule.Test(maps.Map{
			"Hello": "World",
		}))
		a.IsFalse(rule.Test(maps.Map{
			"Hello1": "World",
		}))
	}

	{
		rule := NewRule()
		rule.Operator = RuleOperatorHasKey
		rule.Value = "3"
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(rule.Test("Hello, World"))
		a.IsFalse(rule.Test(maps.Map{
			"Hello": "World",
		}))
		a.IsTrue(rule.Test([]int{1, 2, 3, 4}))
	}
}

func TestRule_MatchStar(t *testing.T) {
	{
		rule := NewRule()
		rule.Operator = RuleOperatorMatch
		rule.Value = `/\*(!|\x00)`
		err := rule.Init()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(rule.Test("/*!"))
		t.Log(rule.Test(url.QueryEscape("/*!")))
		t.Log(url.QueryEscape("/*!"))
	}
}

func TestRule_SetCheckpointFinder(t *testing.T) {
	{
		rule := NewRule()
		rule.Param = "${arg.abc}"
		rule.Operator = RuleOperatorMatch
		rule.Init()
		t.Logf("%#v", rule.singleCheckpoint)
	}

	{
		rule := NewRule()
		rule.Param = "${arg.abc}"
		rule.Operator = RuleOperatorMatch
		rule.checkpointFinder = func(prefix string) checkpoints.CheckpointInterface {
			return new(checkpoints.SampleRequestCheckpoint)
		}
		rule.Init()
		t.Logf("%#v", rule.singleCheckpoint)
	}
}
