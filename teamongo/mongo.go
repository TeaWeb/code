package teamongo

import (
	"context"
	"errors"
	"github.com/TeaWeb/code/teaweb/configs"
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
		defer reader.Close()

		config := &Config{}
		err = reader.ReadYAML(config)
		if err != nil {
			logs.Fatal(err)
			return nil
		}

		clientOptions := options.Client().ApplyURI(config.URI)
		sharedConfig := configs.SharedMongoConfig()

		if sharedConfig != nil && len(sharedConfig.AuthMechanism) > 0 {
			clientOptions.SetAuth(options.Credential{
				Username:                sharedConfig.Username,
				Password:                sharedConfig.Password,
				AuthMechanism:           sharedConfig.AuthMechanism,
				AuthMechanismProperties: sharedConfig.AuthMechanismPropertiesMap(),
				AuthSource:              DatabaseName,
			})
		}

		sharedClient, err = mongo.NewClient(clientOptions)
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
	defer reader.Close()

	config := &Config{}
	err = reader.ReadYAML(config)
	if err != nil {
		logs.Fatal(err)
		return nil
	}

	clientOptions := options.Client().ApplyURI(config.URI)
	sharedConfig := configs.SharedMongoConfig()

	if sharedConfig != nil && len(sharedConfig.AuthMechanism) > 0 {
		clientOptions.SetAuth(options.Credential{
			Username:                sharedConfig.Username,
			Password:                sharedConfig.Password,
			AuthMechanism:           sharedConfig.AuthMechanism,
			AuthMechanismProperties: sharedConfig.AuthMechanismPropertiesMap(),
			AuthSource:              DatabaseName,
		})
	}

	client, err := mongo.NewClient(clientOptions)
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
	defer reader.Close()

	clientOptions := options.Client().ApplyURI(config.URI)
	sharedConfig := configs.SharedMongoConfig()

	if sharedConfig != nil && len(sharedConfig.AuthMechanism) > 0 {
		clientOptions.SetAuth(options.Credential{
			Username:                sharedConfig.Username,
			Password:                sharedConfig.Password,
			AuthMechanism:           sharedConfig.AuthMechanism,
			AuthMechanismProperties: sharedConfig.AuthMechanismPropertiesMap(),
			AuthSource:              DatabaseName,
		})
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err = client.Database(DatabaseName).Collection("logs").Find(ctx, map[string]interface{}{}, options.Find().SetLimit(1))

	if err == nil {
		client.Disconnect(context.Background())
	}

	if err == context.DeadlineExceeded {
		err = errors.New("connection timeout")
	}

	return err
}
