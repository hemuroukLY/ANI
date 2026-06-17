# OpenAPI Schema 补全跟踪表

> 对齐 `docs/console-modules/governance/YAML-EXPANSION-SUMMARY-2026-06-17.md` + `GAP-REPORT-2026-06-17.md`，跟踪 OpenAPI 声明与 handler 实现。  
> 更新频率：每次 YAML 扩充或 handler 合并后刷新。

**最后更新**：2026-06-17（65 项 backlog 文档收口 + P0 YAML 草案）

## P0 YAML 草案（**待定**，暂缓评审）

| 缺口 | 草案路径 | TASK | 状态 |
|---|---|---|---|
| 告警事件 list | `GET /observability/alerts` | TASK-CORE-014 | **待定** |
| 任务 list | `GET /tasks` | TASK-CORE-014 | **待定** |
| 推理 events | `GET .../inference-services/{id}/events` | TASK-SVC-017 | **待定** |

草案全文：`docs/console-modules/openapi-drafts/p0/openapi-p0-yaml-draft.md`（评审恢复后再合入 ANI-main）

## Phase 2 Handler 文档验收

| 交付物 | 路径 |
|---|---|
| 验收记录（curl + 签核） | `tasks/phase2/PHASE2-HANDLER-ACCEPTANCE-RECORD.md` |
| Services 指南 | `tasks/execution/SERVICES-HANDLER-IMPLEMENTATION-GUIDE.md` |
| Core 指南 | `tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` |

文档验收：**✅**（2026-06-17）· 运行时 Handler：**☐** 待 PR

## Phase 3 整域规划（文档层）

`docs/console-modules/openapi-drafts/phase3/openapi-phase3-domain-draft.md` — AI 原生、知识智能、模型/安全待补路径

## 汇总

| 来源 | OpenAPI ops（active） | OpenAPI 层 schema | Handler 实现 | 说明 |
|---|---|---|---|---|
| Core `v1.yaml` | **111** | ✅ 111/111 挂 2xx/422 | ⚠️ 大量 stub/demo | Phase 2 **+19** ops |
| Services `services/v1.yaml` | **41** active + 19 deprecated | ✅ active 全挂 schema | ❌ 均为 501 stub | Phase 2 **+12** active ops |
| Console 模块文档 | **86** 套详文（21 基线 + 65 backlog） | 结构 100% Core 标准 | ⚠️ 依赖 handler | 2026-06-17 全部文档已补 |

> OpenAPI 层「完成率 100%」≠ handler 已实现。Phase 4 才验通 handler。

## Phase 2 新增 Core 路径（handler 待实现）

| 路径组 | ops | TASK 引用 |
|---|---|---|
| `/gpu-inventory*` | 2 | TASK-CORE-003 |
| `/instances/{id}/logs\|events\|metrics\|exec` | 4 | TASK-CORE-005 |
| `/instances/{id}/security-events` | 1 | TASK-CORE-006 |
| `/sandbox-templates` | 1 | TASK-CORE-006 |
| `/networks/routes` | 2 | TASK-CORE-007 |
| `/volumes/{id}/snapshots` | 2 | TASK-CORE-008 |
| `/buckets`、`/objects/upload`、`/objects/{id}/download` | 4 | TASK-CORE-009 |
| `/filesystems/{id}/mount-targets` | 1 | TASK-CORE-010 |
| `/vector-stores/{id}/documents` | 1 | TASK-CORE-011 |
| `/k8s-clusters/{id}/workloads` | 1 | TASK-CORE-012 |

## Phase 2 新增 Services 路径（handler 待实现）

| 路径组 | ops | TASK 引用 |
|---|---|---|
| `/inference-services/{id}/logs\|test\|policies` | 3 | TASK-SVC-013 |
| `/knowledge-bases/{id}/citations\|sessions\|permissions` | 3 | TASK-SVC-014 |
| `/tenant/members/{id}` GET | 1 | TASK-SVC-015 |
| `/tenant/roles/{id}` PUT | 1 | TASK-SVC-015 |
| `/tenant/webhooks/{id}/deliveries` | 1 | TASK-SVC-015 |
| `/integrations*` | 3 | TASK-SVC-016 |

## Core 422 扩充（Phase 2）

| 路径 | YAML | Handler 验通 |
|---|---|---|
| `POST /k8s-clusters` | ✅ | ❌ → TASK-CORE-012 |
| `POST /vector-stores/{id}/search` | ✅ | ❌ → TASK-CORE-011 |
| `POST /instances`、lifecycle | ✅ | ⚠️ → TASK-CORE-002 |

## 仍 stub / 未实现（非 Phase 2 范围）

| 路径组 | 状态 | 优先级 |
|---|---|---|
| Branding / Tasks | Core stub | P3 |
| `/metering/*` handler 验通 | 部分 | P2 |
| `GET /api/v1/alerts` | 见 TASK-CORE-014 草案 | P0 |

## 文档层与 YAML 对齐项（Phase 2 后）

| 项 | 文档 | YAML | Handler |
|---|---|---|---|
| `/instances?kind=sandbox` | ✅ | ✅ | 文档 + TASK-CORE-001 |
| `POST .../lifecycle` 422 | ✅ | ✅ | 文档 + TASK-CORE-002 |
| `/gpu-inventory*` | ✅ | ✅ | 文档 + TASK-CORE-003 |
| instances 观测子路径 | ✅ | ✅ | 文档 + TASK-CORE-005 |
| k8s/search 422 | ✅ | ✅ | ❌ |
| tenant webhook 422 | ✅ | ✅ | ❌ |
| deprecated gpu/sandbox svc | 拒绝 | deprecated | 不实现 |

## 下一步

1. ~~Phase 2 handler 文档验收~~ ✅ → `tasks/phase2/PHASE2-HANDLER-ACCEPTANCE-RECORD.md`
2. **运行时 handler 验通**：Core/Services 按指南 PR 后回填本表 Handler 列
3. **P0 YAML 评审**：**待定** — `openapi-p0-yaml-draft.md`
4. **Phase 3 AI 原生 ×7**：详文 + 草案 ✅ → `ai-native/README.md`；YAML 合入 ☐
5. **Phase 3 知识库智能 ×3**：详文 + 草案 ✅ → `knowledge/README.md`；YAML 合入 ☐
6. **Phase 3 模型增强 ×3**：详文 + 草案 ✅ → `inference/README.md`；YAML 合入 ☐
7. **Phase 3 安全/合规 ×4**：详文 + 草案 ✅ → `security/README.md`；YAML 合入 ☐
8. **Phase 3 文档层**：整域 17 套详化草案 ✅
9. **Phase 3 整域联评**：`tasks/phase3/PHASE3-JOINT-REVIEW-AGENDA.md`（待召开）→ YAML 合入
