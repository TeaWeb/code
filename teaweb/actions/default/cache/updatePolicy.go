package cache

import (
	"fmt"
	"github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type UpdatePolicyAction actions.Action

// 修改缓存策略
func (this *UpdatePolicyAction) Run(params struct {
	Filename string
}) {
	policy := teaconfigs.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到要修改的缓存策略")
	}

	this.Data["types"] = teacache.AllCacheTypes()

	policy.Validate()

	this.Data["policy"] = maps.Map{
		"filename":    policy.Filename,
		"name":        policy.Name,
		"key":         policy.Key,
		"type":        policy.Type,
		"options":     policy.Options,
		"hasAdvanced": policy.CapacitySize() > 0 || (policy.Life != "72h" && policy.LifeDuration() > 0) || (len(policy.Status) != 1 || !lists.Contains(policy.Status, 200)) || policy.MaxDataSize() > 0,
		"life":        policy.Life,
		"status":      policy.Status,
		"maxSize":     policy.MaxSize,
		"capacity":    policy.Capacity,
	}

	this.Show()
}

func (this *UpdatePolicyAction) RunPost(params struct {
	Filename string
	Name     string
	Key      string
	Type     string

	IsAdvanced   bool
	Capacity     float64
	CapacityUnit string
	Life         int
	LifeUnit     string
	StatusList   []int
	MaxSize      float64
	MaxSizeUnit  string

	Must *actions.Must
}) {
	policy := teaconfigs.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到要修改的缓存策略")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称").

		Field("key", params.Key).
		Require("请输入缓存Key")

	policy.Name = params.Name
	policy.Key = params.Key
	policy.Type = params.Type

	if params.IsAdvanced {
		policy.Capacity = fmt.Sprintf("%.2f%s", params.Capacity, params.CapacityUnit)
		policy.Life = fmt.Sprintf("%d%s", params.Life, params.LifeUnit)
		for _, status := range params.StatusList {
			i := types.Int(status)
			if i >= 0 {
				policy.Status = append(policy.Status, i)
			}
		}
		policy.MaxSize = fmt.Sprintf("%.2f%s", params.MaxSize, params.MaxSizeUnit)
		policy.Status = params.StatusList
	} else {
		policy.Capacity = "0.00g"
		policy.Life = "72h"
		policy.MaxSize = "0.00m"
		policy.Status = []int{200}
	}

	// 选项
	if policy.Type == "file" {
		policy.Options = map[string]interface{}{
			"dir": this.ParamString("options_dir"),
		}
	} else if policy.Type == "redis" {
		policy.Options = map[string]interface{}{
			"network":  this.ParamString("options_network"),
			"host":     this.ParamString("options_host"),
			"port":     this.ParamString("options_port"),
			"password": this.ParamString("options_password"),
			"sock":     this.ParamString("options_sock"),
		}
	}

	err := policy.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Next("/cache", nil)
	this.Success("保存成功")
}
