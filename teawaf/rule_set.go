package teawaf

import (
	"github.com/TeaWeb/code/teawaf/actions"
	"github.com/iwind/TeaGo/maps"
	"net/http"
)

type RuleConnector = string

const (
	RuleConnectorAnd = "and"
	RuleConnectorOr  = "or"
)

type RuleSet struct {
	On        bool          `yaml:"on" json:"on"`
	Name      string        `yaml:"name" json:"name"`
	Rules     []*Rule       `yaml:"rules" json:"rules"`
	Connector RuleConnector `yaml:"connector" json:"connector"` // rules connector

	Action        actions.ActionString `yaml:"action" json:"action"`
	ActionOptions maps.Map             `yaml:"actionOptions" json:"actionOptions"` // TODO TO BE IMPLEMENTED

	hasRules bool
}

func NewRuleSet() *RuleSet {
	return &RuleSet{}
}

func (this *RuleSet) Init() error {
	this.hasRules = len(this.Rules) > 0
	if this.hasRules {
		for _, rule := range this.Rules {
			err := rule.Init()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *RuleSet) AddRule(rule ...*Rule) {
	this.Rules = append(this.Rules, rule...)
}

func (this *RuleSet) MatchRequest(req *http.Request) (b bool, err error) {
	if !this.hasRules {
		return false, nil
	}
	switch this.Connector {
	case RuleConnectorAnd:
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchRequest(req)
			if err1 != nil {
				return false, err1
			}
			if !b1 {
				return false, nil
			}
		}
		return true, nil
	case RuleConnectorOr:
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchRequest(req)
			if err1 != nil {
				return false, err1
			}
			if b1 {
				return true, nil
			}
		}
	default: // same as And
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchRequest(req)
			if err1 != nil {
				return false, err1
			}
			if !b1 {
				return false, nil
			}
		}
		return true, nil
	}
	return
}

func (this *RuleSet) MatchResponse(req *http.Request, resp *http.Response) (b bool, err error) {
	if !this.hasRules {
		return false, nil
	}
	switch this.Connector {
	case RuleConnectorAnd:
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchResponse(req, resp)
			if err1 != nil {
				return false, err1
			}
			if !b1 {
				return false, nil
			}
		}
		return true, nil
	case RuleConnectorOr:
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchResponse(req, resp)
			if err1 != nil {
				return false, err1
			}
			if b1 {
				return true, nil
			}
		}
	default: // same as And
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchResponse(req, resp)
			if err1 != nil {
				return false, err1
			}
			if !b1 {
				return false, nil
			}
		}
		return true, nil
	}
	return
}
