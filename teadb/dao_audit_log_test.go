package teadb

import (
	"github.com/TeaWeb/code/teaconfigs/audits"
	"github.com/TeaWeb/code/teadb/shared"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"testing"
	"time"
)

func TestAuditLogDAO_CountAllAuditLogs(t *testing.T) {
	dao := AuditLogDAO()
	count, err := dao.CountAllAuditLogs()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestAuditLogDAO_ListAuditLogs(t *testing.T) {
	dao := AuditLogDAO()
	result, err := dao.ListAuditLogs(0, 5)
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range result {
		t.Log(timeutil.Format("Y-m-d H:i:s", time.Unix(r.Timestamp, 0)), r)
	}
}

func TestAuditLogDAO_InsertOne(t *testing.T) {
	dao := AuditLogDAO()
	err := dao.InsertOne(&audits.Log{
		Id:          shared.NewObjectId(),
		Username:    "test",
		Action:      "test",
		Description: "test",
		Options:     nil,
		Timestamp:   time.Now().Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("ok")
}

func TestAuditLogDAO_InsertOne2(t *testing.T) {
	dao := AuditLogDAO()
	err := dao.InsertOne(&audits.Log{
		Id:          shared.NewObjectId(),
		Username:    "test",
		Action:      "test",
		Description: "test",
		Options: map[string]string{
			"name": "lu",
			"age":  "100",
		},
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("ok")
}
