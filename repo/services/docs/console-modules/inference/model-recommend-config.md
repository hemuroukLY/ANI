# 模型 — 部署推荐配置

## 页面定位

根据模型规格与租户配额，**推荐推理部署参数**（GPU 型号、副本数、量化）的只读页；「一键部署」为 Console 预填 `inference-service` 表单，无独立 apply API。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-model-recommend-config-draft.md`。

## 文档管理规则

- 本文是 **模型 — 部署推荐配置** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-recommend-config-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- PRD/SPEC：`tasks/modules/prd/console/inference/prd-console-model-recommend-config.md`、`tasks/modules/spec/console/inference/spec-console-model-recommend-config.md`
- TASK：`TASK-SVC-018` 子项（模型增强）

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/models/{model_id}/recommendations` | `getModelDeploymentRecommendations` | 3a |

Query（草案）：`profile`（chat|batch|embedding，required）、`target_latency_ms`、`az`。

Schema（草案）：`ModelDeploymentRecommendation`。

RBAC（草案）：`scope:models:read`（只读）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 读权限 | models:read | `403` |
| model 存在 | 有效 id | `404` |
| profile 合法 | 枚举值 | `400` |
| model 就绪 | `status=ready` | `422`（建议语义：`MODEL_NOT_READY`；**待 YAML 声明后定稿**） |
| 无合适 GPU | 配额/库存 | `422`（建议语义：`NO_SUITABLE_GPU_PROFILE`；**待 YAML 声明后定稿**） |

## 页面职责

- 按 workload profile 展示多档推荐（cost_optimized / balanced / performance）
- 展示 `explain[]` 推荐理由、`quota_sufficient`
- 「一键应用」→ 跳转推理服务部署并预填 gpu/replicas/quantization
- 不新增 POST apply 路径

## 操作可用性矩阵

| 操作 | 只读用户 | 模型/推理管理员 | 说明 |
|---|---|---|---|
| 查看推荐 | 可用 | 可用 | GET |
| 切换 profile | 可用 | 可用 | query |
| 预填部署 | 不可用 | 可用 | UI → inference-service |
| 实际 deploy | 不可用 | 可用 | POST inference-services |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-model-recommend-config-draft.md` §4。

### `GET .../models/{model_id}/recommendations`（规划）

- 成功：`200 + ModelDeploymentRecommendation`
- 错误：`400`、`404`、`422`

## 待补边界

- `algorithm_version` 与 A/B 推荐 — BOSS 侧
- 量化字段与 inference create body 映射表 — 3b ADR
- GPU 库存实时性 — 与 Core gpu-inventory 只读引用

## 相关模块

- `inference/model-center.md`
- `inference/inference-service.md`

## 验收标准

- [ ] 只读 GET，无 apply API
- [ ] 字段可映射 InferenceService create
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
