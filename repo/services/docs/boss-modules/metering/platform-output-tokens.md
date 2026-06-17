# 平台 Output Tokens

## 页面定位

`平台 Output Tokens` 是 BOSS **平台计量与结算** 域下的 **跨租户输出 Token 消耗** 专页：按时间范围展示各租户 output token 排行、趋势与占比，供平台运营、财务做推理成本分析与异常租户识别。

本页属于 **Core / Metering** 视角下的 **平台 RBAC** 页面，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`。

Console [`usage-report.md`](../../console-modules/tenant/usage-report.md) 的「Output Tokens」仅为 **当前 JWT 租户** 的 UI 预设视角；本页 **不得** 逐租户切换 JWT 轮询冒充平台大盘。写入侧 `POST /api/v1/metering/token-usage`（`reportTokenUsage`）由 Services 上报 input+output，**不是** 本页查询接口。

## 文档管理规则

- 本文是 `平台 Output Tokens` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-platform-output-tokens.md`](../../tasks/modules/prd/boss/metering/prd-boss-platform-output-tokens.md) 与 [`spec-boss-platform-output-tokens.md`](../../tasks/modules/spec/boss/metering/spec-boss-platform-output-tokens.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- 计量归属 **Core / Metering**；正式租户查询路径为 `GET /api/v1/metering/usage`
- 本页 UI 名 **Output Tokens** 对应 YAML 冻结枚举 `resource_type=token_output`，**不是** 独立 API
- 当前权威源 **未** 给 `/metering/usage` 声明 `operationId`，正文不得自造
- `POST /api/v1/metering/token-usage`（`reportTokenUsage`）为 **写入侧**，`202`，`scope:metering:write`，须 `idempotency_key`；**不是** 本页查询接口
- `token_total` 在 enum 中存在；本页 **锁定** `token_output`，不与 input/total 混页
- BOSS 跨租户 list/aggregate — **TODO-YAML**；不得把单租户 GET 直接写成 BOSS 正式契约
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422` 仅 YAML 已声明 operation 可写冻结（§2.10）；本页查询侧暂无 422

## 页面职责

- 提供 **全平台** 时间范围内的 output token（`token_output`）消耗排行与趋势
- 提供按 day/hour 时间桶的聚合视图（平台 API 合入后）
- 支持从排行行钻取到单租户用量、推理审计
- 提供与 Input Tokens 联合查看入口（跳转，非混页）
- 为 Phase 2 账单/导出预留边界，当前不写结算金额

## 页面结构

- 首屏至少包含：`时间筛选区`、`平台 KPI 区`、`Output Token 趋势图`、`租户排行表`、`边界说明`
- 趋势图与排行表 **共享同一查询上下文**（start/end、group_by）
- 本页锁定 `resource_type=token_output`，不与 `token_input` / `token_total` 混页
- 无数据态、无权限态、平台 API 未就绪态、查询失败态须可区分

```text
平台 Output Tokens
├── 时间筛选（start_time / end_time）
├── 可选筛选（租户状态 — 平台 API 待 YAML）
├── KPI：平台 output tokens 总量、环比、Top N 租户占比
├── 趋势图（group_by=day|hour）
├── 租户排行表（tenant_id / quantity / unit / 占比）
└── 行内钻取 → 租户用量 / inference 审计 / Input Tokens 联合入口
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | `GET /api/v1/metering/usage` | 租户上下文 **只读参考**；`resource_type=token_output` |
| Core | `POST /api/v1/metering/token-usage` | **写入侧**；一次上报含 input+output；**非查询** |
| Core | 平台 metering aggregate **TODO-YAML** | BOSS 正式 list/排行数据源 |
| Services | — | 本页 **不** 承接 Services 推理 list 作为 metering 契约 |

### 关键边界

- 「Output Tokens」是 **UI 展示名**；正式枚举为 `token_output`
- `token_total` 在 enum 中存在；本页 **锁定** `token_output`
- 单租户 `/metering/usage` **不能** 直接作为 BOSS 跨租户正式契约
- `POST /metering/token-usage` 为写入侧；查询侧按 `resource_type` 拆分 input/output
- 账单金额、发票、对账不在本页（Phase 2）

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 时间筛选区 | UI + Core | `start_time`、`end_time` | 刷新 KPI/图表/表格 |
| 平台 KPI | 平台 aggregate **待 YAML** | 全平台 output token 汇总 | — |
| Output Token 趋势图 | 平台 aggregate **待 YAML** | `group_by=day\|hour` | 排行表联动 |
| 租户排行表 | 平台 aggregate **待 YAML** | `tenant_id` + `total_quantity` | 租户钻取 |
| 边界说明 | 规划项 | token-usage POST 为写入侧 | platform-input-tokens |

## BOSS 与 Console 分工

| 维度 | BOSS 平台 Output Tokens | Console Output Tokens 视角 |
|---|---|---|
| 范围 | 全平台多租户排行 | 单租户 |
| 查询 API | 平台 aggregate **待 YAML** | `GET /metering/usage` + UI filter |
| 写入 API | 不涉及 | 不涉及 |
| RBAC | 平台 metering 读 | 租户 JWT |
| 钻取 | inference 审计 / 租户 billing | 同页图表+表格 |
| 账单 | Phase 2 | 不在 usage-report |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/metering/usage` | — | 单租户；`resource_type=token_output` → `200` |
| POST | `/api/v1/metering/token-usage` | `reportTokenUsage` | 写入侧；`202`；**非查询** |

| 字段 | 冻结值 |
|---|---|
| `MeteringUsageRecord.resource_type` | `token_output` |
| `TokenUsageReport.state` | `accepted` / `duplicate` |
| `group_by`（租户 API） | `resource_type` / `az` / `day` / `hour` |

| 能力 | 状态 |
|---|---|
| 跨租户 output token 排行 | **TODO-YAML** |
| `group_by=tenant_id` | **TODO-YAML** |
| 按 model 筛选 | Phase 3 |
| 平台 CSV 导出 | **TODO-YAML** |

## 字段级定义

### 查询字段（平台 API 目标 · 待 YAML 合入后对齐）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `start_time` | query | ✅ | date-time |
| `end_time` | query | ✅ | date-time |
| `resource_type` | query | 固定 | `token_output` |
| `tenant_id` | query | 可选 | 平台 RBAC 下筛选；**不得** 越权 |
| `group_by` | query | 可选 | 平台扩展 `tenant_id` / `day` / `hour` |
| `tenant_status` | query | 可选 | 产品筛选；**待 YAML** |

### 返回字段

| 字段 | 来源 | 说明 |
|---|---|---|
| `items[].tenant_id` | 平台 aggregate **待 YAML** | 租户标识 |
| `items[].resource_type` | `MeteringUsageRecord` | 固定 `token_output` |
| `items[].total_quantity` | `MeteringUsageRecord` | 聚合 output token 数 |
| `items[].unit` | `MeteringUsageRecord` | 通常 `tokens` |
| `items[].period` | `MeteringUsageRecord` | 时间桶，可空 |
| `total` | list response | 条数 |
| `dev_profile` | `CoreDevProfileInfo` | 联调标记；非主展示 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `output_tokens_display` | 直接展示 `total_quantity`（当 unit=tokens） |
| `share_pct` | 租户占平台总量百分比 |
| `rank` | 排行序号 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | KPI + 趋势图 + 排行表 | 共享同一查询上下文 |
| 无数据态 | 「当前时间范围内暂无 Output Token 用量」 | **不** 伪造全 0 图表 |
| 平台 API 未就绪 | 说明「跨租户 aggregate 待 Core 合入」 | 不得伪装为生产已上线 |
| 查询失败态 | KPI/图表/表格分别失败提示 + 重试 | 保留筛选条件 |
| 无权限态 | 403 提示，不渲染假数据 | 平台 RBAC 拒绝 |
| 钻取无租户权 | 行内禁用或 403 提示 | 不得越权查看 tenant 明细 |
| 大数值展示 | 可用 K/M 缩写；tooltip 展示精确值 | 不改变 API 原值 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `resource_type` | 必须为 `token_output` | YAML enum |
| `total_quantity` | 后端聚合原值 | number (double) 或 integer |
| `unit` | 以后端为准；常见 `tokens` | string |
| Output Tokens（UI） | 与 `total_quantity` 一致 | 整数或保留 0 位小数 |
| `period` | 与 `group_by` 时间桶一致 | string，可空 |
| `start_time` / `end_time` | ISO 8601 | date-time |

## 状态与能力口径

本页为 **只读查询页**，无资源状态机。`TokenUsageReport.state`（`accepted`/`duplicate`）属于写入侧，不在本页展示。

| 能力 | 说明 |
|---|---|
| 查询 | 只读 |
| 上报 token | **不在本页**；由 Services 调用 POST token-usage |
| 导出 | **待 YAML** |
| 查看单条 report 明细 | **待 YAML** |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 metering aggregate 读 RBAC | 已授权（scope **TODO-YAML** 待 Core 定义） | `403 FORBIDDEN` |
| 租户 `GET /metering/usage` 只读参考 | JWT 租户上下文；**非** BOSS 跨租户正式方案 | `403`（operation 已声明，**无** `x-ani-rbac-scope`） |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `start_time` / `end_time` | 合法且 start < end | `400 BAD_REQUEST` |
| 钻取单租户 | 具备该租户平台查看权 | `403` |
| 导出（待 YAML） | export RBAC | **待 YAML** |

## 操作可用性矩阵

| 操作 | 平台只读 | 平台运营 | 财务只读 | 说明 |
|---|---|---|---|---|
| 查看跨租户 output token 排行 | ✅ | ✅ | ✅ | aggregate 待 YAML |
| 切换 group_by | ✅ | ✅ | ✅ | 平台 API 待 YAML |
| 钻取租户 | ✅ | ✅ | ✅ | 须 RBAC |
| 查看单条 report 明细 | ❌ | **待 YAML** | ❌ | 非本页 |
| 导出 CSV | ❌ | ✅ **待 YAML** | ✅ **待 YAML** | — |
| 上报 token | ❌ | ❌ | ❌ | Services 写入侧 |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。

## 接口冻结规则

### `GET /api/v1/metering/usage`（租户上下文 · 只读参考 · **非 BOSS 正式契约**）

- `operationId`：权威源 **未声明**
- `summary`：`查询租户用量统计`
- `tags`：`["Metering"]`
- `query.required`：`start_time`、`end_time`
- `query.optional`：`resource_type`、`group_by`
- 本页固定：`resource_type=token_output`
- `group_by enum`：`resource_type`、`az`、`day`、`hour`（**不含** `tenant_id`）
- `success`：`200 + MeteringUsageResponse`
- `errors`：`400`、`401`、`403`

### `POST /api/v1/metering/token-usage`（写入侧 · **非本页查询**）

- `operationId`：`reportTokenUsage`
- `summary`：`上报 Token 用量`
- `tags`：`["Metering"]`
- `x-ani-rbac-scope`：`scope:metering:write`
- `requestBody.required`：`idempotency_key`、`source`、`model`、`input_tokens`、`output_tokens`
- `success`：`202 + TokenUsageReport`
- `errors`：`400`、`401`、`403`
- **禁止** 把此 POST 暴露为本页查询入口

### 平台 aggregate（待补）

<!-- TODO-YAML: GET /api/v1/metering/usage/platform 或等价 -->

- 须支持 `group_by=tenant_id` 与平台 RBAC
- 合入前不得写入「已冻结」正文

## 使用规则

- 页面 **不得** 把 `POST /metering/token-usage` 暴露为查询入口
- **不得** 混淆 `token_output` 与 `token_total`（enum 均存在；本页锁定 output）
- 不得把 UI 名「Output Tokens」写入 API `resource_type`
- 图表与排行表必须同参刷新，避免口径漂移
- 平台 aggregate 未上线时，**禁止** 生产环境用逐租户 JWT 轮询作为正式方案

## 待补能力边界

- 平台跨租户 token aggregate — **ADDED-TO-YAML** P1
- 平台 CSV 导出 — **TODO-YAML**
- 按 model / inference_service 筛选 — Phase 3
- 账单金额 / 发票 — Phase 2
- 单条 token report 明细查询 — **待 YAML**

## 响应示例

### 单租户 usage 查询成功（只读参考 · 非 BOSS 正式契约）

```json
{
  "items": [
    {
      "resource_type": "token_output",
      "total_quantity": 800000,
      "unit": "tokens",
      "period": "day"
    }
  ],
  "total": 1,
  "dev_profile": {
    "mode": "real",
    "provider": "kubernetes_rest",
    "real_provider": true,
    "reason": null
  }
}
```

## 错误示例

### 时间范围非法

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be earlier than end_time",
  "request_id": "req-boss-output-400-001"
}
```

### 无平台 metering 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-output-403-001"
}
```

> **注**：适用于 **平台 metering aggregate（TODO-YAML）** 场景。`message` 为示例文案，**非** YAML 已冻结 `x-ani-rbac-scope` 名。租户 `GET /metering/usage` 的 `403` 为 operation 已声明返回码，**无** scope 字段。

## 相关模块

- [`platform-input-tokens.md`](platform-input-tokens.md)、[`tenant-usage-billing.md`](../tenant/tenant-usage-billing.md)
- [`../audit/platform-inference-audit.md`](../audit/platform-inference-audit.md)
- Console：[usage-report.md](../../console-modules/tenant/usage-report.md)

## 回填验收标准

- [x] 满配章节齐全（对照 `boss-full-depth-checklist.md`）
- [x] 区分 query usage 与 token-usage POST
- [x] `resource_type=token_output` 与 enum 一致
- [x] 含响应示例与错误示例
- [ ] 平台 aggregate YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
