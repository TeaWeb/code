package mongo

import (
	"context"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"time"
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
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, uri)
	if err != nil {
		this.Message = "有错误需要修复：" + err.Error()
		this.Fail()
	}

	ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
	_, err = client.
		Database("teaweb").
		Collection("logs").
		Find(ctx, nil, findopt.Limit(1))
	if err != nil {
		this.Message = "有错误需要修复：" + err.Error()
		this.Fail()
	}

	client.Disconnect(context.Background())

	this.Success()
}
