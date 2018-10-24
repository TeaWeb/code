package teaapps

// 内存使用
type MemoryUsage struct {
	RSS     uint64  `json:"rss"`     // RSS
	VMS     uint64  `json:"vms"`     // VMS
	Percent float64 `json:"percent"` // 百分比
}
