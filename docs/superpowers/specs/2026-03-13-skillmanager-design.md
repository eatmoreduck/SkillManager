# SkillManager 设计文档

> 跨平台 AI Agent Skills 管理工具

## 1. 项目概述

### 1.1 背景

SkillDeck 是一个优秀的 macOS 原生 Skills 管理工具，但仅支持 macOS。本项目旨在使用 Golang 生态构建一个跨平台（Windows/macOS/Linux）的 Skills 管理工具，让更多开发者受益。

### 1.2 目标用户

个人开发者，需要管理多个 AI Agent（Claude Code、Codex、Gemini CLI 等）的 Skills。

### 1.3 核心价值

- **跨平台支持**：一套工具，三个操作系统
- **统一管理**：一个界面管理所有 Agent 的 Skills
- **便捷安装**：从 Registry 一键安装，自动配置符号链接
- **灵活配置**：支持自定义 Agent 和 Registry

## 2. 技术选型

| 维度 | 选择 | 理由 |
|------|------|------|
| GUI 框架 | Wails v3 | Go 原生，跨平台，性能优秀 |
| 前端框架 | Vue 3 + TypeScript | 渐进式框架，Wails 官方支持 |
| UI 组件库 | Naive UI | Vue 3 原生，中文友好，组件丰富 |
| 数据存储 | 文件系统 | 与 SkillDeck 一致，无额外数据库 |
| 配置格式 | YAML | 人类可读，Go 标准库支持 |
| 后端语言 | Go 1.21+ | 跨平台编译，并发友好 |

## 3. 系统架构

### 3.1 分层架构

```
┌─────────────────────────────────────────┐
│         Frontend (Vue 3 + Naive UI)      │
├─────────────────────────────────────────┤
│            Wails Bindings (Go)           │
├─────────────────────────────────────────┤
│  Service Layer                           │
│  ├── SkillService                        │
│  ├── RegistryService                     │
│  ├── AgentService                        │
│  └── ConfigService                       │
├─────────────────────────────────────────┤
│  Repository Layer                        │
│  ├── SkillRepository (文件系统)          │
│  ├── RegistryRepository (HTTP)           │
│  └── ConfigRepository (YAML)             │
└─────────────────────────────────────────┘
```

### 3.2 数据流

```
用户操作 → Vue 组件 → Pinia Store → Wails Binding → Service → Repository → 文件系统/网络
```

## 4. 项目结构

```
SkillManager/
├── app.go                    # Wails 应用入口
├── main.go                   # 程序入口
├── internal/
│   ├── model/                # 数据模型
│   │   ├── skill.go
│   │   ├── agent.go
│   │   ├── registry.go
│   │   └── config.go
│   ├── repository/           # 数据访问层
│   │   ├── skill_repo.go
│   │   ├── registry_repo.go
│   │   └── config_repo.go
│   ├── service/              # 业务逻辑层
│   │   ├── skill_service.go
│   │   ├── registry_service.go
│   │   ├── agent_service.go
│   │   └── config_service.go
│   └── binding/              # Wails 绑定层
│       ├── skill_binding.go
│       ├── registry_binding.go
│       ├── agent_binding.go
│       └── config_binding.go
├── frontend/                 # Vue 3 前端
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
│   │   │   └── RegistrySelector.vue
│   │   ├── stores/
│   │   │   ├── skillStore.ts
│   │   │   ├── agentStore.ts
│   │   │   └── configStore.ts
│   │   ├── api/
│   │   │   └── wails.ts
│   │   ├── App.vue
│   │   ├── main.ts
│   │   └── router.ts
│   └── wailsjs/              # Wails 自动生成
├── configs/                  # 默认配置
│   └── config.yaml
├── go.mod
├── go.sum
├── wails.json
├── package.json
└── README.md
```

## 5. 数据模型

### 5.1 Skill 模型

```go
// internal/model/skill.go
package model

import "time"

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

### 5.2 Agent 模型

```go
// internal/model/agent.go
package model

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

### 5.3 Registry 模型

```go
// internal/model/registry.go
package model

type Registry struct {
    ID        string `json:"id"`        // 唯一标识
    Name      string `json:"name"`      // 显示名称
    URL       string `json:"url"`       // Registry API 地址
    IsDefault bool   `json:"isDefault"` // 是否默认
}

type RegistrySkill struct {
    ID          string   `json:"id"`          // skill 唯一标识
    Name        string   `json:"name"`        // skill 名称
    Description string   `json:"description"` // 描述
    Author      string   `json:"author"`      // 作者
    Stars       int      `json:"stars"`       // 星标数
    Tags        []string `json:"tags"`        // 标签
    InstallURL  string   `json:"installUrl"`  // 安装 URL（Git）
    Category    string   `json:"category"`    // 分类（all-time, trending, hot）
}
```

### 5.4 Config 模型

```go
// internal/model/config.go
package model

type Config struct {
    Version    string        `yaml:"version" json:"version"`
    Proxy      ProxyConfig   `yaml:"proxy" json:"proxy"`
    Registries []Registry    `yaml:"registries" json:"registries"`
    Agents     []AgentConfig `yaml:"agents" json:"agents"`
}

type ProxyConfig struct {
    Enabled  bool   `yaml:"enabled" json:"enabled"`     // 默认 false
    Type     string `yaml:"type" json:"type"`           // http, https, socks5
    Host     string `yaml:"host" json:"host"`           // 127.0.0.1
    Port     int    `yaml:"port" json:"port"`           // 7890
    Username string `yaml:"username" json:"username"`   // 可选
    Password string `yaml:"password" json:"password"`   // 可选
}

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

## 6. 核心功能

### 6.1 Skill 管理

| 功能 | 描述 | API |
|------|------|-----|
| 列表查看 | 显示所有已安装 skills | `SkillService.ListInstalled()` |
| 详情查看 | 查看 skill 详情和内容 | `SkillService.GetDetail(id)` |
| 安装 Skill | 从 URL 克隆并分配 agent | `SkillService.Install(url, agents)` |
| 卸载 Skill | 删除目录和符号链接 | `SkillService.Uninstall(id)` |
| 更新 Skill | 拉取远程更新 | `SkillService.Update(id)` |
| 编辑 Skill | 修改 SKILL.md | `SkillService.UpdateContent(id, content)` |
| 分配 Agent | 管理 skill 适用的 agent | `SkillService.AssignAgents(id, agents)` |

### 6.2 Registry 浏览

| 功能 | 描述 | API |
|------|------|-----|
| 浏览排行榜 | All Time / Trending / Hot | `RegistryService.Browse(registryID, category)` |
| 搜索 Skill | 按关键词搜索 | `RegistryService.Search(query)` |
| 管理 Registry | 添加/删除/切换 Registry | `RegistryService.AddRegistry()` / `RemoveRegistry()` |

### 6.3 Agent 配置

| 功能 | 描述 | API |
|------|------|-----|
| 自动检测 | 扫描系统已安装 agent | `AgentService.DetectInstalled()` |
| 列出 Agent | 显示所有配置的 agent | `AgentService.ListAgents()` |
| 添加自定义 Agent | 用户自定义 agent | `AgentService.AddCustomAgent(agent)` |
| 启用/禁用 | 控制 agent 管理状态 | `AgentService.ToggleAgent(id, enabled)` |

### 6.4 代理配置

| 功能 | 描述 | API |
|------|------|-----|
| 获取配置 | 读取当前代理设置 | `ConfigService.GetConfig()` |
| 更新代理 | 设置代理参数 | `ConfigService.UpdateProxy(proxy)` |
| 获取 HTTP 客户端 | 返回配置好代理的客户端 | `ConfigService.GetHTTPClient()` |

## 7. 服务层接口

```go
// internal/service/skill_service.go
type SkillService struct {
    skillRepo    repository.SkillRepository
    agentService *AgentService
    configService *ConfigService
}

func (s *SkillService) ListInstalled() ([]model.Skill, error)
func (s *SkillService) GetDetail(id string) (*model.Skill, error)
func (s *SkillService) Install(sourceURL string, agents []string) error
func (s *SkillService) Uninstall(id string) error
func (s *SkillService) Update(id string) (*model.Skill, error)
func (s *SkillService) UpdateContent(id string, content string) error
func (s *SkillService) AssignAgents(id string, agents []string) error

// internal/service/registry_service.go
type RegistryService struct {
    registryRepo repository.RegistryRepository
    configService *ConfigService
}

func (s *RegistryService) ListRegistries() ([]model.Registry, error)
func (s *RegistryService) AddRegistry(registry model.Registry) error
func (s *RegistryService) RemoveRegistry(id string) error
func (s *RegistryService) Browse(registryID string, category string) ([]model.RegistrySkill, error)
func (s *RegistryService) Search(query string) ([]model.RegistrySkill, error)

// internal/service/agent_service.go
type AgentService struct {
    configRepo repository.ConfigRepository
}

func (s *AgentService) ListAgents() ([]model.Agent, error)
func (s *AgentService) DetectInstalled() ([]model.Agent, error)
func (s *AgentService) AddCustomAgent(agent model.Agent) error
func (s *AgentService) RemoveAgent(id string) error
func (s *AgentService) ToggleAgent(id string, enabled bool) error

// internal/service/config_service.go
type ConfigService struct {
    configRepo repository.ConfigRepository
}

func (s *ConfigService) GetConfig() (*model.Config, error)
func (s *ConfigService) UpdateProxy(proxy model.ProxyConfig) error
func (s *ConfigService) GetHTTPClient() *http.Client
```

## 8. Repository 层接口

```go
// internal/repository/skill_repo.go
type SkillRepository interface {
    ScanSkills(agents []model.Agent) ([]model.Skill, error)
    ReadSkill(path string) (*model.Skill, error)
    WriteSkill(path string, skill *model.Skill) error
    DeleteSkill(path string) error
    CloneSkill(sourceURL, targetPath string, proxy *model.ProxyConfig) error
    PullSkill(path string, proxy *model.ProxyConfig) error
    CreateSymlink(skillPath, agentDir string) error
    RemoveSymlink(agentDir, skillName string) error
}

// internal/repository/registry_repo.go
type RegistryRepository interface {
    FetchRegistry(url string, client *http.Client) ([]model.RegistrySkill, error)
    FetchSkillMeta(sourceURL string, client *http.Client) (*model.RegistrySkill, error)
}

// internal/repository/config_repo.go
type ConfigRepository interface {
    LoadConfig() (*model.Config, error)
    SaveConfig(config *model.Config) error
    GetConfigPath() string
}
```

## 9. 文件存储结构

### 9.1 应用数据目录（跨平台）

使用平台特定的标准目录：

| 平台 | 配置目录 | 缓存/日志目录 |
|------|----------|---------------|
| Windows | `%APPDATA%/skillmanager/` | `%LOCALAPPDATA%/skillmanager/` |
| macOS | `~/Library/Application Support/skillmanager/` | `~/Library/Caches/skillmanager/` |
| Linux | `~/.config/skillmanager/` | `~/.local/share/skillmanager/` |

```
${CONFIG_DIR}/                    # 平台特定配置目录
├── config.yaml                   # 主配置文件
└── skills/                       # 本地 skill 存储
    ├── claude-code-tdd/
    │   ├── SKILL.md
    │   └── .git/
    ├── python-patterns/
    │   ├── SKILL.md
    │   └── .git/
    └── ...

${CACHE_DIR}/                     # 平台特定缓存目录
└── logs/
    └── skillmanager.log
```

**Go 实现参考**：

```go
// internal/pkg/paths/paths.go
package paths

import (
    "os"
    "path/filepath"
)

// GetConfigDir 返回平台特定的配置目录
func GetConfigDir() (string, error) {
    dir, err := os.UserConfigDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(dir, "skillmanager"), nil
}

// GetCacheDir 返回平台特定的缓存目录
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
```

### 9.2 Agent Skills 目录（符号链接）

```
~/.claude/skills/
├── claude-code-tdd -> ~/.skillmanager/skills/claude-code-tdd
└── python-patterns -> ~/.skillmanager/skills/python-patterns

~/.codex/skills/
├── claude-code-tdd -> ~/.skillmanager/skills/claude-code-tdd
└── golang-patterns -> ~/.skillmanager/skills/golang-patterns
```

### 9.3 配置文件示例

```yaml
# ${CONFIG_DIR}/config.yaml (平台特定)
version: "1.0"

proxy:
  enabled: false
  type: "http"
  host: "127.0.0.1"
  port: 7890
  username: ""
  password: ""

registries:
  - id: "skills-sh"
    name: "skills.sh"
    url: "https://api.skills.sh"
    isDefault: true

agents:
  - id: "claude"
    name: "Claude Code"
    skillsDir: "~/.claude/skills"
    binaryName: "claude"
    priorityPaths:
      - "~/.claude/skills"
    isEnabled: true
    isCustom: false

  - id: "codex"
    name: "Codex"
    skillsDir: "~/.codex/skills"
    binaryName: "codex"
    priorityPaths:
      - "~/.codex/skills"
      - "~/.agents/skills"
    isEnabled: true
    isCustom: false
```

## 10. 前端设计

### 10.1 页面结构

```
┌─────────────────────────────────────────────────────────────┐
│  SkillManager                              ─ □ ×  [设置]    │
├──────────────┬──────────────────────────────────────────────┤
│              │                                              │
│  📁 My Skills │    [搜索框]  [筛选]  [排序]                  │
│              │    ┌────────────────────────────────────┐    │
│  ☁️ Registry  │    │  Skill 列表 / 卡片视图              │    │
│              │    │                                    │    │
│  ⚙️ Agents    │    │  ┌──────────┐ ┌──────────┐       │    │
│              │    │  │ Skill 1  │ │ Skill 2  │       │    │
│  ─────────── │    │  │ 描述...  │ │ 描述...  │       │    │
│  Agents      │    │  └──────────┘ └──────────┘       │    │
│  ✓ Claude    │    │                                    │    │
│  ✓ Codex     │    │  ┌──────────┐ ┌──────────┐       │    │
│  ○ Gemini    │    │  │ Skill 3  │ │ Skill 4  │       │    │
│  + 自定义    │    │  │ 描述...  │ │ 描述...  │       │    │
│              │    │  └──────────┘ └──────────┘       │    │
│              │    └────────────────────────────────────┘    │
├──────────────┴──────────────────────────────────────────────┤
│  [状态栏: 12 skills | 3 agents | Proxy: Off]                │
└─────────────────────────────────────────────────────────────┘
```

### 10.2 路由设计

| 路由 | 页面 | 功能 |
|------|------|------|
| `/skills` | SkillsView | 本地 skill 列表 |
| `/skills/:id` | SkillDetailView | Skill 详情与编辑 |
| `/registry` | RegistryView | Registry 浏览与搜索 |
| `/agents` | AgentsView | Agent 管理 |
| `/settings` | SettingsView | 代理与 Registry 配置 |

## 11. 技术依赖

### 11.1 Go 后端依赖

```go
require (
    github.com/wailsapp/wails/v3 v3.0.0
    gopkg.in/yaml.v3 v3.0.1
    github.com/gomarkdown/markdown v0.0.0-20231222211730-1d6d59f0b5d2
    github.com/mitchellh/go-homedir v1.1.0
    github.com/sirupsen/logrus v1.9.3
)
```

### 11.2 前端依赖

```json
{
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
    "typescript": "^5.4.0",
    "vite": "^5.2.0",
    "@vitejs/plugin-vue": "^5.0.0",
    "sass": "^1.72.0"
  }
}
```

## 12. 开发计划

### 12.1 MVP 里程碑（4-6 周）

**Phase 1: 基础框架（1 周）**
- 项目初始化（Wails + Vue 3 + Naive UI）
- 基础目录结构搭建
- 配置文件读写
- 代理配置功能
- 基础 UI 框架

**Phase 2: Skill 管理核心（2 周）**
- Agent 自动检测
- Skill 扫描与列表展示
- Skill 详情查看
- Skill 安装（git clone）
- Skill 卸载
- 符号链接管理

**Phase 3: Registry 浏览（1.5 周）**
- Registry API 对接
- 排行榜展示（All Time / Trending / Hot）
- 搜索功能
- 多 Registry 支持
- 一键安装

**Phase 4: 完善与发布（1.5 周）**
- SKILL.md 编辑器
- Skill 更新功能
- 自定义 Agent 添加
- 跨平台测试
- 打包与发布（Windows/macOS/Linux）
- 文档编写

### 12.2 后续版本

**v1.1 - Issue 挖掘**
- GitHub Issue 扫描
- 需求分析与分类
- 痛点聚合展示
- 导出报告

**v1.2 - 增强功能**
- Skill 版本管理
- 批量操作
- 导入/导出配置
- 自动更新检查

## 13. 风险与应对

| 风险 | 影响 | 应对策略 |
|------|------|----------|
| Wails v3 稳定性 | 中 | 使用稳定版本，关注官方更新 |
| 跨平台符号链接 | 高 | Windows 需管理员权限，提供降级方案（复制替代） |
| Registry API 变更 | 中 | 抽象 Repository 层，支持多版本 API |
| Git 命令依赖 | 低 | 检测 git 是否安装，提供友好提示 |
| 大文件/网络超时 | 中 | 添加进度条、超时配置、错误重试 |
| 配置迁移 | 低 | 版本化配置，自动升级逻辑 |

## 14. 发布计划

| 平台 | 输出格式 | 构建命令 |
|------|----------|----------|
| Windows | `.exe` + NSIS 安装包 | `wails build -platform windows/amd64` |
| macOS | `.app` + DMG | `wails build -platform darwin/universal` |
| Linux | `.deb` + `.rpm` + AppImage | `wails build -platform linux/amd64` |

---

**文档版本**: 1.0
**创建日期**: 2026-03-13
**作者**: Claude Code
