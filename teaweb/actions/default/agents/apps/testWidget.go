package apps

import (
	"encoding/json"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"time"
)

type TestWidgetAction actions.Action

// 测试Widget
func (this *TestWidgetAction) Run(params struct {
	AgentId      string
	AppId        string
	Name         string
	Code         string
	Author       string
	Version      string
	Description  string
	ParamNames   []string
	ParamValues  []string
	ValueYears   []string
	ValueMonths  []string
	ValueDays    []string
	ValueHours   []string
	ValueMinutes []string
	ValueSeconds []string
	Values       []string
	ChartCode    string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	// 删除已经存在的测试数据
	err := teamongo.NewValueQuery().
		Agent(params.AgentId).
		App(params.AppId).
		Attr("isTesting", true).
		Delete()
	if err != nil {
		this.Fail("删除测试数据失败：" + err.Error())
	}

	// 写入数据
	for k := len(params.ValueYears) - 1; k >= 0; k -- {
		if k < len(params.ValueMonths) && k < len(params.ValueDays) && k < len(params.ValueHours) && k < len(params.ValueMinutes) && k < len(params.ValueSeconds) && k < len(params.Values) {
			year := types.Int(params.ValueYears[k])
			month := types.Int(params.ValueMonths[k])
			day := types.Int(params.ValueDays[k])
			hour := types.Int(params.ValueHours[k])
			minute := types.Int(params.ValueMinutes[k])
			second := types.Int(params.ValueSeconds[k])

			if year <= 0 {
				year = time.Now().Year()
			}

			if month < 1 {
				month = 1
			} else if month > 12 {
				month = 12
			}

			if day < 1 {
				day = 1
			} else if day > 31 {
				day = 31
			}

			if hour < 0 {
				hour = 0
			} else if hour > 23 {
				hour = 23
			}

			if minute < 0 {
				minute = 0
			} else if minute > 59 {
				minute = 59
			}

			if second < 0 {
				second = 0
			} else if second > 59 {
				second = 59
			}

			valueJSON := params.Values[k]
			var value interface{} = nil
			err := json.Unmarshal([]byte(valueJSON), &value)
			if err != nil {
				logs.Error(err)
				continue
			}

			testingValue := &agents.Value{
				AgentId:   params.AgentId,
				AppId:     params.AppId,
				Value:     value,
				IsTesting: true,
			}

			t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
			testingValue.SetTime(t)

			err = teamongo.NewValueQuery().Agent(params.AgentId).Insert(testingValue)
			if err != nil {
				this.Fail("测试数据添加失败：" + err.Error())
			}
		}
	}

	engine := scripts.NewEngine()
	engine.SetContext(&scripts.Context{
		Agent: agent,
		App:   app,
	})

	options := map[string]string{}
	for index, name := range params.ParamNames {
		if index < len(params.ParamValues) {
			options[name] = params.ParamValues[index]
		}
	}

	widgetCode := `var widget = new widgets.Widget({
	"name": ` + stringutil.JSONEncode(params.Name) + `,
	"code": ` + stringutil.JSONEncode(params.Code) + `,
	"author": ` + stringutil.JSONEncode(params.Author) + `,
	"version": ` + stringutil.JSONEncode(params.Version) + `,
	"description": ` + stringutil.JSONEncode(params.Description) + `,
	"options": ` + stringutil.JSONEncode(options) + `
});

widget.run = function () {
` + params.ChartCode + `	
};`

	err = engine.RunCode(widgetCode)
	if err != nil {
		this.Fail("测试失败：" + err.Error() + "\n~~~\n" + widgetCode)
	}

	this.Data["charts"] = engine.Charts()
	this.Data["output"] = engine.Output()

	this.Success()
}
