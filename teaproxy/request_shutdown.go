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
		return this.callURL(writer, http.MethodGet, this.shutdownPage, "")
	} else {
		file := Tea.Root + Tea.DS + this.shutdownPage
		fp, err := os.Open(file)
		if err != nil {
			logs.Error(err)
			msg := "404 page not found: '" + this.shutdownPage + "'"

			writer.WriteHeader(http.StatusNotFound)
			_, err = writer.Write([]byte(msg))
			if err != nil {
				logs.Error(err)
			}
			return err
		}

		// 自定义响应Headers
		this.WriteResponseHeaders(writer, http.StatusOK)

		writer.WriteHeader(http.StatusOK)
		buf := bytePool1k.Get()
		_, err = io.CopyBuffer(writer, fp, buf)
		bytePool1k.Put(buf)
		err = fp.Close()
		if err != nil {
			logs.Error(err)
		}

		return err
	}
}
