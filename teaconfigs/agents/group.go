package agents

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"io/ioutil"
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

// 默认的分组
// 一定会有返回值
func LoadDefaultGroup() *Group {
	data, err := ioutil.ReadFile(DefaultGroupFile())
	if err != nil {
		return &Group{
			Id:   "default",
			On:   true,
			Name: "默认分组",
		}
	}
	group := new(Group)
	err = yaml.Unmarshal(data, group)
	if err != nil {
		logs.Error(err)
		return &Group{
			Id:   "default",
			On:   true,
			Name: "默认分组",
		}
	}

	return group
}

// 默认的分组文件名
func DefaultGroupFile() string {
	return Tea.ConfigFile("agents/group.default.conf")
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

// 查找一个或多个级别对应的接收者，并合并相同的接收者
func (this *Group) FindAllNoticeReceivers(level ...notices.NoticeLevel) []*notices.NoticeReceiver {
	if len(level) == 0 {
		return []*notices.NoticeReceiver{}
	}

	m := maps.Map{} // mediaId_user => bool
	result := []*notices.NoticeReceiver{}
	for _, l := range level {
		receivers, ok := this.NoticeSetting[l]
		if !ok {
			continue
		}
		for _, receiver := range receivers {
			if !receiver.On {
				continue
			}
			key := receiver.Key()
			if m.Has(key) {
				continue
			}
			m[key] = true
			result = append(result, receiver)
		}
	}
	return result
}

// 写入文件
func (this *Group) WriteToFile(path string) error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0666)
}
