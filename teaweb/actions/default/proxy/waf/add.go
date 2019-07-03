package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teawaf"
	actions2 "github.com/TeaWeb/code/teawaf/actions"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
)

type AddAction actions.Action

// 添加策略
func (this *AddAction) RunGet(params struct{}) {
	this.Data["groups"] = lists.Map(teawaf.Template().Inbound, func(k int, v interface{}) interface{} {
		g := v.(*rules.RuleGroup)
		return maps.Map{
			"name":      g.Name,
			"code":      g.Code,
			"isChecked": g.On,
		}
	})

	this.Show()
}

// 保存提交
func (this *AddAction) RunPost(params struct {
	Name       string
	GroupCodes []string

	On bool

	BlockStatusCode string
	BlockBody       string
	BlockURL        string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	if len(params.BlockStatusCode) > 0 && !regexp.MustCompile(`^\d{3}$`).MatchString(params.BlockStatusCode) {
		this.FailField("blockStatusCode", "请输入正确的HTTP状态码")
	}
	statusCode := types.Int(params.BlockStatusCode)

	waf := teawaf.NewWAF()
	waf.Name = params.Name
	waf.On = params.On
	waf.ActionBlock = &actions2.BlockAction{
		StatusCode: statusCode,
		Body:       params.BlockBody,
		URL:        params.BlockURL,
	}

	template := teawaf.Template()

	for _, g := range template.Inbound {
		newGroup := rules.NewRuleGroup()
		newGroup.Id = stringutil.Rand(16)
		newGroup.On = lists.ContainsString(params.GroupCodes, g.Code)
		newGroup.Code = g.Code
		newGroup.Name = g.Name
		newGroup.RuleSets = g.RuleSets
		newGroup.IsInbound = g.IsInbound
		newGroup.Description = g.Description
		waf.AddRuleGroup(newGroup)
	}

	filename := "waf." + waf.Id + ".conf"
	err := waf.Save(Tea.ConfigFile(filename))
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	wafList := teaconfigs.SharedWAFList()
	wafList.AddFile(filename)
	err = wafList.Save()
	if err != nil {
		err1 := files.NewFile(Tea.ConfigFile(filename)).Delete()
		if err1 != nil {
			logs.Error(err1)
		}

		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
