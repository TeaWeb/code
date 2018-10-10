package teastats

import (
	"testing"
	"github.com/TeaWeb/code/tealogs"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/iwind/TeaGo/assert"
	"time"
)

func TestProcess(t *testing.T) {
	log := &tealogs.AccessLog{
		ServerId:    "123456",
		RequestTime: 0.023,
		RemoteAddr:  "183.131.156.10",
		UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
		RequestURI:  "/",

		Scheme: "http",
		Host:   "localhost",

		SentHeader: map[string][]string{
			"Content-Type": {"text/html; charset=utf-8"},
		},
	}
	log.Parse()
	new(Processor).Process(log)
}

func TestProcessorCost(t *testing.T) {
	log := &tealogs.AccessLog{
		Id: objectid.ObjectID{0x5b, 0xbc, 0x78, 0xeb, 0xb6, 0x93, 0xb1, 0x20, 0xe2, 0xef, 0xc0, 0x9b}, ServerId: "lb001",
		BackendId:       "",
		LocationId:      "",
		FastcgiId:       "123456",
		RewriteId:       "",
		TeaVersion:      "0.0.1",
		RemoteAddr:      "127.0.0.1",
		RemotePort:      60212,
		RemoteUser:      "",
		RequestURI:      "/index.php?__ACTION__=/test",
		RequestPath:     "/index.php",
		RequestLength:   0,
		RequestTime:     0.004772714,
		RequestMethod:   "GET",
		RequestFilename: "",
		Scheme:          "http",
		Proto:           "HTTP/1.1",
		BytesSent:       4, BodyBytesSent: 4, Status: 200, StatusMessage: "200 OK", SentHeader: map[string][]string{"Content-Type": []string{"text/html; charset=UTF-8"}, "X-Powered-By": []string{"PHP/7.0.9"}}, TimeISO8601: "2018-10-09T17:46:18.945+08:00", TimeLocal: "9/Oct/2018:17:46:18 +0800", Msec: 1.5390783789453979e+09, Timestamp: 1539078378, Host: "127.0.0.1:8880", Referer: "", UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", Request: "GET /index.php?__ACTION__=/test HTTP/1.1", ContentType: "", Cookie: map[string]string{"sid": "T8J27ireA0bSZAvpf7AweQ54HBgGQuI3"}, Arg: map[string][]string(nil), Args: "__ACTION__=/test", QueryString: "__ACTION__=/test", Header: map[string][]string{"Upgrade-Insecure-Requests": []string{"1"}, "Accept-Encoding": []string{"gzip, deflate, br"}, "Connection": []string{"keep-alive"}, "Cache-Control": []string{"max-age=0"}, "Cookie": []string{"sid=T8J27ireA0bSZAvpf7AweQ54HBgGQuI3"}, "User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"}, "Accept": []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"}, "Accept-Language": []string{"zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6,de;q=0.5,ja;q=0.4"}}, ServerName: "shop.balefm.cn",
		ServerPort: 8880, ServerProtocol: "HTTP/1.1",
		BackendAddress: "",
		FastcgiAddress: "127.0.0.1:9000",
		Extend: struct {
			File   tealogs.AccessLogFile   "bson:\"file\" json:\"file\"";
			Client tealogs.AccessLogClient "bson:\"client\" json:\"client\"";
			Geo    tealogs.AccessLogGeo    "bson:\"geo\" json:\"geo\""
		}{
			File: tealogs.AccessLogFile{MimeType: "", Extension: "php", Charset: ""}, Client: tealogs.AccessLogClient{OS: tealogs.AccessLogClientOS{Family: "Mac OS X", Major: "10", Minor: "14", Patch: "0", PatchMinor: ""}, Device: tealogs.AccessLogClientDevice{Family: "Other", Brand: "", Model: ""}, Browser: tealogs.AccessLogClientBrowser{Family: "Chrome", Major: "69", Minor: "0", Patch: "3497"}}, Geo: tealogs.AccessLogGeo{Region: "", State: "", City: "", Location: tealogs.AccessLogGeoLocation{Latitude: 0, Longitude: 0, TimeZone: "", AccuracyRadius: 0x0, MetroCode: 0x0,}}},
	}

	t.Log(log)

	stat := new(DailyPVStat)
	findCollection("stats.pv.daily", stat.Init) //find collection

	a := assert.NewAssertion(t)
	for i := 0; i < 1000; i ++ {
		stat.Process(log)
	}
	a.Cost()

	time.Sleep(2 * time.Second)
}
