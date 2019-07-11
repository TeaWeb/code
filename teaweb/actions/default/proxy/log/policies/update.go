package policies

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/cmd"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction actions.Action

// 修改策略
func (this *UpdateAction) RunGet(params struct {
	PolicyId string
}) {
	policy := teaconfigs.NewAccessLogStoragePolicyFromId(params.PolicyId)
	if policy == nil {
		this.Fail("找不到Policy")
	}

	this.Data["policy"] = maps.Map{
		"id":      policy.Id,
		"name":    policy.Name,
		"type":    policy.Type,
		"options": policy.Options,
		"on":      policy.On,
	}

	this.Data["selectedMenu"] = "list"

	this.Data["storages"] = tealogs.AllStorages()
	this.Data["formats"] = tealogs.AllStorageFormats()

	this.Show()
}

// 保存提交
func (this *UpdateAction) RunPost(params struct {
	PolicyId string

	Name            string
	StorageFormat   string
	StorageTemplate string
	StorageType     string

	// file
	FilePath string

	// es
	EsEndpoint    string
	EsIndex       string
	EsMappingType string

	// mysql
	MysqlHost     string
	MysqlPort     int
	MysqlUsername string
	MysqlPassword string
	MysqlDatabase string
	MysqlTable    string
	MysqlLogField string

	// tcp
	TcpNetwork string
	TcpAddr    string

	// command
	CommandCommand string
	CommandArgs    string
	CommandDir     string

	On bool

	Must *actions.Must
}) {
	policy := teaconfigs.NewAccessLogStoragePolicyFromId(params.PolicyId)
	if policy == nil {
		this.Fail("找不到要修改的策略")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入日志策略的名称").
		Field("storageType", params.StorageType).
		Require("请选择存储类型")

	var instance interface{} = nil
	switch params.StorageType {
	case tealogs.StorageTypeFile:
		params.Must.
			Field("filePath", params.FilePath).
			Require("请输入日志文件路径")

		storage := new(tealogs.FileStorage)
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Path = params.FilePath
		instance = storage
	case tealogs.StorageTypeES:
		params.Must.
			Field("esEndpoint", params.EsEndpoint).
			Require("请输入Endpoint").
			Field("esIndex", params.EsIndex).
			Require("请输入Index名称").
			Field("esMappingType", params.EsMappingType).
			Require("请输入Mapping名称")

		storage := new(tealogs.ESStorage)
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Endpoint = params.EsEndpoint
		storage.Index = params.EsIndex
		storage.MappingType = params.EsMappingType
		instance = storage
	case tealogs.StorageTypeMySQL:
		params.Must.
			Field("mysqlHost", params.MysqlHost).
			Require("请输入主机地址").
			Field("mysqlDatabase", params.MysqlDatabase).
			Require("请输入数据库名称").
			Field("mysqlTable", params.MysqlTable).
			Require("请输入数据表名称").
			Field("mysqlLogField", params.MysqlLogField).
			Require("请输入日志存储字段名")

		storage := new(tealogs.MySQLStorage)
		storage.AutoCreateTable = true
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Host = params.MysqlHost
		storage.Port = params.MysqlPort
		storage.Username = params.MysqlUsername
		storage.Password = params.MysqlPassword
		storage.Database = params.MysqlDatabase
		storage.Table = params.MysqlTable
		storage.LogField = params.MysqlLogField
		instance = storage
	case tealogs.StorageTypeTCP:
		params.Must.
			Field("tcpNetwork", params.TcpNetwork).
			Require("请选择网络协议").
			Field("tcpAddr", params.TcpAddr).
			Require("请输入网络地址")

		storage := new(tealogs.TCPStorage)
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Network = params.TcpNetwork
		storage.Addr = params.TcpAddr
		instance = storage
	case tealogs.StorageTypeCommand:
		params.Must.
			Field("commandCommand", params.CommandCommand).
			Require("请输入可执行命令")

		storage := new(tealogs.CommandStorage)
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Command = params.CommandCommand
		storage.Args = cmd.ParseArgs(params.CommandArgs)
		storage.Dir = params.CommandDir
		instance = storage
	}

	if instance == nil {
		this.Fail("请选择存储类型")
	}

	policy.Type = params.StorageType
	policy.Name = params.Name
	policy.On = params.On

	options := map[string]interface{}{}
	teautils.ObjectToMapJSON(instance, &options)
	policy.Options = options

	err := policy.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 重置缓存策略
	tealogs.ResetPolicyStorage(policy.Id)

	this.Success()
}
