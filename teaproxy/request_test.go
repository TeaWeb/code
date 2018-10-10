package teaproxy

import (
	"testing"
	"github.com/iwind/TeaGo/assert"
	"net/http"
	"github.com/iwind/TeaGo/Tea"
	"bytes"
	"github.com/TeaWeb/code/teaconfigs"
)

type testResponseWriter struct {
	a    *assert.Assertion
	data []byte
}

func testNewResponseWriter(a *assert.Assertion) *testResponseWriter {
	return &testResponseWriter{
		a: a,
	}
}

func (this *testResponseWriter) Header() http.Header {
	return http.Header{}
}

func (this *testResponseWriter) Write(data []byte) (int, error) {
	this.data = append(this.data, data ...)
	return len(data), nil
}

func (this *testResponseWriter) WriteHeader(statusCode int) {
}

func (this *testResponseWriter) Close() {
	this.a.Log(string(this.data))
}

func TestRequest_Call(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	request := NewRequest(nil)
	err := request.Call(writer)
	a.IsNotNil(err)
	if err != nil {
		a.Log(err.Error())
	}
}

func TestRequest_CallRoot(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	request := NewRequest(nil)
	request.root = Tea.ViewsDir() + "/@default"
	request.uri = "/layout.css"
	err := request.Call(writer)
	a.IsNil(err)
	writer.Close()

	a.Log("requestTime:", request.requestTime)
	a.Log("bytes send:", request.responseBytesSent, request.responseBodyBytesSent)
}

func TestRequest_CallBackend(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx", nil)
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"
	request.backend = &teaconfigs.ServerBackendConfig{
		Address: "127.0.0.1",
	}
	request.backend.Validate()
	err = request.Call(writer)
	a.IsNil(err)
	writer.Close()

	a.Log("status:", request.responseStatus, request.responseStatusMessage)
	a.Log("requestTime:", request.requestTime)
	a.Log("bytes send:", request.responseBytesSent, request.responseBodyBytesSent)
}

func TestRequest_CallProxy(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx", nil)
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"

	proxy := teaconfigs.NewServerConfig()
	proxy.AddBackend(&teaconfigs.ServerBackendConfig{
		Address: "127.0.0.1:80",
	})
	/**proxy.AddBackend(&teaconfigs.ServerBackendConfig{
		Address: "127.0.0.1:81",
	})**/
	request.proxy = proxy

	err = request.Call(writer)
	a.IsNil(err)
	writer.Close()

	a.Log("status:", request.responseStatus, request.responseStatusMessage)
	a.Log("requestTime:", request.requestTime)
	a.Log("bytes send:", request.responseBytesSent, request.responseBodyBytesSent)
}

func TestRequest_CallFastcgi(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx/box/version", bytes.NewBuffer([]byte("hello=world")))
	//req, err := http.NewRequest("GET", "/index.php", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"
	request.serverAddr = "127.0.0.1:80"

	request.fastcgi = &teaconfigs.FastcgiConfig{
		Params: map[string]string{
			"SCRIPT_FILENAME": "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk/index.php",
			//"DOCUMENT_ROOT":   "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk",
		},
		Pass: "127.0.0.1:9000",
	}
	request.fastcgi.Validate()
	err = request.Call(writer)
	a.IsNil(err)
	writer.Close()

	a.Log("status:", request.responseStatus, request.responseStatusMessage)
	a.Log("requestTime:", request.requestTime)
	a.Log("bytes send:", request.responseBytesSent, request.responseBodyBytesSent)
}

func TestRequest_Format(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}
	rawReq.RemoteAddr = "127.0.0.1:1234"
	rawReq.Header.Add("Content-Type", "text/plain")

	req := NewRequest(rawReq)
	req.uri = "/hello/world?name=Lu&age=20"
	req.method = "GET"
	req.filePath = "hello.go"
	req.scheme = "http"

	a.IsTrue(req.requestRemoteAddr() == "127.0.0.1:1234")
	a.IsTrue(req.requestRemotePort() == 1234)
	a.IsTrue(req.requestURI() == req.uri)
	a.IsTrue(req.requestPath() == "/hello/world")
	a.IsTrue(req.requestMethod() == "GET")
	a.IsTrue(req.requestLength() > 0)
	a.IsTrue(req.requestFilename() == req.filePath)
	a.IsTrue(req.requestProto() == "HTTP/1.1")
	a.IsTrue(req.requestQueryString() == "name=Lu&age=20")
	a.IsTrue(req.requestQueryParam("name") == "Lu")

	t.Log(req.format("hello ${teaVersion} remoteAddr:${remoteAddr} name:${arg.name} header:${header.Content-Type} test:${test}"))
}

func TestRequest_Index(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}

	req := NewRequest(rawReq)
	req.index = []string{}
	t.Log(req.findIndexFile(Tea.Root))

	req.index = []string{"main.go", "main2.go", "run.sh"}
	a.Equals(req.findIndexFile(Tea.Root), "main.go")

	req.index = []string{"main.*"}
	a.Equals(req.findIndexFile(Tea.Root), "main.go")
}
