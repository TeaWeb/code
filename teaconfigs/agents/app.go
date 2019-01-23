package agents

import "github.com/iwind/TeaGo/utils/string"

// App定义
type AppConfig struct {
	Id       string        `yaml:"id" json:"id"`             // ID
	On       bool          `yaml:"on" json:"on"`             // 是否启用
	Tasks    []*TaskConfig `yaml:"tasks" json:"tasks"`       // 任务设置
	Items    []*Item       `yaml:"item" json:"items"`        // 监控项
	Name     string        `yaml:"name" json:"name"`         // 名称
	IsSystem bool          `yaml:"isSystem" json:"isSystem"` // 是否为系统定义
}

// 获取新对象
func NewAppConfig() *AppConfig {
	return &AppConfig{
		Id: stringutil.Rand(16),
		On: true,
	}
}

// 获取非用户定义对象
func NewSystemAppConfig(id string) *AppConfig {
	app := NewAppConfig()
	app.IsSystem = true
	app.Id = id
	return app
}

// 校验
func (this *AppConfig) Validate() error {
	// 任务
	for _, t := range this.Tasks {
		err := t.Validate()
		if err != nil {
			return err
		}
	}

	// 监控项
	for _, item := range this.Items {
		err := item.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Schedule Tasks
func (this *AppConfig) FindSchedulingTasks() []*TaskConfig {
	result := []*TaskConfig{}
	for _, t := range this.Tasks {
		if len(t.Schedule) > 0 {
			result = append(result, t)
		}
	}
	return result
}

// Boot Tasks
func (this *AppConfig) FindBootingTasks() []*TaskConfig {
	result := []*TaskConfig{}
	for _, t := range this.Tasks {
		if t.IsBooting {
			result = append(result, t)
		}
	}
	return result
}

// Manual Tasks
func (this *AppConfig) FindManualTasks() []*TaskConfig {
	result := []*TaskConfig{}
	for _, t := range this.Tasks {
		if t.IsManual {
			result = append(result, t)
		}
	}
	return result
}

// 添加任务
func (this *AppConfig) AddTask(task *TaskConfig) {
	this.Tasks = append(this.Tasks, task)
}

// 删除任务
func (this *AppConfig) RemoveTask(taskId string) {
	result := []*TaskConfig{}
	for _, t := range this.Tasks {
		if t.Id == taskId {
			continue
		}
		result = append(result, t)
	}
	this.Tasks = result
}

// 查找任务
func (this *AppConfig) FindTask(taskId string) *TaskConfig {
	for _, t := range this.Tasks {
		if t.Id == taskId {
			return t
		}
	}
	return nil
}

// 添加监控项
func (this *AppConfig) AddItem(item *Item) {
	this.Items = append(this.Items, item)
}

// 删除监控项
func (this *AppConfig) RemoveItem(itemId string) {
	result := []*Item{}
	for _, item := range this.Items {
		if item.Id == itemId {
			continue
		}
		result = append(result, item)
	}
	this.Items = result
}

// 查找监控项
func (this *AppConfig) FindItem(itemId string) *Item {
	for _, item := range this.Items {
		if item.Id == itemId {
			item.Validate()
			return item
		}
	}
	return nil
}
