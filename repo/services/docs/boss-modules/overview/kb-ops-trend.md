# 知识库运营态势

## 页面定位

`知识库运营态势` 是 BOSS **平台运营总览** 域下的 **全平台知识库运营** 趋势明细页，展示知识库规模、查询量趋势、解析失败与跨租户活跃排行。

Console 对照：[`home-kb-trend.md`](../../console-modules/home/home-kb-trend.md)。深度专页：[`kb-monitoring.md`](../health/kb-monitoring.md)、[`platform-kb-queries.md`](../metering/platform-kb-queries.md)。

## 文档管理规则

- 本文是 `知识库运营态势` 的主维护源
- PRD/SPEC 为辅助
- 流程：ANI-14 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- **Services** `GET /api/v1/svc/knowledge-bases`（`listKnowledgeBases`）— 租户 **只读参考**
- **Core** `GET /api/v1/observability/query` — PromQL
- Core metering KB `resource_type` — metering 域 **TODO-YAML**；`GET /metering/usage` 租户 **只读参考**
- BOSS 跨租户 KB aggregate — **TODO-YAML**
- **不得** 逐租户 JWT 轮询

## 页面职责

- 展示全平台知识库总数、活跃数、查询量趋势
- 解析失败/索引异常摘要（aggregate 或 Services 推断）
- 跨租户 KB 活跃 Top N（**TODO-YAML**）
- 跳转知识库监控、平台 KB Queries 计量
- 标明时间窗口与刷新时间

## 页面结构

```text
知识库运营态势
├── 顶部 KPI（KB 总数 / 活跃 / 查询量 / 失败数）
├── 查询量趋势（7d / 30d）
├── 跨租户活跃排行（待 YAML）
└── 跳转
    ├── 知识库监控
    └── 平台 KB Queries
```

## 数据来源与分层约束

| 层 | 路径 | 用法 |
|---|---|---|
| Services | `GET /api/v1/svc/knowledge-bases` | 租户 **只读参考** |
| Core | `GET /api/v1/observability/query` | PromQL |
| Core | `GET /api/v1/metering/usage` | KB queries enum **待 YAML** |
| Core | 平台 KB aggregate **TODO-YAML** | BOSS 正式数据源 |

### 关键边界

- metering 域 `kb_queries` resource_type enum 待 YAML 确认
- Services list 非 BOSS 正式契约
- 文档解析失败须以 `KBDocument` 等 YAML 字段为准，禁止自造

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| 规模 KPI | aggregate **待 YAML** | kb-monitoring |
| 查询趋势 | metering + PromQL | platform-kb-queries |
| 失败摘要 | Services / aggregate | kb-monitoring |
| 租户排行 | **TODO-YAML** | tenant-list |

## BOSS 与 Console 分工

| 维度 | BOSS | Console home-kb-trend |
|---|---|---|
| 范围 | 全平台 KB 运营 | 当前租户 KB |
| 排行 | 跨租户 | 无 |
| 计量 | 平台 aggregate 待 YAML | 租户 metering |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/svc/knowledge-bases` | `listKnowledgeBases` | 租户 Services |
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL |
| GET | `/api/v1/metering/usage` | — | 租户；KB enum **待 YAML** |

| 能力 | 状态 |
|---|---|
| 平台 KB aggregate | **ADDED-TO-YAML** `getPlatformKBTrend` |
| metering `kb_queries` resource_type | **TODO-YAML**（metering 域） |

## 字段级定义

| 字段 | 说明 |
|---|---|
| `kb_total` | 全平台知识库数 |
| `kb_active` | 窗口内有查询的 KB 数 |
| `queries_window` | 窗口查询总量 |
| `parse_failure_count` | 解析失败文档数 |
| `success_rate` | 查询成功率 |
| `tenant_rank[]` | **待 YAML** |
| `last_refreshed_at` | UI |
| `time_window` | UI |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 正常 | KPI + 趋势 |
| metering enum 未冻结 | 查询量卡片「待 YAML」 |
| aggregate 未就绪 | 标注 TODO-YAML |
| `parse_failure_count` > 0 | 高亮 |
| 403 | 无权限 |

## 字段口径与单位

| 字段 | 口径 |
|---|---|
| `queries_window` | 整数；须标注窗口 |
| `success_rate` | 0–100 % |

## 状态与能力口径

本页 **只读**；KB 生命周期状态以 Services schema 为准。

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| 平台读 RBAC（**TODO-YAML**） | `403` |
| 未认证 | `401` |
| metering 时间范围无效 | `400` |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 管理员 |
|---|---|---|---|
| 查看态势 | ✅ | ✅ | ✅ |
| 跳转 KB 监控/计量 | ✅ | ✅ | ✅ |
| 在本页重建索引 | ❌ | Phase 2 | Phase 2 |

## 删除前置校验

**N/A**

## 接口冻结规则

### `GET /api/v1/svc/knowledge-bases`（租户 · **只读参考**）

- `success`：`200`；当前 YAML GET 未声明额外错误 responses
- 认证/授权错误为 Services 网关通用推断，**未** 在本 operation responses 中冻结

### `GET /api/v1/metering/usage`（租户 · **只读参考**）

- `group_by` **不含** `tenant_id`
- KB `resource_type` enum — **TODO-YAML**

### 平台 KB aggregate（待补）

<!-- TODO-YAML: GET /api/v1/observability/knowledge-bases/platform 或等价 -->

## 使用规则

- 查询趋势须标注 `time_window`
- 禁止跨租户 JWT 轮询
- metering 与 monitoring 口径分工见专页

## 待补能力边界

- 平台 KB aggregate — **ADDED-TO-YAML**
- metering KB resource_type — **TODO-YAML**
- reindex Skill — Phase 2

## 响应示例

### 平台 aggregate 目标（**待 YAML**）

```json
{
  "kb_total": 89,
  "kb_active": 72,
  "queries_window": 120000,
  "parse_failure_count": 15,
  "success_rate": 99.2,
  "time_window": "7d",
  "top_tenants": [{ "tenant_id": "t-001", "queries": 45000 }]
}
```

## 错误示例

### 无平台读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-kot-403-001"
}
```

### metering 时间参数缺失

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time and end_time are required",
  "request_id": "req-boss-kot-400-001"
}
```

## 相关模块

- [`operations-overview.md`](operations-overview.md)、[`kb-monitoring.md`](../health/kb-monitoring.md)
- [`platform-kb-queries.md`](../metering/platform-kb-queries.md)

## 回填验收标准

- [x] 满配章节齐全
- [x] Services 租户路径 vs 平台 aggregate 分层
- [x] metering KB enum 标注 TODO-YAML
- [x] 400 + 403 错误示例
- [ ] platform aggregate YAML 合入后回写
- [x] PRD/SPEC/HTML 与本文同步
