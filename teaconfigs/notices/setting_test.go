package notices

import (
	"testing"
	"time"
)

func TestNoticeSetting_Notify(t *testing.T) {
	setting := SharedNoticeSetting()
	receiverIds := setting.Notify(NoticeLevelInfo, "消息内容第1行\n消息内容第2行", func(receiverId string, minutes int) int {
		return 0
	})
	t.Log(receiverIds)
	time.Sleep(3 * time.Second)
}
