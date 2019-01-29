package notices

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AddReceiverAction actions.Action

// 添加接收者
func (this *AddReceiverAction) Run(params struct {
	Level notices.NoticeLevel
}) {
	level := notices.FindNoticeLevel(params.Level)
	if level == nil {
		this.Fail("找不到Level信息")
	}
	this.Data["level"] = level

	setting := notices.SharedNoticeSetting()
	mediaMaps := []maps.Map{}
	for _, media := range setting.Medias {
		if !media.On {
			continue
		}
		mediaType := notices.FindNoticeMediaType(media.Type)
		if mediaType == nil {
			continue
		}
		mediaMaps = append(mediaMaps, maps.Map{
			"id":               media.Id,
			"name":             media.Name,
			"typeName":         mediaType["name"],
			"type":             media.Type,
			"mediaDescription": mediaType["description"],
			"userDescription":  mediaType["user"],
		})
	}
	this.Data["medias"] = mediaMaps

	this.Show()
}

// 提交保存
func (this *AddReceiverAction) RunPost(params struct {
	Level   notices.NoticeLevel
	On      bool
	Name    string
	MediaId string
	User    string
	Must    actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入接收人名称").
		Field("mediaId", params.MediaId).
		Require("请选择使用的媒介").
		Field("user", params.User).
		Require("请输入接收人标识")

	receiver := notices.NewNoticeReceiver()
	receiver.On = params.On
	receiver.Name = params.Name
	receiver.MediaId = params.MediaId
	receiver.User = params.User

	setting := notices.SharedNoticeSetting()
	levelConfig := setting.LevelConfig(params.Level)
	levelConfig.AddReceiver(receiver)

	err := setting.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
