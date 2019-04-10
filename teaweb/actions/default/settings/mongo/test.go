package mongo

import (
	"context"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"regexp"
	"strings"
	"time"
)

type TestAction actions.Action

func (this *TestAction) Run(params struct {
	Host                    string
	Port                    uint
	Username                string
	Password                string
	AuthMechanism           string
	AuthMechanismProperties string
}) {
	oldConfig := configs.SharedMongoConfig()

	config := configs.MongoConnectionConfig{
		Host:          params.Host,
		Port:          params.Port,
		Username:      params.Username,
		Password:      params.Password,
		AuthMechanism: params.AuthMechanism,
	}

	if len(params.AuthMechanismProperties) > 0 {
		properties := regexp.MustCompile("\\s*,\\s*").Split(params.AuthMechanismProperties, -1)
		for _, property := range properties {
			if strings.Contains(property, ":") {
				pieces := strings.Split(property, ":")
				config.AuthMechanismProperties = append(config.AuthMechanismProperties, shared.NewVariable(pieces[0], pieces[1]))
			}
		}
	}

	if len(config.Password) == 0 && oldConfig != nil {
		config.Password = oldConfig.Password
	}

	uri := config.URI()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	clientOptions := options.Client().ApplyURI(uri)
	if len(config.AuthMechanism) > 0 {
		clientOptions.SetAuth(options.Credential{
			Username:                config.Username,
			Password:                config.Password,
			AuthMechanism:           config.AuthMechanism,
			AuthMechanismProperties: config.AuthMechanismPropertiesMap(),
			AuthSource:              teamongo.DatabaseName,
		})
	}
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		this.Message = "[连接]有错误需要修复：" + err.Error()
		this.Fail()
	}

	ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
	_, err = client.
		Database(teamongo.DatabaseName).
		Collection("logs").
		Find(ctx, map[string]interface{}{}, options.Find().SetLimit(1))
	if err != nil {
		this.Message = "[查询]有错误需要修复：" + err.Error()
		this.Fail()
	}

	client.Disconnect(context.Background())

	this.Success()
}
