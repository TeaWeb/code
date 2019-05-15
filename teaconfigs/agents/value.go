package agents

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/utils/time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// 应用指标值定义
type Value struct {
	Id          primitive.ObjectID  `bson:"_id" json:"id"`                  // 数据库存储的ID
	NodeId      string              `bson:"nodeId" json:"nodeId"`           // 节点ID
	AgentId     string              `bson:"agentId" json:"agentId"`         // Agent ID
	AppId       string              `bson:"appId" json:"appId"`             // App ID
	ItemId      string              `bson:"itemId" json:"itemId"`           // 监控项ID
	Timestamp   int64               `bson:"timestamp" json:"timestamp"`     // Agent时间戳
	CreatedAt   int64               `bson:"createdAt" json:"createdAt"`     // Master时间戳
	Value       interface{}         `bson:"value" json:"value"`             // 值，可以是个标量，或者一个组合的值
	Error       string              `bson:"error" json:"error"`             // 错误信息
	NoticeLevel notices.NoticeLevel `bson:"noticeLevel" json:"noticeLevel"` // 通知级别
	IsNotified  bool                `bson:"isNotified" json:"isNotified"`   // 是否已通知
	ThresholdId string              `bson:"thresholdId" json:"thresholdId"` // 阈值ID
	Threshold   string              `bson:"threshold" json:"threshold"`     // 阈值描述
	TimeFormat  struct {
		Year   string `bson:"year" json:"year"`
		Month  string `bson:"month" json:"month"`
		Day    string `bson:"day" json:"day"`
		Hour   string `bson:"hour" json:"hour"`
		Minute string `bson:"minute" json:"minute"`
		Second string `bson:"second" json:"second"`
	} `bson:"timeFormat" json:"timeFormat"`
}

// 获取新对象
func NewValue() *Value {
	return &Value{
		Id: primitive.NewObjectID(),
	}
}

// 设置时间
func (this *Value) SetTime(t time.Time) {
	this.Timestamp = t.Unix()
	this.TimeFormat.Year = timeutil.Format("Y", t)
	this.TimeFormat.Month = timeutil.Format("Ym", t)
	this.TimeFormat.Day = timeutil.Format("Ymd", t)
	this.TimeFormat.Hour = timeutil.Format("YmdH", t)
	this.TimeFormat.Minute = timeutil.Format("YmdHi", t)
	this.TimeFormat.Second = timeutil.Format("YmdHis", t)
}
