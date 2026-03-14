package model

import "fmt"

// Config stores persisted application settings.
type Config struct {
	Version    string        `yaml:"version" json:"version"`
	Proxy      ProxyConfig   `yaml:"proxy" json:"proxy"`
	Registries []Registry    `yaml:"registries" json:"registries"`
	Agents     []AgentConfig `yaml:"agents" json:"agents"`
}

// ProxyConfig stores outbound proxy settings.
type ProxyConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Type     string `yaml:"type" json:"type"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

// URL returns the proxy URL when the proxy is enabled.
func (p ProxyConfig) URL() string {
	if !p.Enabled {
		return ""
	}

	auth := ""
	if p.Username != "" {
		auth = p.Username + ":" + p.Password + "@"
	}

	return fmt.Sprintf("%s://%s%s:%d", p.Type, auth, p.Host, p.Port)
}

// AgentConfig stores persisted agent settings.
type AgentConfig struct {
	ID            string   `yaml:"id" json:"id"`
	Name          string   `yaml:"name" json:"name"`
	SkillsDir     string   `yaml:"skillsDir" json:"skillsDir"`
	BinaryName    string   `yaml:"binaryName" json:"binaryName"`
	PriorityPaths []string `yaml:"priorityPaths" json:"priorityPaths"`
	IsEnabled     bool     `yaml:"isEnabled" json:"isEnabled"`
	IsCustom      bool     `yaml:"isCustom" json:"isCustom"`
}

// ToAgent converts a persisted agent config into runtime shape.
func (a AgentConfig) ToAgent() Agent {
	return Agent{
		ID:            a.ID,
		Name:          a.Name,
		SkillsDir:     a.SkillsDir,
		BinaryName:    a.BinaryName,
		PriorityPaths: append([]string(nil), a.PriorityPaths...),
		IsEnabled:     a.IsEnabled,
		IsCustom:      a.IsCustom,
	}
}
