package repository

import (
	"context"
	"websocket-chat/internal/models"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	DB  *mongo.Client
	Log *log.Logger
}

type RepositoryInterface interface {
	InsertMessage(c context.Context, msg models.Message) error
	GetChatHistory(sender, recipient string) ([]models.Message, error)
	InsertLogClientNeedSupport(c context.Context, client string) error
}

func NewRepository(db *mongo.Client, log *log.Logger) RepositoryInterface {
	return &Repository{
		DB:  db,
		Log: log,
	}
}
