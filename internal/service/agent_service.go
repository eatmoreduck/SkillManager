package service

import (
	"log"
	"os/exec"

	"skillmanager/internal/model"
	"skillmanager/internal/pkg/paths"
	"skillmanager/internal/repository"
)

// AgentService manages configured agents and installation detection.
type AgentService struct {
	configRepo repository.ConfigRepository
}

func NewAgentService(configRepo repository.ConfigRepository) *AgentService {
	return &AgentService{configRepo: configRepo}
}

func (s *AgentService) ListAgents() ([]model.Agent, error) {
	cfg, err := s.configRepo.LoadConfig()
	if err != nil {
		return nil, err
	}

	agents := make([]model.Agent, 0, len(cfg.Agents))
	for _, item := range cfg.Agents {
		agent := item.ToAgent()
		if agent.SkillsDir, err = paths.Expand(agent.SkillsDir); err != nil {
			return nil, err
		}
		for i, p := range agent.PriorityPaths {
			agent.PriorityPaths[i], err = paths.Expand(p)
			if err != nil {
				return nil, err
			}
		}
		agent.IsInstalled = isBinaryAvailable(agent.BinaryName)
		agents = append(agents, agent)
	}

	log.Printf("[agent] loaded %d configured agents", len(agents))
	return agents, nil
}

func (s *AgentService) DetectInstalled() ([]model.Agent, error) {
	agents, err := s.ListAgents()
	if err != nil {
		return nil, err
	}

	filtered := make([]model.Agent, 0, len(agents))
	for _, agent := range agents {
		if agent.IsInstalled {
			filtered = append(filtered, agent)
		}
	}

	log.Printf("[agent] detected %d installed agents", len(filtered))
	return filtered, nil
}

func (s *AgentService) AddCustomAgent(agent model.Agent) error {
	cfg, err := s.configRepo.LoadConfig()
	if err != nil {
		return err
	}

	for _, existing := range cfg.Agents {
		if existing.ID == agent.ID {
			return model.ErrAlreadyExists
		}
	}

	agent.IsCustom = true
	cfg.Agents = append(cfg.Agents, agent.ToConfig())
	return s.configRepo.SaveConfig(cfg)
}

func (s *AgentService) RemoveAgent(id string) error {
	cfg, err := s.configRepo.LoadConfig()
	if err != nil {
		return err
	}

	next := make([]model.AgentConfig, 0, len(cfg.Agents))
	found := false
	for _, agent := range cfg.Agents {
		if agent.ID == id {
			found = true
			continue
		}
		next = append(next, agent)
	}
	if !found {
		return model.ErrNotFound
	}

	cfg.Agents = next
	return s.configRepo.SaveConfig(cfg)
}

func (s *AgentService) ToggleAgent(id string, enabled bool) error {
	cfg, err := s.configRepo.LoadConfig()
	if err != nil {
		return err
	}

	for i := range cfg.Agents {
		if cfg.Agents[i].ID == id {
			cfg.Agents[i].IsEnabled = enabled
			return s.configRepo.SaveConfig(cfg)
		}
	}

	return model.ErrNotFound
}

func isBinaryAvailable(binary string) bool {
	if binary == "" {
		return false
	}
	_, err := exec.LookPath(binary)
	return err == nil
}
