# Sandbox 模板

## 页面定位

查询 Sandbox 创建可选**模板列表**，服务于 Sandbox 实例创建向导。

## 文档管理规则

- 本文是 Sandbox 模板子页的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- Handler 契约：`tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` §TASK-CORE-006

## Core 层要求

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/sandbox-templates` | `listSandboxTemplates` |
| GET | `/api/v1/instances/{instance_id}/security-events` | `listInstanceSecurityEvents` |

响应：`SandboxTemplateListResponse`、`InstanceSecurityEventListResponse`。

RBAC：模板列表 `scope:instances:read`；安全事件同实例读权限。

## 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户认证 | 已登录 | `401` |
| 实例读权限 | 安全事件查询 | `403` / `404` |

模板列表无写操作前置条件。

## 操作可用性矩阵

| 操作 | 只读 | 编辑/管理员 |
|---|---|---|
| 查看模板列表 | ✅ | ✅ |
| 查看安全事件 | ✅（实例可见） | ✅ |
| 模板 CRUD | ❌ | ❌（YAML 未声明） |

## 页面职责

- 模板选择器（创建 Sandbox 时引用 `sandbox_config`）
- 实例详情 Tab：安全事件时间线

## 接口冻结规则

### `GET /api/v1/sandbox-templates`

- Query：`limit`、`cursor`
- 成功：`200 + SandboxTemplateListResponse`

### `GET /api/v1/instances/{instance_id}/security-events`

- 成功：`200 + InstanceSecurityEventListResponse`
- 实例不存在：`404`

## 待补边界

- Sandbox 模板 CRUD — YAML 未声明
- 延长/暂停专属 lifecycle action — 见 `sandbox-instance-management.md`

## 相关模块

- `sandbox-instance-management.md`
- TASK：`TASK-CORE-006`

## 验收标准

- [ ] 路径与 Phase 2 YAML 一致
- [ ] 不把模板 CRUD 写成冻结能力
