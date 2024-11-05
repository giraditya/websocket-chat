package agent

type SupportAgent struct{}

type SupportAgentInterface interface{}

func NewSupportAgent() SupportAgentInterface {
	return &SupportAgent{}
}
