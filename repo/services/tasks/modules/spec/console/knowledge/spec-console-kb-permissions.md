# SPEC: Console kb-permissions

> Source: `tasks/modules/prd/console/knowledge/prd-console-kb-permissions.md`  
> Revised: 2026-06-17

## 1. Summary

配置知识库 **readers / editors** 列表的写能力子页。与 `knowledge-base.md` CRUD 分层：本页仅负责 ACL。属于 **Services / KnowledgeBases**。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | 权限 |
|---|---|---|---|---|
| PUT | `/api/v1/svc/knowledge-bases/{kb_id}/permissions` | `updateKnowledgeBasePermissions` | `200 + KnowledgeBase` | editor 或租户管理员 |

### 2.3 Verified Schemas

- `KnowledgeBasePermissionsUpdateRequest`（必填 `idempotency_key`、`readers[]`、`editors[]`）
- `KnowledgeBase`（更新后响应，以 YAML 为准）

## 3. Page Scope

- 权限矩阵 UI：readers / editors 主体选择（用户/组，以 schema 为准）
- 保存后展示更新后的权限列表
- 仅租户管理员或 KB 编辑者可修改

## 4. Non-Goals

- 继承租户默认角色 — **YAML 未声明**
- 权限变更审计导出 — 待补
- 细粒度 document 级 ACL — 非当前冻结能力
- 若 YAML 未为 `updateKnowledgeBasePermissions` 声明 `422`，则 **当前 YAML 未声明 `422`**

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 知识库编辑权限 | 当前用户为 editor 或租户管理员 | `403 FORBIDDEN` |
| 知识库存在 | 当前租户可见 | `404 NOT_FOUND` |
| 幂等键 | PUT 请求携带 `idempotency_key` | `400`（若 YAML 声明必填缺失） |

## 6. 操作可用性矩阵

| 操作 | 只读 | 编辑者 | 管理员 |
|---|---|---|---|
| 查看权限 | ✅ | ✅ | ✅ |
| 更新权限 | ❌ | ✅ | ✅ |

## 7. 主维护源

- `docs/console-modules/knowledge/kb-permissions.md`
- 相关：`knowledge-base.md`、`tenant/role-permission-edit.md`（租户角色，非 KB ACL）
- TASK：`TASK-SVC-014`

## 8. Handler 验收（Services 团队）

```bash
curl -X PUT -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"idempotency_key":"...","readers":[],"editors":[]}' \
  "$BASE/api/v1/svc/knowledge-bases/{kb_id}/permissions"
```

- 幂等：相同 `idempotency_key` 重复提交行为以 YAML 为准
- 只读用户 UI 禁用保存
- OpenAPI 已声明 ≠ handler 已实现
