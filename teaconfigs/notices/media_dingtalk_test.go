package notices

import "testing"

func TestNoticeDingTalkMedia_Send(t *testing.T) {
	media := NewNoticeDingTalkMedia()
	media.WebhookURL = "https://oapi.dingtalk.com/robot/send?access_token=xxx"
	resp, err := media.Send("186xxx", "服务器好像出了点小问题", "IP:192.168.1.xxx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(resp))
}
