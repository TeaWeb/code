package ssl

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/utils/string"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["selectedTab"] = "https"
	this.Data["server"] = server

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	ServerId string
	HttpsOn  bool
	Listen   []string
	CertFile *actions.File
	KeyFile  *actions.File
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.SSL == nil {
		server.SSL = teaconfigs.NewSSLConfig()
	}
	server.SSL.On = params.HttpsOn
	server.SSL.Listen = params.Listen

	if params.CertFile != nil {
		data, err := params.CertFile.Read()
		if err != nil {
			this.Fail(err.Error())
		}

		certFilename := "ssl." + stringutil.Rand(16) + params.CertFile.Ext
		configFile := files.NewFile(Tea.ConfigFile(certFilename))
		err = configFile.Write(data)
		if err != nil {
			this.Fail(err.Error())
		}

		server.SSL.Certificate = certFilename
	}

	if params.KeyFile != nil {
		data, err := params.KeyFile.Read()
		if err != nil {
			this.Fail(err.Error())
		}

		keyFilename := "ssl." + stringutil.Rand(16) + params.KeyFile.Ext
		configFile := files.NewFile(Tea.ConfigFile(keyFilename))
		err = configFile.Write(data)
		if err != nil {
			this.Fail(err.Error())
		}

		server.SSL.CertificateKey = keyFilename
	}

	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
