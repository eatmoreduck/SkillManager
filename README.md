<div align="center">

<img src="frontend/src/assets/logo.png" alt="SkillManager Logo" width="128" height="128">

# SkillManager

**Cross-platform AI Agent Skills Manager**

Manage skills for Claude Code, Gemini CLI, Codex, Cursor, Windsurf and more from one place.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Windows%20%7C%20Linux-lightgrey)](https://github.com/eatmoreduck/skillmanager)

**[English](#overview) | [简体中文](./README.zh-CN.md)**

</div>

---

## Overview

SkillManager is a desktop application that helps you manage AI coding assistant skills across multiple agents. Instead of manually managing `SKILL.md` files in different directories for different tools, SkillManager provides a unified GUI to:

- 🔍 **Discover** skills from remote registries
- 📦 **Install** skills with one click
- 🔄 **Sync** skills across multiple agents (Claude Code, Gemini CLI, Cursor, etc.)
- 📊 **Monitor** which skills are active for each agent

## Supported Agents

| Agent | Status | Skills Directory |
|-------|--------|------------------|
| Claude Code | ✅ | `~/.claude/skills/` |
| Gemini CLI | ✅ | `~/.gemini/skills/` |
| Codex | ✅ | `~/.codex/skills/` |
| Cursor | ✅ | `~/.cursor/rules/` |
| Windsurf | ✅ | `~/.windsurf/rules/` |

## Installation

### Download

Download the latest release for your platform from [Releases](https://github.com/eatmoreduck/skillmanager/releases).

### Build from Source

```bash
# Clone the repository
git clone https://github.com/eatmoreduck/skillmanager.git
cd skillmanager

# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Build
task build

# Or package as .app/.dmg (macOS)
task package
```

## Usage

### GUI Mode

Run the application and manage your skills visually:

```bash
# macOS
open bin/SkillManager.app

# Or run directly
task run
```

### CLI Mode

For quick diagnostics from the command line:

```bash
# Full diagnostics
go run . doctor

# List agents
go run . agents

# Show skill inventory
go run . skills

# Display config
go run . config
```

## Development

```bash
# Start development mode with hot reload
task dev

# Run tests
go test ./...
cd frontend && npm test

# Build frontend only
cd frontend && npm run build
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend (Vue 3)                      │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│  │ Skills  │ │ Agents  │ │Registry │ │ Config  │       │
│  │  View   │ │  View   │ │  View   │ │  View   │       │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘       │
│       └───────────┼───────────┼───────────┘            │
│                   ▼           ▼                         │
│              ┌─────────────────────┐                    │
│              │   Pinia Stores      │                    │
│              └──────────┬──────────┘                    │
└─────────────────────────┼───────────────────────────────┘
                          │ Wails Bindings
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    Backend (Go)                          │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│  │ Binding │ │ Service │ │  Repo   │ │  Model  │       │
│  │  Layer  │─▶  Layer  │─▶  Layer  │─▶  Layer  │       │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │
└─────────────────────────────────────────────────────────┘
```

## Configuration

Configuration is stored in platform-specific locations:

| Platform | Path |
|----------|------|
| macOS | `~/Library/Application Support/skillmanager/config.yaml` |
| Linux | `~/.config/skillmanager/config.yaml` |
| Windows | `%APPDATA%\skillmanager\config.yaml` |

Override with `SKILLMANAGER_CONFIG` environment variable.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

Made with ❤️ by [eatmoreduck](https://github.com/eatmoreduck)

</div>
