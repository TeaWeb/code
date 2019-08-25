package teadb

import (
	"github.com/iwind/TeaGo/logs"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"testing"
)

func TestAccessLogDAO_FindAccessLogCookie(t *testing.T) {
	dao := SharedDB().AccessLogDAO()
	accessLog, err := dao.FindAccessLogCookie("20190608", "5cfbbecd79c023a965148da9")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(stringutil.JSONEncodePretty(accessLog.Cookie))
}

func TestAccessLogDAO_FindResponseHeaderAndBody(t *testing.T) {
	dao := SharedDB().AccessLogDAO()
	accessLog, err := dao.FindResponseHeaderAndBody("20190608", "5cfbbecd79c023a965148da9")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(accessLog.SentHeader)
	t.Log(accessLog.ResponseBodyData)
	logs.PrintAsJSON(accessLog, t)
}

func TestAccessLogDAO_FindRequestHeaderAndBody(t *testing.T) {
	dao := SharedDB().AccessLogDAO()
	accessLog, err := dao.FindRequestHeaderAndBody("20190608", "5cfbbecd79c023a965148da9")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(accessLog.Header)
	if len(accessLog.RequestData) == 0 {
		t.Log(accessLog.RequestData)
	} else {
		t.Log(string(accessLog.RequestData))
	}
	logs.PrintAsJSON(accessLog, t)
}

func TestAccessLogDAO_ListAccessLogs(t *testing.T) {
	dao := SharedDB().AccessLogDAO()
	accessLogs, err := dao.ListAccessLogs("20190608", "5W8NLAoMYo6iJ78V", "5cfbc98141a7eae69097db99", true, "", 0, 5)
	if err != nil {
		t.Fatal(err)
	}

	for _, accessLog := range accessLogs {
		t.Log(accessLog.Id, accessLog.ServerId, accessLog.Errors, accessLog.RemoteAddr)
	}
}

func TestAccessLogDAO_HasNextAccessLog(t *testing.T) {
	dao := SharedDB().AccessLogDAO()
	b, err := dao.HasNextAccessLog("20190608", "5W8NLAoMYo6iJ78V", "5cfbbc918e6b5df25169a432", false, "")
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Log("has next")
	} else {
		t.Log("has no next")
	}
}

func TestAccessLogDAO_ListLatestAccessLogs(t *testing.T) {
	dao := SharedDB().AccessLogDAO()
	accessLogs, err := dao.ListLatestAccessLogs("20190608", "5W8NLAoMYo6iJ78V", "5cfbc98141a7eae69097db95", false, 5)
	if err != nil {
		t.Fatal(err)
	}
	for _, accessLog := range accessLogs {
		t.Log(accessLog.Id, accessLog.ServerId, accessLog.Errors, accessLog.RemoteAddr)
	}
}

func TestAccessLogDAO_QueryAccessLogs(t *testing.T) {
	dao := SharedDB().AccessLogDAO()

	query := NewQuery("")
	query.Limit(5)
	query.Debug()

	accessLogs, err := dao.QueryAccessLogs("20190608", "5W8NLAoMYo6iJ78V", query)
	if err != nil {
		t.Fatal(err)
	}
	for _, accessLog := range accessLogs {
		t.Log(accessLog.Id, accessLog.ServerId, accessLog.Errors, accessLog.RemoteAddr)
	}
}
