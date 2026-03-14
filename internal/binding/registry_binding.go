package binding

import (
	"skillmanager/internal/model"
	"skillmanager/internal/service"
)

type RegistryBinding struct {
	service *service.RegistryService
}

func NewRegistryBinding(svc *service.RegistryService) *RegistryBinding {
	return &RegistryBinding{service: svc}
}

func (b *RegistryBinding) ListRegistries() ([]model.Registry, error) {
	return b.service.ListRegistries()
}

func (b *RegistryBinding) AddRegistry(registry model.Registry) error {
	return b.service.AddRegistry(registry)
}

func (b *RegistryBinding) RemoveRegistry(id string) error {
	return b.service.RemoveRegistry(id)
}

func (b *RegistryBinding) Browse(registryID string, category string) ([]model.RegistrySkill, error) {
	return b.service.Browse(registryID, category)
}

func (b *RegistryBinding) Search(query string) ([]model.RegistrySkill, error) {
	return b.service.Search(query)
}
