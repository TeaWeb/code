package agents

import (
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
)

// Agent定义
type AgentConfig struct {
	Id                  string       `yaml:"id" json:"id"`                                   // ID
	On                  bool         `yaml:"on" json:"on"`                                   // 是否启用
	Name                string       `yaml:"name" json:"name"`                               // 名称
	Host                string       `yaml:"host" json:"host"`                               // 主机地址
	Key                 string       `yaml:"key" json:"key"`                                 // 密钥
	AllowAll            bool         `yaml:"allowAll" json:"allowAll"`                       // 是否允许所有的IP
	Allow               []string     `yaml:"allow" json:"allow"`                             // 允许的IP地址
	Apps                []*AppConfig `yaml:"apps" json:"apps"`                               // Apps
	Version             uint         `yaml:"version" json:"version"`                         // 版本
	CheckDisconnections bool         `yaml:"checkDisconnections" json:"checkDisconnections"` // 是否检查离线
	CountDisconnections int          `yaml:"countDisconnections" json:"countDisconnections"` // 错误次数
	GroupIds            []string     `yaml:"groupIds" json:"groupIds"`                       // 分组IDs
	AutoUpdates         bool         `yaml:"autoUpdates" json:"autoUpdates"`                 // 是否开启自动更新
	AppsIsInitialized   bool         `yaml:"appsIsInitialized" json:"appsIsInitialized"`     // 是否已经初始化App

	NoticeSetting map[notices.NoticeLevel][]*notices.NoticeReceiver `yaml:"noticeSetting" json:"noticeSetting"`
}

// 获取新对象
func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		On:                  true,
		Id:                  stringutil.Rand(16),
		CheckDisconnections: true,
	}
}

// 本地Agent
var localAgentConfig *AgentConfig = nil

func LocalAgentConfig() *AgentConfig {
	if localAgentConfig == nil {
		localAgentConfig = &AgentConfig{
			On:       true,
			Id:       "local",
			Name:     "本地",
			Key:      stringutil.Rand(32),
			AllowAll: false,
			Allow:    []string{"127.0.0.1"},
		}
	}
	return localAgentConfig
}

// 从文件中获取对象
func NewAgentConfigFromFile(filename string) *AgentConfig {
	reader, err := files.NewReader(Tea.ConfigFile("agents/" + filename))
	if err != nil {
		return nil
	}
	defer reader.Close()
	agent := &AgentConfig{}
	err = reader.ReadYAML(agent)
	if err != nil {
		return nil
	}
	return agent
}

// 根据ID获取对象
func NewAgentConfigFromId(agentId string) *AgentConfig {
	if len(agentId) == 0 {
		return nil
	}
	agent := NewAgentConfigFromFile("agent." + agentId + ".conf")
	if agent != nil {
		if agent.Id == "local" && len(agent.Name) == 0 {
			agent.Name = "本地"
		}

		return agent
	}

	if agentId == "local" {
		return LocalAgentConfig()
	}

	return nil
}

// 判断是否为Local Agent
func (this *AgentConfig) IsLocal() bool {
	return this.Id == "local"
}

// 校验
func (this *AgentConfig) Validate() error {
	for _, a := range this.Apps {
		err := a.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// 文件名
func (this *AgentConfig) Filename() string {
	return "agent." + this.Id + ".conf"
}

// 保存
func (this *AgentConfig) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	defer func() {
		NotifyAgentsChange() // 标记列表改变
	}()

	dirFile := files.NewFile(Tea.ConfigFile("agents"))
	if !dirFile.Exists() {
		dirFile.Mkdir()
	}

	writer, err := files.NewWriter(Tea.ConfigFile("agents/" + this.Filename()))
	if err != nil {
		return err
	}
	defer writer.Close()
	this.Version++
	_, err = writer.WriteYAML(this)
	return err
}

// 删除
func (this *AgentConfig) Delete() error {
	defer func() {
		NotifyAgentsChange() // 标记列表改变
	}()

	// 删除board
	{
		f := files.NewFile(Tea.ConfigFile("agents/board." + this.Id + ".conf"))
		if f.Exists() {
			err := f.Delete()
			if err != nil {
				return err
			}
		}
	}

	f := files.NewFile(Tea.ConfigFile("agents/" + this.Filename()))
	return f.Delete()
}

// 添加App
func (this *AgentConfig) AddApp(app *AppConfig) {
	this.Apps = append(this.Apps, app)
}

// 替换App，如果不存在则增加
func (this *AgentConfig) ReplaceApp(app *AppConfig) {
	found := false
	for index, a := range this.Apps {
		if a.Id == app.Id {
			this.Apps[index] = app
			found = true
			break
		}
	}
	if !found {
		this.Apps = append(this.Apps, app)
	}
}

// 添加一组App
func (this *AgentConfig) AddApps(apps []*AppConfig) {
	this.Apps = append(this.Apps, apps ...)
}

// 删除App
func (this *AgentConfig) RemoveApp(appId string) {
	result := []*AppConfig{}
	for _, a := range this.Apps {
		if a.Id == appId {
			continue
		}
		result = append(result, a)
	}
	this.Apps = result
}

// 移动App位置
func (this *AgentConfig) MoveApp(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.Apps) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.Apps) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	location := this.Apps[fromIndex]
	newList := []*AppConfig{}
	for i := 0; i < len(this.Apps); i++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			newList = append(newList, location)
		}
		newList = append(newList, this.Apps[i])
		if fromIndex < toIndex && i == toIndex {
			newList = append(newList, location)
		}
	}

	this.Apps = newList
}

// 查找App
func (this *AgentConfig) FindApp(appId string) *AppConfig {
	for _, a := range this.Apps {
		if a.Id == appId {
			return a
		}
	}
	return nil
}

// 判断是否有某个App
func (this *AgentConfig) HasApp(appId string) bool {
	for _, a := range this.Apps {
		if a.Id == appId {
			return true
		}
	}
	return false
}

// YAML编码
func (this *AgentConfig) EncodeYAML() ([]byte, error) {
	return yaml.Marshal(this)
}

// 查找任务
func (this *AgentConfig) FindTask(taskId string) (appConfig *AppConfig, taskConfig *TaskConfig) {
	for _, app := range this.Apps {
		for _, task := range app.Tasks {
			if task.Id == taskId {
				return app, task
			}
		}
	}
	return nil, nil
}

// 查找监控项
func (this *AgentConfig) FindItem(itemId string) (appConfig *AppConfig, item *Item) {
	for _, app := range this.Apps {
		for _, item := range app.Items {
			if item.Id == itemId {
				return app, item
			}
		}
	}
	return nil, nil
}

// 添加分组
func (this *AgentConfig) AddGroup(groupId string) {
	if lists.ContainsString(this.GroupIds, groupId) {
		return
	}
	this.GroupIds = append(this.GroupIds, groupId)
}

// 删除分组
func (this *AgentConfig) RemoveGroup(groupId string) {
	result := []string{}
	for _, g := range this.GroupIds {
		if g == groupId {
			continue
		}
		result = append(result, g)
	}
	this.GroupIds = result
}

// 判断是否有某些分组
func (this *AgentConfig) InGroups(groupIds []string) bool {
	if len(this.GroupIds) == 0 && len(groupIds) == 0 {
		return true
	}
	for _, groupId := range groupIds {
		b := lists.ContainsString(this.GroupIds, groupId)
		if b {
			return true
		}
	}
	return false
}

// 添加内置的App
func (this *AgentConfig) AddDefaultApps() {
	this.AppsIsInitialized = true
	{
		app := NewAppConfig()
		app.Id = "system"
		app.Name = "系统"
		this.AddApp(app)

		board := NewAgentBoard(this.Id)

		// 添加到看板
		defer func() {
			board.Save()
		}()

		// cpu
		{
			// item
			item := NewItem()
			item.Id = "cpu.usage"
			item.Name = "CPU使用量（%）"
			item.Interval = "60s"

			source := NewCPUSource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)
			app.AddItem(item)

			// 阈值
			threshold1 := NewThreshold()
			threshold1.Param = "${usage.avg}"
			threshold1.Value = "80"
			threshold1.NoticeLevel = notices.NoticeLevelWarning
			threshold1.Operator = ThresholdOperatorGte
			item.AddThreshold(threshold1)

			// chart
			chart := widgets.NewChart()
			chart.Id = "cpu.chart1"
			chart.Name = "CPU使用量（%）"
			chart.Columns = 2
			chart.Type = "javascript"
			chart.Options = maps.Map{
				"code": `
var chart = new charts.LineChart();
chart.max = 100;

var query = new values.Query();
query.limit(30)
var ones = query.desc().cache(60).findAll();
ones.reverse();

var lines = [];

{
	var line = new charts.Line();
	line.color = colors.ARRAY[0];
	line.isFilled = true;
	line.values = [];
	lines.push(line);
}

ones.$each(function (k, v) {
	lines[0].values.push(v.value.usage.avg);
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLines(lines);
chart.render();
`,
			}
			item.AddChart(chart)
			board.AddChart(app.Id, item.Id, chart.Id)
		}

		// load
		{
			// item
			item := NewItem()
			item.Id = "cpu.load"
			item.Name = "负载（Load）"
			item.Interval = "60s"

			source := NewLoadSource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)

			app.AddItem(item)

			// 阈值
			{
				threshold1 := NewThreshold()
				threshold1.Param = "${load5}"
				threshold1.Value = "10"
				threshold1.NoticeLevel = notices.NoticeLevelWarning
				threshold1.Operator = ThresholdOperatorGte
				item.AddThreshold(threshold1)
			}

			{
				threshold2 := NewThreshold()
				threshold2.Param = "${load5}"
				threshold2.Value = "20"
				threshold2.NoticeLevel = notices.NoticeLevelError
				threshold2.Operator = ThresholdOperatorGte
				item.AddThreshold(threshold2)
			}

			// chart
			chart := widgets.NewChart()
			chart.Id = "cpu.load.chart1"
			chart.Name = "负载（Load）"
			chart.Columns = 2
			chart.Type = "javascript"
			chart.Options = maps.Map{
				"code": `
var chart = new charts.LineChart();

var query = new values.Query();
query.limit(30)
var ones = query.desc().cache(60).findAll();
ones.reverse();

var lines = [];

{
	var line = new charts.Line();
	line.name = "1分钟";
	line.color = colors.ARRAY[0];
	line.isFilled = true;
	line.values = [];
	lines.push(line);
}

{
	var line = new charts.Line();
	line.name = "5分钟";
	line.color = colors.BROWN;
	line.isFilled = false;
	line.values = [];
	lines.push(line);
}

{
	var line = new charts.Line();
	line.name = "15分钟";
	line.color = colors.RED;
	line.isFilled = false;
	line.values = [];
	lines.push(line);
}

var maxValue = 1;

ones.$each(function (k, v) {
	lines[0].values.push(v.value.load1);
	lines[1].values.push(v.value.load5);
	lines[2].values.push(v.value.load15);

	if (v.value.load1 > maxValue) {
		maxValue = Math.ceil(v.value.load1 / 2) * 2;
	}
	if (v.value.load5 > maxValue) {
		maxValue = Math.ceil(v.value.load5 / 2) * 2;
	}
	if (v.value.load15 > maxValue) {
		maxValue = Math.ceil(v.value.load15 / 2) * 2;
	}
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLines(lines);
chart.max = maxValue;
chart.render();
`,
			}
			item.AddChart(chart)
			board.AddChart(app.Id, item.Id, chart.Id)
		}

		// memory usage
		{
			//item
			item := NewItem()
			item.Id = "memory.usage"
			item.Name = "内存使用量"
			item.Interval = "60s"

			source := NewMemorySource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)

			app.AddItem(item)

			// 阈值
			{
				threshold1 := NewThreshold()
				threshold1.Param = "${usage.virtualPercent}"
				threshold1.Value = "80"
				threshold1.NoticeLevel = notices.NoticeLevelWarning
				threshold1.Operator = ThresholdOperatorGte
				item.AddThreshold(threshold1)
			}

			// chart
			{
				chart := widgets.NewChart()
				chart.Id = "memory.usage.chart1"
				chart.Name = "内存使用量（%）"
				chart.Columns = 2
				chart.Type = "javascript"
				chart.Options = maps.Map{
					"code": `
var chart = new charts.LineChart();

var query = new values.Query();
query.limit(30)
var ones = query.desc().cache(60).findAll();
ones.reverse();

var lines = [];

{
	var line = new charts.Line();
	line.color = colors.ARRAY[0];
	line.isFilled = true;
	line.values = [];
	lines.push(line);
}

ones.$each(function (k, v) {
	lines[0].values.push(v.value.usage.virtualPercent);

	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLines(lines);
chart.max = 100;
chart.render();
`,
				}
				item.AddChart(chart)
				board.AddChart(app.Id, item.Id, chart.Id)
			}

			{
				chart := widgets.NewChart()
				chart.Id = "memory.usage.chart2"
				chart.Name = "当前内存使用量"
				chart.Columns = 1
				chart.Type = "javascript"
				chart.Options = maps.Map{
					"code": `
var chart = new charts.StackBarChart();

var latest = new values.Query().latest(1);
var hasWarning = false;
if (latest.length > 0) {
	hasWarning = (latest[0].value.usage.swapPercent > 50) || (latest[0].value.usage.virtualPercent > 80);
	chart.values = [ 
		[latest[0].value.usage.swapUsed, latest[0].value.usage.swapTotal - latest[0].value.usage.swapUsed],
		[latest[0].value.usage.virtualUsed, latest[0].value.usage.virtualTotal - latest[0].value.usage.virtualUsed]
	];
	chart.labels = [ "虚拟内存（" +  (Math.round(latest[0].value.usage.swapUsed * 10) / 10) + "G/" + Math.round(latest[0].value.usage.swapTotal) + "G"  + "）", "物理内存（" + (Math.round(latest[0].value.usage.virtualUsed * 10) / 10)+ "G/" + Math.round(latest[0].value.usage.virtualTotal)  + "G"  + "）"];
} else {
	chart.values = [ [0, 0], [0, 0] ];
	chart.labels = [ "虚拟内存", "物理内存" ];
}
if (hasWarning) {
	chart.colors = [ colors.RED, colors.GREEN ];
} else {
	chart.colors = [ colors.BROWN, colors.GREEN ];
}
chart.render();
`,
				}
				item.AddChart(chart)
				board.AddChart(app.Id, item.Id, chart.Id)
			}
		}

		// clock
		{
			// item
			item := NewItem()
			item.Id = "clock"
			item.Name = "时钟"
			item.Interval = "60s"

			source := NewDateSource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)

			app.AddItem(item)

			// 时钟
			{
				chart := widgets.NewChart()
				chart.Id = "clock"
				chart.Name = "时钟"
				chart.Columns = 1
				chart.Type = "javascript"
				chart.Options = maps.Map{
					"code": `
var chart = new charts.Clock();
var latest = new values.Query().latest(1);
if (latest.length > 0) {
	chart.timestamp = parseInt(new Date().getTime() / 1000) - (latest[0].createdAt - latest[0].value.timestamp);
}
chart.render();
`,
				}
				item.AddChart(chart)
				board.AddChart(app.Id, item.Id, chart.Id)
			}
		}

		// network out && network in
		{
			// item
			item := NewItem()
			item.Id = "network.usage"
			item.Name = "网络相关"
			item.Interval = "60s"

			{
				threshold := NewThreshold()
				threshold.Param = "${stat.avgSentBytes}"
				threshold.Operator = ThresholdOperatorGte
				threshold.Value = "13107200"
				threshold.NoticeLevel = notices.NoticeLevelWarning
				threshold.NoticeMessage = "当前出口流量超过100MBit/s"
				item.AddThreshold(threshold)
			}

			source := NewNetworkSource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)

			app.AddItem(item)

			// 图表
			{
				chart := widgets.NewChart()
				chart.Id = "network.usage.received"
				chart.Name = "出口带宽（M/s）"
				chart.Columns = 2
				chart.Type = "javascript"
				chart.Options = maps.Map{
					"code": `
var chart = new charts.LineChart();

var line = new charts.Line();
line.isFilled = true;

var ones = new values.Query().cache(60).latest(60);
ones.reverse();
ones.$each(function (k, v) {
	line.values.push(Math.round(v.value.stat.avgSentBytes / 1024 / 1024 * 100) / 100);
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});
var maxValue = line.values.$max();
if (maxValue < 1) {
	chart.max = 1;
} else if (maxValue < 5) {
	chart.max = 5;
} else if (maxValue < 10) {
	chart.max = 10;
}

chart.addLine(line);
chart.render();
`,
				}
				item.AddChart(chart)
				board.AddChart(app.Id, item.Id, chart.Id)
			}

			{
				chart := widgets.NewChart()
				chart.Id = "network.usage.sent"
				chart.Name = "入口带宽（M/s）"
				chart.Columns = 2
				chart.Type = "javascript"
				chart.Options = maps.Map{
					"code": `
var chart = new charts.LineChart();

var line = new charts.Line();
line.isFilled = true;

var ones = new values.Query().cache(60).latest(60);
ones.reverse();
ones.$each(function (k, v) {
	line.values.push(Math.round(v.value.stat.avgReceivedBytes / 1024 / 1024 * 100) / 100);
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});
var maxValue = line.values.$max();
if (maxValue < 1) {
	chart.max = 1;
} else if (maxValue < 5) {
	chart.max = 5;
} else if (maxValue < 10) {
	chart.max = 10;
}

chart.addLine(line);
chart.render();
`,
				}
				item.AddChart(chart)
				board.AddChart(app.Id, item.Id, chart.Id)
			}
		}

		// disk
		{
			// item
			item := NewItem()
			item.Id = "disk.usage"
			item.Name = "文件系统"
			item.Interval = "120s"

			source := NewDiskSource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)

			{
				threshold := NewThreshold()
				threshold.Param = "${partitions.$.percent}"
				threshold.Operator = ThresholdOperatorGt
				threshold.Value = "80"
				threshold.NoticeLevel = notices.NoticeLevelWarning
				threshold.NoticeMessage = "${ROW.name}分区已使用80%"
				item.AddThreshold(threshold)
			}

			app.AddItem(item)

			// 图表
			{
				chart := widgets.NewChart()
				chart.Id = "disk.usage.chart1"
				chart.Name = "文件系统"
				chart.Columns = 2
				chart.Type = "javascript"
				chart.Options = maps.Map{
					"code": `
var chart = new charts.StackBarChart();
chart.values = [];
chart.labels = [];

var latest = new values.Query().cache(120).latest(1);
if (latest.length > 0) {
	var partitions = latest[0].value.partitions;
	partitions.$each(function (k, v) {
		chart.values.push([v.used, v.total - v.used]);
		chart.labels.push(v.name + "（" + (Math.round(v.used / 1024 / 1024 / 1024 * 100) / 100)+ "G/" + (Math.round(v.total / 1024 / 1024 / 1024 * 100) / 100) +"G）");
	});

	chart.options.height = partitions.length * 4;
}

chart.colors = [ colors.BROWN, colors.GREEN ];
chart.render();
`,
				}
				item.AddChart(chart)
				board.AddChart(app.Id, item.Id, chart.Id)
			}
		}
	}
}

// 添加通知接收者
func (this *AgentConfig) AddNoticeReceiver(level notices.NoticeLevel, receiver *notices.NoticeReceiver) {
	if this.NoticeSetting == nil {
		this.NoticeSetting = map[notices.NoticeLevel][]*notices.NoticeReceiver{}
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		receivers = []*notices.NoticeReceiver{}
	}
	receivers = append(receivers, receiver)
	this.NoticeSetting[level] = receivers
}

// 删除通知接收者
func (this *AgentConfig) RemoveNoticeReceiver(level notices.NoticeLevel, receiverId string) {
	if this.NoticeSetting == nil {
		return
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		return
	}

	result := []*notices.NoticeReceiver{}
	for _, r := range receivers {
		if r.Id == receiverId {
			continue
		}
		result = append(result, r)
	}
	this.NoticeSetting[level] = result
}

// 获取通知接收者数量
func (this *AgentConfig) CountNoticeReceivers() int {
	count := 0
	for _, receivers := range this.NoticeSetting {
		count += len(receivers)
	}
	return count
}

// 删除媒介
func (this *AgentConfig) RemoveMedia(mediaId string) (found bool) {
	for level, receivers := range this.NoticeSetting {
		result := []*notices.NoticeReceiver{}
		for _, receiver := range receivers {
			if receiver.MediaId == mediaId {
				found = true
				continue
			}
			result = append(result, receiver)
		}
		this.NoticeSetting[level] = result
	}
	return
}

// 查找一个或多个级别对应的接收者，并合并相同的接收者
func (this *AgentConfig) FindAllNoticeReceivers(level ...notices.NoticeLevel) []*notices.NoticeReceiver {
	if len(level) == 0 {
		return []*notices.NoticeReceiver{}
	}

	m := maps.Map{} // mediaId_user => bool
	result := []*notices.NoticeReceiver{}
	for _, l := range level {
		receivers, ok := this.NoticeSetting[l]
		if !ok {
			continue
		}
		for _, receiver := range receivers {
			if !receiver.On {
				continue
			}
			key := receiver.Key()
			if m.Has(key) {
				continue
			}
			m[key] = true
			result = append(result, receiver)
		}
	}
	return result
}

// 获取分组名
func (this *AgentConfig) GroupName() string {
	if len(this.GroupIds) == 0 {
		return "默认分组"
	}
	groupId := this.GroupIds[0]
	if len(groupId) == 0 {
		return "默认分组"
	}

	group := SharedGroupConfig().FindGroup(groupId)
	if group == nil {
		return "默认分组"
	}
	return group.Name
}
