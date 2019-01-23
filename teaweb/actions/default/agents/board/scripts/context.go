package scripts

import "github.com/TeaWeb/code/teaconfigs/agents"

type Context struct {
	Agent *agents.AgentConfig
	App   *agents.AppConfig
	Item  *agents.Item
}
