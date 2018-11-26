package api

import (
	"errors"
	"github.com/iwind/TeaGo/lists"
)

// 服务的API配置
type APIConfig struct {
	On         bool         `yaml:"apiOn" json:"apiOn"`               // 是否开启API功能
	Files      []string     `yaml:"apiFiles" json:"apiFiles"`         // API文件列表
	Groups     []string     `yaml:"apiGroups" json:"apiGroups"`       // API分组
	Versions   []string     `yaml:"apiVersions" json:"apiVersions"`   // API版本
	TestPlans  []string     `yaml:"apiTestPlans" json:"apiTestPlans"` // API测试计划
	Limit      *APILimit    `yaml:"apiLimit" json:"apiLimit"`         // API全局的限制 TODO
	StatusList []*APIStatus `yaml:"status" json:"status"`             // 状态码列表

	pathMap    map[string]*API // path => api
	patternMap map[string]*API // path => api
}

// 获取新对象
func NewAPIConfig() *APIConfig {
	return &APIConfig{}
}

// 校验
func (this *APIConfig) Validate() error {
	this.pathMap = map[string]*API{}
	this.patternMap = map[string]*API{}
	for _, apiFilename := range this.Files {
		api := NewAPIFromFile(apiFilename)
		if api == nil {
			continue
		}
		err := api.Validate()
		if err != nil {
			return err
		}
		if api.pathReg == nil {
			this.pathMap[api.Path] = api
		} else {
			this.patternMap[api.Path] = api
		}
	}

	// api limit
	if this.Limit != nil {
		err := this.Limit.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// 添加API
func (this *APIConfig) AddAPI(api *API) {
	if api == nil {
		return
	}

	// 分析API
	if this.pathMap != nil {
		err := api.Validate()
		if err == nil {
			if api.pathReg == nil {
				this.pathMap[api.Path] = api
			} else {
				this.patternMap[api.Path] = api
			}
		}
	}

	// 如果已包含文件名则不重复添加
	if lists.Contains(this.Files, api.Filename) {
		return
	}
	this.Files = append(this.Files, api.Filename)
}

// 获取所有APIs
func (this *APIConfig) FindAllAPIs() []*API {
	apis := []*API{}
	for _, filename := range this.Files {
		api := NewAPIFromFile(filename)
		if api == nil {
			continue
		}
		apis = append(apis, api)
	}
	return apis
}

// 获取单个API信息
func (this *APIConfig) FindAPI(path string) *API {
	for _, api := range this.FindAllAPIs() {
		if api.Path == path {
			return api
		}
	}
	return nil
}

// 查找激活状态中的API
func (this *APIConfig) FindActiveAPI(path string, method string) (api *API, params map[string]string) {
	api, found := this.pathMap[path]
	if !found {
		// 寻找pattern
		for _, api := range this.patternMap {
			params, found := api.Match(path)
			if !found || api.IsDeprecated || !api.On || !api.AllowMethod(method) {
				continue
			}
			return api, params
		}

		return nil, nil
	}

	// 检查是否过期或者失效
	if api.IsDeprecated || !api.On || !api.AllowMethod(method) {
		return nil, nil
	}

	return api, nil
}

// 删除API
func (this *APIConfig) DeleteAPI(api *API) {
	this.Files = lists.Delete(this.Files, api.Filename).([]string)

	delete(this.pathMap, api.Path)
	delete(this.patternMap, api.Path)
}

// 添加API分组
func (this *APIConfig) AddAPIGroup(name string) {
	this.Groups = append(this.Groups, name)
}

// 删除API分组
func (this *APIConfig) RemoveAPIGroup(name string) {
	result := []string{}
	for _, groupName := range this.Groups {
		if groupName != name {
			result = append(result, groupName)
		}
	}

	for _, filename := range this.Files {
		api := NewAPIFromFile(filename)
		if api == nil {
			continue
		}
		api.RemoveGroup(name)
		api.Save()
	}

	this.Groups = result
}

// 修改API分组
func (this *APIConfig) ChangeAPIGroup(oldName string, newName string) {
	result := []string{}
	for _, groupName := range this.Groups {
		if groupName == oldName {
			result = append(result, newName)
		} else {
			result = append(result, groupName)
		}
	}

	for _, filename := range this.Files {
		api := NewAPIFromFile(filename)
		if api == nil {
			continue
		}
		api.ChangeGroup(oldName, newName)
		api.Save()
	}

	this.Groups = result
}

// 把API分组往上调整
func (this *APIConfig) MoveUpAPIGroup(name string) {
	index := lists.Index(this.Groups, name)
	if index <= 0 {
		return
	}
	this.Groups[index], this.Groups[index-1] = this.Groups[index-1], this.Groups[index]
}

// 把API分组往下调整
func (this *APIConfig) MoveDownAPIGroup(name string) {
	index := lists.Index(this.Groups, name)
	if index < 0 {
		return
	}
	this.Groups[index], this.Groups[index+1] = this.Groups[index+1], this.Groups[index]
}

// 添加API版本
func (this *APIConfig) AddAPIVersion(name string) {
	this.Versions = append(this.Versions, name)
}

// 删除API版本
func (this *APIConfig) RemoveAPIVersion(name string) {
	result := []string{}
	for _, versionName := range this.Versions {
		if versionName != name {
			result = append(result, versionName)
		}
	}

	for _, filename := range this.Files {
		api := NewAPIFromFile(filename)
		if api == nil {
			continue
		}
		api.RemoveVersion(name)
		api.Save()
	}

	this.Versions = result
}

// 修改API版本
func (this *APIConfig) ChangeAPIVersion(oldName string, newName string) {
	result := []string{}
	for _, versionName := range this.Versions {
		if versionName == oldName {
			result = append(result, newName)
		} else {
			result = append(result, versionName)
		}
	}

	for _, filename := range this.Files {
		api := NewAPIFromFile(filename)
		if api == nil {
			continue
		}
		api.ChangeVersion(oldName, newName)
		api.Save()
	}

	this.Versions = result
}

// 把API版本往上调整
func (this *APIConfig) MoveUpAPIVersion(name string) {
	index := lists.Index(this.Versions, name)
	if index <= 0 {
		return
	}
	this.Versions[index], this.Versions[index-1] = this.Versions[index-1], this.Versions[index]
}

// 把API版本往下调整
func (this *APIConfig) MoveDownAPIVersion(name string) {
	index := lists.Index(this.Versions, name)
	if index < 0 {
		return
	}
	this.Versions[index], this.Versions[index+1] = this.Versions[index+1], this.Versions[index]
}

// 添加测试计划
func (this *APIConfig) AddTestPlan(filename string) {
	this.TestPlans = append(this.TestPlans, filename)
}

// 查找所有测试计划
func (this *APIConfig) FindTestPlans() []*APITestPlan {
	result := []*APITestPlan{}
	for _, filename := range this.TestPlans {
		plan := NewAPITestPlanFromFile(filename)
		if plan != nil {
			result = append(result, plan)
		}
	}
	return result
}

// 删除某个测试计划
func (this *APIConfig) DeleteTestPlan(filename string) error {
	if len(filename) == 0 {
		return errors.New("filename should not be empty")
	}

	plan := NewAPITestPlanFromFile(filename)
	if plan != nil {
		err := plan.Delete()
		if err != nil {
			return err
		}
	}

	this.TestPlans = lists.Delete(this.TestPlans, filename).([]string)

	return nil
}
