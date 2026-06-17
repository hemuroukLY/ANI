# SPEC: Console 开放与集成总入口

> Technical specification derived from: `tasks/modules/prd/console/integration/prd-console-open-integration-overview.md`
> Source of truth: `ANI-main/repo/api/openapi/v1.yaml`
> Source of truth: `ANI-main/repo/api/openapi/services/v1.yaml`
> Scope: `Console / 开放与集成`

## 1. Summary

本 SPEC 定义 `Console / 开放与集成` 的技术边界、入口分类、权威源映射和总入口约束。该页面属于租户侧接入聚合页，只汇总当前已冻结 API 面、文档工件入口和 CLI 入口，不扩写新的后端资源，也不扩写 OpenAI / SDK / Webhook 的正式配置页。

## 2. Source of Truth

### 2.1 Primary Authorities

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Frozen API Sources

- `POST /api/v1/auth/oidc/begin`
- `POST /api/v1/auth/token`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/api-keys`
- `POST /api/v1/auth/api-keys`
- `DELETE /api/v1/auth/api-keys/{key_id}`
- `GET /api/v1/svc/models`
- `GET /api/v1/svc/inference-services`
- `GET /api/v1/svc/knowledge-bases`

说明：

- `Services` 路径在本页只作为“REST API 文档可覆盖的已冻结业务面”引用
- 总入口不在本页重写模型、推理服务、知识库的完整接口契约

### 2.3 Frozen Supporting Entrypoints

- `ANI-main/repo/docs/api/index.html`
- `ANI-main/repo/docs/api/core.html`
- `ANI-main/repo/docs/api/services.html`
- `ANI-main/repo/cli/ani/main.go`

### 2.4 Non-Frozen Areas

- OpenAI 兼容 API 独立接入页
- Go SDK / Python SDK / TypeScript Client / Java SDK 独立交付页
- Webhook
- 企业微信 / 钉钉 Bot
- 第三方业务系统集成

## 3. Entry Design

### 3.1 Entry Matrix

| 类型 | 路径或工件 | 用途 | 当前结论 |
|---|---|---|---|
| Auth | `/api/v1/auth/oidc/begin` | 发起 OIDC 登录 | 已冻结 |
| Auth | `/api/v1/auth/token` | OIDC callback 换取 TokenPair | 已冻结 |
| Auth | `/api/v1/auth/refresh` | 刷新 AccessToken | 已冻结 |
| Auth | `/api/v1/auth/logout` | 吊销当前 JWT JTI | 已冻结 |
| Auth | `/api/v1/auth/api-keys*` | 管理 API Key | 已冻结 |
| API docs | `ANI-main/repo/docs/api/index.html` | API 文档目录入口 | 已存在工件 |
| API docs | `ANI-main/repo/docs/api/core.html` | Core API 文档 | 已存在工件 |
| API docs | `ANI-main/repo/docs/api/services.html` | Services API 文档 | 已存在工件 |
| CLI source | `ANI-main/repo/cli/ani/main.go` | `ani CLI` 仓内入口 | 已存在入口 |
| Integration | `OpenAI / SDK / Webhook / Bot / 第三方集成` | 接入能力扩展 | 待补边界 |

### 3.2 Entry Rules

- 文档工件入口不是业务 API 路径
- CLI 源码入口不是下载页或发布页
- 总入口只负责“去哪准备接入”和“现在哪些能力可用”
- 子页完整契约分别回指 `安全与身份概览`、`API Key 管理`、服务模块主文档

## 4. Data Model

### 4.1 Authentication Preparation

| 场景 | 权威源 | 说明 |
|---|---|---|
| 交互式登录 | `Core / Auth` OIDC & token 链路 | 适合登录、控制台交互 |
| 程序化调用 | `Core / Auth` API Key 链路 | 适合脚本、SDK、CLI、服务调用 |

### 4.2 API Surface Coverage

| 接入面 | 权威源 | 说明 |
|---|---|---|
| Core API | `v1.yaml` | 基础设施与鉴权能力 |
| Services API | `services/v1.yaml` | 模型、推理服务、知识库等业务能力 |
| API docs 工件 | `docs/api/*.html` | 文档入口，不是 API 契约本体 |
| CLI 入口 | `cli/ani/main.go` | 仓内客户端源码入口 |

## 5. Page Design Constraints

### 5.1 Entry Page

- 总入口默认展示“推荐接入路径 + 接入准备项 + 已冻结接入面 + 待补边界”
- 页面要先回答用户“先去哪里准备”，再给链接入口
- 页面不复制 `API Key 管理`、`安全与身份概览`、`模型中心` 等文档的完整正文

### 5.2 Navigation Constraints

- 当前冻结导航未单列“开放与集成概览”页时，HTML 摘要允许落在 `REST API` 首屏说明中
- 该 HTML 落点只是一种摘要承载方式，不等于新增业务资源页
- 文档层主维护源仍然是 `docs/console-modules/integration/open-integration-overview.md`

## 6. Rules

- 本页本质是接入聚合页，不单独发明新的后端资源路径
- 已冻结 API 路径优先引用 `Core / Auth` 和现有 `Services` 文档入口覆盖面
- 页面不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- 页面不能把 OpenAI / SDK / Webhook / Bot / 第三方集成写成当前已冻结入口
- 页面不能把仓内工件路径写成对外业务 API 路径

## 7. Error Handling

### 7.1 Unified Error Shape

```json
{"code":"UPPER_SNAKE","message":"...","request_id":"..."}
```

### 7.2 Documentation Rule

- 总入口不新增自己的错误码表
- 若引用鉴权能力，错误结构沿用 `Core / Auth`
- 若引用服务能力，错误结构沿用对应服务模块

## 8. Acceptance

- 文档中所有正式 API 路径都能回指到现有权威源
- 文档中所有工件入口都能回指到仓内真实文件
- 文档入口与业务 API 路径不混写
- 待补接入面没有被回填成正式入口、下载页或配置页

## 9. Backfill Dependencies

- 如未来补 `OpenAI 兼容 API`，需先冻结正式路径与返回码
- 如未来补 `SDK / Webhook / Bot / 第三方业务系统集成`，需先补独立契约来源
- 如未来补 CLI 下载与发布页，需先明确工件来源、版本策略和入口关系
