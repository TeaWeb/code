package tealogs

import (
	"github.com/TeaWeb/uaparser"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/oschwald/geoip2-golang"
	"reflect"
)

var userAgentParser *uaparser.Parser
var geoDB *geoip2.Reader
var accessLogVars = map[string]string{}

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		accessLogger = NewAccessLogger()
	})
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		logs.Println("[proxy]start user-agent parser")
		var err error
		userAgentParser, err = uaparser.NewParser(Tea.Root + Tea.DS + "resources" + Tea.DS + "regexes.yaml")
		if err != nil {
			logs.Error(err)
		}

		logs.Println("[proxy]start geo db")
		geoDB, err = geoip2.Open(Tea.Root + "/resources/GeoLite2-City/GeoLite2-City.mmdb")
		if err != nil {
			logs.Error(err)
		}

		// 变量
		reflectType := reflect.TypeOf(AccessLog{})
		countField := reflectType.NumField()
		for i := 0; i < countField; i ++ {
			field := reflectType.Field(i)
			value := field.Tag.Get("var")
			if len(value) > 0 {
				accessLogVars[value] = field.Name
			}
		}
	})
}
