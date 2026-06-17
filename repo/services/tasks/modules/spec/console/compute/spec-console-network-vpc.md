# SPEC: Console 网络 VPC

> Revised: 2026-06-17

## 1. Summary

Core NetworkVPC CRUD 子模块。

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/networks/vpcs` | `listNetworkVPCs` | `200 + NetworkVPCListResponse` | `401`,`403` |
| POST | `/api/v1/networks/vpcs` | `createNetworkVPC` | `201 + NetworkVPC` | `400`,`401`,`403` |
| GET | `/api/v1/networks/vpcs/{vpc_id}` | `getNetworkVPC` | `200 + NetworkVPC` | `401`,`403`,`404` |
| DELETE | `/api/v1/networks/vpcs/{vpc_id}` | `deleteNetworkVPC` | `200 + NetworkVPC` | `401`,`403`,`404` |

Schemas: `NetworkVPC`, `CreateNetworkVPCRequest`（含 idempotency_key）

## 3. Page Scope / Non-Goals

见详文。

## 4. References

- `docs/console-modules/compute/network/vpc.md`
