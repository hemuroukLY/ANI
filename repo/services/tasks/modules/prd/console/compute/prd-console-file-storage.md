# PRD: Console 文件存储

## 1. Introduction/Overview

`Console / 算力与云资源 / 存储管理 / 文件存储` 用于帮助租户用户查看、创建、删除和理解当前权限范围内的文件存储资源，并为 VM、容器、GPU 容器等工作负载提供共享文件能力入口。

## 2. Goals

- 让租户在文件存储页面内理解当前文件存储资源的边界、状态和基础能力
- 严格以现有 `Core v1.yaml` 为准，沉淀文件存储页面可直接用于 Core 对齐的主维护材料
- 清楚区分哪些文件存储能力已经在 `Core` 冻结，哪些仍为后续待补能力

## 3. User Stories

### US-001: 查看文件存储总览
**Description:** 作为租户用户，我希望在文件存储页面中快速理解当前文件存储资源的范围和边界，以便确认哪些能力已经可用。

**Acceptance Criteria:**
- [ ] 页面明确这是 `Console` 租户侧文件存储资源页
- [ ] 页面明确当前承接列表、详情、创建、删除
- [ ] 页面明确标出挂载目标、挂载关系详情、访问策略细化仍为待补能力

### US-002: 查看文件存储列表和详情
**Description:** 作为租户用户，我希望查看文件存储列表和详情，以便了解名称、协议、容量和 endpoint 等基础属性。

**Acceptance Criteria:**
- [ ] 页面提供文件存储列表入口
- [ ] 页面提供文件存储详情入口
- [ ] 查询接口对齐 `GET /api/v1/filesystems`
- [ ] 详情接口对齐 `GET /api/v1/filesystems/{filesystem_id}`

### US-003: 创建文件存储
**Description:** 作为租户用户，我希望创建文件存储，以便为工作负载准备共享文件能力。

**Acceptance Criteria:**
- [ ] 页面提供创建文件存储入口
- [ ] 创建接口对齐 `POST /api/v1/filesystems`
- [ ] 创建请求体包含 `name`、`size_gib`、`idempotency_key`
- [ ] 成功返回口径对齐 `201 + StorageFilesystem`

### US-004: 删除文件存储
**Description:** 作为租户用户，我希望删除不再使用的文件存储，以便清理租户内资源目录。

**Acceptance Criteria:**
- [ ] 页面提供删除文件存储入口
- [ ] 删除接口对齐 `DELETE /api/v1/filesystems/{filesystem_id}`
- [ ] 页面在删除前提示活跃挂载风险
- [ ] 文档不把 `409 CONFLICT` 写成当前已冻结删除契约

## 4. Functional Requirements

- FR-1: 系统必须提供文件存储页面总览，并明确已冻结能力与待补能力边界
- FR-2: 系统必须提供文件存储的查询、详情、创建和删除入口
- FR-3: 系统必须从认证上下文获取租户边界，不依赖前端传入 `tenant_id`
- FR-4: 系统必须对创建接口应用 `idempotency_key`

## 5. Non-Goals (Out of Scope)

- 不在本轮把挂载目标写成已冻结 Core 契约
- 不在本轮把访问策略细化写成已冻结 Core 契约
- 不在本轮实现 `BOSS` 侧平台文件存储池运营总览

## 6. Success Metrics

- 用户能在 30 秒内识别当前租户下有哪些文件存储资源
- 用户能在同一模块内完成文件存储的基础创建和查询
- 文档不再把挂载目标或访问策略写成已经冻结的 Core 契约
