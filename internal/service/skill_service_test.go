package service

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"skillmanager/internal/model"
	"skillmanager/internal/repository"
)

func TestSkillServiceGetDetailPreservesAssignedAgents(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	managedSkillsDir := filepath.Join(tempDir, "managed-skills")
	agentADir := filepath.Join(tempDir, "agent-a")
	agentBDir := filepath.Join(tempDir, "agent-b")

	for _, dir := range []string{managedSkillsDir, agentADir, agentBDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", dir, err)
		}
	}

	configRepo, err := repository.NewFileConfigRepository(configPath)
	if err != nil {
		t.Fatalf("NewFileConfigRepository() error = %v", err)
	}

	cfg, err := configRepo.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	cfg.Agents = []model.AgentConfig{
		{
			ID:            "agent-a",
			Name:          "Agent A",
			SkillsDir:     agentADir,
			BinaryName:    "sh",
			PriorityPaths: []string{agentADir},
			IsEnabled:     true,
			IsCustom:      false,
		},
		{
			ID:            "agent-b",
			Name:          "Agent B",
			SkillsDir:     agentBDir,
			BinaryName:    "sh",
			PriorityPaths: []string{agentBDir},
			IsEnabled:     true,
			IsCustom:      false,
		},
	}
	if err := configRepo.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	skillPath := filepath.Join(managedSkillsDir, "demo-skill")
	if err := os.MkdirAll(skillPath, 0o755); err != nil {
		t.Fatalf("MkdirAll(skillPath) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillPath, "SKILL.md"), []byte("# Demo Skill\n\nA demo skill."), 0o644); err != nil {
		t.Fatalf("WriteFile(SKILL.md) error = %v", err)
	}

	if err := os.Symlink(skillPath, filepath.Join(agentADir, "demo-skill")); err != nil {
		t.Fatalf("Symlink(agent A) error = %v", err)
	}
	if err := os.Symlink(skillPath, filepath.Join(agentBDir, "demo-skill")); err != nil {
		t.Fatalf("Symlink(agent B) error = %v", err)
	}

	agentService := NewAgentService(configRepo)
	configService := NewConfigService(configRepo)
	skillService := NewSkillService(repository.NewFileSkillRepository(), agentService, configService)

	detail, err := skillService.GetDetail("demo-skill")
	if err != nil {
		t.Fatalf("GetDetail() error = %v", err)
	}

	sort.Strings(detail.Agents)
	wantAgents := []string{"agent-a", "agent-b"}
	if !reflect.DeepEqual(detail.Agents, wantAgents) {
		t.Fatalf("GetDetail() agents = %v, want %v", detail.Agents, wantAgents)
	}
}
