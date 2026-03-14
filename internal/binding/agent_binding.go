package binding

import (
	"skillmanager/internal/model"
	"skillmanager/internal/service"
)

type AgentBinding struct {
	service *service.AgentService
}

func NewAgentBinding(svc *service.AgentService) *AgentBinding {
	return &AgentBinding{service: svc}
}

func (b *AgentBinding) ListAgents() ([]model.Agent, error) {
	return b.service.ListAgents()
}

func (b *AgentBinding) DetectInstalled() ([]model.Agent, error) {
	return b.service.DetectInstalled()
}

func (b *AgentBinding) AddCustomAgent(agent model.Agent) error {
	return b.service.AddCustomAgent(agent)
}

func (b *AgentBinding) RemoveAgent(id string) error {
	return b.service.RemoveAgent(id)
}

func (b *AgentBinding) ToggleAgent(id string, enabled bool) error {
	return b.service.ToggleAgent(id, enabled)
}
