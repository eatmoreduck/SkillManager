# GitHub Actions 使用指南

本文档介绍项目的 CI/CD 配置和使用方法。

## 目录

- [概述](#概述)
- [Workflow 文件](#workflow-文件)
- [触发条件](#触发条件)
- [CI Workflow](#ci-workflow)
- [Release Workflow](#release-workflow)
- [完整发布流程](#完整发布流程)
- [Git Tag 基础](#git-tag-基础)
- [常见问题](#常见问题)

---

## 概述

项目使用 GitHub Actions 实现持续集成和自动发布：

- **CI (持续集成)**：代码提交时自动测试和构建检查
- **Release (发布)**：创建 Release 时自动打包并上传
- **版本固定**：Go 跟随 `go.mod`，Node.js 与 Wails CLI 使用固定稳定版本，避免 `latest` 带来的漂移风险
- **构建可复现**：前端依赖安装统一使用 `npm ci`，CLI 通过 `GITHUB_PATH` 显式暴露，减少 runner 环境差异

---

## Workflow 文件

```
.github/workflows/
├── ci.yml        # 持续集成
└── release.yml   # 自动发布
```

---

## 触发条件

### CI Workflow 触发条件

| 操作 | 是否触发 |
|------|:--------:|
| `git push origin main` | ✅ 触发 |
| `git push origin master` | ✅ 触发 |
| 创建 Pull Request | ✅ 触发 |
| `git push origin v1.0.0` (tag) | ❌ 不触发 |
| 创建 GitHub Release | ❌ 不触发 |

### Release Workflow 触发条件

| 操作 | 是否触发 |
|------|:--------:|
| `git push origin main` | ❌ 不触发 |
| 创建 Pull Request | ❌ 不触发 |
| `git push origin v1.0.0` (tag) | ❌ 不触发 |
| **创建 GitHub Release** | ✅ 触发 |

### 流程图

```
git push origin main
       │
       └──► CI 运行（测试、构建检查）


创建 GitHub Release (v1.0.0)
       │
       └──► Release 运行（打包、上传）
```

---

## CI Workflow

### 触发配置

```yaml
on:
  push:
    branches: [main, master]
    tags-ignore: ['v*']  # 版本 tag 不触发 CI
  pull_request:
    branches: [main, master]
```

### 执行内容

1. **后端测试** (`backend` job)
   - Go 依赖下载
   - 运行 `go test ./...`
   - 构建检查 `go build ./...`

2. **前端构建** (`frontend` job)
   - Node.js 依赖安装
   - TypeScript 类型检查
   - 前端构建

3. **跨平台构建**（仅在 main 分支 push 时）
   - macOS 构建
   - Windows 构建
   - Linux 构建

### 查看状态

- GitHub → Actions → CI
- 也可以在 PR 页面看到 CI 状态

---

## Release Workflow

### 触发配置

```yaml
on:
  release:
    types: [published]
  workflow_dispatch:
    inputs:
      release_tag:
        description: Existing release tag to build and upload assets to
        required: true
        type: string
```

### 执行内容

1. **macOS 构建** - 生成 `.dmg` 文件
2. **Windows 构建** - 生成 NSIS 安装器 `.exe`
3. **Linux 构建** - 生成 `.AppImage`、`.deb`、`.rpm`

所有构建产物会自动上传到 GitHub Release 页面。

如果某个平台资产缺失，也可以手动补传：
- GitHub → Actions → `Release`
- 点击 `Run workflow`
- 输入已有 tag，例如 `v0.0.1`
- workflow 会重新构建并覆盖上传该 release 的资产

### 固定版本

- Go：跟随项目根目录的 `go.mod`
- Node.js：`20.20.1`
- Wails CLI：`v3.0.0-alpha.74`
- 发布架构：macOS 生成 `amd64` + `arm64`，Windows / Linux 当前固定为 `amd64`

### 产物命名规范

- macOS Intel：`SkillManager-<version>-macos-amd64.dmg`
- macOS Apple Silicon：`SkillManager-<version>-macos-arm64.dmg`
- Windows：`SkillManager-<version>-windows-amd64-installer.exe`
- Linux AppImage：`SkillManager-<version>-linux-amd64.AppImage`
- Linux Debian：`skillmanager_<version>_amd64.deb`
- Linux RPM：`skillmanager-<version>-1.x86_64.rpm`

---

## 完整发布流程

### 方式一：命令行 + 网页（推荐新手）

```bash
# 1. 确保代码已提交并推送到 main
git add .
git commit -m "feat: 准备发布 v1.0.0"
git push origin main

# 2. 等待 CI 通过（可选但推荐）

# 3. 创建并推送 tag
git tag v1.0.0
git push origin v1.0.0

# 4. 去 GitHub 网页创建 Release
#    进入仓库 → Releases → Draft a new release
#    选择 tag v1.0.0
#    填写标题和更新日志
#    点击 Publish release

# 5. 等待 Release workflow 完成
#    产物会自动上传到 Release 页面
```

### 补传已有 Release 资产

适用于已经存在 `v0.0.1`，但想补 `macos-arm64` 这类漏掉的资产：

1. GitHub → `Actions` → 选择 `Release`
2. 点击右上角 `Run workflow`
3. `release_tag` 填 `v0.0.1`
4. 等待 workflow 跑完
5. 回到 `Releases` 页面确认新资产已经出现

### 方式二：纯命令行（使用 gh CLI）

```bash
# 确保已安装 gh CLI: brew install gh
gh auth login  # 首次使用需要登录

# 一键发布
git tag v1.0.0 && \
git push origin v1.0.0 && \
gh release create v1.0.0 --title "v1.0.0" --notes "首个正式版本"
```

### 方式三：带更新日志的发布

```bash
# 创建更详细的 Release
gh release create v1.0.0 \
  --title "v1.0.0 - 首个正式版本" \
  --notes-file release-notes.md
```

`release-notes.md` 示例：
```markdown
## 🎉 新功能
- 支持 Claude Code 技能管理
- 支持 Gemini CLI 技能管理

## 🐛 Bug 修复
- 修复配置文件路径问题

## 📦 下载
- macOS Intel: SkillManager-1.0.0-macos-amd64.dmg
- macOS Apple Silicon: SkillManager-1.0.0-macos-arm64.dmg
- Windows: SkillManager-1.0.0-windows-amd64-installer.exe
- Linux: SkillManager-1.0.0-linux-amd64.AppImage / skillmanager_1.0.0_amd64.deb / skillmanager-1.0.0-1.x86_64.rpm
```

---

## Git Tag 基础

### 什么是 Tag？

Tag 是 Git 中用于标记特定提交的"书签"，通常用于标记版本发布。

### 常用命令

```bash
# 查看所有 tag
git tag

# 创建轻量 tag
git tag v1.0.0

# 创建带注释的 tag（推荐）
git tag -a v1.0.0 -m "首个正式版本"

# 推送单个 tag
git push origin v1.0.0

# 推送所有 tag
git push origin --tags

# 删除本地 tag
git tag -d v1.0.0

# 删除远程 tag
git push origin --delete v1.0.0

# 查看 tag 详情
git show v1.0.0

# 基于 tag 创建分支
git checkout -b fix-branch v1.0.0
```

### 版本号规范（Semantic Versioning）

```
v1.2.3
│ │ │
│ │ └── PATCH：Bug 修复，向后兼容
│ └──── MINOR：新功能，向后兼容
└────── MAJOR：重大更新，可能不兼容
```

示例：
- `v0.1.0` → 初始开发版本
- `v1.0.0` → 首个正式版
- `v1.0.1` → 修复 bug
- `v1.1.0` → 新增功能
- `v2.0.0` → 重大更新（可能有破坏性变更）

### 预发布版本

```
v1.0.0-alpha.1   # 内部测试
v1.0.0-beta.1    # 公开测试
v1.0.0-rc.1      # 候选版本
```

---

## GitHub 页面配置

### 必需配置

1. **Workflow 权限**

   Settings → Actions → General → Workflow permissions

   选择：`Read and write permissions`

### 可选配置

2. **Actions 权限**

   Settings → Actions → General

   选择：`Allow all actions and reusable workflows`

3. **macOS 签名（高级）**

   如果要给 macOS 应用签名，需要配置 Secrets：

   | Secret | 说明 |
   |--------|------|
   | `APPLE_CERTIFICATE` | 开发者证书（base64） |
   | `APPLE_CERTIFICATE_PASSWORD` | 证书密码 |
   | `APPLE_ID` | Apple ID |
   | `APPLE_TEAM_ID` | Team ID |
   | `APPLE_PASSWORD` | App-specific password |

---

## 常见问题

### Q: CI 和 Release 会同时运行吗？

A: 不会。CI 在 push 到分支时运行，Release 在创建 GitHub Release 时运行。tag push 不会触发 CI（已配置 `tags-ignore`）。

### Q: 如何跳过 CI？

A: 在 commit message 中添加 `[skip ci]` 或 `[ci skip]`：

```bash
git commit -m "docs: 更新文档 [skip ci]"
```

### Q: Release 失败了怎么办？

A:
1. 去 Actions 页面查看错误日志
2. 修复问题后，删除远程 tag 重新发布：
   ```bash
   git push origin --delete v1.0.0
   git tag -d v1.0.0
   # 修复问题...
   git tag v1.0.0
   git push origin v1.0.0
   gh release create v1.0.0 --title "v1.0.0" --notes "..."
   ```

### Q: 如何查看构建产物？

A:
- CI 产物：Actions → 对应的 workflow run → Artifacts
- Release 产物：Releases → 对应版本 → Assets
  - macOS：`SkillManager-<version>-macos-amd64.dmg`
  - Windows：`SkillManager-<version>-windows-amd64-installer.exe`
  - Linux：`SkillManager-<version>-linux-amd64.AppImage`、`skillmanager_<version>_amd64.deb`、`skillmanager-<version>-1.x86_64.rpm`

### Q: 如何回滚到之前的版本？

A:
```bash
# 查看所有 tag
git tag

# 切换到特定版本
git checkout v1.0.0
```

---

## 参考链接

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)
- [gh CLI 手册](https://cli.github.com/manual/)


## build 发版build 失败后的操作

# 1. 删除远程 tag
git push origin --delete v0.0.1
# 2. 删除本地 tag
git tag -d v0.0.1
# 3. 重新 tag
git tag v0.0.1
# 4. 重新推送 tag
git push origin v0.0.1
# 5. 重新创建 Release（如果之前创建了）
gh release create v0.0.1 --title "v0.0.1" --notes "初始版本"

  ---
如果还没创建 Release，只需要：

# 删除远程 tag
git push origin --delete v0.0.1

# 重新推送
git push origin v0.0.1

# 创建 Release
gh release create v0.0.1 --title "v0.0.1" --notes "初始版本"

  ---
如果 Release 已经存在但失败了：

# 删除 Release 和 tag
gh release delete v0.0.1 --yes
git push origin --delete v0.0.1

# 重新来过
git push origin v0.0.1
gh release create v0.0.1 --title "v0.0.1" --notes "初始版本"
