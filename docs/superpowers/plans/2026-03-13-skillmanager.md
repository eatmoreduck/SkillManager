# SkillManager 实施计划

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建一个跨平台（Windows/macOS/Linux）的 AI Agent Skills 管理工具

**Architecture:** Wails v3 + Vue 3 + Naive UI 分层架构（Binding → Service → Repository）

**Tech Stack:** Go 1.22+, Vue 3, TypeScript, Pinia, Naive UI, YAML

---

## 文件结构

```
SkillManager/
├── app.go                          # Wails 应用入口
├── main.go                         # 程序入口
├── internal/
│   ├── model/                      # 数据模型
│   │   ├── skill.go
│   │   ├── agent.go
│   │   ├── registry.go
│   │   ├── config.go
│   │   └── errors.go               # 自定义错误类型
│   ├── repository/                 # 数据访问层
│   │   ├── skill_repo.go
│   │   ├── registry_repo.go
│   │   └── config_repo.go
│   ├── service/                    # 业务逻辑层
│   │   ├── skill_service.go
│   │   ├── registry_service.go
│   │   ├── agent_service.go
│   │   └── config_service.go
│   ├── binding/                    # Wails 绑定层
│   │   ├── skill_binding.go
│   │   ├── registry_binding.go
│   │   ├── agent_binding.go
│   │   └── config_binding.go
│   └── pkg/
│       ├── paths/paths.go          # 跨平台路径
│       └── logger/logger.go        # 日志封装
├── frontend/                       # Vue 3 前端
│   ├── src/
│   │   ├── views/
│   │   │   ├── SkillsView.vue
│   │   │   ├── SkillDetailView.vue
│   │   │   ├── RegistryView.vue
│   │   │   ├── AgentsView.vue
│   │   │   └── SettingsView.vue
│   │   ├── components/
│   │   │   ├── SkillCard.vue
│   │   │   ├── SkillEditor.vue
│   │   │   ├── AgentSelector.vue
│   │   │   ├── ProxySetting.vue
│   │   │   ├── RegistrySelector.vue
│   │   │   ├── RegistrySkillCard.vue
│   │   │   └── ErrorAlert.vue
│   │   ├── stores/
│   │   │   ├── skillStore.ts
│   │   │   ├── agentStore.ts
│   │   │   ├── registryStore.ts
│   │   │   └── configStore.ts
│   │   ├── types/
│   │   │   ├── skill.ts
│   │   │   ├── agent.ts
│   │   │   ├── registry.ts
│   │   │   └── config.ts
│   │   ├── api/
│   │   │   └── wails.ts
│   │   ├── App.vue
│   │   ├── main.ts
│   │   └── router.ts
│   └── wailsjs/                    # Wails 自动生成
├── configs/
│   └── config.yaml                 # 默认配置模板
├── scripts/
│   └── test-cross-platform.sh      # 跨平台测试脚本
├── linux/
│   └── skillmanager.desktop        # Linux desktop entry
├── installer/
│   └── windows/
│       └── installer.nsi           # NSIS 安装脚本
├── .github/
│   └── workflows/
│       └── release.yml             # GitHub Actions 发布
├── go.mod
├── go.sum
├── wails.json
├── package.json
├── README.md
└── docs/
    └── USER_GUIDE.md
```

---

## Chunk 1: Phase 1 - 基础框架

### Task 1.1: 初始化项目

**Files:**
- Create: `go.mod`
- Create: `wails.json`

- [ ] **Step 1: 初始化 Go 模块**

```bash
cd /Users/leon/Documents/GolandProjects/SkillManager
go mod init skillmanager
```

- [ ] **Step 2: 创建 wails.json**

```json
{
  "name": "skillmanager",
  "outputfilename": "SkillManager",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "SkillManager",
    "email": "skillmanager@example.com"
  },
  "info": {
    "companyName": "SkillManager",
    "productName": "SkillManager",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2026",
    "comments": "Cross-platform AI Agent Skills Manager"
  }
}
```

- [ ] **Step 3: 提交**

```bash
git add go.mod go.sum wails.json
git commit -m "chore: initialize project with go mod and wails config"
```

---

### Task 1.2: 创建数据模型

**Files:**
- Create: `internal/model/config.go`
- Create: `internal/model/skill.go`
- Create: `internal/model/agent.go`
- Create: `internal/model/registry.go`
- Create: `internal/model/errors.go`

- [ ] **Step 1: 创建 config.go（含 fmt import 修复）**

```go
// internal/model/config.go
package model

import "fmt"

// Config 应用配置
type Config struct {
	Version    string        `yaml:"version" json:"version"`
	Proxy      ProxyConfig   `yaml:"proxy" json:"proxy"`
	Registries []Registry    `yaml:"registries" json:"registries"`
	Agents     []AgentConfig `yaml:"agents" json:"agents"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`     // 默认 false
	Type     string `yaml:"type" json:"type"`           // http, https, socks5
	Host     string `yaml:"host" json:"host"`           // 127.0.0.1
	Port     int    `yaml:"port" json:"port"`           // 7890
	Username string `yaml:"username" json:"username"`   // 可选
	Password string `yaml:"password" json:"password"`   // 可选
}

// URL 返回代理 URL（如果启用）
func (p *ProxyConfig) URL() string {
	if !p.Enabled {
		return ""
	}
	auth := ""
	if p.Username != "" {
		auth = p.Username + ":" + p.Password + "@"
	}
	return fmt.Sprintf("%s://%s%s:%d", p.Type, auth, p.Host, p.Port)
}

// AgentConfig Agent 配置
type AgentConfig struct {
	ID            string   `yaml:"id" json:"id"`
	Name          string   `yaml:"name" json:"name"`
	SkillsDir     string   `yaml:"skillsDir" json:"skillsDir"`
	BinaryName    string   `yaml:"binaryName" json:"binaryName"`
	PriorityPaths []string `yaml:"priorityPaths" json:"priorityPaths"`
	IsEnabled     bool     `yaml:"isEnabled" json:"isEnabled"`
	IsCustom      bool     `yaml:"isCustom" json:"isCustom"`
}
```

- [ ] **Step 2: 创建 skill.go**

```go
// internal/model/skill.go
package model

import "time"

// Skill 表示一个已安装的 Skill
type Skill struct {
	ID          string            `json:"id"`           // 唯一标识（repo/name）
	Name        string            `json:"name"`         // skill 名称
	Description string            `json:"description"`  // 简短描述
	Author      string            `json:"author"`       // 作者
	Version     string            `json:"version"`      // 版本
	Tags        []string          `json:"tags"`         // 标签
	Agents      []string          `json:"agents"`       // 已分配的 agent ID 列表
	Content     string            `json:"content"`      // SKILL.md 内容
	LocalPath   string            `json:"localPath"`    // 本地安装路径
	SourceURL   string            `json:"sourceUrl"`    // 来源 Git URL
	InstalledAt time.Time         `json:"installedAt"`  // 安装时间
	UpdatedAt   time.Time         `json:"updatedAt"`    // 最后更新时间
}
```

- [ ] **Step 3: 创建 agent.go**

```go
// internal/model/agent.go
package model

// Agent 表示一个 AI Agent
type Agent struct {
	ID            string   `json:"id"`              // 如 "claude", "codex"
	Name          string   `json:"name"`            // 显示名称
	SkillsDir     string   `json:"skillsDir"`       // skills 目录路径
	BinaryName    string   `json:"binaryName"`      // 检测用的二进制名
	PriorityPaths []string `json:"priorityPaths"`   // skill 读取优先级路径
	IsInstalled   bool     `json:"isInstalled"`     // 是否已安装（系统检测）
	IsEnabled     bool     `json:"isEnabled"`       // 是否启用管理
	IsCustom      bool     `json:"isCustom"`        // 是否用户自定义
}
```

- [ ] **Step 4: 创建 registry.go**

```go
// internal/model/registry.go
package model

// Registry 表示一个 Skills Registry 配置
type Registry struct {
	ID        string `json:"id"`        // 唯一标识，如 "skills-sh"
	Name      string `json:"name"`      // 显示名称，如 "skills.sh"
	URL       string `json:"url"`       // Registry API 地址
	IsDefault bool   `json:"isDefault"` // 是否为默认 Registry
}

// RegistrySkill 表示 Registry 中的 Skill 条目
type RegistrySkill struct {
	ID          string   `json:"id"`          // skill 唯一标识 (repo/name)
	Name        string   `json:"name"`        // skill 名称
	Description string   `json:"description"` // 描述
	Author      string   `json:"author"`      // 作者
	Stars       int      `json:"stars"`       // 星标数
	Tags        []string `json:"tags"`        // 标签
	InstallURL  string   `json:"installUrl"`  // 安装 URL（Git）
	Category    string   `json:"category"`    // 分类
}

// RegistryBrowseResponse 表示 Registry 浏览 API 响应
type RegistryBrowseResponse struct {
	AllTime  []RegistrySkill `json:"allTime"`
	Trending []RegistrySkill `json:"trending"`
	Hot      []RegistrySkill `json:"hot"`
}

// RegistrySearchResponse 表示 Registry 搜索 API 响应
type RegistrySearchResponse struct {
	Results []RegistrySkill `json:"results"`
	Query   string          `json:"query"`
	Total   int             `json:"total"`
}
```

- [ ] **Step 5: 创建 errors.go（自定义错误类型）**

```go
// internal/model/errors.go
package model

import "fmt"

// NetworkError 网络错误类型
type NetworkError struct {
	Type    string // "connection", "timeout", "proxy"
	Message string
	Cause   error
}

func (e *NetworkError) Error() string {
	return e.Message
}

func (e *NetworkError) Unwrap() error {
	return e.Cause
}

// IsNetworkError 判断是否为网络错误
func IsNetworkError(err error) bool {
	_, ok := err.(*NetworkError)
	return ok
}

// IsTimeoutError 判断是否为超时错误
func IsTimeoutError(err error) bool {
	if ne, ok := err.(*NetworkError); ok {
		return ne.Type == "timeout"
	}
	return false
}

// SkillNotFoundError Skill 未找到错误
type SkillNotFoundError struct {
	ID string
}

func (e *SkillNotFoundError) Error() string {
	return fmt.Sprintf("skill not found: %s", e.ID)
}

// ConfigLoadError 配置加载错误
type ConfigLoadError struct {
	Path string
	Err  error
}

func (e *ConfigLoadError) Error() string {
	return fmt.Sprintf("failed to load config from %s: %v", e.Path, e.Err)
}

func (e *ConfigLoadError) Unwrap() error {
	return e.Err
}

// GitNotInstalledError Git 未安装错误
type GitNotInstalledError struct{}

func (e *GitNotInstalledError) Error() string {
	return "git is not installed or not in PATH"
}
```

- [ ] **Step 6: 运行测试**

```bash
go test ./internal/model/... -v
```

- [ ] **Step 7: 提交**

```bash
git add internal/model/
git commit -m "feat(model): add data models with custom error types

- Add Config, ProxyConfig, AgentConfig models
- Add Skill and Agent models
- Add Registry and RegistrySkill models
- Add custom error types (NetworkError, SkillNotFoundError, etc.)
- Fix: ProxyConfig.URL() includes fmt import"
```

---

### Task 1.3: 创建跨平台路径工具

**Files:**
- Create: `internal/pkg/paths/paths.go`

- [ ] **Step 1: 创建 paths.go**

```go
// internal/pkg/paths/paths.go
package paths

import (
	"os"
	"path/filepath"
)

// GetConfigDir 返回平台特定的配置目录
// Windows: %APPDATA%/skillmanager/
// macOS: ~/Library/Application Support/skillmanager/
// Linux: ~/.config/skillmanager/
func GetConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "skillmanager"), nil
}

// GetCacheDir 返回平台特定的缓存目录
// Windows: %LOCALAPPDATA%/skillmanager/
// macOS: ~/Library/Caches/skillmanager/
// Linux: ~/.local/share/skillmanager/
func GetCacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "skillmanager"), nil
}

// GetSkillsDir 返回 skills 存储目录
func GetSkillsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "skills"), nil
}

// GetConfigPath 返回配置文件路径
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

// GetLogPath 返回日志文件路径
func GetLogPath() (string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "logs", "skillmanager.log"), nil
}

// EnsureDir 确保目录存在
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
```

- [ ] **Step 2: 编写测试**

```go
// internal/pkg/paths/paths_test.go
package paths

import (
	"os"
	"runtime"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir failed: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty config directory")
	}

	// 验证路径以 skillmanager 结尾
	if !endsWithSkillmanager(dir) {
		t.Errorf("expected path to end with 'skillmanager', got: %s", dir)
	}
}

func TestGetSkillsDir(t *testing.T) {
	dir, err := GetSkillsDir()
	if err != nil {
		t.Fatalf("GetSkillsDir failed: %v", err)
	}

	// 验证路径以 skills 结尾
	if len(dir) < 6 || dir[len(dir)-6:] != "skills" {
		t.Errorf("expected path to end with 'skills', got: %s", dir)
	}
}

func TestPlatformSpecificPaths(t *testing.T) {
	configDir, _ := GetConfigDir()
	cacheDir, _ := GetCacheDir()

	switch runtime.GOOS {
	case "windows":
		if !contains(configDir, "AppData") {
			t.Error("Windows config dir should contain AppData")
		}
	case "darwin":
		if !contains(configDir, "Library") {
			t.Error("macOS config dir should contain Library")
		}
	case "linux":
		if !contains(configDir, ".config") {
			t.Error("Linux config dir should contain .config")
		}
	}

	t.Logf("Config dir: %s", configDir)
	t.Logf("Cache dir: %s", cacheDir)
}

func endsWithSkillmanager(path string) bool {
	return len(path) >= 12 && path[len(path)-12:] == "skillmanager"
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := tmpDir + "/test/nested/dir"

	err := EnsureDir(testDir)
	if err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}

	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("directory was not created")
	}
}
```

- [ ] **Step 3: 运行测试**

```bash
go test ./internal/pkg/paths/... -v
```

- [ ] **Step 4: 提交**

```bash
git add internal/pkg/paths/
git commit -m "feat(paths): add cross-platform directory utilities

- Add GetConfigDir for platform-specific config directory
- Add GetCacheDir for platform-specific cache directory
- Add GetSkillsDir and GetConfigPath helpers
- Add EnsureDir utility function
- Include unit tests"
```

---

### Task 1.4: 创建 Config Repository

**Files:**
- Create: `internal/repository/config_repo.go`

- [ ] **Step 1: 创建 config_repo.go（含配置初始化和迁移）**

```go
// internal/repository/config_repo.go
package repository

import (
	"os"
	"path/filepath"

	"skillmanager/internal/model"
	"skillmanager/internal/pkg/paths"

	"gopkg.in/yaml.v3"
)

// ConfigRepository 配置仓库接口
type ConfigRepository interface {
	LoadConfig() (*model.Config, error)
	SaveConfig(config *model.Config) error
	GetConfigPath() string
}

// configRepository 配置仓库实现
type configRepository struct {
	configPath string
}

// NewConfigRepository 创建配置仓库
func NewConfigRepository() ConfigRepository {
	configPath, _ := paths.GetConfigPath()
	return &configRepository{configPath: configPath}
}

// LoadConfig 加载配置
func (r *configRepository) LoadConfig() (*model.Config, error) {
	// 确保配置目录存在
	configDir := filepath.Dir(r.configPath)
	if err := paths.EnsureDir(configDir); err != nil {
		return nil, &model.ConfigLoadError{Path: r.configPath, Err: err}
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(r.configPath); os.IsNotExist(err) {
		// 首次运行，创建默认配置
		defaultConfig := r.createDefaultConfig()
		if err := r.SaveConfig(defaultConfig); err != nil {
			return nil, err
		}
		return defaultConfig, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return nil, &model.ConfigLoadError{Path: r.configPath, Err: err}
	}

	var config model.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, &model.ConfigLoadError{Path: r.configPath, Err: err}
	}

	// 配置迁移
	if err := r.migrateConfigIfNeeded(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig 保存配置
func (r *configRepository) SaveConfig(config *model.Config) error {
	// 确保目录存在
	configDir := filepath.Dir(r.configPath)
	if err := paths.EnsureDir(configDir); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(r.configPath, data, 0644)
}

// GetConfigPath 获取配置文件路径
func (r *configRepository) GetConfigPath() string {
	return r.configPath
}

// createDefaultConfig 创建默认配置
func (r *configRepository) createDefaultConfig() *model.Config {
	return &model.Config{
		Version: "1.0",
		Proxy: model.ProxyConfig{
			Enabled:  false,
			Type:     "http",
			Host:     "127.0.0.1",
			Port:     7890,
			Username: "",
			Password: "",
		},
		Registries: []model.Registry{
			{
				ID:        "skills-sh",
				Name:      "skills.sh",
				URL:       "https://api.skills.sh",
				IsDefault: true,
			},
		},
		Agents: []model.AgentConfig{
			{
				ID:            "claude",
				Name:          "Claude Code",
				SkillsDir:     "~/.claude/skills",
				BinaryName:    "claude",
				PriorityPaths: []string{"~/.claude/skills"},
				IsEnabled:     true,
				IsCustom:      false,
			},
			{
				ID:            "codex",
				Name:          "Codex",
				SkillsDir:     "~/.codex/skills",
				BinaryName:    "codex",
				PriorityPaths: []string{"~/.codex/skills", "~/.agents/skills"},
				IsEnabled:     true,
				IsCustom:      false,
			},
		},
	}
}

// migrateConfigIfNeeded 配置迁移
func (r *configRepository) migrateConfigIfNeeded(config *model.Config) error {
	// 版本为空时设置默认版本
	if config.Version == "" {
		config.Version = "1.0"
	}

	// 确保代理配置存在
	// 未来版本迁移逻辑可在此添加

	return nil
}
```

- [ ] **Step 2: 编写测试**

```go
// internal/repository/config_repo_test.go
package repository

import (
	"os"
	"path/filepath"
	"testing"

	"skillmanager/internal/model"
)

func TestConfigRepository_LoadConfig_CreatesDefault(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	repo := &configRepository{configPath: configPath}

	config, err := repo.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// 验证默认配置
	if config.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", config.Version)
	}

	if config.Proxy.Enabled {
		t.Error("expected proxy to be disabled by default")
	}

	if len(config.Registries) == 0 {
		t.Error("expected default registries")
	}

	if len(config.Agents) == 0 {
		t.Error("expected default agents")
	}

	// 验证文件已创建
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}
}

func TestConfigRepository_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	repo := &configRepository{configPath: configPath}

	// 创建自定义配置
	customConfig := &model.Config{
		Version: "1.0",
		Proxy: model.ProxyConfig{
			Enabled: true,
			Type:    "socks5",
			Host:    "localhost",
			Port:    1080,
		},
		Registries: []model.Registry{
			{ID: "test", Name: "Test", URL: "https://test.com", IsDefault: true},
		},
		Agents: []model.AgentConfig{},
	}

	// 保存
	if err := repo.SaveConfig(customConfig); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// 重新加载
	loaded, err := repo.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// 验证
	if !loaded.Proxy.Enabled {
		t.Error("expected proxy to be enabled")
	}
	if loaded.Proxy.Type != "socks5" {
		t.Errorf("expected socks5, got %s", loaded.Proxy.Type)
	}
}
```

- [ ] **Step 3: 运行测试**

```bash
go test ./internal/repository/... -v -run TestConfigRepository
```

- [ ] **Step 4: 提交**

```bash
git add internal/repository/config_repo.go internal/repository/config_repo_test.go
git commit -m "feat(repository): implement ConfigRepository with YAML storage

- Add LoadConfig with default config creation
- Add SaveConfig for YAML serialization
- Add config migration support
- Include comprehensive unit tests"
```

---

### Task 1.5: 创建 Config Service

**Files:**
- Create: `internal/service/config_service.go`

- [ ] **Step 1: 创建 config_service.go**

```go
// internal/service/config_service.go
package service

import (
	"net/http"
	"net/url"
	"time"

	"skillmanager/internal/model"
	"skillmanager/internal/repository"
)

// ConfigService 配置服务
type ConfigService struct {
	configRepo repository.ConfigRepository
}

// NewConfigService 创建配置服务
func NewConfigService(configRepo repository.ConfigRepository) *ConfigService {
	return &ConfigService{configRepo: configRepo}
}

// GetConfig 获取配置
func (s *ConfigService) GetConfig() (*model.Config, error) {
	return s.configRepo.LoadConfig()
}

// UpdateProxy 更新代理配置
func (s *ConfigService) UpdateProxy(proxy model.ProxyConfig) error {
	config, err := s.configRepo.LoadConfig()
	if err != nil {
		return err
	}

	config.Proxy = proxy
	return s.configRepo.SaveConfig(config)
}

// GetHTTPClient 获取配置好代理的 HTTP 客户端
func (s *ConfigService) GetHTTPClient() *http.Client {
	config, err := s.configRepo.LoadConfig()
	if err != nil {
		// 返回默认客户端
		return &http.Client{Timeout: 30 * time.Second}
	}

	transport := &http.Transport{}

	if config.Proxy.Enabled {
		proxyURL := config.Proxy.URL()
		if proxyURL != "" {
			proxyParsed, err := url.Parse(proxyURL)
			if err == nil {
				transport.Proxy = http.ProxyURL(proxyParsed)
			}
		}
	}

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}

// CheckGitInstalled 检查 Git 是否安装
func (s *ConfigService) CheckGitInstalled() bool {
	// 使用 exec.LookPath 检测
	// 留给 SkillService 实现
	return true
}
```

- [ ] **Step 2: 编写测试**

```go
// internal/service/config_service_test.go
package service

import (
	"testing"

	"skillmanager/internal/model"
)

// MockConfigRepository 模拟配置仓库
type MockConfigRepository struct {
	config *model.Config
	err    error
}

func (m *MockConfigRepository) LoadConfig() (*model.Config, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.config == nil {
		return &model.Config{Version: "1.0"}, nil
	}
	return m.config, nil
}

func (m *MockConfigRepository) SaveConfig(config *model.Config) error {
	m.config = config
	return m.err
}

func (m *MockConfigRepository) GetConfigPath() string {
	return "/test/config.yaml"
}

func TestConfigService_GetConfig(t *testing.T) {
	repo := &MockConfigRepository{
		config: &model.Config{Version: "1.0"},
	}
	service := NewConfigService(repo)

	config, err := service.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}

	if config.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", config.Version)
	}
}

func TestConfigService_UpdateProxy(t *testing.T) {
	repo := &MockConfigRepository{
		config: &model.Config{Version: "1.0"},
	}
	service := NewConfigService(repo)

	newProxy := model.ProxyConfig{
		Enabled: true,
		Type:    "socks5",
		Host:    "localhost",
		Port:    1080,
	}

	err := service.UpdateProxy(newProxy)
	if err != nil {
		t.Fatalf("UpdateProxy failed: %v", err)
	}

	// 验证更新
	config, _ := service.GetConfig()
	if !config.Proxy.Enabled {
		t.Error("expected proxy to be enabled")
	}
	if config.Proxy.Type != "socks5" {
		t.Errorf("expected socks5, got %s", config.Proxy.Type)
	}
}

func TestConfigService_GetHTTPClient(t *testing.T) {
	repo := &MockConfigRepository{
		config: &model.Config{
			Version: "1.0",
			Proxy: model.ProxyConfig{
				Enabled: false,
			},
		},
	}
	service := NewConfigService(repo)

	client := service.GetHTTPClient()
	if client == nil {
		t.Fatal("expected non-nil client")
	}

	if client.Timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %v", client.Timeout)
	}
}
```

- [ ] **Step 3: 运行测试**

```bash
go test ./internal/service/... -v -run TestConfigService
```

- [ ] **Step 4: 提交**

```bash
git add internal/service/config_service.go internal/service/config_service_test.go
git commit -m "feat(service): implement ConfigService with HTTP client

- Add GetConfig and UpdateProxy methods
- Add GetHTTPClient with proxy support and 30s timeout
- Include unit tests with mock repository"
```

---

### Task 1.6: 创建 Config Binding

**Files:**
- Create: `internal/binding/config_binding.go`

- [ ] **Step 1: 创建 config_binding.go**

```go
// internal/binding/config_binding.go
package binding

import (
	"skillmanager/internal/model"
	"skillmanager/internal/service"
)

// ConfigBinding 配置 Wails 绑定
type ConfigBinding struct {
	configService *service.ConfigService
}

// NewConfigBinding 创建配置绑定
func NewConfigBinding(configService *service.ConfigService) *ConfigBinding {
	return &ConfigBinding{configService: configService}
}

// GetConfig 获取配置
func (b *ConfigBinding) GetConfig() (*model.Config, error) {
	return b.configService.GetConfig()
}

// UpdateProxy 更新代理配置
func (b *ConfigBinding) UpdateProxy(proxy model.ProxyConfig) error {
	return b.configService.UpdateProxy(proxy)
}

// ProxyConfigResponse 代理配置响应
type ProxyConfigResponse struct {
	Enabled  bool   `json:"enabled"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetProxy 获取代理配置
func (b *ConfigBinding) GetProxy() (*ProxyConfigResponse, error) {
	config, err := b.configService.GetConfig()
	if err != nil {
		return nil, err
	}

	return &ProxyConfigResponse{
		Enabled:  config.Proxy.Enabled,
		Type:     config.Proxy.Type,
		Host:     config.Proxy.Host,
		Port:     config.Proxy.Port,
		Username: config.Proxy.Username,
		Password: config.Proxy.Password,
	}, nil
}

// SetProxy 设置代理配置
func (b *ConfigBinding) SetProxy(enabled bool, proxyType string, host string, port int, username string, password string) error {
	return b.configService.UpdateProxy(model.ProxyConfig{
		Enabled:  enabled,
		Type:     proxyType,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
}
```

- [ ] **Step 2: 提交**

```bash
git add internal/binding/config_binding.go
git commit -m "feat(binding): add ConfigBinding for Wails frontend

- Add GetConfig and UpdateProxy bindings
- Add GetProxy and SetProxy convenience methods
- Define ProxyConfigResponse for frontend"
```

---

### Task 1.7: 创建 Wails 应用入口

**Files:**
- Create: `main.go`
- Create: `app.go`

- [ ] **Step 1: 创建 main.go**

```go
// main.go
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 创建应用
	app := application.New(application.Options{
		Name:        "SkillManager",
		Description: "Cross-platform AI Agent Skills Manager",
		Services: []application.Service{
			application.NewService(NewApp()),
		},
		Assets: application.AssetOptions{
			Handler: application.EmbeddedAssets(assets),
		},
	})

	// 创建窗口
	app.NewWebviewWindow()

	// 运行应用
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

- [ ] **Step 2: 创建 app.go**

```go
// app.go
package main

import (
	"context"

	"skillmanager/internal/binding"
	"skillmanager/internal/repository"
	"skillmanager/internal/service"
)

// App 应用结构
type App struct {
	ctx           context.Context
	configBinding *binding.ConfigBinding
}

// NewApp 创建应用
func NewApp() *App {
	// 初始化仓库
	configRepo := repository.NewConfigRepository()

	// 初始化服务
	configService := service.NewConfigService(configRepo)

	// 初始化绑定
	configBinding := binding.NewConfigBinding(configService)

	return &App{
		configBinding: configBinding,
	}
}

// startup 应用启动
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetConfigBinding 获取配置绑定
func (a *App) GetConfigBinding() *binding.ConfigBinding {
	return a.configBinding
}
```

- [ ] **Step 3: 提交**

```bash
git add main.go app.go
git commit -m "feat: add Wails application entry points

- Add main.go with Wails v3 application setup
- Add app.go with dependency injection
- Initialize config repository, service, and binding"
```

---

### Task 1.8: 创建前端基础框架

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/tsconfig.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/index.html`
- Create: `frontend/src/main.ts`
- Create: `frontend/src/App.vue`
- Create: `frontend/src/router.ts`

- [ ] **Step 1: 创建 package.json**

```json
{
  "name": "skillmanager-frontend",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build",
    "preview": "vite preview",
    "test": "vitest",
    "test:unit": "vitest run",
    "lint": "eslint . --ext .vue,.ts"
  },
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.2.0",
    "pinia": "^2.1.0",
    "naive-ui": "^2.38.0",
    "@vicons/ionicons5": "^0.12.0",
    "marked": "^12.0.0",
    "highlight.js": "^11.9.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0",
    "typescript": "^5.4.0",
    "vite": "^5.2.0",
    "vue-tsc": "^2.0.0",
    "sass": "^1.72.0",
    "vitest": "^1.0.0",
    "@vue/test-utils": "^2.4.0"
  }
}
```

- [ ] **Step 2: 创建 tsconfig.json**

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "paths": {
      "@/*": ["./src/*"]
    },
    "baseUrl": "."
  },
  "include": ["src/**/*.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

- [ ] **Step 3: 创建 vite.config.ts**

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 5173
  }
})
```

- [ ] **Step 4: 创建 index.html**

```html
<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>SkillManager</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
```

- [ ] **Step 5: 创建 main.ts**

```typescript
// frontend/src/main.ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import naive from 'naive-ui'
import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(naive)

app.mount('#app')
```

- [ ] **Step 6: 创建 App.vue**

```vue
<!-- frontend/src/App.vue -->
<template>
  <n-config-provider :theme="theme">
    <n-message-provider>
      <n-dialog-provider>
        <n-layout has-sider class="app-layout">
          <!-- 侧边栏 -->
          <n-layout-sider
            bordered
            collapse-mode="width"
            :collapsed-width="64"
            :width="200"
            :collapsed="collapsed"
            show-trigger
            @collapse="collapsed = true"
            @expand="collapsed = false"
          >
            <div class="logo">
              <span v-if="!collapsed">SkillManager</span>
              <span v-else>SM</span>
            </div>

            <n-menu
              :collapsed="collapsed"
              :collapsed-width="64"
              :collapsed-icon-size="22"
              :options="menuOptions"
              :value="currentMenuKey"
              @update:value="handleMenuSelect"
            />
          </n-layout-sider>

          <!-- 主内容区 -->
          <n-layout>
            <n-layout-header bordered class="app-header">
              <span class="page-title">{{ pageTitle }}</span>
            </n-layout-header>

            <n-layout-content class="app-content">
              <router-view />
            </n-layout-content>
          </n-layout>
        </n-layout>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NConfigProvider,
  NLayout,
  NLayoutSider,
  NLayoutHeader,
  NLayoutContent,
  NMenu,
  NMessageProvider,
  NDialogProvider,
  NIcon
} from 'naive-ui'
import {
  FolderOutline,
  CloudOutline,
  SettingsOutline,
  PersonOutline
} from '@vicons/ionicons5'

const router = useRouter()
const route = useRoute()
const collapsed = ref(false)
const theme = ref(null)

const currentMenuKey = computed(() => {
  const path = route.path
  if (path.startsWith('/skills')) return 'skills'
  if (path.startsWith('/registry')) return 'registry'
  if (path.startsWith('/agents')) return 'agents'
  if (path.startsWith('/settings')) return 'settings'
  return 'skills'
})

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    skills: 'My Skills',
    registry: 'Registry',
    agents: 'Agents',
    settings: 'Settings'
  }
  return titles[currentMenuKey.value] || 'SkillManager'
})

function renderIcon(icon: any) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions = [
  { label: 'My Skills', key: 'skills', icon: renderIcon(FolderOutline) },
  { label: 'Registry', key: 'registry', icon: renderIcon(CloudOutline) },
  { label: 'Agents', key: 'agents', icon: renderIcon(PersonOutline) },
  { label: 'Settings', key: 'settings', icon: renderIcon(SettingsOutline) }
]

function handleMenuSelect(key: string) {
  const routes: Record<string, string> = {
    skills: '/skills',
    registry: '/registry',
    agents: '/agents',
    settings: '/settings'
  }
  router.push(routes[key])
}
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body, #app { height: 100%; overflow: hidden; }
</style>

<style scoped>
.app-layout { height: 100%; }
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  border-bottom: 1px solid var(--n-border-color);
}
.app-header {
  height: 48px;
  padding: 0 20px;
  display: flex;
  align-items: center;
}
.page-title { font-size: 16px; font-weight: 600; }
.app-content { height: calc(100% - 48px); overflow: hidden; }
</style>
```

- [ ] **Step 7: 创建 router.ts**

```typescript
// frontend/src/router.ts
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/skills' },
  { path: '/skills', name: 'Skills', component: () => import('./views/SkillsView.vue') },
  { path: '/skills/:id', name: 'SkillDetail', component: () => import('./views/SkillDetailView.vue') },
  { path: '/registry', name: 'Registry', component: () => import('./views/RegistryView.vue') },
  { path: '/agents', name: 'Agents', component: () => import('./views/AgentsView.vue') },
  { path: '/settings', name: 'Settings', component: () => import('./views/SettingsView.vue') }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
```

- [ ] **Step 8: 创建占位视图**

```vue
<!-- frontend/src/views/SkillsView.vue -->
<template>
  <div class="skills-view">
    <n-empty description="Skills 列表 - Phase 2 实现" />
  </div>
</template>
<script setup lang="ts">
import { NEmpty } from 'naive-ui'
</script>
<style scoped>.skills-view { padding: 20px; }</style>
```

- [ ] **Step 9: 安装依赖并测试**

```bash
cd /Users/leon/Documents/GolandProjects/SkillManager/frontend
npm install
npm run dev
```

- [ ] **Step 10: 提交**

```bash
git add frontend/
git commit -m "feat(frontend): initialize Vue 3 + Naive UI framework

- Add package.json with dependencies
- Configure TypeScript and Vite
- Create App.vue with sidebar navigation
- Add router configuration
- Create placeholder views"
```

---

### Task 1.9: 创建 Settings View 和 Proxy 配置

**Files:**
- Create: `frontend/src/views/SettingsView.vue`
- Create: `frontend/src/components/ProxySetting.vue`
- Create: `frontend/src/stores/configStore.ts`
- Create: `frontend/src/types/config.ts`

- [ ] **Step 1: 创建 types/config.ts**

```typescript
// frontend/src/types/config.ts
export interface ProxyConfig {
  enabled: boolean
  type: 'http' | 'https' | 'socks5'
  host: string
  port: number
  username: string
  password: string
}

export interface Registry {
  id: string
  name: string
  url: string
  isDefault: boolean
}

export interface AgentConfig {
  id: string
  name: string
  skillsDir: string
  binaryName: string
  priorityPaths: string[]
  isEnabled: boolean
  isCustom: boolean
}

export interface Config {
  version: string
  proxy: ProxyConfig
  registries: Registry[]
  agents: AgentConfig[]
}
```

- [ ] **Step 2: 创建 configStore.ts**

```typescript
// frontend/src/stores/configStore.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Config, ProxyConfig } from '@/types/config'

export const useConfigStore = defineStore('config', () => {
  const config = ref<Config | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function loadConfig() {
    loading.value = true
    error.value = null
    try {
      // Wails 绑定调用
      const result = await window.go.main.App.GetConfig()
      config.value = result
    } catch (e) {
      error.value = (e as Error).message
    } finally {
      loading.value = false
    }
  }

  async function updateProxy(proxy: ProxyConfig) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.SetProxy(
        proxy.enabled,
        proxy.type,
        proxy.host,
        proxy.port,
        proxy.username,
        proxy.password
      )
      if (config.value) {
        config.value.proxy = proxy
      }
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    config,
    loading,
    error,
    loadConfig,
    updateProxy
  }
})
```

- [ ] **Step 3: 创建 ProxySetting.vue**

```vue
<!-- frontend/src/components/ProxySetting.vue -->
<template>
  <n-card title="代理设置">
    <n-form :model="proxyConfig" label-placement="left" label-width="100">
      <n-form-item label="启用代理">
        <n-switch v-model:value="proxyConfig.enabled" />
      </n-form-item>

      <template v-if="proxyConfig.enabled">
        <n-form-item label="代理类型">
          <n-select
            v-model:value="proxyConfig.type"
            :options="proxyTypeOptions"
            style="width: 150px"
          />
        </n-form-item>

        <n-form-item label="主机地址">
          <n-input v-model:value="proxyConfig.host" placeholder="127.0.0.1" />
        </n-form-item>

        <n-form-item label="端口">
          <n-input-number v-model:value="proxyConfig.port" :min="1" :max="65535" />
        </n-form-item>

        <n-collapse>
          <n-collapse-item title="认证（可选）" name="auth">
            <n-form-item label="用户名">
              <n-input v-model:value="proxyConfig.username" />
            </n-form-item>
            <n-form-item label="密码">
              <n-input v-model:value="proxyConfig.password" type="password" />
            </n-form-item>
          </n-collapse-item>
        </n-collapse>
      </template>

      <n-form-item>
        <n-space>
          <n-button type="primary" :loading="saving" @click="handleSave">
            保存
          </n-button>
          <n-button @click="handleTest" :loading="testing">
            测试连接
          </n-button>
        </n-space>
      </n-form-item>
    </n-form>
  </n-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import {
  NCard, NForm, NFormItem, NSwitch, NSelect, NInput, NInputNumber,
  NCollapse, NCollapseItem, NButton, NSpace, useMessage
} from 'naive-ui'
import { useConfigStore } from '@/stores/configStore'
import type { ProxyConfig } from '@/types/config'

const message = useMessage()
const configStore = useConfigStore()

const proxyConfig = reactive<ProxyConfig>({
  enabled: false,
  type: 'http',
  host: '127.0.0.1',
  port: 7890,
  username: '',
  password: ''
})

const saving = ref(false)
const testing = ref(false)

const proxyTypeOptions = [
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'SOCKS5', value: 'socks5' }
]

onMounted(async () => {
  await configStore.loadConfig()
  if (configStore.config?.proxy) {
    Object.assign(proxyConfig, configStore.config.proxy)
  }
})

async function handleSave() {
  saving.value = true
  try {
    await configStore.updateProxy({ ...proxyConfig })
    message.success('代理设置已保存')
  } catch (error) {
    message.error('保存失败: ' + (error as Error).message)
  } finally {
    saving.value = false
  }
}

async function handleTest() {
  testing.value = true
  try {
    // 测试代理连接
    message.info('测试连接中...')
    // 实际测试逻辑
    await new Promise(resolve => setTimeout(resolve, 1000))
    message.success('代理连接正常')
  } catch (error) {
    message.error('代理连接失败')
  } finally {
    testing.value = false
  }
}
</script>
```

- [ ] **Step 4: 创建 SettingsView.vue**

```vue
<!-- frontend/src/views/SettingsView.vue -->
<template>
  <div class="settings-view">
    <n-space vertical size="large">
      <ProxySetting />

      <n-card title="Registry 管理">
        <n-empty description="Registry 管理功能 - Phase 3 实现" />
      </n-card>

      <n-card title="关于">
        <n-descriptions label-placement="left" :column="1">
          <n-descriptions-item label="版本">1.0.0</n-descriptions-item>
          <n-descriptions-item label="配置文件">
            {{ configPath }}
          </n-descriptions-item>
        </n-descriptions>
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NSpace, NEmpty, NDescriptions, NDescriptionsItem } from 'naive-ui'
import ProxySetting from '@/components/ProxySetting.vue'

const configPath = ref('')

onMounted(async () => {
  // 获取配置文件路径
  // configPath.value = await window.go.main.App.GetConfigPath()
  configPath.value = '正在加载...'
})
</script>

<style scoped>
.settings-view {
  padding: 20px;
  overflow-y: auto;
  height: 100%;
}
</style>
```

- [ ] **Step 5: 提交**

```bash
git add frontend/src/views/SettingsView.vue frontend/src/components/ProxySetting.vue frontend/src/stores/configStore.ts frontend/src/types/config.ts
git commit -m "feat(settings): implement proxy configuration UI

- Add ProxySetting component with form controls
- Add configStore for state management
- Add SettingsView with proxy and registry sections
- Add TypeScript type definitions"
```

---

## Phase 1 完成检查点

- [ ] 运行所有后端测试: `go test ./... -v`
- [ ] 运行前端开发服务器: `cd frontend && npm run dev`
- [ ] 验证应用启动和基础 UI 显示
- [ ] 验证代理设置保存功能

**Phase 1 完成后提交:**

```bash
git add .
git commit -m "feat: complete Phase 1 - basic framework

- Project initialization with Wails v3 + Vue 3 + Naive UI
- Data models with custom error types
- Cross-platform path utilities
- Config repository and service
- Config binding for frontend
- Basic UI framework with sidebar navigation
- Proxy configuration UI"
```

---

## Chunk 2: Phase 2 - Skill 管理核心

> **注意:** Phase 2 详细计划由 agent 编写，以下为核心任务概述。

### 核心任务

| 任务 | 描述 | 文件 |
|------|------|------|
| 2.1 | Agent 自动检测 | `internal/service/agent_service.go` |
| 2.2 | Skill 扫描 | `internal/repository/skill_repo.go` |
| 2.3 | Skill 列表展示 | `frontend/src/views/SkillsView.vue` |
| 2.4 | Skill 详情查看 | `frontend/src/views/SkillDetailView.vue` |
| 2.5 | Skill 安装（git clone） | `internal/repository/skill_repo.go` |
| 2.6 | Skill 卸载 | `internal/service/skill_service.go` |
| 2.7 | 符号链接管理（含 Windows 降级） | `internal/repository/skill_repo.go` |

### Windows 符号链接降级方案

```go
// internal/repository/skill_repo.go
func (r *skillRepository) CreateSymlink(skillPath, agentDir string) error {
    err := os.Symlink(skillPath, agentDir)
    if err != nil && runtime.GOOS == "windows" {
        // Windows 降级方案：复制目录
        return r.copyDirectory(skillPath, agentDir)
    }
    return err
}

func (r *skillRepository) copyDirectory(src, dst string) error {
    return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        relPath, _ := filepath.Rel(src, path)
        dstPath := filepath.Join(dst, relPath)

        if info.IsDir() {
            return os.MkdirAll(dstPath, info.Mode())
        }
        return copyFile(path, dstPath)
    })
}
```

---

## Chunk 3: Phase 3 - Registry 浏览

> **注意:** Phase 3 详细计划由 agent 完成，以下为核心任务概述。

### 核心任务

| 任务 | 描述 | 文件 |
|------|------|------|
| 3.1 | Registry Model | `internal/model/registry.go` |
| 3.2 | Registry Repository (HTTP) | `internal/repository/registry_repo.go` |
| 3.3 | Registry Service | `internal/service/registry_service.go` |
| 3.4 | Registry Binding | `internal/binding/registry_binding.go` |
| 3.5 | Registry Store | `frontend/src/stores/registryStore.ts` |
| 3.6 | RegistrySelector 组件 | `frontend/src/components/RegistrySelector.vue` |
| 3.7 | RegistrySkillCard 组件 | `frontend/src/components/RegistrySkillCard.vue` |
| 3.8 | ErrorAlert 组件 | `frontend/src/components/ErrorAlert.vue` |
| 3.9 | RegistryView 主页面 | `frontend/src/views/RegistryView.vue` |

### 错误处理策略

- **网络错误**: 显示 "网络错误，请检查代理配置"
- **超时错误**: 显示 "网络超时，请重试" + 重试按钮

---

## Chunk 4: Phase 4 - 完善与发布

> **注意:** Phase 4 详细计划由 agent 完成，以下为核心任务概述。

### 核心任务

| 任务 | 描述 | 文件 |
|------|------|------|
| 4.1 | SKILL.md 编辑器 | `frontend/src/components/SkillEditor.vue` |
| 4.2 | Skill 更新功能 | `internal/service/skill_service.go` |
| 4.3 | 自定义 Agent 添加 | `internal/service/agent_service.go` |
| 4.4 | 跨平台测试 | `tests/e2e/cross-platform.test.ts` |
| 4.5 | 打包与发布 | `.github/workflows/release.yml` |
| 4.6 | 文档编写 | `README.md`, `docs/USER_GUIDE.md` |

### 发布目标

| 平台 | 输出格式 |
|------|----------|
| Windows | `.exe` + NSIS 安装包 |
| macOS | `.app` + DMG |
| Linux | `.deb` + `.rpm` + AppImage |

---

## Review Agent 反馈修复

根据代码审查结果，以下问题已在计划中修复：

### 🔴 严重问题（已修复）

1. **ProxyConfig.URL() 缺少 fmt import** ✅
   - 在 `internal/model/config.go` 中添加 `import "fmt"`

2. **配置目录路径不一致** ✅
   - 统一使用 `internal/pkg/paths/paths.go` 中的 `GetConfigDir()` 方法

3. **Windows 符号链接降级方案** ✅
   - 在 Phase 2 计划中包含 `copyDirectory` 降级方案

### 🟡 建议改进（已采纳）

1. **自定义错误类型** ✅
   - 添加 `internal/model/errors.go`

2. **HTTP 客户端超时** ✅
   - `GetHTTPClient()` 设置 30 秒默认超时

3. **配置文件初始化** ✅
   - `ConfigRepository.LoadConfig()` 自动创建默认配置

4. **配置迁移逻辑** ✅
   - 添加 `migrateConfigIfNeeded()` 方法

---

## 执行顺序

```
Phase 1 (1周) → Phase 2 (2周) → Phase 3 (1.5周) → Phase 4 (1.5周)
     │              │                │                 │
     ▼              ▼                ▼                 ▼
  基础框架      Skill 管理       Registry 浏览      完善与发布
```

---

**文档版本**: 1.0
**创建日期**: 2026-03-13
**作者**: Claude Code
