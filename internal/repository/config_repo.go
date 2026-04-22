package repository

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"skillmanager/internal/model"
	"skillmanager/internal/pkg/paths"
)

// ConfigRepository provides access to persisted configuration.
type ConfigRepository interface {
	LoadConfig() (*model.Config, error)
	SaveConfig(config *model.Config) error
	GetConfigPath() string
}

type FileConfigRepository struct {
	configPath string
}

func NewFileConfigRepository(configPath string) (*FileConfigRepository, error) {
	if configPath == "" {
		defaultPath, err := paths.GetConfigPath()
		if err != nil {
			return nil, err
		}
		configPath = defaultPath
	}

	return &FileConfigRepository{configPath: configPath}, nil
}

func (r *FileConfigRepository) LoadConfig() (*model.Config, error) {
	if err := paths.EnsureDir(filepath.Dir(r.configPath)); err != nil {
		return nil, &model.ConfigLoadError{Path: r.configPath, Err: err}
	}

	data, err := os.ReadFile(r.configPath)
	if errors.Is(err, os.ErrNotExist) {
		cfg := defaultConfig()
		if err := r.SaveConfig(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if err != nil {
		return nil, &model.ConfigLoadError{Path: r.configPath, Err: err}
	}

	var cfg model.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, &model.ConfigLoadError{Path: r.configPath, Err: err}
	}

	changed := migrateConfigIfNeeded(&cfg)
	if changed {
		if err := r.SaveConfig(&cfg); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func (r *FileConfigRepository) SaveConfig(config *model.Config) error {
	if config == nil {
		return model.ErrInvalidConfig
	}

	if err := paths.EnsureDir(filepath.Dir(r.configPath)); err != nil {
		return &model.ConfigSaveError{Path: r.configPath, Err: err}
	}

	payload, err := yaml.Marshal(config)
	if err != nil {
		return &model.ConfigSaveError{Path: r.configPath, Err: err}
	}

	if err := os.WriteFile(r.configPath, payload, 0o644); err != nil {
		return &model.ConfigSaveError{Path: r.configPath, Err: err}
	}

	return nil
}

func (r *FileConfigRepository) GetConfigPath() string {
	return r.configPath
}

func defaultConfig() *model.Config {
	return &model.Config{
		Version:  "1.0",
		Language: "zh-CN",
		Proxy: model.ProxyConfig{
			Enabled: false,
			Type:    "http",
			Host:    "127.0.0.1",
			Port:    7890,
		},
		Registries: []model.Registry{
			{
				ID:        "skills-sh",
				Name:      "skills.sh",
				URL:       "https://skills.sh",
				IsDefault: true,
			},
		},
		Agents: defaultAgents(),
	}
}

func defaultAgents() []model.AgentConfig {
	claudeSkills := paths.HomePath(".claude", "skills")
	codexSkills := paths.HomePath(".codex", "skills")
	geminiSkills := paths.HomePath(".gemini", "skills")
	sharedAgentSkills := paths.HomePath(".agents", "skills")

	return []model.AgentConfig{
		{
			ID:            "claude",
			Name:          "Claude Code",
			SkillsDir:     claudeSkills,
			BinaryName:    "claude",
			PriorityPaths: []string{claudeSkills},
			IsEnabled:     true,
			IsCustom:      false,
		},
		{
			ID:            "codex",
			Name:          "Codex",
			SkillsDir:     codexSkills,
			BinaryName:    "codex",
			PriorityPaths: []string{codexSkills, sharedAgentSkills},
			IsEnabled:     true,
			IsCustom:      false,
		},
		{
			ID:            "gemini",
			Name:          "Gemini CLI",
			SkillsDir:     geminiSkills,
			BinaryName:    "gemini",
			PriorityPaths: []string{geminiSkills},
			IsEnabled:     false,
			IsCustom:      false,
		},
	}
}

func migrateConfigIfNeeded(cfg *model.Config) bool {
	if cfg == nil {
		return false
	}

	changed := false
	defaultCfg := defaultConfig()

	if cfg.Version == "" {
		cfg.Version = defaultCfg.Version
		changed = true
	}
	if cfg.Language == "" {
		cfg.Language = defaultCfg.Language
		changed = true
	}

	if cfg.Proxy.Type == "" {
		cfg.Proxy.Type = defaultCfg.Proxy.Type
		changed = true
	}
	if cfg.Proxy.Host == "" {
		cfg.Proxy.Host = defaultCfg.Proxy.Host
		changed = true
	}
	if cfg.Proxy.Port == 0 {
		cfg.Proxy.Port = defaultCfg.Proxy.Port
		changed = true
	}

	if len(cfg.Registries) == 0 {
		cfg.Registries = append([]model.Registry(nil), defaultCfg.Registries...)
		changed = true
	} else {
		hasDefault := false
		for i := range cfg.Registries {
			if cfg.Registries[i].ID == "" {
				cfg.Registries[i].ID = cfg.Registries[i].Name
				changed = true
			}
			if shouldMigrateSkillsSHRegistry(cfg.Registries[i]) && cfg.Registries[i].URL != defaultCfg.Registries[0].URL {
				cfg.Registries[i].URL = defaultCfg.Registries[0].URL
				changed = true
			}
			if cfg.Registries[i].IsDefault {
				hasDefault = true
			}
		}
		if !hasDefault {
			cfg.Registries[0].IsDefault = true
			changed = true
		}
	}

	if len(cfg.Agents) == 0 {
		cfg.Agents = append([]model.AgentConfig(nil), defaultCfg.Agents...)
		return true
	}

	defaultAgentsByID := make(map[string]model.AgentConfig, len(defaultCfg.Agents))
	for _, agent := range defaultCfg.Agents {
		defaultAgentsByID[agent.ID] = agent
	}

	for i := range cfg.Agents {
		if def, ok := defaultAgentsByID[cfg.Agents[i].ID]; ok {
			if cfg.Agents[i].Name == "" {
				cfg.Agents[i].Name = def.Name
				changed = true
			}
			if cfg.Agents[i].SkillsDir == "" {
				cfg.Agents[i].SkillsDir = def.SkillsDir
				changed = true
			}
			if cfg.Agents[i].BinaryName == "" {
				cfg.Agents[i].BinaryName = def.BinaryName
				changed = true
			}
			if len(cfg.Agents[i].PriorityPaths) == 0 {
				cfg.Agents[i].PriorityPaths = append([]string(nil), def.PriorityPaths...)
				changed = true
			}
		}
	}

	return changed
}

func shouldMigrateSkillsSHRegistry(registry model.Registry) bool {
	url := strings.ToLower(strings.TrimSpace(registry.URL))
	id := strings.ToLower(strings.TrimSpace(registry.ID))
	name := strings.ToLower(strings.TrimSpace(registry.Name))

	return strings.Contains(url, "api.skills.sh") || id == "skills-sh" || name == "skills.sh"
}
