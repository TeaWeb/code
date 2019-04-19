package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
)

type IndexAction actions.Action

// SSL设置
func (this *IndexAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["selectedTab"] = "https"
	this.Data["server"] = server
	this.Data["errs"] = teaproxy.SharedManager.FindServerErrors(params.ServerId)

	this.Data["error"] = ""
	this.Data["info"] = []maps.Map{}

	if server.SSL != nil && len(server.SSL.Certificate) > 0 && len(server.SSL.CertificateKey) > 0 {
		cert, err := tls.LoadX509KeyPair(Tea.ConfigFile(server.SSL.Certificate), Tea.ConfigFile(server.SSL.CertificateKey))
		if err != nil {
			this.Data["error"] = err.Error()
		} else {
			info := []maps.Map{}
			for _, data := range cert.Certificate {
				c, err := x509.ParseCertificate(data)
				if err != nil {
					this.Data["error"] = err.Error()
				} else {
					info = append(info, maps.Map{
						"subject": c.Subject.CommonName,
						"issuer":  c.Issuer.CommonName,
						"before":  timeutil.Format("Y-m-d", c.NotBefore),
						"after":   timeutil.Format("Y-m-d", c.NotAfter),
					})
				}
			}
			lists.Reverse(info)
			this.Data["info"] = info
		}
	}

	this.Show()
}
