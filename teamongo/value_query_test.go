package teamongo

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestInsertValues(t *testing.T) {
	query := NewValueQuery()

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

func TestNewQuery(t *testing.T) {
	query := NewValueQuery()
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
