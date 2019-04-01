package notices

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/actions"
)

type BadgeAction actions.Action

// 计算未读数量
func (this *BadgeAction) Run(params struct{}) {
	this.Data["count"] = noticeutils.CountUnreadNotices()
	this.Data["soundOn"] = notices.SharedNoticeSetting().SoundOn

	this.Success()
}
