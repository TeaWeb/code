package headers

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type DataAction actions.Action

// Header数据
func (this *DataAction) Run(params struct {
	ServerId   string
	LocationId string
	RewriteId  string
	FastcgiId  string
	BackendId  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	headerList, err := server.FindHeaderList(params.LocationId, params.BackendId, params.RewriteId, params.FastcgiId)
	if err != nil {
		this.Fail(err.Error())
	}
	this.Data["headers"] = lists.Map(headerList.AllHeaders(), func(k int, v interface{}) interface{} {
		header := v.(*shared.HeaderConfig)
		return maps.Map{
			"on":     header.On,
			"id":     header.Id,
			"always": header.Always,
			"status": header.Status,
			"name":   header.Name,
			"value":  header.Value,
		}
	})
	this.Data["ignoreHeaders"] = headerList.AllIgnoreHeaders()

	this.Success()
}
