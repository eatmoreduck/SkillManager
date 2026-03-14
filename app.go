package main

import (
	"skillmanager/internal/binding"
	"skillmanager/internal/repository"
	"skillmanager/internal/service"
)

type App struct {
	ConfigPath       string
	SkillBinding     *binding.SkillBinding
	RegistryBinding  *binding.RegistryBinding
	AgentBinding     *binding.AgentBinding
	ConfigBinding    *binding.ConfigBinding
	InventoryBinding *binding.InventoryBinding
}

func NewApp(configPath string) (*App, error) {
	configRepo, err := repository.NewFileConfigRepository(configPath)
	if err != nil {
		return nil, err
	}

	configService := service.NewConfigService(configRepo)
	agentService := service.NewAgentService(configRepo)
	skillRepo := repository.NewFileSkillRepository()
	skillService := service.NewSkillService(skillRepo, agentService, configService)
	registryService := service.NewRegistryService(repository.NewHTTPRegistryRepository(), configRepo)
	inventoryService := service.NewInventoryService(configRepo, skillRepo, agentService, skillService)

	return &App{
		ConfigPath:       configRepo.GetConfigPath(),
		SkillBinding:     binding.NewSkillBinding(skillService),
		RegistryBinding:  binding.NewRegistryBinding(registryService),
		AgentBinding:     binding.NewAgentBinding(agentService),
		ConfigBinding:    binding.NewConfigBinding(configService),
		InventoryBinding: binding.NewInventoryBinding(inventoryService),
	}, nil
}
