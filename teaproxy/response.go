package teaproxy

import (
	"bytes"
	"net/http"
)

// 响应Writer
type ResponseWriter struct {
	writer http.ResponseWriter

	statusCode    int
	sentBodyBytes int64

	bodyCopying bool
	body        []byte
}

// 包装对象
func NewResponseWriter(httpResponseWriter http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		writer: httpResponseWriter,
	}
}

// 获取Header
func (this *ResponseWriter) Header() http.Header {
	return this.writer.Header()
}

// 添加一组Header
func (this *ResponseWriter) AddHeaders(header http.Header) {
	for key, value := range header {
		for _, v := range value {
			this.writer.Header().Add(key, v)
		}
	}
}

// 写入数据
func (this *ResponseWriter) Write(data []byte) (n int, err error) {
	n, err = this.writer.Write(data)
	if n > 0 {
		this.sentBodyBytes += int64(n)
	}
	if this.bodyCopying {
		this.body = append(this.body, data ...)
	}
	return
}

// 读取发送的字节数
func (this *ResponseWriter) SentBodyBytes() int64 {
	return this.sentBodyBytes
}

// 写入状态码
func (this *ResponseWriter) WriteHeader(statusCode int) {
	this.writer.WriteHeader(statusCode)
	this.statusCode = statusCode
}

// 读取状态码
func (this *ResponseWriter) StatusCode() int {
	if this.statusCode == 0 {
		return http.StatusOK
	}
	return this.statusCode
}

// 设置拷贝Body数据
func (this *ResponseWriter) SetBodyCopying(b bool) {
	this.bodyCopying = b
}

// 判断是否在拷贝Body数据
func (this *ResponseWriter) BodyIsCopying() bool {
	return this.bodyCopying
}

// 读取拷贝的Body数据
func (this *ResponseWriter) Body() []byte {
	return this.body
}

// 读取Header二进制数据
func (this *ResponseWriter) HeaderData() []byte {
	resp := &http.Response{}
	resp.Header = this.Header()
	if this.statusCode == 0 {
		this.statusCode = http.StatusOK
	}
	resp.StatusCode = this.statusCode
	resp.ProtoMajor = 1
	resp.ProtoMinor = 1

	resp.ContentLength = 1 // Trick：这样可以屏蔽Content-Length

	writer := bytes.NewBuffer([]byte{})
	resp.Write(writer)
	return writer.Bytes()
}
