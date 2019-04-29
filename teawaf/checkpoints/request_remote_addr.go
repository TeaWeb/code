package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"strings"
)

type RequestRemoteAddrCheckpoint struct {
	Checkpoint
}

func (this *RequestRemoteAddrCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	// X-Forwarded-For
	forwardedFor := req.Header.Get("X-Forwarded-For")
	if len(forwardedFor) > 0 {
		index := strings.LastIndex(forwardedFor, ":")
		if index < 0 {
			value = forwardedFor
			return
		} else {
			value = forwardedFor[:index]
			return
		}
	}

	// Real-IP
	realIP := req.Header.Get("X-Real-IP")
	if len(realIP) > 0 {
		index := strings.LastIndex(realIP, ":")
		if index < 0 {
			value = realIP
		} else {
			value = realIP[:index]
		}
		return
	}

	// Real-Ip
	realIP = req.Header.Get("X-Real-Ip")
	if len(realIP) > 0 {
		index := strings.LastIndex(realIP, ":")
		if index < 0 {
			value = realIP
		} else {
			value = realIP[:index]
		}
		return
	}

	// Remote-Addr
	remoteAddr := req.RemoteAddr
	index := strings.LastIndex(remoteAddr, ":")
	if index < 0 {
		value = remoteAddr
	} else {
		value = remoteAddr[:index]
	}
	return
}

func (this *RequestRemoteAddrCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
