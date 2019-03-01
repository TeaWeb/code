package agentutils

import "go.mongodb.org/mongo-driver/bson/primitive"

// 进程日志
type ProcessLog struct {
	Id         primitive.ObjectID `var:"id" bson:"_id" json:"id"` // 数据库存储的ID
	AgentId    string             `bson:"agentId" json:"agentId"`
	TaskId     string             `bson:"taskId" json:"taskId"`
	ProcessId  string             `bson:"processId" json:"processId"`
	ProcessPid int                `bson:"processPid" json:"processPid"`
	EventType  string             `bson:"eventType" json:"eventType"` // start, log, stop
	Data       string             `bson:"data" json:"data"`
	Timestamp  int64              `bson:"timestamp" json:"timestamp"` // unix时间戳，单位为秒
	TimeFormat struct {
		Year   string `bson:"year" json:"year"`
		Month  string `bson:"month" json:"month"`
		Day    string `bson:"day" json:"day"`
		Hour   string `bson:"hour" json:"hour"`
		Minute string `bson:"minute" json:"minute"`
		Second string `bson:"second" json:"second"`
	} `bson:"timeFormat" json:"timeFormat"`
}
