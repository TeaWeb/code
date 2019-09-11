package teadb

// 表格统计
type TableStat struct {
	Count         int64  `json:"count"`
	Size          int64  `json:"size"`
	FormattedSize string `json:"formattedSize"`
}
