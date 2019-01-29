package notices

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type AddMediaAction actions.Action

// 添加媒介
func (this *AddMediaAction) Run(params struct{}) {
	this.Data["mediaTypes"] = notices.AllNoticeMediaTypes()
	this.Data["methods"] = []string{http.MethodGet, http.MethodPost}

	this.Show()
}

// 提交保存
func (this *AddMediaAction) RunPost(params struct {
	Name string
	Type string
	On   bool

	EmailSmtp     string
	EmailUsername string
	EmailPassword string
	EmailFrom     string

	WebhookURL    string
	WebhookMethod string

	ScriptType      string
	ScriptPath      string
	ScriptLang      string
	ScriptCode      string
	ScriptCwd       string
	ScriptEnvNames  []string
	ScriptEnvValues []string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入媒介名称")

	if notices.FindNoticeMediaType(params.Type) == nil {
		this.Fail("找不到此媒介类型")
	}

	mediaConfig := notices.NewNoticeMediaConfig()
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
	}

	setting := notices.SharedNoticeSetting()
	setting.AddMedia(mediaConfig)
	err := setting.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
