package mongodb

import (
	"budgeting-service/config"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Database, error) {
	cfg := config.Load()
	client, err := mongo.Connect(context.Background(), options.Client().
		ApplyURI(cfg.MONGODB_URI))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.MONGODB_NAME), nil
}
