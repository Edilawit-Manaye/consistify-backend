package database

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}
func NewMongoClient() (*MongoClient, error) {
	mongoURI := viper.GetString("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in environment variables")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB!")

	return &MongoClient{
		Client: client,
		DB:     client.Database("consistify_db"), 
	}, nil
}


func (mc *MongoClient) Disconnect(ctx context.Context) error {
	if mc.Client == nil {
		return nil
	}
	log.Println("Disconnecting from MongoDB...")
	return mc.Client.Disconnect(ctx)
}
