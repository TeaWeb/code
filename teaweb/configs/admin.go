package configs

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"io/ioutil"
	"sync"
)

// 管理员配置
type AdminConfig struct {
	// 安全设置
	Security *AdminSecurity `yaml:"security" json:"security"`

	// 角色
	Roles []*AdminRole `yaml:"roles" json:"roles"`

	// 权限
	Grant []*AdminGrant `yaml:"grant" json:"grant"`

	// 用户
	Users []*AdminUser `yaml:"users" json:"users"`
}

var adminConfig *AdminConfig
var adminConfigLocker sync.Mutex

// 读取全局的管理员配置
func SharedAdminConfig() *AdminConfig {
	adminConfigLocker.Lock()
	defer adminConfigLocker.Unlock()

	if adminConfig != nil {
		return adminConfig
	}

	adminConfig = &AdminConfig{}

	configFile := Tea.ConfigFile("admin.conf")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		logs.Error(err)
		return adminConfig
	}

	err = yaml.Unmarshal(data, adminConfig)
	if err != nil {
		logs.Error(err)
		return adminConfig
	}

	return adminConfig
}

// 写回配置文件
func (this *AdminConfig) Save() error {
	adminConfigLocker.Lock()
	defer adminConfigLocker.Unlock()

	writer, err := files.NewWriter(Tea.ConfigFile("admin.conf"))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}

// 是否包含某个激活的用户名
func (this *AdminConfig) ContainsActiveUser(username string) bool {
	adminConfigLocker.Lock()
	defer adminConfigLocker.Unlock()

	for _, user := range this.Users {
		if user.Username == username && !user.IsDisabled {
			return true
		}
	}
	return false
}

// 是否包含某个用户名
func (this *AdminConfig) ContainsUser(username string) bool {
	adminConfigLocker.Lock()
	defer adminConfigLocker.Unlock()

	for _, user := range this.Users {
		if user.Username == username {
			return true
		}
	}
	return false
}

// 使用用户名查找激活的用户
func (this *AdminConfig) FindActiveUser(username string) *AdminUser {
	adminConfigLocker.Lock()
	defer adminConfigLocker.Unlock()

	for _, user := range this.Users {
		if user.Username == username && !user.IsDisabled {
			return user
		}
	}
	return nil
}

// 使用用户名查找用户
func (this *AdminConfig) FindUser(username string) *AdminUser {
	adminConfigLocker.Lock()
	defer adminConfigLocker.Unlock()

	for _, user := range this.Users {
		if user.Username == username {
			return user
		}
	}
	return nil
}

// 使用Key查找用户
func (this *AdminConfig) FindUserWithKey(key string) *AdminUser {
	adminConfigLocker.Lock()
	defer adminConfigLocker.Unlock()

	if len(key) == 0 {
		return nil
	}

	for _, user := range this.Users {
		if user.Key == key {
			return user
		}
	}
	return nil
}

// 添加用户
func (this *AdminConfig) AddUser(user *AdminUser) {
	this.Users = append(this.Users, user)
}

// 根据代号查找激活的角色
func (this *AdminConfig) FindActiveRole(roleCode string) *AdminRole {
	for _, role := range this.Roles {
		if role.Code == roleCode && !role.IsDisabled {
			return role
		}
	}
	return nil
}

// 根据代号查找角色
func (this *AdminConfig) FindRole(roleCode string) *AdminRole {
	for _, role := range this.Roles {
		if role.Code == roleCode {
			return role
		}
	}
	return nil
}

// 查找激活的角色
func (this *AdminConfig) FindAllActiveRoles() []*AdminRole {
	result := []*AdminRole{}
	for _, role := range this.Roles {
		if !role.IsDisabled {
			result = append(result, role)
		}
	}
	return result
}

// 添加新角色
func (this *AdminConfig) AddRole(role *AdminRole) {
	this.Roles = append(this.Roles, role)
}

// 根据代号查找权限
func (this *AdminConfig) FindActiveGrant(grantCode string) *AdminGrant {
	for _, grant := range this.FindAllActiveGrants() {
		if grant.Code == grantCode && !grant.IsDisabled {
			return grant
		}
	}
	return nil
}

// 取得所有内置的权限
func (this *AdminConfig) FindAllActiveGrants() []*AdminGrant {
	grants := []*AdminGrant{
		NewAdminGrant("[超级权限]", AdminGrantAll),
		NewAdminGrant("代理", AdminGrantProxy),
		NewAdminGrant("日志", AdminGrantLog),
		NewAdminGrant("本地服务", AdminGrantAgent),
		NewAdminGrant("插件", AdminGrantPlugin),
	}

	if teaconst.PlusEnabled {
		grants = append(grants, []*AdminGrant{
			NewAdminGrant("测试小Q", AdminGrantQ),
			NewAdminGrant("API", AdminGrantApi),
			NewAdminGrant("团队", AdminGrantTeam),
		} ...)
	}

	return grants
}

// 添加授权
func (this *AdminConfig) AddGrant(grant *AdminGrant) {
	this.Grant = append(this.Grant, grant)
}

// 检查是否允许IP
func (this *AdminConfig) AllowIP(ip string) bool {
	// deny
	if len(this.Security.Deny) > 0 && lists.ContainsString(this.Security.Deny, ip) {
		return false
	}

	// allow
	if lists.ContainsString(this.Security.Allow, "all") || lists.ContainsString(this.Security.Allow, "0.0.0.0") || lists.ContainsString(this.Security.Allow, ip) {
		return true
	}

	return false
}

// 重置状态
func (this *AdminConfig) Reset() {
	adminConfigLocker.Lock()
	defer adminConfigLocker.Unlock()

	for _, u := range this.Users {
		u.Reset()
	}
}
