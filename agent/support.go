package agent

import (
	"fmt"
	"websocket-chat/models"
)

type SupportAgentInterface interface {
	SendMessage(msg models.Message, wsCon models.WebsocketConnection) error
	NotifyAllSupportAgent(msg models.Message, wsCons []models.WebsocketConnection) error
}

func NewSupportAgent() SupportAgentInterface {
	return &SupportAgent{}
}

type SupportAgent struct {
	Username    string
	ChatHistory []models.Message
}

func (a *SupportAgent) SendMessage(msg models.Message, wsCon models.WebsocketConnection) error {
	err := wsCon.Conn.WriteJSON(msg)
	if err != nil {
		return fmt.Errorf("error send message to user: %v", err)
	}
	return nil
}

func (a *SupportAgent) NotifyAllSupportAgent(msg models.Message, wsCons []models.WebsocketConnection) error {
	for _, wsCon := range wsCons {
		err := wsCon.Conn.WriteJSON(msg)
		if err != nil {
			return fmt.Errorf("error send message to support agent: %v", err)
		}
	}
	return nil
}
