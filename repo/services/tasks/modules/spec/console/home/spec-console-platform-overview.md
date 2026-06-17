# SPEC: Console 平台概览

> Technical specification derived from: `tasks/modules/prd/console/home/prd-console-platform-overview.md`
> Generated: 2026-06-14

## 1. Summary

### 1.1 What This SPEC Covers

本 SPEC 定义 `Console / 首页与总览 / 平台概览` 的技术边界、页面结构、数据来源约束和回填前置依赖。该页面是租户侧 Console 聚合首页，不新增底层资源归属，不在当前阶段定义独立首页聚合契约。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/home/prd-console-platform-overview.md`
- User Stories covered: `US-001 ~ US-004`
- Functional Requirements covered: `FR-1 ~ FR-7`

### 1.3 Design Decisions Summary

| Decision | Choice | Rationale |
|---|---|---|
| 页面归属 | Console 聚合首页 | 首页承接摘要和导航，不是独立资源管理页 |
| 数据来源 | Core + Services 子模块 | 首页只消费各来源模块的已冻结事实 |
| 独立首页接口 | 当前未冻结 | 不能把待补聚合能力写成正式契约 |
| 租户边界 | 认证上下文获取 | 不信任前端显式传 `tenant_id` |
| 错误结构 | 统一错误结构 | 与子模块和 Core 标准保持一致 |

## 2. Architecture

### 2.1 System Context

```text
浏览器 Console
    -> 已冻结子模块接口
        -> Core 资源摘要
        -> Services 业务摘要
    -> 页面层轻聚合
    -> 平台概览 5 大主题区块
```

约束说明：

1. 当前平台概览不增加新的冻结路径
2. 首页只对既有模块事实做摘要，不改写子模块状态机
3. 未来如需后端聚合，必须先冻结只读契约，再更新本 SPEC

### 2.2 Component Design

- 页面壳：顶部总判断、5 大主题区块、快捷跳转入口
- 聚合层：对来源模块结果做轻量拼装，不创造新资源对象
- 状态层：负责空态、局部失败态、无权限态和刷新时间
- 跳转层：只指向当前可用的 Console 模块入口

### 2.3 Module Interactions

```text
用户进入平台概览
  -> 加载当前工作区状态与各主题区块摘要
  -> 某区块正常
      -> 展示摘要值 + 刷新时间 + 查看详情入口
  -> 某区块失败
      -> 展示局部失败态 + 刷新入口
  -> 用户点击区块入口
      -> 跳转到对应子模块继续处理
```

## 3. Frozen Facts

### 3.1 Frozen Paths

平台概览当前**没有独立冻结路径**。本页只能引用当前可用模块自身已冻结路径作为事实来源。

### 3.2 Frozen Schemas

平台概览当前**没有独立冻结 schema**。本页不能把首页返回体写成 `ConsolePlatformOverviewResponse` 等自造 schema。

### 3.3 Frozen Response Codes

平台概览当前**没有独立冻结返回码**。页面只能复用各来源模块已冻结的返回码与统一错误结构。

### 3.4 Non-Frozen Capabilities

- 首页后端只读聚合接口
- 首页专属返回体
- 首页专属状态枚举
- 首页专属 top items 结构
- 首页统一 drill-down 查询参数规范

### 3.5 Known Risky Assumptions

- 假设前端可以先通过并发调用各子模块接口完成首版首页聚合
- 假设首页摘要字段不会反向变成来源模块的正式资源字段
- 假设告警与待处理事项当前仍属于聚合摘要，而非独立模块契约

## 4. API Design

### 4.1 Independent Endpoint Status

| 项 | 当前结论 |
|---|---|
| 独立首页聚合 API | 未冻结 |
| 独立首页 schema | 未冻结 |
| 独立首页 operationId | 未冻结 |
| 独立首页返回码 | 未冻结 |

### 4.2 Source Module Dependencies

平台概览依赖以下当前可用模块的冻结事实，但本 SPEC 不重复改写这些子模块的正式接口定义：

| 区块 | 来源模块 | 依赖方式 |
|---|---|---|
| 资源使用概览 | `VM / GPU / 网络 / 存储 / K8s` | 读取子模块已冻结摘要事实 |
| GPU 利用率 | `GPU 算力管理` | 读取已冻结 GPU 摘要事实 |
| 推理服务状态 | `推理服务` | 读取 `Services` 已冻结服务摘要事实 |
| 知识库调用趋势 | `知识库管理` | 读取 `Services` 已冻结知识库摘要事实 |
| 告警与待处理事项 | 现有子模块异常摘要 | 仅聚合为首页摘要，不新增独立资源契约 |
| 用量相关指标 | `租户用量报表` | 读取 `Core / Metering` 已冻结聚合读能力 |

### 4.3 Error Handling

- 首页局部失败只影响对应区块，不阻断整页
- 页面层错误展示必须沿用统一错误结构
- 无权限区块不伪造摘要值，可展示隐藏、置灰或无权限提示

## 5. View Model Rules

### 5.1 Required View Sections

- `顶部总判断`
- `资源使用概览`
- `GPU 利用率`
- `推理服务状态`
- `知识库调用趋势`
- `告警与待处理事项`
- `快捷跳转入口`

### 5.2 Display Rules

- 首页只展示摘要，不展示子模块完整明细
- 区块必须附带刷新时间或统计窗口
- 无数据与真实值为 0 必须区分
- 跳转目标只能是当前可用模块

## 6. Security and Tenant Boundary

- 首页不接受前端显式传 `tenant_id / X-Tenant-Id`
- 首页所有数据必须基于当前登录租户的认证上下文
- 无权限模块不得在首页伪装成可用入口

## 7. Review Constraints

- 不允许把首页聚合路径写成当前已冻结事实
- 不允许自造首页专属 schema、字段、状态和错误码
- 不允许把子模块待补能力提前升级为首页正式能力
- 不允许把首页跳转入口写成已经存在的独立资源接口

## 8. Acceptance

| 验收项 | 通过标准 |
|---|---|
| 聚合定位 | 明确为聚合页，不伪装成独立资源页 |
| 路径边界 | 不出现旧 `/api/v1/console/home/*` 路径 |
| 契约边界 | 不自造独立首页接口、schema 和返回码 |
| 状态处理 | 具备空态、局部失败态、无权限态 |
| 跳转闭环 | 每个区块都有对应明细入口 |

## 9. Backfill Dependencies

- 如未来新增首页后端聚合接口，必须先冻结只读路径、response schema 和返回码
- 如未来将告警与待处理事项独立成模块，必须先明确其资源归属和冻结路径
- 如未来要求首页统一下钻参数，必须先与对应子模块统一筛选参数口径
