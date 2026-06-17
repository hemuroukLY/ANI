# 平台 Storage-GBDays

## 页面定位

`平台 Storage-GBDays` 是 BOSS **平台计量与结算** 域下的 **跨租户存储用量** 专页：按时间范围展示各租户对象/块存储占用（GB·天）排行、趋势与占比，供平台运营、SRE、财务做存储容量规划与成本分析。

本页属于 **Core / Metering** 视角下的 **平台 RBAC** 页面，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`。

Console [`usage-report.md`](../../console-modules/tenant/usage-report.md) 的「Storage-GBDays」仅为 **当前 JWT 租户** 的 UI 预设视角；**当前 `MeteringUsageRecord.resource_type` enum 不含 storage 类**，Console 与本页在 Core 合入前均可能呈现 **规划空态**。多指标聚合入口见 [`tenant-usage-billing.md`](../tenant/tenant-usage-billing.md)。

## 文档管理规则

- 本文是 `平台 Storage-GBDays` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-platform-storage-gbdays.md`](../../tasks/modules/prd/boss/metering/prd-boss-platform-storage-gbdays.md) 与 [`spec-boss-platform-storage-gbdays.md`](../../tasks/modules/spec/boss/metering/spec-boss-platform-storage-gbdays.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- 计量归属 **Core / Metering**（storage 类 resource_type **待 YAML 入 enum**）
- `GET /api/v1/metering/usage` 路径已声明，但 `MeteringUsageRecord.resource_type` **当前不含** storage
- **禁止** 把 UI 名「Storage-GBDays」硬编码为已冻结 `resource_type`
- 当前权威源 **未** 给 `/metering/usage` 声明 `operationId`，正文不得自造
- BOSS 跨租户 list/aggregate — **TODO-YAML**；不得把单租户 GET 直接写成 BOSS 正式契约
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422` 仅 YAML 已声明 operation 可写冻结（§2.10）；本页查询侧暂无 422

<!-- TODO-YAML: resource_type 如 storage_gb_days -->

## 页面职责

- 提供 **全平台** 时间范围内的存储（GB·天）消耗排行与趋势（Core enum 合入后）
- 提供按 AZ、day/hour 时间桶的聚合视图（平台 API 合入后）
- 支持从排行行钻取到存储基础设施、租户配额策略
- 在 Core 未合入前，提供 **诚实的规划空态**，不伪造计量数据
- 为 Phase 2 账单/导出预留边界，当前不写结算金额

## 页面结构

- 首屏至少包含：`时间筛选区`、`平台 KPI 区`、`Storage 趋势图`、`租户排行表`、`边界说明`
- 趋势图与排行表 **共享同一查询上下文**（start/end、group_by）
- 本页锁定 storage 类 `resource_type`（**待 YAML**，目标如 `storage_gb_days`），不在此页切换至 CPU/GPU/Token 等指标
- **enum 未合入时**：KPI/图表/表格展示规划空态，**不** 伪造全 0 数据

```text
平台 Storage-GBDays
├── 时间筛选（start_time / end_time）
├── 可选筛选（AZ / 租户状态 — 平台 API 待 YAML）
├── KPI：平台 Storage GB·Days 总量（enum 合入后）
├── 趋势图（group_by=day|hour）
├── 租户排行表（tenant_id / quantity / unit / 占比）
└── 行内钻取 → ops/storage-infrastructure / tenant-quota-policy
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | `GET /api/v1/metering/usage` | 路径已声明；**当前无** storage resource_type 映射 |
| Core | 平台 metering aggregate **TODO-YAML** | BOSS 正式 list/排行数据源 |
| Services | — | 本页 **不** 承接 Services 存储 list 作为 metering 契约 |

### 关键边界

- 「Storage-GBDays」是 **UI 展示名**；正式 `resource_type` **待 YAML**（目标如 `storage_gb_days`）
- **不得** 用 `/metering/usage` 隐式证明 storage 已计量
- 单租户 `/metering/usage` **不能** 直接作为 BOSS 跨租户正式契约
- 账单金额、发票、对账不在本页（Phase 2）
- [`../ops/storage-infrastructure.md`](../ops/storage-infrastructure.md) 为运维视角，**不是** 本页 metering 数据源

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 时间筛选区 | UI | `start_time`、`end_time` | 刷新 KPI/图表/表格 |
| 平台 KPI | 平台 aggregate **待 YAML** | 全平台 GB·Days 汇总 | — |
| Storage 趋势图 | 平台 aggregate **待 YAML** | `group_by=day\|hour` | 排行表联动 |
| 租户排行表 | 平台 aggregate **待 YAML** | `tenant_id` + `total_quantity` | 租户钻取 |
| 边界说明 | 规划项 | storage enum 待 Core 合入 | storage-infrastructure |

## BOSS 与 Console 分工

| 维度 | BOSS 平台 Storage-GBDays | Console Storage-GBDays 视角 |
|---|---|---|
| 范围 | 全平台多租户排行 | 单租户 |
| API | 平台 aggregate **待 YAML** | usage-report UI 预设（当前可能空态） |
| RBAC | 平台 metering 读 | 租户 JWT |
| 钻取 | storage-infrastructure / quota | storage-management |
| 账单 | Phase 2 | 不在 usage-report |

## 当前冻结事实

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/metering/usage` | 路径已声明；**当前无** storage resource_type |

| 字段 | 冻结值 |
|---|---|
| `MeteringUsageRecord.resource_type` | **不含** storage — **TODO-YAML**（目标如 `storage_gb_days`） |
| `group_by`（租户 API） | `resource_type` / `az` / `day` / `hour`（storage 合入后适用） |

| 能力 | 状态 |
|---|---|
| storage resource_type | **ADDED-TO-YAML** `storage_gb_days` |
| 跨租户 Storage 排行 | **TODO-YAML** |
| `group_by=tenant_id` | **TODO-YAML** |
| 平台 CSV 导出 | **TODO-YAML** |

## 字段级定义

### 查询字段（平台 API 目标 · 待 YAML 合入后对齐）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `start_time` | query | ✅ | date-time |
| `end_time` | query | ✅ | date-time |
| `resource_type` | query | 固定 | **待 YAML**（目标如 `storage_gb_days`） |
| `tenant_id` | query | 可选 | 平台 RBAC 下筛选；**不得** 越权 |
| `group_by` | query | 可选 | 平台扩展 `tenant_id` / `az` / `day` / `hour` |
| `tenant_status` | query | 可选 | 产品筛选；**待 YAML** |

### 返回字段

| 字段 | 来源 | 说明 |
|---|---|---|
| `items[].tenant_id` | 平台 aggregate **待 YAML** | 租户标识 |
| `items[].resource_type` | `MeteringUsageRecord` **待 YAML** | 目标如 `storage_gb_days` |
| `items[].total_quantity` | `MeteringUsageRecord` **待 YAML** | 聚合 GB·Days |
| `items[].unit` | `MeteringUsageRecord` **待 YAML** | 以后端为准 |
| `items[].period` | `MeteringUsageRecord` | 时间桶，可空 |
| `total` | list response | 条数 |
| `dev_profile` | `CoreDevProfileInfo` | 联调标记；非主展示 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `storage_gbdays_display` | 直接展示 `total_quantity`（当 unit 为 GB·Days 时） |
| `share_pct` | 租户占平台总量百分比 |
| `rank` | 排行序号 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态（enum 已合入） | KPI + 趋势图 + 排行表 | 共享同一查询上下文 |
| **规划空态（enum 未合入）** | 「Storage 计量 resource_type 待 Core 合入，当前暂无数据」 | **不** 伪造全 0 图表或假排行 |
| 无数据态（enum 已合入但无用量） | 「当前时间范围内暂无 Storage 用量」 | 与规划空态文案区分 |
| 平台 API 未就绪 | 说明「跨租户 aggregate 待 Core 合入」 | 不得伪装为生产已上线 |
| 查询失败态 | KPI/图表/表格分别失败提示 + 重试 | 保留筛选条件 |
| 无权限态 | 403 提示，不渲染假数据 | 平台 RBAC 拒绝 |
| 钻取无租户权 | 行内禁用或 403 提示 | 不得越权查看 tenant 明细 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `resource_type` | **待 YAML**（目标如 `storage_gb_days`） | 合入前不得写入已冻结 enum |
| `total_quantity` | 后端聚合原值 | number (double) |
| `unit` | 以后端为准；目标 GB·Days | string |
| Storage-GBDays（UI） | 与 `total_quantity` 一致（当 unit 匹配时） | 保留 2 位小数 |
| `period` | 与 `group_by` 时间桶一致 | string，可空 |
| `start_time` / `end_time` | ISO 8601 | date-time |

## 状态与能力口径

本页为 **只读查询页**，无资源状态机。AsyncTask / 租户 lifecycle 不在本页变更。

| 能力 | 说明 |
|---|---|
| 查询 | 只读（enum 合入后） |
| 导出 | **待 YAML** |
| 调整配额 | 跳转 [`tenant-quota-policy.md`](../tenant/tenant-quota-policy.md)，不在本页 PATCH |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 metering aggregate 读 RBAC | 已授权（scope **TODO-YAML** 待 Core 定义） | `403 FORBIDDEN` |
| 租户 `GET /metering/usage` 只读参考 | JWT 租户上下文；**非** BOSS 跨租户正式方案 | `403`（operation 已声明，**无** `x-ani-rbac-scope`） |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `start_time` / `end_time` | 合法且 start < end | `400 BAD_REQUEST` |
| storage enum 已合入 | Core 扩展完成 | 规划空态（非 HTTP 错误） |
| 钻取单租户 | 具备该租户平台查看权 | `403` |
| 导出（待 YAML） | export RBAC | **待 YAML** |

## 操作可用性矩阵

| 操作 | 平台只读 | 平台运营 | 财务只读 | 说明 |
|---|---|---|---|---|
| 查看跨租户 Storage 排行 | ✅（可空态） | ✅ | ✅ | aggregate 待 YAML |
| 切换 group_by | ✅ **待 YAML** | ✅ **待 YAML** | ✅ **待 YAML** | enum 合入后 |
| 钻取租户 | ✅ | ✅ | ✅ | 须 RBAC |
| 导出 CSV | ❌ | ✅ **待 YAML** | ✅ **待 YAML** | — |
| 调整配额 | ❌ | ✅ 跳转 | ❌ | 非本页 API |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。

## 接口冻结规则

### `GET /api/v1/metering/usage`（租户上下文 · 只读参考 · **非 BOSS 正式契约**）

- `operationId`：权威源 **未声明**
- `summary`：`查询租户用量统计`
- `tags`：`["Metering"]`
- `query.required`：`start_time`、`end_time`
- `query.optional`：`resource_type`、`group_by`
- **当前**：`MeteringUsageRecord.resource_type` enum **不含** storage；本页 **不得** 固定 `resource_type` 为已冻结值
- `group_by enum`：`resource_type`、`az`、`day`、`hour`（**不含** `tenant_id`）
- `success`：`200 + MeteringUsageResponse`
- `errors`：`400`、`401`、`403`

### 平台 aggregate（待补）

<!-- TODO-YAML: GET /api/v1/metering/usage/platform 或等价 -->

- 须支持 storage resource_type 与 `group_by=tenant_id`
- 合入前不得写入「已冻结」正文

## 使用规则

- **不得** 自造 `storage_gb_days` 为当前已冻结 enum 值
- **不得** 用 `/metering/usage` 隐式证明 storage 已计量
- enum 未合入时，页面须展示规划空态，**禁止** 伪造全 0 排行
- 平台 aggregate 未上线时，**禁止** 生产环境用逐租户 JWT 轮询作为正式方案
- 不得把 UI 名「Storage-GBDays」写入 API `resource_type`（合入前）

## 待补能力边界

- `MeteringUsageRecord.storage_gb_days` — **ADDED-TO-YAML**
- 平台跨租户 metering API — **ADDED-TO-YAML** P1 (`GET /api/v1/metering/usage/platform`)
- 平台 CSV 导出 — **TODO-YAML**
- 账单金额 / 发票 — Phase 2
- 按存储类型（对象/块）拆分 — Phase 3

## 响应示例

### enum 未合入时的页面空态说明（非 API 响应）

当前 Core 无 storage resource_type，BOSS 页面应展示：

> Storage 计量 resource_type 待 Core 合入（目标如 `storage_gb_days`），当前暂无跨租户 Storage 排行数据。

### 单租户 usage 查询成功（storage enum 合入后 · 只读参考 · 非 BOSS 正式契约）

```json
{
  "items": [
    {
      "resource_type": "storage_gb_days",
      "total_quantity": 5120,
      "unit": "gib_days",
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

> 注：`storage_gb_days` 为 **TODO-YAML** 目标值，合入前不得作为已冻结 enum 引用。

## 错误示例

### 时间范围非法

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be earlier than end_time",
  "request_id": "req-boss-storage-400-001"
}
```

### 无平台 metering 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-storage-403-001"
}
```

> **注**：适用于 **平台 metering aggregate（TODO-YAML）** 场景。`message` 为示例文案，**非** YAML 已冻结 `x-ani-rbac-scope` 名。租户 `GET /metering/usage` 的 `403` 为 operation 已声明返回码，**无** scope 字段。

## 相关模块

- [`../ops/storage-infrastructure.md`](../ops/storage-infrastructure.md)、[`tenant-quota-policy.md`](../tenant/tenant-quota-policy.md)
- [`tenant-usage-billing.md`](../tenant/tenant-usage-billing.md)
- Console：[usage-report.md](../../console-modules/tenant/usage-report.md)

## 回填验收标准

- [x] 满配章节齐全（对照 `boss-full-depth-checklist.md`）
- [x] 明确 storage 不在当前 enum
- [x] 规划空态与无数据态区分清晰
- [x] 含响应示例与错误示例
- [ ] storage resource_type YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
