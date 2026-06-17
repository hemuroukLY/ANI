# 租户用量报表

## 页面定位

`租户用量报表` 是 `Console / 用量与计量` 下的租户侧聚合报表页面，用于帮助用户查看指定时间范围内的资源用量统计。

当前模块属于 `Core / Metering`，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`，正式路径前缀为 `/api/v1`。

本页是租户报表页，不是平台结算与对账页。

## 文档管理规则

- 本文是 `租户用量报表` 的主维护文档，页面定位、参数口径、字段映射和验收标准统一以本文为准
- `prototypes/ani-services-prototype-console.html` 只保留页面摘要、模块边界和详细材料入口
- `tasks/modules/prd/console/tenant/prd-console-usage-report.md` 与 `tasks/modules/spec/console/tenant/spec-console-usage-report.md` 作为辅助材料保留，不替代本文
- 如本文与辅助材料出现差异，先对照 `ANI-main/repo/api/openapi/v1.yaml`，再统一回写

## Core 层要求

- 租户用量统计属于 `Core / Metering`
- 正式查询路径为 `GET /api/v1/metering/usage`
- 当前页面不把 `/api/v1/metering/token-usage` 写成租户查询接口
- 当前页面不把 `GPU-Hours`、`CPU-Hours`、`Memory-GBHours`、`Storage-GBDays`、`Input Tokens`、`Output Tokens`、`KB Queries` 写成独立 API
- 当前页面不扩写账单、发票、对账与平台运营分析
- 当前权威源未给 `/metering/usage` 声明 `operationId`，正文不得自造
- 错误响应格式统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`

## 页面职责

- 提供时间范围维度的租户用量查询
- 提供资源类型筛选
- 提供分组维度切换
- 为多个报表入口提供统一数据口径
- 作为首页和服务模块的趋势钻取承接页

## 页面结构

- 首屏至少包含 `时间筛选区`、`预设视角区`、`趋势图`、`明细表格`、`边界说明`
- 图表与表格必须共享同一查询上下文
- 预设视角切换必须在同一页内完成，不拆出新的正式报表资源
- 无数据态、无权限态、查询失败态都必须有不同反馈

## 数据来源与分层约束

### 数据来源划分

- `Core` 数据
  - 租户聚合用量查询 `GET /api/v1/metering/usage`
- `Services` 数据
  - 当前页面不承接 `Services` 计量正式资源

### 关键边界

- `/metering/usage` 是租户查询入口
- `/metering/token-usage` 是写入侧上报接口，不是租户查询入口
- 预设视角属于 UI 侧过滤和展示口径，不等于新的 `Core` API
- 账单、发票、结算、对账、平台分析只保留为待补边界

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 时间筛选区 | Core | `start_time`、`end_time` | 报表刷新 |
| 预设视角区 | UI + Core | 基于当前查询结果做视角切换 | 趋势图 / 明细表格 |
| 用量趋势图 | Core | 聚合后的 `items[]` | 明细表格 |
| 用量表格 | Core | 展示资源类型、数量、单位、周期 | 同页分析 |
| 边界说明 | 规划项 | 说明账单与结算不在本页内 | 后续平台模块 |

## 字段级定义

### 返回字段

| 字段 | 来源 | 说明 |
|---|---|---|
| `items[].resource_type` | inline response / `MeteringUsageRecord.resource_type` | 资源类型 |
| `items[].total_quantity` | inline response / `MeteringUsageRecord.total_quantity` | 聚合数量 |
| `items[].unit` | inline response / `MeteringUsageRecord.unit` | 单位 |
| `items[].period` | inline response / `MeteringUsageRecord.period` | 统计周期，可空 |
| `total` | inline response | 结果总数 |
| `dev_profile` | inline response / `CoreDevProfileInfo` | 开发者画像或上下文信息 |

### 查询字段

| 字段 | 来源 | 说明 |
|---|---|---|
| `start_time` | query | 起始时间，必填 |
| `end_time` | query | 结束时间，必填 |
| `resource_type` | query | 资源类型筛选，可选 |
| `group_by` | query | `resource_type / az / day / hour`，可选 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 展示趋势图 + 明细表格 + 时间范围 | 图表与表格共享同一数据源 |
| 无数据态 | 展示“当前时间范围内暂无用量” | 不伪造全 0 图表 |
| 查询失败态 | 图表和表格分别提示失败，并保留重试入口 | 不清空筛选上下文 |
| 无权限态 | 展示无权限提示，不渲染误导性占位结果 | 不伪造空数据 |
| 预设视角切换 | 保持同页切换，不生成新 API 路径 | 属于 UI 交互 |

## 字段口径与单位

| 字段 | 口径建议 | 单位/格式 |
|---|---|---|
| `resource_type` | 以后端返回值为准，不用 UI 预设名替代正式枚举 | 字符串 |
| `total_quantity` | 聚合后的数值结果 | number |
| `unit` | 以后端返回单位为准 | 字符串 |
| `period` | 当前聚合周期或时间桶 | 字符串，可空 |
| `start_time` / `end_time` | 查询时间边界 | date-time |

## 已冻结的 Core 能力

- `GET /api/v1/metering/usage`

## 视角说明

- `GPU-Hours`
- `CPU-Hours`
- `Memory-GBHours`
- `Storage-GBDays`
- `Input Tokens`
- `Output Tokens`
- `KB Queries`

以上项是同一租户用量报表页的预设视角，不代表新增独立 `Core` API，也不等于 `/metering/usage` 当前已声明的正式 `resource_type` 枚举。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录与租户上下文 | 已认证 | `401` / `403` |
| `start_time` / `end_time` | 合法时间范围且 `start_time < end_time` | `400 BAD_REQUEST` |
| 报表查看权限 | 已授权 | `403 FORBIDDEN` |

## 操作可用性矩阵

| 操作 | 未登录 | 只读用户 | 租户管理员 | 说明 |
|---|---|---|---|---|
| 查看用量报表 | ❌ | ✅ | ✅ | `GET /api/v1/metering/usage` |
| 切换预设视角 | ❌ | ✅ | ✅ | UI 侧过滤，不新增 API |
| 导出/账单（待补） | ❌ | ❌ | 待补 | 不在当前冻结范围 |

## 接口冻结规则

### `GET /api/v1/metering/usage`

- `operationId`: 当前权威源未声明
- `summary`: `查询租户用量统计`
- `tags`: `["Metering"]`
- `query.required`: `start_time`、`end_time`
- `query.optional`: `resource_type`、`group_by`
- `group_by enum`: `resource_type`、`az`、`day`、`hour`
- `success`: `200 + {items,total,dev_profile}`
- `errors`: `400`、`401`、`403`

## 使用规则

- `start_time` 与 `end_time` 为必填
- `group_by` 只允许 `resource_type`、`az`、`day`、`hour`
- 页面不得把写入侧 `/metering/token-usage` 暴露为租户查询接口
- 推理用量统计若展示，应作为本页过滤视角，而非独立数据源
- 页面不得把 UI 预设视角名称直接硬编码成正式资源类型枚举

## 待补能力边界

- 账单金额
- 发票
- 对账单
- 平台级多租户运营分析
- 租户可见的 token 明细查询
- 按 `model_id` / `inference_service_id` 筛选 — Phase 3 规划见 `inference/model-usage-stats.md`、`../openapi-drafts/phase3/openapi-phase3-model-usage-stats-draft.md`（合入前非冻结）

## 响应示例

### 用量查询成功

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
      "total_quantity": 120000,
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

## 错误示例

### 查询参数错误

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be earlier than end_time",
  "request_id": "req-usage-400-001"
}
```

### 无权限查看用量

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-usage-403-001"
}
```

## 验收标准

- 正文路径、参数、字段与现有 `ANI-main/repo/api/openapi/v1.yaml` 一致
- 正文明确 `/metering/usage` 成功返回体为内联结构
- 各用量视角未被误写成独立 `Core` API
- HTML 摘要、PRD、SPEC 与本文一致
- 本文可以独立作为 `Console / 租户用量报表` 的主维护材料

## 产品经理补充定义

### 目标用户与权限视角

- 租户管理员：查看全租户时间范围内的用量聚合结果
- 业务成员：查看被授权范围内的聚合结果和预设视角
- 只读用户：查看图表与表格，但不承担结算或运营分析动作

### 默认视图与页面状态

- 首屏默认带入产品预设时间窗口并立即查询 `GET /api/v1/metering/usage`；预设是页面配置，不是新的 `Core` 契约
- 资源类型和 `group_by` 切换必须保持在同一页内完成，不拆出新的报表页
- 无数据时展示“当前时间范围内暂无用量”而不是直接展示全 0 图表
- 当查询失败或无权限时，图表区和表格区需要给出不同的可读反馈

### 核心任务流

1. 用户进入报表页后先确认当前时间范围与默认聚合视角
2. 用户切换 `resource_type`、`group_by` 或预设视角，查看图表和表格同步刷新
3. 用户从首页或服务模块回跳到本页时，沿用已有筛选上下文继续分析

### 跨模块协同

- 与 `平台概览` 协同，用于承接首页中的资源或服务趋势钻取
- 与 `推理服务`、`知识库管理` 协同，仅通过预设视角和筛选跳转连接，不重写业务统计契约
- 账单、发票、对账等未冻结能力只保留边界说明，不在本页出现伪入口

### 产品验收补充

- 用户必须能在一个页面内完成时间范围切换、预设视角切换和分组分析
- 图表与表格必须共享同一查询上下文，不能出现口径漂移
- 无数据态、无权限态、查询失败态必须可区分
- 本页不得把预设视角写成新的独立 `Core` API
