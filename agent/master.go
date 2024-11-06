package agent

import (
	"context"
	"fmt"
	"websocket-chat/constants"
	"websocket-chat/helpers"
	"websocket-chat/models"
	"websocket-chat/repository"

	log "github.com/sirupsen/logrus"
)

type MasterAgent struct {
	Connections       map[string]models.WebsocketConnection
	BondedConnections map[string]models.BondedConnection
	Log               *log.Logger
	Repo              repository.RepositoryInterface
}

type MasterAgentInterface interface {
	NotifyUser(msg models.Message, wsConn models.WebsocketConnection) error
	NotifyAllSupportAgent(c context.Context, msg models.Message, wsCons []models.WebsocketConnection) error
	ForwardMessage(c context.Context, msg models.Message, bondedConn *models.BondedConnection, from string) error
	SaveConnection(c context.Context, identifier string, connection models.WebsocketConnection)
	GetConnections(c context.Context) map[string]models.WebsocketConnection
	GetSupportAgentConnections(c context.Context) []models.WebsocketConnection
	GetUserConnections(c context.Context) []models.WebsocketConnection
	GetUserConnection(c context.Context, identifier string) models.WebsocketConnection
	SaveBondedConnection(c context.Context, identifier string, bondedConn models.BondedConnection)
	GetBondedConnection(c context.Context, identifier string) *models.BondedConnection
	IsBondedConnectionExistAndActive(c context.Context, identifier string, connSupportActive *models.WebsocketConnection, connUserActive *models.WebsocketConnection) bool
	RemoveConnection(c context.Context, identifier string)
	RemoveBondedConnection(c context.Context, identifier string)
}

func NewMasterAgent(log *log.Logger, repo repository.RepositoryInterface) MasterAgentInterface {
	return &MasterAgent{
		Connections:       make(map[string]models.WebsocketConnection),
		BondedConnections: make(map[string]models.BondedConnection),
		Log:               log,
		Repo:              repo,
	}
}

func (m *MasterAgent) NotifyUser(msg models.Message, wsConn models.WebsocketConnection) error {
	err := wsConn.Conn.WriteJSON(msg)
	if err != nil {
		return fmt.Errorf("error notifying user: %v", err)
	}
	return nil
}

func (m *MasterAgent) NotifyAllSupportAgent(c context.Context, msg models.Message, wsCons []models.WebsocketConnection) error {
	for _, wsCon := range wsCons {
		err := wsCon.Conn.WriteJSON(msg)
		if err != nil {
			return fmt.Errorf("error send message to support agent: %v", err)
		}
	}
	go m.Repo.InsertLogClientNeedSupport(c, msg.Sender)

	return nil
}

func (m *MasterAgent) ForwardMessage(c context.Context, msg models.Message, bondedConn *models.BondedConnection, from string) error {
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

	go m.Repo.InsertMessage(c, msg)

	return nil
}

func (m *MasterAgent) SaveConnection(c context.Context, identifier string, connection models.WebsocketConnection) {
	m.RemoveConnection(c, identifier)
	m.Connections[identifier] = connection
}

func (m *MasterAgent) GetConnections(c context.Context) map[string]models.WebsocketConnection {
	return m.Connections
}

func (m *MasterAgent) GetSupportAgentConnections(c context.Context) []models.WebsocketConnection {
	var supportConnections []models.WebsocketConnection
	for _, conn := range m.Connections {
		if conn.Type == constants.SUPPORT_AGENT_WS {
			supportConnections = append(supportConnections, conn)
		}
	}
	return supportConnections
}

func (m *MasterAgent) GetUserConnections(c context.Context) []models.WebsocketConnection {
	var userConnections []models.WebsocketConnection
	for _, conn := range m.Connections {
		if conn.Type == constants.USER_AGENT_WS {
			userConnections = append(userConnections, conn)
		}
	}
	return userConnections
}

func (m *MasterAgent) GetUserConnection(c context.Context, identifier string) models.WebsocketConnection {
	return m.Connections[identifier]
}

func (m *MasterAgent) SaveBondedConnection(c context.Context, identifier string, bondedConn models.BondedConnection) {
	m.BondedConnections[identifier] = bondedConn
}

func (m *MasterAgent) GetBondedConnection(c context.Context, identifier string) *models.BondedConnection {
	conn := m.BondedConnections[identifier]
	return &conn
}

func (m *MasterAgent) IsBondedConnectionExistAndActive(c context.Context, identifier string, connSupportActive *models.WebsocketConnection, connUserActive *models.WebsocketConnection) bool {
	conn := m.BondedConnections[identifier]

	if helpers.IsStructEmpty(conn) {
		return false
	}

	if !helpers.IsStructEmpty(connSupportActive) {
		conn.ConnSupport = connSupportActive
		m.BondedConnections[identifier] = conn
		return true
	}

	if !helpers.IsStructEmpty(connUserActive) {
		conn.ConnUser = connUserActive
		m.BondedConnections[identifier] = conn
		return true
	}
	return false
}

func (m *MasterAgent) RemoveConnection(c context.Context, identifier string) {
	conn := m.Connections[identifier]
	if !helpers.IsStructEmpty(conn) {
		delete(m.Connections, identifier)
	}
}

func (m *MasterAgent) RemoveBondedConnection(c context.Context, identifier string) {
	delete(m.BondedConnections, identifier)
}
