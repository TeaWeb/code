package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"regexp"
	"strings"
)

// 添加新的服务
type AddAction actions.Action

func (this *AddAction) Run(params struct {
}) {
	this.Show()
}

// 提交保存
func (this *AddAction) RunPost(params struct {
	Description string
	ServiceType uint
	Name        string
	Listen      string
	Backend     string
	Root        string
	Must        *actions.Must
}) {
	if len(params.Description) == 0 {
		params.Description = "新代理服务"
	}

	server := teaconfigs.NewServerConfig()
	server.Http = true
	server.Description = params.Description
	server.Charset = "utf-8"
	server.Index = []string{"index.html", "index.htm", "index.php"}

	if len(params.Name) > 0 {
		for _, name := range regexp.MustCompile("\\s+").Split(params.Name, -1) {
			name = strings.TrimSpace(name)
			if len(name) > 0 {
				server.AddName(name)
			}
		}
	}

	if len(params.Listen) > 0 {
		for _, listen := range regexp.MustCompile("\\s+").Split(params.Listen, -1) {
			listen = strings.TrimSpace(listen)
			if len(listen) > 0 {
				server.AddListen(listen)
			}
		}
	}

	if params.ServiceType == 1 { // 代理服务
		for _, backend := range regexp.MustCompile("\\s+").Split(params.Backend, -1) {
			backend = strings.TrimSpace(backend)
			if len(backend) > 0 {
				backendObject := teaconfigs.NewBackendConfig()
				backendObject.Address = backend
				server.AddBackend(backendObject)
			}
		}
	} else if params.ServiceType == 2 { // 普通服务
		server.Root = params.Root
	}

	err := server.Validate()
	if err != nil {
		this.Fail("添加时有问题发生：" + err.Error())
	}

	filename := "server." + server.Id + ".proxy.conf"
	server.Filename = filename
	configPath := Tea.ConfigFile(filename)
	err = server.WriteToFile(configPath)
	if err != nil {
		this.Fail(err.Error())
	}

	proxyutils.NotifyChange()

	this.Next("/proxy/detail", map[string]interface{}{
		"serverId": server.Id,
	}, "").Success("添加成功，现在去查看详细信息")
}
