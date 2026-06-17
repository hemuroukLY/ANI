# 推理服务 — 限流与访问策略

## 页面定位

单个推理服务的限流 RPM 与 allowed_scopes 配置页。

父级：`inference-service.md`。

## 文档管理规则

- 本文是推理策略子模块主维护源
- 一级权威源：`services/v1.yaml`

## Services 层要求

<!-- ADDED-TO-YAML: PUT .../policies (Phase 2) -->

| 方法 | 路径 | operationId |
|---|---|---|
| PUT | `/api/v1/svc/inference-services/{service_id}/policies` | `updateInferenceServicePolicies` |

请求：`InferenceServicePolicyUpdateRequest`（必填 `idempotency_key`；可选 `rate_limit_rpm`、`allowed_scopes`）。

响应：`200 + InferenceService`。

错误：`400`、`401`、`403`、`404`、`422`。

## 页面职责

- 服务详情「策略」Tab：编辑 RPM、scope 列表
- 保存后刷新服务详情中的策略字段
- 不提供策略历史审计（待补）

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 服务存在 | 有效 `service_id` | `404 NOT_FOUND` |
| 写权限 | 服务管理员 | `403 FORBIDDEN` |
| PUT 请求体 | `idempotency_key` | `400 BAD_REQUEST` |
| 服务可更新策略 | 建议 `status=running` | `422 PRECONDITION_FAILED`（具体 `code` 待 Services 冻结） |

## 操作可用性矩阵

| 操作 | 只读 | 服务管理员 |
|---|---|---|
| 查看策略 | ✅ | ✅ |
| 更新策略 | ❌ | ✅ |

### 按服务状态（产品建议）

| 操作 | running | stopped | failed |
|---|---|---|---|
| 更新策略 | ✅ | ❌ → 422 | ❌ → 422 |

## 接口冻结规则

### `PUT /api/v1/svc/inference-services/{service_id}/policies`

- 成功：`200 + InferenceService`
- 错误：`400`、`401`、`403`、`404`、`422`
- requestBody.required：`idempotency_key`

## 待补边界

- `422` code 举例 — 待 Services YAML description
- 策略变更审计 log — **TODO-YAML**
- allowed_scopes 枚举权威清单 — 待 RBAC 注册表

## 验收标准

- [ ] 路径与 services/v1.yaml 一致
- [ ] PUT 必填 `idempotency_key`
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
