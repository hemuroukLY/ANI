# 推理服务可观测性（日志 / 事件 / 指标）

## 页面定位

`推理服务可观测性` 是 `Console / 推理服务` 下的运维观测页，用于查看单个推理服务的**日志**、**指标**与（待补）**事件**。

本页跨 **Services**（日志）与 **Core**（PromQL 指标代理）两层，必须在 UI 与文档中保持分层清晰。

## 文档管理规则

- 本文是 `推理服务可观测性` 的主维护文档
- Services 日志路径以 `services/v1.yaml` 为准
- Core 指标以 `GET /api/v1/observability/query` 为准
- 不得自造 `inference-services/{id}/metrics` 或 `/events` 路径（当前 YAML 未声明）
- OpenAPI 已声明 ≠ handler 已实现

## Services 层要求 — 日志

- 正式路径：`GET /api/v1/svc/inference-services/{service_id}/logs`
- `operationId`: `getInferenceServiceLogs`
- 响应：`InferenceServiceLogListResponse`（`items` + `next_cursor`）
- Query：`limit`（默认 100，最大 500）、`cursor`
- 错误：`401`、`403`、`404`

<!-- ADDED-TO-YAML: GET .../logs (Services v1.yaml, Phase 2 2026-06-17) -->

### InferenceServiceLogEntry（展示字段）

以 YAML `InferenceServiceLogEntry` 为准；典型展示：时间戳、级别、消息正文。

## Core 层要求 — 指标

- 正式路径：`GET /api/v1/observability/query`
- `operationId`: `queryObservability`
- Query 必填：`query`（PromQL 字符串）
- Query 可选：`time`、`timeout`
- 响应：`ObservabilityQueryResponse`
- RBAC：`scope:observability:read`

Console 须构造与服务相关的 PromQL（具体语句待运维模板冻结）；**不得**在前端硬编码未文档化的 label 集合为冻结事实。

## 待补 — 事件

<!-- TODO-YAML: 推理服务独立 events 路径未在 services/v1.yaml 声明；可暂用 logs 或 Core 审计（若未来冻结） -->

- 当前无 `GET .../events` 冻结路径
- 页面「事件」Tab 应显示「待接入」或隐藏，不得伪造事件列表

## 页面职责

- 按服务展示日志 tail / 分页加载
- 按时间窗口展示 QPS、延迟、错误率等指标图（经 PromQL 代理）
- 提供服务上下文与刷新时间
- 不替代集群级 observability 控制台

## 页面结构

```text
日志 / 事件 / 指标
├── 服务上下文条
├── Tab：日志
│   └── 分页列表 + cursor
├── Tab：指标
│   └── 时间范围 + PromQL 图表
└── Tab：事件（待补）
    └── 空态 / 待接入
```

## 创建前置条件

| 依赖项 | 要求 | 未满足响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 服务存在 | 有效 `service_id` | `404 NOT_FOUND` |
| 读日志权限 | Services RBAC | `403 FORBIDDEN` |
| 读指标权限 | `scope:observability:read` | `403 FORBIDDEN` |
| PromQL 合法 | `query` 非空 | `400 BAD_REQUEST`（指标 Tab） |

本页无 POST/PUT；不涉及 `idempotency_key`。

## 操作可用性矩阵

| 操作 | 只读用户 | 服务管理员 |
|---|---|---|
| 查看日志 | ✅ | ✅ |
| 查看指标 | ✅（有 observability 权限） | ✅ |
| 查看事件 | ❌（未冻结） | ❌（未冻结） |
| 导出日志 | ❌ | ❌（待补 API） |

## 接口冻结规则

### `GET /api/v1/svc/inference-services/{service_id}/logs`

- 成功：`200 + InferenceServiceLogListResponse`
- 错误：`401`、`403`、`404`

### `GET /api/v1/observability/query`

- 成功：`200 + ObservabilityQueryResponse`
- 错误：`400`、`401`、`403`
- Query 必填：`query`

## 待补边界

- 推理服务 events 独立路径 — **TODO-YAML**
- PromQL 模板与服务 label 映射 — 待运维文档冻结
- 日志导出 / 长期归档 — **当前 YAML 未声明**

## 与相关模块的关系

- `inference-service.md`：服务状态与 endpoint
- `inference-call-test.md`：连通性调试，非持续观测
- Core `observability/alert-rules`：告警规则，非本页日志

## 验收标准

- [ ] 日志路径与 Phase 2 YAML 一致
- [ ] 指标走 Core `/observability/query`，不写入 Services YAML
- [ ] 事件 Tab 标注 `TODO-YAML`
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
