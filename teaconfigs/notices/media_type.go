package notices

import "github.com/iwind/TeaGo/maps"

// 通知媒介类型
type NoticeMediaType = string

const (
	NoticeMediaTypeEmail   = "email"
	NoticeMediaTypeWebhook = "webhook"
	NoticeMediaTypeScript  = "script"
)

// 所有媒介
func AllNoticeMediaTypes() []maps.Map {
	return []maps.Map{
		{
			"name":         "邮件",
			"code":         NoticeMediaTypeEmail,
			"supportsHTML": true,
			"instance":     new(NoticeEmailMedia),
			"description":  "通过邮件发送通知",
			"user":         "接收人邮箱地址",
		},
		{
			"name":         "Webhook",
			"code":         NoticeMediaTypeWebhook,
			"supportsHTML": false,
			"instance":     new(NoticeWebhookMedia),
			"description":  "通过HTTP请求发送通知",
			"user":         "通过${NoticeUser}参数传递到URL上",
		},
		{
			"name":         "脚本",
			"code":         NoticeMediaTypeScript,
			"supportsHTML": false,
			"instance":     new(NoticeScriptMedia),
			"description":  "通过运行脚本发送通知",
			"user":         "可以在脚本中使用${NoticeUser}来获取这个标识",
		},
	}
}

// 查找媒介类型
func FindNoticeMediaType(mediaType string) maps.Map {
	for _, m := range AllNoticeMediaTypes() {
		if m["code"] == mediaType {
			return m
		}
	}
	return nil
}

// 查找媒介类型名称
func FindNoticeMediaTypeName(mediaType string) string {
	m := FindNoticeMediaType(mediaType)
	if m == nil {
		return ""
	}
	return m["name"].(string)
}

// 媒介接口
type NoticeMediaInterface interface {
	Send(user string, subject string, body string) (resp []byte, err error)
}
