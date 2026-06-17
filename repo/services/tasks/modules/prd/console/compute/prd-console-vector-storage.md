# PRD: Console 向量存储

## 1. Introduction/Overview

`Console / 算力与云资源 / 存储管理 / 向量存储` 用于帮助租户用户查看、创建、删除和搜索当前权限范围内的向量存储资源，并为知识库、语义检索等场景提供统一的向量资源入口。

## 2. Goals

- 让租户在向量存储页面内理解当前向量资源的边界、状态和基础能力
- 严格以现有 `Core v1.yaml` 为准，沉淀向量存储页面可直接用于 Core 对齐的主维护材料
- 清楚区分哪些向量存储能力已经在 `Core` 冻结，哪些仍为后续待补能力

## 3. User Stories

### US-001: 查看向量存储总览
**Description:** 作为租户用户，我希望在向量存储页面中快速理解当前向量资源的范围和边界，以便确认哪些能力已经可用。

**Acceptance Criteria:**
- [ ] 页面明确这是 `Console` 租户侧向量存储资源页
- [ ] 页面明确当前承接列表、详情、创建、删除和搜索
- [ ] 页面明确标出向量写入、索引维护、批量导入仍为待补能力

### US-002: 查看向量存储列表和详情
**Description:** 作为租户用户，我希望查看向量存储列表和详情，以便了解名称、维度和距离度量等基础属性。

**Acceptance Criteria:**
- [ ] 页面提供向量存储列表入口
- [ ] 页面提供向量存储详情入口
- [ ] 查询接口对齐 `GET /api/v1/vector-stores`
- [ ] 详情接口对齐 `GET /api/v1/vector-stores/{vector_store_id}`

### US-003: 创建向量存储
**Description:** 作为租户用户，我希望创建向量存储，以便为知识库和语义检索场景准备底层向量资源。

**Acceptance Criteria:**
- [ ] 页面提供创建向量存储入口
- [ ] 创建接口对齐 `POST /api/v1/vector-stores`
- [ ] 创建请求体包含 `name`、`dimension`、`idempotency_key`
- [ ] 成功返回口径对齐 `201 + VectorStore`

### US-004: 搜索向量存储
**Description:** 作为租户用户，我希望在指定向量存储中执行搜索，以便验证检索能力是否可用。

**Acceptance Criteria:**
- [ ] 页面提供向量搜索入口
- [ ] 搜索接口对齐 `POST /api/v1/vector-stores/{vector_store_id}/search`
- [ ] 搜索请求体至少包含 `vector`
- [ ] 页面不把向量写入或索引维护写成现有冻结动作

## 4. Functional Requirements

- FR-1: 系统必须提供向量存储页面总览，并明确已冻结能力与待补能力边界
- FR-2: 系统必须提供向量存储的查询、详情、创建、删除和搜索入口
- FR-3: 系统必须从认证上下文获取租户边界，不依赖前端传入 `tenant_id`
- FR-4: 系统必须对创建接口应用 `idempotency_key`

## 5. Non-Goals (Out of Scope)

- 不在本轮把向量写入写成已冻结 Core 契约
- 不在本轮把索引维护或批量导入写成已冻结 Core 契约
- 不在本轮实现 `BOSS` 侧平台向量基础设施运营总览

## 6. Success Metrics

- 用户能在 30 秒内识别当前租户下有哪些向量存储资源
- 用户能在同一模块内完成向量存储的基础创建和搜索
- 文档不再把向量写入、索引维护、批量导入写成已经冻结的 Core 契约
