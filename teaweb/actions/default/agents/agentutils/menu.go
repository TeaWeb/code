package agentutils

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/TeaWeb/code/teaweb/utils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"net/http"
)

func AddTabbar(actionWrapper actions.ActionWrapper) {
	if actionWrapper.Object().Request.Method != http.MethodGet {
		return
	}

	action := actionWrapper.Object()
	action.Data["teaMenu"] = "agents"
	action.Data["teaSubHeader"] = "Agent主机"

	// 子菜单
	subMenu := utils.NewSubMenu()
	agentId := action.ParamString("agentId")
	if len(agentId) == 0 {
		agentId = "local"
	}

	actionCode := "board"
	if action.HasPrefix("/agents/apps") {
		actionCode = "apps"
	} else if action.HasPrefix("/agents/settings") {
		actionCode = "settings"
	} else if action.HasPrefix("/agents/delete") {
		actionCode = "delete"
	} else if action.HasPrefix("/agents/notices") {
		actionCode = "notices"
	}

	isWaiting := CheckAgentIsWaiting("local")
	if isWaiting {
		subMenu.Add("本地", "已连接", "/agents/"+actionCode+"?agentId=local", agentId == "local" && !action.HasPrefix("/agents/addAgent"))
	} else {
		subMenu.Add("本地", "", "/agents/"+actionCode+"?agentId=local", agentId == "local" && !action.HasPrefix("/agents/addAgent"))
	}

	// agent列表
	agentList, err := agents.SharedAgentList()
	if err != nil {
		logs.Error(err)
	} else {
		for _, agent := range agentList.FindAllAgents() {
			isWaiting := CheckAgentIsWaiting(agent.Id)
			if isWaiting {
				subMenu.Add(agent.Name, "已连接", "/agents/"+actionCode+"?agentId="+agent.Id, agentId == agent.Id)
			} else {
				subMenu.Add(agent.Name, "", "/agents/"+actionCode+"?agentId="+agent.Id, agentId == agent.Id)
			}
		}
	}

	subMenu.Add("[添加新主机]", "", "/agents/addAgent", action.HasPrefix("/agents/addAgent"))
	utils.SetSubMenu(action, subMenu)

	// Tabbar
	if !action.HasPrefix("/agents/addAgent") {
		agent := agents.NewAgentConfigFromId(agentId)

		tabbar := utils.NewTabbar()

		// 看板和Apps
		tabbar.Add("看板", "", "/agents/board?agentId="+agentId, "dashboard", action.HasPrefix("/agents/board"))
		tabbar.Add("Apps", fmt.Sprintf("%d", len(agent.Apps)+len(FindAgentRuntime(agent).FindSystemApps())), "/agents/apps?agentId="+agentId, "gem outline", action.HasPrefix("/agents/apps"))

		// 通知
		countUnreadNotices := noticeutils.CountUnreadNoticesForAgent(agentId)
		if countUnreadNotices > 0 {
			tabbar.Add("通知", fmt.Sprintf("%d", countUnreadNotices), "/agents/notices?agentId="+agentId, "bell blink orange", action.HasPrefix("/agents/notices"))
		} else {
			tabbar.Add("通知", fmt.Sprintf("%d", countUnreadNotices), "/agents/notices?agentId="+agentId, "bell", action.HasPrefix("/agents/notices"))
		}

		// 设置和删除
		if agentId != "local" {
			tabbar.Add("设置", "", "/agents/settings?agentId="+agentId, "setting", action.HasPrefix("/agents/settings"))
			tabbar.Add("删除", "", "/agents/delete?agentId="+agentId, "trash", action.HasPrefix("/agents/delete"))
		}
		utils.SetTabbar(actionWrapper, tabbar)
	}
}
