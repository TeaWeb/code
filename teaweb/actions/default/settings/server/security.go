package server

import (
	"github.com/TeaWeb/code/teaweb/actions/default/settings"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
)

type SecurityAction actions.Action

// 安全设置
func (this *SecurityAction) Run(params struct{}) {
	admin := configs.SharedAdminConfig()

	if admin.Security == nil {
		admin.Security = configs.NewAdminSecurity()
	}

	this.Data["security"] = admin.Security
	this.Data["allowAll"] = lists.ContainsString(admin.Security.Allow, "all")
	this.Data["userIP"] = this.RequestRemoteIP()

	this.Show()
}

func (this *SecurityAction) RunPost(params struct {
	AllowIPs        []string
	DenyIPs         []string
	AllowAll        bool
	DirAutoComplete bool
	LoginURL        string
}) {
	admin := configs.SharedAdminConfig()
	if admin.Security == nil {
		admin.Security = configs.NewAdminSecurity()
	}

	if params.AllowAll {
		admin.Security.Allow = []string{"all"}
	} else {
		ips := []string{}
		for _, ip := range params.AllowIPs {
			if len(ip) > 0 {
				ips = append(ips, ip)
			}
		}

		if len(ips) == 0 {
			this.Fail("至少要有一个允许访问的IP")
		}

		admin.Security.Allow = ips
	}

	{
		ips := []string{}
		for _, ip := range params.DenyIPs {
			if len(ip) > 0 {
				ips = append(ips, ip)
			}
		}
		admin.Security.Deny = ips
	}

	admin.Security.DirAutoComplete = params.DirAutoComplete

	isServerChanged := false
	if admin.Security.LoginURL != params.LoginURL {
		isServerChanged = true
	}
	admin.Security.LoginURL = params.LoginURL

	err := admin.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	if isServerChanged {
		settings.NotifyServerChange()
	}

	this.Next("/settings", map[string]interface{}{})
	this.Success("保存成功")
}
