# 文件存储

## 页面定位

`文件存储` 是 `Console / 算力与云资源 / 存储管理` 下的租户侧资源子页面，用于帮助用户查看、创建、删除和理解当前权限范围内的文件存储资源。

本页属于 `Console` 页面，不是 `BOSS` 的平台存储池运营页。

## 文档管理规则

- 本文件是 `文件存储` 的主维护源
- `prototypes/ani-services-prototype-console.html` 只保留摘要、边界和详细材料入口
- `tasks/modules/prd/console/compute/prd-console-file-storage.md` 与 `tasks/modules/spec/console/compute/spec-console-file-storage.md` 作为辅助材料保留
- 如正文与辅助材料冲突，以本文件为准

## Core 层要求

- `StorageFilesystem` 资源对象属于 `Core`
- 查询与动作接口必须使用 `/api/v1/filesystems*`
- 不允许把文件存储资源对象写入 `Services /api/v1/svc/*`
- 不允许继续使用旧 `/api/v1/storage/*` 或 `/api/v1/console/*` 路径
- 不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- 创建接口必须使用现有 `Core v1.yaml` 已冻结的 `idempotency_key`

## 冻结事实表

### Frozen Paths

- `GET /api/v1/filesystems`
- `POST /api/v1/filesystems`
- `GET /api/v1/filesystems/{filesystem_id}`
- `DELETE /api/v1/filesystems/{filesystem_id}`

### Frozen Schemas

- `StorageFilesystem`
- `StorageFilesystemListResponse`
- `CreateStorageFilesystemRequest`
- `StorageResourceState`

### Non-Frozen Capabilities

- 挂载目标
- 活跃挂载关系详情
- 访问策略细化

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 冻结能力概览 | Core | 当前 `Filesystems` 路径组 | 列表区 |
| 文件存储列表区 | Core | 展示名称、协议、容量、endpoint、状态 | 详情区 |
| 文件存储创建区 | Core | 对齐 `CreateStorageFilesystemRequest` | 创建抽屉 |
| 待补能力说明 | 规划项 | 仅说明当前未冻结能力 | 后续模块文档 |

## 字段级定义

| 字段 | 说明 | 展示建议 |
|---|---|---|
| id | 文件系统 ID | 文本 |
| tenant_id | 租户 ID | 仅视为后端返回字段，不作为前端必传 |
| name | 文件系统名称 | 文本 |
| protocol | 协议 | `nfs / cephfs` 标签 |
| size_gib | 容量（GiB） | 数值 + 单位 |
| endpoint | 访问端点 | 文本；为空时展示 `-` |
| state | 资源状态 | 标签展示 |
| reason | 状态原因 | 详情页展示 |
| created_at | 创建时间 | 时间 |
| updated_at | 更新时间 | 时间 |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录与租户上下文 | 已认证 | `401` / `403` |
| `name` | 租户内唯一 | `409 CONFLICT` |
| `size_gib` | ≥ 1 | `400 BAD_REQUEST` |

## 操作可用性矩阵

基于 `StorageResourceState`：`pending` / `available` / `failed` / `deleting` / `deleted`。

| 操作 | pending | available | failed | deleting | deleted |
|---|---|---|---|---|---|
| 查看详情 | ✅ | ✅ | ✅ | ⚠️ 只读 | ❌ |
| 创建 | — | — | — | — | — |
| 删除 | ❌ | ✅ | ✅ | ❌ | ❌ |

`pending` 期间删除置灰；删除成功返回 `200 + StorageFilesystem`。

## 接口冻结规则

### `GET /api/v1/filesystems`

- `operationId`: `listStorageFilesystems`
- `success`: `200 + StorageFilesystemListResponse`
- `error responses`: `401 UNAUTHORIZED`、`403 FORBIDDEN`

### `POST /api/v1/filesystems`

- `operationId`: `createStorageFilesystem`
- `requestBody.required`: `name`、`size_gib`、`idempotency_key`
- `success`: `201 + StorageFilesystem`
- `error responses`: `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN`

### `GET /api/v1/filesystems/{filesystem_id}`

- `operationId`: `getStorageFilesystem`
- `success`: `200 + StorageFilesystem`
- `error responses`: `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND`

### `DELETE /api/v1/filesystems/{filesystem_id}`

- `operationId`: `deleteStorageFilesystem`
- `success`: `200 + StorageFilesystem`
- `error responses`: `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND`

## 删除前置校验与当前契约边界

- 页面可以提示用户先确认是否存在活跃挂载影响，但当前权威源未冻结“挂载冲突校验”作为正式后端契约
- 当前 `openapi/v1.yaml` 未显式冻结 `409 CONFLICT`
- 因此本轮正文不把 `409` 写成已冻结删除契约

## 挂载目标（Phase 2 已声明）

- 详文：`filesystem-mount-targets.md`
- `GET /api/v1/filesystems/{filesystem_id}/mount-targets` <!-- ADDED-TO-YAML: Core v1.yaml, Phase 2 2026-06-17 -->

## 回填前置依赖

- 挂载目标路径已在 Phase 2 声明；handler 仍为 stub
- 后续若要增加更细粒度的挂载关系或访问策略，必须先补充正式 schema 与路径

## 回填验收标准

- 正文明确区分已冻结能力与 Phase 2 已声明但 stub 的能力
- 正文不再出现旧 `/api/v1/storage/*`、旧 `/api/v1/console/*` 路径
- 正文不要求前端显式传 `tenant_id / X-Tenant-Id`
- 正文不把挂载目标写成**已实现 handler**的正式接口（YAML 已声明，见 `filesystem-mount-targets.md`）
- `PRD`、`SPEC`、HTML 摘要与本文件一致

## 产品经理补充定义

### 目标用户与权限视角

- 租户管理员：查看、创建、删除文件存储资源
- 业务成员：查看协议、容量和 endpoint，并在授权前提下执行创建或删除
- 只读用户：仅查看列表与详情

### 默认视图与页面状态

- 默认展示文件存储列表、协议和容量摘要
- 当无文件存储时，页面展示空态并提示当前只承接文件系统本体能力
- 若 endpoint 暂不可用，详情中应显示 `-` 或待就绪说明，而不是伪造访问地址
- 当前未冻结的挂载目标、活跃挂载关系等能力只以说明态出现

### 核心任务流

1. 用户创建文件存储后进入详情确认协议、容量和可见 endpoint
2. 用户通过列表筛选目标文件系统，查看状态和更新时间
3. 用户在满足条件时删除文件存储，并通过结果反馈确认完成

### 跨模块协同

- 与 `存储管理` 协同，作为文件存储子模块下钻页
- 与实例和容器模块仅通过关联引用和跳转协同，不承诺当前已支持挂载操作
- 挂载目标与访问策略细化需待 Core 契约补充后再扩展

### 产品验收补充

- 页面必须把“文件系统本体”与“挂载使用”区分清楚
- endpoint 缺失时要有明确文案，不能误导用户已经可挂载
- 创建、详情、删除三条主流转都必须可被完整追踪
- 本页不得把挂载目标或活跃挂载写成已冻结能力
