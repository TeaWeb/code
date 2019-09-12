package teautils

import (
	"strings"
	"unsafe"
)

// convert bytes to string
func BytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// format address
func FormatAddress(addr string) string {
	if strings.HasSuffix(addr, "unix:") {
		return addr
	}
	addr = strings.ReplaceAll(addr, " ", "")
	addr = strings.ReplaceAll(addr, "\t", "")
	addr = strings.ReplaceAll(addr, "：", ":")
	addr = strings.TrimSpace(addr)
	return addr
}

// format address list
func FormatAddressList(addrList []string) []string {
	result := []string{}
	for _, addr := range addrList {
		result = append(result, FormatAddress(addr))
	}
	return result
}
