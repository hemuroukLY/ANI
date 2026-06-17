# YAML EXPANSION SUMMARY — BOSS Phase 2（P1 + P1-fix）

> **日期**：2026-06-17 · **状态**：✅ 严格门禁 120/120 ALL PASS  
> **合入**：`apply-boss-phase2-yaml.py` + `apply-boss-phase2-fixes.py`  
> **同步**：`sync-boss-phase2-full.py`（详文 + HTML）  
> **门禁**：`validate-boss-phase2-yaml-gate.py`

## 架构拍板

| 决策项 | 结论 |
|---|---|
| 平台 aggregate | `/platform/*` |
| 租户治理 | `/tenants*` + `/platform/tenant-admins` |
| 审计 | `/audit/events*` + `/platform/audit/*` |
| 可观测 | `/observability/{logs,traces,alert-events}` |
| 租户 task 单查 | `GET /tasks/{task_id}` → `scope:tasks:read` |
| 平台 task list | `GET /tasks` → `scope:platform:tasks:read` |
| Metering enum | `storage_gb_days` · `kb_query_count` |

## 新增到 Core v1.yaml

### 批次 P1（23 path · 30 op）

见首版表格：`/tenants*` · `/tasks` · `/metering/usage/platform` · `/audit/*` · `/platform/webhooks*` · `/platform/overview|resource-pools|networks|storage|capacity|nodes|health` · `/observability/alert-events*` · `/observability/traces*`

### 批次 P1-fix（+11 path · +13 op）

| 路径 | HTTP | operationId | 模块 |
|---|---|---|---|
| /observability/logs | GET | listObservabilityLogs | platform-logs |
| /platform/pending | GET | getPlatformPendingSummary | platform-alerts-pending |
| /platform/tenant-admins | GET | listTenantAdmins | tenant-admin |
| /tenants/{tenant_id}/admin | GET/PUT | getTenantAdmin/updateTenantAdmin | tenant-admin |
| /tenants/{tenant_id}/admin/reset | POST | resetTenantAdmin | tenant-admin |
| /platform/monitoring/gpu | GET | getPlatformGPUMonitoring | gpu-monitoring |
| /platform/monitoring/inference | GET | getPlatformInferenceMonitoring | inference-monitoring |
| /platform/monitoring/kb | GET | getPlatformKBMonitoring | kb-monitoring |
| /platform/trends/gpu | GET | getPlatformGPUTrend | gpu-pool-trend |
| /platform/trends/inference | GET | getPlatformInferenceTrend | inference-ops-trend |
| /platform/trends/kb | GET | getPlatformKBTrend | kb-ops-trend |

### 补丁（非新 path）

| 项 | 修正 |
|---|---|
| GET /tasks/{task_id} | `scope:tasks:read` + security + 401/403 |
| GET /metering/usage | security + getMeteringUsage |
| DELETE /platform/webhooks/{id} | +422 PreconditionFailed |
| 全部 Phase2 op | 显式 `security: [BearerAuth, ApiKeyAuth]` |

## Schema 合计

P1：39 个 · P1-fix：+13 个 · **合计 52 个** BOSS Phase2 schema

## 新增到 services/v1.yaml

0 条（平台 BOSS API 归属 Core）

## 刻意未合入（正确边界）

| 类别 | Phase |
|---|---|
| compliance export | P2 |
| branding PATCH | P2 |
| Registry 平台 API | P2 |
| task retry/cancel | P2 |
| maint-skills / incident | P2 |
| ops-system-integration 平台 aggregate | P2 |
| settings 交付安装 REST | N/A |

## 验收命令

```bash
python3 scripts/apply-boss-phase2-yaml.py
python3 scripts/apply-boss-phase2-fixes.py
python3 scripts/sync-boss-phase2-full.py
python3 scripts/validate-boss-phase2-yaml-gate.py   # 120/120
python3 scripts/validate-boss-audit-gate.py
python3 scripts/validate-boss-integration-gate.py
```

## Phase 3

待用户指令；输入为本 SUMMARY + 扩充后 `v1.yaml`。
