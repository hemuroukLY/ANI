# BOSS audit 域 Phase 0 GAP 摘要

> **用途**：audit 域 Phase 2 YAML 批次输入；完整 BOSS GAP 待全模块汇总。  
> **日期**：2026-06-17  
> **权威源**：`docs/boss-modules/audit/*.md` + `v1.yaml`

## 已确认：YAML 无 audit 契约

| 检查项 | 结论 |
|---|---|
| `/audit*` paths | **不存在** |
| `AuditEvent` schema | **不存在** |
| Tag `Audit` | 仅在 `info.description` 提及 |

## 租户 path 只读参考（不可冒充 BOSS 大盘）

| 模块 | 路径 | operationId | 用途 |
|---|---|---|---|
| API Key 审计 | `GET/POST /auth/api-keys` | `listAPIKeys` / `createAPIKey` | Key 元数据 |
| API Key 审计 | `DELETE /auth/api-keys/{key_id}` | `revokeAPIKey` | 吊销事件源 |
| 推理调用审计 | `POST /metering/token-usage` | `reportTokenUsage` | **写入**计量，非 audit list |
| 推理调用审计 | `GET /metering/usage` | **未声明** | 租户用量聚合参考 |
| 合规导出 | `GET /tasks/{task_id}` | **未声明** | 异步任务单查 |

## Phase 2 待新增（BOSS 平台 · TODO-YAML）

| 优先级 | 模块 | 建议能力 | 归属 |
|---|---|---|---|
| P1 | 平台审计日志 | `GET /audit/events`（跨租户 list + 筛选） | Core |
| P1 | API Key 审计 | 平台 Key 调用事件 list / aggregate | Core Auth |
| P1 | 推理调用审计 | 平台 inference audit list（无 prompt） | Core 或 Services |
| P1 | 共用 | `GET /tasks` 跨租户 list | Core |
| P2 | 合规导出 | `POST /compliance/exports` + export job schema | Core |
| P2 | 合规导出 | `AsyncTask.task_type` 扩展 `compliance.export` | Core |

## 架构待决（Phase 2 前必须拍板）

1. 平台 RBAC scope 命名（`scope:platform:audit:read` 等）
2. `AuditEvent` / export job schema 归属 Core components
3. Console 草案 path（api-key audit-events、svc compliance）与 BOSS 平台 path 是否统一前缀

## 文档状态

audit 域 **4 模块** 已满配（Core）+ PRD/SPEC/HTML 已同步。  
门禁：`python3 scripts/validate-boss-audit-gate.py`
