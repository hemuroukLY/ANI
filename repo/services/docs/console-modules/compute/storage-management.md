# 存储管理

## 页面定位

`存储管理` 是 `Console / 算力与云资源` 下的租户侧资源管理页面，用于帮助用户查看、创建、删除和理解当前权限范围内的存储资源。

本页属于 `Console` 页面，不是 `BOSS` 的平台存储池运营页。

## 文档管理规则

- 本文件是 `存储管理` 的主维护源
- `prototypes/ani-services-prototype-console.html` 只保留页面摘要、模块边界和详细材料入口
- `tasks/modules/prd/console/compute/prd-console-storage-management.md` 和 `tasks/modules/spec/console/compute/spec-console-storage-management.md` 作为阶段性产物保留，不替代本文件
- 如正文与辅助材料冲突，以本文件为准，并同步回修辅助材料

## Core 层要求

- `StorageVolume`、`StorageFilesystem`、`StorageObject`、`VectorStore` 资源对象属于 `Core`
- 查询与动作接口必须使用 `/api/v1/*`
- 不允许把存储资源对象写入 `Services /api/v1/svc/*`
- 不允许继续使用旧的 `/api/v1/storage/*` 或 `/api/v1/console/*` 路径
- 不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- 租户边界必须从认证上下文或后端中间件获取
- 创建接口必须使用现有 `Core v1.yaml` 已冻结的 `idempotency_key`
- 新增接口说明必须具备 `operationId`、中文 `summary`、`tags`、`security`、完整 `responses`
- 错误响应格式统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- 本模块必须严格区分“已冻结 Core 能力”和“待补 Core 契约能力”

## 页面职责

- 展示存储资源对象关系和租户侧使用边界
- 展示当前租户下的块存储、文件存储、对象元数据、向量存储
- 提供这些资源的列表、详情、创建和删除入口
- 提供向量存储搜索入口
- 为 VM、容器、GPU 容器等页面提供存储关联回指入口
- 明确标出尚未进入 `Core v1.yaml` 的存储能力，避免误用旧路径和旧描述

## 页面结构

```text
存储管理
├── 页面定位与边界说明
├── 冻结能力概览
│   ├── 块存储
│   ├── 文件存储
│   ├── 对象存储
│   └── 向量存储
├── 待补能力说明
│   ├── 快照
│   ├── 卷挂载/卸载
│   ├── 对象上传/下载
│   ├── 桶策略
│   └── 文件系统挂载目标
├── 资源关系说明
└── 详细资源区
    ├── 块存储
    ├── 文件存储
    ├── 对象存储
    └── 向量存储
```

## 交互与使用规则

- 页面默认展示当前租户、当前权限范围内的存储资源，不展示平台全量存储池
- 对已冻结能力，页面可以展示列表、详情、创建和删除入口
- 对待补能力，只能展示规划说明、禁用入口或待补提示，不能伪造 API
- 实例详情中的“卷 / 文件 / 对象 / 向量库”区域只展示当前关联关系和跳转入口，不重写底层存储资源契约
- 删除动作必须明确提示引用关系和风险影响

## 数据来源与分层约束

### 数据来源划分

- `Core` 数据
  - 块存储卷列表、详情、创建、删除
  - 文件存储列表、详情、创建、删除
  - 对象元数据列表、详情、创建、删除
  - 向量存储列表、详情、创建、删除、搜索
- `Services` 数据
  - 本页不定义存储资源对象，不定义新的存储聚合资源契约

### 关键边界

- 本页只承接 `Core v1.yaml` 已冻结的存储资源对象
- 本页不把 `快照`、`卷挂载/卸载`、`对象上传/下载`、`桶策略`、`文件系统挂载目标` 写成已对齐接口
- 对象存储当前只承接对象元数据资源，不代表内容传输链路已经冻结
- 如后续需要生成 `Core v1.yaml` 扩展草稿，应从本文件的待补能力清单继续下沉，而不是回到 HTML 旧稿重抄

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 冻结能力概览 | Core | 当前 `Core v1.yaml` 已冻结的存储资源组 | 各资源区 |
| 块存储 | Core | 展示当前租户下的块存储卷列表、详情、创建、删除 | 块存储详情 |
| 文件存储 | Core | 展示当前租户下的文件存储列表、详情、创建、删除 | 文件存储详情 |
| 对象存储 | Core | 展示当前租户下的对象元数据列表、详情、创建、删除 | 对象元数据详情 |
| 向量存储 | Core | 展示当前租户下的向量存储列表、详情、创建、删除、搜索 | 向量存储详情 |
| 待补能力说明 | 规划项 | 仅说明当前还未进入 `Core v1.yaml` 的存储能力 | 待补说明 |

## 模块区块详细说明

### 冻结能力概览

- 展示重点：`块存储 / 文件存储 / 对象存储 / 向量存储`
- 主要来源层：Core
- 展示口径：对齐 `Core v1.yaml` 现有 `Volumes / Filesystems / ObjectStore / VectorStore` 路径组
- 异常/空态：无资源时展示空态，不伪造统计值
- 跳转目标：对应资源区

### 块存储

- 展示重点：名称、容量、存储类、状态、创建时间、更新时间
- 主要来源层：Core
- 展示口径：对齐 `StorageVolume`
- 异常/空态：找不到资源时展示 `NOT_FOUND`
- 跳转目标：块存储详情

### 文件存储

- 展示重点：名称、协议、容量、endpoint、状态、创建时间
- 主要来源层：Core
- 展示口径：对齐 `StorageFilesystem`
- 异常/空态：找不到资源时展示 `NOT_FOUND`
- 跳转目标：文件存储详情

### 对象存储

- 展示重点：bucket、key、大小、内容类型、状态、创建时间
- 主要来源层：Core
- 展示口径：对齐 `StorageObject`
- 异常/空态：找不到资源时展示 `NOT_FOUND`
- 跳转目标：对象元数据详情

### 向量存储

- 展示重点：名称、维度、距离度量、状态、创建时间
- 主要来源层：Core
- 展示口径：对齐 `VectorStore`
- 异常/空态：找不到资源时展示 `NOT_FOUND`
- 跳转目标：向量存储详情

### 待补能力说明

- 展示重点：`快照`、`卷挂载/卸载`、`对象上传/下载`、`桶策略`、`文件系统挂载目标`
- 主要来源层：规划项
- 展示口径：当前仅保留规划说明，不形成正式接口契约
- 异常/空态：若后端无能力，入口保持禁用或仅展示说明
- 跳转目标：后续 `Core v1.yaml` 扩充

## 字段级定义

### 块存储字段

| 字段 | 说明 | 展示建议 |
|---|---|---|
| id | 卷 ID | 文本；详情页保留复制能力 |
| name | 卷名称 | 文本；列表可点击 |
| size_gib | 容量（GiB） | 数值 + 单位 |
| storage_class | 存储类 | 标签或文本 |
| state | 资源状态 | 标签展示 |
| reason | 状态原因 | 详情页展示 |
| created_at | 创建时间 | 时间 |
| updated_at | 更新时间 | 时间 |

### 文件存储字段

| 字段 | 说明 | 展示建议 |
|---|---|---|
| id | 文件系统 ID | 文本 |
| name | 文件系统名称 | 文本 |
| protocol | 协议 | `nfs / cephfs` 标签 |
| size_gib | 容量（GiB） | 数值 + 单位 |
| endpoint | 访问端点 | 文本；为空时展示 `-` |
| state | 资源状态 | 标签展示 |
| reason | 状态原因 | 详情页展示 |
| created_at | 创建时间 | 时间 |
| updated_at | 更新时间 | 时间 |

### 对象元数据字段

| 字段 | 说明 | 展示建议 |
|---|---|---|
| id | 对象 ID | 文本 |
| bucket | 桶名 | 文本 |
| key | 对象 key | 文本；支持截断显示 |
| size_bytes | 对象大小 | 数值 + 单位 |
| content_type | 内容类型 | 文本 |
| state | 资源状态 | 标签展示 |
| reason | 状态原因 | 详情页展示 |
| created_at | 创建时间 | 时间 |
| updated_at | 更新时间 | 时间 |

### 向量存储字段

| 字段 | 说明 | 展示建议 |
|---|---|---|
| id | 向量库 ID | 文本 |
| name | 向量库名称 | 文本 |
| dimension | 向量维度 | 数值 |
| metric | 距离度量 | `cosine / l2 / ip` 标签 |
| state | 资源状态 | 标签展示 |
| reason | 状态原因 | 详情页展示 |
| created_at | 创建时间 | 时间 |
| updated_at | 更新时间 | 时间 |

## 字段展示规则

- 页面展示的存储资源字段必须直接映射现有 `Core v1.yaml` schema
- 可以在前端 view model 中转为 camelCase，但正文里的契约口径以 `Core v1.yaml` 为准
- 对 `tenant_id` 只视为后端返回字段，不视为前端必传字段
- 对象存储首版只展示元数据，不展示实际文件上传/下载链路
- 向量搜索首版只展示请求向量、筛选条件和命中结果摘要，不扩写写入/索引维护动作

## 字段口径与单位

| 字段 | 口径 |
|---|---|
| `size_gib` | 使用 GiB 为单位 |
| `size_bytes` | 使用 bytes 为底层单位，前端可换算展示 |
| `metric` | 仅允许 `cosine`、`l2`、`ip` |
| `protocol` | 仅允许 `nfs`、`cephfs` |
| 时间字段 | 使用 ISO 8601 日期时间 |

## 状态与能力口径

### 已冻结的 Core 能力

- `GET /api/v1/volumes`
- `POST /api/v1/volumes`
- `GET /api/v1/volumes/{volume_id}`
- `DELETE /api/v1/volumes/{volume_id}`
- `GET /api/v1/filesystems`
- `POST /api/v1/filesystems`
- `GET /api/v1/filesystems/{filesystem_id}`
- `DELETE /api/v1/filesystems/{filesystem_id}`
- `GET /api/v1/objects`
- `POST /api/v1/objects`
- `GET /api/v1/objects/{object_id}`
- `DELETE /api/v1/objects/{object_id}`
- `GET /api/v1/vector-stores`
- `POST /api/v1/vector-stores`
- `GET /api/v1/vector-stores/{vector_store_id}`
- `DELETE /api/v1/vector-stores/{vector_store_id}`
- `POST /api/v1/vector-stores/{vector_store_id}/search`

### 待补 Core 契约能力

- 卷快照
- 卷挂载/卸载
- 对象上传/下载
- 桶策略
- 文件系统挂载目标

## 创建前置条件

### `POST /api/v1/volumes`

- 调用方已通过标准鉴权
- 请求体必须包含 `name`、`size_gib` 和 `idempotency_key`
- `size_gib` 必须满足最小值约束

### `POST /api/v1/filesystems`

- 调用方已通过标准鉴权
- 请求体必须包含 `name`、`size_gib` 和 `idempotency_key`
- 若提供 `protocol`，必须满足 `nfs / cephfs` 枚举约束

### `POST /api/v1/objects`

- 调用方已通过标准鉴权
- 请求体必须包含 `bucket`、`key` 和 `idempotency_key`
- `bucket`、`key` 长度必须满足 schema 约束

### `POST /api/v1/vector-stores`

- 调用方已通过标准鉴权
- 请求体必须包含 `name`、`dimension` 和 `idempotency_key`
- `dimension` 必须为大于 0 的整数
- 若提供 `metric`，必须满足 `cosine / l2 / ip` 枚举约束

### `POST /api/v1/vector-stores/{vector_store_id}/search`

- 调用方已通过标准鉴权
- 路径参数 `vector_store_id` 必须存在且当前租户可见
- 请求体必须包含 `vector`
- `vector` 维度应与目标向量库配置一致

> 本页为存储总入口，各子模块详文见 `storage/` 下对应文件。

## 产品经理补充定义

### 目标用户与权限视角

- 租户管理员：查看并管理当前租户可见的块、文件、对象元数据和向量存储资源
- 业务成员：查看存储资源与关联关系，并在授权前提下进入子模块执行操作
- 只读用户：仅查看资源摘要、详情与边界说明

### 默认视图与页面状态

- 首屏默认展示 `冻结能力概览 + 资源关系说明 + 四类资源区`，作为统一存储总入口
- 当租户尚无存储资源时，页面展示按资源类型创建的入口，不伪造总览统计值
- 当某个子资源区失败时，只影响对应资源区，其他区块继续可用
- 对 `快照 / 挂载 / 上传下载 / 桶策略 / 挂载目标` 等待补能力，入口只展示说明和禁用态

### 核心任务流

1. 用户根据使用场景先选择资源类型，再进入对应子模块完成创建、查看或删除
2. 用户从实例、知识库或首页回跳到本页，确认当前需要下钻到哪类存储资源
3. 用户在删除前查看引用关系或待补边界，避免把尚未冻结的能力当成既有操作流

### 跨模块协同

- 与 `VM / 容器实例 / GPU 容器实例` 协同，只展示关联资源摘要和跳转入口
- 与 `知识库管理` 协同，仅通过向量存储摘要和跳转衔接，不重新定义业务资源关系
- 与首页协同，用于承接存储异常或容量摘要的钻取入口

### 产品验收补充

- 用户必须能明确区分四类存储资源，而不是把本页误解为统一 CRUD 页面
- 每个子资源都必须有清晰的下钻入口和回流关系
- 待补能力必须以禁用态或说明态呈现，不能出现“可点击但无契约”的假入口
- 对象存储必须持续强调“当前仅承接对象元数据”

## 操作可用性矩阵

| 资源 | list/get | create | delete | action | 前置校验 |
|---|---|---|---|---|---|
| 块存储卷 | 可用 | 可用 | 可用 | 无额外冻结动作 | 页面可以提示用户先确认是否存在潜在引用影响；当前冻结返回码以 `Core v1.yaml` 为准 |
| 文件存储 | 可用 | 可用 | 可用 | 无额外冻结动作 | 页面可以提示用户先确认是否存在活跃挂载影响；当前冻结返回码以 `Core v1.yaml` 为准 |
| 对象元数据 | 可用 | 可用 | 可用 | 无额外冻结动作 | 页面可以提示用户先确认是否存在潜在引用影响；当前冻结返回码以 `Core v1.yaml` 为准 |
| 向量存储 | 可用 | 可用 | 可用 | `search` 可用 | 搜索前目标向量库必须存在且可访问 |

## 删除前置校验与当前契约边界

- 页面可以提示用户先确认是否存在潜在依赖影响，例如卷引用、文件系统活跃挂载、对象关键引用、向量库活动依赖，但当前权威源未冻结这些冲突校验作为统一正式后端契约
- 当前 `openapi/v1.yaml` 对存储删除接口显式冻结的错误响应为 `401`、`403`、`404`
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

### `GET /api/v1/filesystems`

| 项 | 值 |
|---|---|
| operationId | `listStorageFilesystems` |
| summary | `查询文件存储列表` |
| tags | `["Filesystems"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| query params | `limit?`、`cursor?` |
| success | `200 + StorageFilesystemListResponse` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN` |

### `POST /api/v1/filesystems`

| 项 | 值 |
|---|---|
| operationId | `createStorageFilesystem` |
| summary | `创建文件存储` |
| tags | `["Filesystems"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| requestBody.required | `name`、`size_gib`、`idempotency_key` |
| success | `201 + StorageFilesystem` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN` |

### `GET /api/v1/filesystems/{filesystem_id}`

| 项 | 值 |
|---|---|
| operationId | `getStorageFilesystem` |
| summary | `查询文件存储` |
| tags | `["Filesystems"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `filesystem_id` |
| success | `200 + StorageFilesystem` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

### `DELETE /api/v1/filesystems/{filesystem_id}`

| 项 | 值 |
|---|---|
| operationId | `deleteStorageFilesystem` |
| summary | `删除文件存储` |
| tags | `["Filesystems"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `filesystem_id` |
| success | `200 + StorageFilesystem` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

### `GET /api/v1/objects`

| 项 | 值 |
|---|---|
| operationId | `listStorageObjects` |
| summary | `查询对象元数据列表` |
| tags | `["ObjectStore"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| query params | `limit?`、`cursor?` |
| success | `200 + StorageObjectListResponse` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN` |

### `POST /api/v1/objects`

| 项 | 值 |
|---|---|
| operationId | `createStorageObject` |
| summary | `创建对象元数据` |
| tags | `["ObjectStore"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| requestBody.required | `bucket`、`key`、`idempotency_key` |
| success | `201 + StorageObject` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN` |

### `GET /api/v1/objects/{object_id}`

| 项 | 值 |
|---|---|
| operationId | `getStorageObject` |
| summary | `查询对象元数据` |
| tags | `["ObjectStore"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `object_id` |
| success | `200 + StorageObject` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

### `DELETE /api/v1/objects/{object_id}`

| 项 | 值 |
|---|---|
| operationId | `deleteStorageObject` |
| summary | `删除对象元数据` |
| tags | `["ObjectStore"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `object_id` |
| success | `200 + StorageObject` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

### `GET /api/v1/vector-stores`

| 项 | 值 |
|---|---|
| operationId | `listVectorStores` |
| summary | `查询向量存储列表` |
| tags | `["VectorStore"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| query params | `limit?`、`cursor?` |
| success | `200 + VectorStoreListResponse` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN` |

### `POST /api/v1/vector-stores`

| 项 | 值 |
|---|---|
| operationId | `createVectorStore` |
| summary | `创建向量存储` |
| tags | `["VectorStore"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| requestBody.required | `name`、`dimension`、`idempotency_key` |
| success | `201 + VectorStore` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN` |

### `GET /api/v1/vector-stores/{vector_store_id}`

| 项 | 值 |
|---|---|
| operationId | `getVectorStore` |
| summary | `查询向量存储` |
| tags | `["VectorStore"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `vector_store_id` |
| success | `200 + VectorStore` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

### `DELETE /api/v1/vector-stores/{vector_store_id}`

| 项 | 值 |
|---|---|
| operationId | `deleteVectorStore` |
| summary | `删除向量存储` |
| tags | `["VectorStore"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `vector_store_id` |
| success | `200 + VectorStore` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

### `POST /api/v1/vector-stores/{vector_store_id}/search`

| 项 | 值 |
|---|---|
| operationId | `searchVectorStore` |
| summary | `搜索向量存储` |
| tags | `["VectorStore"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `vector_store_id` |
| requestBody.required | `vector` |
| success | `200 + VectorStoreSearchResponse` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

## 成功响应示例

### 块存储创建成功

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

### 向量搜索成功

```json
{
  "items": [
    {
      "id": "doc-001",
      "score": 0.98,
      "metadata": {
        "kb_id": "kb-001",
        "chunk_id": "chunk-01"
      }
    }
  ],
  "total": 1
}
```

## 错误返回示例

```json
{
  "code": "NOT_FOUND",
  "message": "storage volume not found",
  "request_id": "req-20260613-001"
}
```

## 回填前置依赖

- 后续若要把 `快照` 写入正式接口，必须先扩充 `Core v1.yaml`
- 后续若要增加 `挂载/卸载`、`对象上传/下载`、`桶策略`、`文件系统挂载目标`，必须先补充正式 schema 与路径
- 后续若要扩展向量存储写入、索引刷新等动作，也必须先补充 `Core v1.yaml`

## 待确认项

- 卷快照未来是独立资源组，还是仍附着在实例或卷动作之下
- 文件系统挂载目标是否进入下一轮 `Core v1.yaml` 扩充
- 对象上传/下载未来是否由独立网关或对象服务承担

## 回填验收标准

- 正文明确区分已冻结能力与待补能力
- 正文不再出现旧 `/api/v1/storage/*`、旧 `/api/v1/console/*` 路径
- 正文不要求前端显式传 `tenant_id / X-Tenant-Id`
- 正文中的路径、返回码、请求字段与现有 `Core v1.yaml` 一致
- `PRD`、`SPEC`、HTML 摘要和本文件一致
- 本文件可以独立作为 `Console / 存储管理` 对齐 `Core` 的主维护材料
