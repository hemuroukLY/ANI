# SPEC: Console gpu-inventory-ui

> Technical specification — Core handler 契约见 `tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` §TASK-CORE-003  
> Source: `tasks/modules/prd/console/compute/prd-console-gpu-inventory-ui.md`  
> Revised: 2026-06-17

## 1. Summary

GPU 设备清单与占用摘要的**只读** Console 子页；不等同于独立 GPU CRUD 资源域。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | `200 + GPUInventoryListResponse` | `scope:gpu-inventory:read` |
| GET | `/api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | `200 + GPUOccupancySummary` | `scope:gpu-inventory:read` |

### 2.3 Verified Schemas

- `GPUInventoryDevice`（`state`: available / allocated / unavailable / maintenance）
- `GPUInventoryListResponse`（`items`、`next_cursor`）
- `GPUOccupancySummary`（`total`、`allocated`、`available`、`avg_utilization_pct`）

## 3. Page Scope

- 设备列表筛选（vendor、state）、占用摘要卡片
- 跳转：`kind=gpu_container` 实例详情
- **Non-Goals**：GPU 分配/回收写操作（YAML 未声明）

## 4. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户认证 | 已登录 | `401` |
| 读权限 | `scope:gpu-inventory:read` | `403` |

## 5. 操作可用性矩阵

| 操作 | 只读用户 | 运维/管理员 |
|---|---|---|
| 查看清单 | ✅ | ✅ |
| 查看占用摘要 | ✅ | ✅ |
| GPU 分配/回收 | ❌ | ❌（无 API） |

## 6. 主维护源

- `docs/console-modules/compute/gpu-inventory-ui.md`
- 父模块：`docs/console-modules/compute/gpu-management.md`

## 7. Handler 验收（Core 团队）

```bash
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/gpu-inventory"
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/gpu-inventory/occupancy"
```

OpenAPI 已声明 ≠ handler 已实现。
