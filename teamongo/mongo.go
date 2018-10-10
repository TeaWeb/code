package teamongo

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"context"
	"time"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"errors"
)

var sharedClient *mongo.Client

func RestartClient() {
	sharedClient = nil
}

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

		sharedClient, err = mongo.NewClient(config.URI)
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

	client, err := mongo.NewClient(config.URI)
	if err != nil {
		return err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err = client.Database("teaweb").Collection("logs").Find(ctx, nil, findopt.Limit(1))

	if err == nil {
		client.Disconnect(context.Background())
	}

	return err
}
