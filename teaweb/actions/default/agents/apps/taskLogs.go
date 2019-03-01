package apps

import (
	"context"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TaskLogsAction actions.Action

// 任务日志
func (this *TaskLogsAction) Run(params struct {
	AgentId string
	AppId   string
	TaskId  string
	Tabbar  string
}) {
	this.Data["tabbar"] = params.Tabbar

	agentutils.InitAppData(this, params.AgentId, params.AppId, params.Tabbar)

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	task := app.FindTask(params.TaskId)
	if task == nil {
		this.Fail("找不到要修改的任务")
	}

	this.Data["task"] = maps.Map{
		"id":        task.Id,
		"name":      task.Name,
		"on":        task.On,
		"script":    task.Script,
		"cwd":       task.Cwd,
		"isBooting": task.IsBooting,
		"isManual":  task.IsManual,
		"env":       task.Env,
		"schedules": lists.Map(task.Schedule, func(k int, v interface{}) interface{} {
			s := v.(*agents.ScheduleConfig)
			return maps.Map{
				"summary": s.Summary(),
			}
		}),
	}

	this.Show()
}

// 日志数据
func (this *TaskLogsAction) RunPost(params struct {
	AgentId string
	TaskId  string
	LastId  string
}) {
	filter := map[string]interface{}{
		"taskId": params.TaskId,
	}
	if len(params.LastId) > 0 {
		lastObjectId, err := primitive.ObjectIDFromHex(params.LastId)
		if err != nil {
			logs.Error(err)
		} else {
			filter["_id"] = map[string]interface{}{
				"$gt": lastObjectId,
			}
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	cursor, err := teamongo.FindCollection("logs.agent." + params.AgentId).Find(ctx, filter, options.Find().SetSort(map[string]interface{}{
		"_id": -1,
	}), options.Find().SetLimit(100))
	if err != nil {
		this.Fail("查询数据库出错：" + err.Error())
	}
	taskLogs := []*agentutils.ProcessLog{}
	for cursor.Next(context.Background()) {
		m := &agentutils.ProcessLog{}
		err = cursor.Decode(&m)
		if err != nil {
			logs.Error(err)
		} else {
			taskLogs = append(taskLogs, m)
		}
	}
	err = cursor.Close(context.Background())
	if err != nil {
		logs.Error(err)
	}

	this.Data["logs"] = taskLogs
	this.Success()
}
