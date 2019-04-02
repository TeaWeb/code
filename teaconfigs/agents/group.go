package agents

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/utils/string"
)

// Agent分组
type Group struct {
	Id            string                                            `yaml:"id" json:"id"`
	On            bool                                              `yaml:"on" json:"on"`
	Name          string                                            `yaml:"name" json:"name"`
	Index         int                                               `yaml:"index" json:"index"`
	NoticeSetting map[notices.NoticeLevel][]*notices.NoticeReceiver `yaml:"noticeSetting" json:"noticeSetting"`
}

// 获取新分组
func NewGroup(name string) *Group {
	return &Group{
		Id:   stringutil.Rand(16),
		On:   true,
		Name: name,
	}
}

// 添加通知接收者
func (this *Group) AddNoticeReceiver(level notices.NoticeLevel, receiver *notices.NoticeReceiver) {
	if this.NoticeSetting == nil {
		this.NoticeSetting = map[notices.NoticeLevel][]*notices.NoticeReceiver{}
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		receivers = []*notices.NoticeReceiver{}
	}
	receivers = append(receivers, receiver)
	this.NoticeSetting[level] = receivers
}

// 删除通知接收者
func (this *Group) RemoveNoticeReceiver(level notices.NoticeLevel, receiverId string) {
	if this.NoticeSetting == nil {
		return
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		return
	}

	result := []*notices.NoticeReceiver{}
	for _, r := range receivers {
		if r.Id == receiverId {
			continue
		}
		result = append(result, r)
	}
	this.NoticeSetting[level] = result
}

// 删除媒介
func (this *Group) RemoveMedia(mediaId string) (found bool) {
	for level, receivers := range this.NoticeSetting {
		result := []*notices.NoticeReceiver{}
		for _, receiver := range receivers {
			if receiver.MediaId == mediaId {
				found = true
				continue
			}
			result = append(result, receiver)
		}
		this.NoticeSetting[level] = result
	}
	return
}
