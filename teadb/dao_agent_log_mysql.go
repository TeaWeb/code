package teadb

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/logs"
	"golang.org/x/net/context"
)

type MySQLAgentLogDAO struct {
}

// 初始化
func (this *MySQLAgentLogDAO) Init() {

}

func (this *MySQLAgentLogDAO) TableName(agentId string) string {
	table := "logs.agent." + agentId
	this.initTable(table)
	return table
}

// 插入一条数据
func (this *MySQLAgentLogDAO) InsertOne(agentId string, log *agents.ProcessLog) error {
	return NewQuery(this.TableName(agentId)).InsertOne(log)
}

// 获取最新任务的日志
func (this *MySQLAgentLogDAO) FindLatestTaskLogs(agentId string, taskId string, fromId string, size int) ([]*agents.ProcessLog, error) {
	query := NewQuery(this.TableName(agentId)).
		Attr("taskId", taskId).
		Desc("_id")
	if len(fromId) > 0 {
		query.Gt("_id", fromId)
	}
	if size > 0 {
		query.Limit(size)
	}
	ones, err := query.FindOnes(new(agents.ProcessLog))
	if err != nil {
		return nil, err
	}
	result := []*agents.ProcessLog{}
	for _, one := range ones {
		result = append(result, one.(*agents.ProcessLog))
	}
	return result, nil
}

// 获取任务最后一次的执行日志
func (this *MySQLAgentLogDAO) FindLatestTaskLog(agentId string, taskId string) (*agents.ProcessLog, error) {
	query := NewQuery(this.TableName(agentId)).
		Attr("taskId", taskId).
		Desc("_id")
	one, err := query.FindOne(new(agents.ProcessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*agents.ProcessLog), nil
}

func (this *MySQLAgentLogDAO) initTable(table string) {
	if isInitializedTable(table) {
		return
	}

	conn, err := SharedDB().(*MySQLDriver).connect()
	if err != nil {
		return
	}

	_, err = conn.ExecContext(context.Background(), "SHOW CREATE TABLE `"+table+"`")
	if err != nil {
		s := "CREATE TABLE `" + table + "` (" +
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`_id` varchar(24) DEFAULT NULL," +
			"`agentId` varchar(64) DEFAULT NULL," +
			"`taskId` varchar(64) DEFAULT NULL," +
			"`processId` varchar(64) DEFAULT NULL," +
			"`processPid` int(11) unsigned DEFAULT '0'," +
			"`eventType` varchar(32) DEFAULT NULL," +
			"`data` text," +
			"`timestamp` int(11) unsigned DEFAULT '0'," +
			"`timeFormat_year` varchar(4) DEFAULT NULL," +
			"`timeFormat_month` varchar(6) DEFAULT NULL," +
			"`timeFormat_week` varchar(6) DEFAULT NULL," +
			"`timeFormat_day` varchar(8) DEFAULT NULL," +
			"`timeFormat_hour` varchar(10) DEFAULT NULL," +
			"`timeFormat_minute` varchar(12) DEFAULT NULL," +
			"`timeFormat_second` varchar(14) DEFAULT NULL," +
			"PRIMARY KEY (`id`)," +
			"UNIQUE KEY `_id` (`_id`)," +
			"KEY `taskId` (`taskId`) USING BTREE" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
		_, err = conn.ExecContext(context.Background(), s)
		if err != nil {
			logs.Error(err)
		}
	}
}
