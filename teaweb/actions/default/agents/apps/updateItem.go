package apps

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net/http"
	"regexp"
)

type UpdateItemAction actions.Action

// 添加监控项
func (this *UpdateItemAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
	From    string
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "monitor")

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到要修改的监控项")
	}
	this.Data["item"] = item

	this.Data["from"] = params.From
	this.Data["sources"] = agents.AllDataSources()
	this.Data["methods"] = []string{http.MethodGet, http.MethodPost, http.MethodPut}
	this.Data["dataFormats"] = agents.AllSourceDataFormats()
	this.Data["operators"] = agents.AllThresholdOperators()
	this.Data["noticeLevels"] = notices.AllNoticeLevels()

	this.Show()
}

// 提交保存
func (this *UpdateItemAction) RunPost(params struct {
	AgentId string
	AppId   string
	ItemId  string

	Name       string
	SourceCode string
	On         bool

	ScriptType      string
	ScriptPath      string
	ScriptLang      string
	ScriptCode      string
	ScriptCwd       string
	ScriptEnvNames  []string
	ScriptEnvValues []string

	WebhookURL     string
	WebhookMethod  string
	WebhookTimeout uint

	FilePath string

	DataFormat uint8
	Interval   uint

	CondParams         []string
	CondOps            []string
	CondValues         []string
	CondNoticeLevels   []uint
	CondNoticeMessages []string

	Must *actions.Must
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入监控项名称").
		Field("sourceCode", params.SourceCode).
		Require("请选择数据源类型")

	item.On = params.On
	item.Name = params.Name

	// 数据源
	item.SourceCode = params.SourceCode
	item.SourceOptions = map[string]interface{}{}

	switch params.SourceCode {
	case "script":
		if params.ScriptType == "path" {
			params.Must.
				Field("scriptPath", params.ScriptPath).
				Require("请输入脚本路径")
		} else if params.ScriptType == "code" {
			params.Must.
				Field("scriptCode", params.ScriptCode).
				Require("请输入脚本代码")
		} else {
			params.Must.
				Field("scriptPath", params.ScriptPath).
				Require("请输入脚本路径")
		}

		source := agents.NewScriptSource()
		source.ScriptType = params.ScriptType
		source.Path = params.ScriptPath
		source.ScriptLang = params.ScriptLang
		source.Script = params.ScriptCode
		source.Cwd = params.ScriptCwd
		source.DataFormat = params.DataFormat

		for index, envName := range params.ScriptEnvNames {
			if index < len(params.ScriptEnvValues) {
				source.AddEnv(envName, params.ScriptEnvValues[index])
			}
		}

		err := teautils.ObjectToMapJSON(source, &item.SourceOptions)
		if err != nil {
			logs.Error(err)
		}
	case "webhook":
		params.Must.
			Field("webhookURL", params.WebhookURL).
			Require("请输入URL").
			Match("(?i)^(http|https)://", "URL地址必须以http或https开头").
			Field("webhookMethod", params.WebhookMethod).
			Require("请选择请求方法")

		source := agents.NewWebHookSource()
		source.URL = params.WebhookURL
		source.Method = params.WebhookMethod
		source.Timeout = fmt.Sprintf("%ds", params.WebhookTimeout)
		source.DataFormat = params.DataFormat

		err := teautils.ObjectToMapJSON(source, &item.SourceOptions)
		if err != nil {
			logs.Error(err)
		}
	case "file":
		params.Must.
			Field("filePath", params.FilePath).
			Require("请输入数据文件路径")

		source := agents.NewFileSource()
		source.Path = params.FilePath
		source.DataFormat = params.DataFormat

		err := teautils.ObjectToMapJSON(source, &item.SourceOptions)
		if err != nil {
			logs.Error(err)
		}
	}

	// 刷新间隔
	item.Interval = fmt.Sprintf("%ds", params.Interval)

	// 阈值设置
	item.Thresholds = []*agents.Threshold{}
	for index, param := range params.CondParams {
		if index < len(params.CondValues) && index < len(params.CondOps) && index < len(params.CondValues) && index < len(params.CondNoticeLevels) && index < len(params.CondNoticeMessages) {
			// 校验
			op := params.CondOps[index]
			value := params.CondValues[index]
			if op == agents.ThresholdOperatorRegexp || op == agents.ThresholdOperatorNotRegexp {
				_, err := regexp.Compile(value)
				if err != nil {
					this.Fail("阈值" + param + "正则表达式" + value + "校验失败：" + err.Error())
				}
			}

			t := agents.NewThreshold()
			t.Param = param
			t.Operator = op
			t.Value = value
			t.NoticeLevel = types.Uint8(params.CondNoticeLevels[index])
			t.NoticeMessage = params.CondNoticeMessages[index]
			item.AddThreshold(t)
		}
	}

	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("UPDATE_ITEM", maps.Map{
		"appId":  app.Id,
		"itemId": params.ItemId,
	}))

	this.Success()
}
