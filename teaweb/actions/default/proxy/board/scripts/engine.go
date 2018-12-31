package scripts

import (
	"encoding/json"
	"errors"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/robertkrimen/otto"
	"reflect"
	"time"
)

// 脚本引擎
type Engine struct {
	vm           *otto.Otto
	chartOptions []maps.Map
	widgetCodes  map[string]maps.Map // "code" => { name, ..., definition:FUNCTION CODE }
}

func NewEngine() *Engine {
	engine := &Engine{
		chartOptions: []maps.Map{},
		widgetCodes:  map[string]maps.Map{},
	}
	engine.init()
	return engine
}

func (this *Engine) SetContext(context *Context) {
	if context.Server != nil {
		runningServer, _ := teaproxy.FindServer(context.Server.Id)

		options := map[string]interface{}{
			"isOn":        context.Server.On,
			"id":          context.Server.Id,
			"name":        context.Server.Name,
			"filename":    context.Server.Filename,
			"description": context.Server.Description,
			"listen":      context.Server.Listen,
			"backends": lists.Map(context.Server.Backends, func(k int, v interface{}) interface{} {
				backend := v.(*teaconfigs.ServerBackendConfig)

				if runningServer != nil {
					runningBackend := runningServer.FindBackend(backend.Id)
					if runningBackend != nil {
						backend.IsDown = runningBackend.IsDown
					}
				}

				return map[string]interface{}{
					"on":       backend.On,
					"weight":   backend.Weight,
					"id":       backend.Id,
					"isDown":   backend.IsDown,
					"isBackup": backend.IsBackup,
					"name":     backend.Name,
					"address":  backend.Address,
				}
			}),
			"locations": lists.Map(context.Server.Locations, func(k int, v interface{}) interface{} {
				location := v.(*teaconfigs.LocationConfig)
				location.Validate()
				return map[string]interface{}{
					"id":          location.Id,
					"on":          location.On,
					"pattern":     location.PatternString(),
					"cachePolicy": location.CachePolicy,
					"fastcgi": lists.Map(location.Fastcgi, func(k int, v interface{}) interface{} {
						fastcgi := v.(*teaconfigs.FastcgiConfig)
						return map[string]interface{}{
							"id":   fastcgi.Id,
							"on":   fastcgi.On,
							"pass": fastcgi.Pass,
						}
					}),
					"rewrite": lists.Map(location.Rewrite, func(k int, v interface{}) interface{} {
						rewrite := v.(*teaconfigs.RewriteRule)
						return map[string]interface{}{
							"id":      rewrite.Id,
							"on":      rewrite.On,
							"pattern": rewrite.Pattern,
							"replace": rewrite.Replace,
						}
					}),
					"root":    location.Root,
					"index":   location.Index,
					"headers": location.Headers,
				}
			}),
		}
		this.vm.Run(`context.server = new http.Server(` + this.jsonEncode(options) + `);`)
	}

	// 可供使用的特性
	features := []string{}
	if teamongo.Test() == nil {
		features = append(features, "mongo")
	}
	this.vm.Run(`context.features=` + this.jsonEncode(features) + `;`)
}

func (this *Engine) init() {
	this.vm = otto.New()
	this.loadLib("libs/array.js")
	this.loadLib("libs/times.js")
	this.loadLib("libs/mutex.js")
	this.loadLib("libs/logs.js")
	this.loadLib("libs/http.js")
	this.loadLib("libs/colors.js")
	this.loadLib("libs/widgets.js")
	this.loadLib("libs/charts.js")
	this.loadLib("libs/charts.gauge.js")
	this.loadLib("libs/charts.html.js")
	this.loadLib("libs/charts.line.js")
	this.loadLib("libs/charts.pie.js")
	this.loadLib("libs/charts.progress.js")
	this.loadLib("libs/context.js")

	this.loadWidgets()

	this.vm.Set("callChartRender", this.renderChart)
	this.vm.Set("callExecuteQuery", this.executeQuery)
}

func (this *Engine) RunConfig(configFile string, options maps.Map) error {
	reader, err := files.NewReader(configFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	m := maps.Map{}
	err = reader.ReadYAML(&m)
	if err != nil {
		return err
	}

	widgets := m.Get("widgets")
	if widgets == nil {
		return nil
	}
	if reflect.TypeOf(widgets).Kind() != reflect.Slice {
		return errors.New("'widgets' should be array")
	}

	arr, ok := widgets.([]interface{})
	if !ok {
		return errors.New("'widgets' format not valid")
	}

	for _, item := range arr {
		m := maps.NewMap(item)
		code := m.GetString("code")
		if len(code) == 0 {
			return errors.New("'code' should not be empty")
		}

		widget, found := this.widgetCodes[code]
		if !found {
			return errors.New("widget with code '" + code + "' not found")
		}
		err = this.RunCode(widget.GetString("definition"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Engine) RunCode(code string) error {
	_, err := this.vm.Run(`(function () {` + code + `
	
	widget.callRun();
})();`)
	return err
}

func (this *Engine) Charts() []maps.Map {
	return this.chartOptions
}

func (this *Engine) renderChart(call otto.FunctionCall) otto.Value {
	obj := call.Argument(0)
	v, err := obj.Export()
	if err != nil {
		logs.Error(err)
		return otto.UndefinedValue()
	}
	m := maps.NewMap(v)

	options, err := obj.Object().Get("options")
	if err != nil {
		logs.Error(err)
	} else {
		v, err := options.Export()
		if err != nil {
			logs.Error(err)
		} else {
			m["options"] = maps.NewMap(v)
		}
	}

	this.chartOptions = append(this.chartOptions, m)
	return otto.UndefinedValue()
}

func (this *Engine) executeQuery(call otto.FunctionCall) otto.Value {
	arg, err := call.Argument(0).Export()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	m := maps.NewMap(arg)

	action := m.GetString("action")
	if len(action) == 0 {
		this.throw(errors.New("'action' should not be empty"))
		return otto.UndefinedValue()
	}

	query := tealogs.NewQuery()

	// for
	forField := m.GetString("for")
	query.For(forField)

	// group
	group := m.Get("group")
	if group != nil {
		groupKind := reflect.TypeOf(group).Kind()
		if groupKind == reflect.String {
			query.Group([]string{group.(string)})
		} else if groupKind == reflect.Slice {
			groupSlice, ok := group.([]interface{})
			if ok {
				groupFields := []string{}
				for _, v := range groupSlice {
					groupFields = append(groupFields, types.String(v))
				}
				query.Group(groupFields)
			}
		}
	}

	// duration
	duration := m.GetString("duration")
	query.Duration(duration)

	// cond
	cond := m.Get("cond")
	if cond != nil && reflect.TypeOf(cond).Kind() == reflect.Map {
		m, ok := cond.(map[string]interface{})
		if ok {
			for field, ops := range m {
				opsMap, ok := ops.(map[string]interface{})
				if ok {
					for op, v := range opsMap {
						query.Op(op, field, v)
					}
				}
			}
		}
	}

	// timeFrom
	timeFrom := m.GetInt64("timeFrom")
	if timeFrom > 0 {
		query.From(time.Unix(timeFrom, 0))
	}

	// timeTo
	timeTo := m.GetInt64("timeTo")
	if timeTo > 0 {
		query.To(time.Unix(timeTo, 0))
	}

	// offset & size
	query.Offset(m.GetInt64("offset"))
	query.Limit(m.GetInt64("size"))

	// sort
	sorts := m.Get("sorts")
	if sorts != nil {
		sortsMap, ok := sorts.([]map[string]interface{})
		if ok {
			for _, m := range sortsMap {
				for k, v := range m {
					vInt := types.Int(v)
					if vInt < 0 {
						query.Desc(k)
					} else {
						query.Asc(k)
					}
				}
			}
		}
	}

	// 开始执行
	query.Action(action)
	v, err := query.Execute()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}

	jsValue, err := this.vm.ToValue(v)
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	return jsValue
}

// 加载widgets
func (this *Engine) loadWidgets() {
	widgetFiles := files.NewFile(Tea.Root + Tea.DS + "libs" + Tea.DS + "widgets").Glob("*.js")
	for _, file := range widgetFiles {
		s, err := file.ReadAllString()
		if err != nil {
			logs.Error(err)
			continue
		}

		widgetValue, err := this.vm.Run(`(function () {` + s + `
	return widget;
})();`)
		if err != nil {
			logs.Error(errors.New("[" + file.Name() + "]" + err.Error()))
			continue
		}
		w, err := widgetValue.Export()
		if err != nil {
			logs.Error(errors.New("[" + file.Name() + "]" + err.Error()))
			continue
		}
		m := maps.NewMap(w)
		code := m.GetString("code")
		if len(code) == 0 {
			logs.Error(errors.New("[" + file.Name() + "]'code' should not be empty"))
			continue
		}
		m["definition"] = s
		this.widgetCodes[code] = m
	}
}

// 加载JS库文件
func (this *Engine) loadLib(file string) {
	path := Tea.Root + Tea.DS + file
	s, err := files.NewFile(path).ReadAllString()
	if err != nil {
		logs.Error(err)
		return
	}
	_, err = this.vm.Run(s)
	if err != nil {
		logs.Error(err)
		return
	}
}

func (this *Engine) jsonEncode(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		logs.Error(err)
		return "null"
	}

	return string(data)
}

func (this *Engine) throw(err error) {
	if err != nil {
		value, _ := this.vm.Call("new Error", nil, err.Error())
		panic(value)
	}
}
