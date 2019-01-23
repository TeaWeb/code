package utils

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

// 子菜单定义
type SubMenu struct {
	items []maps.Map
}

// 获取新对象
func NewSubMenu() *SubMenu {
	return &SubMenu{
		items: []maps.Map{},
	}
}

// 添加菜单项
func (this *SubMenu) Add(name string, subName string, url string, active bool) {
	this.items = append(this.items, maps.Map{
		"name":    name,
		"subName": subName,
		"url":     url,
		"active":  active,
	})
}

// 取得所有的Items
func (this *SubMenu) Items() []maps.Map {
	return this.items
}

// 设置子菜单
func SetSubMenu(action actions.ActionWrapper, subMenu *SubMenu) {
	action.Object().Data["teaSubMenus"] = subMenu.Items()
}
