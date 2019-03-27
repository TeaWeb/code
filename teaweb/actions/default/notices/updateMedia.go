package notices

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type UpdateMediaAction actions.Action

// 修改媒介
func (this *UpdateMediaAction) Run(params struct {
	MediaId string
	From    string
}) {
	setting := notices.SharedNoticeSetting()
	media := setting.FindMedia(params.MediaId)
	if media == nil {
		this.Fail("找不到Media")
	}

	this.Data["from"] = params.From
	this.Data["media"] = media
	this.Data["mediaTypes"] = notices.AllNoticeMediaTypes()
	this.Data["methods"] = []string{http.MethodGet, http.MethodPost}

	this.Show()
}

// 提交修改
func (this *UpdateMediaAction) RunPost(params struct {
	MediaId string

	Name string
	Type string
	On   bool

	EmailSmtp     string
	EmailUsername string
	EmailPassword string
	EmailFrom     string

	WebhookURL          string
	WebhookMethod       string
	WebhookHeaderNames  []string
	WebhookHeaderValues []string
	WebhookContentType  string
	WebhookParamNames   []string
	WebhookParamValues  []string
	WebhookBody         string

	ScriptType      string
	ScriptPath      string
	ScriptLang      string
	ScriptCode      string
	ScriptCwd       string
	ScriptEnvNames  []string
	ScriptEnvValues []string

	DingTalkWebhookURL string

	QyWeixinCorporateId string
	QyWeixinAgentId     string
	QyWeixinAppSecret   string

	TimeFromHour   int
	TimeFromMinute int
	TimeFromSecond int
	TimeToHour     int
	TimeToMinute   int
	TimeToSecond   int
	RateCount      int
	RateMinutes    int

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入媒介名称")

	if notices.FindNoticeMediaType(params.Type) == nil {
		this.Fail("找不到此媒介类型")
	}

	setting := notices.SharedNoticeSetting()
	mediaConfig := setting.FindMedia(params.MediaId)
	if mediaConfig == nil {
		this.Fail("找不到Media")
	}

	mediaConfig.Name = params.Name
	mediaConfig.Type = params.Type
	mediaConfig.On = params.On

	switch params.Type {
	case notices.NoticeMediaTypeEmail:
		params.Must.
			Field("emailSmtp", params.EmailSmtp).
			Require("请输入SMTP地址").
			Field("emailUsername", params.EmailUsername).
			Require("请输入邮箱账号").
			Field("emailPassword", params.EmailPassword).
			Require("请输入密码或授权码")

		media := notices.NewNoticeEmailMedia()
		media.SMTP = params.EmailSmtp
		media.Username = params.EmailUsername
		media.Password = params.EmailPassword
		media.From = params.EmailFrom
		teautils.ObjectToMapJSON(media, &mediaConfig.Options)
	case notices.NoticeMediaTypeWebhook:
		params.Must.
			Field("webhookURL", params.WebhookURL).
			Require("请输入URL地址").
			Match("(?i)^(http|https)://", "URL地址必须以http或https开头").
			Field("webhookMethod", params.WebhookMethod).
			Require("请选择请求方法")

		media := notices.NewNoticeWebhookMedia()
		media.URL = params.WebhookURL
		media.Method = params.WebhookMethod

		media.ContentType = params.WebhookContentType
		if len(params.WebhookHeaderNames) > 0 {
			for index, name := range params.WebhookHeaderNames {
				if index < len(params.WebhookHeaderValues) {
					media.AddHeader(name, params.WebhookHeaderValues[index])
				}
			}
		}

		if params.WebhookContentType == "params" {
			for index, name := range params.WebhookParamNames {
				if index < len(params.WebhookParamValues) {
					media.AddParam(name, params.WebhookParamValues[index])
				}
			}
		} else if params.WebhookContentType == "body" {
			media.Body = params.WebhookBody
		}
		teautils.ObjectToMapJSON(media, &mediaConfig.Options)
	case notices.NoticeMediaTypeScript:
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

		media := notices.NewNoticeScriptMedia()
		media.ScriptType = params.ScriptType
		media.Path = params.ScriptPath
		media.ScriptLang = params.ScriptLang
		media.Script = params.ScriptCode
		media.Cwd = params.ScriptCwd

		for index, envName := range params.ScriptEnvNames {
			if index < len(params.ScriptEnvValues) {
				media.AddEnv(envName, params.ScriptEnvValues[index])
			}
		}

		teautils.ObjectToMapJSON(media, &mediaConfig.Options)
	case notices.NoticeMediaTypeDingTalk:
		params.Must.
			Field("dingTalkWebhookURL", params.DingTalkWebhookURL).
			Require("请输入Hook地址").
			Match("^https:", "Hook地址必须以https://开头")

		media := notices.NewNoticeDingTalkMedia()
		media.WebhookURL = params.DingTalkWebhookURL
		teautils.ObjectToMapJSON(media, &mediaConfig.Options)
	case notices.NoticeMediaTypeQyWeixin:
		params.Must.
			Field("qyWeixinCorporateId", params.QyWeixinCorporateId).
			Require("请输入企业ID").
			Field("qyWeixinAgentId", params.QyWeixinAgentId).
			Require("请输入应用AgentId").
			Field("qyWeixinSecret", params.QyWeixinAppSecret).
			Require("请输入应用Secret")

		media := notices.NewNoticeQyWeixinMedia()
		media.CorporateId = params.QyWeixinCorporateId
		media.AgentId = params.QyWeixinAgentId
		media.AppSecret = params.QyWeixinAppSecret
		teautils.ObjectToMapJSON(media, &mediaConfig.Options)
	}

	// 时间
	params.Must.
		Field("timeFromHour", params.TimeFromHour).
		Require("请输入正确的小时数").
		Gte(0, "请输入正确的小时数").
		Lte(23, "请输入正确的小时数").
		Field("timeFromMinute", params.TimeFromMinute).
		Require("请输入正确的分钟数").
		Gte(0, "请输入正确的分钟数").
		Lte(59, "请输入正确的分钟数").
		Field("timeFromSecond", params.TimeFromSecond).
		Require("请输入正确的秒数").
		Gte(0, "请输入正确的秒数").
		Lte(59, "请输入正确的秒数").

		Field("timeToHour", params.TimeToHour).
		Require("请输入正确的小时数").
		Gte(0, "请输入正确的小时数").
		Lte(23, "请输入正确的小时数").
		Field("timeToMinute", params.TimeToMinute).
		Require("请输入正确的分钟数").
		Gte(0, "请输入正确的分钟数").
		Lte(59, "请输入正确的分钟数").
		Field("timeToSecond", params.TimeToSecond).
		Require("请输入正确的秒数").
		Gte(0, "请输入正确的秒数").
		Lte(59, "请输入正确的秒数")

	mediaConfig.TimeFrom = fmt.Sprintf("%02d:%02d:%02d", params.TimeFromHour, params.TimeFromMinute, params.TimeFromSecond)
	mediaConfig.TimeTo = fmt.Sprintf("%02d:%02d:%02d", params.TimeToHour, params.TimeToMinute, params.TimeToSecond)
	mediaConfig.RateCount = params.RateCount
	mediaConfig.RateMinutes = params.RateMinutes

	err := setting.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
