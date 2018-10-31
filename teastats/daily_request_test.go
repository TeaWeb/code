package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/utils/time"
	"testing"
	"time"
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
	total := stat.SumDayRequests("123456", []string{timeutil.Format("Ymd"), "20181026"})
	t.Log(total)
}
