# BOSS settings 域 Phase 0 GAP 摘要

> **日期**：2026-06-17 · **详文**：`docs/boss-modules/settings/*.md`

## 交付类（10 模块 · REST N/A）

| 模块 | 实际操作面 | OpenAPI |
|---|---|---|
| ani-installer / baremetal / vm / hw / gpu-driver / ca / ip-https / offline / acceptance | `ani-installer` CLI/TUI | **TODO-YAML: N/A** |
| upgrade-backup（Phase 1） | `ani-installer upgrade` + 备份 SOP | 平台 BOSS patch UI — **TODO-YAML P2** |

**不得** 自造 `/api/v1/boss/*` 或 `/api/v1/installer/*`。

## Core 已声明（只读参考）

| 路径 | operationId | 模块 | 用法 |
|---|---|---|---|
| `GET /api/v1/branding` | **未声明** | brand-config | 公开只读；`security: []` |
| `GET/POST /api/v1/k8s-clusters` | listK8sClusters / createK8sCluster | k8s-attach | **注册后**只读展示；attach 本身走 CLI |
| `GET/DELETE /api/v1/k8s-clusters/{cluster_id}` | getK8sCluster / deleteK8sCluster | k8s-attach | 登记详情/删除 |

## Core 缺口

| 项 | 结论 |
|---|---|
| `PATCH /api/v1/branding` | **不存在** — brand-config **TODO-YAML P2** |
| 平台 patch upload / dry-run / history | **不存在** — upgrade-backup **TODO-YAML P2** |
| installer / preflight / deploy REST | **不存在** — 正确（N/A） |
| CA / offline / acceptance REST | **不存在** — 正确（N/A） |

## 边界提醒

- `POST /api/v1/k8s-clusters`（createK8sCluster）**≠** `ani-installer attach-k8s`
- `POST /api/v1/k8s-clusters/{cluster_id}/upgrade` 为**租户集群升级**，**≠** 平台 `.anipatch` 升级页
- Console **无** 交付域对等页

## Phase 2 建议批次

| 优先级 | 能力 |
|---|---|
| **P2** | `PATCH /api/v1/branding` + Logo 上传 + RBAC |
| **P2** | 平台 patch BOSS API（upload / dry-run / SSE history） |

门禁：`python3 scripts/validate-boss-settings-gate.py`
