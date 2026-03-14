package repository

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"gopkg.in/yaml.v3"

	"skillmanager/internal/model"
	"skillmanager/internal/pkg/paths"
)

// SkillRepository manages installed skills in the filesystem.
type SkillRepository interface {
	ScanSkills(agents []model.Agent) ([]model.Skill, error)
	ScanManagedSkills(baseDir string) ([]model.Skill, error)
	InspectAgentSkills(agent model.Agent, managedRoot string) (model.AgentInventory, error)
	ReadSkill(path string) (*model.Skill, error)
	WriteSkill(path string, skill *model.Skill) error
	DeleteSkill(path string) error
	CloneSkill(sourceURL, targetPath string, proxy *model.ProxyConfig) error
	PullSkill(path string, proxy *model.ProxyConfig) error
	CreateSymlink(skillPath, agentDir string) error
	RemoveSymlink(agentDir, skillName string) error
}

type FileSkillRepository struct{}

type skillFileMetadata struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Version     string   `yaml:"version"`
	Tags        []string `yaml:"tags"`
	SourceURL   string   `yaml:"source_url"`
}

func NewFileSkillRepository() *FileSkillRepository {
	return &FileSkillRepository{}
}

func (r *FileSkillRepository) ScanSkills(agents []model.Agent) ([]model.Skill, error) {
	byRealPath := map[string]*model.Skill{}

	for _, agent := range agents {
		if !agent.IsEnabled {
			continue
		}

		skillsDir, err := paths.Expand(agent.SkillsDir)
		if err != nil {
			return nil, err
		}

		log.Printf("[skill] scanning agent=%s dir=%s", agent.ID, skillsDir)
		err = r.walkDiscoveredSkills(skillsDir, func(entryPath string) error {
			resolvedPath, resolveErr := filepath.EvalSymlinks(entryPath)
			if resolveErr != nil {
				if errors.Is(resolveErr, fs.ErrNotExist) {
					return nil
				}
				resolvedPath = entryPath
			}

			skill, exists := byRealPath[resolvedPath]
			if !exists {
				skill, err = r.ReadSkill(resolvedPath)
				if err != nil {
					log.Printf("[skill] skipping unreadable skill path=%s resolved=%s err=%v", entryPath, resolvedPath, err)
					return nil
				}
				skill.IsManaged = isManagedSkillPath(resolvedPath)
				byRealPath[resolvedPath] = skill
			}
			skill.Agents = appendIfMissing(skill.Agents, agent.ID)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	skills := make([]model.Skill, 0, len(byRealPath))
	for _, skill := range byRealPath {
		sort.Strings(skill.Agents)
		skills = append(skills, *skill)
	}

	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})

	log.Printf("[skill] scan complete total=%d", len(skills))
	return skills, nil
}

func (r *FileSkillRepository) ScanManagedSkills(baseDir string) ([]model.Skill, error) {
	baseDir, err := paths.Expand(baseDir)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(baseDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	skills := make([]model.Skill, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillPath := filepath.Join(baseDir, entry.Name())
		skill, err := r.ReadSkill(skillPath)
		if err != nil {
			continue
		}
		skill.IsManaged = true
		skills = append(skills, *skill)
	}

	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})

	return skills, nil
}

func (r *FileSkillRepository) InspectAgentSkills(agent model.Agent, managedRoot string) (model.AgentInventory, error) {
	skillsDir, err := paths.Expand(agent.SkillsDir)
	if err != nil {
		return model.AgentInventory{}, err
	}

	managedRoot, err = paths.Expand(managedRoot)
	if err != nil {
		return model.AgentInventory{}, err
	}

	result := model.AgentInventory{Agent: agent}
	if _, statErr := os.Stat(skillsDir); errors.Is(statErr, os.ErrNotExist) {
		return result, nil
	} else if statErr != nil {
		return result, statErr
	}

	seen := make(map[string]struct{})
	err = r.walkDiscoveredSkills(skillsDir, func(currentPath string) error {
		item, ok, err := r.inspectAgentSkillEntry(currentPath, managedRoot)
		if err != nil {
			return err
		}
		if ok {
			appendInventoryItem(&result, item, seen)
		}
		return nil
	})
	if err != nil {
		return result, err
	}

	sortInventoryItems(result.Managed)
	sortInventoryItems(result.External)
	sortInventoryItems(result.Broken)

	return result, nil
}

func appendInventoryItem(result *model.AgentInventory, item model.SkillInventoryItem, seen map[string]struct{}) {
	key := item.Path
	if item.ResolvedPath != "" {
		key = item.ResolvedPath + "::" + item.Path
	}
	if _, ok := seen[key]; ok {
		return
	}
	seen[key] = struct{}{}

	switch item.Source {
	case model.SkillSourceManaged:
		result.Managed = append(result.Managed, item)
	case model.SkillSourceExternal:
		result.External = append(result.External, item)
	case model.SkillSourceBroken:
		result.Broken = append(result.Broken, item)
	}
}

func (r *FileSkillRepository) ReadSkill(path string) (*model.Skill, error) {
	skillDir, err := paths.Expand(path)
	if err != nil {
		return nil, err
	}

	contentPath := filepath.Join(skillDir, "SKILL.md")
	content, err := os.ReadFile(contentPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(skillDir)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(skillDir)
	metadata := parseSkillMetadata(string(content))
	description := extractDescription(string(content))
	if metadata.Name != "" {
		name = metadata.Name
	}
	if metadata.Description != "" {
		description = metadata.Description
	}
	now := info.ModTime()
	sourceURL := metadata.SourceURL
	if sourceURL == "" {
		sourceURL = findGitRemoteURL(skillDir)
	}
	tags := normalizeStrings(metadata.Tags)

	return &model.Skill{
		ID:          name,
		Name:        name,
		Description: description,
		Author:      metadata.Author,
		Version:     metadata.Version,
		Tags:        tags,
		Agents:      []string{},
		Content:     string(content),
		LocalPath:   skillDir,
		SourceURL:   sourceURL,
		InstalledAt: now,
		UpdatedAt:   now,
	}, nil
}

func (r *FileSkillRepository) WriteSkill(path string, skill *model.Skill) error {
	if skill == nil {
		return model.ErrInvalidConfig
	}

	skillDir, err := paths.Expand(path)
	if err != nil {
		return err
	}

	if err := paths.EnsureDir(skillDir); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skill.Content), 0o644); err != nil {
		return err
	}

	return nil
}

func (r *FileSkillRepository) DeleteSkill(path string) error {
	resolved, err := paths.Expand(path)
	if err != nil {
		return err
	}
	return os.RemoveAll(resolved)
}

func (r *FileSkillRepository) CloneSkill(sourceURL, targetPath string, proxy *model.ProxyConfig) error {
	targetPath, err := paths.Expand(targetPath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(targetPath); err == nil {
		return model.ErrAlreadyExists
	}

	if err := paths.EnsureDir(filepath.Dir(targetPath)); err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", sourceURL, targetPath)
	cmd.Env = append(os.Environ(), proxyEnv(proxy)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

func (r *FileSkillRepository) PullSkill(path string, proxy *model.ProxyConfig) error {
	path, err := paths.Expand(path)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "-C", path, "pull", "--ff-only")
	cmd.Env = append(os.Environ(), proxyEnv(proxy)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

func (r *FileSkillRepository) CreateSymlink(skillPath, agentDir string) error {
	skillPath, err := paths.Expand(skillPath)
	if err != nil {
		return err
	}

	agentDir, err = paths.Expand(agentDir)
	if err != nil {
		return err
	}

	if err := paths.EnsureDir(agentDir); err != nil {
		return err
	}

	linkPath := filepath.Join(agentDir, filepath.Base(skillPath))
	if _, err := os.Lstat(linkPath); err == nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		return createWindowsDirectoryJunction(linkPath, skillPath)
	}

	if err := os.Symlink(skillPath, linkPath); err != nil {
		if supportsWindowsLinkFallback(err) {
			return createWindowsDirectoryJunction(linkPath, skillPath)
		}
		return err
	}

	return nil
}

func (r *FileSkillRepository) RemoveSymlink(agentDir, skillName string) error {
	agentDir, err := paths.Expand(agentDir)
	if err != nil {
		return err
	}

	linkPath := filepath.Join(agentDir, skillName)
	if err := os.Remove(linkPath); errors.Is(err, os.ErrNotExist) {
		return nil
	} else {
		return err
	}
}

func extractDescription(content string) string {
	inFrontMatter := false
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if trimmed == "---" {
			inFrontMatter = !inFrontMatter
			continue
		}
		if inFrontMatter {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			trimmed = strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
			if trimmed != "" {
				return trimmed
			}
			continue
		}
		return trimmed
	}
	return ""
}

func parseSkillMetadata(content string) skillFileMetadata {
	lines := strings.Split(content, "\n")
	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return skillFileMetadata{}
	}

	var frontMatter []string
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "---" {
			break
		}
		frontMatter = append(frontMatter, line)
	}
	if len(frontMatter) == 0 {
		return skillFileMetadata{}
	}

	var metadata skillFileMetadata
	if err := yaml.Unmarshal([]byte(strings.Join(frontMatter, "\n")), &metadata); err != nil {
		return skillFileMetadata{}
	}
	return metadata
}

func proxyEnv(proxy *model.ProxyConfig) []string {
	if proxy == nil || !proxy.Enabled {
		return nil
	}

	url := proxy.URL()
	return []string{
		"HTTP_PROXY=" + url,
		"HTTPS_PROXY=" + url,
		"ALL_PROXY=" + url,
	}
}

func appendIfMissing(items []string, value string) []string {
	for _, item := range items {
		if item == value {
			return items
		}
	}
	return append(items, value)
}

func normalizeStrings(values []string) []string {
	result := make([]string, 0, len(values))
	result = append(result, values...)
	return result
}

func (r *FileSkillRepository) walkDiscoveredSkills(root string, visit func(entryPath string) error) error {
	root, err := paths.Expand(root)
	if err != nil {
		return err
	}

	if _, statErr := os.Stat(root); errors.Is(statErr, os.ErrNotExist) {
		return nil
	} else if statErr != nil {
		return statErr
	}

	return filepath.WalkDir(root, func(currentPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if currentPath == root {
			return nil
		}

		if d.Type()&os.ModeSymlink != 0 {
			if err := visit(currentPath); err != nil {
				return err
			}
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		if _, err := os.Stat(filepath.Join(currentPath, "SKILL.md")); err == nil {
			if err := visit(currentPath); err != nil {
				return err
			}
			return filepath.SkipDir
		}

		return nil
	})
}

func (r *FileSkillRepository) inspectAgentSkillEntry(entryPath string, managedRoot string) (model.SkillInventoryItem, bool, error) {
	info, err := os.Lstat(entryPath)
	if errors.Is(err, os.ErrNotExist) {
		return model.SkillInventoryItem{}, false, nil
	}
	if err != nil {
		return model.SkillInventoryItem{}, false, err
	}

	item := model.SkillInventoryItem{
		Name:      filepath.Base(entryPath),
		Path:      entryPath,
		IsSymlink: info.Mode()&os.ModeSymlink != 0,
	}

	resolvedPath, err := filepath.EvalSymlinks(entryPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			item.Source = model.SkillSourceBroken
			return item, true, nil
		}
		resolvedPath = entryPath
	}
	item.ResolvedPath = resolvedPath

	skill, readErr := r.ReadSkill(resolvedPath)
	if readErr != nil {
		item.Source = model.SkillSourceBroken
		return item, true, nil
	}

	item.Name = skill.Name
	item.Description = skill.Description
	if pathWithinBase(resolvedPath, managedRoot) {
		item.Source = model.SkillSourceManaged
	} else {
		item.Source = model.SkillSourceExternal
	}

	return item, true, nil
}

func sortInventoryItems(items []model.SkillInventoryItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
}

func pathWithinBase(path string, base string) bool {
	if path == "" || base == "" {
		return false
	}

	cleanPath := canonicalPath(path)
	cleanBase := canonicalPath(base)
	rel, err := filepath.Rel(cleanBase, cleanPath)
	if err != nil {
		return false
	}

	return rel == "." || (!strings.HasPrefix(rel, "..") && rel != "..")
}

func canonicalPath(path string) string {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return filepath.Clean(resolved)
	}
	return filepath.Clean(path)
}

func findGitRemoteURL(skillDir string) string {
	for current := filepath.Clean(skillDir); current != "." && current != string(filepath.Separator); current = filepath.Dir(current) {
		gitPath := filepath.Join(current, ".git")
		info, err := os.Stat(gitPath)
		if err != nil {
			parent := filepath.Dir(current)
			if parent == current {
				break
			}
			continue
		}

		if !info.IsDir() {
			return ""
		}

		if remote := parseGitConfigRemote(filepath.Join(gitPath, "config")); remote != "" {
			return remote
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}

	return ""
}

func parseGitConfigRemote(configPath string) string {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	inOrigin := false
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "[remote ") && strings.Contains(trimmed, "\"origin\""):
			inOrigin = true
		case strings.HasPrefix(trimmed, "["):
			inOrigin = false
		case inOrigin && strings.HasPrefix(trimmed, "url ="):
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "url ="))
		}
	}

	return ""
}

func isManagedSkillPath(path string) bool {
	managedRoot, err := paths.GetSkillsDir()
	if err != nil {
		return false
	}
	managedRoot, err = paths.Expand(managedRoot)
	if err != nil {
		return false
	}
	return pathWithinBase(path, managedRoot)
}

func createWindowsDirectoryJunction(linkPath string, targetPath string) error {
	cmd := exec.Command("cmd", "/c", "mklink", "/J", linkPath, targetPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mklink /J failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func supportsWindowsLinkFallback(err error) bool {
	if runtime.GOOS != "windows" {
		return false
	}

	var linkErr *os.LinkError
	if !errors.As(err, &linkErr) {
		return false
	}

	return errors.Is(linkErr.Err, syscall.EPERM)
}
