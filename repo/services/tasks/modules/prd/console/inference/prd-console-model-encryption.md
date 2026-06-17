# PRD: Console model-encryption

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/inference/model-encryption.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-encryption-draft.md`

## 目标

- 定义模型 artifact 与 Core Encryption Key **绑定**的产品边界
- 输出 Services OpenAPI 草案（GET/PUT encryption binding）
- 为 TASK-SVC-018 模型增强子项提供合入 YAML 的前置材料

## 用户故事

- 作为模型管理员，我希望为 ready 模型绑定 SM4/AES 密钥并查看绑定状态
- 作为安全管理员，我希望密钥 CRUD 仍在 Core encryption 页，模型页只引用 key_id
- 作为开发者，我希望 key revoke 时 PUT 返回 422 `ENCRYPTION_KEY_REVOKED`

## 范围

- GET/PUT `/models/{id}/encryption`；Core keys 只读引用
- Console 绑定状态、算法选择、跳转 crypto-sm

## 非目标

- 重复 Core encryption keys CRUD
- 合入 ANI-main（本阶段仅文档）

## 成功标准

- [ ] 草案通过评审；422 与 model-center 一致
- [ ] 合入 YAML 后详文切换为冻结口径
