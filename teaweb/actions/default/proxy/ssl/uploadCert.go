package ssl

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/utils/string"
)

type UploadCertAction actions.Action

func (this *UploadCertAction) Run(params struct {
	Filename string
	CertFile *actions.File
}) {
	// @TODO 校验证书文件格式

	if params.CertFile == nil {
		this.Fail("请选择证书文件")
	}

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

	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		configFile.Delete()
		this.Fail(err.Error())
	}

	if server.SSL == nil {
		server.SSL = new(teaconfigs.SSLConfig)
	}

	server.SSL.Certificate = certFilename
	server.WriteBack()

	if server.SSL.On && len(server.SSL.CertificateKey) > 0 {
		global.NotifyChange()
	}

	this.Refresh().Success("保存成功")
}
