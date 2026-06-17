# Phase 3 AI 原生域 — 文档验收记录

> **生成日期**：2026-06-17  
> **用途**：AI 原生 ×7 子模块 **文档层验收**（详文 + PRD/SPEC + OpenAPI 草案对齐）。  
> **约束**：本阶段**不修改** `ANI-main/**`；YAML 合入与 handler 验通在 ANI-main PR 后回填。  
> **索引**：`docs/console-modules/ai-native/README.md` · **TASK**：`TASK-SVC-018`

---

## 总表

| 模块 | 详文 | PRD/SPEC | OpenAPI 草案 | 文档验收 | YAML 合入 | Handler |
|---|---|---|---|---|---|---|
| Agent Session | `agent-session.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-agent-session-draft.md` | ✅ | ☐ | ☐ |
| Tool Permission | `agent-tool-permission.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md` | ✅ | ☐ | ☐ |
| Agent Audit | `agent-audit.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-agent-audit-draft.md` | ✅ | ☐ | ☐ |
| Prompt Guard | `prompt-injection-guard.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md` | ✅ | ☐ | ☐ |
| MCP Tool Market | `mcp-tool-market.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-mcp-tool-market-draft.md` | ✅ | ☐ | ☐ |
| Sandbox Security | `ai-native-sandbox-security.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-ai-native-sandbox-security-draft.md` | ✅ | ☐ | ☐ |
| Agent Orchestration | `agent-orchestration.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-agent-orchestration-draft.md` | ✅ | ☐ | ☐ |

**文档验收 ✅ 定义**

- 详文含：页面定位、Services/Core 分层、创建前置条件、操作矩阵、接口冻结规则（标注**规划**）、待补边界、验收标准
- PRD/SPEC 链到草案路径；SPEC §2 为 Planned Facts，非 Frozen Facts
- 路径前缀统一 `/api/v1/svc/agent/*`（MCP 为 `/agent/mcp/tools*`）
- Sandbox Security 与 Core `kind=sandbox` 分域；禁止 Services `/sandboxes*`

---

> **整域联评**：`tasks/phase3/PHASE3-JOINT-REVIEW-AGENDA.md`（场次 A：AI 原生）

---

## 联评议程（域内摘要 · 完整版见整域议程 §4.1）

1. **主流程**：Session → Tool Permission → Audit → Prompt Guard  
   - 确认 RBAC scope 命名、422 口径、`idempotency_key` 必填项
2. **安装与隔离**：MCP Market → Sandbox Security  
   - tool id `mcp:{package}/{name}` 与 permission 矩阵；egress 与 audit 事件扩展
3. **编排**：Orchestration  
   - run `202 + task_id` 与 Core AsyncTask；workflow `active` 前置条件
4. **合入批次**：建议 3a（Session/Permission/Audit/Guard/MCP/Sandbox profiles）→ 3b（messages/uninstall/summary/workflow-run 详情）

---

## 文档层 curl 预览（合入 YAML 后执行）

> 以下路径来自**草案**，合入前勿对生产 Gateway 断言。

```bash
# Session
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/agent/sessions"

# Tool Permission
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/agent/tool-permissions"

# MCP Market
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/agent/mcp/tools"

# Sandbox Profiles
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/agent/sandbox-profiles"

# Workflows
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/agent/workflows"
```

---

## 签核

- [x] ×7 详文达 Core 文档标准（2026-06-17）
- [x] ×7 PRD/SPEC 同步草案
- [x] `ai-native/README.md`、`../openapi-drafts/phase3/openapi-phase3-domain-draft.md` §1 更新
- [ ] Services 架构联评通过 — **整域**：`PHASE3-JOINT-REVIEW-AGENDA.md` §7
- [ ] 分批合入 `services/v1.yaml`
- [ ] Handler 运行时验通（ANI-main）

---

## 相关文件

- `docs/console-modules/openapi-drafts/phase3/openapi-phase3-domain-draft.md`
- `tasks/execution/SERVICES-TEAM-TASKS.md` — TASK-SVC-018
- `tasks/execution/TASK-DEPENDENCY-MAP.md` — Phase 5 批次 8
