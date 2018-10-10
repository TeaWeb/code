package install

import (
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/maps"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/TeaWeb/code/teamongo"
)

type MongoAction actions.Action

func (this *MongoAction) Run(params struct {
	Auth *helpers.UserMustAuth
}) {
	errorMessage := ""
	uriString := ""

	configFile := files.NewFile(Tea.ConfigFile("mongo.conf"))
	if !configFile.Exists() {
		errorMessage = "'conf/mongo.conf' 文件不存在"
	} else {
		reader, err := configFile.Reader()
		if err != nil {
			errorMessage = "无法读取文件 'conf/mongo.conf'"
		} else {
			defer reader.Close()

			m := maps.Map{}
			err := reader.ReadYAML(m)
			if err != nil {
				errorMessage = "文件读取错误：" + err.Error()
			} else {
				uriString = m.GetString("uri")

				if len(uriString) == 0 {
					errorMessage = "配置缺少uri字段"
				}
			}
		}
	}

	this.Data["error"] = errorMessage
	this.Data["uri"] = uriString

	if len(uriString) > 0 {
		_, err := mongo.NewClient(uriString)
		if err != nil {
			this.Data["error"] = "uri错误：" + err.Error()
		} else {
			err := teamongo.Test()
			if err != nil {
				this.Data["error"] = "mongodb连接错误：" + err.Error()
			}
		}
	}

	this.Show()
}

func (this *MongoAction) RunPost(params struct {
	Auth *helpers.UserMustAuth
	URI  string `alias:"uri"`
	Must *actions.Must
}) {
	params.Must.
		Field("uri", params.URI).
		Require("请输入MongoDB连接DSN")

	reader, err := files.NewReader(Tea.ConfigFile("mongo.conf"))
	if err != nil {
		this.Fail(err.Error())
	}
	defer reader.Close()

	config := maps.Map{}
	err = reader.ReadYAML(&config)
	if err != nil {
		this.Fail(err.Error())
	}

	config["uri"] = params.URI
	writer, err := files.NewWriter(Tea.ConfigFile("mongo.conf"))
	if err != nil {
		this.Fail(err.Error())
	}
	defer writer.Close()

	_, err = writer.WriteYAML(config)
	if err != nil {
		this.Fail(err.Error())
	}

	// 重新连接
	teamongo.RestartClient()

	this.Refresh().Success("保存成功")
}
