package repository

import (
	"context"
	"fmt"
	"time"
	"websocket-chat/models"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *Repository) InsertMessage(c context.Context, msg models.Message) error {
	collection := d.DB.Database("chat_app").Collection("chat_messages")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}

	log.Println("Message inserted successfully")
	return nil
}

func (d *Repository) InsertLogClientNeedSupport(c context.Context, client string) error {
	collection := d.DB.Database("chat_app").Collection("support_requested")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	supportRequested := models.SupportRequested{
		ClientName:  client,
		Timestamp:   time.Now(),
		IsSupported: false,
	}

	_, err := collection.InsertOne(ctx, supportRequested)
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}

	log.Println("Message inserted successfully")
	return nil
}

func (d *Repository) GetChatHistory(sender, recipient string) ([]models.Message, error) {
	collection := d.DB.Database("chat_app").Collection("chat_messages")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"sender": sender, "recipient": recipient},
			{"sender": recipient, "recipient": sender},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chat history: %v", err)
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	for cursor.Next(ctx) {
		var msg models.Message
		if err := cursor.Decode(&msg); err != nil {
			return nil, fmt.Errorf("error decoding message: %v", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
