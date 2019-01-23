package agentutils

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"sync"
)

var agentRuntimeMap = map[string]*agents.AgentConfig{} // 当前在运行的 agent id => agent
var agentRuntimeLocker = sync.Mutex{}

func FindAgentRuntime(agentConfig *agents.AgentConfig) *agents.AgentConfig {
	if agentConfig == nil {
		return nil
	}
	agentRuntimeLocker.Lock()
	defer agentRuntimeLocker.Unlock()

	agent, found := agentRuntimeMap[agentConfig.Id]
	if found {
		return agent
	} else {
		agentRuntimeMap[agentConfig.Id] = agentConfig
	}
	return agentConfig
}

func FindAgentApp(agent *agents.AgentConfig, appId string) *agents.AppConfig {
	app := agent.FindApp(appId)
	if app != nil {
		return app
	}
	app = FindAgentRuntime(agent).FindApp(appId)
	if app != nil {
		app.IsSystem = true
	}
	return app
}
