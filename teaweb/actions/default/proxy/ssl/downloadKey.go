package ssl

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"io/ioutil"
	"path/filepath"
)

type DownloadKeyAction actions.Action

func (this *DownloadKeyAction) RunGet(params struct {
	ServerId string
	View     bool
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.SSL == nil {
		this.Fail("还没有配置SSL")
	}

	cert := server.SSL.CertificateKey
	if len(cert) == 0 {
		this.Fail("没有设置密钥文件")
	}

	data, err := ioutil.ReadFile(Tea.ConfigFile(cert))
	if err != nil {
		this.WriteString(err.Error())
		return
	}

	if params.View { // 在线浏览
		this.Write(data)
	} else { // 下载
		this.AddHeader("Content-Disposition", "attachment; filename=\""+filepath.Base(cert)+"\";")
		this.Write(data)
	}
}
