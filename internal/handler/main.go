package handler

import (
	"websocket-chat/internal/agent"
	"websocket-chat/internal/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	MasterAgent agent.MasterAgentInterface
	SupporAgent agent.SupportAgentInterface
	Repo        repository.RepositoryInterface
}

type HandlerInterface interface {
	ClientWs(c *gin.Context)
	SupportAgentWs(c *gin.Context)
}

func NewHandler(masterAgent agent.MasterAgentInterface, supportAgent agent.SupportAgentInterface, repo repository.RepositoryInterface) HandlerInterface {
	return &Handler{
		MasterAgent: masterAgent,
		SupporAgent: supportAgent,
		Repo:        repo,
	}
}
