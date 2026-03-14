package binding

import (
	"skillmanager/internal/model"
	"skillmanager/internal/service"
)

type SkillBinding struct {
	service *service.SkillService
}

func NewSkillBinding(svc *service.SkillService) *SkillBinding {
	return &SkillBinding{service: svc}
}

func (b *SkillBinding) ListInstalled() ([]model.Skill, error) {
	return b.service.ListInstalled()
}

func (b *SkillBinding) GetDetail(id string) (*model.Skill, error) {
	return b.service.GetDetail(id)
}

func (b *SkillBinding) Install(sourceURL string, agents []string) error {
	return b.service.Install(sourceURL, agents)
}

func (b *SkillBinding) Uninstall(id string) error {
	return b.service.Uninstall(id)
}

func (b *SkillBinding) Update(id string) (*model.Skill, error) {
	return b.service.Update(id)
}

func (b *SkillBinding) UpdateContent(id string, content string) error {
	return b.service.UpdateContent(id, content)
}

func (b *SkillBinding) AssignAgents(id string, agents []string) error {
	return b.service.AssignAgents(id, agents)
}
