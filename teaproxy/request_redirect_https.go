package teaproxy

import (
	"net/http"
	"strings"
)

func (this *Request) callRedirectToHttps(writer *ResponseWriter) {
	// 是否需要跳转到HTTPS
	if this.redirectToHttps && this.rawScheme == "http" {
		host := this.raw.Host
		hostPortIndex := strings.LastIndex(host, ":")
		if hostPortIndex > -1 {
			host = host[:hostPortIndex]
		}

		// 是否有HTTPS
		if this.server.SSL != nil && this.server.SSL.On && len(this.server.SSL.Listen) > 0 {
			listen := this.server.SSL.Listen[0]
			portIndex := strings.LastIndex(listen, ":")
			if portIndex > -1 {
				port := listen[portIndex+1:]
				if port == "443" {
					u := "https://" + host + this.raw.RequestURI
					http.Redirect(writer, this.raw, u, http.StatusTemporaryRedirect)
					return
				} else {
					u := "https://" + host + ":" + port + this.raw.RequestURI
					http.Redirect(writer, this.raw, u, http.StatusTemporaryRedirect)
					return
				}
			} else {
				u := "https://" + host + this.raw.RequestURI
				http.Redirect(writer, this.raw, u, http.StatusTemporaryRedirect)
				return
			}
		}

		u := "https://" + host + this.raw.RequestURI
		http.Redirect(writer, this.raw, u, http.StatusTemporaryRedirect)
		return
	}
}
