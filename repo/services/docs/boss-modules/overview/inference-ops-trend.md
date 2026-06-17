# AI 服务运营态势

## 页面定位

`AI 服务运营态势` 是 BOSS **平台运营总览** 域下的 **全平台推理服务运营** 趋势明细页，展示推理服务规模、部署态分布、QPS/P99/错误率趋势与跨租户活跃排行。

Console 对照：[`home-inference-status.md`](../../console-modules/home/home-inference-status.md)。深度专页：[`inference-monitoring.md`](../health/inference-monitoring.md)。

## 文档管理规则

- 本文是 `AI 服务运营态势` 的主维护源
- PRD/SPEC 为辅助；冲突以本文 + OpenAPI 为准
- 流程：ANI-14 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- **Services** `GET /api/v1/svc/inference-services`（`listInferenceServices`）— 租户 JWT **只读参考**；**非** BOSS 正式契约
- **Core** `GET /api/v1/observability/query` — PromQL；**非** 推理 list
- BOSS 跨租户 inference aggregate — **TODO-YAML**
- **不得** 逐租户切换 JWT 轮询冒充平台大盘
- 平台 RBAC — **TODO-YAML**
- `AsyncTask.status` 用 `completed`，禁止 UI 使用 `succeeded`

## 页面职责

- 展示全平台推理服务总数、运行中/失败/部署中分布
- QPS、P99、错误率趋势（PromQL 或 aggregate 待 YAML）
- 跨租户推理活跃 Top N（**TODO-YAML**）
- 跳转推理监控、推理调用审计、平台 Input/Output Tokens 计量
- 状态命名与 Services `InferenceService` schema 一致

## 页面结构

```text
AI 服务运营态势
├── 顶部 KPI（总数 / 运行中 / 失败 / 部署中）
├── 趋势图（QPS / P99 / error_rate）
├── 跨租户活跃排行（待 YAML）
├── 筛选（时间窗口、状态、模型）
└── 跳转
    ├── 推理监控
    ├── 推理调用审计
    └── 平台 Input/Output Tokens
```

## 数据来源与分层约束

| 层 | 路径 | 用法 |
|---|---|---|
| Services | `GET /api/v1/svc/inference-services` | 租户 **只读参考** |
| Core | `GET /api/v1/observability/query` | PromQL |
| Core | 平台 inference aggregate **TODO-YAML** | BOSS 正式数据源 |
| Core | `GET /api/v1/tasks/{task_id}` | 部署任务 **路径已声明** |

### 关键边界

- Services list 为租户边界；BOSS 须 platform API
- PromQL **不等于** inference list
- 失败服务计数须与 `InferenceService.status` 一致，禁止自造状态

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| 规模 KPI | aggregate **待 YAML** | inference-monitoring |
| 状态分布 | Services 参考 / aggregate | inference-monitoring |
| QPS/P99 趋势 | PromQL 或 aggregate | observability/query |
| 租户排行 | **TODO-YAML** | tenant-list |
| 失败服务列表摘要 | aggregate | inference-monitoring |

## BOSS 与 Console 分工

| 维度 | BOSS | Console home-inference-status |
|---|---|---|
| 范围 | 全平台推理运营 | 当前租户推理摘要 |
| 排行 | 跨租户 | 无 |
| API | platform aggregate 待 YAML | `listInferenceServices` |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/svc/inference-services` | `listInferenceServices` | 租户 Services |
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL |
| GET | `/api/v1/tasks/{task_id}` | — | **路径已声明** |

| 能力 | 状态 |
|---|---|
| 平台 inference aggregate | **ADDED-TO-YAML** `getPlatformInferenceTrend` |
| 平台 curated PromQL（QPS/P99） | **TODO-YAML** |

## 字段级定义

| 字段 | 说明 |
|---|---|
| `inference_total` | 全平台推理服务数 |
| `running_count` | `status=running`（以 YAML enum 为准） |
| `failed_count` | 失败态计数 |
| `deploying_count` | 部署中计数 |
| `qps_avg` | 窗口平均 QPS |
| `p99_latency_ms` | P99 毫秒 |
| `error_rate` | 0–1 或 % |
| `tenant_rank[]` | **待 YAML** |
| `last_refreshed_at` | UI |
| `time_window` | UI |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 正常 | KPI + 趋势 |
| aggregate 未就绪 | 标注 TODO-YAML |
| `failed_count` > 0 | 高亮 |
| 单 PromQL 失败 | 该曲线失败，KPI 可保留 |
| 403 | 无权限 |

## 字段口径与单位

| 字段 | 口径 |
|---|---|
| `p99_latency_ms` | 毫秒 |
| `error_rate` | 错误/总请求 |
| `qps_avg` | 请求/秒，窗口平均 |

## 状态与能力口径

`InferenceService.status` 以 `services/v1.yaml` enum 为准；本页只汇总，不改写状态机。本页 **只读**。

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| 平台读 RBAC（**TODO-YAML**） | `403` |
| 未认证 | `401` |
| PromQL 非空 | `400` |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 管理员 |
|---|---|---|---|
| 查看态势 | ✅ | ✅ | ✅ |
| 跳转推理监控 | ✅ | ✅ | ✅ |
| 在本页部署/回滚 | ❌ | ❌ | ❌ |

## 删除前置校验

**N/A**

## 接口冻结规则

### `GET /api/v1/svc/inference-services`（租户 · **只读参考**）

- `success`：`200`；当前 YAML GET 未声明额外错误 responses
- 认证/授权错误为 Services 网关通用推断，**未** 在本 operation responses 中冻结
- **非** BOSS 正式跨租户契约

### `GET /api/v1/observability/query`

- `x-ani-rbac-scope`：`scope:observability:read`
- `errors`：`400`、`401`、`403`

### 平台 inference aggregate（待补）

<!-- TODO-YAML: GET /api/v1/observability/inference/platform 或等价 -->

## 使用规则

- 状态命名与推理监控专页一致
- 禁止跨租户 JWT 轮询
- 不得把 PromQL 结果伪造为服务 list

## 待补能力边界

- 平台 inference aggregate — **ADDED-TO-YAML**
- curated PromQL — **TODO-YAML**
- 自动扩缩容 — Phase 2

## 响应示例

### 平台 aggregate 目标（**待 YAML**）

```json
{
  "inference_total": 340,
  "running_count": 310,
  "failed_count": 8,
  "deploying_count": 12,
  "qps_avg": 1250.5,
  "p99_latency_ms": 89,
  "error_rate": 0.003,
  "time_window": "24h",
  "top_tenants": [{ "tenant_id": "t-001", "inference_count": 42 }]
}
```

## 错误示例

### 无 observability 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied: scope:observability:read required",
  "request_id": "req-boss-iot-403-001"
}
```

### PromQL 无效

```json
{
  "code": "BAD_REQUEST",
  "message": "invalid PromQL expression",
  "request_id": "req-boss-iot-400-001"
}
```

## 相关模块

- [`operations-overview.md`](operations-overview.md)、[`inference-monitoring.md`](../health/inference-monitoring.md)
- [`platform-input-tokens.md`](../metering/platform-input-tokens.md)、[`platform-output-tokens.md`](../metering/platform-output-tokens.md)

## 回填验收标准

- [x] 满配章节齐全
- [x] Services 租户路径 vs 平台 aggregate 分层
- [x] 400 + 403 错误示例
- [ ] platform aggregate YAML 合入后回写
- [x] PRD/SPEC/HTML 与本文同步
