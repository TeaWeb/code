package teamongo

import (
	"context"
	"errors"
	"github.com/TeaWeb/code/teaconfigs/db"
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

	// 清空collection缓存
	collLocker.Lock()
	collMap = map[string]*Collection{}
	collLocker.Unlock()
}

// 获取共享的Client
func SharedClient() *mongo.Client {
	if sharedClient != nil {
		return sharedClient
	}
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
	defer func() {
		err := reader.Close()
		if err != nil {
			logs.Error(err)
		}
	}()

	config := &db.MongoConfig{}
	err = reader.ReadYAML(config)
	if err != nil {
		logs.Fatal(err)
		return nil
	}

	clientOptions := options.Client().ApplyURI(config.URI)
	clientOptions.SetMaxPoolSize(10).
		SetConnectTimeout(5 * time.Second)
	sharedConfig := db.SharedMongoConfig()

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
	logs.Println("[mongo]create new client")
	if err != nil {
		logs.Fatal(err)
		return nil
	}

	err = sharedClient.Connect(context.Background())
	if err != nil {
		logs.Error(err)
		return nil
	}

	return sharedClient
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
	defer func() {
		err = reader.Close()
		if err != nil {
			logs.Error(err)
		}
	}()

	config := &db.MongoConfig{}
	err = reader.ReadYAML(config)
	if err != nil {
		return err
	}

	clientOptions := options.Client().ApplyURI(config.URI)
	clientOptions.SetMaxPoolSize(1).
		SetConnectTimeout(1 * time.Second)
	sharedConfig := db.SharedMongoConfig()

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

	// 尝试查询
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err = client.Database(DatabaseName).
		Collection("logs").
		Find(ctx, map[string]interface{}{}, options.Find().SetLimit(1))

	// 关闭连接
	ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
	err1 := client.Disconnect(ctx)
	if err1 != nil {
		logs.Error(err1)
	}

	if err == context.DeadlineExceeded {
		err = errors.New("connection timeout")
	}

	return err
}
