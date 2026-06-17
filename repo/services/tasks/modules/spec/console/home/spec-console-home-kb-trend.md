# SPEC: Console home-kb-trend

> Source: `tasks/modules/prd/console/home/prd-console-home-kb-trend.md`  
> Revised: 2026-06-17

## 1. Summary

平台概览第四主题区块明细：知识库调用量、成功率与活跃库数趋势。**Services + Metering 聚合页**，无独立 KB 趋势 API。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/services/v1.yaml`（knowledge-bases）
- `ANI-main/repo/api/openapi/v1.yaml`（metering/usage）

### 2.2 Verified Paths

本页不新增 API。依赖以下只读接口：

| Method | Path | operationId | 成功响应 | 用途 |
|---|---|---|---|---|
| GET | `/api/v1/svc/knowledge-bases` | `listKnowledgeBases` | `200 + KnowledgeBaseListResponse` | 活跃知识库数（客户端计数） |
| GET | `/api/v1/metering/usage` | （以 YAML 为准） | `200 + UsageResponse` | 调用量趋势（knowledge 维度 filter 待确认） |

**无**独立 `GET /knowledge-trend` 类 API。

### 2.3 Verified Schemas

- `KnowledgeBaseListResponse`（以 YAML 为准）
- `UsageResponse`（以 YAML 为准）
- 问答成功率、异常解析数 — **无专用 API**，待产品口径

## 3. Page Scope

- 趋势图 + 时间窗口说明（7d/30d 等产品默认）
- 跳转：知识库管理 / 问答流程 / 用量报表
- 对「待确认」指标展示占位或隐藏，不伪造数据

## 4. Non-Goals

- 专用 KB 趋势 API — **YAML 未声明**
- 问答成功率/失败原因聚合 — 待 Services + Metering 联合口径
- 跨租户平台运营视图 — 非 Console 租户页职责
- 导出报表 — 见 `usage-report.md`
- **当前 YAML 未声明 `422`**

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 首页读权限 | 租户成员 | `403 FORBIDDEN` |

## 6. 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 查看趋势摘要 | ✅ | ✅ |
| 导出报表 | ❌ | ❌（见 `usage-report.md`） |
| 修改 KB | ❌ | 跳转 `knowledge-base.md` |

## 7. 主维护源

- `docs/console-modules/home/home-kb-trend.md`
- 父级：`platform-overview.md` 第四区块
- 相关：`knowledge-base.md`、`kb-qa-flow.md`、`tenant/usage-report.md`

## 8. 验收要点

- 待确认指标在 UI 有明确标注
- 不伪造 metering 不支持的维度曲线
- knowledge 维度 filter 使用前须对照 OpenAPI
