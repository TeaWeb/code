package teawaf

import (
	"github.com/TeaWeb/code/teawaf/actions"
	"github.com/pkg/errors"
	"net/http"
)

type WAF struct {
	RuleGroups []*RuleGroup `yaml:"ruleGroups" json:"ruleGroups"`

	hasRuleGroups bool
}

func NewWAF() *WAF {
	return &WAF{}
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

func (this *WAF) AddRuleGroup(ruleGroup *RuleGroup) {
	this.RuleGroups = append(this.RuleGroups, ruleGroup)
}

func (this *WAF) MatchRequest(req *http.Request, writer http.ResponseWriter) (goNext bool, set *RuleSet, err error) {
	if !this.hasRuleGroups {
		return true, nil, nil
	}
	for _, group := range this.RuleGroups {
		b, set, err := group.MatchRequest(req)
		if err != nil {
			return true, nil, err
		}
		if b {
			actionObject := actions.FindAction(set.Action)
			if actionObject == nil {
				return true, set, errors.New("no action called '" + set.Action + "'")
			}
			goNext := actionObject.Perform(writer)
			return goNext, set, nil
		}
	}
	return true, nil, nil
}

func (this *WAF) MatchResponse(req *http.Request, resp *http.Response, writer http.ResponseWriter) (goNext bool, set *RuleSet, err error) {
	if !this.hasRuleGroups {
		return true, nil, nil
	}
	for _, group := range this.RuleGroups {
		b, set, err := group.MatchResponse(req, resp)
		if err != nil {
			return true, nil, err
		}
		if b {
			actionObject := actions.FindAction(set.Action)
			if actionObject == nil {
				return true, set, errors.New("no action called '" + set.Action + "'")
			}
			goNext := actionObject.Perform(writer)
			return goNext, set, nil
		}
	}
	return true, nil, nil
}
