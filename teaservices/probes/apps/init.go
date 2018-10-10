package apps

import (
	"time"
)

func init() {
	go func() {
		time.Sleep(1 * time.Second)

		new(MySQLProbe).Run()
		new(MongoDBProbe).Run()
		new(RedisProbe).Run()
		new(NginxProbe).Run()
		new(PHPFPMProbe).Run()
	}()
}
