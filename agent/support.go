package agent

import (
	"fmt"
	"websocket-chat/models"
)

type SupportAgent struct {
	Username    string
	ChatHistory []models.Message
}

type SupportAgentInterface interface {
	SendMessage(msg models.Message, wsCon models.WebsocketConnection) error
}

func NewSupportAgent() SupportAgentInterface {
	return &SupportAgent{}
}

func (a *SupportAgent) SendMessage(msg models.Message, wsCon models.WebsocketConnection) error {
	err := wsCon.Conn.WriteJSON(msg)
	if err != nil {
		return fmt.Errorf("error send message to user: %v", err)
	}
	return nil
}
