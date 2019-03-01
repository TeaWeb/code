package audits

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// 动作
type Action = string

const (
	ActionLogin = "LOGIN" // 登录 {ip}
)

// 审计日志
type Log struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"` // 数据库存储的ID
	Username    string             `bson:"username" json:"username"`
	Action      Action             `bson:"action" json:"action"`           // 类型
	Description string             `bson:"description" json:"description"` // 描述
	Options     map[string]string  `bson:"options" json:"options"`         // 选项
	Timestamp   int64              `bson:"timestamp" json:"timestamp"`     // 时间戳
}

// 获取新审计日志对象
func NewLog(username string, action Action, description string, options map[string]string) *Log {
	return &Log{
		Id:          primitive.NewObjectID(),
		Username:    username,
		Action:      action,
		Description: description,
		Timestamp:   time.Now().Unix(),
		Options:     options,
	}
}

// 审计日志类型
func (this *Log) ActionName() string {
	switch this.Action {
	case ActionLogin:
		return "登录"
	}
	return ""
}
