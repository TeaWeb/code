package mongo

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/mongodb/mongo-go-driver/mongo"
	"context"
	"github.com/iwind/TeaGo/logs"
	"github.com/TeaWeb/code/teamongo"
	"time"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

type TestAction actions.Action

func (this *TestAction) Run(params struct {
	Host     string
	Port     uint
	Username string
	Password string
}) {
	config := configs.MongoConfig{
		Host:     params.Host,
		Port:     params.Port,
		Username: params.Username,
		Password: params.Password,
	}

	uri := config.URI()
	logs.Println(uri)
	client, err := mongo.Connect(context.Background(), uri)
	if err != nil {
		this.Message = "有错误需要修复：" + err.Error()
		this.Fail()
	}
	defer client.Disconnect(context.Background())

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err = client.
		Database("teaweb").
		Collection("logs").
		Find(ctx, nil, findopt.Limit(1))
	if err != nil {
		this.Message = "有错误需要修复：" + err.Error()
		this.Fail()
	}

	teamongo.Test()

	this.Success()
}
