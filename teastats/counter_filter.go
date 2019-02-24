package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
	"strings"
	"sync"
	"time"
)

// 数值增长型的
type CounterFilter struct {
	looper   *timers.Looper
	queue    *Queue
	code     string
	Period   ValuePeriod
	params   map[string]string
	value    maps.Map
	lastTime time.Time
	locker   sync.Mutex

	IncreaseFunc func(value maps.Map, inc maps.Map) maps.Map
}

// 启动筛选器
func (this *CounterFilter) StartFilter(code string, period ValuePeriod) {
	this.code = code
	this.Period = period
	this.lastTime = time.Now()

	// 自动导入
	duration := 1 * time.Second
	switch this.Period {
	case ValuePeriodSecond:
		duration = 1 * time.Second
	case ValuePeriodMinute:
		duration = 1 * time.Minute
	case ValuePeriodHour:
		duration = 5 * time.Minute
	case ValuePeriodDay:
		duration = 10 * time.Minute
	case ValuePeriodWeek:
		duration = 15 * time.Minute
	case ValuePeriodMonth:
		duration = 20 * time.Minute
	case ValuePeriodYear:
		duration = 30 * time.Minute
	}
	this.looper = timers.Loop(duration, func(looper *timers.Looper) {
		this.locker.Lock()
		defer this.locker.Unlock()

		this.commit()
	})
}

// 应用筛选器
func (this *CounterFilter) ApplyFilter(accessLog *tealogs.AccessLog, params map[string]string, value map[string]interface{}) {
	this.locker.Lock()
	defer this.locker.Unlock()

	switch this.Period {
	case ValuePeriodSecond:
		if accessLog.Timestamp == this.lastTime.Unix() && this.equalParams(this.params, params) {
			this.increase(value)
		} else {
			this.commit()
			this.params = params
			this.increase(value)
			this.lastTime = accessLog.Time()
		}
	case ValuePeriodMinute:
		if accessLog.Timestamp/60 == this.lastTime.Unix()/60 && this.equalParams(this.params, params) {
			this.increase(value)
		} else {
			this.commit()
			this.params = params
			this.increase(value)
			this.lastTime = accessLog.Time()
		}
	case ValuePeriodHour:
		if accessLog.Timestamp/3600 == this.lastTime.Unix()/3600 && this.equalParams(this.params, params) {
			this.increase(value)
		} else {
			this.commit()
			this.params = params
			this.increase(value)
			this.lastTime = accessLog.Time()
		}
	case ValuePeriodDay:
		at := accessLog.Time()
		if at.Year() == this.lastTime.Year() && at.Month() == this.lastTime.Month() && at.Day() == this.lastTime.Day() && this.equalParams(this.params, params) {
			this.increase(value)
		} else {
			this.commit()
			this.params = params
			this.increase(value)
			this.lastTime = accessLog.Time()
		}
	case ValuePeriodWeek:
		at := accessLog.Time()
		year1, week1 := at.ISOWeek()
		year2, week2 := this.lastTime.ISOWeek()
		if year1 == year2 && week1 == week2 && this.equalParams(this.params, params) {
			this.increase(value)
		} else {
			this.commit()
			this.params = params
			this.increase(value)
			this.lastTime = accessLog.Time()
		}
	case ValuePeriodMonth:
		at := accessLog.Time()
		if at.Year() == this.lastTime.Year() && at.Month() == this.lastTime.Month() && this.equalParams(this.params, params) {
			this.increase(value)
		} else {
			this.commit()
			this.params = params
			this.increase(value)
			this.lastTime = accessLog.Time()
		}
	case ValuePeriodYear:
		at := accessLog.Time()
		if at.Year() == this.lastTime.Year() && this.equalParams(this.params, params) {
			this.increase(value)
		} else {
			this.commit()
			this.params = params
			this.increase(value)
			this.lastTime = accessLog.Time()
		}
	}
}

// 停止筛选器
func (this *CounterFilter) StopFilter() {
	if this.looper != nil {
		this.looper.Stop()
		this.looper = nil
	}

	this.commit()
}

// 检查新UV
func (this *CounterFilter) CheckNewUV(accessLog *tealogs.AccessLog, attachKey string) bool {
	// 从cookie中检查UV是否存在
	uid, ok := accessLog.Cookie["TeaUID"]
	if !ok || len(uid) == 0 {
		ip := accessLog.RemoteAddr
		if len(ip) == 0 {
			return false
		}
		lastIndex := strings.LastIndex(ip, ":")
		if lastIndex > -1 {
			ip = ip[:lastIndex]
		}
		uid = ip
	}
	l := len(uid)
	if l == 0 {
		return false
	}
	if l > 32 {
		uid = uid[:32]
	}

	key := ""
	life := time.Second
	switch this.Period {
	case ValuePeriodSecond:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("YmdHis", accessLog.Time())
	case ValuePeriodMinute:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("YmdHi", accessLog.Time())
		life = 2 * time.Minute
	case ValuePeriodHour:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("YmdH", accessLog.Time())
		life = 2 * time.Hour
	case ValuePeriodDay:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("Ymd", accessLog.Time())
		life = 2 * 24 * time.Hour
	case ValuePeriodWeek:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("YW", accessLog.Time())
		life = 8 * 24 * time.Hour
	case ValuePeriodMonth:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("Ym", accessLog.Time())
		life = 32 * 24 * time.Hour
	case ValuePeriodYear:
		key = this.code + "_" + this.queue.ServerId + "_" + uid + "_" + timeutil.Format("Y", accessLog.Time())
		life = 370 * 24 * time.Hour
	}

	if len(attachKey) > 0 {
		key += attachKey
	}
	hasKey, err := sharedKV.Has(key)
	if err != nil {
		logs.Error(err)
		return false
	}
	if hasKey {
		return false
	}

	err = sharedKV.Set(key, "1", life)
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}

// 检查新IP
func (this *CounterFilter) CheckNewIP(accessLog *tealogs.AccessLog, attachKey string) bool {
	// IP是否存在
	ip := accessLog.RemoteAddr
	if len(ip) == 0 {
		return false
	}
	lastIndex := strings.LastIndex(ip, ":")
	if lastIndex > -1 {
		ip = ip[:lastIndex]
	}
	key := ""
	life := time.Second
	switch this.Period {
	case ValuePeriodSecond:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("YmdHis", accessLog.Time())
	case ValuePeriodMinute:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("YmdHi", accessLog.Time())
		life = 2 * time.Minute
	case ValuePeriodHour:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("YmdH", accessLog.Time())
		life = 2 * time.Hour
	case ValuePeriodDay:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("Ymd", accessLog.Time())
		life = 2 * 24 * time.Hour
	case ValuePeriodWeek:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("YW", accessLog.Time())
		life = 8 * 24 * time.Hour
	case ValuePeriodMonth:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("Ym", accessLog.Time())
		life = 32 * 24 * time.Hour
	case ValuePeriodYear:
		key = this.code + "_" + this.queue.ServerId + "_" + ip + "_" + timeutil.Format("Y", accessLog.Time())
		life = 370 * 24 * time.Hour
	}

	if len(attachKey) > 0 {
		key += attachKey
	}
	hasKey, err := sharedKV.Has(key)
	if err != nil {
		logs.Error(err)
		return false
	}
	if hasKey {
		return false
	}

	err = sharedKV.Set(key, "1", life)
	if err != nil {
		logs.Error(err)
		return false
	}

	return true
}

// 提交
func (this *CounterFilter) commit() {
	if this.value != nil && len(this.value) > 0 {
		if this.IncreaseFunc != nil {
			this.value["$increase"] = this.IncreaseFunc
		}
		this.queue.Add(this.code, this.lastTime, this.Period, this.params, this.value)
		this.value = nil
	}
}

// 判断参数是否相等
func (this *CounterFilter) equalParams(params1 map[string]string, params2 map[string]string) bool {
	if params1 == nil && params2 == nil {
		return true
	}
	if len(params1) != len(params2) {
		return false
	}
	for k, v := range params1 {
		v1, ok := params2[k]
		if !ok {
			return false
		}
		if v != v1 {
			return false
		}
	}
	return true
}

// 增加值
// 只支持int, int64, float32, float64
func (this *CounterFilter) increase(inc maps.Map) {
	if inc == nil {
		return
	}
	if this.value == nil {
		this.value = inc
		return
	}
	for k, v := range inc {
		v1, ok := this.value[k]
		if !ok {
			this.value[k] = v
			continue
		}
		switch v2 := v1.(type) {
		case int:
			v1 = v2 + types.Int(v)
		case int64:
			v1 = v2 + types.Int64(v)
		case float32:
			v1 = v2 + types.Float32(v)
		case float64:
			v1 = v2 + types.Float64(v)
		}
		this.value[k] = v1
	}
}
