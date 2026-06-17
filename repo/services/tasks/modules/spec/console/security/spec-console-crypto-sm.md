# SPEC: Console crypto-sm

> Technical specification derived from: `tasks/modules/prd/console/security/prd-console-crypto-sm.md`  
> Revised: 2026-06-17

## 1. Summary

### 1.1 What This SPEC Covers

收口 Console **安全与身份 / 国密加解密** 的正式边界、字段映射口径与待补依赖。只对齐当前 OpenAPI 已冻结能力，不重新设计 API。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/security/prd-console-crypto-sm.md`
- 主维护源: `docs/console-modules/security/crypto-sm.md`

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/encryption/keys` | `listEncryptionKeys` | `200 + EncryptionKeyListResponse` | `scope:encryption:read` |
| POST | `/api/v1/encryption/keys` | `createEncryptionKey` | `201 + EncryptionKey` | `scope:encryption:create` |
| POST | `/api/v1/encryption/keys/{key_id}/rotate` | `rotateEncryptionKey` | `200 + EncryptionKeyRotationResponse` | `scope:encryption:manage` |
| POST | `/api/v1/encryption/seal` | `sealEncryptionObject` | `200 + EncryptionSealResponse` | `scope:encryption:use` |


## 3. Page Scope

### 3.1 Page Responsibilities

- **加密密钥生命周期** 与 seal/unseal 操作的 Console 页（含国密 SM4 规划语义）；local dev profile 与真实 KMS 响应需保持兼容（见 YAML description）。…

### 3.2 Non-Goals

- 不改变未冻结的规划路径
- 不把 handler stub 标注为已实现
- 非 Console API 页不自造 REST 契约

## 4. 创建前置条件与操作矩阵

见主维护源 `docs/console-modules/security/crypto-sm.md`（§创建前置条件、§操作可用性矩阵）。

## 5. 待补边界

- 真实 SM4/HSM provider — handler 待接入
- 与 `model-encryption.md` 模型绑定联动

## 6. 主维护源

- `docs/console-modules/security/crypto-sm.md`
- 流程：`docs/console-modules/governance/module-delivery-workflow.md`

OpenAPI 已声明 ≠ handler 已实现。
