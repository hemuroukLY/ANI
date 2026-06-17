# SPEC: Console 模型中心

> Technical specification derived from: `tasks/modules/prd/console/inference/prd-console-model-center.md`
> Source of truth: `ANI-main/repo/api/openapi/services/v1.yaml`
> Scope: `Console / 模型中心`

## 1. Summary

本 SPEC 定义 `Console / 模型中心` 的技术边界、字段映射、页面承接范围和接口冻结规则。模型中心属于 `Services / Models`，正式前缀为 `/api/v1/svc`，不扩写为 `Core` 资源。

## 2. Source of Truth

### 2.1 Primary Authority

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Frozen Path Group

- `/models`
- `/models/import`
- `/models/{model_id}`
- `/models/{model_id}/versions`

### 2.3 Frozen Schemas

- `ErrorResponse`
- `AsyncTask`
- `Model`
- `ModelVersion`

### 2.4 Shared Response Objects

- `BadRequest`
- `Unauthorized`
- `Forbidden`
- `NotFound`
- `Conflict`

## 3. API Design

### 3.1 Endpoint Matrix

| Method | Path | operationId | summary | Success | Error Responses |
|---|---|---|---|---|---|
| GET | `/api/v1/svc/models` | `listModels` | 获取模型列表 | `200` | `401`、`403` |
| POST | `/api/v1/svc/models` | `createModel` | 创建模型元数据记录 | `201 + Model` | `400`、`401`、`409` |
| POST | `/api/v1/svc/models/import` | `importModel` | 从 HuggingFace 或 ModelScope 异步导入模型 | `202 + AsyncTask` | `400`、`401` |
| GET | `/api/v1/svc/models/{model_id}` | `getModel` | 获取模型详情 | `200 + Model` | `404` |
| DELETE | `/api/v1/svc/models/{model_id}` | `deleteModel` | 删除模型 | `204` | `404` |
| POST | `/api/v1/svc/models/{model_id}/versions` | `createModelVersion` | 上传模型版本 | `201 + ModelVersion` | 当前权威源未单独列出错误返回码 |

### 3.2 Query and Request Rules

#### `GET /api/v1/svc/models`

- 查询参数：
  - `status`: `pending | downloading | ready | error | deleted`
  - `limit`: `1-100`，默认 `20`
  - `cursor`: 游标分页参数
- 成功响应：
  - `CursorPage + items: Model[]`

#### `POST /api/v1/svc/models`

- 必填字段：
  - `idempotency_key`
  - `name`
  - `display_name`
- 可选字段：
  - `description`
  - `capabilities`
- 字段约束：
  - `idempotency_key`: `uuid`
  - `name`: 正则 `^[a-z0-9][a-z0-9.-]{0,62}$`
  - `display_name`: `1-128` 字符
  - `description`: 最长 `1024`
  - `capabilities`: 枚举项为 `text-generation | embedding | speech-to-text`

#### `POST /api/v1/svc/models/import`

- 必填字段：
  - `source`
  - `repo_id`
  - `idempotency_key`
- 可选字段：
  - `revision`
  - `webhook_url`
- 字段约束：
  - `source`: `huggingface | modelscope`
  - `repo_id`: `3-256` 字符
  - `revision`: 默认 `main`
  - `idempotency_key`: 最大长度 `128`
- 成功响应：
  - `202 + AsyncTask`
  - `Location` 响应头存在

#### `POST /api/v1/svc/models/{model_id}/versions`

- 必填字段：
  - `idempotency_key`
  - `version`
  - `format`
  - `storage_path`
  - `checksum_sha256`
  - `size_bytes`
- 可选字段：
  - `is_encrypted`
- 字段约束：
  - `format`: `safetensors | gguf | pytorch`
  - `idempotency_key`: `uuid`

## 4. Data Model

### 4.1 Model

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `uuid` | 模型主键 |
| `name` | `string` | 模型系统名 |
| `display_name` | `string` | 模型展示名 |
| `source` | `string` | `upload | huggingface | modelscope | builtin` |
| `capabilities` | `string[]` | 模型能力标签 |
| `status` | `string` | `pending | downloading | ready | error | deleted` |
| `total_size_bytes` | `integer` | 模型总大小 |
| `created_at` | `date-time` | 创建时间 |

### 4.2 ModelVersion

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `uuid` | 版本主键 |
| `model_id` | `uuid` | 所属模型 |
| `version` | `string` | 版本标识 |
| `format` | `string` | `safetensors | gguf | pytorch` |
| `is_encrypted` | `boolean` | 是否加密 |
| `size_bytes` | `integer` | 版本大小 |
| `created_at` | `date-time` | 创建时间 |

### 4.3 AsyncTask

导入任务反馈复用 `AsyncTask`：

- `id`
- `idempotency_key`
- `task_type`
- `resource_type`
- `resource_id`
- `status`
- `attempt_count`
- `max_attempts`
- `progress_pct`
- `result`
- `error_message`
- `created_at`
- `completed_at`

## 5. Page Design Constraints

### 5.1 List Page

- 承接列表读取、筛选、创建入口、导入入口、详情入口
- 列表字段必须来自 `Model`
- 不在列表页承诺模型加密、推荐配置、调用统计子页的正式后端支撑

### 5.2 Detail Page

- 承接模型基础信息、版本信息、最近导入反馈
- 不把未冻结能力扩写成新资源页
- 删除结果与版本创建结果回流到当前详情上下文

### 5.3 State Handling

- 空态：无模型时展示创建和导入入口
- 局部失败态：列表、详情、导入反馈区独立失败
- 异步态：导入必须显示 `AsyncTask.status`
- 删除态：删除成功前不能伪造“已删除”

## 6. Frozen / Non-Frozen Boundary

### 6.1 Frozen Capabilities

- 模型列表
- 模型元数据创建
- 模型导入任务提交
- 模型详情
- 模型删除
- 模型版本创建

### 6.2 Non-Frozen Capabilities

- 模型加密与密钥绑定
- 模型部署推荐配置
- 模型调用与使用情况
- 平台审核流
- 模型版本历史独立资源页

## 7. Error Handling

### 7.1 Unified Error Shape

```json
{"code":"UPPER_SNAKE","message":"...","request_id":"..."}
```

### 7.2 Expected Errors

- 参数错误：`400 + ErrorResponse`
- 未认证：`401 + ErrorResponse`
- 无权限：`403 + ErrorResponse`
- 资源冲突：`409 + ErrorResponse`
- 资源不存在：`404 + ErrorResponse`

### 7.3 Documentation Rule

- 不允许写成 `201/200`、`202/200`、`200/204` 这类模糊成功码
- 只允许引用权威源真实存在的 schema 名称

## 8. Acceptance

- 文档中所有正式路径均位于 `/api/v1/svc/models*`
- 文档中所有成功返回码与 `services/v1.yaml` 保持一致
- 文档中不再出现 `[Assumption]` 和模糊成功码
- 文档中不把待补能力写成正式接口或正式 schema

## 9. Backfill Dependencies

- 如未来新增模型加密、推荐配置、调用统计，需先冻结路径与 schema
- 如未来新增模型审核流，需先明确资源归属
- 如未来新增任务历史页，需先明确与 `AsyncTask` 的关系
