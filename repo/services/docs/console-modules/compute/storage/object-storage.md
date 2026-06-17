# 对象存储

## 页面定位

`对象存储` 是 `Console / 算力与云资源 / 存储管理` 下的租户侧资源子页面，用于帮助用户查看、创建、删除和理解当前权限范围内的对象元数据资源。

本页属于 `Console` 页面，不是 `BOSS` 的平台对象存储池运营页。

## 文档管理规则

- 本文件是 `对象存储` 的主维护源
- `prototypes/ani-services-prototype-console.html` 只保留页面摘要、模块边界和详细材料入口
- `tasks/modules/prd/console/compute/prd-console-object-storage.md` 和 `tasks/modules/spec/console/compute/spec-console-object-storage.md` 作为阶段性产物保留，不替代本文件
- 如正文与辅助材料冲突，以本文件为准，并同步回修辅助材料

## Core 层要求

- `StorageObject` 资源对象属于 `Core`
- 查询与动作接口必须使用 `/api/v1/*`
- 不允许把对象资源对象写入 `Services /api/v1/svc/*`
- 不允许继续使用旧的 `/api/v1/storage/*` 或 `/api/v1/console/*` 路径
- 不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- 租户边界必须从认证上下文或后端中间件获取
- 创建接口必须使用现有 `Core v1.yaml` 已冻结的 `idempotency_key`
- 新增接口说明必须具备 `operationId`、中文 `summary`、`tags`、`security`、完整 `responses`
- 错误响应格式统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- 本模块必须严格区分“已冻结 Core 能力”和“待补 Core 契约能力”
- 本模块当前只承接对象元数据，不承接对象内容上传/下载链路

## 页面职责

- 展示当前租户下的对象元数据列表、详情、创建和删除入口
- 展示对象资源的 bucket、key、大小、内容类型和基础时间字段
- 为 VM、容器、GPU 容器等页面提供对象资源回指入口
- 明确标出尚未进入 `Core v1.yaml` 的对象上传、对象下载、桶资源、桶策略等能力
- 避免误用旧路径和旧描述

## 页面结构

```text
对象存储
├── 页面定位与边界说明
├── 冻结能力概览
├── 对象元数据列表区
├── 对象元数据详情区
├── 对象元数据创建区
├── 待补能力说明
│   ├── 对象上传
│   ├── 对象下载
│   ├── 桶资源
│   └── 桶策略
└── 删除风险提示
```

## 交互与使用规则

- 页面默认展示当前租户、当前权限范围内的对象元数据，不展示平台全量对象存储池
- 对已冻结能力，页面可以展示列表、详情、创建和删除入口
- 对待补能力，只能展示规划说明、禁用入口或待补提示，不能伪造 API
- 工作负载详情中的“对象”区域只提供跳转入口，不重写底层对象资源契约
- 删除动作必须明确提示风险影响
- 创建对象当前只表示对象元数据登记，不表示文件内容已经上传完成

## 数据来源与分层约束

### 数据来源划分

- `Core` 数据
  - 对象元数据列表、详情、创建、删除
- `Services` 数据
  - 本页不定义对象资源对象，不定义新的对象聚合资源契约

### 关键边界

- 本页只承接 `Core v1.yaml` 已冻结的对象元数据资源对象
- 本页不把 `对象上传`、`对象下载`、`桶资源`、`桶策略` 写成已对齐接口
- `bucket` 当前只作为 `StorageObject` 的已冻结字段，不代表已存在独立桶资源管理契约
- 如后续需要生成 `Core v1.yaml` 扩展草稿，应从本文件的待补能力清单继续下沉，而不是回到 HTML 旧稿重抄

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 冻结能力概览 | Core | 当前 `Core v1.yaml` 已冻结的对象元数据路径组 | 对象元数据列表区 |
| 对象元数据列表区 | Core | 展示当前租户下的对象元数据列表、状态、大小和类型 | 对象元数据详情 |
| 对象元数据详情区 | Core | 展示单个对象的基础字段和状态原因 | 对象元数据详情 |
| 对象元数据创建区 | Core | 展示元数据创建入口和请求约束 | 创建抽屉 |
| 待补能力说明 | 规划项 | 仅说明当前还未进入 `Core v1.yaml` 的对象能力 | 待补说明 |

## 模块区块详细说明

### 冻结能力概览

- 展示重点：`对象元数据列表 / 对象详情 / 创建 / 删除`
- 主要来源层：Core
- 展示口径：对齐 `Core v1.yaml` 现有 `ObjectStore` 路径组
- 异常/空态：无资源时展示空态，不伪造统计值
- 跳转目标：对象元数据列表区

### 对象元数据列表区

- 展示重点：bucket、key、大小、内容类型、状态、创建时间、更新时间
- 主要来源层：Core
- 展示口径：对齐 `StorageObject`
- 异常/空态：找不到资源时展示 `NOT_FOUND`
- 跳转目标：对象元数据详情

### 对象元数据详情区

- 展示重点：对象 ID、bucket、key、大小、内容类型、状态、状态原因、时间字段
- 主要来源层：Core
- 展示口径：对齐 `StorageObject`
- 异常/空态：对象不存在时展示 `NOT_FOUND`
- 跳转目标：无

### 对象元数据创建区

- 展示重点：`bucket`、`key`、`size_bytes`、`content_type`
- 主要来源层：Core
- 展示口径：对齐 `CreateStorageObjectRequest`
- 请求约束：前端内部生成并随请求携带 `idempotency_key`，但不作为用户可见展示字段
- 能力边界：当前只创建对象元数据，不表示文件上传已经完成
- 异常/空态：请求不合法时展示 `BAD_REQUEST`
- 跳转目标：创建结果刷新列表

### 待补能力说明

- 展示重点：`对象上传`、`对象下载`、`桶资源`、`桶策略`
- 主要来源层：规划项
- 展示口径：当前仅保留规划说明，不形成正式接口契约
- 异常/空态：若后端无能力，入口保持禁用或仅展示说明
- 跳转目标：后续 `Core v1.yaml` 扩充

## 字段级定义

| 字段 | 说明 | 展示建议 |
|---|---|---|
| id | 对象 ID | 文本；详情页保留复制能力 |
| tenant_id | 租户 ID | 仅视为后端返回字段，不作为前端必传 |
| bucket | 桶名 | 文本；当前只作为对象字段展示 |
| key | 对象 key | 文本；支持截断显示 |
| size_bytes | 对象大小 | 数值 + 单位 |
| content_type | 内容类型 | 文本 |
| state | 资源状态 | 标签展示 |
| reason | 状态原因 | 详情页展示 |
| dev_profile | 开发态辅助信息 | 默认不作为首屏展示字段；仅开发态按 Core 返回值直出 |
| created_at | 创建时间 | 时间 |
| updated_at | 更新时间 | 时间 |

## 字段展示规则

- 页面展示的对象资源字段必须直接映射现有 `Core v1.yaml` schema
- 可以在前端 view model 中转为 camelCase，但正文里的契约口径以 `Core v1.yaml` 为准
- 对 `tenant_id` 只视为后端返回字段，不视为前端必传字段
- `dev_profile` 属于可选开发态辅助字段，默认不作为首屏展示字段；如需展示，必须直接映射 Core 返回值
- `bucket` 当前只作为对象字段展示，不能扩写成已冻结桶资源列表或桶详情能力
- 页面不扩写对象上传、下载为当前已冻结动作

## 字段口径与单位

| 字段 | 口径 |
|---|---|
| `size_bytes` | 使用 bytes 为底层单位，前端可换算展示 |
| `content_type` | 使用当前 schema 字符串口径 |
| `bucket` | 使用当前 schema 字符串口径 |
| `key` | 使用当前 schema 字符串口径 |
| 时间字段 | 使用 ISO 8601 日期时间 |

## 状态与能力口径

### 已冻结的 Core 能力

- `GET /api/v1/objects`
- `POST /api/v1/objects`
- `GET /api/v1/objects/{object_id}`
- `DELETE /api/v1/objects/{object_id}`

### 待补 Core 契约能力

- 对象上传/下载/桶 — 详文 `object-storage-upload.md`（Phase 2 YAML 已声明 upload、buckets、download）
- 桶策略

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录与租户上下文 | 已认证 | `401 UNAUTHORIZED` / `403 FORBIDDEN` |
| `bucket`、`key` | 满足 schema 长度约束 | `400 BAD_REQUEST` |
| `idempotency_key` | 必填 | `400 BAD_REQUEST` |

### `POST /api/v1/objects` 补充

- 若提供 `size_bytes`，必须满足非负整数约束

## 操作可用性矩阵

| 资源 | list/get | create | delete | action | 前置校验 |
|---|---|---|---|---|---|
| 对象元数据 | 可用 | 可用 | 可用 | 无额外冻结动作 | 页面可提示用户关注潜在引用影响；当前冻结返回码以 `Core v1.yaml` 为准 |

## 删除前置校验与当前契约边界

- 页面可以提示用户先确认是否存在潜在引用影响，但当前权威源未冻结“引用冲突校验”作为正式后端契约
- 当前 `openapi/v1.yaml` 对 `DELETE /api/v1/objects/{object_id}` 显式冻结的错误响应为 `401`、`403`、`404`
- 因此本轮正文不把 `409 CONFLICT` 写成已冻结删除契约
- 若后续需要显式增加冲突错误码，必须先更新 `Core v1.yaml`，再同步回写本文件和辅助材料

## 接口冻结规则

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

## 成功响应示例

```json
{
  "id": "obj-12ab",
  "tenant_id": "t-demo",
  "bucket": "tenant-assets",
  "key": "datasets/train/file-001.jsonl",
  "size_bytes": 2048,
  "content_type": "application/jsonl",
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
  "message": "storage object not found",
  "request_id": "req-20260613-201"
}
```

## 回填前置依赖

- 后续若要把 `对象上传/下载` 写入正式接口，必须先扩充 `Core v1.yaml`
- 后续若要增加 `桶资源` 或 `桶策略`，必须先补充正式 schema 与路径
- 后续若要扩展更细粒度的对象引用关系，也必须先确认 `StorageObject` schema 是否需要增补

## 待确认项

- 对象上传/下载未来是否由对象网关独立承接
- 桶资源与桶策略是否进入下一轮 `Core v1.yaml` 扩充
- 删除依赖冲突是否进入下一轮 `Core v1.yaml` 扩充

## 回填验收标准

- 正文明确区分已冻结能力与待补能力
- 正文不再出现旧 `/api/v1/storage/*`、旧 `/api/v1/console/*` 路径
- 正文不要求前端显式传 `tenant_id / X-Tenant-Id`
- 正文中的路径、返回码、请求字段与现有 `Core v1.yaml` 一致
- 正文明确当前只承接对象元数据，不把上传/下载写成现有能力
- `PRD`、`SPEC`、HTML 摘要和本文件一致
- 本文件可以独立作为 `Console / 对象存储` 对齐 `Core` 的主维护材料

## 产品经理补充定义

### 目标用户与权限视角

- 租户管理员：查看、创建、删除对象元数据记录
- 业务成员：查看对象元数据详情，并在授权前提下执行元数据创建或删除
- 只读用户：仅查看列表与详情，不显示删除入口

### 默认视图与页面状态

- 默认展示对象元数据列表和创建入口，首屏必须显式提示“当前不承接上传/下载”
- 当无对象元数据时，页面展示空态，但不能把空态文案写成“暂无文件”
- 当详情或删除失败时，必须保留 `request_id` 与对象标识，便于定位问题
- `bucket` 仅作为对象字段展示，不得被包装成可管理的独立资源入口

### 核心任务流

1. 用户创建对象元数据记录并进入详情确认 `bucket / key / size / content_type` 等字段
2. 用户在列表中查找对象元数据并根据详情决定是否删除
3. 用户从存储总入口或其他模块跳入本页时，明确此页只承接元数据管理，不承接内容传输

### 跨模块协同

- 与 `存储管理` 协同，作为对象存储子模块下钻页
- 与使用对象引用的业务模块协同，仅通过对象 ID 或 key 跳转，不承诺上传下载链路
- 桶、桶策略、上传、下载等能力待冻结后再单独扩充

### 产品验收补充

- 产品文案必须持续强调“对象元数据”而不是“文件上传”
- 空态、成功态、删除确认文案都不能把 `bucket` 表述为独立资源中心
- 列表、详情、创建、删除四条链路都必须存在清晰回流点
- 本页不得出现任何上传、下载、桶策略的伪入口
