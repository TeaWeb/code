package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
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

	if server.SSL != nil {
		err := server.SSL.Validate()
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}

		for index, certConfig := range server.SSL.Certs {
			info := []maps.Map{}

			cert, err := tls.LoadX509KeyPair(Tea.ConfigFile(certConfig.CertFile), Tea.ConfigFile(certConfig.KeyFile))
			if err != nil {
				errorMessages = append(errorMessages, fmt.Sprintf("证书#%d：", index+1)+err.Error())
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
						warningMessages = append(warningMessages, fmt.Sprintf("证书#%d：", index+1)+"当前代理服务的域名中没有域名可以匹配\""+domain+"\"，请在代理服务的<a href=\"/proxy/update?serverId="+server.Id+"\">基本信息</a>中添加此域名。")
					}
				}
			}
		}
	}

	this.Data["errorMessages"] = errorMessages
	this.Data["warningMessages"] = warningMessages
	this.Data["certs"] = certs

	this.Show()
}
