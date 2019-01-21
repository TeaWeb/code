package agents

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
)

// Agent列表
type AgentList struct {
	Files []string `yaml:"files" json:"files"`
}

// 取得Agent列表
func SharedAgentList() (*AgentList, error) {
	file := files.NewFile(Tea.ConfigFile("agents/agentlist.conf"))
	if !file.Exists() {
		// 创建目录
		dir := files.NewFile(Tea.ConfigFile("agents"))
		if !dir.Exists() {
			dir.MkdirAll()
		}

		return &AgentList{}, nil
	}
	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	agentList := &AgentList{}
	err = reader.ReadYAML(agentList)
	if err != nil {
		return nil, err
	}
	return agentList, nil
}

// 添加Agent
func (this *AgentList) AddAgent(agentFile string) {
	this.Files = append(this.Files, agentFile)
}

// 删除Agent
func (this *AgentList) RemoveAgent(agentFile string) {
	result := []string{}
	for _, f := range this.Files {
		if f == agentFile {
			continue
		}
		result = append(result, f)
	}
	this.Files = result
}

// 查找所有Agents
func (this *AgentList) FindAllAgents() []*AgentConfig {
	result := []*AgentConfig{}
	for _, f := range this.Files {
		agent := NewAgentConfigFromFile(f)
		if agent == nil {
			continue
		}
		result = append(result, agent)
	}
	return result
}

// 保存
func (this *AgentList) Save() error {
	writer, err := files.NewWriter(Tea.ConfigFile("agents/agentlist.conf"))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}
