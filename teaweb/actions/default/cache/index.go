package cache

import "github.com/iwind/TeaGo/actions"

type IndexAction actions.Action

// 缓存首页
func (this *IndexAction) Run(params struct{}) {


	this.Show()
}
