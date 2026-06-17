# 模型 — 加密与密钥绑定

## 页面定位

模型 artifact **加密存储与密钥绑定** 配置页；Core 密钥生命周期见 `security/crypto-sm.md`，本页聚焦 Services 模型 ↔ `encryption_key_id` 绑定。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-model-encryption-draft.md`。

## 文档管理规则

- 本文是 **模型 — 加密与密钥绑定** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-encryption-draft.md`
- 合入后 Services 权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- Core 只读：`ANI-main/repo/api/openapi/v1.yaml`（encryption/keys）
- PRD/SPEC：`tasks/modules/prd/console/inference/prd-console-model-encryption.md`、`tasks/modules/spec/console/inference/spec-console-model-encryption.md`
- TASK：`TASK-SVC-018` 子项（模型增强）

## Core 层要求（只读引用）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/encryption/keys` | 选 key — `crypto-sm.md` |
| GET | `/api/v1/encryption/keys/{key_id}` | key 详情 |

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/models/{model_id}/encryption` | `getModelEncryptionBinding` | 3a |
| PUT | `/api/v1/svc/models/{model_id}/encryption` | `updateModelEncryptionBinding` | 3a |

Schema（草案）：`ModelEncryptionBinding`、`UpdateModelEncryptionRequest`。

RBAC（草案）：`scope:models:read`、`scope:models:write`。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 读/写权限 | models read/write | `403` |
| PUT 幂等键 | 必填 | `400` |
| model 存在 | 有效 id | `404` |
| model 就绪 | `status=ready`（绑定时） | `422`（建议语义：`MODEL_NOT_READY`；**待 YAML 声明后定稿**） |
| key 有效 | 未 revoke | `422`（建议语义：`ENCRYPTION_KEY_REVOKED`；**待 YAML 声明后定稿**） |

## 页面职责

- 展示绑定状态、`encryption_key_id`、algorithm（aes256_gcm / sm4_gcm）
- 绑定/解绑/轮换（大 artifact 可能 202 AsyncTask）
- 跳转 **国密加解密** 密钥管理
- `ModelVersion.is_encrypted` 由 binding 成功后回写（UI 只读展示）

## 操作可用性矩阵

| 操作 | 只读用户 | 模型管理员 | YAML 合入后 |
|---|---|---|---|
| 查看绑定 | 可用 | 可用 | GET |
| 绑定/解绑 | 不可用 | 可用 | PUT |
| 密钥 CRUD | 跳转 | 跳转 | Core encryption |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-model-encryption-draft.md` §4。

### `PUT .../models/{model_id}/encryption`（规划）

- 成功：`200 + ModelEncryptionBinding` 或 `202 + AsyncTask`（re-encrypt）
- 错误：`400`、`404`、`422`

## 待补边界

- binding 粒度：model 级 vs version 级 — `version_id` 可选
- re-encrypt 异步 task_type — 建议 `model.reencrypt`
- HSM/SM4 provider — `crypto-sm.md`

## 相关模块

- `inference/model-center.md`
- `security/crypto-sm.md`

## 验收标准

- [ ] 不重复 Core keys CRUD
- [ ] 422 与 model-center 口径一致
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
