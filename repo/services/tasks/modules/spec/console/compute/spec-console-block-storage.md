# SPEC: Console 块存储

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-block-storage.md`
> Generated: 2026-06-13 | Scope: `Console / 算力与云资源 / 存储管理 / 块存储`

## 1. Summary

### 1.1 What This SPEC Covers

本 SPEC 说明 `Console / 算力与云资源 / 存储管理 / 块存储` 的技术边界、数据模型、接口冻结规则、错误处理方式和待补能力边界。目标是把块存储模块收口成一份可直接用于对齐 `Core v1.yaml` 的主维护材料，并明确区分“当前 Core 已冻结能力”与“仍待补充到 Core 的规划能力”。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/compute/prd-console-block-storage.md`
- User Stories covered: `US-001` ~ `US-005`
- Functional Requirements covered: `FR-1` ~ `FR-7`

### 1.3 Design Decisions Summary

| Decision | Choice | Rationale |
|---|---|---|
| 模块归属 | `Console` 租户侧资源管理页 | 与 `BOSS` 平台存储池运营页分离 |
| 资源归属 | `StorageVolume` 归 `Core` | 现有 `Core v1.yaml` 已冻结对应 schema 和路径 |
| 能力范围 | 只承接列表、详情、创建、删除 | 现有 `Core` 只冻结 `/volumes` 和 `/volumes/{volume_id}` |
| 快照处理 | 标记为待补 `Core` 能力 | 当前 `Core v1.yaml` 未冻结快照路径 |
| 卷挂载/卸载处理 | 标记为待补 `Core` 能力 | 当前 `Core v1.yaml` 未冻结相关动作路径 |
| 删除错误码 | 以当前 `Core v1.yaml` 为准 | 现有 `openapi/v1.yaml` 未显式冻结 `409 CONFLICT` |

## 2. Architecture

### 2.1 System Context

- 页面属于 `Console / 算力与云资源 / 存储管理`
- 页面直接消费 `Core /api/v1/*`
- 页面不定义新的 `Services /api/v1/svc/*` 块存储资源契约
- VM、容器、GPU 容器等页面中的卷资源摘要和跳转只消费本模块已冻结的卷资源定义

### 2.2 Component Design

- `块存储总览`
  - 展示卷资源边界和当前可承接能力
  - 展示哪些能力已经 Core 冻结，哪些仍待补
- `块存储列表区`
  - 承接卷列表查询和状态展示
- `块存储详情区`
  - 承接卷基础信息展示
- `块存储创建区`
  - 承接创建表单和请求约束说明
- `待补能力说明区`
  - 快照
  - 卷挂载
  - 卷卸载

### 2.3 Module Interactions

1. 用户打开块存储页
2. 页面展示卷资源边界和能力状态
3. 用户查询卷列表并打开卷详情
4. 用户通过创建抽屉或表单发起卷创建
5. 用户删除卷前查看风险提示
6. 页面对待补能力只展示说明，不伪造调用路径

### 2.4 File Structure

```text
docs/
└── console-modules/
    └── compute/
        └── storage/
            └── block-storage.md             [NEW]

tasks/
├── prd-console-block-storage.md            [NEW]
└── spec-console-block-storage.md           [NEW]

root/
├── prototypes/ani-services-prototype-console.html     [MODIFY]
├── prototypes/ani-services-prototype.html             [MODIFY]
└── docs/console-modules/README.md          [MODIFY]
```

## 3. Data Model

### 3.1 Core Frozen Schemas

已冻结并可直接引用 `Core v1.yaml` 的 schema：

- `StorageVolume`
- `StorageVolumeListResponse`
- `CreateStorageVolumeRequest`
- `StorageResourceState`

### 3.2 Entity Definitions

#### StorageVolume

- 关键字段：`id`、`tenant_id`、`name`、`size_gib`、`storage_class`、`state`、`reason`、`created_at`、`updated_at`
- 可选开发态辅助字段：`dev_profile`

### 3.3 Non-Frozen Data Areas

以下能力在当前 `Core v1.yaml` 中未冻结，不应写成正式 Core 契约：

- 卷快照能力
- 卷挂载能力
- 卷卸载能力

### 3.4 Naming Convention

- 页面展示字段允许使用 camelCase 作为前端 view model
- API 路径参数、query 参数、请求体和响应 JSON schema 统一以现有 `Core v1.yaml` 命名为准
- 不把旧 HTML 中的自定义卷动作字段直接写成已冻结接口参数

## 4. API Design

### 4.1 Frozen Endpoints

| Method | Path | Description | Auth | Request | Response |
|---|---|---|---|---|---|
| GET | `/api/v1/volumes` | 查询块存储卷列表 | Bearer / ApiKey | `limit?`、`cursor?` | `200 + StorageVolumeListResponse` |
| POST | `/api/v1/volumes` | 创建块存储卷 | Bearer / ApiKey | `CreateStorageVolumeRequest` | `201 + StorageVolume` |
| GET | `/api/v1/volumes/{volume_id}` | 查询块存储卷 | Bearer / ApiKey | path: `volume_id` | `200 + StorageVolume` |
| DELETE | `/api/v1/volumes/{volume_id}` | 删除块存储卷 | Bearer / ApiKey | path: `volume_id` | `200 + StorageVolume` |

### 4.2 Per-Endpoint Frozen Rules

#### `GET /api/v1/volumes`

| 项 | 值 |
|---|---|
| operationId | `listStorageVolumes` |
| summary | `查询块存储卷列表` |
| tags | `["Volumes"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| query params | `limit?`、`cursor?` |
| success | `200 + StorageVolumeListResponse` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN` |

#### `POST /api/v1/volumes`

| 项 | 值 |
|---|---|
| operationId | `createStorageVolume` |
| summary | `创建块存储卷` |
| tags | `["Volumes"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| preconditions | 已通过标准鉴权；`name`、`size_gib`、`idempotency_key` 非空；`size_gib` 满足 schema 约束 |
| requestBody.required | `name`、`size_gib`、`idempotency_key` |
| success | `201 + StorageVolume` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN` |

#### `GET /api/v1/volumes/{volume_id}`

| 项 | 值 |
|---|---|
| operationId | `getStorageVolume` |
| summary | `查询块存储卷` |
| tags | `["Volumes"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `volume_id` |
| success | `200 + StorageVolume` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

#### `DELETE /api/v1/volumes/{volume_id}`

| 项 | 值 |
|---|---|
| operationId | `deleteStorageVolume` |
| summary | `删除块存储卷` |
| tags | `["Volumes"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `volume_id` |
| success | `200 + StorageVolume` |
| error responses | `401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

### 4.3 Non-Frozen Endpoints

以下路径当前**不能**写成已对齐 Core 的正式接口：

- `/api/v1/volumes/{volume_id}/snapshots`
- `/api/v1/volumes/{volume_id}:attach`
- `/api/v1/volumes/{volume_id}:detach`

如果后续需要补这些路径，应先新增到 `Core v1.yaml`，再回写模块正文。

### 4.4 Operation Availability Matrix

| 资源 | list/get | create | delete | action | 前置校验 |
|---|---|---|---|---|---|
| 块存储卷 | 可用 | 可用 | 可用 | 无额外冻结动作 | 页面可以提示用户先确认是否存在潜在引用影响；当前冻结返回码以 `Core v1.yaml` 为准 |

### 4.5 Deletion Dependency Handling

- 页面可以提示用户先确认是否存在潜在引用影响，但当前权威源未冻结“引用冲突校验”作为正式后端契约
- 当前 `openapi/v1.yaml` 与 `core-v1-compatibility-baseline.yaml` 对 `DELETE /api/v1/volumes/{volume_id}` 显式冻结的返回码为 `200`、`401`、`403`、`404`
- 因此本轮文档**不能**把 `409 CONFLICT` 写成已冻结块存储删除契约
- 若后续要把冲突语义显式收入口径，必须先更新 `Core v1.yaml`，再同步回写本模块正文和辅助材料

### 4.6 Request/Response Examples

#### CreateStorageVolumeRequest

```json
{
  "idempotency_key": "idem-volume-001",
  "name": "volume-data-01",
  "size_gib": 200,
  "storage_class": "standard"
}
```

#### StorageVolume

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

## 5. Business Logic

### 5.1 Core Page Logic

- 块存储页默认展示当前租户可见卷资源
- 总览区负责说明卷资源边界与能力状态
- 详情区只展示 `StorageVolume` 当前已冻结字段
- 创建区只承接现有 `CreateStorageVolumeRequest`
- 对未冻结能力，只记录为“待补 Core 契约”，不能写成正式接口表

### 5.2 Validation Rules

- 创建卷必须提供 `name`、`size_gib` 和 `idempotency_key`
- `size_gib` 必须符合最小值约束
- `storage_class` 若填写，必须遵循当前 schema 允许的字符串口径
- 前端不可要求用户主动填写 `tenant_id`
- 页面不得继续使用旧 `/api/v1/storage/*` 或 `/api/v1/console/*` 路径

### 5.3 State Handling

- 卷资源状态统一使用 `StorageResourceState`
- 页面只按已有 `state` 和 `reason` 展示标签与异常提示
- 不在本轮引入新的卷状态机

### 5.4 Edge Cases

- 卷详情查询时目标卷不存在
- 卷删除前仍存在活动引用
- 卷创建时容量非法
- 卷创建时 `idempotency_key` 为空

## 6. Error Handling

### 6.1 Error Taxonomy

| Error Code | HTTP Status | Condition | User Message |
|---|---|---|---|
| `BAD_REQUEST` | 400 | 参数格式错误、容量非法 | 请求参数不合法，请检查后重试 |
| `UNAUTHORIZED` | 401 | 未登录或凭证失效 | 登录状态已失效，请重新登录 |
| `FORBIDDEN` | 403 | 无访问或操作权限 | 当前账号无权访问该块存储资源 |
| `NOT_FOUND` | 404 | `volume_id` 不存在 | 目标卷不存在或已删除 |

### 6.2 Shared Error Format

```json
{
  "code": "UPPER_SNAKE",
  "message": "error message",
  "request_id": "req-xxx"
}
```

### 6.3 Failure Modes

- 删除类动作当前按 `Core v1.yaml` 已冻结返回码处理，不额外声明未冻结的 `409`
- 查询失败时仅影响块存储区，不阻断整个页面框架
- 待补能力若无后端契约，页面只能展示占位说明或跳转禁用，不得伪造接口

## 7. Security

### 7.1 Authentication & Authorization

- 全部块存储接口使用 `security: [{BearerAuth: []}, {ApiKeyAuth: []}]`
- 后端必须从认证上下文获取租户边界
- 前端不可信任也不显式传 `tenant_id / X-Tenant-Id`

### 7.2 Input Validation

- `idempotency_key` 必须非空
- `name` 必须符合长度要求
- `size_gib` 必须符合 schema 约束

### 7.3 Data Protection

- 不展示平台侧全局存储池信息
- 只展示当前租户可见卷资源
- 页面不暴露平台内部调度实现细节

## 8. Performance

### 8.1 Expected Load

- 卷列表以租户级分页查询为主
- 卷详情为按需加载

### 8.2 Optimization Strategy

- 列表统一沿用 `limit + cursor`
- 详情按需加载
- 删除前只展示风险提示，不扩写为独立引用详情接口

### 8.3 Database Considerations

- 本文档不新增数据库设计
- 数据模型以现有 `Core` schema 为准

## 9. Testing Strategy

### 9.1 Unit Tests

- 校验创建请求字段映射
- 校验待补能力不会误显示为已冻结接口

### 9.2 Integration Tests

- 卷列表、详情、创建、删除接口映射正确
- 创建请求带有 `idempotency_key`
- 删除契约不额外声明未冻结 `409`

### 9.3 Edge Case Tests

- 查询不存在卷时展示 `NOT_FOUND`
- 创建卷容量非法时展示 `BAD_REQUEST`
- 未鉴权访问时展示 `UNAUTHORIZED`

### 9.4 Acceptance Criteria Mapping

| US/FR | Test | Type | Description |
|---|---|---|---|
| `US-001` | 块存储总览边界校验 | integration | 验证已冻结能力与待补能力被正确区分 |
| `US-002` | 卷查询契约校验 | integration | 验证 `/volumes` 和 `/volumes/{volume_id}` 路径、字段、返回码对齐 |
| `US-003` | 卷创建契约校验 | integration | 验证 `CreateStorageVolumeRequest` 必填字段和返回结构 |
| `US-004` | 卷删除契约校验 | integration | 验证删除返回码口径与当前 `Core v1.yaml` 一致 |
| `US-005` | Core 边界校验 | unit | 验证未冻结能力不会写入正式接口表 |

## 10. Implementation Plan

### 10.1 Phases

1. 明确块存储已冻结与未冻结边界
2. 生成 PRD 和 SPEC
3. 收口主维护源
4. 回填 HTML 摘要
5. 执行最终 Core 合规复审

### 10.2 Issue Mapping

| Issue | SPEC Sections | Priority | Depends On |
|---|---|---|---|
| 块存储边界收口 | 1, 2, 4.3 | high | — |
| Core 已冻结卷接口整理 | 3, 4.1, 4.2 | high | 块存储边界收口 |
| 主维护源落盘 | 5, 6, 7 | high | Core 已冻结卷接口整理 |
| HTML 摘要回填 | 2.4, 10.1 | medium | 主维护源落盘 |

### 10.3 Incremental Delivery

- 本轮先收口主维护材料，不直接扩写新的 `Core v1.yaml`
- 后续若要补“快照”“挂载/卸载”等能力，再进入 `YAML` 扩充阶段

## 11. Open Questions & Risks

### 11.1 Unresolved Questions

- 卷快照后续是否作为独立路径组进入 `Core v1.yaml`
- 卷挂载/卸载是否需要成为下一个块存储侧补充点
- 删除依赖冲突是否后续需要显式扩写为 `409`

### 11.2 Technical Risks

| Risk | Impact | Mitigation |
|---|---|---|
| 旧 HTML 路径仍沿用 `/api/v1/storage/*` | 会造成前后端继续对错路径 | 统一改为现有 `/api/v1/volumes*` 路径 |
| 旧文档继续要求 `X-Tenant-Id` | 会与 Core 网关边界冲突 | 文档统一改为认证上下文获取 |
| 把未冻结能力或未冻结返回码写进正式契约 | 后续 Core 对齐会失真 | 在正文和 SPEC 中单列“待补能力”，并以 `openapi/v1.yaml` 为准回修 |

### 11.3 Frozen Baseline

- 当前 `Core v1.yaml` 就是块存储资源的唯一权威来源
- 页面不新增独立 `Services` 块存储资源聚合接口
- 本轮只收口 `Console / 块存储` 子模块，不同步细化卷快照或挂载流程
