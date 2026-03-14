package service

import (
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"skillmanager/internal/model"
	"skillmanager/internal/pkg/paths"
	"skillmanager/internal/repository"
)

// SkillService manages installed local skills.
type SkillService struct {
	skillRepo     repository.SkillRepository
	agentService  *AgentService
	configService *ConfigService
}

func NewSkillService(
	skillRepo repository.SkillRepository,
	agentService *AgentService,
	configService *ConfigService,
) *SkillService {
	return &SkillService{
		skillRepo:     skillRepo,
		agentService:  agentService,
		configService: configService,
	}
}

func (s *SkillService) ListInstalled() ([]model.Skill, error) {
	agents, err := s.agentService.ListAgents()
	if err != nil {
		return nil, err
	}

	skills, err := s.skillRepo.ScanSkills(agents)
	if err != nil {
		return nil, err
	}

	log.Printf("[skill] ListInstalled discovered %d skills across %d configured agents", len(skills), len(agents))
	return skills, nil
}

func (s *SkillService) GetDetail(id string) (*model.Skill, error) {
	skills, err := s.ListInstalled()
	if err != nil {
		return nil, err
	}

	for _, skill := range skills {
		if skill.ID == id {
			detail, err := s.skillRepo.ReadSkill(skill.LocalPath)
			if err != nil {
				return nil, err
			}
			detail.Agents = append([]string{}, skill.Agents...)
			detail.SourceURL = skill.SourceURL
			detail.Version = skill.Version
			detail.Tags = append([]string{}, skill.Tags...)
			return detail, nil
		}
	}

	return nil, model.ErrNotFound
}

func (s *SkillService) Install(sourceURL string, agents []string) error {
	if len(agents) == 0 {
		return model.ErrSkillNotAssigned
	}

	enabledAgents, err := s.enabledAgentsByID()
	if err != nil {
		return err
	}

	for _, agentID := range agents {
		if _, ok := enabledAgents[agentID]; !ok {
			return model.ErrAgentNotEnabled
		}
	}

	skillsDir, err := paths.GetSkillsDir()
	if err != nil {
		return err
	}
	if err := paths.EnsureDir(skillsDir); err != nil {
		return err
	}

	proxy, err := s.proxyConfig()
	if err != nil {
		return err
	}

	targetPath := filepath.Join(skillsDir, baseNameFromSource(sourceURL))
	if err := s.skillRepo.CloneSkill(sourceURL, targetPath, proxy); err != nil {
		return err
	}

	for _, agentID := range agents {
		agent := enabledAgents[agentID]
		if err := s.skillRepo.CreateSymlink(targetPath, agent.SkillsDir); err != nil {
			return err
		}
	}

	return nil
}

func (s *SkillService) Uninstall(id string) error {
	skills, err := s.ListInstalled()
	if err != nil {
		return err
	}

	enabledAgents, err := s.enabledAgentsByID()
	if err != nil {
		return err
	}

	for _, skill := range skills {
		if skill.ID != id {
			continue
		}

		for _, agentID := range skill.Agents {
			agent, ok := enabledAgents[agentID]
			if !ok {
				continue
			}
			if err := s.skillRepo.RemoveSymlink(agent.SkillsDir, filepath.Base(skill.LocalPath)); err != nil {
				return err
			}
		}

		return s.skillRepo.DeleteSkill(skill.LocalPath)
	}

	return model.ErrNotFound
}

func (s *SkillService) Update(id string) (*model.Skill, error) {
	skill, err := s.GetDetail(id)
	if err != nil {
		return nil, err
	}

	proxy, err := s.proxyConfig()
	if err != nil {
		return nil, err
	}

	if err := s.skillRepo.PullSkill(skill.LocalPath, proxy); err != nil {
		return nil, err
	}

	return s.skillRepo.ReadSkill(skill.LocalPath)
}

func (s *SkillService) UpdateContent(id string, content string) error {
	skill, err := s.GetDetail(id)
	if err != nil {
		return err
	}

	skill.Content = content
	skill.UpdatedAt = time.Now()

	return s.skillRepo.WriteSkill(skill.LocalPath, skill)
}

func (s *SkillService) AssignAgents(id string, agents []string) error {
	skill, err := s.GetDetail(id)
	if err != nil {
		return err
	}

	enabledAgents, err := s.enabledAgentsByID()
	if err != nil {
		return err
	}

	targets := make(map[string]struct{}, len(agents))
	for _, agentID := range agents {
		agent, ok := enabledAgents[agentID]
		if !ok {
			return model.ErrAgentNotEnabled
		}
		targets[agentID] = struct{}{}
		if err := s.skillRepo.CreateSymlink(skill.LocalPath, agent.SkillsDir); err != nil {
			return err
		}
	}

	for _, current := range skill.Agents {
		if _, keep := targets[current]; keep {
			continue
		}
		if agent, ok := enabledAgents[current]; ok {
			if err := s.skillRepo.RemoveSymlink(agent.SkillsDir, filepath.Base(skill.LocalPath)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SkillService) enabledAgentsByID() (map[string]model.Agent, error) {
	agents, err := s.agentService.ListAgents()
	if err != nil {
		return nil, err
	}

	result := make(map[string]model.Agent, len(agents))
	for _, agent := range agents {
		if agent.IsEnabled {
			result[agent.ID] = agent
		}
	}

	return result, nil
}

func (s *SkillService) proxyConfig() (*model.ProxyConfig, error) {
	cfg, err := s.configService.GetConfig()
	if err != nil {
		return nil, err
	}
	return &cfg.Proxy, nil
}

func baseNameFromSource(sourceURL string) string {
	base := filepath.Base(strings.TrimSuffix(sourceURL, ".git"))
	return strings.TrimSpace(base)
}

func sortStrings(values []string) []string {
	items := append([]string(nil), values...)
	sort.Strings(items)
	return items
}
