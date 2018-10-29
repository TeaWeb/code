package configs

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

// 管理员配置
type AdminConfig struct {
	// 安全设置
	Security *AdminSecurity `yaml:"security" json:"security"`

	// 角色
	Roles []*AdminRole `yaml:"roles" json:"roles"`

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
func (this *AdminConfig) WriteBack() error {
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
