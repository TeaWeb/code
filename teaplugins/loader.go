package teaplugins

import (
	"encoding/binary"
	"errors"
	"github.com/TeaWeb/code/teaapps"
	"github.com/TeaWeb/code/teacharts"
	"github.com/TeaWeb/plugin/apps"
	"github.com/TeaWeb/plugin/messages"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/processes"
	"os"
	"reflect"
	"time"
)

type Loader struct {
	path   string
	plugin *Plugin

	methods   map[string]reflect.Method
	thisValue reflect.Value

	writer *os.File

	debug bool
}

func NewLoader(path string) *Loader {
	loader := &Loader{
		path:    path,
		methods: map[string]reflect.Method{},
	}

	// 当前methods
	t := reflect.TypeOf(loader)
	for i := 0; i < t.NumMethod(); i ++ {
		method := t.Method(i)
		loader.methods[method.Name] = method
	}

	loader.thisValue = reflect.ValueOf(loader)

	return loader
}

func (this *Loader) Debug() {
	this.debug = true
}

func (this *Loader) Load() error {
	reader, w /** 子进程写入器 **/, err := os.Pipe()
	if err != nil {
		return err
	}

	r2 /** 子进程读取器 **/, writer, err := os.Pipe()
	if err != nil {
		return err
	}

	this.writer = writer

	p := processes.NewProcess(this.path)
	p.AppendFile(r2, w)

	go func() {
		buf := make([]byte, 1024)
		msgData := []byte{}
		for {
			if this.debug {
				logs.Println("[plugin][loader]try to read buf")
			}

			n, err := reader.Read(buf)

			if n > 0 {
				msgData = append(msgData, buf[:n] ...)

				if this.debug {
					logs.Println("[plugin][loader]len:", len(msgData), ",", "read msg data:", string(msgData))
				}

				msgLen := uint32(len(msgData))
				h := uint32(24) // header length

				if msgLen > h { // 数据组成方式： | actionLen[8] | dataLen[8] | action | data[len-8]
					id := binary.BigEndian.Uint32(msgData[:8])
					l1 := binary.BigEndian.Uint32(msgData[8:16])
					l2 := binary.BigEndian.Uint32(msgData[16:24])

					if msgLen >= h+l1+l2 { // 数据已经完整了
						action := string(msgData[h : h+l1])
						valueData := msgData[h+l1 : h+l1+l2]

						msgData = msgData[h+l1+l2:]

						ptr, err := messages.Unmarshal(action, valueData)
						if err != nil {
							logs.Println("[plugin][loader]unmarshal message error:", err.Error())
							continue
						}

						err = this.CallAction(ptr, id)
						if err != nil {
							logs.Println("[plugin][loader]call action error:", err.Error())
							continue
						}
					}
				}
			}

			if err != nil {
				logs.Println(err.Error())
				break
			}
		}
	}()

	err = p.Start()
	if err != nil {
		return err
	}

	err = p.Wait()
	if err != nil {
		reader.Close()

		// 重新加载
		time.Sleep(1 * time.Second)
		return this.Load()
	}

	return nil
}

func (this *Loader) CallAction(ptr interface{}, messageId uint32) error {
	action, ok := ptr.(messages.ActionInterface)
	if !ok {
		return errors.New("ptr should be an action")
	}
	action.SetMessageId(messageId)

	method, found := this.methods["Action"+action.Name()]
	if !found {
		return errors.New("[plugin]handler for '" + action.Name() + "' not found")

	}
	method.Func.Call([]reflect.Value{this.thisValue, reflect.ValueOf(action)})
	return nil
}

func (this *Loader) ActionRegisterPlugin(action *messages.RegisterPluginAction) {
	if this.plugin != nil {
		logs.Println("[plugin][loader]load only one plugin from one file")
		return
	}

	// 添加到插件中
	if action.Plugin != nil {
		p := action.Plugin
		p2 := NewPlugin()
		p2.IsExternal = true
		p2.Name = p.Name
		p2.Description = p.Description
		p2.Code = p.Code
		p2.Site = p.Site
		p2.Developer = p.Developer
		p2.Date = p.Date
		p2.Version = p.Version

		// request filter
		if p.HasRequestFilter {
			requestFilters = append(requestFilters, func(data []byte) (result []byte, willContinue bool) {
				action := &messages.FilterRequestAction{
					Data: data,
				}
				this.Write(action)

				respAction := messages.ActionQueue.Wait(action)
				r, ok := respAction.(*messages.FilterRequestAction)
				if ok {
					return r.Data, r.Continue
				} else {
					return action.Data, true
				}
			})
			hasRequestFilters = true
		}

		// widget
		for _, w := range p.Widgets {
			w2 := NewWidget()
			w2.Id = w.Id
			w2.Name = w.Name
			w2.Icon = w.Icon
			w2.Title = w.Title
			w2.URL = w.URL
			w2.MoreURL = w.MoreURL
			w2.TopBar = w.TopBar
			w2.MenuBar = w.MenuBar
			w2.HelperBar = w.HelperBar
			w2.Dashboard = w.Dashboard
			w2.Group = w.Group
			w2.OnForceReload(func() {
				action := new(messages.ReloadWidgetAction)
				action.WidgetId = w2.Id
				this.Write(action)
			})
			w2.OnReload(func() {
				action := new(messages.ReloadWidgetAction)
				action.WidgetId = w2.Id
				this.Write(action)
			})

			// chart
			for _, c := range w.Charts {
				c2 := teacharts.ConvertInterface(c)
				if c2 == nil {
					continue
				}
				w2.AddChart(c2)
			}

			p2.AddWidget(w2)
		}

		// apps
		for _, a := range p.Apps {
			a2 := teaapps.NewApp()
			a2.LoadFromInterface(a)
			p2.AddApp(a2)
		}

		Register(p2)

		this.plugin = p2
	}

	// 发送启动信息
	this.Write(&messages.StartAction{})
}

func (this *Loader) ActionReloadChart(action *messages.ReloadChartAction) {
	chart := teacharts.ConvertInterface(action.Chart)
	if chart == nil {
		return
	}

	// 查找
	for _, w := range this.plugin.Widgets {
		for index, c := range w.Charts {
			if c.ChartId() == chart.ChartId() {
				w.Charts[index] = chart
				break
			}
		}
	}
}

func (this *Loader) ActionFilterRequest(action *messages.FilterRequestAction) {
	messages.ActionQueue.Notify(action)
}

func (this *Loader) ActionReloadApps(action *messages.ReloadAppsAction) {
	this.plugin.ResetApps()

	for _, a := range action.Apps {
		a2 := teaapps.NewApp()
		a2.LoadFromInterface(a)
		a2.OnReload(func() {
			this.Write(&messages.ReloadAppAction{
				App: &apps.App{
					Id: a.Id,
				},
			})
		})
		this.plugin.AddApp(a2)
	}
}

func (this *Loader) ActionReloadApp(action *messages.ReloadAppAction) {
	a := action.App
	if a != nil {
		app := this.plugin.AppWithId(a.Id)
		if app != nil {
			app.LoadFromInterface(a)
		}
	}
}

func (this *Loader) Write(action messages.ActionInterface) error {
	msg := messages.NewActionMessage(action)
	msg.Id = action.MessageId()
	data, err := msg.Marshal()
	if err != nil {
		return err
	}
	action.SetMessageId(msg.Id)
	_, err = this.writer.Write(data)
	return err
}
