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
	addr = strings.ReplaceAll(addr, " ", "")
	addr = strings.ReplaceAll(addr, "\t", "")
	addr = strings.ReplaceAll(addr, "：", ":")
	return addr
}
