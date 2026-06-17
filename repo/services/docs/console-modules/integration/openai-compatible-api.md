# OpenAI 兼容 API

## 页面定位

面向接入方与 Console 调试的 **OpenAI 兼容调用** 说明页，覆盖 Gateway 根路径 `/v1/chat/completions`（非 `/api/v1/svc`）。

本页归属 `integration/`，为**文档与示例**模块，非 Services OpenAPI 资源。

## 文档管理规则

- 本文是 OpenAI 兼容接入的主维护源
- Gateway 路由以 `ani-gateway` 实现为准；**不在** `services/v1.yaml` 内
- 与 `inference-call-test.md`（Services `/test`）明确分流
- 区分 Gateway `/v1` 与 Core `/api/v1`

<!-- Gateway POST /v1/chat/completions — 无正式 OpenAPI YAML；文档口径 Phase 2 -->

## Services 层要求

- 推理服务资源：`InferenceService.endpoint_url`（`services/v1.yaml`）
- Console 内调试优先 `POST .../test`，生产接入走 Gateway

## Core 层要求

- 鉴权：Bearer / API Key（Core Auth）
- 计量上报可能经 `POST /api/v1/metering/token-usage`（Services 侧调用）

## Gateway 层（文档口径）

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/v1/chat/completions` | OpenAI 兼容 chat（`inferenceProxy` stub） |
| GET | `/v1/inference/stream` | 流式（若 Gateway 支持，对照代码） |

**当前无** Gateway 专用 OpenAPI YAML — 成功/错误码以实现与 OpenAI 兼容惯例为准，不得自造 ANI 专属 schema 写入 services/v1.yaml。

## 页面职责

- 文档：请求格式、模型/header 路由、流式 SSE
- 示例 curl / SDK 片段
- 明确 stub ≠ 生产就绪
- 链接 `InferenceService` 详情与 API Key 管理

## 创建前置条件

| 依赖项 | 要求 | 说明 |
|---|---|---|
| 目标服务 | `InferenceService.status=running` | 否则 Gateway 应拒绝或 502 |
| 鉴权 | 有效 API Key 或 Bearer | Core Auth |
| 模型路由 | `X-Model-Name` 或 body.model | 对照 gateway 路由规则 |

不涉及 Console 页面向 Gateway 发 PUT；客户端调用无 `idempotency_key` 要求（OpenAI 惯例）。

## 操作可用性矩阵

| 操作 | 只读用户 | 接入开发者 |
|---|---|---|
| 阅读文档 | ✅ | ✅ |
| 复制示例 | ✅ | ✅ |
| Gateway 生产调用 | 视 API Key | ✅ |
| 在 Console 编辑 Gateway 路由 | ❌ | ❌ |

## 接口冻结规则

### `POST /v1/chat/completions`（Gateway）

- 成功：`200` + OpenAI chat completion JSON（或 SSE 流）
- 错误：参考 OpenAI 惯例（401/429/502 等）；**ANI 未在 repo OpenAPI 冻结**
- **标注**：handler stub，见 `GAP-REPORT` / `CORE-TEAM-TASKS`

## 待补边界

- Gateway OpenAPI 正式 YAML — **TODO-YAML**
- 流式 SSE 事件格式 — 以实现为准
- `/v1/embeddings` 等扩展路径 — 待产品规划

## 与相关模块

- `inference-call-test.md` — Console 内测试
- `open-integration-overview.md` — 接入总入口
- `api-key-management.md` — 鉴权凭据

## 验收标准

- [ ] 不写进 services/v1.yaml 路径表
- [ ] 区分 Gateway `/v1` 与 Core `/api/v1`
- [ ] 标注 Gateway stub / TODO-YAML
