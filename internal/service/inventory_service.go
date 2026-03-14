package service

import (
	"sort"

	"skillmanager/internal/model"
	"skillmanager/internal/pkg/paths"
	"skillmanager/internal/repository"
)

// InventoryService produces a local diagnostic view of configured agents and visible skills.
type InventoryService struct {
	configRepo   repository.ConfigRepository
	skillRepo    repository.SkillRepository
	agentService *AgentService
	skillService *SkillService
}

func NewInventoryService(
	configRepo repository.ConfigRepository,
	skillRepo repository.SkillRepository,
	agentService *AgentService,
	skillService *SkillService,
) *InventoryService {
	return &InventoryService{
		configRepo:   configRepo,
		skillRepo:    skillRepo,
		agentService: agentService,
		skillService: skillService,
	}
}

func (s *InventoryService) BuildReport() (*model.InventoryReport, error) {
	agents, err := s.agentService.ListAgents()
	if err != nil {
		return nil, err
	}

	managerSkillsDir, err := paths.GetSkillsDir()
	if err != nil {
		return nil, err
	}
	managerSkillsDir, err = paths.Expand(managerSkillsDir)
	if err != nil {
		return nil, err
	}

	managedSkills, err := s.skillRepo.ScanManagedSkills(managerSkillsDir)
	if err != nil {
		return nil, err
	}

	linkedSkills, err := s.skillService.ListInstalled()
	if err != nil {
		return nil, err
	}

	agentsBySkillPath := make(map[string][]string, len(linkedSkills))
	for _, skill := range linkedSkills {
		agentsBySkillPath[skill.LocalPath] = append([]string(nil), skill.Agents...)
	}

	for i := range managedSkills {
		managedSkills[i].IsManaged = true
		managedSkills[i].Agents = append([]string(nil), agentsBySkillPath[managedSkills[i].LocalPath]...)
		sort.Strings(managedSkills[i].Agents)
	}

	inventories := make([]model.AgentInventory, 0, len(agents))
	for _, agent := range agents {
		inventory, err := s.skillRepo.InspectAgentSkills(agent, managerSkillsDir)
		if err != nil {
			return nil, err
		}
		inventories = append(inventories, inventory)
	}

	return &model.InventoryReport{
		Platform:         paths.PlatformName(),
		ConfigPath:       s.configRepo.GetConfigPath(),
		ManagerSkillsDir: managerSkillsDir,
		Agents:           agents,
		ManagedSkills:    managedSkills,
		AgentInventories: inventories,
	}, nil
}
