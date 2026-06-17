# SPEC: Console 租户用量报表

> Technical specification derived from: `tasks/modules/prd/console/tenant/prd-console-usage-report.md`
> Source of truth: `ANI-main/repo/api/openapi/v1.yaml`
> Scope: `Console / 用量与计量 / 租户用量报表`

## 1. Summary

本 SPEC 定义 `Console / 用量与计量 / 租户用量报表` 的技术边界、查询参数、返回体形状、字段映射和页面约束。该页面直接承接 `Core / Metering` 的租户聚合用量查询能力，不扩写为账单、结算或平台运营分析模块。

## 2. Source of Truth

### 2.1 Primary Authority

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Frozen Path Group

- `/metering/usage`

### 2.3 Related Schemas

- `MeteringUsageRecord`
- `MeteringUsageResponse`
- `CoreDevProfileInfo`
- `ErrorResponse`
- `ReportTokenUsageRequest`
- `TokenUsageReport`

说明：

- `/metering/usage` 的 `200` 响应在权威源中以内联结构描述
- `MeteringUsageRecord` 与 `MeteringUsageResponse` 可作为字段映射参考，但正文不得误写为该路径直接 `$ref`
- `/metering/token-usage` 是写入侧上报接口，不属于租户查询能力

## 3. API Design

### 3.1 Endpoint Matrix

| Method | Path | operationId | summary | Success | Error Responses |
|---|---|---|---|---|---|
| GET | `/api/v1/metering/usage` | 当前权威源未声明 | 查询租户用量统计 | `200 + {items,total,dev_profile}` | `400`、`401`、`403` |

### 3.2 Query Parameters

| 参数 | 必填 | 类型 | 说明 |
|---|---|---|---|
| `start_time` | 是 | `date-time` | 起始时间 |
| `end_time` | 是 | `date-time` | 结束时间 |
| `resource_type` | 否 | `string` | 资源类型筛选 |
| `group_by` | 否 | `string` | `resource_type` / `az` / `day` / `hour` |

### 3.3 Response Shape

| 字段 | 类型 | 说明 |
|---|---|---|
| `items` | `array<object>` | 聚合用量结果数组 |
| `items[].resource_type` | `string` | 资源类型 |
| `items[].total_quantity` | `number` | 聚合数量 |
| `items[].unit` | `string` | 计量单位 |
| `items[].period` | `string` | 统计周期，可空 |
| `total` | `integer` | 结果总数 |
| `dev_profile` | `CoreDevProfileInfo` | 开发者画像或上下文信息 |

## 4. Data Model

### 4.1 MeteringUsageRecord Mapping

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `resource_type` | `MeteringUsageRecord.resource_type` | 当前组件 schema 允许的枚举包括 `instance_cpu_seconds`、`instance_memory_gib_seconds`、`instance_gpu_seconds`、`token_input`、`token_output`、`token_total` |
| `total_quantity` | `MeteringUsageRecord.total_quantity` | 聚合数量 |
| `unit` | `MeteringUsageRecord.unit` | 单位 |
| `period` | `MeteringUsageRecord.period` | 聚合周期，可空 |
| `tenant_id` | `MeteringUsageRecord.tenant_id` | 当前页面默认不要求显式展示 |

### 4.2 View Rules

- `GPU-Hours`、`CPU-Hours`、`Memory-GBHours`、`Storage-GBDays`、`Input Tokens`、`Output Tokens`、`KB Queries` 都是页面预设视角
- 这些视角属于 UI 侧筛选和命名，不等于新增正式资源类型枚举
- 页面不得把预设视角直接等同为 `/metering/usage` 当前已声明的 `resource_type` 枚举

### 4.3 Non-Frozen Capabilities

- 账单金额
- 发票
- 对账单
- 平台租户对比分析
- 用户直接调用 `/metering/token-usage`

## 5. Page Design Constraints

### 5.1 Query Page

- 页面默认带入时间范围并触发同一条查询链路
- 图表与表格共享同一查询上下文
- 资源视角切换只影响筛选与展示，不形成新 API

### 5.2 Drilldown Relationship

- `推理用量统计` 只能作为本页的推理侧过滤入口
- 首页或服务模块跳转到报表页时，应保留筛选上下文
- 本页不承担账单、结算、发票或对账入口

## 6. Rules

- 租户用量报表归属 `Core / Metering`
- 正式查询路径必须使用 `/api/v1/metering/usage`
- 页面不把 `/api/v1/metering/token-usage` 写成租户查询接口
- 页面不把预设视角写成新的独立 API 或新的正式资源类型枚举
- 当前权威源未给 `/metering/usage` 声明 `operationId`，正文不得自造

## 7. Error Handling

### 7.1 Unified Error Shape

```json
{"code":"UPPER_SNAKE","message":"...","request_id":"..."}
```

### 7.2 Expected Errors

- `GET /api/v1/metering/usage`: `400 + ErrorResponse`、`401 + ErrorResponse`、`403 + ErrorResponse`

## 8. Acceptance

- 文档中所有正式路径均位于 `/api/v1/metering/usage`
- 文档中明确 `/metering/token-usage` 是写入侧上报接口
- 文档中明确预设视角属于同一报表页，不是独立 API
- 文档中不再出现 `[Assumption]` 或把内联响应误写成直接 `$ref`

## 9. Backfill Dependencies

- 如未来扩展账单、发票、对账，需先冻结正式路径与 schema
- 如未来扩展更细的 AI 服务计量视角，需先明确是否仍复用 `/metering/usage`
- 如未来需要租户查询 token 明细，需先明确新查询接口与写入接口边界
