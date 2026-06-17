# SPEC: Console secrets-management

> Technical specification derived from: `tasks/modules/prd/console/security/prd-console-secrets-management.md`  
> Revised: 2026-06-17

## 1. Summary

### 1.1 What This SPEC Covers

收口 Console **安全与身份 / 密钥管理** 的正式边界、字段映射口径与待补依赖。只对齐当前 OpenAPI 已冻结能力，不重新设计 API。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/security/prd-console-secrets-management.md`
- 主维护源: `docs/console-modules/security/secrets-management.md`

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/secrets` | `listSecrets` | `200 + SecretListResponse` | `scope:secrets:read` |
| POST | `/api/v1/secrets` | `createSecret` | `201 + Secret` | `scope:secrets:create` |
| GET | `/api/v1/secrets/{secret_id}` | `getSecret` | `200 + Secret` | `scope:secrets:read` |
| DELETE | `/api/v1/secrets/{secret_id}` | `deleteSecret` | `200 + Secret` | `scope:secrets:delete` |
| POST | `/api/v1/secrets/{secret_id}/bindings` | `bindSecret` | `201 + SecretBinding` | `scope:secrets:bind` |


## 3. Page Scope

### 3.1 Page Responsibilities

- 租户 **Secret 元数据** CRUD 与工作负载绑定页；API **不返回明文值**（见 YAML description）。…

### 3.2 Non-Goals

- 不改变未冻结的规划路径
- 不把 handler stub 标注为已实现
- 非 Console API 页不自造 REST 契约

## 4. 创建前置条件与操作矩阵

见主维护源 `docs/console-modules/security/secrets-management.md`（§创建前置条件、§操作可用性矩阵）。

## 5. 待补边界

- Secret 值轮换 PUT — **YAML 未声明**
- K8s provider 真实注入 — handler 边界

## 6. 主维护源

- `docs/console-modules/security/secrets-management.md`
- 流程：`docs/console-modules/governance/module-delivery-workflow.md`

OpenAPI 已声明 ≠ handler 已实现。
