package teaplugins

import "sync"

var plugins = []*Plugin{}
var pluginsLocker = &sync.Mutex{}

func Register(plugin *Plugin) {
	pluginsLocker.Lock()
	plugins = append(plugins, plugin)
	pluginsLocker.Unlock()
}

func TopBarWidgets() []*Widget {
	pluginsLocker.Lock()
	defer pluginsLocker.Unlock()

	result := []*Widget{}
	for _, plugin := range plugins {
		for _, widget := range plugin.Widgets {
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
	for _, plugin := range plugins {
		for _, widget := range plugin.Widgets {
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
	for _, plugin := range plugins {
		for _, widget := range plugin.Widgets {
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
	for _, plugin := range plugins {
		for _, widget := range plugin.Widgets {
			if widget.Dashboard && widget.Group == group {
				result = append(result, widget)
			}
		}
	}
	return result
}
