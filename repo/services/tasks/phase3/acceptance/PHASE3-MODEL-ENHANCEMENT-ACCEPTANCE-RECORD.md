# Phase 3 模型增强 ×3 — 文档验收记录

> **生成日期**：2026-06-17  
> **用途**：模型增强三子模块 **文档层验收**（详文 + PRD/SPEC + OpenAPI 草案）。  
> **约束**：本阶段**不修改** `ANI-main/**`；YAML 合入与 handler 在 ANI-main PR 后回填。  
> **索引**：`docs/console-modules/inference/README.md` · **TASK**：`TASK-SVC-018` 子项

---

## 总表

| 模块 | 详文 | PRD/SPEC | OpenAPI 草案 | 文档验收 | YAML 合入 | Handler |
|---|---|---|---|---|---|---|
| 模型加密 | `model-encryption.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-model-encryption-draft.md` | ✅ | ☐ | ☐ |
| 部署推荐 | `model-recommend-config.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-model-recommend-config-draft.md` | ✅ | ☐ | ☐ |
| 调用统计 | `model-usage-stats.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-model-usage-stats-draft.md` | ✅ | ☐ | ☐ |

**文档验收 ✅ 定义**

- 详文含 Services/Core 分层、创建前置条件、操作矩阵、规划口径接口冻结规则
- 加密：Services binding + Core encryption 只读；不重复 keys CRUD
- 推荐：只读 GET；无 POST apply
- 统计：优先 Core metering `model_id` filter；Services summary 标可选

---

> **整域联评**：`tasks/phase3/PHASE3-JOINT-REVIEW-AGENDA.md`（场次 B：KB + 模型）

---

## 联评议程（域内摘要 · 完整版见整域议程 §4.3）

1. **Metering 扩展**（usage-stats + usage-report 同场）
   - `model_id` / `inference_service_id` query；`group_by` 扩展；resource_type 枚举
2. **推荐 → 部署**
   - `ModelDeploymentRecommendation` 字段 ↔ `CreateInferenceServiceRequest`
3. **加密绑定**
   - model/version 粒度；re-encrypt 202；`ENCRYPTION_KEY_REVOKED` 422
4. **合入批次**
   - Core v1.yaml（metering）与 services/v1.yaml（encryption/recommendations/usage-stats）可分两 PR

---

## 文档层 curl 预览（合入 YAML 后）

```bash
curl .../svc/models/$ID/encryption
curl -X PUT .../svc/models/$ID/encryption -d '{"idempotency_key":"...","encryption_key_id":"..."}'
curl ".../svc/models/$ID/recommendations?profile=chat"
curl ".../metering/usage?start_time=...&end_time=...&model_id=$ID&group_by=day"
curl ".../svc/models/$ID/usage-stats?start_time=...&end_time=..."
```

---

## 签核

- [x] ×3 详文 + PRD/SPEC + 草案（2026-06-17）
- [x] `inference/README.md`、`../openapi-drafts/phase3/openapi-phase3-domain-draft.md` §3 更新
- [ ] Metering 扩展与 usage-report 联评通过 — **整域**：`PHASE3-JOINT-REVIEW-AGENDA.md` §5 X-06
- [ ] 分批合入 v1.yaml / services/v1.yaml
- [ ] Handler 运行时验通（ANI-main）

---

## 相关文件

- `docs/console-modules/openapi-drafts/phase3/openapi-phase3-domain-draft.md`
- `tasks/phase3/acceptance/PHASE3-KB-INTELLIGENCE-ACCEPTANCE-RECORD.md`
- `tasks/execution/SERVICES-TEAM-TASKS.md` — TASK-SVC-018
