package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestLocationConfig_Match(t *testing.T) {
	location := NewLocation()
	err := location.Validate()
	if err != nil {
		t.Fatal(err)
	}

	a := assert.NewAssertion(t).Quiet()

	location.Pattern = "/hell"
	a.IsNotError(location.Validate())

	_, b := location.Match("/hello")
	a.IsTrue(b)

	location.Pattern = "/hello"
	a.IsNotError(location.Validate())

	_, b = location.Match("/hello")
	a.IsTrue(b)

	location.Pattern = "~ ^/\\w+$"
	a.IsNotError(location.Validate())
	_, b = location.Match("/hello")
	a.IsTrue(b)

	location.Pattern = "!~ ^/HELLO$"
	a.IsNotError(location.Validate())
	_, b = location.Match("/hello")
	a.IsTrue(b)

	location.Pattern = "~* ^/HELLO$"
	a.IsNotError(location.Validate())

	_, b = location.Match("/hello")
	a.IsTrue(b)

	location.Pattern = "!~* ^/HELLO$"
	a.IsNotError(location.Validate())
	_, b = location.Match("/hello")
	a.IsFalse(b)

	location.Pattern = "= /hello"
	a.IsNotError(location.Validate())
	_, b = location.Match("/hello")
	a.IsTrue(b)
}
