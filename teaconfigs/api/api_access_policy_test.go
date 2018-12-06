package api

import (
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/utils/time"
	"log"
	"sync"
	"testing"
	"time"
)

func TestAPIAccessPolicy(t *testing.T) {
	printTime := func(t1 time.Time) {
		log.Println(timeutil.Format("Y-m-d H:i:s", t1))
	}

	printTime2 := func(t1 time.Time, t2 time.Time) {
		log.Println(timeutil.Format("Y-m-d H:i:s", t1), timeutil.Format("Y-m-d H:i:s", t2))
	}

	p := APIAccessPolicy{}
	p.Traffic.On = true
	p.Traffic.Second.On = true
	p.Traffic.Second.Duration = 1

	p.Traffic.Minute.On = true
	p.Traffic.Minute.Duration = 1

	p.Traffic.Hour.On = true
	p.Traffic.Hour.Duration = 1

	p.Traffic.Day.On = true
	p.Traffic.Day.Duration = 1

	p.Traffic.Month.On = true
	p.Traffic.Month.Duration = 1

	for {
		time.Sleep(1000 * time.Millisecond)

		log.Println("===")
		p.IncreaseTraffic()

		log.Println("seconds:", p.Traffic.Second.Used)
		printTime(p.Traffic.Second.fromTime)

		log.Println("minutes:", p.Traffic.Minute.Used)
		printTime2(p.Traffic.Minute.fromTime, p.Traffic.Minute.toTime)

		/**log.Println("hours:", p.Traffic.Hour.Used)
		printTime2(p.Traffic.Hour.fromTime, p.Traffic.Hour.toTime)

		log.Println("days:", p.Traffic.Minute.Used)
		printTime2(p.Traffic.Day.fromTime, p.Traffic.Day.toTime)

		log.Println("months:", p.Traffic.Minute.Used)
		printTime2(p.Traffic.Month.fromTime, p.Traffic.Month.toTime)**/

		break
	}
}

func TestAPIAccessPolicySecond(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	p := APIAccessPolicy{}
	p.Traffic.On = true
	p.Traffic.Second.On = false
	a.IsTrue(p.AllowTraffic())

	p.Traffic.Second.On = true
	a.IsFalse(p.AllowTraffic())

	p.Traffic.Second.On = true
	p.Traffic.Second.Duration = 1
	p.Traffic.Second.Total = 1
	a.IsTrue(p.AllowTraffic())

	p.Traffic.Second.On = true
	p.Traffic.Second.Duration = 1
	p.Traffic.Second.Total = 2
	p.IncreaseTraffic()
	p.IncreaseTraffic()
	a.IsFalse(p.AllowTraffic())

	time.Sleep(1 * time.Second)
	a.IsTrue(p.AllowTraffic())
}

func TestAPIAccessPolicyMinute(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	p := APIAccessPolicy{}
	p.Traffic.On = true
	p.Traffic.Minute.On = true
	a.IsFalse(p.AllowTraffic())

	p.Traffic.Minute.Total = 1
	a.IsFalse(p.AllowTraffic())

	p.Traffic.Minute.Duration = 1
	a.IsTrue(p.AllowTraffic())

	p.IncreaseTraffic()
	a.IsFalse(p.AllowTraffic())

	//time.Sleep(61 * time.Second)
	//a.IsTrue(p.AllowTraffic())
}

func TestAPIAccessPolicyPerformance(t *testing.T) {
	times := 100000
	before := time.Now()

	locker := sync.Mutex{}
	for i := 0; i < times; i ++ {
		locker.Lock()

		p := APIAccessPolicy{}
		p.Traffic.On = true
		p.Traffic.Second.On = true
		p.Traffic.Second.Duration = 1
		p.Traffic.Second.Total = 1
		p.Traffic.Minute.On = true
		p.Traffic.Minute.Duration = 1
		p.Traffic.Minute.Total = 1
		p.Traffic.Hour.On = true
		p.Traffic.Hour.Duration = 1
		p.Traffic.Hour.Total = 1
		p.Traffic.Day.On = true
		p.Traffic.Day.Duration = 1
		p.Traffic.Day.Total = 1
		p.Traffic.Month.On = true
		p.Traffic.Month.Duration = 1
		p.Traffic.Month.Total = 1
		p.AllowTraffic()
		p.IncreaseTraffic()

		locker.Unlock()
	}

	t.Log(int(float64(times) / time.Since(before).Seconds()))
}
