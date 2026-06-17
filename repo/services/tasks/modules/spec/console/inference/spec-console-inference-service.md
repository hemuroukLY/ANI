# SPEC: Console 推理服务

> Technical specification derived from: `tasks/modules/prd/console/inference/prd-console-inference-service.md`
> Source of truth: `ANI-main/repo/api/openapi/services/v1.yaml`
> Main doc: `docs/console-modules/inference/inference-service.md`
> Scope: `Console / 推理服务`

## 1. Summary

本 SPEC 定义 `Console / 推理服务` 的技术边界、字段映射、页面承接范围和接口冻结规则。推理服务属于 `Services / InferenceServices`，正式前缀为 `/api/v1/svc`，不扩写为 `Core` 资源。

## 2. Source of Truth

### 2.1 Primary Authority

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Frozen Path Group

- `/inference-services`
- `/inference-services/{service_id}`

### 2.3 Frozen Schemas

- `InferenceService`
- `ErrorResponse`

## 3. API Design

### 3.1 Endpoint Matrix

| Method | Path | operationId | summary | Success | Error Responses |
|---|---|---|---|---|---|
| GET | `/api/v1/svc/inference-services` | `listInferenceServices` | 获取推理服务列表 | `200 + {items: InferenceService[]}` | 当前权威源未单独列出 |
| POST | `/api/v1/svc/inference-services` | `createInferenceService` | 部署推理服务 | `202 + InferenceService` | 当前权威源未单独列出 |
| GET | `/api/v1/svc/inference-services/{service_id}` | `getInferenceService` | 获取推理服务详情 | `200 + InferenceService` | `404` |
| PATCH | `/api/v1/svc/inference-services/{service_id}` | `patchInferenceService` | 扩缩容/变配 | `202 + InferenceService` | `400/401/403/404/422` |
| DELETE | `/api/v1/svc/inference-services/{service_id}` | `deleteInferenceService` | 停止并删除推理服务 | `202` | `404` |

### 3.1.1 创建前置条件

| 依赖项 | 要求 | 未满足响应 |
|---|---|---|
| 上游模型 | `Model.status=ready` | `422 PRECONDITION_FAILED`（Services YAML 已举例 `MODEL_NOT_READY`） |
| `name` | 租户内唯一 | `409 CONFLICT` |
| GPU 资源（若指定） | 可满足调度 | `422 PRECONDITION_FAILED`（具体 `code` 待 Services 冻结；建议语义：GPU 不可用） |

> 错误码定稿句式见 `docs/console-modules/governance/module-delivery-workflow.md` §2.10；完整矩阵见 `docs/console-modules/inference/inference-service.md`。

### 3.1.2 操作可用性（按 status）

| 操作 | running | 其他状态 |
|---|---|---|
| PATCH 变配 | ✅ | ❌ 产品建议 `422 PRECONDITION_FAILED`（具体 `code` 待 Services 冻结；建议语义：服务状态不允许变配） |
| DELETE | ⚠️ | 视状态 |

### 3.2 Request Rules

#### `POST /api/v1/svc/inference-services`

- 必填字段：
  - `idempotency_key`
  - `name`
  - `model`
- 可选字段：
  - `replicas`
  - `gpu_type`
  - `gpu_count_per_pod`
  - `max_concurrency`
- 字段约束：
  - `idempotency_key`: `uuid`
  - `replicas`: 默认 `1`
  - `gpu_count_per_pod`: 默认 `1`
  - `max_concurrency`: 默认 `8`
- 成功响应：
  - `202 + InferenceService`

#### `GET /api/v1/svc/inference-services/{service_id}`

- 路径参数：
  - `service_id`: `uuid`

#### `DELETE /api/v1/svc/inference-services/{service_id}`

- 路径参数：
  - `service_id`: `uuid`
- 成功响应：
  - `202`
  - 语义为“停止任务已提交”

## 4. Data Model

### 4.1 InferenceService

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `uuid` | 服务唯一标识 |
| `name` | `string` | 服务名称 |
| `model` | `string` | 上游模型标识 |
| `replicas` | `integer` | 副本数，默认 `1` |
| `gpu_type` | `string \| null` | GPU 类型 |
| `gpu_count_per_pod` | `integer` | 每个 Pod 的 GPU 数，默认 `1` |
| `max_concurrency` | `integer` | 最大并发，默认 `8` |
| `status` | `string` | `pending | deploying | running | stopping | stopped | failed` |
| `endpoint_url` | `string \| null` | 服务访问地址 |
| `created_at` | `date-time` | 创建时间 |

## 5. Page Design Constraints

### 5.1 List Page

- 承接服务列表、部署入口、详情入口
- 列表字段必须来自 `InferenceService`
- 不在列表页承诺 OpenAI、限流、观测资源的正式子页

### 5.2 Detail Page

- 承接服务基础信息、配置摘要、状态反馈、endpoint 信息
- 删除结果回流到当前详情或列表
- 不把未冻结能力扩写成独立资源页

### 5.3 State Handling

- 空态：无服务时展示部署入口和模型准备说明
- 进行中态：`pending`、`deploying`、`stopping`
- 稳态：`running`、`stopped`
- 失败态：`failed`
- 删除请求提交成功仅表示任务已提交，不等于资源已立即消失

## 6. Rules

- 推理服务业务对象归 `Services`
- 生命周期动作仍视为主资源动作，不扩写独立资源域
- OpenAI 兼容 API、限流策略、观测资源、调用测试当前只作为模块边界，不写成已冻结独立接口
- 部署成功返回 `202 + InferenceService`，不能写成 `201`
- 删除成功返回 `202`，不能写成 `204`

## 7. Error Handling

### 7.1 Unified Error Shape

```json
{"code":"UPPER_SNAKE","message":"...","request_id":"..."}
```

### 7.2 Expected Errors

- `GET /api/v1/svc/inference-services/{service_id}`: `404 + ErrorResponse`
- `DELETE /api/v1/svc/inference-services/{service_id}`: `404 + ErrorResponse`

### 7.3 Documentation Rule

- 当前权威源未列出的错误返回码，不允许擅自补写为正式冻结事实
- 不允许把部署写成同步成功，也不允许把删除写成同步完成
- 前置条件 `code` 口径见 `docs/console-modules/governance/module-delivery-workflow.md` §2.10

## 8. Acceptance

- 文档中所有正式路径均位于 `/api/v1/svc/inference-services*`
- 文档中所有成功返回码与 `services/v1.yaml` 保持一致
- 文档中不再把部署写成 `201`、把删除写成 `204`
- 文档中不把待补能力写成正式接口或正式 schema

## 9. Backfill Dependencies

- 如未来新增 OpenAI 兼容协议、限流策略、观测资源，需先冻结路径与 schema
- 如未来新增删除任务历史或部署事件流，需先明确资源归属
