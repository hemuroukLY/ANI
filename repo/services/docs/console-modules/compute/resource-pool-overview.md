# 资源池概览

## 页面定位

`资源池概览` 是算力域入口页，展示租户可见资源池/配额摘要（非 BOSS 全平台池运营）。

本页为 **Console UI 聚合页**，当前无独立资源池 CRUD 契约。

## 文档管理规则

- 本文是资源池概览主维护源
- 与 `home-resource-overview.md` 字段口径不得冲突
- 不得自造 `ResourcePool` schema

## Core 层要求

**当前无**独立 `GET /api/v1/resource-pools` 或 `/compute/pool/overview` 冻结路径。

<!-- TODO-YAML: 只读资源池/配额聚合 API；或明确复用 gpu-inventory + instances list 页面聚合 -->

### 当前可行页面聚合来源

| 摘要 | 建议 API | RBAC |
|---|---|---|
| GPU 库存 | `GET /api/v1/gpu-inventory`、`/gpu-inventory/occupancy` | `scope:gpu-inventory:read` |
| 实例规模 | `GET /api/v1/instances`（分 kind 计数） | `scope:instances:read` |
| 用量 | `GET /api/v1/metering/usage` | 待确认 scope |

## Services 层要求

- 本页不消费 Services 作为资源池事实来源
- BOSS 平台池运营不在 Console 范围

## 页面职责

- 算力域 landing：资源池摘要卡片 + 跳转 VM/GPU/实例/K8s
- 标明聚合时间与数据来源
- 单来源失败时局部降级

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 各 Core 读权限 | 对应 scope | `403 FORBIDDEN` |

本页无 POST/PUT；不涉及 `idempotency_key`。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员 |
|---|---|---|
| 查看摘要 | ✅ | ✅ |
| 跳转子模块 | ✅ | ✅ |
| 在本页创建资源池 | ❌ | ❌ |

## 接口冻结规则

本页**无独立冻结 API**。引用 operation 见 `home-resource-overview.md` §接口冻结规则，另加：

### `GET /api/v1/metering/usage`

- 成功：`200` + items 数组
- 错误：`400`、`401`、`403`
- Query 必填：`start_time`、`end_time`

## 待补边界

- 独立资源池/配额 API — **TODO-YAML**
- 配额超限与创建前置 422 的统一口径 — 见各资源 create operation
- 与 BOSS 池运营视图边界 — 产品文档待补

## 验收标准

- [ ] 标注 TODO-YAML
- [ ] 与 `home-resource-overview.md` 不重复定义冲突口径
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
