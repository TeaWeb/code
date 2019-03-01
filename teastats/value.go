package teastats

import (
	"github.com/iwind/TeaGo/utils/time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// 值周期
type ValuePeriod = string

const (
	ValuePeriodSecond = "second"
	ValuePeriodMinute = "minute"
	ValuePeriodHour   = "hour"
	ValuePeriodDay    = "day"
	ValuePeriodWeek   = "week"
	ValuePeriodMonth  = "month"
	ValuePeriodYear   = "year"
)

// 统计指标值定义
type Value struct {
	Id         primitive.ObjectID     `bson:"_id" json:"id"`              // 数据库存储的ID
	Item       string                 `bson:"item" json:"item"`           // 指标代号
	Period     ValuePeriod            `bson:"period" json:"period"`       // 周期
	Value      map[string]interface{} `bson:"value" json:"value"`         // 数据内容
	Params     map[string]string      `bson:"params" json:"params"`       // 参数
	Timestamp  int64                  `bson:"timestamp" json:"timestamp"` // 时间戳
	TimeFormat struct {
		Year   string `bson:"year" json:"year"`
		Month  string `bson:"month" json:"month"`
		Week   string `bson:"week" json:"week"`
		Day    string `bson:"day" json:"day"`
		Hour   string `bson:"hour" json:"hour"`
		Minute string `bson:"minute" json:"minute"`
		Second string `bson:"second" json:"second"`
	} `bson:"timeFormat" json:"timeFormat"`                               // 时间信息
}

// 获取新对象
func NewItemValue() *Value {
	return &Value{}
}

func (this *Value) SetTime(t time.Time) {
	this.Timestamp = t.Unix()
	this.TimeFormat.Year = timeutil.Format("Y", t)
	this.TimeFormat.Month = timeutil.Format("Ym", t)
	this.TimeFormat.Week = timeutil.Format("YW", t)
	this.TimeFormat.Day = timeutil.Format("Ymd", t)
	this.TimeFormat.Hour = timeutil.Format("YmdH", t)
	this.TimeFormat.Minute = timeutil.Format("YmdHi", t)
	this.TimeFormat.Second = timeutil.Format("YmdHis", t)
}
