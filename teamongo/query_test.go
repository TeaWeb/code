package teamongo

import (
	"github.com/TeaWeb/code/teaconfigs/audits"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestQuery_Insert(t *testing.T) {
	log := audits.NewLog("zhangsan", audits.ActionLogin, "登录", map[string]string{
		"name": "value",
	})

	query := NewAuditsQuery()
	err := query.Insert(log)
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuery_Count(t *testing.T) {
	query := NewAuditsQuery()
	query.Gte("timestamp", 1548468624)
	count, err := query.Count()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result:", count)
}

func TestQuery_FindAll(t *testing.T) {
	query := NewAuditsQuery()
	ones, err := query.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, one := range ones {
		t.Logf("%#v", one.(*audits.Log))
	}
}

func TestQuery_Find(t *testing.T) {
	query := NewAuditsQuery()
	one, err := query.Find()
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(one, t)
}
