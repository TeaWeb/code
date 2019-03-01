package teamongo

import (
	"encoding/json"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestInsertValues(t *testing.T) {
	query := NewAgentValueQuery()

	{
		value := agents.NewValue()
		value.AgentId = "local"
		value.AppId = "mysql"
		value.SetTime(time.Now())
		value.Value = 3
		err := query.Insert(value)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAgentValueQuery_Insert(t *testing.T) {
	jsonString := `
{
    "code": 500,
    "message": "\u8bf7\u8f93\u5165\u6b63\u786e\u7684\u4ee4\u724c\uff08001\uff09",
    "data": {},
    "next": null,
    "errors": []
}`
	v := map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonString), &v)
	if err != nil {
		t.Fatal(err)
	}
	value := &agents.Value{
		AppId:       "1",
		AgentId:     "1",
		ItemId:      "1",
		Value:       v,
		Error:       "",
		NoticeLevel: notices.NoticeLevelWarning,
	}
	value.SetTime(time.Now())

	err = NewAgentValueQuery().Insert(value)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("success")
	}
}

func TestNewQuery(t *testing.T) {
	query := NewAgentValueQuery()
	query.Action(ValueQueryActionFindAll)
	query.Agent("local")
	result, err := query.
		Desc("value").
		Limit(1).
		Execute()
	if err != nil {
		t.Fatal(err)
	} else {
		logs.PrintAsJSON(result)
	}
}
