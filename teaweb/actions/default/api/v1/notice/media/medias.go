package media

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
)

type MediasAction actions.Action

// 媒介列表
func (this *MediasAction) RunGet(params struct{}) {
	setting := notices.SharedNoticeSetting()

	apiutils.Success(this, setting.Medias)
}
