# 推理与模型模块索引

Console **模型中心 / 推理服务**域子模块主维护源。主资源 `model-center.md`、`inference-service.md` 已冻结 CRUD；**模型增强 ×3** 为 Phase 3 **详化草案**（2026-06-17）。

## 模块清单

| 模块 | 详文 | OpenAPI | 说明 |
|---|---|---|---|
| 模型中心 | [model-center.md](./model-center.md) | **已冻结** | Models CRUD + import |
| 推理服务 | [inference-service.md](./inference-service.md) | **已冻结** | InferenceServices deploy |
| 模型加密 | [model-encryption.md](./model-encryption.md) | **详化草案** | [openapi-phase3-model-encryption-draft.md](../openapi-drafts/phase3/openapi-phase3-model-encryption-draft.md) |
| 部署推荐 | [model-recommend-config.md](./model-recommend-config.md) | **详化草案** | [openapi-phase3-model-recommend-config-draft.md](../openapi-drafts/phase3/openapi-phase3-model-recommend-config-draft.md) |
| 调用统计 | [model-usage-stats.md](./model-usage-stats.md) | **详化草案** | [openapi-phase3-model-usage-stats-draft.md](../openapi-drafts/phase3/openapi-phase3-model-usage-stats-draft.md) |

## 维护规则

- Services 路径前缀：`/api/v1/svc/models*`、`/api/v1/svc/inference-services*`
- 模型增强 ×3 合入前详文「接口冻结规则」保持**规划**标记
- 用量统计 **优先 Core** `GET /metering/usage` 扩展；Services summary 标注可选
- 加密绑定引用 Core `/encryption/keys` — 见 `security/crypto-sm.md`
- TASK：主 CRUD → SVC-005/006；增强 ×3 → `TASK-SVC-018` 子项

## 评审建议顺序

1. model-usage-stats（与 `usage-report.md` 联评 metering 扩展）
2. model-recommend-config（与 inference-service create 字段映射）
3. model-encryption（与 crypto-sm 职责切分）

## 相关

- Backlog：`../governance/console-undefined-features-backlog.md` P2-15～P2-17
- 整域索引：`../openapi-drafts/phase3/openapi-phase3-domain-draft.md` §3
- 租户用量：`../tenant/usage-report.md`
