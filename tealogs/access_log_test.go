package tealogs

import (
	"testing"
	"github.com/ua-parser/uap-go/uaparser"
	"log"
	"github.com/iwind/TeaGo/Tea"
	"time"
	"github.com/TeaWeb/code/teamongo"
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

func TestLogParseMimeType(t *testing.T) {
	accessLog := &AccessLog{
		ContentType: "text/html;charset=utf-8",
	}
	accessLog.parseMime()
	t.Log(accessLog.Extend.File.MimeType)
	if accessLog.Extend.File.MimeType != "text/html" {
		t.Error("[ERROR]")
	}

	if accessLog.Extend.File.Charset != "UTF-8" {
		t.Error("[ERROR]", accessLog.Extend.File.Charset)
	}

	accessLog.ContentType = "text/html"
	accessLog.parseMime()
	if accessLog.Extend.File.MimeType != "text/html" {
		t.Error("[ERROR]", accessLog.Extend.File.MimeType)
	}

	if accessLog.Extend.File.Charset != "" {
		t.Error("[ERROR]", accessLog.Extend.File.Charset)
	}

	accessLog.ContentType = "text/html; CHARSET=UTF-8"
	accessLog.parseMime()
	if accessLog.Extend.File.Charset != "UTF-8" {
		t.Error("[ERROR]", accessLog.Extend.File.Charset)
	}
}

func TestLogParseExtension(t *testing.T) {
	accessLog := &AccessLog{
		RequestPath: "/hello",
	}

	accessLog.parseExtension()
	if accessLog.Extend.File.Extension != "" {
		t.Error("[ERROR]", accessLog.Extend.File.Extension)
	}

	accessLog.RequestPath = "/hello.js"
	accessLog.parseExtension()
	if accessLog.Extend.File.Extension != "js" {
		t.Error("[ERROR]", accessLog.Extend.File.Extension)
	}

	accessLog.RequestPath = "/hello.JS"
	accessLog.parseExtension()
	if accessLog.Extend.File.Extension != "js" {
		t.Error("[ERROR]", accessLog.Extend.File.Extension)
	}

	accessLog.RequestPath = "/hello.tar.gz"
	accessLog.parseExtension()
	if accessLog.Extend.File.Extension != "gz" {
		t.Error("[ERROR]", accessLog.Extend.File.Extension)
	}
}

func TestLogOSParser1(t *testing.T) {
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.59 Safari/537.36"
	parser, err := uaparser.New(Tea.Root + Tea.DS + "resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	client := parser.Parse(userAgent)

	t.Log(client.UserAgent.Family) // "Amazon Silk"
	t.Log(client.UserAgent.Major)  // "1"
	t.Log(client.UserAgent.Minor)  // "1"
	t.Log(client.UserAgent.Patch)  // "0-80"
	t.Log(client.Os.Family)        // "Android"
	t.Log(client.Os.Major)         // ""
	t.Log(client.Os.Minor)         // ""
	t.Log(client.Os.Patch)         // ""
	t.Log(client.Os.PatchMinor)    // ""
	t.Log(client.Device.Family)    // "Kindle Fire"

	cost := float64(time.Since(now).Nanoseconds()) / 1000000000
	t.Log("cost:", cost)
	t.Log("QPS", 1/cost)
}

func TestLogOSParser2(t *testing.T) {
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.59 Safari/537.36"
	parser, err := uaparser.New(Tea.Root + Tea.DS + "resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	agent := parser.ParseUserAgent(userAgent)
	t.Logf("%#v", agent)

	cost := float64(time.Since(now).Nanoseconds()) / 1000000000
	t.Log("cost:", cost)
	t.Log("QPS", 1/cost)
}

func TestLogOSParser3(t *testing.T) {
	userAgent := " Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)"
	parser, err := uaparser.New(Tea.Root + Tea.DS + "resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	agent := parser.ParseUserAgent(userAgent)
	t.Logf("%#v", agent)

	os := parser.ParseOs(userAgent)
	t.Logf("%#v", os)

	cost := float64(time.Since(now).Nanoseconds()) / 1000000000
	t.Log("cost:", cost)
	t.Log("QPS", 1/cost)
}

func TestLogOSParser4(t *testing.T) {
	userAgent := "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"
	parser, err := uaparser.New(Tea.Root + Tea.DS + "resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	agent := parser.ParseUserAgent(userAgent)
	t.Logf("%#v", agent)

	os := parser.ParseOs(userAgent)
	t.Logf("%#v", os)

	cost := float64(time.Since(now).Nanoseconds()) / 1000000000
	t.Log("cost:", cost)
	t.Log("QPS", 1/cost)
}

func TestLogParse5(t *testing.T) {
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.59 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:60.0) Gecko/20100101 Firefox/60.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)",
		"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36",
		"",
		"hello",
	}

	beforeTime1 := time.Now()
	for _, userAgent := range userAgents {
		accessLog := &AccessLog{
			UserAgent: userAgent,
		}
		accessLog.parseUserAgent()

		t.Log("=======")
		t.Log(userAgent)
		t.Logf("%#v", accessLog.Extend.Client)
	}

	beforeTime2 := time.Now()
	for i := 0; i < 10000; i ++ {
		for _, userAgent := range userAgents {
			accessLog := &AccessLog{
				UserAgent: userAgent,
			}
			accessLog.parseUserAgent()

			//t.Log("=======")
			//t.Log(userAgent)
			//t.Logf("%#v", accessLog.Extend.Client)
		}
	}
	t.Log(float64(time.Since(beforeTime1).Nanoseconds())/1000000, "ms")
	t.Log(float64(time.Since(beforeTime2).Nanoseconds())/1000000, "ms")
}

func TestAccessLogger_DB(t *testing.T) {
	client := teamongo.SharedClient()
	if client == nil {
		t.Fatal("client=nil")
	}

	accessLog := AccessLog{
		Id:   time.Now().UnixNano(),
		Args: "a=b",
		Arg: map[string][]string{
			"name": {"liu", "lu"},
		},
		Cookie: map[string]string{
			"sid": "123456",
		},
		RemoteAddr:    "127.0.0.1",
		RemotePort:    80,
		TimeLocal:     "23/Jul/2018:22:23:35 +0800",
		TimeISO8601:   "2018-07-23T22:23:35+08:00",
		Status:        200,
		BodyBytesSent: 1048,
		Request:       "GET / HTTP/1.1",
	}

	r, err := client.
		Database("teaweb").
		Collection("accessLogs").
		InsertOne(context.Background(), accessLog)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(r)
}

func TestAccessLog_Format(t *testing.T) {
	accessLog := AccessLog{
		Args: "a=b",
		Arg: map[string][]string{
			"name": {"liu", "lu"},
		},
		Cookie: map[string]string{
			"sid": "123456",
		},
		RemoteAddr:    "127.0.0.1",
		RemotePort:    80,
		TimeLocal:     "23/Jul/2018:22:23:35 +0800",
		TimeISO8601:   "2018-07-23T22:23:35+08:00",
		Status:        200,
		BodyBytesSent: 1048,
		Request:       "GET / HTTP/1.1",
		Header: map[string][]string{
			"User-Agent": {
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.68 Safari/537.36",
			},
			"Referer": {
				"https://www.baidu.com/",
			},
		},
	}

	format := "${args} ${arg.name} ${cookie.sid} ${remoteAddr} - [${timeLocal}] \"${request}\" ${status} ${bodyBytesSent} \"${http.Referer}\" \"${http.UserAgent}\""
	t.Log(accessLog.Format(format))

	format = "Extend:${extend.File} ${extend.Geo}"
	t.Log(accessLog.Format(format))
}

func TestAccessLog_Decode(t *testing.T) {
	client := teamongo.SharedClient()
	if client == nil {
		t.Fatal("client=nil")
	}

	r, err := client.
		Database("teaweb").
		Collection("accessLogs").
		Find(context.Background(), bson.NewDocument(
			//bson.EC.String("remoteAddr", "127.0.0.1"),
			bson.EC.SubDocument("id", bson.NewDocument(bson.EC.Int64("$gt", 1535886567943382000))),
		), findopt.Skip(0), findopt.Limit(2), findopt.Sort(bson.NewDocument(

			bson.EC.Int32("_id", 1),
		)))
	if err != nil {
		t.Fatal(err)
	}

	for r.Next(context.Background()) {
		accessLog := AccessLog{}
		err = r.Decode(&accessLog)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("mongoId:", accessLog.Id)
		t.Log(accessLog)
	}

	client.Disconnect(context.Background())
}

func TestAccessLog_ParseGEO(t *testing.T) {
	accessLog := &AccessLog{
		RemoteAddr: "183.131.156.10",
	}

	//ip := net.ParseIP("183.131.156.10")
	//ip := net.ParseIP("111.197.204.174")
	accessLog.parseGeoIP()

	t.Logf("%#v", accessLog.Extend.Geo)
}
