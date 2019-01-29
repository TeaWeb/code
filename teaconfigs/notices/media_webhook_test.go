package notices

import (
	"testing"
)

func TestNoticeMediaWebhook_Send(t *testing.T) {
	media := NewNoticeWebhookMedia()
	media.URL = "http://baidu.com/s?subject=${NoticeSubject}&body=${NoticeBody}"
	t.Log(media.Send("zhangsan", "this is subject", "this is body"))
}
