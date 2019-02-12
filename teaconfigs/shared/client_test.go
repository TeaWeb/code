package shared

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestClientConfig_Allow(t *testing.T) {
	a := assert.NewAssertion(t)

	client := NewClientConfig()
	client.Validate()
	a.IsFalse(client.Match("127.0.0.1"))

	client.IP = "192.168.1.100"
	client.Validate()
	a.IsTrue(client.Match("192.168.1.100"))
	a.IsFalse(client.Match("192.168.1.101"))

	client.IP = "192.168.1.*"
	client.Validate()
	a.IsTrue(client.Match("192.168.1.100"))
	a.IsTrue(client.Match("192.168.1.101"))
	a.IsFalse(client.Match("192.168.2.100"))

	client.IP = "192.168.*.*"
	client.Validate()
	a.IsTrue(client.Match("192.168.1.100"))
	a.IsTrue(client.Match("192.168.1.101"))
	a.IsTrue(client.Match("192.168.2.100"))
}
