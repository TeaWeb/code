package teadb

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/assert"
	"testing"
	"time"
)

func TestNoticeDAO_InsertOne(t *testing.T) {
	notice := notices.NewNotice()
	notice.Message = "this is test"
	notice.Hash()

	dao := NoticeDAO()
	err := dao.InsertOne(notice)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_NotifyProxyMessage(t *testing.T) {
	dao := NoticeDAO()
	err := dao.NotifyProxyMessage(notices.ProxyCond{
		ServerId: "test",
	}, "this is test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_NotifyProxyServerMessage(t *testing.T) {
	dao := NoticeDAO()
	err := dao.NotifyProxyServerMessage("test2", notices.NoticeLevelWarning, "Hello")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_CountAllCountUnreadNotices(t *testing.T) {
	dao := NoticeDAO()
	count, err := dao.CountAllUnreadNotices()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestNoticeDAO_CountAllCountReadNotices(t *testing.T) {
	dao := NoticeDAO()
	count, err := dao.CountAllReadNotices()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestNoticeDAO_CountUnreadNoticesForAgent(t *testing.T) {
	dao := NoticeDAO()
	count, err := dao.CountUnreadNoticesForAgent("local")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestNoticeDAO_CountReadNoticesForAgent(t *testing.T) {
	dao := NoticeDAO()
	count, err := dao.CountReadNoticesForAgent("local")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestNoticeDAO_CountReceivedNotices(t *testing.T) {
	dao := NoticeDAO()
	count, err := dao.CountReceivedNotices("EpBqPQMqpRlvFh9Q", map[string]interface{}{}, 86400)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestNoticeDAO_ExistNoticesWithHash(t *testing.T) {
	dao := NoticeDAO()
	b, err := dao.ExistNoticesWithHash("4157704578", map[string]interface{}{}, 86400*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Log("exists")
	} else {
		t.Log("not exists")
	}
}

func TestNoticeDAO_ListNotices(t *testing.T) {
	result, err := NoticeDAO().ListNotices(true, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(result), "notices")
	for _, n := range result {
		t.Log(n.Id, n.Message)
	}
}

func TestNoticeDAO_ListNotices_Unread(t *testing.T) {
	result, err := NoticeDAO().ListNotices(false, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(result), "notices")
	for _, n := range result {
		t.Log(n.Id, n.Message)
	}
}

func TestNoticeDAO_ListAgentNotices(t *testing.T) {
	result, err := NoticeDAO().ListAgentNotices("local", true, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(result), "notices")
	for _, n := range result {
		t.Log(n.Id, n.Message)
	}
}

func TestNoticeDAO_ListAgentNotices_Notfound(t *testing.T) {
	a := assert.NewAssertion(t)
	result, err := NoticeDAO().ListAgentNotices("local123", true, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	a.IsTrue(len(result) == 0)
}

func TestNoticeDAO_DeleteNoticesForAgent(t *testing.T) {
	dao := NoticeDAO()
	err := dao.DeleteNoticesForAgent("1TABzdF0uAIFPGkr")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_UpdateNoticeReceivers(t *testing.T) {
	dao := NoticeDAO()
	err := dao.UpdateNoticeReceivers("5d6284e5e31cd55d472761c4", []string{"a", "b", "c"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_UpdateAllNoticesRead(t *testing.T) {
	dao := NoticeDAO()
	err := dao.UpdateAllNoticesRead()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_UpdateNoticesRead(t *testing.T) {
	dao := NoticeDAO()
	err := dao.UpdateNoticesRead([]string{"5c92676949d81a09b2f1a760", "5c925b0086e004b594e9a90c"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_UpdateAgentNoticesRead(t *testing.T) {
	err := NoticeDAO().UpdateAgentNoticesRead("local", []string{"5c925b0086e004b594e9a90c"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestNoticeDAO_UpdateAllAgentNoticesRead(t *testing.T) {
	err := NoticeDAO().UpdateAllAgentNoticesRead("local")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
