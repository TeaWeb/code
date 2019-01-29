package notices

import "github.com/iwind/TeaGo/utils/string"

// 接收者
type NoticeReceiver struct {
	Id      string `yaml:"id" json:"id"`
	On      bool   `yaml:"on" json:"on"`
	Name    string `yaml:"name" json:"name"`
	MediaId string `yaml:"mediaId" json:"mediaId"`
	User    string `yaml:"user" json:"user"` // 用户标识
}

// 获取新对象
func NewNoticeReceiver() *NoticeReceiver {
	return &NoticeReceiver{
		On: true,
		Id: stringutil.Rand(16),
	}
}
