package agents

import (
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/utils/string"
)

// Agent定义
type AgentConfig struct {
	Id       string       `yaml:"id" json:"id"`             // ID
	On       bool         `yaml:"on" json:"on"`             // 是否启用
	Name     string       `yaml:"name" json:"name"`         // 名称
	Host     string       `yaml:"host" json:"host"`         // 主机地址
	Key      string       `yaml:"key" json:"key"`           // 密钥
	AllowAll bool         `yaml:"allowAll" json:"allowAll"` // 是否允许所有的IP
	Allow    []string     `yaml:"allow" json:"allow"`       // 允许的IP地址
	Apps     []*AppConfig `yaml:"apps" json:"apps"`         // Apps
	Version  uint         `yaml:"version" json:"version"`   // 版本
}

// 获取新对象
func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		On: true,
		Id: stringutil.Rand(16),
	}
}

// 本地Agent
var localAgentConfig *AgentConfig = nil

func LocalAgentConfig() *AgentConfig {
	if localAgentConfig == nil {
		localAgentConfig = &AgentConfig{
			On:       true,
			Id:       "local",
			Name:     "本地",
			Key:      stringutil.Rand(32),
			AllowAll: false,
			Allow:    []string{"127.0.0.1"},
		}
	}
	return localAgentConfig
}

// 从文件中获取对象
func NewAgentConfigFromFile(filename string) *AgentConfig {
	reader, err := files.NewReader(Tea.ConfigFile("agents/" + filename))
	if err != nil {
		return nil
	}
	defer reader.Close()
	agent := &AgentConfig{}
	err = reader.ReadYAML(agent)
	if err != nil {
		return nil
	}
	return agent
}

// 根据ID获取对象
func NewAgentConfigFromId(agentId string) *AgentConfig {
	if len(agentId) == 0 {
		return nil
	}
	agent := NewAgentConfigFromFile("agent." + agentId + ".conf")
	if agent != nil {
		if agent.Id == "local" && len(agent.Name) == 0 {
			agent.Name = "本地"
		}

		return agent
	}

	if agentId == "local" {
		return LocalAgentConfig()
	}

	return nil
}

// 判断是否为Local Agent
func (this *AgentConfig) IsLocal() bool {
	return this.Id == "local"
}

// 校验
func (this *AgentConfig) Validate() error {
	for _, a := range this.Apps {
		err := a.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// 文件名
func (this *AgentConfig) Filename() string {
	return "agent." + this.Id + ".conf"
}

// 保存
func (this *AgentConfig) Save() error {
	dirFile := files.NewFile(Tea.ConfigFile("agents"))
	if !dirFile.Exists() {
		dirFile.Mkdir()
	}

	writer, err := files.NewWriter(Tea.ConfigFile("agents/" + this.Filename()))
	if err != nil {
		return err
	}
	defer writer.Close()
	this.Version ++
	_, err = writer.WriteYAML(this)
	return err
}

// 删除
func (this *AgentConfig) Delete() error {
	f := files.NewFile(Tea.ConfigFile("agents/" + this.Filename()))
	return f.Delete()
}

// 添加App
func (this *AgentConfig) AddApp(app *AppConfig) {
	this.Apps = append(this.Apps, app)
}

// 添加一组App
func (this *AgentConfig) AddApps(apps []*AppConfig) {
	this.Apps = append(this.Apps, apps ...)
}

// 删除App
func (this *AgentConfig) RemoveApp(appId string) {
	result := []*AppConfig{}
	for _, a := range this.Apps {
		if a.Id == appId {
			continue
		}
		result = append(result, a)
	}
	this.Apps = result
}

// 查找App
func (this *AgentConfig) FindApp(appId string) *AppConfig {
	for _, a := range this.Apps {
		if a.Id == appId {
			return a
		}
	}
	return nil
}

// YAML编码
func (this *AgentConfig) EncodeYAML() ([]byte, error) {
	return yaml.Marshal(this)
}

// 查找任务
func (this *AgentConfig) FindTask(taskId string) (appConfig *AppConfig, taskConfig *TaskConfig) {
	for _, app := range this.Apps {
		for _, task := range app.Tasks {
			if task.Id == taskId {
				return app, task
			}
		}
	}
	return nil, nil
}

// 查找监控项
func (this *AgentConfig) FindItem(itemId string) (appConfig *AppConfig, item *Item) {
	for _, app := range this.Apps {
		for _, item := range app.Items {
			if item.Id == itemId {
				return app, item
			}
		}
	}
	return nil, nil
}

// 清除系统App
func (this *AgentConfig) ResetSystemApps() {
	result := []*AppConfig{}
	for _, app := range this.Apps {
		if app.IsSystem {
			continue
		}
		result = append(result, app)
	}
	this.Apps = result
}

// 取得系统App列表
func (this *AgentConfig) FindSystemApps() []*AppConfig {
	result := []*AppConfig{}
	for _, app := range this.Apps {
		if !app.IsSystem {
			continue
		}
		result = append(result, app)
	}
	return result
}
