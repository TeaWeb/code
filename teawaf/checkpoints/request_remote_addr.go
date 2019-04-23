package checkpoints

import (
	"net/http"
	"strings"
)

type RequestRemoteAddrCheckpoint struct {
	Checkpoint
}

func (this *RequestRemoteAddrCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
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
