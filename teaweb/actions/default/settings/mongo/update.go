package mongo

import (
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"strings"
)

type UpdateAction actions.Action

// 修改连接
func (this *UpdateAction) Run(params struct{}) {
	config := configs.SharedMongoConfig()
	this.Data["config"] = maps.Map{
		"scheme":                  config.Scheme,
		"username":                config.Username,
		"password":                "",
		"host":                    config.Host,
		"port":                    config.Port,
		"authMechanism":           config.AuthMechanism,
		"authMechanismProperties": config.AuthMechanismPropertiesString(),
		"requestURI":              config.RequestURI,
	}

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	Host                    string
	Port                    uint
	Username                string
	Password                string
	AuthMechanism           string
	AuthMechanismProperties string

	Must *actions.Must
}) {
	params.Must.
		Field("host", params.Host).
		Require("请输入主机地址").
		Field("port", params.Port).
		Require("请输入端口").
		Gt(0, "请输入正确的端口")

	config := configs.SharedMongoConfig()
	config.Host = params.Host
	config.Port = params.Port
	config.AuthMechanism = params.AuthMechanism
	config.AuthMechanismProperties = []*shared.Variable{}

	if len(params.AuthMechanismProperties) > 0 {
		properties := regexp.MustCompile("\\s*,\\s*").Split(params.AuthMechanismProperties, -1)
		for _, property := range properties {
			if strings.Contains(property, ":") {
				pieces := strings.Split(property, ":")
				config.AuthMechanismProperties = append(config.AuthMechanismProperties, shared.NewVariable(pieces[0], pieces[1]))
			}
		}
	}

	config.Username = params.Username
	if len(params.Username) > 0 {
		if len(params.Password) > 0 {
			config.Password = params.Password
		}
	} else {
		config.Password = ""
	}
	err := config.Save()

	if err != nil {
		this.Fail("文件写入失败，请检查'configs/mongo.conf'写入权限")
	}

	// 重新连接
	teamongo.RestartClient()

	this.Next("/settings/mongo", nil).Success("保存成功")
}
