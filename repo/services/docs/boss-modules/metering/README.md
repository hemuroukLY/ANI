# 平台计量与结算（BOSS）

> **PRD/SPEC/HTML 同步**：tenant / health / metering 三域 21 模块已于 2026-06-17 与满配详文对齐。

平台运营/财务域：按资源类型查看 **跨租户** 用量排行、趋势与导出；单页聚焦一种计量视角。

| 页面 | 详文 | Core resource_type（冻结） |
|---|---|---|
| 平台 GPU-Hours | [platform-gpu-hours.md](platform-gpu-hours.md) | `instance_gpu_seconds` |
| 平台 CPU-Hours | [platform-cpu-hours.md](platform-cpu-hours.md) | `instance_cpu_seconds` |
| 平台 Memory-GBHours | [platform-memory-gbhours.md](platform-memory-gbhours.md) | `instance_memory_gib_seconds` |
| 平台 Storage-GBDays | [platform-storage-gbdays.md](platform-storage-gbdays.md) | **待 YAML** |
| 平台 Input Tokens | [platform-input-tokens.md](platform-input-tokens.md) | `token_input` |
| 平台 Output Tokens | [platform-output-tokens.md](platform-output-tokens.md) | `token_output` |
| 平台 KB Queries | [platform-kb-queries.md](platform-kb-queries.md) | **待 YAML** |

## 聚合入口

跨指标切换与租户钻取见 [`../tenant/tenant-usage-billing.md`](../tenant/tenant-usage-billing.md)。

Console 对照：[`usage-report.md`](../../console-modules/tenant/usage-report.md)（单租户 `GET /metering/usage`）。

## Core 已冻结（租户上下文 · 只读参考）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/metering/usage` | 单租户用量；`group_by`: resource_type / az / day / hour |
| POST | `/api/v1/metering/token-usage` | Token 上报（`reportTokenUsage`）；非查询 |

`MeteringUsageRecord.resource_type` 冻结枚举：`instance_cpu_seconds`、`instance_memory_gib_seconds`、`instance_gpu_seconds`、`token_input`、`token_output`、`token_total`。

## BOSS 待补（TODO-YAML）

- 平台跨租户用量 list/aggregate（含 `tenant_id` group_by）
- `storage_*`、`kb_query_*` 等资源类型入枚举
- 平台 CSV 导出、账单金额字段

## Phase 1 / 满配状态

**✅ 本域 7 模块已达 Console 满配（Core）**（2026-06-17，阶段 1）。检查清单：[`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)。
