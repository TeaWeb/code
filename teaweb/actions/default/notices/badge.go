package notices

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type BadgeAction actions.Action

// 计算未读数量
func (this *BadgeAction) Run(params struct{}) {
	count, err := teadb.NoticeDAO().CountAllUnreadNotices()
	if err != nil {
		logs.Error(err)
	}
	this.Data["count"] = count
	this.Data["soundOn"] = notices.SharedNoticeSetting().SoundOn

	this.Success()
}
