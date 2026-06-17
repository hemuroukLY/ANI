# 租户计费与用量

## 页面定位

`租户计费与用量` 是 BOSS **租户与客户管理** 域下的 **跨租户计量与消费分析** 页：按时间范围展示各租户 GPU/CPU/内存/Token 等资源用量排行、趋势与占比，供平台运营与财务做对账前分析。

本页属于 **Core / Metering** 视角下的 **平台 RBAC** 页面，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`。

Console [`usage-report.md`](../../console-modules/tenant/usage-report.md) 仅查询 **当前 JWT 租户** 的 `GET /api/v1/metering/usage`；本页须 **平台 RBAC + 平台 metering aggregate**，**不得** 逐租户切换 JWT 轮询。各指标专页见 [`../metering/`](../metering/)。

本页是 BOSS 运营分析页，**不是** Console 账单导出页；结算/发票归 Phase 2。

## 文档管理规则

- 本文是 `租户计费与用量` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- [`prd-boss-tenant-usage-billing.md`](../../tasks/modules/prd/boss/tenant/prd-boss-tenant-usage-billing.md) 与 [`spec-boss-tenant-usage-billing.md`](../../tasks/modules/spec/boss/tenant/spec-boss-tenant-usage-billing.md) 为辅助材料，不替代本文
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- 计量归属 **Core / Metering**；租户上下文查询：`GET /api/v1/metering/usage`
- `/metering/usage` **无 operationId** — 正文不得自造
- 写入侧：`POST /api/v1/metering/token-usage`（`reportTokenUsage`）— **不是** 本页查询接口
- BOSS 跨租户 list/aggregate — **TODO-YAML**；不得把单租户 GET 直接写成 BOSS 正式契约
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422` 仅 YAML 已声明 operation 可写冻结（§2.10）；本页查询侧暂无 422
- 禁止自造 metering platform schema / operationId

## 页面职责

- 提供 **全平台** 多指标视角切换（GPU-Hours、CPU-Hours、Memory-GBHours、Token 等）下的租户用量 **排行与趋势**
- 提供时间筛选、租户筛选（平台 RBAC）、`group_by` 切换（平台 API 合入后含 `tenant_id`）
- 支持钻取单租户明细（口径对齐 Console usage-report）
- 为 [`../metering/platform-gpu-hours.md`](../metering/platform-gpu-hours.md) 等专页提供 **多指标聚合入口**
- 为 [`tenant-list.md`](tenant-list.md)、[`tenant-quota-policy.md`](tenant-quota-policy.md) 提供用量上下文深链
- **不写** 账单金额、发票、对账（Phase 2）

## 页面结构

- 首屏至少包含：`时间筛选区`、`指标视角切换`、`租户用量排行表`、`趋势图`、`边界说明`
- 排行表与趋势图 **共享同一查询上下文**
- 预设视角在同一页切换，**不** 拆成独立 REST 资源
- 无数据态、无权限态、平台 API 未就绪态、查询失败态须可区分

```text
租户计费与用量
├── 时间筛选（start_time / end_time）
├── 租户筛选（tenant_id / status — 平台 API 待 YAML）
├── 指标视角切换（UI 预设 · 映射 resource_type）
├── 租户用量排行表
├── 趋势图（group_by=day|hour|tenant_id 待 YAML）
└── 钻取抽屉 → 单租户明细（Console usage-report 口径）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | `GET /api/v1/metering/usage` | 租户上下文 **只读参考** |
| Core | 平台 metering aggregate **TODO-YAML** | BOSS 正式跨租户数据源 |
| Core | `POST /api/v1/metering/token-usage` | 写入侧；**非查询** |
| Services | — | 本页 **不** 承接 Services 计量 |

### 关键边界

- 预设视角（GPU-Hours 等）是 **UI 名**，映射 YAML `MeteringUsageRecord.resource_type` — **不是** 独立 API
- `group_by` 租户 API enum：`resource_type` / `az` / `day` / `hour` — **不含** `tenant_id`；平台 API 须扩展 **TODO-YAML**
- `Storage-GBDays`、`KB Queries` **不在** 当前 `MeteringUsageRecord.resource_type` enum — UI 须标注 **TODO-YAML** 或禁用
- `/metering/token-usage` **不得** 作为查询入口
- 账单金额字段 **待补**；Console `billing-export.md` 不含 BOSS 结算

### 预设视角 ↔ resource_type 映射（YAML 已冻结部分）

| UI 视角 | resource_type | 状态 |
|---|---|---|
| GPU-Hours | `instance_gpu_seconds` | ✅ YAML |
| CPU-Hours | `instance_cpu_seconds` | ✅ YAML |
| Memory-GBHours | `instance_memory_gib_seconds` | ✅ YAML |
| Input Tokens | `token_input` | ✅ YAML |
| Output Tokens | `token_output` | ✅ YAML |
| Token Total | `token_total` | ✅ YAML |
| Storage-GBDays | — | **TODO-YAML** enum |
| KB Queries | — | **TODO-YAML** enum |

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 时间筛选区 | UI + Core | start/end 必填 | 刷新排行/趋势 |
| 指标视角区 | UI + Core | 切换 resource_type | 联动图表 |
| 排行表 | 平台 aggregate **待 YAML** | tenant_id + total_quantity | 钻取抽屉 |
| 趋势图 | aggregate；group_by day/hour | 与排行同参 | — |
| 钻取抽屉 | 单租户 usage 参考 | 对齐 usage-report | quota / metering 专页 |
| 边界说明 | 规划项 | 单租户 API 仅为参考 | billing-export |

## BOSS 与 Console 分工

| 维度 | BOSS 租户计费与用量 | Console usage-report |
|---|---|---|
| 范围 | 多租户排行、对比 | 单租户 |
| API | 平台 aggregate **待 YAML** | `GET /metering/usage` |
| group_by | 含 tenant_id **待 YAML** | resource_type/az/day/hour |
| 账单/发票 | Phase 2 | billing-export 草案 |
| RBAC | 平台 metering 读 | 租户 JWT |
| 专页深链 | ../metering/* | 同页预设视角 |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/metering/usage` | **未声明** | 单租户；query 见下 |
| POST | `/api/v1/metering/token-usage` | `reportTokenUsage` | 写入；**非查询** |

| 参数 | 冻结值 |
|---|---|
| `start_time` / `end_time` | required，date-time |
| `resource_type` | optional；见 `MeteringUsageRecord` enum |
| `group_by` | `resource_type` / `az` / `day` / `hour` |

| `MeteringUsageRecord.resource_type` | `instance_cpu_seconds` / `instance_memory_gib_seconds` / `instance_gpu_seconds` / `token_input` / `token_output` / `token_total` |

| 能力 | 状态 |
|---|---|
| 跨租户用量 list / 排行 | **ADDED-TO-YAML** `/metering/usage/platform` |
| `group_by=tenant_id` | **TODO-YAML** |
| 平台 CSV 导出 | **TODO-YAML** |
| 账单金额 | **待补** Phase 2 |

### 建议 TODO-YAML（非冻结）

| 建议能力 | 建议路径 |
|---|---|
| 平台用量聚合 | `GET /api/v1/metering/usage/platform` 或 `GET /api/v1/metering/tenants` |
| 平台钻取单租户 | query `tenant_id` + **后端 RBAC 校验** |

## 字段级定义

### 查询字段 — 租户 `/metering/usage`（只读参考 · YAML 已声明）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `start_time` | query | ✅ | date-time |
| `end_time` | query | ✅ | date-time |
| `resource_type` | query | ❌ | YAML enum |
| `group_by` | query | ❌ | 四 enum；**无 tenant_id** |

### 查询字段（平台 aggregate 目标 · **TODO-YAML**）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `start_time` / `end_time` | query | ✅ | 同租户 API |
| `resource_type` | query | ❌ | 视角映射 |
| `tenant_id` | query | ❌ | 平台 RBAC 筛选 |
| `group_by` | query | ❌ | 扩展 `tenant_id` |
| `tenant_status` | query | ❌ | 产品筛选 **待 YAML** |

### 返回字段 — `MeteringUsageRecord`（YAML 已冻结）

| 字段 | 说明 |
|---|---|
| `resource_type` | enum |
| `total_quantity` | number (double) |
| `unit` | string |
| `period` | string，nullable |
| `tenant_id` | string，nullable（schema 有字段；租户 API inline 200 未必返回 — 平台 API 目标） |

### 返回字段 — `MeteringUsageResponse`

| 字段 | 说明 |
|---|---|
| `items[]` | MeteringUsageRecord 数组 |
| `total` | integer |
| `dev_profile` | CoreDevProfileInfo；联调标记 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `gpu_hours_display` | total_quantity/3600 when unit=seconds |
| `share_pct` | 租户占平台总量 % |
| `rank` | 排行序号 |
| `mom_change_pct` | 环比 **待 YAML** |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 视角 + 排行 + 趋势 | 共享查询上下文 |
| 无数据态 | 「当前时间范围内暂无用量」 | **不** 伪造全 0 |
| 平台 API 未就绪 | TODO-YAML 横幅 | 不得伪装生产 |
| Storage/KB 视角 | disabled +「enum 待 YAML」 | 禁止伪造查询 |
| 查询失败态 | 分块失败 + 重试 | 保留筛选 |
| 无权限态 | 403 | 不渲染假数据 |
| 钻取无租户权 | 403 / 行 disabled | RBAC |
| dev_profile | 联调角标 | 不进财务报表 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `instance_gpu_seconds` | YAML enum | 展示可换算 GPU-Hours |
| `token_input` / `token_output` | YAML enum | tokens |
| `total_quantity` | 后端聚合 | number |
| `unit` | 以后端为准 | seconds / tokens 等 |
| `period` | 与 group_by 一致 | string |
| GPU-Hours（UI） | quantity/3600 if unit=seconds | 2 位小数 |

## 状态与能力口径

本页 **只读查询**，无资源状态机。

| 能力 | 说明 |
|---|---|
| 查询排行/趋势 | 只读 |
| 导出 CSV | **待 YAML** |
| 生成账单 | Phase 2 |
| 调整配额 | 跳转 tenant-quota-policy |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 metering aggregate 读 RBAC | 已授权（scope **TODO-YAML** 待 Core 定义） | `403 FORBIDDEN` |
| 租户 `GET /metering/usage` 只读参考 | JWT 租户上下文；**非** BOSS 跨租户正式方案 | `403`（operation 已声明，**无** `x-ani-rbac-scope`） |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `start_time` / `end_time` | 合法且 start < end | `400 BAD_REQUEST` |
| 钻取单租户 | 平台 tenant 查看权（**TODO-YAML** 与 aggregate 一并定义） | `403 FORBIDDEN` |
| 导出 | export RBAC | **待 YAML** |

## 操作可用性矩阵

| 操作 | 平台只读 | 平台运营 | 财务只读 | 说明 |
|---|---|---|---|---|
| 查看跨租户排行 | ✅（待 YAML） | ✅ | ✅ | aggregate |
| 切换指标视角 | ✅ | ✅ | ✅ | UI |
| 切换 group_by | ✅ | ✅ | ✅ | 平台 API 待 YAML |
| 钻取单租户 | ✅ | ✅ | ✅ | RBAC |
| 跳转 metering 专页 | ✅ | ✅ | ✅ | ../metering/* |
| 导出 CSV | ❌ | ✅ **待 YAML** | ✅ **待 YAML** | — |
| 生成账单/发票 | ❌ | Phase 2 | Phase 2 | — |

## 删除前置校验与当前契约边界

本页 **无 DELETE** — **N/A**。

## 接口冻结规则

### `GET /api/v1/metering/usage`（租户上下文 · 只读参考 · **非 BOSS 正式契约**）

- `operationId`：**未声明**
- `summary`：`查询租户用量统计`
- `tags`：`["Metering"]`
- `query.required`：`start_time`、`end_time`
- `query.optional`：`resource_type`、`group_by`
- `group_by enum`：`resource_type`、`az`、`day`、`hour`
- `success`：`200 + MeteringUsageResponse`（inline + schema 对齐）
- `errors`：`400`、`401`、`403`
- **BOSS 不得**写「本页正式 API = 上路径」而不标注平台扩展待补

### `POST /api/v1/metering/token-usage`（写入 · **非查询**）

- `operationId`：`reportTokenUsage`
- `x-ani-rbac-scope`：`scope:metering:write`
- `success`：`202 + TokenUsageReport`
- **不得** 暴露为本页查询入口

### 平台 metering aggregate（待补）

<!-- TODO-YAML: GET /api/v1/metering/usage/platform 或等价 -->

- 须支持 `group_by=tenant_id` 与跨租户排行
- 须平台 RBAC
- 合入前不得写入「已冻结」正文

## 使用规则

- **不得** 把 `/metering/token-usage` 当查询接口
- **不得** 把 UI 视角名写入 API `resource_type`
- 排行与趋势 **必须** 同参刷新
- 平台 aggregate 未上线时，**禁止** 生产环境逐租户 JWT 轮询
- Storage-GBDays / KB Queries 视角 **必须** disabled 直至 enum 合入
- `dev_profile` 不进入财务对账

## 待补能力边界

- 平台跨租户 metering API — **ADDED-TO-YAML**
- `resource_type` 扩展 storage / kb — **TODO-YAML**
- 与 [`../metering/`](../metering/) 七专页口径统一 — Phase 1
- 平台 CSV 导出 — **TODO-YAML**
- 账单、对账、发票 — Phase 2

## 响应示例

### 单租户 usage 成功（只读参考 · 非 BOSS 正式契约）

```json
{
  "items": [
    {
      "resource_type": "instance_gpu_seconds",
      "total_quantity": 86400,
      "unit": "seconds",
      "period": "day"
    },
    {
      "resource_type": "token_input",
      "total_quantity": 1250000,
      "unit": "tokens",
      "period": "day"
    }
  ],
  "total": 2,
  "dev_profile": {
    "mode": "real",
    "provider": "kubernetes_rest",
    "real_provider": true,
    "reason": null
  }
}
```

### 平台 aggregate 目标响应（**待 YAML · 非已冻结**）

```json
{
  "items": [
    {
      "tenant_id": "ten-acme",
      "resource_type": "instance_gpu_seconds",
      "total_quantity": 604800,
      "unit": "seconds",
      "period": "week"
    },
    {
      "tenant_id": "ten-beta",
      "resource_type": "instance_gpu_seconds",
      "total_quantity": 302400,
      "unit": "seconds",
      "period": "week"
    }
  ],
  "total": 2
}
```

## 错误示例

### 时间范围非法

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be earlier than end_time",
  "request_id": "req-boss-billing-400-001"
}
```

### 无平台 metering 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-billing-403-001"
}
```

> **注**：适用于 **平台 metering aggregate（TODO-YAML）** 场景。`message` 为示例文案，**非** YAML 已冻结 scope 名。租户 `GET /metering/usage` 的 `403` 为 operation 已声明返回码，**无** scope 字段。

## 相关模块

- [tenant-list.md](tenant-list.md)、[tenant-quota-policy.md](tenant-quota-policy.md)
- [`../metering/platform-gpu-hours.md`](../metering/platform-gpu-hours.md) 等七专页
- Console：[usage-report.md](../../console-modules/tenant/usage-report.md)、[billing-export.md](../../console-modules/tenant/billing-export.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] 区分单租户 `/metering/usage` 与 BOSS 平台 API
- [x] 预设视角映射 YAML enum；Storage/KB 标注 TODO-YAML
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（查询/返回/展示）
- [ ] 平台 metering YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
