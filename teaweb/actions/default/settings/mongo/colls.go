package mongo

import (
	"context"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"sort"
	"time"
)

type CollsAction actions.Action

// 集合列表
func (this *CollsAction) Run(params struct{}) {
	db := teamongo.SharedClient().Database(teamongo.DatabaseName)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := db.ListCollections(ctx, maps.Map{})
	if err != nil {
		logs.Error(err)
		this.Fail("读取集合列表失败：" + err.Error())
	}
	defer cursor.Close(context.Background())

	names := []string{}
	for cursor.Next(context.Background()) {
		m := maps.Map{}
		err := cursor.Decode(&m)
		if err != nil {
			logs.Error(err)
			this.Fail("读取集合列表失败：" + err.Error())
		}
		name := m.GetString("name")
		if len(name) == 0 {
			continue
		}
		names = append(names, name)
	}

	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	// 排序
	result := []maps.Map{}
	recognizedNames := []string{}

	// 日志
	{
		reg := regexp.MustCompile("^logs\\.\\d{8}$")
		for _, name := range names {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "代理访问日志",
				"canDelete": true,
				"subName":   name[5:9] + "-" + name[9:11] + "-" + name[11:],
			})
		}
	}

	// 统计
	{
		reg := regexp.MustCompile("^values\\.server\\.")
		for _, name := range names {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "代理统计数据",
				"canDelete": true,
			})
		}
	}

	// 监控数据
	{
		reg := regexp.MustCompile("^values\\.agent\\.")
		for _, name := range names {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "主机监控数据",
				"canDelete": true,
			})
		}
	}

	// 监控数据
	{
		reg := regexp.MustCompile("^logs\\.agent\\.")
		for _, name := range names {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "主机任务运行日志",
				"canDelete": true,
			})
		}
	}

	// 通知
	{
		reg := regexp.MustCompile("^notices$")
		for _, name := range names {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "通知提醒",
				"canDelete": true,
			})
		}
	}

	// 审计日志
	{
		reg := regexp.MustCompile("^logs\\.audit$")
		for _, name := range names {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "审计日志（操作日志）",
				"canDelete": true,
			})
		}
	}

	// 旧的统计数据
	{
		reg := regexp.MustCompile("^stats\\.")
		for _, name := range names {
			if !reg.MatchString(name) {
				continue
			}
			recognizedNames = append(recognizedNames, name)
			result = append(result, maps.Map{
				"name":      name,
				"type":      "旧的统计数据",
				"canDelete": true,
				"warning":   true,
			})
		}
	}

	// 其他
	for _, name := range names {
		if lists.ContainsString(recognizedNames, name) {
			continue
		}
		result = append(result, maps.Map{
			"name":      name,
			"type":      "无法识别，请报告官方",
			"canDelete": false,
			"warning":   true,
		})
	}

	this.Data["colls"] = result

	this.Success()
}
