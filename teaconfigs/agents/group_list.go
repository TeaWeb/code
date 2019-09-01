package agents

import (
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaconst"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"io/ioutil"
)

const (
	groupListFilename = "agents/grouplist.conf"
)

// 分组配置
type GroupList struct {
	Groups     []*Group `yaml:"groups" json:"groups"`
	TeaVersion string   `yaml:"teaVersion" json:"teaVersion"`
}

// 取得公用的配置
// 一定会返回一个不为nil的GroupConfig
func SharedGroupList() *GroupList {
	shared.Locker.Lock()
	groupList := &GroupList{
		Groups: []*Group{},
	}
	data, err := ioutil.ReadFile(Tea.ConfigFile(groupListFilename))
	if err != nil {
		// 默认分组
		defaultGroup := &Group{
			On:        true,
			Id:        "default",
			IsDefault: true,
			Name:      "默认分组",
		}
		groupList.AddGroup(defaultGroup)

		// 老的默认分组
		oldDefault := loadOldDefaultGroup()
		if oldDefault != nil {
			defaultGroup.Name = oldDefault.Name
		}

		// 升级
		oldList := oldGroupConfig()
		for _, g := range oldList.Groups {
			if g.Id == "" || g.Id == "default" {
				defaultGroup.NoticeSetting = g.NoticeSetting
				continue
			}
			groupList.AddGroup(g)
		}

		shared.Locker.ReadUnlock()
		err = groupList.Save()
		if err != nil {
			logs.Error(err)
		}

		return groupList
	}

	err = yaml.Unmarshal(data, groupList)
	if err != nil {
		logs.Error(err)
	}

	shared.Locker.ReadUnlock()
	return groupList
}

// 0.1.7版本之前的GroupConfig
// deprecated in 0.1.8
func oldGroupConfig() *GroupList {
	config := &GroupList{
		Groups: []*Group{},
	}
	file := files.NewFile(Tea.ConfigFile("agents/group.conf"))
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
func (this *GroupList) FindAllGroups() []*Group {
	return this.Groups
}

// 添加分组
func (this *GroupList) AddGroup(group *Group) {
	this.Groups = append(this.Groups, group)
}

// 删除分组
func (this *GroupList) RemoveGroup(groupId string) {
	// 默认分组不能删除
	if groupId == "default" {
		return
	}
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
func (this *GroupList) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	this.TeaVersion = teaconst.TeaVersion
	writer, err := files.NewWriter(Tea.ConfigFile(groupListFilename))
	if err != nil {
		return err
	}
	defer func() {
		err1 := writer.Close()
		if err1 != nil {
			logs.Error(err1)
		}
	}()
	_, err = writer.WriteYAML(this)
	return err
}

// 查找分组
func (this *GroupList) FindGroup(groupId string) *Group {
	if groupId == "" {
		groupId = "default"
	}
	for index, g := range this.Groups {
		if g.Id == groupId {
			g.Index = index
			return g
		}
	}
	return nil
}

// 查找默认的分组
func (this *GroupList) FindDefaultGroup() *Group {
	return this.FindGroup("default")
}

// 移动位置
func (this *GroupList) Move(fromIndex int, toIndex int) {
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
