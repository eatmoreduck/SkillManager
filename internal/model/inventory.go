package model

// SkillSource describes where an agent-visible skill comes from.
type SkillSource string

const (
	SkillSourceManaged  SkillSource = "managed"
	SkillSourceExternal SkillSource = "external"
	SkillSourceBroken   SkillSource = "broken"
)

// SkillInventoryItem describes a discovered skill in either the manager store or an agent directory.
type SkillInventoryItem struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Path         string      `json:"path"`
	ResolvedPath string      `json:"resolvedPath"`
	Source       SkillSource `json:"source"`
	IsSymlink    bool        `json:"isSymlink"`
	Agents       []string    `json:"agents,omitempty"`
}

// AgentInventory groups discovered skills for a single agent.
type AgentInventory struct {
	Agent    Agent                `json:"agent"`
	Managed  []SkillInventoryItem `json:"managed"`
	External []SkillInventoryItem `json:"external"`
	Broken   []SkillInventoryItem `json:"broken"`
}

// InventoryReport is the aggregated local diagnostic view used by CLI and future UIs.
type InventoryReport struct {
	Platform         string           `json:"platform"`
	ConfigPath       string           `json:"configPath"`
	ManagerSkillsDir string           `json:"managerSkillsDir"`
	Agents           []Agent          `json:"agents"`
	ManagedSkills    []Skill          `json:"managedSkills"`
	AgentInventories []AgentInventory `json:"agentInventories"`
}
