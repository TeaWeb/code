package teaplugins

import (
	"bufio"
	"bytes"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"net/http/httputil"
	"sync"
)

var plugins = []*Plugin{}
var pluginsLocker = &sync.Mutex{}

var requestFilters = []func(req []byte) (result []byte, willContinue bool){}
var HasRequestFilters = false
var responseFilters = []func(resp []byte) (result []byte, willContinue bool){}
var HasResponseFilters = false

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

func FilterRequest(request *http.Request) (resultReq *http.Request, willContinue bool) {
	if !HasRequestFilters {
		return request, true
	}

	data, err := httputil.DumpRequest(request, true)
	if err != nil {
		logs.Error(err)
		return request, true
	}

	defer func() {
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(data)))
		if err != nil {
			logs.Error(err)
			return
		}

		resultReq = req
	}()

	for _, f := range requestFilters {
		result, willContinue := f(data)

		data = result

		if !willContinue {
			return resultReq, false
		}
	}

	return resultReq, true
}

func FilterResponse(response *http.Response) (resultResp *http.Response) {
	if !HasResponseFilters {
		return response
	}

	data, err := httputil.DumpResponse(response, true)
	if err != nil {
		logs.Error(err)
		return response
	}

	defer func() {
		resp, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(data)), nil)
		if err != nil {
			logs.Error(err)
			return
		}

		resultResp = resp
	}()

	for _, f := range responseFilters {
		result, willContinue := f(data)

		data = result

		if !willContinue {
			return resultResp
		}
	}

	return resultResp
}

func load() {
	logs.Println("[plugin]load plugins")
	dir := Tea.Root + Tea.DS + "plugins"
	files.NewFile(dir).Range(func(file *files.File) {
		if file.Ext() != ".tea" {
			return
		}

		logs.Println("[plugin][loader]load plugin '" + file.Name() + "'")
		go NewLoader(file.Path()).Load()
	})
}
