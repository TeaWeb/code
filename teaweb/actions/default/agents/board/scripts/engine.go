package scripts

import (
	"encoding/json"
	"errors"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/caches"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/robertkrimen/otto"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var engineCache = caches.NewFactory()

// 脚本引擎
type Engine struct {
	vm           *otto.Otto
	chartOptions []maps.Map
	widgetCodes  map[string]maps.Map // "code" => { name, ..., definition:FUNCTION CODE }
	context      *Context
	output       []string
	mongoEnabled bool
}

// 获取新引擎
func NewEngine() *Engine {
	engine := &Engine{
		chartOptions: []maps.Map{},
		widgetCodes:  map[string]maps.Map{},
	}
	engine.init()
	return engine
}

// 设置MongoDB是否可用
func (this *Engine) SetMongo(b bool) {
	this.mongoEnabled = b
}

// 设置上下文信息
func (this *Engine) SetContext(context *Context) {
	this.context = context

	if context.Agent != nil {
		options := map[string]interface{}{
			"isOn":    context.Agent.On,
			"id":      context.Agent.Id,
			"isLocal": context.Agent.IsLocal(),
			"name":    context.Agent.Name,
			"host":    context.Agent.Host,
			"apps": lists.Map(context.Agent.Apps, func(k int, v interface{}) interface{} {
				app := v.(*agents.AppConfig)
				return maps.Map{
					"id":   app.Id,
					"isOn": app.On,
					"name": app.Name,
					"tasks": lists.Map(app.Tasks, func(k int, v interface{}) interface{} {
						task := v.(*agents.TaskConfig)
						return maps.Map{
							"id":           task.Id,
							"isOn":         task.On,
							"name":         task.Name,
							"isBooting":    task.IsBooting,
							"isManual":     task.IsManual,
							"isScheduling": len(task.Schedule) > 0,
						}
					}),
				}
			}),
		}

		_, err := this.vm.Run(`context.agent = new agents.Agent(` + stringutil.JSONEncode(options) + `);`)
		if err != nil {
			logs.Error(err)
		}
	}

	if context.App != nil {
		options := map[string]interface{}{
			"id":   context.App.Id,
			"isOn": context.App.On,
			"name": context.App.Name,
			"tasks": lists.Map(context.App.Tasks, func(k int, v interface{}) interface{} {
				task := v.(*agents.TaskConfig)
				return maps.Map{
					"id":           task.Id,
					"isOn":         task.On,
					"name":         task.Name,
					"isBooting":    task.IsBooting,
					"isManual":     task.IsManual,
					"isScheduling": len(task.Schedule) > 0,
				}
			}),
		}

		_, err := this.vm.Run(`context.app = new agents.App(` + stringutil.JSONEncode(options) + `);`)
		if err != nil {
			logs.Error(err)
		}
	}

	if context.Item != nil {
		options := map[string]interface{}{
			"id":       context.Item.Id,
			"isOn":     context.Item.On,
			"name":     context.Item.Name,
			"interval": context.Item.IntervalDuration().Seconds(),
		}
		_, err := this.vm.Run(`context.item = new agents.Item(` + stringutil.JSONEncode(options) + `);`)
		if err != nil {
			logs.Error(err)
		}
	}

	// 可供使用的特性
	features := []string{}
	if this.mongoEnabled {
		features = append(features, "mongo")
	}
	features = append(features, runtime.GOOS)
	features = append(features, runtime.GOARCH)
	this.vm.Run(`context.features=` + stringutil.JSONEncode(features) + `;`)
}

// 初始化
func (this *Engine) init() {
	this.vm = otto.New()

	this.vm.Set("callConsoleLog", this.callConsoleLog)
	this.vm.Run("console.log = callConsoleLog;")

	this.loadLib("libs/array.js")
	this.loadLib("libs/times.js")
	this.loadLib("libs/caches.js")
	this.loadLib("libs/mutex.js")
	this.loadLib("libs/agent.values.js")
	this.loadLib("libs/colors.js")
	this.loadLib("libs/widgets.js")
	this.loadLib("libs/charts.js")
	this.loadLib("libs/charts.gauge.js")
	this.loadLib("libs/charts.html.js")
	this.loadLib("libs/charts.line.js")
	this.loadLib("libs/charts.pie.js")
	this.loadLib("libs/charts.url.js")
	this.loadLib("libs/charts.progress.js")
	this.loadLib("libs/charts.stackbar.js")
	this.loadLib("libs/charts.clock.js")
	this.loadLib("libs/charts.table.js")
	this.loadLib("libs/context.js")
	this.loadLib("libs/agent.js")

	this.vm.Set("callSetCache", this.callSetCache)
	this.vm.Set("callGetCache", this.callGetCache)
	this.vm.Set("callChartRender", this.callRenderChart)
	this.vm.Set("callExecuteQuery", this.callExecuteQuery)
}

// 运行Widget代码
func (this *Engine) RunCode(code string) error {
	_, err := this.vm.Run(`(function () {` + code + `
	
	widget.callRun();
})();`)
	return err
}

// 获取Widget中的图表对象
func (this *Engine) Charts() []maps.Map {
	return this.chartOptions
}

// 获取控制台输出
func (this *Engine) Output() []string {
	if this.output == nil {
		return []string{}
	}
	return this.output
}

func (this *Engine) callConsoleLog(call otto.FunctionCall) otto.Value {
	values := []string{}
	for _, v := range call.ArgumentList {
		i, err := v.Export()
		if err != nil {
			values = append(values, v.String())
		} else {
			values = append(values, stringutil.JSONEncodePretty(i))
		}
	}
	s := strings.Join(values, ", ")
	//logs.Println("[console]", s)

	this.output = append(this.output, s)

	return otto.UndefinedValue()
}

func (this *Engine) callRenderChart(call otto.FunctionCall) otto.Value {
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

func (this *Engine) callSetCache(call otto.FunctionCall) otto.Value {
	key, err := call.Argument(0).ToString()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	key = stringutil.Md5(key)

	value, err := call.Argument(1).Export()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}

	life, err := call.Argument(2).Export()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	lifeSeconds := types.Int64(life)
	engineCache.Set(key, value, time.Duration(lifeSeconds)*time.Second)

	return otto.UndefinedValue()
}

// 获取缓存
func (this *Engine) callGetCache(call otto.FunctionCall) otto.Value {
	key, err := call.Argument(0).ToString()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	key = stringutil.Md5(key)

	value, found := engineCache.Get(key)
	if !found {
		return otto.UndefinedValue()
	}
	v, err := this.vm.ToValue(value)
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	return v
}

// 加载JS库文件
func (this *Engine) loadLib(file string) {
	path := Tea.Root + Tea.DS + "web" + Tea.DS + file
	cacheKey := "libfile://" + path
	code, found := engineCache.Get(cacheKey)
	if !found {
		var err error = nil
		code, err = files.NewFile(path).ReadAllString()
		if err != nil {
			logs.Error(err)
			return
		}
		engineCache.Set(cacheKey, code)
	}

	_, err := this.vm.Run(code)
	if err != nil {
		logs.Error(err)
		return
	}
}

func (this *Engine) callExecuteQuery(call otto.FunctionCall) otto.Value {
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

	query := teamongo.NewAgentValueQuery()
	if this.context != nil {
		if this.context.Agent != nil {
			query.Agent(this.context.Agent.Id)
		}
		if this.context.App != nil {
			query.App(this.context.App.Id)
		}
		if this.context.Item != nil {
			query.Item(this.context.Item.Id)
		}
	}

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

	jsValue, err := this.toValue(v)
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	return jsValue
}

func (this *Engine) toValue(data interface{}) (v otto.Value, err error) {
	if data == nil {
		return this.vm.ToValue(data)
	}

	// *Value
	if _, ok := data.(*agents.Value); ok {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return this.vm.ToValue(data)
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			logs.Error(err)
			return this.vm.ToValue(data)
		}
		return this.vm.ToValue(m)
	}

	// []*Value
	if _, ok := data.([]*agents.Value); ok {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return this.vm.ToValue(data)
		}
		m := []map[string]interface{}{}
		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			logs.Error(err)
			return this.vm.ToValue(data)
		}
		return this.vm.ToValue(m)
	}

	return this.vm.ToValue(data)
}

func (this *Engine) throw(err error) {
	if err != nil {
		value, _ := this.vm.Call("new Error", nil, err.Error())
		panic(value)
	}
}
