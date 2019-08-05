package certutils

import "github.com/TeaWeb/code/teaconfigs"

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
