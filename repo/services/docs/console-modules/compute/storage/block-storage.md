# 块存储

## 页面定位

`块存储` 是 `Console / 算力与云资源 / 存储管理` 下的租户侧资源子页面，用于帮助用户查看、创建、删除和理解当前权限范围内的块存储卷资源。

本页属于 `Console` 页面，不是 `BOSS` 的平台存储池运营页。

## 文档管理规则

- 本文件是 `块存储` 的主维护源
- `prototypes/ani-services-prototype-console.html` 只保留页面摘要、模块边界和详细材料入口
- `tasks/modules/prd/console/compute/prd-console-block-storage.md` 和 `tasks/modules/spec/console/compute/spec-console-block-storage.md` 作为阶段性产物保留，不替代本文件
- 如正文与辅助材料冲突，以本文件为准，并同步回修辅助材料

## Core 层要求

- `StorageVolume` 资源对象属于 `Core`
- 查询与动作接口必须使用 `/api/v1/*`
- 不允许把块存储资源对象写入 `Services /api/v1/svc/*`
- 不允许继续使用旧的 `/api/v1/storage/*` 或 `/api/v1/console/*` 路径
- 不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- 租户边界必须从认证上下文或后端中间件获取
- 创建接口必须使用现有 `Core v1.yaml` 已冻结的 `idempotency_key`
- 新增接口说明必须具备 `operationId`、中文 `summary`、`tags`、`security`、完整 `responses`
- 错误响应格式统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- 本模块必须严格区分“已冻结 Core 能力”和“待补 Core 契约能力”

## 页面职责

- 展示当前租户下的块存储卷列表、详情、创建和删除入口
- 展示卷资源的状态、容量、存储类和基础时间字段
- 为 VM、容器、GPU 容器等页面提供卷资源回指入口
- 明确标出尚未进入 `Core v1.yaml` 的卷快照、卷挂载、卷卸载等能力
- 避免误用旧路径和旧描述

## 页面结构

```text
块存储
├── 页面定位与边界说明
├── 冻结能力概览
├── 卷列表区
├── 卷详情区
├── 卷创建区
├── 待补能力说明
│   ├── 卷快照
│   ├── 卷挂载
│   └── 卷卸载
└── 删除风险提示
```

## 交互与使用规则

- 页面默认展示当前租户、当前权限范围内的块存储卷，不展示平台全量存储池
- 对已冻结能力，页面可以展示列表、详情、创建和删除入口
- 对待补能力，只能展示规划说明、禁用入口或待补提示，不能伪造 API
- 工作负载详情中的“卷”区域只提供跳转入口，不重写底层卷资源契约
- 删除动作必须明确提示风险影响

## 数据来源与分层约束

### 数据来源划分

- `Core` 数据
  - 块存储卷列表、详情、创建、删除
- `Services` 数据
  - 本页不定义卷资源对象，不定义新的卷聚合资源契约

### 关键边界

- 本页只承接 `Core v1.yaml` 已冻结的卷资源对象
- 本页不把 `快照`、`卷挂载`、`卷卸载` 写成已对齐接口
- 如后续需要生成 `Core v1.yaml` 扩展草稿，应从本文件的待补能力清单继续下沉，而不是回到 HTML 旧稿重抄

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 冻结能力概览 | Core | 当前 `Core v1.yaml` 已冻结的卷资源组 | 卷列表区 |
| 卷列表区 | Core | 展示当前租户下的卷列表、状态、容量和存储类 | 卷详情 |
| 卷详情区 | Core | 展示单个卷的基础字段和状态原因 | 卷详情 |
| 卷创建区 | Core | 展示创建卷入口和请求约束 | 创建抽屉 |
| 待补能力说明 | 规划项 | 仅说明当前还未进入 `Core v1.yaml` 的卷能力 | 待补说明 |

## 模块区块详细说明

### 冻结能力概览

- 展示重点：`卷列表 / 卷详情 / 创建 / 删除`
- 主要来源层：Core
- 展示口径：对齐 `Core v1.yaml` 现有 `Volumes` 路径组
- 异常/空态：无资源时展示空态，不伪造统计值
- 跳转目标：卷列表区

### 卷列表区

- 展示重点：名称、容量、存储类、状态、创建时间、更新时间
- 主要来源层：Core
- 展示口径：对齐 `StorageVolume`
- 异常/空态：找不到资源时展示 `NOT_FOUND`
- 跳转目标：卷详情

### 卷详情区

- 展示重点：卷 ID、名称、容量、存储类、状态、状态原因、时间字段
- 主要来源层：Core
- 展示口径：对齐 `StorageVolume`
- 异常/空态：卷不存在时展示 `NOT_FOUND`
- 跳转目标：无

### 卷创建区

- 展示重点：`name`、`size_gib`、`storage_class`
- 主要来源层：Core
- 展示口径：对齐 `CreateStorageVolumeRequest`
- 请求约束：前端内部生成并随请求携带 `idempotency_key`，但不作为用户可见展示字段
- 异常/空态：请求不合法时展示 `BAD_REQUEST`
- 跳转目标：创建结果刷新列表

### 待补能力说明

- 展示重点：`卷快照`、`卷挂载`、`卷卸载`
- 主要来源层：规划项
- 展示口径：当前仅保留规划说明，不形成正式接口契约
- 异常/空态：若后端无能力，入口保持禁用或仅展示说明
- 跳转目标：后续 `Core v1.yaml` 扩充

## 字段级定义

| 字段 | 说明 | 展示建议 |
|---|---|---|
| id | 卷 ID | 文本；详情页保留复制能力 |
| tenant_id | 租户 ID | 仅视为后端返回字段，不作为前端必传 |
| name | 卷名称 | 文本；列表可点击 |
| size_gib | 容量（GiB） | 数值 + 单位 |
| storage_class | 存储类 | 标签或文本 |
| state | 资源状态 | 标签展示 |
| reason | 状态原因 | 详情页展示 |
| dev_profile | 开发态辅助信息 | 默认不作为首屏展示字段；仅开发态按 Core 返回值直出 |
| created_at | 创建时间 | 时间 |
| updated_at | 更新时间 | 时间 |

## 字段展示规则

- 页面展示的卷资源字段必须直接映射现有 `Core v1.yaml` schema
- 可以在前端 view model 中转为 camelCase，但正文里的契约口径以 `Core v1.yaml` 为准
- 对 `tenant_id` 只视为后端返回字段，不视为前端必传字段
- `dev_profile` 属于可选开发态辅助字段，默认不作为首屏展示字段；如需展示，必须直接映射 Core 返回值
- 不扩写卷绑定关系为新的冻结接口

## 字段口径与单位

| 字段 | 口径 |
|---|---|
| `size_gib` | 使用 GiB 为单位 |
| `storage_class` | 使用当前 schema 字符串口径 |
| 时间字段 | 使用 ISO 8601 日期时间 |

## 状态与能力口径

### 已冻结的 Core 能力

- `GET /api/v1/volumes`
- `POST /api/v1/volumes`
- `GET /api/v1/volumes/{volume_id}`
- `DELETE /api/v1/volumes/{volume_id}`

### 待补 Core 契约能力

- 卷快照 — 详文 `block-storage-snapshot.md`（Phase 2 YAML 已声明）
- 卷挂载
- 卷卸载

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录与租户上下文 | 已认证 | `401 UNAUTHORIZED` / `403 FORBIDDEN` |
| `name` | 租户内唯一 | `409 CONFLICT` |
| `size_gib` | ≥ 1 | `400 BAD_REQUEST` |

### `POST /api/v1/volumes` 补充

- 请求体必须包含 `name`、`size_gib` 和 `idempotency_key`

## 操作可用性矩阵

| 资源 | list/get | create | delete | action | 前置校验 |
|---|---|---|---|---|---|
| 块存储卷 | 可用 | 可用 | 可用 | 无额外冻结动作 | 页面可以提示用户先确认是否存在潜在引用影响；当前冻结返回码以 `Core v1.yaml` 为准 |

## 删除前置校验与当前契约边界

- 页面可以提示用户先确认是否存在潜在引用影响，但当前权威源未冻结“引用冲突校验”作为正式后端契约
- 当前 `openapi/v1.yaml` 对 `DELETE /api/v1/volumes/{volume_id}` 显式冻结的错误响应为 `401`、`403`、`404`
- 因此本轮正文不把 `409 CONFLICT` 写成已冻结删除契约
- 若后续需要显式增加冲突错误码，必须先更新 `Core v1.yaml`，再同步回写本文件和辅助材料

## 接口冻结规则

### `GET /api/v1/volumes`

| 项 | 值 |
|---|---|
| operationId | `listStorageVolumes` |
| summary | `查询块存储卷列表` |
| tags | `["Volumes"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| query params | `limit?`、`cursor?` |
| success | `200 + StorageVolumeListResponse` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN` |

### `POST /api/v1/volumes`

| 项 | 值 |
|---|---|
| operationId | `createStorageVolume` |
| summary | `创建块存储卷` |
| tags | `["Volumes"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| requestBody.required | `name`、`size_gib`、`idempotency_key` |
| success | `201 + StorageVolume` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN` |

### `GET /api/v1/volumes/{volume_id}`

| 项 | 值 |
|---|---|
| operationId | `getStorageVolume` |
| summary | `查询块存储卷` |
| tags | `["Volumes"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `volume_id` |
| success | `200 + StorageVolume` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

### `DELETE /api/v1/volumes/{volume_id}`

| 项 | 值 |
|---|---|
| operationId | `deleteStorageVolume` |
| summary | `删除块存储卷` |
| tags | `["Volumes"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `volume_id` |
| success | `200 + StorageVolume` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

## 成功响应示例

```json
{
  "id": "vol-12ab",
  "tenant_id": "t-demo",
  "name": "volume-data-01",
  "size_gib": 200,
  "storage_class": "standard",
  "state": "available",
  "reason": null,
  "created_at": "2026-06-13T10:00:00Z",
  "updated_at": "2026-06-13T10:00:00Z"
}
```

## 错误返回示例

```json
{
  "code": "NOT_FOUND",
  "message": "storage volume not found",
  "request_id": "req-20260613-101"
}
```

## 回填前置依赖

- 后续若要把 `快照` 写入正式接口，必须先扩充 `Core v1.yaml`
- 后续若要增加 `挂载/卸载`，必须先补充正式 schema 与路径
- 后续若要扩展更细粒度的卷引用关系，也必须先确认 `StorageVolume` schema 是否需要增补

## 待确认项

- 卷快照未来是独立资源组，还是仍附着在卷动作之下
- 卷挂载/卸载未来是否由块存储侧独立承接
- 删除依赖冲突是否进入下一轮 `Core v1.yaml` 扩充

## 回填验收标准

- 正文明确区分已冻结能力与待补能力
- 正文不再出现旧 `/api/v1/storage/*`、旧 `/api/v1/console/*` 路径
- 正文不要求前端显式传 `tenant_id / X-Tenant-Id`
- 正文中的路径、返回码、请求字段与现有 `Core v1.yaml` 一致
- `PRD`、`SPEC`、HTML 摘要和本文件一致
- 本文件可以独立作为 `Console / 块存储` 对齐 `Core` 的主维护材料

## 产品经理补充定义

### 目标用户与权限视角

- 租户管理员：查看、创建、删除块存储卷
- 业务成员：查看卷状态和容量，并在授权前提下执行创建或删除
- 只读用户：仅查看列表与详情

### 默认视图与页面状态

- 默认展示卷列表、容量信息和创建入口
- 当租户无卷时，页面展示空态并提示当前只承接卷本体能力
- 删除冲突或引用中状态必须在删除前提示，不能等提交后才让用户猜测原因
- 当前未冻结的快照、挂载等能力只展示边界说明，不出现在主操作区

### 核心任务流

1. 用户创建卷后进入详情确认容量、存储类和状态
2. 用户从列表筛选或搜索目标卷，查看状态与更新时间
3. 用户在符合条件时删除卷，并通过结果反馈确认是否成功

### 跨模块协同

- 与 `VM` 等实例页仅通过卷摘要和跳转协同，不在本页承诺挂载动作
- 与存储总入口协同，作为块存储的下钻承接页
- 快照和挂载能力待冻结后再扩充到对应模块或子页

### 产品验收补充

- 页面必须突出“卷本体”能力边界，避免用户误以为已支持全套块存储操作
- 删除前提示必须说明是否存在引用或状态阻塞
- 列表、详情、创建、删除四条主流转都必须有明确反馈
- 本页不得把快照、挂载写成当前已冻结能力
