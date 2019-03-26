package teaproxy

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
)

// 调用API Mock
func (this *Request) callMock(writer *ResponseWriter) error {
	if this.api != nil && len(this.api.MockFiles) > 0 {
		mock := this.api.RandMock()
		if mock != nil {
			for _, header := range mock.Headers {
				name := header.GetString("name")
				value := header.GetString("value")
				if len(name) > 0 {
					writer.Header().Set(name, value)
				}
			}

			writer.Header().Set("Tea-API-Mock", "on")

			if len(mock.File) > 0 {
				reader, err := files.NewReader(Tea.ConfigFile(mock.File))
				if err == nil {
					defer reader.Close()
					data := reader.ReadAll()
					writer.Write(data)
				}
			} else {
				writer.Write([]byte(mock.Text))
			}
		}
	} else {
		writer.Write([]byte("mock data not found"))
	}

	return nil
}
