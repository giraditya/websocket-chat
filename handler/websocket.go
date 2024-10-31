package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"websocket-chat/agent"
	"websocket-chat/constants"
	"websocket-chat/helpers"
	"websocket-chat/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Connections       []models.WebsocketConnection
	BondedConnections []models.BondedConnection
	MasterAgent       agent.MasterAgentInterface
	SupporAgent       agent.SupportAgentInterface
}

type HandlerInterface interface {
	ClientWs(c *gin.Context)
	SupportAgentWs(c *gin.Context)
}

func NewHandler(masterAgent agent.MasterAgentInterface, supportAgent agent.SupportAgentInterface) HandlerInterface {
	return &Handler{
		MasterAgent: masterAgent,
		SupporAgent: supportAgent,
	}
}

var WebsocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) ClientWs(c *gin.Context) {
	username := c.Query("username")
	conn, err := WebsocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.WithContext(c).Errorf("Error upgrading connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer conn.Close()

	wsConnection := models.WebsocketConnection{
		ID:   username,
		Type: constants.USER_AGENT_WS,
		Conn: conn,
	}
	h.SaveConnection(c, wsConnection)

	msg := models.Message{
		Username:    "WsMaster",
		Content:     fmt.Sprintf("Hello %s welcome,  you are now connected to the chat", username),
		Timestamp:   time.Now(),
		MessageType: "text",
	}
	h.MasterAgent.NotifyUser(msg, h.GetUserConnection(c, username))

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.WithContext(c).Warnf("Error reading message: %v", err)
			break
		}
		log.WithContext(c).Infof("Message from %s: %s", username, msg.Content)

		switch msg.Content {
		case "NEED SUPPORT":
			msgToSupportAgent := models.Message{
				Username:    msg.Username,
				Content:     "Hi, i need your support, can you help me?",
				Timestamp:   time.Now(),
				MessageType: "text",
			}
			h.SupporAgent.NotifyAllSupportAgent(msgToSupportAgent, h.GetSupportAgentConnections(c))
		default:
			if h.IsBondedConnectionExist(c, msg.Username) {
				err := h.MasterAgent.ForwardMessage(c, msg, h.GetBondedConnection(c, msg.Username), constants.USER_AGENT_WS)
				if err != nil {
					log.WithContext(c).Errorf("Error forwarding message: %v", err)
				}
			}
		}
	}
}

func (h *Handler) SupportAgentWs(c *gin.Context) {
	username := c.Query("username")
	conn, err := WebsocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.WithContext(c).Errorf("Error upgrading connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer conn.Close()

	wsConnection := models.WebsocketConnection{
		ID:   username,
		Type: constants.SUPPORT_AGENT_WS,
		Conn: conn,
	}
	h.SaveConnection(c, wsConnection)

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.WithContext(c).Warnf("Error reading message: %v", err)
			break
		}

		switch msg.Content {
		case "BONDED CONNECTION":
			userConn := h.GetUserConnection(c, "John Doe")
			if !helpers.IsStructEmpty(userConn) {
				h.SaveBondedConnection(c, models.BondedConnection{
					ConnUser:    &userConn,
					ConnSupport: &wsConnection,
					ID:          userConn.ID,
				})
			}
			log.WithContext(c).Info("CONNECTION BONDED")
		default:
			if h.IsBondedConnectionExist(c, "John Doe") {
				err := h.MasterAgent.ForwardMessage(c, msg, h.GetBondedConnection(c, "John Doe"), constants.SUPPORT_AGENT_WS)
				if err != nil {
					log.WithContext(c).Errorf("Error forwarding message: %v", err)
				}
			}
		}
		log.WithContext(c).Infof("Message from %s: %s", "John Doe", msg.Content)
	}
}

func (h *Handler) SaveConnection(c context.Context, identifier models.WebsocketConnection) {
	h.Connections = append(h.Connections, identifier)
}

func (h *Handler) GetConnections(c context.Context) []models.WebsocketConnection {
	return h.Connections
}

func (h *Handler) GetSupportAgentConnections(c context.Context) []models.WebsocketConnection {
	var supportConnections []models.WebsocketConnection
	for _, conn := range h.Connections {
		if conn.Type == constants.SUPPORT_AGENT_WS {
			supportConnections = append(supportConnections, conn)
		}
	}
	return supportConnections
}

func (h *Handler) GetUserConnections(c context.Context) []models.WebsocketConnection {
	var userConnections []models.WebsocketConnection
	for _, conn := range h.Connections {
		if conn.Type == constants.USER_AGENT_WS {
			userConnections = append(userConnections, conn)
		}
	}
	return userConnections
}

func (h *Handler) GetUserConnection(c context.Context, id string) models.WebsocketConnection {
	var userConnections models.WebsocketConnection
	for _, conn := range h.Connections {
		fmt.Println(conn)
		if conn.Type == constants.USER_AGENT_WS && conn.ID == id {
			userConnections = conn
			break
		}
	}
	return userConnections
}

func (h *Handler) SaveBondedConnection(c context.Context, bondedConn models.BondedConnection) {
	h.BondedConnections = append(h.BondedConnections, bondedConn)
}

func (h *Handler) GetBondedConnection(c context.Context, id string) models.BondedConnection {
	var conn models.BondedConnection
	for _, v := range h.BondedConnections {
		if v.ID == id {
			conn = v
			break
		}
	}
	return conn
}

func (h *Handler) IsBondedConnectionExist(c context.Context, identifier string) bool {
	for _, v := range h.BondedConnections {
		if v.ID == identifier {
			return true
		}
	}
	return false
}
