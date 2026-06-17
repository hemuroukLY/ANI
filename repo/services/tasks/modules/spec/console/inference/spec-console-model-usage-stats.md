# SPEC: Console model-usage-stats

> Technical specification — Phase 3 规划详化  
> Source: `tasks/modules/prd/console/inference/prd-console-model-usage-stats.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-usage-stats-draft.md`  
> Revised: 2026-06-17

## 1. Summary

模型维度用量（规划）；**优先 Core metering 扩展**；Services summary 可选。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- Core 扩展: `ANI-main/repo/api/openapi/v1.yaml` — `GET /metering/usage`
- Services 可选: `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Planned Core Query 扩展

| 参数 | 说明 |
|---|---|
| `model_id` | uuid，筛选单模型 |
| `inference_service_id` | uuid |
| `group_by` | 扩展 `model_id`、`inference_service_id`（3b） |

### 2.3 Planned Services Paths（可选）

| Method | Path | operationId | 成功 |
|---|---|---|---|
| GET | `/api/v1/svc/models/{model_id}/usage-stats` | `getModelUsageStatsSummary` | 200 |
| GET | `/api/v1/svc/models/usage-stats` | `listModelsUsageStatsSummaries` | 200（3b） |

### 2.4 Planned Schemas

- `ModelUsageStatsSummary`

## 3. Page Scope

- 模型中心 Dashboard；跳转 usage-report、inference-service
- **Non-Goals**：token-usage 查询；billing-export

## 4. 边界

与 `tenant/usage-report.md` 联评 group_by / resource_type 枚举。

## 5. Handler 验收（合入 YAML 后）

```bash
curl ".../metering/usage?start_time=...&end_time=...&model_id=$ID&group_by=day"
curl ".../svc/models/$ID/usage-stats?start_time=...&end_time=..."
```

## 6. 主维护源

- `docs/console-modules/inference/model-usage-stats.md`
- 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-usage-stats-draft.md`

OpenAPI **未合入** ≠ handler 已实现。
