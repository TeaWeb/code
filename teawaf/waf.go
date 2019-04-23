package teawaf

import (
	"errors"
	"github.com/TeaWeb/code/teawaf/actions"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/utils/string"
	"io/ioutil"
	"net/http"
)

type WAF struct {
	Id         string             `yaml:"id" json:"id"`
	On         bool               `yaml:"on" json:"on"`
	Name       string             `yaml:"name" json:"name"`
	RuleGroups []*rules.RuleGroup `yaml:"ruleGroups" json:"ruleGroups"`

	hasRuleGroups bool
}

func NewWAF() *WAF {
	return &WAF{
		Id: stringutil.Rand(16),
		On: true,
	}
}

func NewWAFFromFile(path string) (waf *WAF, err error) {
	if len(path) == 0 {
		return nil, errors.New("'path' should not be empty")
	}
	file := files.NewFile(path)
	if !file.Exists() {
		return nil, errors.New("'path' not exist")
	}

	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}

	waf = &WAF{}
	err = reader.ReadYAML(waf)
	if err != nil {
		return nil, err
	}
	return waf, nil
}

func (this *WAF) Init() error {
	this.hasRuleGroups = len(this.RuleGroups) > 0

	if this.hasRuleGroups {
		for _, group := range this.RuleGroups {
			err := group.Init()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *WAF) AddRuleGroup(ruleGroup *rules.RuleGroup) {
	this.RuleGroups = append(this.RuleGroups, ruleGroup)
}

func (this *WAF) RemoveRuleGroup(ruleGroupId string) {
	if len(ruleGroupId) == 0 {
		return
	}
	result := []*rules.RuleGroup{}
	for _, group := range this.RuleGroups {
		if group.Id == ruleGroupId {
			continue
		}
		result = append(result, group)
	}
	this.RuleGroups = result
}

func (this *WAF) FindRuleGroup(ruleGroupId string) *rules.RuleGroup {
	if len(ruleGroupId) == 0 {
		return nil
	}
	for _, group := range this.RuleGroups {
		if group.Id == ruleGroupId {
			return group
		}
	}
	return nil
}

func (this *WAF) MoveRuleGroup(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.RuleGroups) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.RuleGroups) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	location := this.RuleGroups[fromIndex]
	result := []*rules.RuleGroup{}
	for i := 0; i < len(this.RuleGroups); i ++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			result = append(result, location)
		}
		result = append(result, this.RuleGroups[i])
		if fromIndex < toIndex && i == toIndex {
			result = append(result, location)
		}
	}

	this.RuleGroups = result
}

func (this *WAF) MatchRequest(req *http.Request, writer http.ResponseWriter) (goNext bool, set *rules.RuleSet, err error) {
	if !this.hasRuleGroups {
		return true, nil, nil
	}
	for _, group := range this.RuleGroups {
		if !group.On {
			continue
		}
		b, set, err := group.MatchRequest(req)
		if err != nil {
			return true, nil, err
		}
		if b {
			actionObject := actions.FindActionInstance(set.Action)
			if actionObject == nil {
				return true, set, errors.New("no action called '" + set.Action + "'")
			}
			goNext := actionObject.Perform(writer)
			return goNext, set, nil
		}
	}
	return true, nil, nil
}

func (this *WAF) MatchResponse(req *http.Request, resp *http.Response, writer http.ResponseWriter) (goNext bool, set *rules.RuleSet, err error) {
	if !this.hasRuleGroups {
		return true, nil, nil
	}
	for _, group := range this.RuleGroups {
		if !group.On {
			continue
		}
		b, set, err := group.MatchResponse(req, resp)
		if err != nil {
			return true, nil, err
		}
		if b {
			actionObject := actions.FindActionInstance(set.Action)
			if actionObject == nil {
				return true, set, errors.New("no action called '" + set.Action + "'")
			}
			goNext := actionObject.Perform(writer)
			return goNext, set, nil
		}
	}
	return true, nil, nil
}

// save to file path
func (this *WAF) Save(path string) error {
	if len(path) == 0 {
		return errors.New("path should not be empty")
	}
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func (this *WAF) ContainsGroupCode(code string) bool {
	if len(code) == 0 {
		return false
	}
	for _, group := range this.RuleGroups {
		if group.Code == code {
			return true
		}
	}
	return false
}

func (this *WAF) Copy() *WAF {
	waf := &WAF{
		Id:         this.Id,
		On:         this.On,
		Name:       this.Name,
		RuleGroups: this.RuleGroups,
	}
	return waf
}

func (this *WAF) CountRuleSets() int {
	count := 0
	for _, group := range this.RuleGroups {
		count += len(group.RuleSets)
	}
	return count
}
