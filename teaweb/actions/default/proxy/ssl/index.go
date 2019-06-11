package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/TeaWeb/code/teautils"
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

	this.Data["minVersion"] = "TLS 1.0"
	if server.SSL != nil && len(server.SSL.MinVersion) > 0 {
		this.Data["minVersion"] = server.SSL.MinVersion
	}

	this.Data["selectedTab"] = "https"
	this.Data["server"] = server
	this.Data["errs"] = teaproxy.SharedManager.FindServerErrors(params.ServerId)

	errorMessages := []string{}
	warningMessages := []string{}
	certs := []maps.Map{}

	notMatchedDomains := []string{}

	if server.SSL != nil {
		err := server.SSL.Validate()
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}

		for index, certConfig := range server.SSL.Certs {
			info := []maps.Map{}

			cert, err := tls.LoadX509KeyPair(certConfig.FullCertPath(), certConfig.FullKeyPath())
			if err != nil {
				if server.SSL.On {
					errorMessages = append(errorMessages, fmt.Sprintf("证书#%d：", index+1)+err.Error())
				}
				certs = append(certs, maps.Map{
					"config": certConfig,
					"info":   info,
				})
				continue
			}

			allDnsNames := []string{}
			for _, data := range cert.Certificate {
				c, err := x509.ParseCertificate(data)
				if err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("证书#%d：", index+1)+err.Error())
					certs = append(certs, maps.Map{
						"config": certConfig,
						"info":   info,
					})
					continue
				}
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

			lists.Reverse(info)
			certs = append(certs, maps.Map{
				"config": certConfig,
				"info":   info,
			})

			// 检查域名是否设置
			if len(allDnsNames) > 0 {
				// 检查domain
				for _, domain := range allDnsNames {
					if !teautils.MatchDomains(server.Name, domain) {
						if !lists.ContainsString(notMatchedDomains, domain) {
							notMatchedDomains = append(notMatchedDomains, fmt.Sprintf("证书#%d：", index+1)+domain)
						}
					}
				}
			}
		}
	}

	if len(notMatchedDomains) > 0 {
		warningMessages = append(warningMessages, "当前代理服务的已设置域名和证书中的域名不匹配：<br/><div class=\"ui segment\" style=\"margin:0.6em 0;line-height: 1.8;padding-top:0;padding-bottom:0\">"+strings.Join(notMatchedDomains, "<br/>")+"</div>请在代理服务的<a href=\"/proxy/update?serverId="+server.Id+"\">基本信息</a>中添加这些域名。")
	}

	this.Data["errorMessages"] = errorMessages
	this.Data["warningMessages"] = warningMessages
	this.Data["certs"] = certs

	this.Show()
}
