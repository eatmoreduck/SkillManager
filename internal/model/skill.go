package model

import "time"

// Skill represents an installed skill in local storage.
type Skill struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Version     string    `json:"version"`
	Tags        []string  `json:"tags"`
	Agents      []string  `json:"agents"`
	Content     string    `json:"content"`
	LocalPath   string    `json:"localPath"`
	SourceURL   string    `json:"sourceUrl"`
	InstalledAt time.Time `json:"installedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsManaged   bool      `json:"isManaged"`
}
