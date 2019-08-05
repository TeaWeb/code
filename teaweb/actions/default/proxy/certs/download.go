package certs

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"path/filepath"
)

type DownloadAction actions.Action

// 下载
func (this *DownloadAction) RunGet(params struct {
	CertId string
	Type   string
}) {
	cert := teaconfigs.SharedSSLCertList().FindCert(params.CertId)
	if cert == nil {
		this.WriteString("找不到要查看的证书")
		return
	}

	// 下载证书
	if params.Type == "cert" {
		data, err := cert.ReadCert()
		if err != nil {
			this.WriteString(err.Error())
			return
		}

		this.AddHeader("Content-Disposition", "attachment; filename=\""+filepath.Base(cert.CertFile)+"\";")
		this.Write(data)
		return
	}

	// 下载私钥
	if params.Type == "key" {
		data, err := cert.ReadKey()
		if err != nil {
			this.WriteString(err.Error())
			return
		}

		this.AddHeader("Content-Disposition", "attachment; filename=\""+filepath.Base(cert.KeyFile)+"\";")
		this.Write(data)
		return
	}

	// 查看证书
	if params.Type == "viewCert" {
		data, err := cert.ReadCert()
		if err != nil {
			this.WriteString(err.Error())
			return
		}

		this.Write(data)
		return
	}

	// 查看私钥
	if params.Type == "viewKey" {
		data, err := cert.ReadKey()
		if err != nil {
			this.WriteString(err.Error())
			return
		}

		this.Write(data)
		return
	}
}
