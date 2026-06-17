# 裸金属 / DPU 实例

## 页面定位

展示与管理 **`kind=bare_metal` / `kind=dpu_node`** 统一实例的 Console 页；复用 Core instances 契约，非独立裸金属 CRUD 域。

## 文档管理规则

- 本文是 **裸金属 / DPU 实例** 的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 不得把规划路径或 handler stub 写成已实现
- PRD/SPEC 同步规则见 `docs/console-modules/governance/module-delivery-workflow.md` §3.5

## Core 层要求

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/instances?kind=bare_metal` | `listInstances` |
| GET | `/api/v1/instances?kind=dpu_node` | `listInstances` |
| GET | `/api/v1/instances/{instance_id}` | `getInstance` |
| POST | `/api/v1/instances` | `createInstance`（`kind` 枚举含 `bare_metal`、`dpu_node`） |
| POST | `/api/v1/instances/{instance_id}/lifecycle` | `applyInstanceLifecycle` |

Schema：`InstanceRecord.kind` / `instance_type` 枚举已含 `bare_metal`、`dpu_node`。

RBAC：`scope:instances:*`（与同 unified instances）。

<!-- TODO-YAML: 独立 dpu-inventory 路径在 v1.yaml 顶部说明中规划，当前 paths 段未声明 -->

## Services 层要求

无 Services 路径。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 模块读/写权限 | 对应 RBAC scope | `403 FORBIDDEN` |

写操作 POST/PUT 的 `idempotency_key` 与 422 口径 — **待 YAML 冻结后按 ../governance/module-delivery-workflow.md §2.10 补充**。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员/运维 | 说明 |
|---|---|---|---|
| 列表/详情 | 可用 | 可用 | kind 过滤 |
| 创建 bare_metal/dpu | 不可用 | 可用 | handler 待实现 |
| 生命周期 | 不可用 | 可用 | lifecycle |

## 页面职责

- 占位 UI + 明确 YAML/OpenAPI 缺口（若适用）
- 跳转关联模块（见「相关模块」）
- 不把 BOSS/平台运营能力写入 Console 冻结契约

## 接口冻结规则

### `GET /api/v1/instances?kind=bare_metal`

- operationId: `listInstances`
- 成功：`200 + InstanceListResponse`
- RBAC：`scope:instances:read`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文

### `GET /api/v1/instances?kind=dpu_node`

- operationId: `listInstances`
- 成功：`200 + InstanceListResponse`
- RBAC：`scope:instances:read`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文

### `POST /api/v1/instances`

- operationId: `createInstance`
- 成功：`201 + CreateInstanceResponse`
- RBAC：`scope:instances:create`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文

<!-- 上表为当前 YAML 已声明子集；模块完整能力仍待 TODO-YAML -->

## 待补边界

- DPU 设备清单 `dpu-inventory` — YAML 顶部规划，paths **未声明**
- 裸金属专有 provisioning 字段 — 待 Core 扩 schema

## 相关模块

- `compute/vm-management.md`
- `compute/gpu-inventory-ui.md`

## 验收标准

- [ ] 路径与 OpenAPI 权威源一致（或明确 TODO-YAML / N/A）
- [ ] 正文不把 handler stub 写成已实现
- [ ] 含创建前置条件与操作可用性矩阵
- [ ] PRD/SPEC/HTML 摘要已同步
