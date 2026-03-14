package service

import (
	"log"
	"strings"

	"skillmanager/internal/model"
	"skillmanager/internal/repository"
)

// RegistryService handles registry browsing and configuration.
type RegistryService struct {
	registryRepo  repository.RegistryRepository
	configService *ConfigService
	configRepo    repository.ConfigRepository
}

func NewRegistryService(registryRepo repository.RegistryRepository, configRepo repository.ConfigRepository) *RegistryService {
	return &RegistryService{
		registryRepo:  registryRepo,
		configService: NewConfigService(configRepo),
		configRepo:    configRepo,
	}
}

func (s *RegistryService) ListRegistries() ([]model.Registry, error) {
	cfg, err := s.configRepo.LoadConfig()
	if err != nil {
		return nil, err
	}
	return append([]model.Registry(nil), cfg.Registries...), nil
}

func (s *RegistryService) AddRegistry(registry model.Registry) error {
	cfg, err := s.configRepo.LoadConfig()
	if err != nil {
		return err
	}
	for _, existing := range cfg.Registries {
		if existing.ID == registry.ID {
			return model.ErrAlreadyExists
		}
	}
	cfg.Registries = append(cfg.Registries, registry)
	return s.configRepo.SaveConfig(cfg)
}

func (s *RegistryService) RemoveRegistry(id string) error {
	cfg, err := s.configRepo.LoadConfig()
	if err != nil {
		return err
	}

	next := make([]model.Registry, 0, len(cfg.Registries))
	found := false
	for _, registry := range cfg.Registries {
		if registry.ID == id {
			found = true
			continue
		}
		next = append(next, registry)
	}
	if !found {
		return model.ErrNotFound
	}

	cfg.Registries = next
	return s.configRepo.SaveConfig(cfg)
}

func (s *RegistryService) Browse(registryID string, category string) ([]model.RegistrySkill, error) {
	registry, err := s.registryByID(registryID)
	if err != nil {
		return nil, err
	}

	client := s.configService.GetHTTPClient()
	rawURL, err := repository.BuildRegistryBrowseURL(registry.URL, category)
	if err != nil {
		return nil, err
	}
	log.Printf("[registry] browse registry=%s url=%s category=%q", registry.ID, rawURL, category)

	return s.registryRepo.FetchRegistry(rawURL, client)
}

func (s *RegistryService) Search(query string) ([]model.RegistrySkill, error) {
	registries, err := s.ListRegistries()
	if err != nil {
		return nil, err
	}
	if len(registries) == 0 {
		return nil, model.ErrNotFound
	}

	client := s.configService.GetHTTPClient()

	var results []model.RegistrySkill
	for _, registry := range registries {
		rawURL, err := repository.BuildRegistrySearchURL(registry.URL, query)
		if err != nil {
			return nil, err
		}
		log.Printf("[registry] search registry=%s url=%s query=%q", registry.ID, rawURL, query)
		items, err := s.registryRepo.FetchRegistry(rawURL, client)
		if err != nil {
			log.Printf("[registry] search registry=%s failed err=%v", registry.ID, err)
			continue
		}
		results = append(results, items...)
	}

	return dedupeRegistrySkills(results), nil
}

func (s *RegistryService) registryByID(registryID string) (*model.Registry, error) {
	registries, err := s.ListRegistries()
	if err != nil {
		return nil, err
	}

	if registryID == "" {
		for _, registry := range registries {
			if registry.IsDefault {
				return &registry, nil
			}
		}
		if len(registries) > 0 {
			return &registries[0], nil
		}
		return nil, model.ErrNotFound
	}

	for _, registry := range registries {
		if strings.EqualFold(registry.ID, registryID) {
			return &registry, nil
		}
	}

	return nil, model.ErrNotFound
}

func dedupeRegistrySkills(skills []model.RegistrySkill) []model.RegistrySkill {
	seen := make(map[string]struct{}, len(skills))
	result := make([]model.RegistrySkill, 0, len(skills))
	for _, item := range skills {
		if _, ok := seen[item.ID]; ok {
			continue
		}
		seen[item.ID] = struct{}{}
		result = append(result, item)
	}
	return result
}
