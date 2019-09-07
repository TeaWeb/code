package teadb

import (
	"context"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/logs"
	"strings"
	"time"
)

type MySQLNoticeDAO struct {
}

// 初始化
func (this *MySQLNoticeDAO) Init() {

}

// 表名
func (this *MySQLNoticeDAO) TableName() string {
	this.initTable("notices")
	return "notices"
}

// 写入一个通知
func (this *MySQLNoticeDAO) InsertOne(notice *notices.Notice) error {
	return NewQuery(this.TableName()).
		InsertOne(notice)
}

// 发送一个代理的通知（形式1）
func (this *MySQLNoticeDAO) NotifyProxyMessage(cond notices.ProxyCond, message string) error {
	notice := notices.NewNotice()
	notice.Message = message
	notice.SetTime(time.Now())
	notice.Proxy = cond
	notice.Hash()
	return NewQuery(this.TableName()).InsertOne(notice)
}

// 发送一个代理的通知（形式2）
func (this *MySQLNoticeDAO) NotifyProxyServerMessage(serverId string, level notices.NoticeLevel, message string) error {
	return this.NotifyProxyMessage(notices.ProxyCond{
		ServerId: serverId,
		Level:    level,
	}, message)
}

// 获取所有未读通知数
func (this *MySQLNoticeDAO) CountAllUnreadNotices() (int, error) {
	count, err := NewQuery(this.TableName()).
		Attr("isRead", 0).
		Count()
	return int(count), err
}

// 获取所有已读通知数
func (this *MySQLNoticeDAO) CountAllReadNotices() (int, error) {
	count, err := NewQuery(this.TableName()).
		Attr("isRead", 1).
		Count()
	return int(count), err
}

// 获取某个Agent的未读通知数
func (this *MySQLNoticeDAO) CountUnreadNoticesForAgent(agentId string) (int, error) {
	count, err := NewQuery(this.TableName()).
		Attr("agentId", agentId).
		Attr("isRead", 0).
		Count()
	return int(count), err
}

// 获取某个Agent已读通知数
func (this *MySQLNoticeDAO) CountReadNoticesForAgent(agentId string) (int, error) {
	count, err := NewQuery(this.TableName()).
		Attr("agentId", agentId).
		Attr("isRead", 1).
		Count()
	return int(count), err
}

// 获取某个接收人在某个时间段内接收的通知数
func (this *MySQLNoticeDAO) CountReceivedNotices(receiverId string, cond map[string]interface{}, minutes int) (int, error) {
	if len(receiverId) == 0 {
		return 0, nil
	}
	if minutes <= 0 {
		return 0, nil
	}
	query := NewQuery(this.TableName()).
		sqlCond("FIND_IN_SET(:receiverId, receivers)>0", map[string]interface{}{
			"receiverId": receiverId,
		}).
		Gte("timestamp", time.Now().Unix()-int64(minutes*60))

	if len(cond) > 0 {
		for k, v := range cond {
			k = this.mapField(k)
			query.Attr(k, v)
		}
	}
	c, err := query.Count()
	return int(c), err
}

// 通过Hash判断是否存在相同的消息
func (this *MySQLNoticeDAO) ExistNoticesWithHash(hash string, cond map[string]interface{}, duration time.Duration) (bool, error) {
	query := NewQuery(this.TableName())
	query.Attr("messageHash", hash)
	for k, v := range cond {
		query.Attr(this.mapField(k), v)
	}
	query.Gt("timestamp", float64(time.Now().Unix())-duration.Seconds())
	query.Desc("_id")
	one, err := query.FindOne(new(notices.Notice))
	if err != nil {
		return false, err
	}
	if one == nil {
		return false, nil
	}
	notice := one.(*notices.Notice)

	// 中间是否有success级别的
	query2 := NewQuery(this.TableName())
	for k, v := range cond {
		query2.Attr(this.mapField(k), v)
	}
	if len(notice.Proxy.ServerId) > 0 {
		query2.Attr("level", notices.NoticeLevelSuccess)
		query2.Gt("_id", notice.Id.Hex())
	} else if len(notice.Agent.AgentId) > 0 {
		query2.Attr("level", notices.NoticeLevelSuccess)
		query2.Gt("_id", notice.Id.Hex())
	}
	result, err := query2.Result("_id").
		FindOne(new(notices.Notice))
	return result == nil, err
}

// 列出消息
func (this *MySQLNoticeDAO) ListNotices(isRead bool, offset int, size int) ([]*notices.Notice, error) {
	ones, err := NewQuery(this.TableName()).
		Attr("isRead", isRead).
		Offset(offset).
		Limit(size).
		Desc("_id").
		FindOnes(new(notices.Notice))
	if err != nil {
		return nil, err
	}

	result := []*notices.Notice{}
	for _, one := range ones {
		result = append(result, one.(*notices.Notice))
	}
	return result, err
}

// 列出某个Agent相关的消息
func (this *MySQLNoticeDAO) ListAgentNotices(agentId string, isRead bool, offset int, size int) ([]*notices.Notice, error) {
	ones, err := NewQuery(this.TableName()).
		Attr("agentId", agentId).
		Attr("isRead", isRead).
		Offset(offset).
		Limit(size).
		Desc("_id").
		FindOnes(new(notices.Notice))
	if err != nil {
		return nil, err
	}

	result := []*notices.Notice{}
	for _, one := range ones {
		result = append(result, one.(*notices.Notice))
	}
	return result, err
}

// 删除Agent相关通知
func (this *MySQLNoticeDAO) DeleteNoticesForAgent(agentId string) error {
	return NewQuery(this.TableName()).
		Attr("agentId", agentId).
		Delete()
}

// 更改某个通知的接收人
func (this *MySQLNoticeDAO) UpdateNoticeReceivers(noticeId string, receiverIds []string) error {
	return SharedDB().(*MySQLDriver).UpdateOnes(NewQuery(this.TableName()).Attr("_id", noticeId), map[string]interface{}{
		"isNotified": 1,
		"receivers":  strings.Join(receiverIds, ","),
	})
}

// 设置全部已读
func (this *MySQLNoticeDAO) UpdateAllNoticesRead() error {
	return SharedDB().(*MySQLDriver).UpdateOnes(NewQuery(this.TableName()), map[string]interface{}{
		"isRead": 1,
	})
}

// 设置一组通知已读
func (this *MySQLNoticeDAO) UpdateNoticesRead(noticeIds []string) error {
	if len(noticeIds) == 0 {
		return nil
	}

	query := NewQuery(this.TableName()).
		Attr("_id", noticeIds)
	return SharedDB().(*MySQLDriver).UpdateOnes(query, map[string]interface{}{
		"isRead": 1,
	})
}

// 设置Agent的一组通知已读
func (this *MySQLNoticeDAO) UpdateAgentNoticesRead(agentId string, noticeIds []string) error {
	if len(noticeIds) == 0 {
		return nil
	}

	query := NewQuery(this.TableName()).
		Attr("_id", noticeIds).
		Attr("agentId", agentId)
	return SharedDB().(*MySQLDriver).UpdateOnes(query, map[string]interface{}{
		"isRead": 1,
	})
}

// 设置Agent所有通知已读
func (this *MySQLNoticeDAO) UpdateAllAgentNoticesRead(agentId string) error {
	query := NewQuery(this.TableName()).
		Attr("agentId", agentId)
	return SharedDB().(*MySQLDriver).UpdateOnes(query, map[string]interface{}{
		"isRead": 1,
	})
}

func (this *MySQLNoticeDAO) initTable(table string) {
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
			"`timestamp` int(11) unsigned DEFAULT '0'," +
			"`message` varchar(1024) DEFAULT NULL," +
			"`messageHash` varchar(64) DEFAULT NULL," +
			"`isRead` tinyint(1) unsigned DEFAULT '0'," +
			"`isNotified` tinyint(1) unsigned DEFAULT '0'," +
			"`receivers` varchar(1024) DEFAULT NULL," +
			"`proxyServerId` varchar(64) DEFAULT NULL," +
			"`proxyWebsocket` tinyint(1) unsigned DEFAULT '0'," +
			"`proxyLocationId` varchar(64) DEFAULT '0'," +
			"`proxyRewriteId` varchar(64) DEFAULT NULL," +
			"`proxyBackendId` varchar(64) DEFAULT NULL," +
			"`proxyFastcgiId` varchar(64) DEFAULT NULL," +
			"`level` tinyint(1) unsigned DEFAULT '0'," +
			"`agentId` varchar(64) DEFAULT NULL," +
			"`agentAppId` varchar(64) DEFAULT NULL," +
			"`agentTaskId` varchar(64) DEFAULT NULL," +
			"`agentItemId` varchar(64) DEFAULT NULL," +
			"`agentThreshold` varchar(1024) DEFAULT NULL," +
			"PRIMARY KEY (`id`)," +
			"UNIQUE KEY `_id` (`_id`)," +
			"KEY `messageHash` (`messageHash`)," +
			"KEY `agentId` (`agentId`)," +
			"KEY `isRead` (`isRead`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
		_, err = conn.ExecContext(context.Background(), s)
		if err != nil {
			logs.Error(err)
		}
	}
}

// 字段映射
func (this *MySQLNoticeDAO) mapField(k string) string {
	switch k {
	case "agent.agentId":
		k = "agentId"
	case "agent.appId":
		k = "agentAppId"
	case "agent.itemId":
		k = "agentItemId"
	case "agent.level":
		k = "level"
	case "agent.threshold":
		k = "agentThreshold"
	case "proxy.serverId":
		k = "proxyServerId"
	case "proxy.websocket":
		k = "proxyWebsocket"
	case "proxy.locationId":
		k = "proxyLocationId"
	case "proxy.rewriteId":
		k = "proxyRewriteId"
	case "proxy.fastcgiId":
		k = "proxyFastcgiId"
	case "proxy.backendId":
		k = "proxyBackendId"
	case "proxy.level":
		k = "level"
	}
	return k
}
