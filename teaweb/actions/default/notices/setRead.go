package notices

import (
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SetReadAction actions.Action

// 设置已读
func (this *SetReadAction) Run(params struct {
	Scope     string
	NoticeIds []string
}) {
	if params.Scope == "page" {
		if len(params.NoticeIds) == 0 {
			this.Success()
		}

		err := noticeutils.NewNoticeQuery().
			Attr("_id", lists.Map(params.NoticeIds, func(k int, v interface{}) interface{} {
				noticeId := v.(string)

				objectId, err := primitive.ObjectIDFromHex(noticeId)
				if err != nil {
					return noticeId
				} else {
					return objectId
				}
			})).
			Update(maps.Map{
				"$set": maps.Map{
					"isRead": true,
				},
			})
		if err != nil {
			this.Fail("操作失败：" + err.Error())
		}
	} else {
		noticeutils.NewNoticeQuery().
			Update(maps.Map{
				"$set": maps.Map{
					"isRead": true,
				},
			})
	}

	this.Success()
}
