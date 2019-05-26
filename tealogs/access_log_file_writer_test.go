package tealogs

import (
	"github.com/TeaWeb/code/teaconfigs"
	"testing"
)

func TestAccessLogFileConfig_Init(t *testing.T) {
	config := &teaconfigs.AccessLogFileConfig{
		Path: "logs/access.log",
	}
	writer := AccessLogFileWriter{
		config: config,
	}
	writer.Init()
	writer.Write(&AccessLog{
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
		}})

	writer.Close()
}
