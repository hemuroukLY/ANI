# Phase 3 整域联评 — 议程与签核包

> **状态**：联评 **待召开**（文档层 17 套草案 ✅ · YAML/handler ☐）  
> **日期建议**：2026-06-18 起，建议 **2～3 场**（见 §3）  
> **约束**：联评通过前不得将草案 path 写入模块详文「接口冻结规则」为已冻结事实；本阶段 service-ani **不改** `ANI-main/**`。  
> **权威索引**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-domain-draft.md`

---

## 1. 联评目标

| 目标 | 产出 |
|---|---|
| 确认 17 套 OpenAPI 草案可进入 YAML PR | 各域 **通过 / 有条件通过 / 驳回** 记录 |
| 冻结跨域口径 | AsyncTask、审计事件、RBAC、422、BOSS 边界（§5） |
| 确定合入批次 | Core `v1.yaml` vs Services `services/v1.yaml` PR 切分（§6） |
| 指派实现 owner | TASK-SVC-018、TASK-CORE-015 子项 handler 负责人 |

**非目标**：本场不实现 handler；不替代 P0 草案（`openapi-p0-yaml-draft.md`，状态 **待定**）评审。

---

## 2. 参会角色与预读

| 角色 | 建议人选 | 预读（必选） |
|---|---|---|
| Services 架构 | Services TL / OpenAPI owner | `../openapi-drafts/phase3/openapi-phase3-domain-draft.md`、§1～§3 草案 |
| Core 架构 | Core TL | §3 usage-stats、§4 api-key-audit / netsec、AsyncTask |
| Console 产品 | Console PM | 各域 README、操作矩阵 |
| 安全 / 合规 | 安全架构 + BOSS 代表 | §4 compliance / billing-export、审计分工 |
| 知识库 / 推理 | 领域 owner | KB ×3、模型 ×3 草案 |
| AI 原生 | Agent 领域 owner | AI ×7 草案、`PHASE3-AI-NATIVE-ACCEPTANCE-RECORD.md` |

**预读包（按域）**

| 域 | 模块数 | 域索引 | 验收记录 |
|---|---|---|---|
| AI 原生 | ×7 | `ai-native/README.md` | `PHASE3-AI-NATIVE-ACCEPTANCE-RECORD.md` |
| 知识库智能 | ×3 | `knowledge/README.md` | `PHASE3-KB-INTELLIGENCE-ACCEPTANCE-RECORD.md` |
| 模型增强 | ×3 | `inference/README.md` | `PHASE3-MODEL-ENHANCEMENT-ACCEPTANCE-RECORD.md` |
| 安全/合规 | ×4 | `security/README.md` | `PHASE3-SECURITY-COMPLIANCE-ACCEPTANCE-RECORD.md` |

**OpenAPI 草案清单（17）**：见 [附录 A](#附录-a--openapi-草案路径一览)。

---

## 3. 会议安排（建议）

### 场次 A — 跨域基础 + AI 原生（≈90 min）

| 时段 | 议题 | 决策点 |
|---|---|---|
| 0:00–0:20 | **跨域口径**（§5 全表） | AsyncTask enum、idempotency_key、422 风格 |
| 0:20–0:50 | AI 原生主流程 | Session / Tool Permission / Audit / Prompt Guard |
| 0:50–1:10 | AI 原生扩展 | MCP Market、Sandbox Security、Orchestration |
| 1:10–1:30 | 签核 | 场次 A 结论写入 §7 |

### 场次 B — 知识库 + 模型（≈75 min）

| 时段 | 议题 | 决策点 |
|---|---|---|
| 0:00–0:30 | KB 智能 ×3 | upload vs analyze、ingest task_type、3a/3b |
| 0:30–0:55 | 模型增强 ×3 | metering 扩展、recommend→deploy 映射、encryption binding |
| 0:55–1:15 | 与 usage-report 对齐 | `model_id` filter、group_by 枚举 |
| 1:15–1:30 | 签核 | 场次 B 结论写入 §7 |

### 场次 C — 安全/合规 + 合入计划（≈60 min）

| 时段 | 议题 | 决策点 |
|---|---|---|
| 0:00–0:25 | 审计三分工 | api-key-audit / audit-log / agent-audit |
| 0:25–0:40 | netsec vs SG | `/networks/policies` 独立资源确认 |
| 0:40–0:50 | BOSS 边界 | compliance 只读、billing-export 无发票金额 |
| 0:50–1:00 | **YAML 合入批次**（§6） | PR 顺序、3a/3b 首发范围 |
| 1:00–1:10 | 总签核 | §7 整域结论、TASK 指派 |

> 若时间紧：场次 A+B 合并为一场「Services 域」；场次 C 单独「Core + BOSS」。

---

## 4. 分域议题速查

### 4.1 AI 原生 ×7

1. RBAC：`agent-sessions:*`、`agent-mcp:*`、`agent-sandbox:*`、`agent-workflows:*` 是否纳入 scope 注册表  
2. MCP tool id：`mcp:{package_id}/{tool_name}` 与 `listAgentTools` 一致  
3. Sandbox：`sandbox-profiles` vs Core `instances?kind=sandbox`  
4. Orchestration：run `202` + `task_id`；`WORKFLOW_NOT_ACTIVE` 422  
5. **3b Defer**：messages list、MCP uninstall、workflow-run 单查、sandbox summary  

### 4.2 知识库智能 ×3

1. `kb.document.analyze` vs upload `kb.parse`/`kb.index` 去重  
2. `DOCUMENT_NOT_INDEXED` 422 触发条件  
3. meeting vs video 输入源与 doc-intelligence 管线共享边界  
4. **3b Defer**：intelligence history、DELETE meeting、transcript 分片  

### 4.3 模型增强 ×3

1. **Core 优先**：`GET /metering/usage?model_id=` 是否首发；Services `usage-stats` 是否 3b  
2. `ModelDeploymentRecommendation` ↔ `CreateInferenceServiceRequest` 字段映射表（ADR-1 页）  
3. encryption：model vs version 粒度；re-encrypt 是否 202  
4. `ENCRYPTION_KEY_REVOKED`、`NO_SUITABLE_GPU_PROFILE` 422  

### 4.4 安全/合规 ×4

1. 审计：三类事件 schema 可对齐字段（occurred_at、actor、resource、action）  
2. netsec：`policy_type` 枚举；scope=vpc 时 vpc_id 必填 422  
3. compliance：**禁止** Console POST 创建报告  
4. billing-export：`billing.export` task_type；export 不含结算金额  

---

## 5. 跨域决策表（联评必决）

| ID | 主题 | 草案现状 | 建议默认（可被联评推翻） | 决议栏 |
|---|---|---|---|---|
| X-01 | AsyncTask `task_type` 命名 | 点分：`kb.document.analyze`、`billing.export` | 统一 **点分小写**；与现有 `model.import` 对齐 | ☐ |
| X-02 | 写操作幂等 | 各草案 POST/PUT 要求 `idempotency_key` | **全部保留** | ☐ |
| X-03 | 422 错误码 | 各草案建议语义码 | 合入 YAML 时必须在 `PreconditionFailed` **举例**；详文引用 §2.10 | ☐ |
| X-04 | Agent vs KB Session | 分 path | **禁止**混用 `/knowledge-bases/.../sessions` 与 `/agent/sessions` | ☐ |
| X-05 | 审计三分工 | 三份草案 | Key 调用 / 平台操作 / Agent 行为 **三 API 并存** | ☐ |
| X-06 | Metering 扩展 | usage-stats 草案 | **Core filter 首发**；Services summary 可 3b | ☐ |
| X-07 | BOSS 边界 | compliance、billing | Console **只读**摘要 + export 用量；发票/取证 BOSS | ☐ |
| X-08 | 3a / 3b 切分 | 各草案标注 3b op | 联评勾选 **3a 首发 op 清单**（§6.2） | ☐ |
| X-09 | RBAC scope 注册 | 各草案 x-ani-rbac-scope | 合入前更新 **scope 总表**（单页 ADR 或 YAML 注释） | ☐ |
| X-10 | 合入 repo | v1 vs services | 见 §6.1；**禁止**自造第三份 yaml | ☐ |

---

## 6. YAML 合入批次建议

### 6.1 PR 切分（ANI-main）

| PR | 文件 | 范围 | 依赖 |
|---|---|---|---|
| **P3-C1** | `v1.yaml` | metering `model_id` filter（若 X-06 通过） | usage-report 联评 |
| **P3-C2** | `v1.yaml` | auth api-key audit-events；networks policies | CORE-015 |
| **P3-S1** | `services/v1.yaml` | Agent 3a：session、tool-permission、audit、guard | SVC-018 |
| **P3-S2** | `services/v1.yaml` | KB 3a：analyze、meetings/videos ingest+list/get | SVC-018 |
| **P3-S3** | `services/v1.yaml` | Model：encryption、recommendations | SVC-018 |
| **P3-S4** | `services/v1.yaml` | MCP、sandbox-profiles、workflows 3a | SVC-018 |
| **P3-S5** | `services/v1.yaml` | compliance read-only、billing exports | SVC-018 |
| **P3-*b** | 按需 | 各域 3b ops | 对应 3a handler 稳定后 |

### 6.2 建议 3a 首发 operation 清单（联评勾选）

**Services — Agent**

- [ ] `listAgentSessions` / `createAgentSession` / `getAgentSession` / `updateAgentSession`
- [ ] `listAgentTools` / `getAgentToolPermissions` / `updateAgentToolPermissions`
- [ ] `listAgentAuditEvents`
- [ ] `listPromptGuardPolicies` / `updatePromptGuardPolicy`
- [ ] `listMCPMarketToolPackages` / `installMCPMarketToolPackage`
- [ ] `listAgentSandboxProfiles` / `updateAgentSandboxProfile`
- [ ] `listAgentWorkflows` / `createAgentWorkflow` / `startAgentWorkflowRun`

**Services — Knowledge**

- [ ] `analyzeKnowledgeBaseDocument` / `getKnowledgeBaseDocumentIntelligence`
- [ ] `ingestKnowledgeBaseMeeting` / `listKnowledgeBaseMeetings` / `getKnowledgeBaseMeeting`
- [ ] `ingestKnowledgeBaseVideo` / `listKnowledgeBaseVideos` / `getKnowledgeBaseVideo`

**Services — Model / Billing / Compliance**

- [ ] `getModelEncryptionBinding` / `updateModelEncryptionBinding`
- [ ] `getModelDeploymentRecommendations`
- [ ] `getComplianceSummary` / `listComplianceReports`
- [ ] `createBillingExport` / `getBillingExport`

**Core**

- [ ] `listApiKeyAuditEvents`
- [ ] `listNetworkSecurityPolicies` / `createNetworkSecurityPolicy` / `getNetworkSecurityPolicy`
- [ ] `GET /metering/usage` 扩展 query（若 X-06 通过）

---

## 7. 签核表

### 7.1 分域结论

| 域 | 模块数 | 文档层 | 联评结论 | 联评日期 | 备注 |
|---|---|---|---|---|---|
| AI 原生 | 7 | ✅ | ☐ 通过 / ☐ 有条件通过 / ☐ 驳回 | | |
| 知识库智能 | 3 | ✅ | ☐ / ☐ / ☐ | | |
| 模型增强 | 3 | ✅ | ☐ / ☐ / ☐ | | |
| 安全/合规 | 4 | ✅ | ☐ / ☐ / ☐ | | |

**有条件通过**定义：仅 §5 决策表中有明确 ADR/跟进项；不阻塞 P3-S* / P3-C* 首批 PR 开干。

### 7.2 整域签核

| 项 | Owner | 日期 | 签核 |
|---|---|---|---|
| Services 架构 | | | ☐ |
| Core 架构 | | | ☐ |
| Console 产品 | | | ☐ |
| 安全/合规（含 BOSS） | | | ☐ |

### 7.3 联评通过后文档动作（service-ani）

- [ ] 更新本文件 §7 结论与 §5 决议栏
- [ ] 回填四份 `PHASE3-*-ACCEPTANCE-RECORD.md` 联评行
- [ ] 更新 `console-document-status-board.md`、`schema-completion-tracker.md`
- [ ] 通过的模块：详文「接口冻结规则」由 **规划** 改为 **合入跟踪**（仍等 YAML PR merged）
- [ ] `TASK-SVC-018` / `TASK-CORE-015` 写入 PR 链接

---

## 8. 联评后 handler 验通（ANI-main，非本场）

合入 YAML 后按域执行 curl（见各 `PHASE3-*-ACCEPTANCE-RECORD.md`），回填：

| 列 | 含义 |
|---|---|
| YAML 合入 | PR merged to ANI-main |
| Handler | 目标环境 curl 通过 |

参考：`tasks/phase2/PHASE2-HANDLER-ACCEPTANCE-RECORD.md`（Phase 2 同样式）。

---

## 附录 A — OpenAPI 草案路径一览

| # | 草案 | 域 |
|---|---|---|
| 1 | `../openapi-drafts/phase3/openapi-phase3-agent-session-draft.md` | AI |
| 2 | `../openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md` | AI |
| 3 | `../openapi-drafts/phase3/openapi-phase3-agent-audit-draft.md` | AI |
| 4 | `../openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md` | AI |
| 5 | `../openapi-drafts/phase3/openapi-phase3-mcp-tool-market-draft.md` | AI |
| 6 | `../openapi-drafts/phase3/openapi-phase3-ai-native-sandbox-security-draft.md` | AI |
| 7 | `../openapi-drafts/phase3/openapi-phase3-agent-orchestration-draft.md` | AI |
| 8 | `../openapi-drafts/phase3/openapi-phase3-kb-doc-intelligence-draft.md` | KB |
| 9 | `../openapi-drafts/phase3/openapi-phase3-kb-meeting-intelligence-draft.md` | KB |
| 10 | `../openapi-drafts/phase3/openapi-phase3-kb-video-intelligence-draft.md` | KB |
| 11 | `../openapi-drafts/phase3/openapi-phase3-model-encryption-draft.md` | Model |
| 12 | `../openapi-drafts/phase3/openapi-phase3-model-recommend-config-draft.md` | Model |
| 13 | `../openapi-drafts/phase3/openapi-phase3-model-usage-stats-draft.md` | Model |
| 14 | `../openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md` | Sec |
| 15 | `../openapi-drafts/phase3/openapi-phase3-netsec-policy-draft.md` | Sec |
| 16 | `../openapi-drafts/phase3/openapi-phase3-compliance-draft.md` | Sec |
| 17 | `../openapi-drafts/phase3/openapi-phase3-billing-export-draft.md` | Sec |

路径前缀：`docs/console-modules/`。

---

## 附录 B — 与 P0 / Phase 2 关系

| 阶段 | 状态 | 联评关系 |
|---|---|---|
| P0 YAML 草案 | **待定** | 不阻塞 Phase 3 联评；若 P0 先合入，Phase 3 PR 需 rebase |
| Phase 2 handler | 文档 ✅ / 运行时 ☐ | Phase 3 合入前建议至少 Sprint A1 路径稳定 |
| Phase 3 文档 | **17/17 ✅** | 本场主题 |

---

## 相关文件

- `docs/console-modules/openapi-drafts/phase3/openapi-phase3-domain-draft.md`
- `tasks/TASK-SVC-018`、`tasks/execution/CORE-TEAM-TASKS.md`（CORE-015）
- `tasks/execution/TASK-DEPENDENCY-MAP.md`
- `docs/console-modules/governance/module-delivery-workflow.md` §2.10（422 口径）
