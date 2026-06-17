# 向量存储

## 页面定位

`向量存储` 是 `Console / 算力与云资源 / 存储管理` 下的租户侧资源子页面，用于帮助用户查看、创建、删除和搜索当前权限范围内的向量存储资源。

本页属于 `Console` 页面，不是 `BOSS` 的平台向量基础设施运营页。

## 文档管理规则

- 本文件是 `向量存储` 的主维护源
- `prototypes/ani-services-prototype-console.html` 只保留摘要、边界和详细材料入口
- `tasks/modules/prd/console/compute/prd-console-vector-storage.md` 与 `tasks/modules/spec/console/compute/spec-console-vector-storage.md` 作为辅助材料保留
- 如正文与辅助材料冲突，以本文件为准

## Core 层要求

- `VectorStore` 资源对象属于 `Core`
- 查询与动作接口必须使用 `/api/v1/vector-stores*`
- 不允许把向量存储资源对象写入 `Services /api/v1/svc/*`
- 不允许继续使用旧 `/api/v1/storage/*` 或 `/api/v1/console/*` 路径
- 不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- 创建接口必须使用现有 `Core v1.yaml` 已冻结的 `idempotency_key`
- 搜索接口当前只承接 `vector`、`top_k`、`filter` 口径

## 冻结事实表

### Frozen Paths

- `GET /api/v1/vector-stores`
- `POST /api/v1/vector-stores`
- `GET /api/v1/vector-stores/{vector_store_id}`
- `DELETE /api/v1/vector-stores/{vector_store_id}`
- `POST /api/v1/vector-stores/{vector_store_id}/search`

### Frozen Schemas

- `VectorStore`
- `VectorStoreListResponse`
- `CreateVectorStoreRequest`
- `VectorStoreSearchRequest`
- `VectorStoreSearchResponse`

### Non-Frozen Capabilities（本模块范围）

- 向量写入（**YAML 已声明**：`POST /vector-stores/{id}/documents`；详文在 `vector-store-write.md`，本模块不承接）
- 索引维护
- 批量导入

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 冻结能力概览 | Core | 当前 `VectorStore` 路径组与 `search` | 列表区 |
| 向量存储列表区 | Core | 展示名称、维度、metric、状态 | 详情区 |
| 向量存储创建区 | Core | 对齐 `CreateVectorStoreRequest` | 创建抽屉 |
| 向量搜索区 | Core | 对齐 `VectorStoreSearchRequest` | 结果摘要 |
| 待补能力说明 | 规划项 | 仅说明当前未冻结能力 | 后续模块文档 |

## 字段级定义

| 字段 | 说明 | 展示建议 |
|---|---|---|
| id | 向量库 ID | 文本 |
| tenant_id | 租户 ID | 仅视为后端返回字段，不作为前端必传 |
| name | 向量库名称 | 文本 |
| dimension | 向量维度 | 数值 |
| metric | 距离度量 | `cosine / l2 / ip` 标签 |
| state | 资源状态 | 标签展示 |
| reason | 状态原因 | 详情页展示 |
| created_at | 创建时间 | 时间 |
| updated_at | 更新时间 | 时间 |

## 搜索请求重点

| 字段 | 说明 | 口径 |
|---|---|---|
| vector | 检索向量 | 必填，`minItems=1` |
| top_k | 返回条数 | `1..100`，默认 `10` |
| filter | 额外筛选条件 | `additionalProperties<string>` |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录与租户上下文 | 已认证 | `401` / `403` |
| `name` | 租户内唯一 | `409 CONFLICT` |
| `dimension` | 正整数 | `400 BAD_REQUEST` |

## 操作可用性矩阵

基于 `VectorStoreState`（`openapi/v1.yaml` 已冻结）：`pending` / `ready` / `failed` / `deleting` / `deleted`。

| 操作 | pending | ready | failed | deleting | deleted |
|---|---|---|---|---|---|
| 查看详情 | ✅ | ✅ | ✅ | ⚠️ 只读 | ❌ |
| 向量搜索 | ❌ | ✅ | ❌ | ❌ | ❌ |
| 删除 | ❌ | ✅ | ✅ | ❌ | ❌ |

非 `ready` 状态执行搜索时，页面应置灰搜索入口；`searchVectorStore` **YAML 已声明 `422 PreconditionFailed`**（Core v1.yaml Phase 2 2026-06-17）；具体 `code` 值待 Core description 补充后再写成冻结契约。<!-- ADDED-TO-YAML: POST /api/v1/vector-stores/{vector_store_id}/search 增 422 PreconditionFailed (Core v1.yaml, Phase 2 2026-06-17) -->

## 接口冻结规则

### `GET /api/v1/vector-stores`

- `operationId`: `listVectorStores`
- `success`: `200 + VectorStoreListResponse`
- `error responses`: `401 UNAUTHORIZED`、`403 FORBIDDEN`

### `POST /api/v1/vector-stores`

- `operationId`: `createVectorStore`
- `requestBody.required`: `name`、`dimension`、`idempotency_key`
- `success`: `201 + VectorStore`
- `error responses`: `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN`

### `GET /api/v1/vector-stores/{vector_store_id}`

- `operationId`: `getVectorStore`
- `success`: `200 + VectorStore`
- `error responses`: `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND`

### `DELETE /api/v1/vector-stores/{vector_store_id}`

- `operationId`: `deleteVectorStore`
- `success`: `200 + VectorStore`
- `error responses`: `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND`

### `POST /api/v1/vector-stores/{vector_store_id}/search`

- `operationId`: `searchVectorStore`
- `requestBody.required`: `vector`
- `success`: `200 + VectorStoreSearchResponse`
- `error responses`: `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND`、`422 PreconditionFailed`（YAML 已声明；具体 code 待 Core description 补充）

## 删除前置校验与当前契约边界

- 页面可以提示用户先确认是否存在潜在依赖影响，但当前权威源未冻结“依赖冲突校验”作为正式后端契约
- 当前 `openapi/v1.yaml` 未显式冻结 `409 CONFLICT`
- 因此本轮正文不把 `409` 写成已冻结删除契约

## 回填前置依赖

- 向量写入 — 详文 `vector-store-write.md`（Phase 2 YAML 已声明 `POST .../documents`）
- 后续若要补 `索引维护 / 批量导入`，必须先补正式 schema 与路径

## 回填验收标准

- 正文明确区分已冻结能力与待补能力
- 正文不再出现旧 `/api/v1/storage/*`、旧 `/api/v1/console/*` 路径
- 正文不要求前端显式传 `tenant_id / X-Tenant-Id`
- 正文不把向量写入、索引维护、批量导入写成当前正式接口
- `PRD`、`SPEC`、HTML 摘要与本文件一致

## 产品经理补充定义

### 目标用户与权限视角

- 租户管理员：查看、创建、搜索、删除向量存储
- 业务成员：查看向量库详情并在授权前提下执行搜索
- 只读用户：仅查看列表、详情和搜索结果

### 默认视图与页面状态

- 默认展示向量库列表，并保留进入详情和搜索的入口
- 当无向量库时，页面展示空态并提示当前不承接写入、导入和索引维护
- 搜索失败时只影响当前搜索区，不阻断列表和详情浏览
- 删除前必须明确提示删除后搜索不可用，但不能假定不存在的索引任务回滚能力

### 核心任务流

1. 用户创建向量库后进入详情确认维度、距离度量与状态
2. 用户在详情或搜索区提交搜索请求，查看命中结果摘要
3. 用户根据使用情况决定删除向量库，并通过结果反馈确认完成

### 跨模块协同

- 与 `知识库管理` 仅通过向量资源摘要和跳转协同，不把知识库流程塞回本页
- 与存储总入口协同，作为向量存储的下钻页
- 写入、批量导入、索引维护等能力待冻结后单独扩展

### 产品验收补充

- 页面必须把“搜索”定义为当前主动作之一，但不得扩写为完整数据写入平台
- 搜索结果必须可回看当前请求上下文，避免用户误判命中来源
- 空态和错误态必须强调当前边界，不能暗示存在导入或索引修复功能
- 本页不得把向量写入、索引维护写成已冻结能力
