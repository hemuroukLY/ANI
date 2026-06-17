# Phase 3 整域 OpenAPI 规划草案（文档层 · 非冻结）

> **状态**：规划草案，供 Phase 3 评审。  
> **覆盖**：AI 原生（×7 ✅）、知识库智能（×3 ✅）、模型增强（×3 ✅）、安全/合规扩展（×4 ✅）。  
> **禁止**：在评审通过前将下文路径写入模块详文「接口冻结规则」为已冻结事实。

---

## 1. AI 原生域（Services）

**索引**：`ai-native/README.md`  
**TASK**：`TASK-SVC-018`

| 建议资源 | 建议路径前缀 | 优先级 |
|---|---|---|
| AgentSession | `/api/v1/svc/agent/sessions` | P2 | **详化** → `../openapi-drafts/phase3/openapi-phase3-agent-session-draft.md` |
| AgentToolGrant | `/api/v1/svc/agent/tool-permissions` | P2 | **详化** → `../openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md` |
| AgentAuditEvent | `/api/v1/svc/agent/audit-events` | P2 | **详化** → `../openapi-drafts/phase3/openapi-phase3-agent-audit-draft.md` |
| PromptGuardPolicy | `/api/v1/svc/agent/guard-policies` | P2 | **详化** → `../openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md` |
| McpToolListing | `/api/v1/svc/agent/mcp/tools` | P3 | **详化** → `../openapi-drafts/phase3/openapi-phase3-mcp-tool-market-draft.md` |
| AgentWorkflow | `/api/v1/svc/agent/workflows` | P3 | **详化** → `../openapi-drafts/phase3/openapi-phase3-agent-orchestration-draft.md` |
| AiNativeSandboxProfile | `/api/v1/svc/agent/sandbox-profiles` | P2 | **详化** → `../openapi-drafts/phase3/openapi-phase3-ai-native-sandbox-security-draft.md` |

**边界**：与 Core `kind=sandbox` 算力实例分离；详见 `ai-native/ai-native-sandbox-security.md`。

---

## 2. 知识库智能（Services）

**索引**：`knowledge/README.md`  
**TASK**：`TASK-SVC-018` 子项  
**验收**：`tasks/phase3/acceptance/PHASE3-KB-INTELLIGENCE-ACCEPTANCE-RECORD.md`

| 模块 | 建议 capability | 说明 |
|---|---|---|
| `kb-doc-intelligence` | `POST .../documents/{doc_id}/analyze` | **详化** → `../openapi-drafts/phase3/openapi-phase3-kb-doc-intelligence-draft.md` |
| | `GET .../documents/{doc_id}/intelligence` | 分析结果只读 |
| `kb-meeting-intelligence` | `POST .../meetings/ingest` | **详化** → `../openapi-drafts/phase3/openapi-phase3-kb-meeting-intelligence-draft.md` |
| | `GET .../meetings`、`GET .../meetings/{id}` | 纪要/行动项 |
| `kb-video-intelligence` | `POST .../videos/ingest` | **详化** → `../openapi-drafts/phase3/openapi-phase3-kb-video-intelligence-draft.md` |
| | `GET .../videos`、`GET .../videos/{id}` | 转写/章节 |

均返回 `202 + AsyncTask`（analyze/ingest）；`task_type` 建议 `kb.document.analyze` / `kb.meeting.ingest` / `kb.video.ingest`；复用 Core `GET /tasks/{task_id}`。

---

## 3. 模型增强（Services + Metering）

**索引**：`inference/README.md`  
**TASK**：`TASK-SVC-018` 子项  
**验收**：`tasks/phase3/acceptance/PHASE3-MODEL-ENHANCEMENT-ACCEPTANCE-RECORD.md`

| 模块 | 建议路径 | 说明 |
|---|---|---|
| `model-encryption` | `GET/PUT .../models/{id}/encryption` | **详化** → `../openapi-drafts/phase3/openapi-phase3-model-encryption-draft.md` |
| `model-recommend-config` | `GET .../models/{id}/recommendations` | **详化** → `../openapi-drafts/phase3/openapi-phase3-model-recommend-config-draft.md` |
| `model-usage-stats` | Core `GET /metering/usage?model_id=` + 可选 `GET .../models/{id}/usage-stats` | **详化** → `../openapi-drafts/phase3/openapi-phase3-model-usage-stats-draft.md` |

加密/推荐合入 **services/v1.yaml**；用量 **优先 Core v1.yaml** metering 扩展。

---

## 4. 安全 / 租户扩展

**索引**：`security/README.md`  
**TASK**：Core 子项 → `TASK-CORE-015`；Services 子项 → `TASK-SVC-018`  
**验收**：`tasks/phase3/acceptance/PHASE3-SECURITY-COMPLIANCE-ACCEPTANCE-RECORD.md`

| 模块 | 建议路径 | 归属 | 草案 |
|---|---|---|---|
| `api-key-audit` | `GET /api/v1/auth/api-keys/{id}/audit-events` | Core Auth | **详化** → `../openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md` |
| `netsec-policy` | `GET/POST /api/v1/networks/policies` + CRUD | Core Networks | **详化** → `../openapi-drafts/phase3/openapi-phase3-netsec-policy-draft.md` |
| `compliance` | `GET /api/v1/svc/compliance/summary`、`.../reports` | Services 只读 | **详化** → `../openapi-drafts/phase3/openapi-phase3-compliance-draft.md` |
| `billing-export` | `POST /api/v1/svc/billing/exports` | Services | **详化** → `../openapi-drafts/phase3/openapi-phase3-billing-export-draft.md` |

---

## 5. 评审顺序建议

1. P0 阻塞：`openapi-p0-yaml-draft.md`（alerts、listTasks、inference events）
2. Phase 2 handler：已声明路径（见 `schema-completion-tracker.md`）
3. AI 原生 ×7 联评（Session → Tool Permission → Audit → Prompt Guard → MCP → Sandbox → Orchestration）
4. ~~知识库智能 ingest（复用 AsyncTask）~~ — **×3 详化草案** ✅ → `knowledge/README.md`
5. ~~模型增强 ×3~~ — **详化草案** ✅ → `inference/README.md`
6. ~~安全合规扩展（§4）~~ — **×4 详化草案** ✅ → `security/README.md`
7. **Phase 3 整域联评**：`tasks/phase3/PHASE3-JOINT-REVIEW-AGENDA.md` → 通过后分批合入 YAML → handler（ANI-main）

---

## 相关文件

- `console-undefined-features-backlog.md` §建议执行顺序 #6
- `tasks/phase3/PHASE3-JOINT-REVIEW-AGENDA.md` — **整域联评议程与签核包**
- `tasks/execution/SERVICES-TEAM-TASKS.md` — TASK-SVC-018
- `tasks/execution/TASK-DEPENDENCY-MAP.md` — Phase 5 批次
