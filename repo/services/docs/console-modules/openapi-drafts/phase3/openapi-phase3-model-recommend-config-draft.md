# Phase 3 — 模型部署推荐配置 OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`inference/model-recommend-config.md`  
> **TASK**：`TASK-SVC-018` 子项（模型增强）  
> **原则**：只读推荐；「一键部署」为 Console 预填 `CreateInferenceServiceRequest`，不新增 apply API。

---

## 1. 与推理服务创建 / GPU 库存的边界

| 维度 | 部署推荐（本草案） | 推理服务 |
|---|---|---|
| 路径 | `GET .../models/{id}/recommendations` | `POST .../inference-services` |
| 职责 | 计算建议 gpu/replicas/量化 | 实际部署 |
| 写操作 | **无** | 202 deploy |
| 输入 | model 规格 + workload profile | 用户确认后的完整 config |

推荐结果字段与 `InferenceService` 创建字段对齐：`gpu_type`、`gpu_count_per_pod`、`replicas`、`max_concurrency`、量化相关扩展字段（3b）。

GPU 配额校验 — Services 实现；推荐可标注 `quota_sufficient: false`。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [Models]
x-ani-rbac-scope:
  read: scope:models:read
```

只读 API；部署写操作走 inference-services。

---

## 3. Schemas（草案）

### ModelDeploymentRecommendation

```yaml
ModelDeploymentRecommendation:
  type: object
  required: [model_id, profile, recommendations]
  properties:
    model_id:    { type: string, format: uuid }
    profile:     { type: string, enum: [chat, batch, embedding], description: 查询参数回显 }
    generated_at: { type: string, format: date-time }
    recommendations:
      type: array
      minItems: 1
      items:
        type: object
        required: [tier, gpu_type, gpu_count_per_pod, replicas]
        properties:
          tier:              { type: string, enum: [cost_optimized, balanced, performance] }
          gpu_type:          { type: string }
          gpu_count_per_pod: { type: integer, minimum: 1 }
          replicas:          { type: integer, minimum: 1 }
          max_concurrency:   { type: integer, minimum: 1, nullable: true }
          quantization:      { type: string, enum: [none, int8, int4, fp8], nullable: true }
          estimated_latency_p50_ms: { type: number, nullable: true }
          quota_sufficient:  { type: boolean }
    explain:
      type: array
      items: { type: string }
      description: 人类可读推荐理由（非结构化）
    algorithm_version: { type: string, nullable: true }
```

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/svc/models/{model_id}/recommendations`

- operationId: `getModelDeploymentRecommendations`
- Query:
  - `profile`（required）：`chat` | `batch` | `embedding`
  - `target_latency_ms`（optional）：性能档筛选 hint
  - `az`（optional）：可用区偏好
- 成功：`200 + ModelDeploymentRecommendation`
- 错误：
  - `401`、`403`
  - `404`：model 不存在
  - `422`：model 非 ready — `MODEL_NOT_READY`；无可用 GPU 类型 — `NO_SUITABLE_GPU_PROFILE`

### 4.2 Console「一键应用」（非 API）

- 将选中 `tier` 字段映射为 `POST /api/v1/svc/inference-services` body 预填值
- 用户仍需确认并提交 deploy（含 `idempotency_key`）

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 模型/推理管理员 | 说明 |
|---|---|---|---|
| 查看推荐 | 可用 | 可用 | GET |
| 切换 profile | 可用 | 可用 | query profile |
| 一键预填部署 | 不可用 | 可用 | UI → inference-service |
| 查看 explain | 可用 | 可用 | recommendations.explain |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/models/$MODEL_ID/recommendations?profile=chat"

curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/models/$MODEL_ID/recommendations?profile=batch&target_latency_ms=500"
```

---

## 7. 评审检查清单

- [ ] 只读，无 POST apply 路径
- [ ] 推荐字段可映射 inference-service create
- [ ] model 非 ready 返回 422
- [ ] 合入后更新 `model-recommend-config.md`

---

## 相关文件

- `docs/console-modules/inference/model-recommend-config.md`
- `docs/console-modules/inference/inference-service.md`
- `docs/console-modules/inference/model-center.md`
