# PRD: Console kb-permissions

详文：`docs/console-modules/knowledge/kb-permissions.md` · Phase 2 YAML 已声明项文档收口（2026-06-17）。

## 目标

- 为 **Services / KnowledgeBases / 权限管理** 提供可维护的产品边界说明
- 明确 `PUT .../permissions` 已冻结、document 级 ACL 非当前范围

## 用户故事（规划）

- 作为 KB 编辑者，我希望配置 readers/editors，以便控制谁可读可写
- 作为只读用户，我希望查看权限但无法保存变更

## 范围

- 权限矩阵 UI、幂等 PUT、只读用户禁用保存 — 见主维护源
- 不修改 ANI-main OpenAPI（由 Services 团队按 TASK-SVC-014 推进）

## 非目标

- 不自造 schema 为冻结事实
- 继承租户默认角色、权限审计导出、document 级 ACL
