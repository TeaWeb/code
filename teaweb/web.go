package teaweb

import (
	"github.com/TeaWeb/code/teaproxy"
	_ "github.com/TeaWeb/code/teaservices"
	_ "github.com/TeaWeb/code/teaweb/actions/default/apps"
	_ "github.com/TeaWeb/code/teaweb/actions/default/dashboard"
	"github.com/TeaWeb/code/teaweb/actions/default/index"
	"github.com/TeaWeb/code/teaweb/actions/default/install"
	_ "github.com/TeaWeb/code/teaweb/actions/default/log"
	_ "github.com/TeaWeb/code/teaweb/actions/default/login"
	"github.com/TeaWeb/code/teaweb/actions/default/logout"
	_ "github.com/TeaWeb/code/teaweb/actions/default/monitor"
	_ "github.com/TeaWeb/code/teaweb/actions/default/plugins"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/backend"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/fastcgi"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/headers"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/locations"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/rewrite"
	_ "github.com/TeaWeb/code/teaweb/actions/default/proxy/ssl"
	_ "github.com/TeaWeb/code/teaweb/actions/default/settings"
	_ "github.com/TeaWeb/code/teaweb/actions/default/settings/login"
	_ "github.com/TeaWeb/code/teaweb/actions/default/settings/mongo"
	_ "github.com/TeaWeb/code/teaweb/actions/default/settings/server"
	_ "github.com/TeaWeb/code/teaweb/actions/default/stat"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/TeaWeb/code/teaweb/utils"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/sessions"
	"time"
)

func Start() {
	// 日志
	writer := new(utils.LogWriter)
	writer.Init()
	logs.SetWriter(writer)

	// 启动代理
	go func() {
		time.Sleep(1 * time.Second)

		// 启动代理
		teaproxy.Start()
	}()

	// 启动管理界面
	TeaGo.NewServer().
		AccessLog(false).
		Get("/", new(index.IndexAction)).
		Get("/logout", new(logout.IndexAction)).

		Helper(new(helpers.UserMustAuth)).
		GetPost("/install/mongo", new(install.MongoAction)).
		EndAll().

		Session(sessions.NewFileSessionManager(
			86400,
			"gSeDQJJ67tAVdnguDAQdGmnDVrjFd2I9",
		)).

		Start()
}
