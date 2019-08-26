package agentutils

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teadb"
	"github.com/iwind/TeaGo/maps"
)

func ActionDeleteAgent(agentId string, onFail func(message string)) (goNext bool) {
	agent := agents.NewAgentConfigFromId(agentId)
	if agent == nil {
		onFail("要删除的主机不存在")
		return
	}

	// 删除通知
	err := teadb.NoticeDAO().DeleteNoticesForAgent(agent.Id)
	if err != nil {
		onFail("通知删除失败：" + err.Error())
		return
	}

	// 删除数值记录
	err = teadb.AgentValueDAO().DropAgentTable(agent.Id)
	if err != nil {
		onFail("数值记录删除失败：" + err.Error())
		return
	}

	// 从列表删除
	agentList, err := agents.SharedAgentList()
	if err != nil {
		onFail("删除失败：" + err.Error())
		return
	}
	agentList.RemoveAgent(agent.Filename())
	err = agentList.Save()
	if err != nil {
		onFail("删除失败：" + err.Error())
		return
	}

	// 删除配置文件
	err = agent.Delete()
	if err != nil {
		onFail("删除失败：" + err.Error())
		return
	}

	// 通知更新
	PostAgentEvent(agent.Id, NewAgentEvent("REMOVE_AGENT", maps.Map{}))
	return true
}
