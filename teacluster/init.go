package teacluster

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
	"time"
)

func init() {
	if !ClusterEnabled {
		return
	}

	// register actions
	RegisterActionType(
		new(SuccessAction),
		new(FailAction),
		new(RegisterAction),
		new(PushAction),
		new(PullAction),
		new(NotifyAction),
		new(SumAction),
	)

	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// build
		BuildSum()

		// start manager
		go func() {
			ticker := time.NewTicker(60 * time.Second)
			for {
				err := ClusterManager.Start()
				if err != nil {
					logs.Println("[cluster]" + err.Error())
				}

				// retry N seconds later
				select {
				case <-ticker.C:
					// every N seconds
				case <-ClusterManager.Context:
					// retry immediately
				}
			}
		}()
	})

	TeaGo.BeforeStop(func(server *TeaGo.Server) {
		if ClusterManager != nil {
			ClusterManager.Stop()
		}
	})
}
