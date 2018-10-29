package configs

import "sync"

type AdminUser struct {
	Username string   `yaml:"username" json:"username"` // 用户名
	Password string   `yaml:"password" json:"password"` // 密码
	Role     []string `yaml:"role" json:"role"`         // 角色

	Name      string `yaml:"name" json:"name"`           // 姓名
	Avatar    string `yaml:"avatar" json:"avatar"`       // 头像
	Tel       string `yaml:"tel" json:"tel"`             // 联系电话
	CreatedAt int64  `yaml:"createdAt" json:"createdAt"` // 创建时间
	LoggedAt  int64  `yaml:"loggedAt" json:"loggedAt"`   // 最后登录时间
	LoggedIP  string `yaml:"loggedIP" json:"loggedIP"`   // 最后登录IP

	countLoginTries uint // 错误登录次数
	locker          sync.Mutex
}

func (this *AdminUser) IncreaseLoginTries() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.countLoginTries ++
}

func (this *AdminUser) CountLoginTries() uint {
	this.locker.Lock()
	defer this.locker.Unlock()
	return this.countLoginTries
}

func (this *AdminUser) ResetLoginTries() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.countLoginTries = 0
}
