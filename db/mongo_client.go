package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/56-Secure/common/logs"
	"github.com/56-Secure/common/observability"
)

var mongoClient *mongo.Database

func NewMongoSingletonClient(
	config MongoConfig,
	observabilityEnabled bool,
) {
	if mongoClient != nil {
		return
	}

	mongoUri := fmt.Sprintf("mongodb://%s", config.URI)
	if config.Client != "local" {
		mongoUri = fmt.Sprintf(
			"mongodb+srv://%s:%s@%s",
			config.Username,
			config.Password,
			config.URI,
		)
	}

	clientOptions := options.Client().ApplyURI(
		mongoUri,
	)
	if observabilityEnabled == true {
		observability.MongoObservability(
			clientOptions,
		)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	client, err := mongo.Connect(ctx, clientOptions)

	defer cancel()

	if err != nil {
		logs.GetClient().Error(err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logs.GetClient().Error(err.Error())
	}

	logs.GetClient().Info("Mongo connected successfully to : " + mongoUri)
	mongoClient = client.Database(config.DbName)
}

func GetMongoClient() *mongo.Database {
	return mongoClient
}
