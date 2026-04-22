# Tasks: fix-skill-list-blank-display

## 实现任务

### Task 1: 添加 `normalizeSkillArray` 数据归一化函数

**状态**: ✅ 已完成

**文件**: `frontend/src/stores/skillStore.ts`

**描述**: 新增 `normalizeSkillArray(data: unknown): Skill[]` 工具函数，处理 Wails v3 `Call.ByName` 可能返回的各种格式：
- `Array.isArray(data)` → 直接返回
- 对象包装（含 `data`/`skills`/`result`/`items` 字段）→ 解包返回
- 单个 skill 对象（含 `id`+`name`）→ 包装为数组
- `null`/`undefined`/其他 → 返回空数组

**验证**: TypeScript 编译通过 (`vue-tsc --noEmit`)

---

### Task 2: 重构 `loadSkills()` 使用归一化函数

**状态**: ✅ 已完成

**文件**: `frontend/src/stores/skillStore.ts`

**描述**:
- 替换 `skills.value = result || []` 为 `skills.value = normalizeSkillArray(result)`
- 在 try 块中添加 `console.info('[skillStore] ListInstalled raw result:', ...)` 记录原始返回值类型和内容
- 添加 `console.info('[skillStore] normalized skills count:', ...)` 记录归一化后数量
- 在 catch 块中添加 `console.error('[skillStore] ListInstalled failed:', e)`

**验证**: Vite 构建通过 (`vite build` 成功)

---

### Task 3: SkillsView 计算属性添加防御性检查

**状态**: ✅ 已完成

**文件**: `frontend/src/views/SkillsView.vue`

**描述**:
- `filteredSkills` computed：添加 `Array.isArray(skillStore.skills)` 检查，filter 回调中 `(skill.name || '')` 和 `(skill.description || '')` 空值保护
- `assignedAgentCount` computed：添加 `Array.isArray` 防御
- `tagCount` computed：添加 `Array.isArray` 防御

**验证**: TypeScript 编译 + Vite 构建通过

---

### Task 6: 清理调试产物

**状态**: ✅ 已完成

**执行时间**: 2026-04-11

**描述**: 清理排查过程中产生的截图文件和 Playwright 调试日志。

**已清理**:
- 24 个 `.png` 调试截图 → 已删除
- `.playwright-mcp/` 目录（52 个 console 日志文件）→ 已删除

**验收结果**:
- [x] 项目根目录无 `.png` 文件
- [x] `.playwright-mcp/` 目录已删除
- [ ] `.gitignore` 已更新（建议后续添加 `*.png` 和 `.playwright-mcp/` 排除规则）

---

### Task 4: 验证修复效果

**状态**: ⬜ 待执行（Task 6 完成后）

**前置条件**: Task 6（清理）已完成，Task 1-3 代码修改已就位

**验证步骤**:

1. 启动应用
   ```bash
   task dev
   ```
2. 导航到 Skill 页面 (`/skills`)
3. 检查列表是否正常渲染 Skill 卡片（应显示 46+ 个已安装技能）
4. 打开浏览器 DevTools Console，确认日志输出：
   - `[skillStore] ListInstalled raw result:` — 原始返回值类型和内容
   - `[skillStore] normalized skills count:` — 归一化后数量
   - `[wails] resolved ...` — wailsCompat 桥接层日志
5. 测试搜索功能：输入关键词，确认过滤生效
6. 验证统计数字：hero 区域的"已安装技能"/"关联 Agent"/"标签总数"是否 > 0
7. 导航到 Agents、Registry、Settings 页面，确认无回归

**预期结果**:
- 列表显示所有已安装 Skill 卡片，每张卡片包含名称、作者、版本、描述、标签、Agent 图标
- 搜索框输入后实时过滤
- hero 统计数字正确（非全部为 0）
- 其他页面功能不受影响

**验收标准**:
- [ ] Skill 卡片网格正常渲染（≥ 1 张卡片）
- [ ] 控制台无 `[skillStore] ListInstalled failed` 错误
- [ ] 搜索过滤功能正常
- [ ] hero 统计数字非零
- [ ] Agents/Registry/Settings 页面正常

---

### Task 5: Go 后端测试确认

**状态**: ✅ 已完成

**描述**: 确认后端单元测试全部通过，无回归

**结果**: `go test ./internal/...` 全部 PASS

---

## 任务依赖关系

```
Task 1 (normalizeSkillArray) ← Task 2 (loadSkills 重构)
Task 2 ← Task 3 (SkillsView 防御性检查)
Task 3 ← Task 4 (验证修复效果, 在 Task 6 之后)
Task 5 (后端测试) — 独立，已完成
Task 6 (清理产物) — ✅ 已完成
```

## 不包含的任务

- 不涉及后端 API 接口变更
- 不涉及数据库/文件系统变更
- 不涉及路由或组件结构调整
- 不包含 Git 版本控制操作

## 执行进度

| Task | 描述 | 优先级 | 状态 |
|------|------|--------|------|
| Task 1 | `normalizeSkillArray` 归一化函数 | - | ✅ 完成 |
| Task 2 | `loadSkills()` 使用归一化函数 + 日志 | - | ✅ 完成 |
| Task 3 | SkillsView 计算属性防御性检查 | - | ✅ 完成 |
| Task 4 | 手动验证修复效果 | 普通 | ⬜ 待执行 |
| Task 5 | Go 后端测试确认 | - | ✅ 完成 |
| Task 6 | 清理截图和 Playwright 调试产物 | - | ✅ 完成 |

**整体进度**: 5/6 完成，剩余 Task 4（手动验证）待执行
