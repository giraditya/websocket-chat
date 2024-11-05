package handler

import (
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
	MasterAgent agent.MasterAgentInterface
	SupporAgent agent.SupportAgentInterface
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
	identifier := c.Query("identifier")
	client := c.Query("client")

	conn, err := WebsocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.WithContext(c).Errorf("Error upgrading connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer conn.Close()

	wsConnection := models.WebsocketConnection{
		Type: constants.USER_AGENT_WS,
		Conn: conn,
	}
	h.MasterAgent.SaveConnection(c, identifier, wsConnection)

	starterContent, err := helpers.ReadHTMLFile("starting-chat.html")
	if err != nil {
		log.WithContext(c).Errorf("Error reading HTML file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading HTML file"})
	}

	msg := models.Message{
		Sender:      "WsMaster",
		Recipient:   client,
		Identifier:  identifier,
		Content:     starterContent,
		Timestamp:   time.Now(),
		MessageType: "html",
	}
	h.MasterAgent.NotifyUser(msg, h.MasterAgent.GetUserConnection(c, identifier))

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.WithContext(c).Warnf("Error reading message: %v", err)
			break
		}
		log.WithContext(c).Infof("Message from user %v: %s", msg.Identifier, msg.Content)

		switch msg.Content {
		case constants.SIGNAL_NEED_SUPPORT:
			msgToSupportAgent := models.Message{
				Sender:      msg.Sender,
				Recipient:   "Support Agent",
				Identifier:  msg.Identifier,
				Content:     "Hi, i need your support, can you help me?",
				Timestamp:   time.Now(),
				MessageType: "text",
			}
			h.MasterAgent.NotifyAllSupportAgent(msgToSupportAgent, h.MasterAgent.GetSupportAgentConnections(c))
		default:
			if h.MasterAgent.IsBondedConnectionExistAndActive(c, msg.Identifier, &models.WebsocketConnection{}, &wsConnection) {
				err := h.MasterAgent.ForwardMessage(c, msg, h.MasterAgent.GetBondedConnection(c, msg.Identifier), constants.USER_AGENT_WS)
				if err != nil {
					log.WithContext(c).Errorf("Error forwarding message: %v", err)
				}
			}
		}
	}
}

func (h *Handler) SupportAgentWs(c *gin.Context) {
	identifier := c.Query("identifier")

	conn, err := WebsocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.WithContext(c).Errorf("Error upgrading connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer conn.Close()

	wsConnection := models.WebsocketConnection{
		Type: constants.SUPPORT_AGENT_WS,
		Conn: conn,
	}
	h.MasterAgent.SaveConnection(c, identifier, wsConnection)

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.WithContext(c).Warnf("Error reading message: %v", err)
			break
		}

		switch msg.Content {
		case constants.SIGNAL_TAKE_SESSION:
			userConn := h.MasterAgent.GetUserConnection(c, msg.Recipient)
			if !helpers.IsStructEmpty(userConn) {
				if !h.MasterAgent.IsBondedConnectionExistAndActive(c, msg.Recipient, &wsConnection, &models.WebsocketConnection{}) {
					h.MasterAgent.SaveBondedConnection(c, msg.Recipient, models.BondedConnection{
						ConnUser:    &userConn,
						ConnSupport: &wsConnection,
						ChatID:      msg.Recipient,
					})
				}
			}
		case constants.SIGNAL_END_SESSION:
			// DO SOMETHING HERE
		case constants.SIGNAL_BANNED:
			// DO SOMETHING HERE
		case constants.SIGNAL_MOVE_SESSION:
			// DO SOMETHING HERE
		default:
			if h.MasterAgent.IsBondedConnectionExistAndActive(c, msg.Recipient, &wsConnection, &models.WebsocketConnection{}) {
				err := h.MasterAgent.ForwardMessage(c, msg, h.MasterAgent.GetBondedConnection(c, msg.Recipient), constants.SUPPORT_AGENT_WS)
				if err != nil {
					log.WithContext(c).Errorf("Error forwarding message: %v", err)
				}
			}
		}
		log.WithContext(c).Infof("Message from support %s: %s", msg.Sender, msg.Content)
	}
}
