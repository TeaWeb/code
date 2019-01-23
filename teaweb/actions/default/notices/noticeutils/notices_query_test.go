package noticeutils

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/logs"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"testing"
	"time"
)

func TestNoticeQuery_Insert(t *testing.T) {
	notice := notices.NewNotice()
	notice.Id = objectid.New()
	notice.SetTime(time.Now())
	notice.Message = "Hello"
	notice.Agent = notices.AgentCond{
		AgentId: "a",
		AppId:   "b",
		ItemId:  "c",
		Level:   2,
	}
	logs.PrintAsJSON(notice)
	err := NewNoticeQuery().Insert(notice)
	if err != nil {
		t.Fatal(err)
	}
}
