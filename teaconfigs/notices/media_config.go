package notices

import (
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// 媒介配置定义
type NoticeMediaConfig struct {
	Id      string                 `yaml:"id" json:"id"`
	On      bool                   `yaml:"on" json:"on"`
	Name    string                 `yaml:"name" json:"name"`
	Type    NoticeMediaType        `yaml:"type" json:"type"`
	Options map[string]interface{} `yaml:"options" json:"options"`
}

// 获取新对象
func NewNoticeMediaConfig() *NoticeMediaConfig {
	return &NoticeMediaConfig{
		On: true,
		Id: stringutil.Rand(16),
	}
}

// 校验
func (this *NoticeMediaConfig) Validate() error {
	teautils.JSONMap(this.Options)
	return nil
}

// 取得原始的媒介
func (this *NoticeMediaConfig) Raw() (NoticeMediaInterface, error) {
	m := FindNoticeMediaType(this.Type)
	if m == nil {
		return nil, errors.New("media type '" + this.Type + "' not found")
	}
	instance := m["instance"]
	err := teautils.MapToObjectJSON(this.Options, instance)
	if err != nil {
		return nil, err
	}
	return instance.(NoticeMediaInterface), nil
}
