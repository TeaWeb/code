package requests

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type Request struct {
	*http.Request
	BodyData []byte
}

func NewRequest(raw *http.Request) *Request {
	return &Request{
		Request: raw,
	}
}

func (this *Request) ReadBody(max int64) (data []byte, err error) {
	data, err = ioutil.ReadAll(io.LimitReader(this.Request.Body, max))
	return
}

func (this *Request) RestoreBody(data []byte) {
	rawReader := bytes.NewBuffer(data)
	io.Copy(rawReader, this.Request.Body)
	this.Request.Body = ioutil.NopCloser(rawReader)
}
