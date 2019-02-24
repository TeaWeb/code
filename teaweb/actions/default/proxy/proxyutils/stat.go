package proxyutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teastats"
	"github.com/iwind/TeaGo/lists"
)

// 刷新服务统计
func ReloadServerStats(serverId string) {
	server := teaconfigs.NewServerConfigFromId(serverId)
	if server == nil || !server.On {
		teastats.RestartServerFilters(serverId, nil)
		return
	}

	codes := []string{}
	for _, board := range []*teaconfigs.Board{server.RealtimeBoard, server.StatBoard} {
		if board == nil {
			continue
		}
		for _, c := range board.Charts {
			_, chart := c.FindChart()
			if chart == nil || !chart.On {
				continue
			}
			for _, r := range chart.Requirements {
				if lists.Contains(codes, r) {
					continue
				}
				codes = append(codes, r)
			}
		}
	}
	teastats.RestartServerFilters(serverId, codes)
}
