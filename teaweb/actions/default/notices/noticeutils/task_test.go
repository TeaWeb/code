package noticeutils

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"testing"
)

func TestAddTask(t *testing.T) {
	AddTask(notices.NoticeLevelWarning, []*notices.NoticeReceiver{}, "subject", "message")
}
