package model

// Agent describes a managed AI agent installation.
type Agent struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	SkillsDir     string   `json:"skillsDir"`
	BinaryName    string   `json:"binaryName"`
	PriorityPaths []string `json:"priorityPaths"`
	IsInstalled   bool     `json:"isInstalled"`
	IsEnabled     bool     `json:"isEnabled"`
	IsCustom      bool     `json:"isCustom"`
}

// ToConfig converts a runtime agent into its persisted configuration form.
func (a Agent) ToConfig() AgentConfig {
	return AgentConfig{
		ID:            a.ID,
		Name:          a.Name,
		SkillsDir:     a.SkillsDir,
		BinaryName:    a.BinaryName,
		PriorityPaths: append([]string(nil), a.PriorityPaths...),
		IsEnabled:     a.IsEnabled,
		IsCustom:      a.IsCustom,
	}
}
