package teastats

import (
	"fmt"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/types"
)

// 计算平均值操作
type AvgOperation struct {
	coll       *teamongo.Collection
	filter     map[string]interface{}
	init       map[string]interface{}
	countField string
	count      int
	sumField   string
	sum        float64
}

func (this *AvgOperation) uniqueId() string {
	keys := []string{}
	for key := range this.filter {
		keys = append(keys, key)
	}
	lists.Sort(keys, func(i int, j int) bool {
		return keys[i] < keys[j]
	})

	uniqueId := fmt.Sprintf("%p", this.coll)
	for _, key := range keys {
		uniqueId += "@" + types.String(this.filter[key])
	}
	return uniqueId + "@" + this.countField + "@" + this.sumField
}
