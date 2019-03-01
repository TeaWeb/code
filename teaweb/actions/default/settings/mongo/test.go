package mongo

import (
	"context"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		this.Message = "有错误需要修复：" + err.Error()
		this.Fail()
	}

	ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
	_, err = client.
		Database("teaweb").
		Collection("logs").
		Find(ctx, map[string]interface{}{}, options.Find().SetLimit(1))
	if err != nil {
		this.Message = "有错误需要修复：" + err.Error()
		this.Fail()
	}

	client.Disconnect(context.Background())

	this.Success()
}
