# 知识库 — 权限管理

## 页面定位

配置知识库的 **readers / editors** 列表。

本页属于 **Services / KnowledgeBases** 写能力子页。

## 文档管理规则

- 本文是 KB 权限子页的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 与 `knowledge-base.md` CRUD 分层：本页仅负责 ACL
- TASK：`TASK-SVC-014`

## Services 层要求

| 方法 | 路径 | operationId |
|---|---|---|
| PUT | `/api/v1/svc/knowledge-bases/{kb_id}/permissions` | `updateKnowledgeBasePermissions` |

请求：`KnowledgeBasePermissionsUpdateRequest`（必填 `idempotency_key`、`readers[]`、`editors[]`）。

成功：返回更新后的 `KnowledgeBase`（以 YAML 为准）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 知识库编辑权限 | 当前用户为 editor 或租户管理员 | `403 FORBIDDEN` |
| 知识库存在 | 当前租户可见 | `404 NOT_FOUND` |
| 幂等键 | PUT 请求携带 `idempotency_key` | `400`（若 YAML 声明必填缺失） |

若 YAML 为 `updateKnowledgeBasePermissions` 声明 `422 PreconditionFailed`，可写「YAML 已举例语义：成员主体不存在」；否则写「**当前 YAML 未声明 `422`**」。

## 页面职责

- 权限矩阵 UI：readers / editors 主体选择（用户/组，以 schema 为准）
- 保存后展示更新后的权限列表
- 仅租户管理员或 KB 编辑者可修改

## 操作可用性矩阵

| 操作 | 只读 | 编辑者 | 管理员 |
|---|---|---|---|
| 查看权限 | 可用 | 可用 | 可用 |
| 更新权限 | 不可用 | 可用 | 可用 |

## 接口冻结规则

### `PUT /api/v1/svc/knowledge-bases/{kb_id}/permissions`

- 成功：`200 + KnowledgeBase`（或 YAML 声明的响应 schema）
- 错误：`401`、`403`、`404`、`400`（校验失败）
- 幂等：相同 `idempotency_key` 重复提交行为以 YAML 为准
- 租户边界：不可修改其他租户 KB

## 待补边界

- 继承租户默认角色 — **YAML 未声明**
- 权限变更审计导出 — 待补
- 细粒度 document 级 ACL — 非当前冻结能力

## 相关模块

- `knowledge-base.md`
- `tenant/role-permission-edit.md`（租户角色，非 KB ACL）

## 验收标准

- [ ] PUT 路径与 services/v1.yaml 一致
- [ ] 必填 `idempotency_key` 在表单层校验
- [ ] 只读用户 UI 禁用保存
- [ ] handler stub ≠ 已实现
