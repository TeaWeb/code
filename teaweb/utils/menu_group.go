package utils

import "github.com/iwind/TeaGo/actions"

// 菜单分组
type MenuGroup struct {
	Menus []*Menu `json:"menus"`
}

// 获取新菜单分组对象
func NewMenuGroup() *MenuGroup {
	return &MenuGroup{
		Menus: []*Menu{},
	}
}

// 查找菜单，如果找不到则自动创建
func (this *MenuGroup) FindMenu(menuId string, menuName string) *Menu {
	for _, m := range this.Menus {
		if m.Id == menuId {
			return m
		}
	}
	menu := NewMenu()
	menu.Id = menuId
	menu.Name = menuName
	menu.Items = []*MenuItem{}
	this.Menus = append(this.Menus, menu)
	return menu
}

// 设置子菜单
func SetSubMenu(action actions.ActionWrapper, menu *MenuGroup) {
	action.Object().Data["teaSubMenus"] = menu
}
