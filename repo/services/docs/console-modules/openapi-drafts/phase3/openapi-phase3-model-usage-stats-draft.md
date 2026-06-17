# Phase 3 — 模型调用与使用统计 OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`inference/model-usage-stats.md`  
> **TASK**：`TASK-SVC-018` 子项（模型增强）  
> **原则**：**优先扩展 Core Metering** 查询 filter；Services 可选 summary 端点便于模型中心 Dashboard。

---

## 1. 与租户用量报表 / Token 上报的边界

| 维度 | 租户用量报表 | 模型使用统计（本草案） |
|---|---|---|
| 页面 | `tenant/usage-report.md` | `inference/model-usage-stats.md` |
| 视角 | 租户全资源聚合 | **单模型 / 模型列表** 视角 |
| 主 API | Core `GET /metering/usage` | 同上 + model filter |
| 写入 | — | `POST /metering/token-usage`（Services 上报，非 Console 查询） |

Console 模型中心页：**不新增**独立计量写入 API；查询复用 Metering。

---

## 2. Core Metering 扩展（草案 · v1.yaml）

> 以下为 **Core 规划扩展**，合入前不得写为已冻结。

### 2.1 扩展 Query 参数

在现有 `GET /api/v1/metering/usage` 上增加：

| 参数 | 类型 | 说明 |
|---|---|---|
| `model_id` | uuid | 筛选单模型相关用量 |
| `inference_service_id` | uuid | 筛选单推理服务 |
| `resource_type` | string | 现有；扩展允许 `model_inference`、`token_input`、`token_output`（**评审固定枚举**） |

`group_by` 扩展（3b）：`model_id`、`inference_service_id`（与现有 `resource_type|az|day|hour` 并列）

### 2.2 响应（沿用 inline items）

```yaml
# 现有 MeteringUsageRecord 扩展建议字段（可选）
MeteringUsageRecord:
  properties:
    model_id:              { type: string, format: uuid, nullable: true }
    inference_service_id:  { type: string, format: uuid, nullable: true }
```

RBAC：沿用 Metering 读 scope（认证上下文 / `scope:metering:read` — **评审与 usage-report 对齐**）。

---

## 3. Services 便捷汇总（草案 · 3a 可选）

### ModelUsageStatsSummary

```yaml
ModelUsageStatsSummary:
  type: object
  required: [model_id, window_start, window_end]
  properties:
    model_id:           { type: string, format: uuid }
    window_start:       { type: string, format: date-time }
    window_end:         { type: string, format: date-time }
    inference_request_count: { type: integer, minimum: 0 }
    token_input_total:  { type: integer, minimum: 0 }
    token_output_total: { type: integer, minimum: 0 }
    avg_latency_ms:     { type: number, nullable: true }
    p95_latency_ms:     { type: number, nullable: true }
    by_inference_service:
      type: array
      items:
        type: object
        required: [inference_service_id, request_count]
        properties:
          inference_service_id: { type: string, format: uuid }
          request_count:        { type: integer }
          token_input_total:    { type: integer, nullable: true }
          token_output_total:   { type: integer, nullable: true }
```

### 4.1 `GET /api/v1/svc/models/{model_id}/usage-stats`（Services · 可选）

- operationId: `getModelUsageStatsSummary`
- Query: `start_time`、`end_time`（required，与 metering 一致）
- 成功：`200 + ModelUsageStatsSummary`
- 说明：Services **聚合** Core metering + 内部 latency 指标；若 Core filter 已足够，本 op 可降为 3b 或省略

### 4.2 `GET /api/v1/svc/models/usage-stats`（3b · 列表）

- operationId: `listModelsUsageStatsSummaries`
- Query: `start_time`、`end_time`、`limit`、`cursor`
- 成功：200 + items（模型中心列表页 sparkline）

---

## 4. Console 数据流（草案）

```
模型中心 Dashboard
  ├─ 单模型详情：GET /metering/usage?model_id=&start_time=&end_time=&group_by=day
  │              或 GET /svc/models/{id}/usage-stats（便捷）
  ├─ 模型列表排行：GET /metering/usage?resource_type=token_input&group_by=model_id（Core 扩展）
  └─ 跳转：inference-service 详情、usage-report 全租户视角
```

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 管理员 | API |
|---|---|---|---|
| 租户级 metering | 可用 | 可用 | Core usage（无 model filter） |
| 按 model 筛选 | 可用 | 可用 | usage + model_id（扩展） |
| 模型 summary 卡片 | 可用 | 可用 | Services usage-stats（可选） |
| 导出账单 | 跳转 | 跳转 | billing-export |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
# Core 扩展（合入 v1.yaml 后）
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/metering/usage?start_time=2026-06-01T00:00:00Z&end_time=2026-06-16T23:59:59Z&model_id=$MODEL_ID&group_by=day"

# Services 便捷（若合入 services/v1.yaml）
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/models/$MODEL_ID/usage-stats?start_time=...&end_time=..."
```

---

## 7. 评审检查清单

- [ ] 优先 Core metering filter，避免重复计量存储
- [ ] 与 `usage-report.md` 枚举/ group_by 不冲突
- [ ] Services summary 标注可选/3b
- [ ] 合入后更新 `model-usage-stats.md`

---

## 相关文件

- `docs/console-modules/inference/model-usage-stats.md`
- `docs/console-modules/tenant/usage-report.md`
- `docs/console-modules/inference/model-center.md`
