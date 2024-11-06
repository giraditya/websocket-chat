package models

import (
	"time"
)

type Message struct {
	Sender      string    `json:"sender"`
	Recipient   string    `json:"recipient"`
	Identifier  string    `json:"identifier"`
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
	MessageType string    `json:"message_type"`
}
