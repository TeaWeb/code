package mongo

import (
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
)

type TestAction actions.Action

// 测试Mongo连接
func (this *TestAction) Run(params struct{}) {
	err := teamongo.Test()
	if err != nil {
		this.Fail()
	} else {
		this.Success()
	}
}
