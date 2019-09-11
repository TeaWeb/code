package mongo

import (
	"github.com/TeaWeb/code/teadb"
	"github.com/iwind/TeaGo/actions"
)

type DeleteCollAction actions.Action

// 删除集合
func (this *DeleteCollAction) Run(params struct {
	CollName string
}) {
	if len(params.CollName) > 0 {
		err := teadb.SharedDB().DropTable(params.CollName)
		if err != nil {
			this.Fail("删除失败：" + err.Error())
		}
	}
	this.Success()
}
