package mongo

import (
	"context"
	"fmt"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type CollStatAction actions.Action

// 集合统计
func (this *CollStatAction) Run(params struct {
	CollNames []string
}) {
	db := teamongo.SharedClient().Database(teamongo.DatabaseName)

	results := maps.Map{}
	for _, collName := range params.CollNames {
		if len(collName) == 0 {
			continue
		}

		result := db.RunCommand(context.Background(), bsonx.Doc{{"collStats", bsonx.String(collName)}, {"verbose", bsonx.Boolean(false)}})
		if result.Err() != nil {
			this.Fail("读取统计信息失败：" + result.Err().Error())
		}

		m1 := maps.Map{}
		err := result.Decode(&m1)
		if err != nil {
			if result.Err() != nil {
				this.Fail("读取统计信息失败：" + result.Err().Error())
			}
		}
		if m1.GetInt("ok") != 1 {
			continue
		}
		size := float64(m1.GetInt("size"))
		formattedSize := ""
		if size < 1024 {
			formattedSize = fmt.Sprintf("%.2fB", size)
		} else if size < 1024*1024 {
			formattedSize = fmt.Sprintf("%.2fKB", size/1024)
		} else if size < 1024*1024*1024 {
			formattedSize = fmt.Sprintf("%.2fMB", size/1024/1024)
		} else {
			formattedSize = fmt.Sprintf("%.2fGB", size/1024/1024/1024)
		}
		results[collName] = maps.Map{
			"count":         m1.GetInt("count"),
			"size":          size,
			"formattedSize": formattedSize,
		}
	}

	this.Data["result"] = results

	this.Success()
}
