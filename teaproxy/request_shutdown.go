package teaproxy

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"io"
	"net/http"
	"os"
)

// 调用临时关闭页面
func (this *Request) callShutdown(writer *ResponseWriter) error {
	if urlPrefixRegexp.MatchString(this.shutdownPage) {
		return this.callURL(writer, http.MethodGet, this.shutdownPage)
	} else {
		file := Tea.Root + Tea.DS + this.shutdownPage
		fp, err := os.Open(file)
		if err != nil {
			logs.Error(err)
			msg := "404 page not found: '" + this.shutdownPage + "'"

			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte(msg))
			return err
		}

		writer.WriteHeader(http.StatusOK)
		_, err = io.Copy(writer, fp)
		fp.Close()

		return err
	}
}
