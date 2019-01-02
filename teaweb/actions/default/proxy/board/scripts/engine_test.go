package scripts

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/maps"
	"testing"
)

func TestEngine_RunJSON(t *testing.T) {
	engine := NewEngine()
	engine.SetContext(&Context{
		Server: &teaconfigs.ServerConfig{
			Id: "123",
		},
	})
	err := engine.RunConfig(Tea.ConfigFile("board.iONhcceoPPB73vYx.conf"), maps.Map{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(engine.Charts())
}

func TestEngine_Run(t *testing.T) {
	engine := NewEngine()
	err := engine.RunCode(`var widget = new widgets.Widget({
	"name": "测试Widget",
	"code": "test_stat@tea",
	"author": "我是测试的",
	"version": "0.0.1"
});

widget.run = function () {
	var chart = new charts.HTMLChart();
	chart.html = "测试HTML Chart";
	chart.render();
};`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(engine.Charts())
}

func TestEngine_Cache(t *testing.T) {
	engine := NewEngine()
	err := engine.RunCode(`
var widget = new widgets.Widget({});
widget.run = function () {
	caches.set("a", "b");
	console.log(caches.get("a"));
};
`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEngine_Log(t *testing.T) {
	engine := NewEngine()
	engine.SetContext(&Context{
		Server: &teaconfigs.ServerConfig{
			Id: "123",
		},
	})
	err := engine.RunCode(`
var widget = new widgets.Widget({});
widget.run = function () {
	var query = new logs.Query();
	query.attr("status", [200]);
	query.count();
};
`)
	if err != nil {
		t.Fatal(err)
	}
}
