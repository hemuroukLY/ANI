# SPEC: Console sandbox-templates

> Technical specification — Core handler 契约见 `tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` §TASK-CORE-006  
> Source: `tasks/modules/prd/console/compute/prd-console-sandbox-templates.md`  
> Revised: 2026-06-17

## 1. Summary

Sandbox 创建可选**模板列表**与实例**安全事件**只读子页；模板 CRUD 不在当前冻结范围。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/sandbox-templates` | `listSandboxTemplates` | `200 + SandboxTemplateListResponse` | `scope:instances:read` |
| GET | `/api/v1/instances/{instance_id}/security-events` | `listInstanceSecurityEvents` | `200 + InstanceSecurityEventListResponse` | 实例读权限 |

### 2.3 Verified Schemas

- `SandboxTemplate`、`SandboxTemplateListResponse`
- `InstanceSecurityEvent`、`InstanceSecurityEventListResponse`

## 3. Page Scope

- 创建 Sandbox 向导中的模板选择器
- 实例详情 Tab：安全事件时间线
- **Non-Goals**：模板 CRUD；延长/暂停 lifecycle（见 `sandbox-instance-management.md`）

## 4. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户认证 | 已登录 | `401` |
| 实例读权限 | 安全事件查询 | `403` / `404` |

模板列表无写操作。**当前 YAML 未声明模板相关 `422`**。

## 5. 操作可用性矩阵

| 操作 | 只读 | 编辑/管理员 |
|---|---|---|
| 查看模板列表 | ✅ | ✅ |
| 查看安全事件 | ✅ | ✅ |
| 模板 CRUD | ❌ | ❌ |

## 6. 主维护源

- `docs/console-modules/compute/sandbox-templates.md`
- 父模块：`docs/console-modules/compute/sandbox-instance-management.md`

## 7. Handler 验收（Core 团队）

```bash
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/sandbox-templates"
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/instances/{id}/security-events"
```

OpenAPI 已声明 ≠ handler 已实现。
