package notices

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/iwind/TeaGo/logs"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Webhook媒介
type NoticeWebhookMedia struct {
	URL    string `yaml:"url" json:"url"` // URL中可以使用${NoticeSubject}, ${NoticeBody}两个变量
	Method string `yaml:"method" json:"method"`
}

// 获取新对象
func NewNoticeWebhookMedia() *NoticeWebhookMedia {
	return &NoticeWebhookMedia{}
}

// 发送
func (this *NoticeWebhookMedia) Send(user string, subject string, body string) (resp []byte, err error) {
	if len(this.URL) == 0 {
		return nil, errors.New("'url' should be specified")
	}

	timeout := 10 * time.Second

	if len(this.Method) == 0 {
		this.Method = http.MethodGet
	}

	this.URL = strings.Replace(this.URL, "${NoticeUser}", url.QueryEscape(user), -1)
	this.URL = strings.Replace(this.URL, "${NoticeSubject}", url.QueryEscape(subject), -1)
	this.URL = strings.Replace(this.URL, "${NoticeBody}", url.QueryEscape(body), -1)

	var req *http.Request
	if this.Method == http.MethodGet {
		logs.Println(this.URL)
		req, err = http.NewRequest(this.Method, this.URL, nil)
	} else {
		params := url.Values{
			"NoticeUser":    []string{user},
			"NoticeSubject": []string{subject},
			"NoticeBody":    []string{body},
		}
		req, err = http.NewRequest(this.Method, this.URL, strings.NewReader(params.Encode()))
	}
	req.Header.Set("User-Agent", "TeaWeb/"+teaconst.TeaVersion)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: timeout,
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	return data, err
}
