# 模型中心

## 页面定位

`模型中心` 是 `Console` 的租户侧模型管理页面，用于管理模型元数据、模型导入和模型版本。

当前模块属于 `Services / Models`，一级权威源为 `ANI-main/repo/api/openapi/services/v1.yaml`，正式路径前缀为 `/api/v1/svc`。

## 文档管理规则

- 本文是 `模型中心` 的主维护文档，页面定义、边界口径和验收标准统一以本文为准
- 当 `PRD`、`SPEC`、HTML 摘要与本文出现差异时，先对照 `ANI-main/repo/api/openapi/services/v1.yaml`，再回写相关材料
- 任何新增能力必须先判断是否已经进入 `Services / Models` 正式冻结范围，未冻结能力只能写入待补边界
- 页面补充只允许增强可读性与交互说明，不得新增 `Core` 资源表述或自造子资源

## Services 层要求

- 模型中心属于 `Services`，不属于 `Core`
- 正式路径使用 `/api/v1/svc/models*`
- 页面不要求前端显式传 `tenant_id`
- 导入模型按异步任务处理
- 当前模块不把模型加密、推荐配置、调用统计写成已冻结能力
- 错误结构统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- 创建、导入、版本创建等写操作必须承接 `idempotency_key`

## 页面职责

- 展示模型列表
- 创建模型元数据
- 发起模型导入
- 创建模型版本
- 查看和删除模型

## 页面结构

- 列表页至少包含 `顶部动作区`、`模型列表`、`筛选区`、`详情入口`
- 详情页至少包含 `基础信息`、`版本信息`、`最近导入任务或结果反馈`
- 导入相关反馈必须在当前页面上下文内闭环呈现，不依赖用户跳转到未冻结的新子页

## 数据来源与分层约束

### 数据来源划分

- `Services` 数据
  - 模型列表
  - 模型创建
  - 模型导入任务
  - 模型详情
  - 模型删除
  - 模型版本创建
- `Core` 数据
  - 当前页面不承接 `Core` 模型资源契约
  - 如后续需要展示对象存储、任务中心等摘要，只能以跳转或引用形式承接，不重写其资源定义

### 关键边界

- 模型中心不把模型对象写成 `Core` 资源
- 模型导入任务返回 `AsyncTask`，不等于模型已进入 `ready`
- 模型版本当前只承接创建结果，不扩写独立版本历史资源页
- 模型加密、推荐配置、调用统计只保留待补边界说明

## 字段级定义

### 模型列表与详情字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `id` | `Model.id` | 模型唯一标识 |
| `name` | `Model.name` | 模型系统名 |
| `display_name` | `Model.display_name` | 模型展示名 |
| `source` | `Model.source` | `upload / huggingface / modelscope / builtin` |
| `capabilities` | `Model.capabilities` | 模型能力标签数组 |
| `status` | `Model.status` | `pending / downloading / ready / error / deleted` |
| `total_size_bytes` | `Model.total_size_bytes` | 模型大小 |
| `created_at` | `Model.created_at` | 创建时间 |

### 模型版本字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `id` | `ModelVersion.id` | 版本唯一标识 |
| `model_id` | `ModelVersion.model_id` | 所属模型 ID |
| `version` | `ModelVersion.version` | 版本标识 |
| `format` | `ModelVersion.format` | `safetensors / gguf / pytorch` |
| `is_encrypted` | `ModelVersion.is_encrypted` | 是否加密 |
| `size_bytes` | `ModelVersion.size_bytes` | 版本大小 |
| `created_at` | `ModelVersion.created_at` | 版本创建时间 |

### 导入任务反馈字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `id` | `AsyncTask.id` | 任务 ID |
| `idempotency_key` | `AsyncTask.idempotency_key` | 幂等键 |
| `task_type` | `AsyncTask.task_type` | 任务类型 |
| `status` | `AsyncTask.status` | `pending / running / completed / failed / cancelled / dead_letter` |
| `progress_pct` | `AsyncTask.progress_pct` | 进度百分比 |
| `error_message` | `AsyncTask.error_message` | 失败原因 |
| `created_at` | `AsyncTask.created_at` | 提交时间 |
| `completed_at` | `AsyncTask.completed_at` | 完成时间 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 模型空态 | 展示创建和导入入口 | 不伪造示例数据 |
| 列表正常态 | 展示模型字段 + 详情入口 | 字段来源必须可回指到 `Model` |
| 导入异步态 | 展示任务状态、进度和失败原因 | 不把“已提交”写成“已完成导入” |
| 详情缺失态 | 返回 `NOT_FOUND` 时展示找不到模型 | 不伪造占位详情 |
| 删除进行态 | 保留当前上下文并提示删除处理中 | 删除成功前不从界面静默移除 |

## 字段口径与单位

| 字段 | 口径建议 | 单位/格式 |
|---|---|---|
| `name` | 系统标识，用于稳定引用 | 小写字符串 |
| `display_name` | 用户可读名称 | 1-128 字符 |
| `total_size_bytes` | 模型总大小 | bytes，前端可格式化为 KB / MB / GB |
| `size_bytes` | 模型版本大小 | bytes，前端可格式化为 KB / MB / GB |
| `created_at` | 服务端创建时间 | date-time |
| `progress_pct` | 异步任务当前进度 | 百分比整数 |

## 状态与能力口径

### 模型状态

- `pending`：模型记录已创建，尚未完成可用准备
- `downloading`：导入或下载处理中
- `ready`：模型可供下游引用
- `error`：模型进入失败态
- `deleted`：模型已删除或不可再用

### 导入任务状态

- `pending`
- `running`
- `completed`
- `failed`
- `cancelled`
- `dead_letter`

### 待补能力边界

- 模型加密与密钥绑定
- 模型部署推荐配置
- 模型调用与使用情况
- 平台审核流
- 模型版本历史独立资源页

## 创建前置条件

- 用户已登录并具备当前租户的模型查看权限
- 创建模型、导入模型、创建版本、删除模型需要具备对应写权限
- 写操作必须包含有效 `idempotency_key`
- 导入来源只允许 `huggingface` 或 `modelscope`

### 下游引用前置条件（供推理服务等消费方）

| 依赖项 | 要求状态 | 未满足时对下游的影响 |
|---|---|---|
| 模型版本 | `Model.status=ready` | 推理服务创建：`422 PRECONDITION_FAILED`（Services YAML 已举例 `MODEL_NOT_READY`） |
| 模型导入任务 | `AsyncTask.status=succeeded` | 模型尚未 ready，不可被引用 |

## 操作可用性矩阵

### 按模型状态（`Model.status`）

| 操作 | pending | downloading | ready | error | deleted |
|---|---|---|---|---|---|
| 查看详情 | ✅ | ✅ | ✅ | ✅ | ❌ |
| 创建元数据 | — | — | — | — | — |
| 发起导入 | ✅ | ❌ | ❌ | ✅ | ❌ |
| 创建版本 | ✅ | ❌ | ✅ | ❌ | ❌ |
| 删除 | ✅ | ⚠️ | ⚠️ 需确认无下游引用 | ✅ | ❌ |

### 按角色（RBAC）

| 操作 | 只读用户 | 模型管理员 |
|---|---|---|
| 查看列表/详情 | ✅ | ✅ |
| 创建/导入/版本/删除 | ❌ | ✅ |

## 删除前置校验与当前契约边界

- 当前正式契约仅确认 `DELETE /api/v1/svc/models/{model_id}` 返回 `204`
- 当前权威源未声明“被推理服务引用时的删除冲突规则”，因此正文不能自造前置校验逻辑
- 若后续权威源补充删除冲突、依赖校验或软删除语义，需再更新本文

## 接口冻结规则

### `GET /api/v1/svc/models`

- `operationId`: `listModels`
- `success`: `200`
- `query`: `status`、`limit`、`cursor`
- `response`: `CursorPage + items: Model[]`
- `errors`: `401`、`403`

### `POST /api/v1/svc/models`

- `operationId`: `createModel`
- `success`: `201 + Model`
- `requestBody.required`: `idempotency_key`、`name`、`display_name`
- `optional fields`: `description`、`capabilities`
- `errors`: `400`、`401`、`409`

### `POST /api/v1/svc/models/import`

- `operationId`: `importModel`
- `success`: `202 + AsyncTask`
- `headers`: `Location`
- `requestBody.required`: `source`、`repo_id`、`idempotency_key`
- `errors`: `400`、`401`

### `GET /api/v1/svc/models/{model_id}`

- `operationId`: `getModel`
- `success`: `200 + Model`
- `errors`: `404`

### `DELETE /api/v1/svc/models/{model_id}`

- `operationId`: `deleteModel`
- `success`: `204`
- `errors`: `404`

### `POST /api/v1/svc/models/{model_id}/versions`

- `operationId`: `createModelVersion`
- `success`: `201 + ModelVersion`
- `requestBody.required`: `idempotency_key`、`version`、`format`、`storage_path`、`checksum_sha256`、`size_bytes`
- `optional fields`: `is_encrypted`
- 当前权威源未单独列出错误返回码，正文不得擅自补写

## 响应示例

### 创建模型成功

```json
{
  "id": "4f1d74a2-8a9f-4f26-b4f6-90d62d8b73d1",
  "name": "llama3-8b",
  "display_name": "LLaMA 3 8B",
  "source": "huggingface",
  "capabilities": ["text-generation"],
  "status": "pending",
  "total_size_bytes": 0,
  "created_at": "2026-06-14T10:00:00Z"
}
```

### 导入任务提交成功

```json
{
  "id": "51debb30-662d-4a14-a1ca-0e16a45f2d4b",
  "idempotency_key": "a4f454c7-16e0-4e88-a3a2-568e6623ddab",
  "task_type": "model_import",
  "resource_type": "model",
  "resource_id": null,
  "status": "pending",
  "attempt_count": 0,
  "max_attempts": 3,
  "progress_pct": 0,
  "result": null,
  "error_message": null,
  "created_at": "2026-06-14T10:05:00Z",
  "completed_at": null
}
```

## 错误示例

### 创建模型参数错误

```json
{
  "code": "BAD_REQUEST",
  "message": "display_name is required",
  "request_id": "req-model-400-001"
}
```

### 创建模型名称冲突

```json
{
  "code": "CONFLICT",
  "message": "model name already exists",
  "request_id": "req-model-409-001"
}
```

### 模型不存在

```json
{
  "code": "NOT_FOUND",
  "message": "model not found",
  "request_id": "req-model-404-001"
}
```

## 回填前置依赖

- 如未来补充模型加密、推荐配置、调用统计，必须先在 `services/v1.yaml` 中冻结路径和 schema
- 如未来补充模型版本历史页，必须先确认是否仍复用 `ModelVersion`
- 如未来补充导入任务历史或详情页，必须先明确与 `AsyncTask` 的关系

## 待确认项

- 列表是否需要未来扩充版本数量或最近版本摘要字段
- 删除模型是否需要未来补充“被推理服务引用时的冲突规则”
- 导入结果是否需要未来独立任务历史页

## 回填验收标准

- 正文路径全部位于 `Services / Models`
- 成功返回码与 `services/v1.yaml` 保持一致：创建 `201`、导入 `202`、详情 `200`、删除 `204`、版本创建 `201`
- 正文只引用权威源真实存在的 schema：`Model`、`ModelVersion`、`AsyncTask`、`ErrorResponse`
- 正文不再出现 `[Assumption]`、`201/200`、`202/200`、`200/204` 等模糊口径
- HTML 摘要、PRD、SPEC 与本文件一致

## 产品经理补充定义

### 目标用户与权限视角

- 模型管理员：查看模型列表、创建模型元数据、发起模型导入、创建模型版本、删除模型
- 服务 Owner：查看模型与版本信息，作为推理服务选型时的上游参考
- 只读用户：仅查看模型与版本信息，不显示创建、导入、删除入口

### 页面结构补充

- 首屏至少包含 `顶部动作区`、`模型列表`、`详情入口`
- 模型详情至少拆成 `基础信息`、`版本信息`、`最近导入任务或结果反馈` 三个区块
- 模型导入必须以异步任务或任务结果反馈呈现，不允许表现成同步完成上传

### 核心任务流

1. 用户先创建模型元数据，再进入详情补充版本或继续发起导入动作
2. 用户发起模型导入后，通过任务反馈查看排队、执行、成功或失败结果，并回流到模型详情
3. 用户基于已存在模型创建新版本，页面刷新版本区块并保留最新结果反馈
4. 用户删除模型时，如后端返回冲突或依赖存在，页面按统一错误结构提示，不自行编造额外约束

### 页面状态与反馈

- 空态：当前无模型时，展示创建和导入两个主入口
- 局部失败态：列表、详情、导入结果区互不阻塞
- 异步态：导入任务必须显示处理中反馈，不允许只有“提交成功”后无后续状态
- 成功回流：创建、导入、版本创建后都必须回流到模型详情或当前上下文

### 跨模块协同

- 与 `推理服务` 协同，模型详情可作为推理服务创建时的上游引用来源，但不在本页直接创建推理服务
- 与首页协同，首页只展示模型摘要或异常，不替代本页详情和导入反馈
- 模型加密、推荐配置、调用统计等待补能力只能保留边界说明

### 产品验收补充

- 用户必须能完成“创建模型元数据 -> 导入或创建版本 -> 在详情中复核结果”的闭环
- 异步导入必须具备明确的过程态和失败反馈
- 模型详情必须能解释当前模型与版本的关系，不能只有路径清单
- 本页不得把模型写成 `Core` 资源，也不得扩写未冻结的模型治理能力
