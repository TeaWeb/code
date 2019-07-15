package sslutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaproxy"
	"strings"
	"testing"
)

func TestReloadACMECert(t *testing.T) {
	server := &teaconfigs.ServerConfig{
		Id: "abc",
	}
	server.SSL = teaconfigs.NewSSLConfig()
	server.SSL.Certs = []*teaconfigs.SSLCertConfig{
		{
			On:       true,
			TaskId:   "123",
			CertFile: "123.pem",
			KeyFile:  "456.key",
		},
	}
	teaproxy.SharedManager.ApplyServer(server)
	errs := ReloadACMECert("abc", "123")
	for _, err := range errs {
		if strings.Contains(err.Error(), "failed:open") { // 符合预期
			continue
		}
		t.Fatal(err)
	}
}
