# BOSS ops 域 Phase 0 GAP 摘要

> **日期**：2026-06-17 · **详文**：`docs/boss-modules/ops/*.md`（8 模块）

## 已声明 path（只读参考）

| 路径 | operationId | 模块 |
|---|---|---|
| `GET /api/v1/gpu-inventory` | `listGPUInventory` | gpu-pool-management |
| `GET /api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | gpu-pool-management |
| `GET /api/v1/volumes` | `listStorageVolumes` | storage-infrastructure、platform-resource-pool |
| `GET /api/v1/volumes/{volume_id}` | `getStorageVolume` | storage-infrastructure |
| `GET /api/v1/networks/vpcs` | `listNetworkVPCs` | network-infrastructure、platform-resource-pool |
| `GET /api/v1/networks/subnets` | `listNetworkSubnets` | network-infrastructure |
| `GET /api/v1/networks/security-groups` | `listNetworkSecurityGroups` | network-infrastructure |
| `GET /api/v1/networks/load-balancers` | `listNetworkLoadBalancers` | network-infrastructure |
| Registry paths | 见 registry 三模块 | registry P2 |

## Phase 0 已修正（2026-06-17）

| 模块 | 问题 | 修正 |
|---|---|---|
| network-infrastructure | 冻结表 `listVPCs` 等与 YAML 不符 | → `listNetworkVPCs` 等 |
| storage-infrastructure | `listVolumes` | → `listStorageVolumes` |
| platform-resource-pool | `listVolumes` / `listVPCs` | → `listStorageVolumes` / `listNetworkVPCs` |
| network-infrastructure | 400 示例不得绑定租户 list 已冻结错误码 | 已改为平台 aggregate TODO-YAML 目标示例 |

## Core 缺口（BOSS 平台 · TODO-YAML）

| 优先级 | 能力 | 模块 |
|---|---|---|
| **P1** | 平台 resource-pool aggregate | platform-resource-pool |
| **P1** | 平台 nodes list | node-status |
| **P1** | 平台 storage/network aggregate | storage、network |
| **P2** | Registry quota / scan / GC 平台 API | registry 三模块 |

## 架构待决

- 平台 infrastructure aggregate 统一前缀（`/platform/*` vs resource 级 aggregate）
- nodes list 与 gpu-inventory 去重边界
- Registry 平台 API 与 Harbor 集成归属

门禁：registry P2 — `python3 scripts/validate-boss-registry-p2-gate.py`
