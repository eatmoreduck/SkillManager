package repository

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"skillmanager/internal/model"
)

func TestFileSkillRepositoryInspectAgentSkillsClassifiesSources(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("broken symlink setup differs on Windows; inventory logic is covered on Unix-like platforms")
	}

	tempDir := t.TempDir()
	managerRoot := filepath.Join(tempDir, "manager")
	agentRoot := filepath.Join(tempDir, "agent")

	for _, dir := range []string{managerRoot, agentRoot} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", dir, err)
		}
	}

	managedSkillDir := filepath.Join(managerRoot, "managed-skill")
	externalSkillDir := filepath.Join(tempDir, "external-skill")
	brokenLink := filepath.Join(agentRoot, "broken-skill")

	for _, dir := range []string{managedSkillDir, externalSkillDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll(skill dir) error = %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# "+filepath.Base(dir)), 0o644); err != nil {
			t.Fatalf("WriteFile(SKILL.md) error = %v", err)
		}
	}

	if err := os.Symlink(managedSkillDir, filepath.Join(agentRoot, "managed-skill")); err != nil {
		t.Fatalf("Symlink(managed) error = %v", err)
	}
	if err := os.Symlink(externalSkillDir, filepath.Join(agentRoot, "external-skill")); err != nil {
		t.Fatalf("Symlink(external) error = %v", err)
	}
	if err := os.Symlink(filepath.Join(tempDir, "missing-target"), brokenLink); err != nil {
		t.Fatalf("Symlink(broken) error = %v", err)
	}

	repo := NewFileSkillRepository()
	inventory, err := repo.InspectAgentSkills(model.Agent{
		ID:        "codex",
		Name:      "Codex",
		SkillsDir: agentRoot,
		IsEnabled: true,
	}, managerRoot)
	if err != nil {
		t.Fatalf("InspectAgentSkills() error = %v", err)
	}

	if len(inventory.Managed) != 1 {
		t.Fatalf("Managed len = %d, want 1", len(inventory.Managed))
	}
	if len(inventory.External) != 1 {
		t.Fatalf("External len = %d, want 1", len(inventory.External))
	}
	if len(inventory.Broken) != 1 {
		t.Fatalf("Broken len = %d, want 1", len(inventory.Broken))
	}
}

func TestFileSkillRepositoryScanSkillsFindsNestedSkillsAndMetadata(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	agentRoot := filepath.Join(tempDir, "agent")
	if err := os.MkdirAll(filepath.Join(agentRoot, ".system", "openai-docs"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	content := `---
name: "OpenAI Docs"
description: "Official OpenAI documentation helper"
author: "OpenAI"
version: "1.2.3"
tags: ["docs", "openai"]
source_url: "https://github.com/openai/docs-skill.git"
---

# Ignored Heading
`
	if err := os.WriteFile(filepath.Join(agentRoot, ".system", "openai-docs", "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(SKILL.md) error = %v", err)
	}

	repo := NewFileSkillRepository()
	skills, err := repo.ScanSkills([]model.Agent{{
		ID:        "codex",
		Name:      "Codex",
		SkillsDir: agentRoot,
		IsEnabled: true,
	}})
	if err != nil {
		t.Fatalf("ScanSkills() error = %v", err)
	}

	if len(skills) != 1 {
		t.Fatalf("ScanSkills() len = %d, want 1", len(skills))
	}

	skill := skills[0]
	if skill.Name != "OpenAI Docs" {
		t.Fatalf("skill.Name = %q, want %q", skill.Name, "OpenAI Docs")
	}
	if skill.Author != "OpenAI" || skill.Version != "1.2.3" {
		t.Fatalf("skill metadata mismatch: %+v", skill)
	}
	if skill.SourceURL != "https://github.com/openai/docs-skill.git" {
		t.Fatalf("skill.SourceURL = %q", skill.SourceURL)
	}
	if len(skill.Agents) != 1 || skill.Agents[0] != "codex" {
		t.Fatalf("skill.Agents = %v, want [codex]", skill.Agents)
	}
}

func TestFileSkillRepositoryReadSkillNormalizesEmptySlices(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	skillDir := filepath.Join(tempDir, "plain-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	content := `# Plain Skill

No front matter tags here.
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(SKILL.md) error = %v", err)
	}

	repo := NewFileSkillRepository()
	skill, err := repo.ReadSkill(skillDir)
	if err != nil {
		t.Fatalf("ReadSkill() error = %v", err)
	}

	if skill.Tags == nil {
		t.Fatalf("skill.Tags is nil, want empty slice")
	}
	if skill.Agents == nil {
		t.Fatalf("skill.Agents is nil, want empty slice")
	}
}
