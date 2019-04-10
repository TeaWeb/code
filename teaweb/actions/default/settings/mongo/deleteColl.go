package mongo

import (
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
	"golang.org/x/net/context"
	"time"
)

type DeleteCollAction actions.Action

// 删除集合
func (this *DeleteCollAction) Run(params struct {
	CollName string
}) {
	if len(params.CollName) > 0 {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err := teamongo.SharedClient().Database(teamongo.DatabaseName).Collection(params.CollName).Drop(ctx)
		if err != nil {
			this.Fail("删除失败：" + err.Error())
		}
	}
	this.Success()
}
