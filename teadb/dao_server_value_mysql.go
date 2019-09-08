package teadb

import (
	"context"
	"encoding/json"
	"github.com/TeaWeb/code/teaconfigs/stats"
	"github.com/TeaWeb/code/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	"time"
)

type MySQLServerValueDAO struct {
}

// 初始化
func (this *MySQLServerValueDAO) Init() {

}

// 表名
func (this *MySQLServerValueDAO) TableName(serverId string) string {
	this.initTable("teaweb.values.server." + serverId)
	return "teaweb.values.server." + serverId
}

// 插入新数据
func (this *MySQLServerValueDAO) InsertOne(serverId string, value *stats.Value) error {
	return NewQuery(this.TableName(serverId)).
		InsertOne(value)
}

// 删除过期的数据
func (this *MySQLServerValueDAO) DeleteExpiredValues(serverId string, period stats.ValuePeriod, life int) error {
	return NewQuery(this.TableName(serverId)).
		Attr("period", period).
		Lt("timestamp", time.Now().Unix()-int64(life)).
		Delete()
}

// 查询相同的数值记录
func (this *MySQLServerValueDAO) FindSameItemValue(serverId string, item *stats.Value) (*stats.Value, error) {
	query := NewQuery(this.TableName(serverId))
	query.Attr("item", item.Item)
	query.Attr("period", item.Period)

	switch item.Period {
	case stats.ValuePeriodSecond:
		query.Attr("timestamp", item.Timestamp)
	case stats.ValuePeriodMinute:
		query.Attr("timeFormat_minute", item.TimeFormat.Minute)
	case stats.ValuePeriodHour:
		query.Attr("timeFormat_hour", item.TimeFormat.Hour)
	case stats.ValuePeriodDay:
		query.Attr("timeFormat_day", item.TimeFormat.Day)
	case stats.ValuePeriodWeek:
		query.Attr("timeFormat_week", item.TimeFormat.Week)
	case stats.ValuePeriodMonth:
		query.Attr("timeFormat_month", item.TimeFormat.Month)
	case stats.ValuePeriodYear:
		query.Attr("timeFormat_year", item.TimeFormat.Year)
	}

	// 参数
	if len(item.Params) > 0 {
		for k, v := range item.Params {
			query.Attr("JSON_EXTRACT(params, '$."+k+"')", v)
		}
	} else {
		query.Attr("JSON_LENGTH(params)", 0)
	}

	one, err := query.FindOne(new(stats.Value))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*stats.Value), nil
}

// 修改值和时间戳
func (this *MySQLServerValueDAO) UpdateItemValueAndTimestamp(serverId string, valueId string, value map[string]interface{}, timestamp int64) error {
	query := NewQuery(this.TableName(serverId)).
		Attr("_id", valueId)
	valuesJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return SharedDB().(*MySQLDriver).UpdateOnes(query, map[string]interface{}{
		"value": valuesJSON,
	})
}

// 创建索引
func (this *MySQLServerValueDAO) CreateIndex(serverId string, fields []*shared.IndexField) error {
	// MySQL无法为JSON字段创建索引，除非是虚拟字段，所以暂时不实现
	return nil
}

// 查询数据
func (this *MySQLServerValueDAO) QueryValues(query *Query) ([]*stats.Value, error) {
	query.fieldMapping = this.mapField
	ones, err := query.FindOnes(new(stats.Value))
	if err != nil {
		return nil, err
	}
	result := []*stats.Value{}
	for _, one := range ones {
		result = append(result, one.(*stats.Value))
	}
	return result, err
}

// 根据item查找一条数据
func (this *MySQLServerValueDAO) FindOneWithItem(serverId string, item string) (*stats.Value, error) {
	one, err := NewQuery(this.TableName(serverId)).
		Attr("item", item).
		FindOne(new(stats.Value))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, err
	}
	return one.(*stats.Value), nil
}

// 删除代理服务相关表
func (this *MySQLServerValueDAO) DropServerTable(serverId string) error {
	return SharedDB().(*MySQLDriver).DropTable(this.TableName(serverId))
}

func (this *MySQLServerValueDAO) initTable(table string) {
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
			"`item` varchar(256) DEFAULT NULL," +
			"`period` varchar(64) DEFAULT NULL," +
			"`value` json DEFAULT NULL," +
			"`params` json DEFAULT NULL," +
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
			"KEY `item_timestamp` (`item`,`timestamp`)," +
			"KEY `item_second` (`item`,`timeFormat_second`)," +
			"KEY `item_minute` (`item`,`timeFormat_minute`)," +
			"KEY `item_hour` (`item`,`timeFormat_hour`)," +
			"KEY `item_day` (`item`,`timeFormat_day`)," +
			"KEY `item_week` (`item`,`timeFormat_week`)," +
			"KEY `item_month` (`item`,`timeFormat_month`)," +
			"KEY `item_year` (`item`,`timeFormat_year`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
		_, err = conn.ExecContext(context.Background(), s)
		if err != nil {
			logs.Error(err)
		}
	}
}

func (this *MySQLServerValueDAO) mapField(field string) string {
	switch field {
	case "timeFormat.year":
		return "timeFormat_year"
	case "timeFormat.month":
		return "timeFormat_month"
	case "timeFormat.week":
		return "timeFormat_week"
	case "timeFormat.day":
		return "timeFormat_day"
	case "timeFormat.hour":
		return "timeFormat_hour"
	case "timeFormat.minute":
		return "timeFormat_minute"
	case "timeFormat.second":
		return "timeFormat_second"
	}
	return field
}
