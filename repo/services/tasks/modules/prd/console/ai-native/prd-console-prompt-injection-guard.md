# PRD: Console prompt-injection-guard

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/ai-native/prompt-injection-guard.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md`

## 目标

- 定义 **Prompt 注入防护**与推理 **RPM/scope 策略**的产品与 API 边界
- 输出 guard-policies CRUD（3a）+ evaluate / inference binding（3b）评审草案
- 与 Agent Audit 的 `policy_hit` 事件对齐

## 用户故事

- 作为安全管理员，我希望对 Agent 与推理入口启用 input/output 注入扫描，并选择 block/warn/log
- 作为运维，我希望在不上线的情况下用样例文本演练策略是否命中
- 作为开发者，我希望 guard 策略独立于 `InferenceServicePolicyUpdateRequest`

## 范围

- GET/PUT `/agent/guard-policies`
- builtin + custom 规则矩阵 UI
- RBAC：`agent-guard:read|write`
- 3b：evaluate、inference guard-binding

## 非目标

- 检测模型训练/更新（BOSS）
- 合并进 inference rate limit PUT
- 本阶段合入 ANI-main

## 成功标准

- [ ] 评审确认与 rate-limit-policy 分资源
- [ ] Audit policy_hit 字段约定一致
- [ ] 合入 YAML 后详文切换冻结口径
