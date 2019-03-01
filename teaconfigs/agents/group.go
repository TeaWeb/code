package agents

import (
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/utils/string"
)

// 默认的config
var sharedGroupConfig *GroupConfig = nil

// Agent分组
type Group struct {
	Id   string `yaml:"id" json:"id"`
	On   bool   `yaml:"on" json:"on"`
	Name string `yaml:"name" json:"name"`
}

// 获取新分组
func NewGroup(name string) *Group {
	return &Group{
		Id:   stringutil.Rand(16),
		On:   true,
		Name: name,
	}
}

// 分组配置
type GroupConfig struct {
	Filename string   `yaml:"filename" json:"filename"`
	Groups   []*Group `yaml:"groups" json:"groups"`
}

// 取得公用的配置
func SharedGroupConfig() *GroupConfig {
	if sharedGroupConfig != nil {
		return sharedGroupConfig
	}
	config := &GroupConfig{
		Filename: "agents/group.conf",
		Groups:   []*Group{},
	}
	sharedGroupConfig = config
	file := files.NewFile(Tea.ConfigFile(config.Filename))
	if !file.Exists() {
		return config
	}
	data, err := file.ReadAll()
	if err != nil {
		logs.Error(err)
		return config
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		logs.Error(err)
	}
	return config
}

// 添加分组
func (this *GroupConfig) AddGroup(group *Group) {
	this.Groups = append(this.Groups, group)
}

// 删除分组
func (this *GroupConfig) RemoveGroup(groupId string) {
	result := []*Group{}
	for _, g := range this.Groups {
		if g.Id == groupId {
			continue
		}
		result = append(result, g)
	}
	this.Groups = result
}

// 保存
func (this *GroupConfig) Save() error {
	writer, err := files.NewWriter(Tea.ConfigFile(this.Filename))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}
