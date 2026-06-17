# 首页 — 知识库调用趋势

## 页面定位

平台概览第四主题区块明细：知识库调用量、成功率与活跃库数趋势。

本页属于 **Services + Metering 聚合**，无独立 KB 趋势 API。

## 文档管理规则

- 本文是首页 KB 趋势区块的主维护源
- 父级聚合入口：`platform-overview.md` 第四区块
- 指标口径必须在 UI 标注「待确认」处显式说明，禁止伪造曲线
- 一级权威源：OpenAPI 只读引用 `services/v1.yaml`（knowledge-bases）与 `v1.yaml`（metering/usage）

## 数据来源与分层

**无**独立 `GET /knowledge-trend` 类 API。页面层聚合：

| 摘要字段 | 建议来源 | 冻结状态 |
|---|---|---|
| 调用量趋势 | `GET /api/v1/metering/usage`（按 knowledge 维度 filter，若 YAML 支持） | **部分** — filter 维度待确认 |
| 活跃知识库数 | `GET /api/v1/svc/knowledge-bases` + 客户端计数 | 可用（列表 API 已声明） |
| 问答成功率 | query 成功/失败比 — **无专用 API** | 待产品口径 |
| 异常解析数 | 文档上传失败 / query 错误 — **无专用 API** | 待产品口径 |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 首页读权限 | 租户成员 | `403 FORBIDDEN` |

本页无写接口。**当前 YAML 未声明 `422`**。

## 页面职责

- 趋势图 + 时间窗口说明（7d/30d 等产品默认）
- 跳转：知识库管理 / 问答流程 / 用量报表
- 对「待确认」指标展示占位或隐藏，不伪造数据

## 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 查看趋势摘要 | 可用 | 可用 |
| 导出报表 | 不可用 | 不可用（见 `usage-report.md`） |
| 修改 KB | 不可用 | 跳转 `knowledge-base.md` |

## 接口冻结规则（聚合读）

本页不新增 API。依赖接口：

### `GET /api/v1/svc/knowledge-bases`

- 成功：`200 + KnowledgeBaseListResponse`
- 错误：`401`、`403`

### `GET /api/v1/metering/usage`

- 成功：`200 + UsageResponse`（以 YAML 为准）
- 错误：`401`、`403`
- **注意**：knowledge 维度 filter 是否在 YAML 冻结 — 使用前须对照 OpenAPI

## 待补边界

- 专用 KB 趋势 API — **YAML 未声明**
- 问答成功率/失败原因聚合 — 待 Services + Metering 联合口径
- 跨租户平台运营视图 — 非 Console 租户页职责

## 相关模块

- `platform-overview.md` — 第四区块
- `knowledge-base.md`、`kb-qa-flow.md`
- `tenant/usage-report.md` — 正式用量报表

## 验收标准

- [ ] 待确认指标在 UI 有明确标注
- [ ] 不伪造 metering 不支持的维度曲线
- [ ] 跳转入口与 IA 一致
