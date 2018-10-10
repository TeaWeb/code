package probes

import "fmt"

func formatBytes(size uint64) string {
	if size == 0 {
		return "-"
	}
	if size < 1024 {
		return fmt.Sprintf("%d字节", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%.2fK", float64(size)/1024)
	}
	if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2fM", float64(size)/1024/1024)
	}
	return fmt.Sprintf("%.2fG", float64(size)/1024/1024/1024)
}
