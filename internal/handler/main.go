package handler

import (
	"websocket-chat/internal/agent"

	"github.com/gin-gonic/gin"
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
