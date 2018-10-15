package teaplugins

import (
	"github.com/TeaWeb/code/teacharts"
	"github.com/TeaWeb/code/teainterfaces"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"plugin"
	"reflect"
	"sync"
)

var plugins = []*Plugin{}
var pluginsLocker = &sync.Mutex{}

func Register(plugin *Plugin) {
	pluginsLocker.Lock()
	plugins = append(plugins, plugin)
	pluginsLocker.Unlock()
}

func Plugins() []*Plugin {
	return plugins
}

func TopBarWidgets() []*Widget {
	pluginsLocker.Lock()
	defer pluginsLocker.Unlock()

	result := []*Widget{}
	for _, p := range plugins {
		for _, widget := range p.Widgets {
			if widget.TopBar {
				result = append(result, widget)
			}
		}
	}
	return result
}

func MenuBarWidgets() []*Widget {
	pluginsLocker.Lock()
	defer pluginsLocker.Unlock()

	result := []*Widget{}
	for _, p := range plugins {
		for _, widget := range p.Widgets {
			if widget.MenuBar {
				result = append(result, widget)
			}
		}
	}
	return result
}

func HelperBarWidgets() []*Widget {
	pluginsLocker.Lock()
	defer pluginsLocker.Unlock()

	result := []*Widget{}
	for _, p := range plugins {
		for _, widget := range p.Widgets {
			if widget.HelperBar {
				result = append(result, widget)
			}
		}
	}
	return result
}

func DashboardWidgets(group WidgetGroup) []*Widget {
	pluginsLocker.Lock()
	defer pluginsLocker.Unlock()

	result := []*Widget{}
	for _, p := range plugins {
		for _, widget := range p.Widgets {
			if widget.Dashboard && widget.Group == group {
				result = append(result, widget)
			}
		}
	}
	return result
}

func load() {
	logs.Println("[plugin]load plugins")
	dir := Tea.Root + Tea.DS + "plugins"
	files.NewFile(dir).Range(func(file *files.File) {
		if file.Ext() == ".so" {
			p, err := plugin.Open(file.Path())
			if err != nil {
				logs.Println("[plugin]"+file.Name()+":", err.Error())
				return
			}

			newFunc, err := p.Lookup("New")
			if err != nil {
				logs.Println("[plugin]"+file.Name()+":", err.Error())
				return
			}

			instance := reflect.ValueOf(newFunc)
			if instance.IsNil() || !instance.IsValid() {
				logs.Println("[plugin]" + file.Name() + ": New() not a function")
				return
			}
			t := instance.Type()
			if t.Kind() != reflect.Func {
				logs.Println("[plugin]" + file.Name() + ": New() not a function")
				return
			}

			if t.NumIn() > 0 {
				logs.Println("[plugin]" + file.Name() + ": New() too many arguments")
				return
			}

			result := instance.Call([]reflect.Value{})
			if len(result) == 0 {
				logs.Println("[plugin]" + file.Name() + ": New() should return a result of 'teainterfaces.PluginInterface'")
				return
			}

			resultInstance := result[0].Interface()
			if resultInstance == nil {
				logs.Println("[plugin]" + file.Name() + ": New() should return a non-nil result")
				return
			}

			p1, ok := resultInstance.(teainterfaces.PluginInterface)
			if !ok {
				logs.Println("[plugin]" + file.Name() + ": New() should return a result of 'teainterfaces.PluginInterface'")
				return
			}

			loadInterface(p1, file.Name())

			logs.Println("[plugin]loaded", "'"+file.Name()+"'")
		}
	})
}

func loadInterface(p1 teainterfaces.PluginInterface, fileName string) {
	p1.OnLoad()
	p1.OnStart()

	p2 := NewPlugin()
	p2.IsExternal = true
	p2.Name = p1.Name()
	p2.Code = p1.Code()
	p2.Date = p1.Date()
	p2.Site = p1.Site()
	p2.Developer = p1.Developer()
	p2.Version = p1.Version()
	p2.Description = p1.Description()

	// widget
	for _, w := range p1.Widgets() {
		w1, ok := w.(teainterfaces.WidgetInterface)
		if !ok {
			logs.Println("[plugin]invalid widget in", fileName)
			continue
		}

		w2 := NewWidget()
		w2.Name = w1.Name()
		w2.URL = w1.URL()
		w2.MoreURL = w1.MoreURL()
		w2.Group = w1.Group()
		w2.TopBar = w1.TopBar()
		w2.MenuBar = w1.MenuBar()
		w2.HelperBar = w1.HelperBar()
		w2.Dashboard = w1.Dashboard()
		w2.OnForceReload(func() {
			// chart
			loadWidgetInterface(w1, w2, fileName)

			w1.OnReload()
		})
		w2.OnReload(func() {
			w1.OnReload()

			// chart
			loadWidgetInterface(w1, w2, fileName)
		})

		// chart
		loadWidgetInterface(w1, w2, fileName)

		p2.AddWidget(w2)
	}

	Register(p2)
}

func loadWidgetInterface(w1 teainterfaces.WidgetInterface, w2 *Widget, fileName string) {
	for _, c := range w1.Charts() {
		c1, ok := c.(teainterfaces.ChartInterface)
		if !ok {
			logs.Println("[plugin]invalid chart in", fileName)
			continue
		}

		c2 := teacharts.ConvertInterface(c1)
		if c2 == nil {
			logs.Println("[plugin]invalid chart in", fileName, "chart type:", c1.Type())
			continue
		}

		if len(c1.Id()) > 0 {
			c2.SetUniqueId(c1.Id())
		}

		w2.AddChart(c2)
		c1.SetId(c2.UniqueId())
	}
}
