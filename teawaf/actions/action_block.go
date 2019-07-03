package actions

import (
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// url client configure
var urlPrefixReg = regexp.MustCompile("^(?i)(http|https)://")
var httpClient = teautils.NewHttpClient(5 * time.Second)

type BlockAction struct {
	StatusCode int    `yaml:"statusCode" json:"statusCode"`
	Body       string `yaml:"body" json:"body"` // supports HTML
	URL        string `yaml:"url" json:"url"`
}

func (this *BlockAction) Perform(writer http.ResponseWriter) (allow bool) {
	if writer != nil {
		if this.StatusCode > 0 {
			writer.WriteHeader(this.StatusCode)
		} else {
			writer.WriteHeader(http.StatusForbidden)
		}
		if len(this.URL) > 0 {
			if urlPrefixReg.MatchString(this.URL) {
				req, err := http.NewRequest(http.MethodGet, this.URL, nil)
				if err != nil {
					logs.Error(err)
					return false
				}
				resp, err := httpClient.Do(req)
				if err != nil {
					logs.Error(err)
					return false
				}
				defer resp.Body.Close()

				for k, v := range resp.Header {
					for _, v1 := range v {
						writer.Header().Add(k, v1)
					}
				}

				io.Copy(writer, resp.Body)
			} else {
				path := this.URL
				if !filepath.IsAbs(this.URL) {
					path = Tea.Root + string(os.PathSeparator) + path
				}

				data, err := ioutil.ReadFile(path)
				if err != nil {
					logs.Error(err)
					return false
				}
				writer.Write(data)
			}
			return false
		}
		if len(this.Body) > 0 {
			writer.Write([]byte(this.Body))
		} else {
			writer.Write([]byte("The request is blocked by TeaWAF"))
		}
	}
	return false
}
