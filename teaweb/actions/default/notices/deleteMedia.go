package notices

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
)

type DeleteMediaAction actions.Action

// 删除媒介
func (this *DeleteMediaAction) Run(params struct {
	MediaId string
}) {
	setting := notices.SharedNoticeSetting()
	setting.RemoveMedia(params.MediaId)
	err := setting.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
