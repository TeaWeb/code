package agents

import (
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

// 分组配置
type GroupConfig struct {
	Filename string   `yaml:"filename" json:"filename"`
	Groups   []*Group `yaml:"groups" json:"groups"`
}

// 取得公用的配置
// 一定会返回一个不为nil的GroupConfig
func SharedGroupConfig() *GroupConfig {
	config := &GroupConfig{
		Filename: "agents/group.conf",
		Groups:   []*Group{},
	}
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

// 获取所有分组，包括默认分组
func (this *GroupConfig) FindAllGroups() []*Group {
	result := []*Group{}
	result = append(result, &Group{
		Name: "默认分组",
		Id:   "",
		On:   true,
	})
	result = append(result, this.Groups...)

	return result
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
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	writer, err := files.NewWriter(Tea.ConfigFile(this.Filename))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}

// 查找分组
func (this *GroupConfig) FindGroup(groupId string) *Group {
	for index, g := range this.Groups {
		if g.Id == groupId {
			g.Index = index
			return g
		}
	}
	return nil
}

// 移动位置
func (this *GroupConfig) Move(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.Groups) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.Groups) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	group := this.Groups[fromIndex]
	newList := []*Group{}
	for i := 0; i < len(this.Groups); i++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			newList = append(newList, group)
		}
		newList = append(newList, this.Groups[i])
		if fromIndex < toIndex && i == toIndex {
			newList = append(newList, group)
		}
	}

	this.Groups = newList
}
