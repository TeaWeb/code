package certutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"strings"
)

// 查找使用服务的证书
func FindAllServersUsingCert(certId string) []*teaconfigs.ServerConfig {
	serverList, err := teaconfigs.SharedServerList()
	if err != nil || serverList == nil {
		return nil
	}
	result := []*teaconfigs.ServerConfig{}
	if serverList != nil {
		for _, server := range serverList.FindAllServers() {
			if server.SSL == nil {
				continue
			}
			for _, c := range server.SSL.Certs {
				if c.Id == certId {
					result = append(result, server)
					break
				}
			}
		}
	}
	return result
}

// 以Map的形式列出所有证书
func ListAllCertsMap() []interface{} {
	return lists.Map(teaconfigs.SharedSSLCertList().Certs, func(k int, v interface{}) interface{} {
		cert := v.(*teaconfigs.SSLCertConfig)
		err := cert.Validate()

		errorString := ""
		if err != nil {
			errorString = err.Error()
		}

		summary := cert.Description
		dnsNames := cert.DNSNames()
		if len(dnsNames) > 0 {
			if len(dnsNames) > 2 {
				summary += " (" + strings.Join(dnsNames[:2], ",") + "等"
			} else {
				summary += " (" + strings.Join(dnsNames, ",")
			}
			summary += " - " + timeutil.Format("Y-m-d H:i:s", cert.TimeAfter())
			summary += ")"
		}

		return maps.Map{
			"id":          cert.Id,
			"error":       errorString,
			"dnsNames":    cert.DNSNames(),
			"description": cert.Description,
			"summary":     summary,
		}
	})
}
