package teaconfigs

import (
	"testing"
	"github.com/iwind/TeaGo/assert"
)

func TestLocationConfig_Match(t *testing.T) {
	location := NewLocationConfig()
	err := location.Validate()
	if err != nil {
		t.Fatal(err)
	}

	a := assert.NewAssertion(t).Quiet()

	location.Pattern = "/hell"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))

	location.Pattern = "/hello"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))

	location.Pattern = "~ ^/\\w+$"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))

	location.Pattern = "!~ ^/HELLO$"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))

	location.Pattern = "~* ^/HELLO$"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))

	location.Pattern = "!~* ^/HELLO$"
	a.IsNotError(location.Validate())
	a.IsFalse(location.Match("/hello"))

	location.Pattern = "= /hello"
	a.IsNotError(location.Validate())
	a.IsTrue(location.Match("/hello"))
}

func TestLocationConfig_RemoveFastcgiAt(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	location := NewLocationConfig()
	location.RemoveFastcgiAt(1)
	t.Log(location.Fastcgi)

	location.AddFastcgi(&FastcgiConfig{
		Pass: "127.0.0.1:9000",
	})
	location.RemoveFastcgiAt(1)
	a.IsTrue(len(location.Fastcgi) == 1)

	location.RemoveFastcgiAt(0)
	a.IsTrue(len(location.Fastcgi) == 0)

	location.Fastcgi = []*FastcgiConfig{}
	location.AddFastcgi(&FastcgiConfig{
		Pass: "127.0.0.1:9001",
	})
	location.AddFastcgi(&FastcgiConfig{
		Pass: "127.0.0.1:9002",
	})
	location.RemoveFastcgiAt(1)
	a.IsTrue(len(location.Fastcgi) == 1)
	a.IsTrue(location.Fastcgi[0].Pass == "127.0.0.1:9001")

	location.Fastcgi = []*FastcgiConfig{}
	location.AddFastcgi(&FastcgiConfig{
		Pass: "127.0.0.1:9001",
	})
	location.AddFastcgi(&FastcgiConfig{
		Pass: "127.0.0.1:9002",
	})
	location.AddFastcgi(&FastcgiConfig{
		Pass: "127.0.0.1:9003",
	})
	location.AddFastcgi(&FastcgiConfig{
		Pass: "127.0.0.1:9004",
	})
	location.RemoveFastcgiAt(1)
	for _, fastcgi := range location.Fastcgi {
		t.Log("fastcgi left:", fastcgi.Pass)
	}
	a.IsTrue(len(location.Fastcgi) == 3)
	a.IsTrue(location.Fastcgi[0].Pass == "127.0.0.1:9001")
}
