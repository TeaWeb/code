package teamongo

import (
	"context"
	"errors"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var sharedClient *mongo.Client

func RestartClient() {
	sharedClient = nil
}

// 获取共享的Client
func SharedClient() *mongo.Client {
	if sharedClient == nil {
		configFile := files.NewFile(Tea.ConfigFile("mongo.conf"))
		if !configFile.Exists() {
			logs.Fatal(errors.New("'mongo.conf' not found"))
			return nil
		}
		reader, err := configFile.Reader()
		if err != nil {
			logs.Fatal(err)
			return nil
		}

		config := &Config{}
		err = reader.ReadYAML(config)
		if err != nil {
			logs.Fatal(err)
			return nil
		}

		sharedClient, err = mongo.NewClient(options.Client().ApplyURI(config.URI))
		if err != nil {
			logs.Fatal(err)
			return nil
		}

		err = sharedClient.Connect(context.Background())
		if err != nil {
			logs.Fatal(err)
			return nil
		}
	}

	return sharedClient
}

// 获取新Client
func NewClient() *mongo.Client {
	configFile := files.NewFile(Tea.ConfigFile("mongo.conf"))
	if !configFile.Exists() {
		logs.Fatal(errors.New("'mongo.conf' not found"))
		return nil
	}
	reader, err := configFile.Reader()
	if err != nil {
		logs.Fatal(err)
		return nil
	}

	config := &Config{}
	err = reader.ReadYAML(config)
	if err != nil {
		logs.Fatal(err)
		return nil
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(config.URI))
	if err != nil {
		logs.Fatal(err)
		return nil
	}

	err = client.Connect(context.Background())
	if err != nil {
		logs.Fatal(err)
		return nil
	}
	return client
}

// 测试连接
func Test() error {
	configFile := files.NewFile(Tea.ConfigFile("mongo.conf"))
	if !configFile.Exists() {
		return errors.New("'mongo.conf' not found")
	}
	reader, err := configFile.Reader()
	if err != nil {
		return err
	}

	config := &Config{}
	err = reader.ReadYAML(config)
	if err != nil {
		return err
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(config.URI))
	if err != nil {
		return err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err = client.Database("teaweb").Collection("logs").Find(ctx, map[string]interface{}{}, options.Find().SetLimit(1))

	if err == nil {
		client.Disconnect(context.Background())
	}

	return err
}
