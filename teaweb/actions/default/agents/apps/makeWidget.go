package apps

import (
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/iwind/TeaGo/utils/time"
	"regexp"
	"strings"
)

type MakeWidgetAction actions.Action

// 制作Widget
func (this *MakeWidgetAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	agentutils.InitAppData(this, params.AgentId, params.AppId, "widget")

	this.Data["testingYear"] = timeutil.Format("Y")
	this.Data["testingMonth"] = timeutil.Format("m")
	this.Data["testingDay"] = timeutil.Format("d")

	this.Show()
}

// 保存Widget
func (this *MakeWidgetAction) RunPost(params struct {
	AgentId           string
	AppId             string
	Name              string
	Code              string
	Author            string
	Version           string
	ParamNames        []string
	ParamCodes        []string
	ParamDefaults     []string
	ParamDescriptions []string
	ChartCode         string
	Must              *actions.Must
}) {
	params.Must.
		Field("agentId", params.AgentId).
		Require("请选择Agent").
		Field("appId", params.AppId).
		Require("请选择App").
		Field("name", params.Name).
		Require("请输入Widget名称").
		Field("code", params.Code).
		Require("请输入Widget代号").
		Expect(func() (message string, success bool) {
			pieces := strings.Split(params.Code, "@")
			if len(pieces) != 2 {
				return "请输入正确的代号", false
			}
			for _, piece := range pieces {
				if !regexp.MustCompile("^\\w+$").MatchString(piece) {
					return "代号中只能包含数字和英文字母", false;
				}
			}
			return "", true
		}).
		Field("author", params.Author).
		Require("请输入作者").
		Field("version", params.Version).
		Expect(func() (message string, success bool) {
			pieces := strings.Split(params.Version, ".")
			if len(pieces) == 0 {
				return "请输入正确的版本号", false
			}
			for _, piece := range pieces {
				if !regexp.MustCompile("^\\d+$").MatchString(piece) {
					return "请输入正确的版本号", false
				}
			}
			return "", true
		}).
		Require("请输入版本")

	widgetParams := []maps.Map{}
	for k, paramName := range params.ParamNames {
		if k < len(params.ParamCodes) && k < len(params.ParamDefaults) && k < len(params.ParamDescriptions) {
			paramCode := params.ParamCodes[k]
			paramDefault := params.ParamDefaults[k]
			paramDescription := params.ParamDescriptions[k]
			widgetParams = append(widgetParams, maps.Map{
				"name":        paramName,
				"code":        paramCode,
				"default":     paramDefault,
				"description": paramDescription,
			})
		}
	}

	widgetCode := `var widget = new widgets.Widget({
	"name": ` + stringutil.JSONEncode(params.Name) + `,
	"code": ` + stringutil.JSONEncode(params.Code) + `,
	"author": ` + stringutil.JSONEncode(params.Author) + `,
	"version": ` + stringutil.JSONEncode(params.Version) + `,
	"params": ` + stringutil.JSONEncode(widgetParams) + `,
	"options": {}
});

widget.run = function () {
` + params.ChartCode + `	
};`

	file := files.NewFile(Tea.Root + "/libs/agent/widget-" + params.Code + ".js")
	err := file.WriteString(widgetCode)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
