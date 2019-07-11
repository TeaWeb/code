package teautils

import "unsafe"

// convert bytes to string
func BytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
