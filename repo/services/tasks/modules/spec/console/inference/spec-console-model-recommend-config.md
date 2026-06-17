# SPEC: Console model-recommend-config

> Technical specification — Phase 3 规划详化  
> Source: `tasks/modules/prd/console/inference/prd-console-model-recommend-config.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-recommend-config-draft.md`  
> Revised: 2026-06-17

## 1. Summary

部署推荐只读 GET（规划）；多档 tier 映射 InferenceService 创建字段。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Planned Paths

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/svc/models/{model_id}/recommendations` | `getModelDeploymentRecommendations` | 200 | `scope:models:read` |

Query: `profile`（required）、`target_latency_ms`、`az`。

### 2.3 Planned Schemas

- `ModelDeploymentRecommendation`（recommendations[]: tier, gpu_type, replicas, quantization, explain）

## 3. Page Scope

- profile 切换、explain、预填 deploy（UI only）
- **Non-Goals**：POST apply

## 4. 创建前置条件

`MODEL_NOT_READY`、`NO_SUITABLE_GPU_PROFILE` — 见主维护源。

## 5. Handler 验收（合入 YAML 后）

```bash
curl ".../svc/models/$ID/recommendations?profile=chat"
```

## 6. 主维护源

- `docs/console-modules/inference/model-recommend-config.md`
- 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-recommend-config-draft.md`

OpenAPI **未合入** ≠ handler 已实现。
