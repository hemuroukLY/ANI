# 国密加解密

## 页面定位

**加密密钥生命周期** 与 seal/unseal 操作的 Console 页（含国密 SM4 规划语义）；local dev profile 与真实 KMS 响应需保持兼容（见 YAML description）。

## 文档管理规则

- 本文是 **国密加解密** 的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 不得把规划路径或 handler stub 写成已实现
- PRD/SPEC 同步规则见 `docs/console-modules/governance/module-delivery-workflow.md` §3.5

## Core 层要求

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/encryption/keys` | `listEncryptionKeys` |
| POST | `/api/v1/encryption/keys` | `createEncryptionKey` |
| GET | `/api/v1/encryption/keys/{key_id}` | `getEncryptionKey` |
| DELETE | `/api/v1/encryption/keys/{key_id}` | `deleteEncryptionKey` |
| POST | `/api/v1/encryption/keys/{key_id}/rotate` | `rotateEncryptionKey` |
| POST | `/api/v1/encryption/keys/{key_id}/revoke` | `revokeEncryptionKey` |
| POST | `/api/v1/encryption/seal` | `sealEncryptionObject` |
| POST | `/api/v1/encryption/unseal-token` | `createEncryptionUnsealToken` |

Schema：`EncryptionKey`、`EncryptionKeyListResponse`、`EncryptionSealRequest/Response` 等（以 YAML 为准）。

## Services 层要求

无 Services 路径。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 模块读/写权限 | 对应 RBAC scope | `403 FORBIDDEN` |

写操作 POST/PUT 的 `idempotency_key` 与 422 口径 — **待 YAML 冻结后按 ../governance/module-delivery-workflow.md §2.10 补充**。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员/运维 | 说明 |
|---|---|---|---|
| 查看密钥 | 可用 | 可用 | list/get |
| 创建/轮换/吊销 | 不可用 | 管理员 | write/manage |
| seal/unseal | 不可用 | 运维 | encryption:use |

## 页面职责

- 占位 UI + 明确 YAML/OpenAPI 缺口（若适用）
- 跳转关联模块（见「相关模块」）
- 不把 BOSS/平台运营能力写入 Console 冻结契约

## 接口冻结规则

### `GET /api/v1/encryption/keys`

- operationId: `listEncryptionKeys`
- 成功：`200 + EncryptionKeyListResponse`
- RBAC：`scope:encryption:read`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文

### `POST /api/v1/encryption/keys`

- operationId: `createEncryptionKey`
- 成功：`201 + EncryptionKey`
- RBAC：`scope:encryption:create`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文

### `POST /api/v1/encryption/keys/{key_id}/rotate`

- operationId: `rotateEncryptionKey`
- 成功：`200 + EncryptionKeyRotationResponse`
- RBAC：`scope:encryption:manage`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文

### `POST /api/v1/encryption/seal`

- operationId: `sealEncryptionObject`
- 成功：`200 + EncryptionSealResponse`
- RBAC：`scope:encryption:use`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文


## 待补边界

- 真实 SM4/HSM provider — handler 待接入
- 与 `model-encryption.md` 模型绑定联动

## 相关模块

- `security/secrets-management.md`
- `inference/model-encryption.md`

## 验收标准

- [ ] 路径与 OpenAPI 权威源一致（或明确 TODO-YAML / N/A）
- [ ] 正文不把 handler stub 写成已实现
- [ ] 含创建前置条件与操作可用性矩阵
- [ ] PRD/SPEC/HTML 摘要已同步
