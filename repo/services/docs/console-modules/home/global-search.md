# 全局搜索

## 页面定位

Console 顶栏全局搜索（跨模块资源与页面导航），对应 HTML `#search` 输入框的产品化后端能力。

## 文档管理规则

- 本文是全局搜索主维护源
- 与 IA 规划一致，不自造 `SearchResult` schema 为冻结事实

## Services 层要求

**当前无** `GET /api/v1/svc/search` 冻结路径。

## Core 层要求

**当前无** `GET /api/v1/search` 冻结路径。

<!-- TODO-YAML: 全局搜索 query API（resource_type、keyword、cursor） -->

## 页面职责

- 定义搜索范围：实例、推理服务、知识库、API Key 等
- 结果分组、跳转 deep link
- 无 API 时不声称服务端全文检索

### 当前可行实现（文档口径）

- **Phase 1**：纯前端导航树过滤（现有 HTML 行为）
- **Phase 2**：并发调用各模块 list + 客户端合并（高成本、无统一排序）
- **Phase 3**：专用只读 search API

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | Phase 1 无后端；Phase 3 为 `401` |
| 各模块读权限 | 对应 scope | 单类型结果为空或隐藏 |

本页 Phase 3 搜索 API 若冻结 POST，须带 `idempotency_key` — **待 YAML**。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员 |
|---|---|---|
| 导航树过滤（Phase 1） | ✅ | ✅ |
| 跨模块 list 聚合（Phase 2） | ✅（有各模块权限） | ✅ |
| 服务端全文检索（Phase 3） | ❌ | ❌（待 YAML） |

## 接口冻结规则

**当前无冻结 search API。**

Phase 2 可引用的 list operation（示例，非 search 契约）：

| 资源 | 路径 | 成功码 |
|---|---|---|
| 实例 | `GET /api/v1/instances` | `200` |
| 推理服务 | `GET /api/v1/svc/inference-services` | `200` |
| 知识库 | `GET /api/v1/svc/knowledge-bases` | `200` |

各 operation 错误码见对应模块详文。

## 待补边界

- `GET /api/v1/search` 或 `/api/v1/svc/search` — **TODO-YAML**
- 搜索结果 ranking / 高亮 — 待产品定稿
- 搜索审计与速率限制 — 待安全评审

## 验收标准

- [ ] TODO-YAML 明确
- [ ] Phase 1/2/3 分阶段口径清晰
- [ ] 不自造 SearchResult schema 为冻结
