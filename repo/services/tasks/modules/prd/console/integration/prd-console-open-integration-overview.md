# PRD: Console 开放与集成总入口

## 1. Introduction/Overview

`Console / 开放与集成` 是租户侧接入总入口，用于帮助用户理解当前 ANI 平台已开放的接入面、鉴权准备项和文档入口，并明确哪些能力已经具备权威源支撑、哪些仍处于待补边界。

当前已确认可安全引用的接入面：

- `Core / Auth` 登录与凭证能力
  - `POST /api/v1/auth/oidc/begin`
  - `POST /api/v1/auth/token`
  - `POST /api/v1/auth/refresh`
  - `POST /api/v1/auth/logout`
  - `GET /api/v1/auth/api-keys`
  - `POST /api/v1/auth/api-keys`
  - `DELETE /api/v1/auth/api-keys/{key_id}`
- API 文档工件入口
  - `ANI-main/repo/docs/api/index.html`
  - `ANI-main/repo/docs/api/core.html`
  - `ANI-main/repo/docs/api/services.html`
- CLI 入口
  - `ANI-main/repo/cli/ani/main.go`

当前仍属于待补边界的接入面：

- OpenAI 兼容 API 独立接入页
- Go / Python / TypeScript / Java SDK 独立交付页
- Webhook
- 企业微信 / 钉钉 Bot
- 第三方业务系统集成

## 2. Goals

- 让租户在一个总入口内快速找到当前可用的接入方式与准备项
- 严格以现有权威源为准，沉淀可直接用于 `Core / Auth` 与接入文档对齐的主维护材料
- 清楚区分已冻结接入面、文档工件入口与待补接入能力
- 保持 `开放与集成` 总入口与 `安全与身份` 子页、服务模块文档口径一致

## 3. User Stories

### US-001: 查看已冻结接入面

作为开发者，我希望在开放与集成总入口里快速看到当前已可用的接入面，以便选择合适的对接方式。

**Acceptance Criteria**

- [ ] 页面明确展示 `REST API 文档 / API Key / OIDC / ani CLI` 的入口与边界
- [ ] 页面明确哪些内容来自现有权威源，哪些只是文档工件入口
- [ ] 页面不发明新的业务路径组

### US-002: 理解接入前置条件

作为开发者，我希望了解使用 REST API 或客户端接入前需要准备的鉴权与文档材料，以便减少试错。

**Acceptance Criteria**

- [ ] 页面明确 `API Key` 与 `OIDC` 在接入链路中的作用
- [ ] 页面给出 API 文档入口与客户端选择建议
- [ ] 页面不把 `tenant_id / X-Tenant-Id` 写成前端必传约束

### US-003: 识别接入路径与子页关系

作为开发者，我希望总入口能告诉我下一步该去哪里准备接入，而不是只给一组平铺名词。

**Acceptance Criteria**

- [ ] 页面能区分“交互式登录”“程序化调用”“文档查阅”“CLI 调试”四类入口
- [ ] 页面明确 `API Key 管理` 和 `安全与身份概览` 是前置准备子页
- [ ] 页面不在总入口重写子页的完整接口契约

### US-004: 识别待补接入边界

作为模块维护人，我希望开放与集成总入口只声明当前已冻结接入面和待补边界，以便后续与 `Core` 或 `Services` 团队对齐时没有歧义。

**Acceptance Criteria**

- [ ] 页面明确 `OpenAI 兼容 API / SDK / Webhook / Bot / 第三方业务系统集成` 当前属于待补边界
- [ ] 页面不把待补接入方式写成正式接口页
- [ ] 页面不把文档工件路径写成业务 API 路径

## 4. Functional Requirements

- FR-1: 系统必须提供开放与集成总入口，并汇总当前可确认的接入面
- FR-2: 系统必须提供鉴权准备与接入方式选择说明
- FR-3: 系统必须明确总入口与 `安全与身份` 子页、服务文档入口的关系
- FR-4: 系统必须明确待补接入能力边界
- FR-5: 系统必须保持文档入口与 API 契约分离

## 5. Business Rules

- 交互式登录走 `OIDC begin / token / refresh / logout`
- 程序化调用优先走 `API Key`
- REST API 文档入口不等于业务 API 路径
- `ani CLI` 当前只确认存在仓内入口，不等于已冻结独立下载与发布页
- OpenAI / SDK / Webhook / Bot / 第三方集成当前只保留边界说明

## 6. Non-Goals

- 不在本轮把 OpenAI 兼容 API 写成已冻结正式入口页
- 不在本轮实现 SDK 单独下载页或发布页
- 不在本轮实现 Webhook、Bot、第三方业务系统集成的正式配置页
- 不在本轮新增独立“接入资源域”后端路径

## 7. Design Considerations

- 首屏优先回答“我该怎么接入”和“我还缺什么准备项”
- 文档入口、凭证入口、CLI 入口需要按任务流组织，而不是平铺链接目录
- 待补接入面必须明显展示为说明态，避免用户误以为已可用
- 总入口只保留摘要与路径选择，不复制安全域和服务域的完整接口正文

## 8. Technical Considerations

- 权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 文档工件：`ANI-main/repo/docs/api/*.html`
- CLI 入口：`ANI-main/repo/cli/ani/main.go`
- 页面不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`

## 9. Success Metrics

- 用户能在 30 秒内识别当前可用接入方式与准备项
- 文档不再把待补接入方式误写成已冻结接口
- 文档不再把文档工件路径误写成业务路径
- 文档入口、凭证入口和 CLI 入口的逻辑关系清晰可复述

## 10. Open Questions

- 后续是否需要把 OpenAI 兼容 API 从推理服务边界提升为独立接入模块
- 后续 SDK 交付是否继续沿用文档工件入口，还是拆成语言级独立模块
- Webhook 和第三方系统集成是否属于同一集成域，还是拆成不同子域

## 11. Backfill Dependencies

- 如未来补 `OpenAI 兼容 API` 正式入口，必须先冻结正式路径与返回码
- 如未来补 `SDK / Webhook / Bot / 第三方业务系统集成` 正式模块，必须先补独立文档和契约来源
- 如未来补 CLI 下载与发布页，必须先明确其工件来源与版本策略
