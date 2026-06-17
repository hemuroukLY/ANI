# SPEC: Console 首页 GPU 利用率

> Revised: 2026-06-17

## 1. Summary

Core GPU 库存与占用摘要页。

## 2. Frozen Facts

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | `200 + GPUInventoryListResponse` | scope:gpu-inventory:read |
| GET | `/api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | `200 + GPUOccupancySummary` | scope:gpu-inventory:read |
| GET | `/api/v1/observability/query` | `queryObservability` | `200 + ObservabilityQueryResponse` | scope:observability:read |

## 3. Page Scope / Non-Goals

见详文；不写 gpu-inventory/{id}/metrics 为冻结。

## 4. References

- `docs/console-modules/home/home-gpu-utilization.md`
