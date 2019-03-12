package notices

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// 通知
type Notice struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"` // 数据库存储的ID
	Proxy      ProxyCond          `bson:"proxy" json:"proxy"`
	Agent      AgentCond          `bson:"agent" json:"agent"`
	Timestamp  int64              `bson:"timestamp" json:"timestamp"` // 时间戳
	Message    string             `bson:"message" json:"message"`
	IsRead     bool               `bson:"isRead" json:"isRead"`         // 已读
	IsNotified bool               `bson:"isNotified" json:"isNotified"` // 是否发送通知
	Receivers  []string           `bson:"receivers" json:"receivers"`   // 接收人ID列表
}

// Proxy条件
type ProxyCond struct {
	ServerId   string `bson:"serverId" json:"serverId"`
	LocationId string `bson:"locationId" json:"serverId"`
	RewriteId  string `bson:"rewriteId" json:"serverId"`
	BackendId  string `bson:"backendId" json:"serverId"`
	FastcgiId  string `bson:"fastcgiId" json:"serverId"`
	Level      uint8  `bson:"level" json:"level"`
}

// Agent条件
type AgentCond struct {
	AgentId   string `bson:"agentId" json:"agentId"`
	AppId     string `bson:"appId" json:"appId"`
	TaskId    string `bson:"taskId" json:"taskId"`
	ItemId    string `bson:"itemId" json:"itemId"`
	Level     uint8  `bson:"level" json:"level"`
	Threshold string `bson:"bson" json:"threshold"`
}

// 获取通知对象
func NewNotice() *Notice {
	return &Notice{}
}

// 设置时间
func (this *Notice) SetTime(t time.Time) {
	this.Timestamp = time.Now().Unix()
}
