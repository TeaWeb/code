package teastats

import (
	"errors"
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/syndtr/goleveldb/leveldb"
	"strings"
	"time"
)

var ipTopRankDb *leveldb.DB = nil

// 请求数统计
type RequestIPPeriodFilter struct {
	code   string
	period string
	queue  *Queue
	rank   *Rank
	timer  *timers.Looper
	db     *leveldb.DB

	lastHour string
}

func (this *RequestIPPeriodFilter) Name() string {
	return "IP请求数排行"
}

func (this *RequestIPPeriodFilter) Codes() []string {
	return []string{
		"request.ip.hour",
		"request.ip.day",
	}
}

func (this *RequestIPPeriodFilter) Indexes() []string {
	return []string{}
}

func (this *RequestIPPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.code = code
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.period = code[strings.LastIndex(code, ".")+1:]
	this.rank = NewRank(20, 200*10000)

	prefix := ""
	duration := 10 * time.Second
	switch this.period {
	case ValuePeriodHour:
		prefix = timeutil.Format("YmdH")
		duration = 10 * time.Minute
	case ValuePeriodDay:
		prefix = timeutil.Format("Ymd")
		duration = 30 * time.Minute
	}
	prefix += this.queue.ServerId

	if ipTopRankDb == nil {
		db, err := leveldb.OpenFile(Tea.Root+"/logs/top.ip.leveldb", nil)
		if err != nil {
			logs.Error(errors.New("logs/top.ip.leveldb:" + err.Error()))
		} else {
			this.db = db
			this.rank.Load(db, prefix)
			ipTopRankDb = this.db
		}
	} else {
		this.db = ipTopRankDb
		this.rank.Load(ipTopRankDb, prefix)
	}

	this.timer = timers.Loop(duration, func(looper *timers.Looper) {
		this.commit()
	})
}

func (this *RequestIPPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	remoteAddr := accessLog.RemoteAddr
	if len(remoteAddr) == 0 {
		return
	}
	hour := timeutil.Format("YmdH")
	if len(this.lastHour) > 0 && hour != this.lastHour {
		this.commit()
		this.rank.Reset()
		this.lastHour = hour
	} else {
		this.lastHour = hour
	}

	this.rank.Add(remoteAddr)
}

func (this *RequestIPPeriodFilter) Stop() {
	if this.timer != nil {
		this.timer.Stop()
		this.timer = nil
	}

	this.commit()
}

func (this *RequestIPPeriodFilter) commit() {
	if this.db != nil {
		prefix := ""
		switch this.period {
		case ValuePeriodDay:
			prefix = timeutil.Format("Ymd")
		case ValuePeriodHour:
			prefix = timeutil.Format("YmdH")
		}
		prefix += this.queue.ServerId

		this.rank.Save(this.db, prefix)

		top := this.rank.Top()
		this.rank.locker.Lock()
		this.queue.Add(this.code, time.Now(), this.period, nil, maps.Map{
			"top": top,
			"$increase": func(value maps.Map, inc maps.Map) maps.Map {
				return inc
			},
		})
		this.rank.locker.Unlock()
	}
}
