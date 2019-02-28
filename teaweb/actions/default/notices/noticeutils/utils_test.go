package noticeutils

import "testing"

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
