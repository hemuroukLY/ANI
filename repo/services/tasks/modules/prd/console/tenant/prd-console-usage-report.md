# PRD: Console 租户用量报表

## 1. Introduction/Overview

`Console / 用量与计量 / 租户用量报表` 是租户侧聚合报表页面，用于查看指定时间范围内的资源用量统计。该页面直接承接 `Core / Metering` 的正式读能力，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`。

当前已确认的正式能力：

- `GET /api/v1/metering/usage`

当前仍属于待补边界的能力：

- 账单金额
- 发票
- 结算与对账
- 平台级多租户运营分析
- 让租户直接调用写入侧 `/api/v1/metering/token-usage`

## 2. Goals

- 让租户查看指定时间范围内的聚合资源用量
- 让页面支持按资源类型和分组维度切换
- 明确 `Console` 报表页与平台结算、运营分析的边界
- 沉淀可直接用于 `Core / Metering` 对齐的主维护材料

## 3. User Stories

### US-001: 查看租户总体用量

作为租户用户，我希望查看当前时间范围内的总体资源用量，以便理解资源消耗情况。

**Acceptance Criteria**

- [ ] 页面支持设置 `start_time` 和 `end_time`
- [ ] 查询接口对齐 `GET /api/v1/metering/usage`
- [ ] 页面至少展示 `resource_type`、`total_quantity`、`unit`、`period`

### US-002: 按资源类型筛选用量

作为租户用户，我希望按资源类型查看用量，以便快速聚焦不同资源消耗。

**Acceptance Criteria**

- [ ] 页面支持 `resource_type` 筛选
- [ ] 页面支持将 `GPU-Hours`、`CPU-Hours`、`Memory-GBHours`、`Storage-GBDays`、`Input Tokens`、`Output Tokens`、`KB Queries` 作为视角切换
- [ ] 页面不把这些视角写成新的独立 `Core` API

### US-003: 按时间或维度聚合

作为租户用户，我希望按天、按小时、按资源类型或按可用区聚合用量，以便观察趋势和结构差异。

**Acceptance Criteria**

- [ ] 页面支持 `group_by=resource_type|az|day|hour`
- [ ] 页面展示趋势图或表格聚合结果
- [ ] 页面不扩写平台结算、对账、计费规则

### US-004: 保持报表边界清晰

作为模块维护人，我希望当前报表页只定义 `Core / Metering` 已冻结的租户用量能力，以便后续与推理用量统计、平台结算拆分维护。

**Acceptance Criteria**

- [ ] 页面正文明确租户用量统计属于 `Core / Metering`
- [ ] 页面正文不把 `/metering/token-usage` 写成租户直接调用的查询接口
- [ ] 页面正文不把账单、发票、对账单写成当前模块正式契约

## 4. Functional Requirements

- FR-1: 系统必须提供租户用量查询视图
- FR-2: 系统必须支持时间范围筛选
- FR-3: 系统必须支持资源类型筛选
- FR-4: 系统必须支持分组维度切换
- FR-5: 系统必须把预设视角控制在同一报表页内，而不是拆成新的 API 契约
- FR-6: 系统必须保持租户报表而非平台结算口径

## 5. Business Rules

- `start_time` 与 `end_time` 为正式必填参数
- `group_by` 只允许 `resource_type`、`az`、`day`、`hour`
- `/api/v1/metering/token-usage` 属于写入侧上报接口，不是租户查询接口
- `GPU-Hours`、`CPU-Hours`、`Memory-GBHours`、`Storage-GBDays`、`Input Tokens`、`Output Tokens`、`KB Queries` 都是同一报表页的预设视角
- 页面展示的是租户侧聚合用量，不等于账单金额或平台级运营分析

## 6. Non-Goals

- 不在本轮实现账单、发票、结算与对账
- 不在本轮实现平台全局运营分析
- 不在本轮把各资源视角拆成独立 `Core` API
- 不在本轮暴露写入侧 `/metering/token-usage` 给租户直接调用

## 7. Design Considerations

- 首屏优先突出时间范围、资源视角和趋势结果
- 图表与表格必须共享同一查询上下文
- 无数据态、无权限态和查询失败态必须区分
- 预设视角只作为同一报表页的切换方式，不制造新的页面契约

## 8. Technical Considerations

- 权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 路径归属：`Core / Metering`
- 正式查询路径：`GET /api/v1/metering/usage`
- 必填查询参数：`start_time`、`end_time`
- 可选查询参数：`resource_type`、`group_by`
- 成功返回体形状：`{items,total,dev_profile}`
- 错误响应格式：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`

## 9. Success Metrics

- 用户能在 1 分钟内获得指定时间范围内的用量趋势
- 页面正文与 `v1.yaml` 中报表路径、参数和返回结构保持一致
- 资源视角切换不再被误写成多条独立 API 契约
- 文档中不再出现 `[Assumption]` 或把写入接口误写成查询接口

## 10. Open Questions

- 后续是否需要增加租户级消费金额口径
- 推理用量统计是否继续作为该页面的过滤视角而不是独立报表资源
- `KB Queries` 视角未来是否需要等待正式 `Metering` 资源补充后再落地

## 11. Backfill Dependencies

- 如未来扩展账单、发票、对账，必须先冻结正式路径与 schema
- 如未来扩展租户可见的 token 详细报表，必须先区分查询接口与写入接口
- 如未来扩展更细的 AI 服务计量视角，必须先明确是否仍复用 `/metering/usage`
