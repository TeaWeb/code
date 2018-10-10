package teastats

import (
	"testing"
	"github.com/TeaWeb/code/tealogs"
	"time"
	"github.com/iwind/TeaGo/utils/time"
)

func TestDailyRequestParse(t *testing.T) {
	log := &tealogs.AccessLog{
		ServerId: "123456",
	}

	stat := new(DailyRequestsStat)
	stat.Process(log)

	time.Sleep(1 * time.Second)
}

func TestDailyPVStat_SumDayPV(t *testing.T) {
	stat := new(DailyRequestsStat)
	total := stat.SumDayRequests([]string{timeutil.Format("Ymd")})
	t.Log(total)
}
