package ssl

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/utils/string"
)

type UploadKeyAction actions.Action

func (this *UploadKeyAction) Run(params struct {
	Filename string
	KeyFile  *actions.File
}) {
	// @TODO 校验证书文件格式

	if params.KeyFile == nil {
		this.Fail("请选择证书文件")
	}

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

	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		configFile.Delete()
		this.Fail(err.Error())
	}

	if server.SSL == nil {
		server.SSL = new(teaconfigs.SSLConfig)
	}

	server.SSL.CertificateKey = keyFilename
	server.WriteBack()

	if server.SSL.On && len(server.SSL.Certificate) > 0 {
		global.NotifyChange()
	}

	this.Refresh().Success("保存成功")
}
