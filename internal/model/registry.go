package model

// Registry describes a remote skill source.
type Registry struct {
	ID        string `yaml:"id" json:"id"`
	Name      string `yaml:"name" json:"name"`
	URL       string `yaml:"url" json:"url"`
	IsDefault bool   `yaml:"isDefault" json:"isDefault"`
}

// RegistrySkill describes a skill returned by a registry.
type RegistrySkill struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Stars       int      `json:"stars"`
	Tags        []string `json:"tags"`
	InstallURL  string   `json:"installUrl"`
	Category    string   `json:"category"`
}
