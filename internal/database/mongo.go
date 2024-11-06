package database

import (
	"context"

	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	Log    *log.Logger
}

type MongoDBInterface interface {
	ConnectToMongoDB()
	GetMongoClient() *mongo.Client
}

func NewMongoDB(log *log.Logger) MongoDBInterface {
	return &MongoDB{
		Log: log,
	}
}

func (d *MongoDB) ConnectToMongoDB() {
	uri := "mongodb://localhost:27017"

	// Create a new MongoDB client
	var err error
	d.Client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = d.Client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	log.Println("Connected to MongoDB!")
}

func (d *MongoDB) GetMongoClient() *mongo.Client {
	return d.Client
}
