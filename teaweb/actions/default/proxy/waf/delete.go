package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/waf/wafutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
)

type DeleteAction actions.Action

// 删除
func (this *DeleteAction) RunPost(params struct {
	WafId string
}) {
	if len(params.WafId) == 0 {
		this.Fail("请输入要删除的WAF ID")
	}

	filename := "waf." + params.WafId + ".conf"
	path := Tea.ConfigFile(filename)
	file := files.NewFile(path)
	if !file.Exists() {
		this.Fail("要删除的WAF不存在")
	}

	wafList := teaconfigs.SharedWAFList()
	wafList.RemoveFile(filename)
	err := wafList.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	err = file.Delete()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	// 通知刷新
	if wafutils.IsPolicyUsed(params.WafId) {
		proxyutils.NotifyChange()
	}

	this.Success()
}
