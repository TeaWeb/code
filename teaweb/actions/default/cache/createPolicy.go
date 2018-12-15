package cache

import (
	"fmt"
	"github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
)

type CreatePolicyAction actions.Action

// 缓存缓存策略
func (this *CreatePolicyAction) Run(params struct{}) {
	this.Data["types"] = teacache.AllCacheTypes()

	this.Show()
}

func (this *CreatePolicyAction) RunPost(params struct {
	Name string
	Key  string
	Type string

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
	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称").

		Field("key", params.Key).
		Require("请输入缓存Key")

	policy := teaconfigs.NewCachePolicy()
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

	config, _ := teaconfigs.SharedCacheConfig()
	config.AddPolicy(policy.Filename)
	err = config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Next("/cache", nil)
	this.Success("保存成功")
}
