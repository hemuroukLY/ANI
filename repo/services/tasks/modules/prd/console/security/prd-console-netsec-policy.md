# PRD: Console netsec-policy

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/security/netsec-policy.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-netsec-policy-draft.md`

## 目标

- 定义 **网络安全策略** 与安全组分资源的 Core Networks 契约
- 输出 `/networks/policies` CRUD 草案
- 为 TASK-CORE-015 / Phase 3 §4 提供合入 YAML 前置材料

## 用户故事

- 作为网络管理员，我希望配置 VPC 级 egress 与微隔离策略
- 作为运维，我希望策略与安全组分 path，避免与 SG rules 混用
- 作为开发者，我希望 create 必填 idempotency_key 与 INVALID_POLICY_RULE 422

## 范围

- policies list/create/get/put/delete
- policy_type: egress_control | micro_segmentation | service_mesh

## 非目标

- CNI 实现细节
- 合入 ANI-main

## 成功标准

- [ ] 与 security-groups 边界清晰
- [ ] 合入 YAML 后详文切换为冻结口径
