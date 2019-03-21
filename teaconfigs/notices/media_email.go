package notices

import (
	"errors"
	"net/smtp"
	"strings"
)

// 邮件媒介
type NoticeEmailMedia struct {
	SMTP     string `yaml:"smtp" json:"smtp"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	From     string `yaml:"from" json:"from"`
}

// 获取新对象
func NewNoticeEmailMedia() *NoticeEmailMedia {
	return &NoticeEmailMedia{}
}

func (this *NoticeEmailMedia) Send(user string, subject string, body string) (resp []byte, err error) {
	if len(this.SMTP) == 0 {
		return nil, errors.New("host address should be specified")
	}

	// 自动加端口
	if strings.Index(this.SMTP, ":") < 0 {
		this.SMTP += ":587"
	}

	if len(this.From) == 0 {
		this.From = this.Username
	}

	hostIndex := strings.Index(this.SMTP, ":")
	auth := smtp.PlainAuth("", this.Username, this.Password, this.SMTP[:hostIndex])

	contentType := "Content-Type: text/html; charset=UTF-8"

	msg := []byte("To: " + user + "\r\nFrom: \"TeaWeb\" <" + this.From + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)

	return nil, smtp.SendMail(this.SMTP, auth, this.From, []string{user}, msg)
}

// 是否需要用户标识
func (this *NoticeEmailMedia) RequireUser() bool {
	return true
}
