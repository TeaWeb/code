package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"strings"
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
	this.Data["warning"] = ""
	this.Data["info"] = []maps.Map{}

	if server.SSL != nil && len(server.SSL.Certificate) > 0 && len(server.SSL.CertificateKey) > 0 {
		cert, err := tls.LoadX509KeyPair(Tea.ConfigFile(server.SSL.Certificate), Tea.ConfigFile(server.SSL.CertificateKey))
		if err != nil {
			this.Data["error"] = err.Error()
		} else {
			info := []maps.Map{}
			allDnsNames := []string{}
			for _, data := range cert.Certificate {
				c, err := x509.ParseCertificate(data)
				if err != nil {
					this.Data["error"] = err.Error()
				} else {
					dnsNames := ""
					if len(c.DNSNames) > 0 {
						dnsNames = "[" + strings.Join(c.DNSNames, ", ") + "]"
						allDnsNames = append(allDnsNames, c.DNSNames...)
					}
					info = append(info, maps.Map{
						"subject":  c.Subject.CommonName + " " + dnsNames,
						"issuer":   c.Issuer.CommonName,
						"before":   timeutil.Format("Y-m-d", c.NotBefore),
						"after":    timeutil.Format("Y-m-d", c.NotAfter),
						"dnsNames": dnsNames,
					})
				}
			}
			lists.Reverse(info)
			this.Data["info"] = info

			// 检查域名是否设置
			if len(allDnsNames) > 0 {
				// 检查domain
				for _, domain := range allDnsNames {
					if !teautils.MatchDomains(server.Name, domain) {
						this.Data["warning"] = "当前代理服务的域名中没有域名可以匹配\"" + domain + "\"，请在代理服务的<a href=\"/proxy/update?serverId=" + server.Id + "\">基本信息</a>中添加此域名。"
					}
				}
			}
		}
	}

	this.Show()
}
