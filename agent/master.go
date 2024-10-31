package agent

import (
	"context"
	"fmt"
	"websocket-chat/constants"
	"websocket-chat/models"
)

type MasterAgent struct{}

type MasterAgentInterface interface {
	NotifyUser(msg models.Message, wsConn models.WebsocketConnection) error
	ForwardMessage(c context.Context, msg models.Message, bondedConn models.BondedConnection, from string) error
}

func NewMasterAgent() MasterAgentInterface {
	return &MasterAgent{}
}

func (m *MasterAgent) NotifyUser(msg models.Message, wsConn models.WebsocketConnection) error {
	err := wsConn.Conn.WriteJSON(msg)
	if err != nil {
		return fmt.Errorf("error notifying user: %v", err)
	}
	return nil
}

func (m *MasterAgent) ForwardMessage(c context.Context, msg models.Message, bondedConn models.BondedConnection, from string) error {
	if from == constants.USER_AGENT_WS {
		err := bondedConn.ConnSupport.Conn.WriteJSON(msg)
		if err != nil {
			return fmt.Errorf("error notifying support: %v", err)
		}
	} else {
		err := bondedConn.ConnUser.Conn.WriteJSON(msg)
		if err != nil {
			return fmt.Errorf("error notifying user: %v", err)
		}
	}

	return nil
}
