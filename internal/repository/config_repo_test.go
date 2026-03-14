package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileConfigRepositoryLoadConfigCreatesDefaultFile(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	repo, err := NewFileConfigRepository(configPath)
	if err != nil {
		t.Fatalf("NewFileConfigRepository() error = %v", err)
	}

	cfg, err := repo.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Version != "1.0" {
		t.Fatalf("LoadConfig() version = %q, want %q", cfg.Version, "1.0")
	}

	if len(cfg.Agents) < 2 {
		t.Fatalf("LoadConfig() agents = %d, want at least 2", len(cfg.Agents))
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file was not created: %v", err)
	}
}

func TestFileConfigRepositoryMigratesLegacyConfig(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	legacy := []byte(`
proxy:
  enabled: false
registries:
  - name: "skills.sh"
    url: "https://api.skills.sh"
agents:
  - id: "codex"
    isEnabled: true
`)
	if err := os.WriteFile(configPath, legacy, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	repo, err := NewFileConfigRepository(configPath)
	if err != nil {
		t.Fatalf("NewFileConfigRepository() error = %v", err)
	}

	cfg, err := repo.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Version != "1.0" {
		t.Fatalf("LoadConfig() version = %q, want %q", cfg.Version, "1.0")
	}
	if len(cfg.Registries) == 0 || cfg.Registries[0].URL != "https://skills.sh" {
		t.Fatalf("LoadConfig() registry url migration failed: %+v", cfg.Registries)
	}
	if len(cfg.Registries) == 0 || !cfg.Registries[0].IsDefault {
		t.Fatalf("LoadConfig() registry default migration failed: %+v", cfg.Registries)
	}
	if len(cfg.Agents) != 1 || cfg.Agents[0].BinaryName == "" || cfg.Agents[0].SkillsDir == "" {
		t.Fatalf("LoadConfig() agent migration failed: %+v", cfg.Agents)
	}
	if cfg.Proxy.Type == "" || cfg.Proxy.Host == "" || cfg.Proxy.Port == 0 {
		t.Fatalf("LoadConfig() proxy defaults missing: %+v", cfg.Proxy)
	}
}
