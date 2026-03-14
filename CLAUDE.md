# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SkillManager is a cross-platform desktop application for managing AI agent skills (Claude Code, Gemini CLI, Codex, Cursor, Windsurf, etc.). Built with Go + Wails v3 backend and Vue 3 + TypeScript + Naive UI frontend.

## Build Commands

```bash
# Development (with hot reload)
task dev

# Build production binary
task build

# Run production build
task run

# Package as macOS .app + DMG
task package
```

## CLI Commands

The project also supports CLI mode (built without wails tag):

```bash
go run . doctor   # Full local diagnostics
go run . agents   # List configured agents
go run . skills   # Show skill inventory
go run . config   # Display config summary
```

## Testing

```bash
# Go tests
go test ./...
go test ./internal/service/... -v  # Run specific package

# Frontend tests
cd frontend && npm test
```

## Architecture

### Backend (Go)

```
internal/
├── binding/     # Wails bindings - exposes services to frontend
├── service/     # Business logic layer
├── repository/  # Data access layer (config, skills, registry)
├── model/       # Domain models (Skill, Agent, Config, Registry)
└── pkg/         # Shared utilities (paths, etc.)
```

**Dependency flow:** `binding -> service -> repository -> model`

Key entry points:
- `main.go` - CLI entry (build tag: `!wails`)
- `main_wails.go` - GUI entry (build tag: `wails`)
- `app.go` - App container, wires up all dependencies

### Frontend (Vue 3 + TypeScript)

```
frontend/src/
├── views/       # Page components (SkillsView, AgentsView, etc.)
├── components/  # Reusable components (SkillCard, etc.)
├── stores/      # Pinia stores (skillStore, agentStore, etc.)
├── types/       # TypeScript interfaces
├── assets/      # Static assets including agent icons
└── router.ts    # Vue Router configuration
```

**Frontend-Backend communication:** Wails bindings are called via `window.go.main.App.<BindingName>.<Method>()`

### Key Concepts

- **Agent**: An AI coding assistant (claude, gemini, codex, cursor, windsurf) with a skills directory
- **Skill**: A SKILL.md file containing prompts/instructions, stored in agent's skills directory
- **Registry**: Remote skill repository for discovering/installing skills
- **Inventory**: Runtime view of all skills visible to each agent

## Configuration

Config file location (platform-specific):
- macOS: `~/Library/Application Support/skillmanager/config.yaml`
- Linux: `~/.config/skillmanager/config.yaml`
- Windows: `%APPDATA%\skillmanager\config.yaml`

Override with `SKILLMANAGER_CONFIG` environment variable.

## Frontend Development

```bash
cd frontend
npm install
npm run dev      # Start Vite dev server
npm run build    # Build for production
```

## Agent Icons

Agent icons are SVG files in `frontend/src/assets/icons/`. They are dynamically loaded via `import.meta.glob`. To add a new agent icon, simply add a `{agent-id}.svg` file to that directory.
