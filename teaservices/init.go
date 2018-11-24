package teaservices

import (
	"github.com/TeaWeb/code/teaservices/probes"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		logs.Println("[services]start service probes")

		go func() {
			new(probes.CPUProbe).Run()
			new(probes.MemoryProbe).Run()
			new(probes.NetworkProbe).Run()
			new(probes.DiskProbe).Run()
		}()
	})
}
