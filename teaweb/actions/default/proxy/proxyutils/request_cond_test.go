package proxyutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"net/url"
	"testing"
)

func TestParseRequestConds(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Form = url.Values{}
	{
		req.Form.Add("request_condParams", "${requestURI}")
		req.Form.Add("request_condOperators", "eq")
		req.Form.Add("request_condValues", "/hello")
	}
	{
		req.Form.Add("request_condParams", "${arg.age}")
		req.Form.Add("request_condOperators", "gt")
		req.Form.Add("request_condValues", "20")
	}
	{
		req.Form.Add("request_condParams", "${arg.name}")
		req.Form.Add("request_condOperators", teaconfigs.RequestCondOperatorRegexp)
		req.Form.Add("request_condValues", "\\w+")
	}
	t.Log(req.Form)

	conds, breakCond, err := ParseRequestConds(req, "request")
	if err != nil {
		t.Log(breakCond)
		t.Fatal(err)
	}

	logs.PrintAsJSON(conds, t)
}
