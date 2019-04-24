package teawaf

import (
	"errors"
	"github.com/TeaWeb/code/teawaf/actions"
	"github.com/TeaWeb/code/teawaf/requests"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/utils/string"
	"io/ioutil"
	"net/http"
)

type WAF struct {
	Id       string             `yaml:"id" json:"id"`
	On       bool               `yaml:"on" json:"on"`
	Name     string             `yaml:"name" json:"name"`
	Inbound  []*rules.RuleGroup `yaml:"inbound" json:"inbound"`
	Outbound []*rules.RuleGroup `yaml:"outbound" json:"outbound"`

	hasInboundRules  bool
	hasOutboundRules bool
	onActionCallback func(action actions.ActionString) (goNext bool)
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
	this.hasInboundRules = len(this.Inbound) > 0
	this.hasOutboundRules = len(this.Outbound) > 0

	if this.hasInboundRules {
		for _, group := range this.Inbound {
			err := group.Init()
			if err != nil {
				return err
			}
		}
	}

	if this.hasOutboundRules {
		for _, group := range this.Outbound {
			err := group.Init()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *WAF) AddRuleGroup(ruleGroup *rules.RuleGroup) {
	if ruleGroup.IsInbound {
		this.Inbound = append(this.Inbound, ruleGroup)
	} else {
		this.Outbound = append(this.Outbound, ruleGroup)
	}
}

func (this *WAF) RemoveRuleGroup(ruleGroupId string) {
	if len(ruleGroupId) == 0 {
		return
	}

	{
		result := []*rules.RuleGroup{}
		for _, group := range this.Inbound {
			if group.Id == ruleGroupId {
				continue
			}
			result = append(result, group)
		}
		this.Inbound = result
	}

	{
		result := []*rules.RuleGroup{}
		for _, group := range this.Outbound {
			if group.Id == ruleGroupId {
				continue
			}
			result = append(result, group)
		}
		this.Outbound = result
	}
}

func (this *WAF) FindRuleGroup(ruleGroupId string) *rules.RuleGroup {
	if len(ruleGroupId) == 0 {
		return nil
	}
	for _, group := range this.Inbound {
		if group.Id == ruleGroupId {
			return group
		}
	}
	for _, group := range this.Outbound {
		if group.Id == ruleGroupId {
			return group
		}
	}
	return nil
}

func (this *WAF) MoveInboundRuleGroup(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.Inbound) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.Inbound) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	group := this.Inbound[fromIndex]
	result := []*rules.RuleGroup{}
	for i := 0; i < len(this.Inbound); i ++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			result = append(result, group)
		}
		result = append(result, this.Inbound[i])
		if fromIndex < toIndex && i == toIndex {
			result = append(result, group)
		}
	}

	this.Inbound = result
}

func (this *WAF) MoveOutboundRuleGroup(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.Outbound) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.Outbound) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	group := this.Outbound[fromIndex]
	result := []*rules.RuleGroup{}
	for i := 0; i < len(this.Outbound); i ++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			result = append(result, group)
		}
		result = append(result, this.Outbound[i])
		if fromIndex < toIndex && i == toIndex {
			result = append(result, group)
		}
	}

	this.Outbound = result
}

func (this *WAF) MatchRequest(req *requests.Request, writer http.ResponseWriter) (goNext bool, set *rules.RuleSet, err error) {
	if !this.hasInboundRules {
		return true, nil, nil
	}
	for _, group := range this.Inbound {
		if !group.On {
			continue
		}
		b, set, err := group.MatchRequest(req)
		if err != nil {
			return true, nil, err
		}
		if b {
			if this.onActionCallback == nil {
				actionObject := actions.FindActionInstance(set.Action)
				if actionObject == nil {
					return true, set, errors.New("no action called '" + set.Action + "'")
				}
				goNext := actionObject.Perform(writer)
				return goNext, set, nil
			} else {
				goNext = this.onActionCallback(set.Action)
			}
			return goNext, set, nil
		}
	}
	return true, nil, nil
}

func (this *WAF) MatchResponse(req *requests.Request, resp *http.Response, writer http.ResponseWriter) (goNext bool, set *rules.RuleSet, err error) {
	if !this.hasOutboundRules {
		return true, nil, nil
	}
	for _, group := range this.Outbound {
		if !group.On {
			continue
		}
		b, set, err := group.MatchResponse(req, resp)
		if err != nil {
			return true, nil, err
		}
		if b {
			if this.onActionCallback == nil {
				actionObject := actions.FindActionInstance(set.Action)
				if actionObject == nil {
					return true, set, errors.New("no action called '" + set.Action + "'")
				}
				goNext := actionObject.Perform(writer)
				return goNext, set, nil
			} else {
				goNext = this.onActionCallback(set.Action)
			}
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
	for _, group := range this.Inbound {
		if group.Code == code {
			return true
		}
	}
	for _, group := range this.Outbound {
		if group.Code == code {
			return true
		}
	}
	return false
}

func (this *WAF) Copy() *WAF {
	waf := &WAF{
		Id:       this.Id,
		On:       this.On,
		Name:     this.Name,
		Inbound:  this.Inbound,
		Outbound: this.Outbound,
	}
	return waf
}

func (this *WAF) CountInboundRuleSets() int {
	count := 0
	for _, group := range this.Inbound {
		count += len(group.RuleSets)
	}
	return count
}

func (this *WAF) CountOutboundRuleSets() int {
	count := 0
	for _, group := range this.Outbound {
		count += len(group.RuleSets)
	}
	return count
}

func (this *WAF) OnAction(onActionCallback func(action actions.ActionString) (goNext bool)) {
	this.onActionCallback = onActionCallback
}
