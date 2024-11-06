package agent

import (
	"websocket-chat/repository"

	log "github.com/sirupsen/logrus"
)

type SupportAgent struct {
	Log  *log.Logger
	Repo repository.RepositoryInterface
}

type SupportAgentInterface interface{}

func NewSupportAgent(log *log.Logger, repo repository.RepositoryInterface) SupportAgentInterface {
	return &SupportAgent{
		Log:  log,
		Repo: repo,
	}
}
