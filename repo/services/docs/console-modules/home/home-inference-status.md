# 首页 — 推理服务状态

## 页面定位

`推理服务状态` 是 `平台概览` 第三主题区块的明细页，用于汇总当前租户下推理服务的健康分布（运行中、部署中、失败等）。

本页为 **Services 聚合视角**，不新增推理服务资源契约。

## 文档管理规则

- 本文是「推理服务状态」区块的主维护文档
- 状态枚举必须与 `InferenceService.status` 完全一致
- 一级权威源：`services/v1.yaml` `InferenceService`
- 与 `platform-overview.md` 第三区块字段口径一致

## Services 层要求

- 正式路径：`GET /api/v1/svc/inference-services`
- `operationId`: `listInferenceServices`
- 页面层对 `items[]` 按 `status` 分组计数
- 错误：`401`（YAML 已声明）
- 路径前缀 `/api/v1/svc/*`

### 冻结状态枚举

`pending` / `deploying` / `running` / `stopping` / `stopped` / `failed`

## Core 层要求

- 本页不调用 Core 实例 API 作为推理服务状态来源
- 异步部署进度可关联 `GET /api/v1/tasks/{task_id}`（若列表项暴露 task 关联）

## 页面职责

- 展示各状态服务数量与最近活跃服务数
- 高亮 `failed` > 0
- 跳转推理服务列表（可带 status 筛选）
- 说明统计时间窗口

## 字段级定义

与 `platform-overview.md` §推理服务状态 一致：

| 字段 | 说明 |
|---|---|
| 运行中服务数 | `status=running` |
| 部署中服务数 | `pending` + `deploying` |
| 失败服务数 | `status=failed` |
| 已停止服务数 | `stopped`（含 `stopping` 进行中可选单独展示） |
| 最近活跃服务数 | **待确认** — 依赖 metering 或 access log，当前无独立冻结 API |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 推理服务读权限 | Services RBAC | **当前 list YAML 未声明 `403`** |

本页无 POST/PUT；不涉及 `idempotency_key`。

## 操作可用性矩阵

| 操作 | 只读用户 | 服务管理员 |
|---|---|---|
| 查看状态摘要 | ✅ | ✅ |
| 刷新摘要 | ✅ | ✅ |
| 跳转推理服务列表 | ✅ | ✅ |
| 在本页部署 / 删除 | ❌ | ❌ |

## 接口冻结规则

### `GET /api/v1/svc/inference-services`

- 成功：`200` + items 数组（`InferenceService`）
- 错误：`401`
- **当前 YAML 未声明** `403`、`404`

### `GET /api/v1/tasks/{task_id}`（可选关联）

- 成功：`200 + AsyncTask`
- 错误：`404`
- 用于部署中服务的进度展示

## 待补边界

- 「最近活跃服务数」口径 — 依赖 `GET /api/v1/metering/usage` 或 access log，**待产品确认**
- list 接口 `status` query 筛选 — **当前 YAML 未在 list 声明** status filter
- 独立首页推理摘要 API — **TODO-YAML**（可选，非必须）

## 与相关模块的关系

- 明细 CRUD：`inference-service.md`
- 首页嵌入：`platform-overview.md`
- 可观测性：`inference-observability.md`

## 验收标准

- [ ] 状态名与 YAML enum 一致
- [ ] 「最近活跃」未冻结处标注待确认
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
