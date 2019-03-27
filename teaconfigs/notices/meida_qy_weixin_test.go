package notices

import "testing"

func TestNewNoticeQyWeixinMedia(t *testing.T) {
	m := NewNoticeQyWeixinMedia()
	m.CorporateId = "xxx"
	m.AppSecret = "xxx"
	m.AgentId = "1000003"
	resp, err := m.Send("", "标题：报警标题", "内容：报警内容/全员都有")
	if err != nil {
		t.Log(string(resp))
		t.Fatal(err)
	}
	t.Log(string(resp))
}
