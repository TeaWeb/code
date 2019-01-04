package locationutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

// 设置公用信息
func SetCommonInfo(action actions.ActionWrapper, serverFilename string, locationId string, subTab string) (server *teaconfigs.ServerConfig, location *teaconfigs.LocationConfig) {
	obj := action.Object()

	server, err := teaconfigs.NewServerConfigFromFile(serverFilename)
	if err != nil {
		obj.Fail(err.Error())
	}

	location = server.FindLocation(locationId)
	if location == nil {
		obj.Fail("找不到要修改的Location")
	}

	obj.Data["location"] = maps.Map{
		"id":          location.Id,
		"pattern":     location.PatternString(),
		"fastcgi":     location.Fastcgi,
		"rewrite":     location.Rewrite,
		"headers":     location.Headers,
		"cachePolicy": location.CachePolicy,
		"websocket":   location.Websocket,
	}

	obj.Data["selectedTab"] = "location"
	obj.Data["selectedSubTab"] = subTab
	obj.Data["filename"] = server.Filename
	obj.Data["proxy"] = server
	obj.Data["server"] = maps.Map{
		"filename": server.Filename,
	}

	return server, location
}
