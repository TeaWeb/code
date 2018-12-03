package teastats

import (
	"fmt"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/types"
)

// 为某个字段加一操作
type IncrementOperation struct {
	coll   *teamongo.Collection
	filter map[string]interface{}
	init   map[string]interface{}
	field  string
}

func (this *IncrementOperation) uniqueId() string {
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
	return uniqueId + "@" + this.field
}
