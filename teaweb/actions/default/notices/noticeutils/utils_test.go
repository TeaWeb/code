package noticeutils

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"testing"
	"time"
)

func TestCountReceivedNotices(t *testing.T) {
	t.Log(CountReceivedNotices("gijBqXS0OxNYmRq0", map[string]interface{}{}, 1))
	t.Log(CountReceivedNotices("gijBqXS0OxNYmRq0", map[string]interface{}{
		"agent.agentId": "local",
		"agent.appId":   "system",
		"agent.itemId":  "cpu.usage",
	}, 1024000))
	t.Log(CountReceivedNotices("gijBqXS0OxNYmRq0", map[string]interface{}{
		"agent.agentId": "ZRc78EC4GJr2AO8B",
	}, 1024000))
}

func TestExistNoticesWithHash(t *testing.T) {
	t.Log(ExistNoticesWithHash("3494563121", map[string]interface{}{
		"agent.agentId": "zlokAzjGVN7ENbC6",
	}, 1*time.Hour))
}

func TestNotifyProxyMessage(t *testing.T) {
	{
		err := NotifyProxyMessage(notices.ProxyCond{
			ServerId: "JnZ03pbDebcOQ7h9",
			Level:    notices.NoticeLevelSuccess,
		}, "HelloWorld")
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err := NotifyProxyMessage(notices.ProxyCond{
			ServerId: "JnZ03pbDebcOQ7h9",
			Level:    notices.NoticeLevelError,
		}, "HelloWorld")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestNoticeServerProxyMessage(t *testing.T) {
	err := NotifyProxyServerMessage("JnZ03pbDebcOQ7h9", notices.NoticeLevelError, "This is a test")
	if err != nil {
		t.Fatal(err)
	}
}
