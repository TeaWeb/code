package teawaf

import (
	"bytes"
	"github.com/TeaWeb/code/teawaf/actions"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/lists"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
)

func Test_Template(t *testing.T) {
	a := assert.NewAssertion(t)

	template := Template()
	err := template.Init()
	if err != nil {
		t.Fatal(err)
	}

	template.OnAction(func(action actions.ActionString) (goNext bool) {
		return action != actions.ActionBlock
	})

	testTemplate1001(a, t, template)
	testTemplate1002(a, t, template)
	testTemplate1003(a, t, template)
	testTemplate2001(a, t, template)
	testTemplate3001(a, t, template)
	testTemplate4001(a, t, template)
	testTemplate5001(a, t, template)
	testTemplate6001(a, t, template)
	testTemplate7001(a, t, template)
}

func testTemplate1001(a *assert.Assertion, t *testing.T, template *WAF) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com/index.php?id=onmousedown%3D123", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "1001")
	}
}

func testTemplate1002(a *assert.Assertion, t *testing.T, template *WAF) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com/index.php?id=eval%28", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "1002")
	}
}

func testTemplate1003(a *assert.Assertion, t *testing.T, template *WAF) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com/index.php?id=<script src=\"123.js\">", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "1003")
	}
}

func testTemplate2001(a *assert.Assertion, t *testing.T, template *WAF) {
	body := bytes.NewBuffer([]byte{})

	writer := multipart.NewWriter(body)

	{
		part, err := writer.CreateFormField("name")
		if err == nil {
			part.Write([]byte("lu"))
		}
	}

	{
		part, err := writer.CreateFormField("age")
		if err == nil {
			part.Write([]byte("20"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile", "hello.txt")
		if err == nil {
			part.Write([]byte("Hello, World!"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile2", "hello.PHP")
		if err == nil {
			part.Write([]byte("Hello, World, PHP!"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile3", "hello.asp")
		if err == nil {
			part.Write([]byte("Hello, World, ASP Pages!"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile4", "hello.asp")
		if err == nil {
			part.Write([]byte("Hello, World, ASP Pages!"))
		}
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "http://teaos.cn/", body)
	if err != nil {
		t.Fatal()
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	_, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "2001")
	}
}

func testTemplate3001(a *assert.Assertion, t *testing.T, template *WAF) {
	req, err := http.NewRequest(http.MethodPost, "http://example.com/index.php?exec1+(", bytes.NewReader([]byte("exec('rm -rf /hello');")))
	if err != nil {
		t.Fatal(err)
	}
	_, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "3001")
	}
}

func testTemplate4001(a *assert.Assertion, t *testing.T, template *WAF) {
	req, err := http.NewRequest(http.MethodPost, "http://example.com/index.php?whoami", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "4001")
	}
}

func testTemplate5001(a *assert.Assertion, t *testing.T, template *WAF) {
	{
		req, err := http.NewRequest(http.MethodPost, "http://example.com/.././..", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNotNil(result)
		if result != nil {
			a.IsTrue(result.Code == "5001")
		}
	}

	{
		req, err := http.NewRequest(http.MethodPost, "http://example.com/..///./", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNotNil(result)
		if result != nil {
			a.IsTrue(result.Code == "5001")
		}
	}
}

func testTemplate6001(a *assert.Assertion, t *testing.T, template *WAF) {
	{
		req, err := http.NewRequest(http.MethodPost, "http://example.com/.svn/123.txt", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNotNil(result)
		if result != nil {
			a.IsTrue(result.Code == "6001")
		}
	}

	{
		req, err := http.NewRequest(http.MethodPost, "http://example.com/123.git", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNil(result)
	}
}

func testTemplate7001(a *assert.Assertion, t *testing.T, template *WAF) {
	for _, id := range []string{
		"union select",
		" and if(",
		"/*!",
		" and select ",
		" and id=123 ",
		"(case when a=1 then ",
		"updatexml (",
		"; delete from table",
	} {
		req, err := http.NewRequest(http.MethodPost, "http://example.com/?id="+url.QueryEscape(id), nil)
		if err != nil {
			t.Fatal(err)
		}
		_, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNotNil(result)
		if result != nil {
			a.IsTrue(lists.ContainsAny([]string{"7001", "7002", "7003", "7004", "7005"}, result.Code))
		} else {
			t.Log("break:", id)
		}
	}
}
