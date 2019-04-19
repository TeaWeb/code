package teawaf

import (
	"net/http"
)

// rule group
type RuleGroup struct {
	On       bool       `yaml:"on" json:"on"`
	Name     string     `yaml:"name" json:"name"` // such as SQL Injection
	RuleSets []*RuleSet `yaml:"ruleSets" json:"ruleSets"`

	hasRuleSets bool
}

func NewRuleGroup() *RuleGroup {
	return &RuleGroup{}
}

func (this *RuleGroup) Init() error {
	this.hasRuleSets = len(this.RuleSets) > 0

	if this.hasRuleSets {
		for _, set := range this.RuleSets {
			err := set.Init()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *RuleGroup) AddRuleSet(ruleSet *RuleSet) {
	this.RuleSets = append(this.RuleSets, ruleSet)
}

func (this *RuleGroup) MatchRequest(req *http.Request) (b bool, set *RuleSet, err error) {
	if !this.hasRuleSets {
		return
	}
	for _, set := range this.RuleSets {
		b, err = set.MatchRequest(req)
		if err != nil {
			return false, nil, err
		}
		if b {
			return true, set, nil
		}
	}
	return
}

func (this *RuleGroup) MatchResponse(req *http.Request, resp *http.Response) (b bool, set *RuleSet, err error) {
	if !this.hasRuleSets {
		return
	}
	for _, set := range this.RuleSets {
		b, err = set.MatchResponse(req, resp)
		if err != nil {
			return false, nil, err
		}
		if b {
			return true, set, nil
		}
	}
	return
}
