# PRD: Console tenant-member-detail

详文：`docs/console-modules/tenant/tenant-member-detail.md` · Phase 2 YAML 已声明项文档收口（2026-06-17）。

## 目标

- 为 **Services / Tenant / 成员详情** 提供可维护的产品边界说明
- 明确 `GET .../members/{member_id}` 只读、PATCH 成员未声明

## 用户故事（规划）

- 作为租户成员，我希望查看成员 profile 与角色，以便了解团队构成
- 作为租户管理员，我希望从详情页跳转移除成员或角色权限配置

## 范围

- 成员只读详情、跳转 tenant-management / role-permission-edit — 见主维护源
- 不修改 ANI-main OpenAPI（由 Services 团队按 TASK-SVC-015 推进）

## 非目标

- 不自造 PATCH member 为冻结能力
- 成员活动审计、与 user-management 合并方案（TODO-YAML）
