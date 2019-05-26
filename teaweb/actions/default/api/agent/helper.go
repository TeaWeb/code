package agent

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"net"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action actions.ActionObject) bool {
	agentId := action.Header("Tea-Agent-Id")
	if len(agentId) == 0 {
		action.Fail("Authenticate Failed 001")
	}

	key := action.Header("Tea-Agent-Key")
	if len(key) == 0 {
		action.Fail("Authenticate Failed 002")
	}

	agent := agents.NewAgentConfigFromId(agentId)
	if agent == nil {
		action.Fail("Authenticate Failed 003")
	}
	if agent.Id != agentId || agent.Key != key {
		action.Fail("Authenticate Failed 004")
	}

	// 检查IP
	addr := action.Request.RemoteAddr
	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		addr = host
	}
	if !agent.IsLocal() && !agent.AllowAll && !lists.ContainsString(agent.Allow, addr) {
		action.Fail("Access Denied 005")
	}

	action.Context.Set("agent", agent)

	return true
}
