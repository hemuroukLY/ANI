# PRD: Console api-key-audit

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/tenant/api-key-audit.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md`

## 目标

- 定义 **API Key 调用审计** 与 Key CRUD、平台审计的产品边界
- 输出 Core Auth OpenAPI 草案（list audit-events）
- 为 TASK-CORE-015 / Phase 3 §4 提供合入 YAML 前置材料

## 用户故事

- 作为 Key 管理员，我希望查看某 Key 的调用时间线（IP、路径、状态码）
- 作为安全管理员，我希望审计 API 不返回 key_value 明文
- 作为开发者，我希望 path 挂在 `/auth/api-keys/{id}/audit-events`

## 范围

- GET audit-events list（3a）；单条 get（3b）
- 从 api-key-management 跳转

## 非目标

- BOSS 取证导出
- 合入 ANI-main（本阶段仅文档）

## 成功标准

- [ ] 与 audit-log 事件模型可映射
- [ ] 合入 YAML 后详文切换为冻结口径
