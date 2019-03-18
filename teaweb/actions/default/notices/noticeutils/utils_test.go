package noticeutils

import (
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
