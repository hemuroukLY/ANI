# BOSS metering 域 Phase 0 GAP 摘要

> **日期**：2026-06-17 · **Phase 2**：✅ P1 已合入 YAML

## 已声明 path

| 路径 | 说明 |
|---|---|
| `GET /api/v1/metering/usage` | `getMeteringUsage` · 租户 scope |
| `GET /api/v1/metering/usage/platform` | `getPlatformMeteringUsage` · 平台 aggregate |
| `POST /api/v1/metering/token-usage` | `reportTokenUsage` |

### resource_type enum

| resource_type | 状态 |
|---|---|
| `instance_*` / `token_*` | ✅ YAML |
| `storage_gb_days` | ✅ **ADDED-TO-YAML** |
| `kb_query_count` | ✅ **ADDED-TO-YAML** |

## 待 P2

| 能力 | Phase |
|---|---|
| 平台 CSV 导出 | P2 |

## 文档状态

metering 域 **7 模块** 满配 + Phase 2 P1 YAML 已合入。
