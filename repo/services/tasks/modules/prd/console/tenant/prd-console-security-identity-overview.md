# PRD: Console 安全与身份概览

## 1. Introduction/Overview

`Console / 安全与身份 / 安全与身份概览` 是租户侧安全域总入口，用于帮助用户理解当前已冻结的认证、Token 与 API Key 能力，并明确安全域中仍待补的管理与合规边界。该页面直接承接 `Core / Auth`，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`。

当前已确认的正式能力：

- `POST /api/v1/auth/oidc/begin`
- `POST /api/v1/auth/token`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/api-keys`
- `POST /api/v1/auth/api-keys`
- `DELETE /api/v1/auth/api-keys/{key_id}`

当前仍属于待补边界的能力：

- 用户管理
- 角色与权限**编辑**（只读列表已在 **租户管理** 模块冻结）
- LDAP 配置中心
- API Key 审计、恢复、轮换、风险分析、合规导出
- 密钥管理、国密加解密、网络安全策略

以下能力已拆至独立模块 **租户管理**（`docs/console-modules/tenant/tenant-management.md`）：成员、角色只读、SSO 读写、Webhook。

## 2. Goals

- 让租户用户在安全域总入口内快速识别当前已可用的认证与凭证能力
- 保证总入口、`API Key 管理` 子页和权威源之间的命名、路径和边界一致
- 把未冻结安全能力明确压回“待补边界”，不再混入正式契约
- 保持 `Console` 租户视角，不混入平台审计、合规或运营口径

## 3. User Stories

### US-001: 查看安全域已冻结能力

作为租户用户，我希望在一个总入口内快速看到当前已冻结的登录、Token 和 API Key 能力，以便判断哪些安全动作已经可用。

**Acceptance Criteria**

- [ ] 页面明确列出 `OIDC begin / token / refresh / logout / API Key` 已冻结能力
- [ ] 页面明确这些能力归属 `Core / Auth`
- [ ] 页面不继续使用旧 `/api/v1/console/*` 路径

### US-002: 进入 API Key 管理子模块

作为租户用户，我希望从安全与身份概览进入 `API Key 管理` 子页，以便查看、创建和吊销 API Key。

**Acceptance Criteria**

- [ ] 页面提供 `API Key 管理` 子模块入口
- [ ] 总入口对 `API Key` 的描述与 `tenant/api-key-management.md` 一致
- [ ] 页面明确 `key_value` 只在创建成功时返回一次，不在总入口常驻回显

### US-003: 区分会话态与凭证态

作为租户用户，我希望页面能把登录会话链路和程序化接入凭证分开说明，以便我快速判断应该去处理登录问题还是接入问题。

**Acceptance Criteria**

- [ ] 页面把 `OIDC / Token / Logout` 作为会话链路说明
- [ ] 页面把 `API Key` 作为程序化接入凭证说明
- [ ] 页面不把两类问题混写成一个模糊的“安全设置页”

### US-004: 识别待补安全边界

作为模块维护人，我希望安全与身份概览只声明当前已冻结能力和待补边界，以便后续与 `Core` 团队对齐时不产生误解。

**Acceptance Criteria**

- [ ] 页面正文明确 **租户管理入口**已分流至 `tenant-management.md`（成员/角色/SSO/Webhook）
- [ ] 页面正文明确 `用户管理 / 角色与权限编辑 / LDAP / 密钥管理 / 国密加解密 / 网络安全策略` 仍为待补边界
- [ ] 页面正文不把待补能力写成正式接口表
- [ ] 页面正文不把 API Key 审计、风险分析、合规导出写成当前已冻结能力

## 4. Functional Requirements

- FR-1: 系统必须提供安全与身份总入口，并汇总当前已冻结的 `Core / Auth` 能力
- FR-2: 系统必须给出 `API Key 管理` 子模块入口与边界说明
- FR-3: 系统必须区分会话态链路与程序化接入凭证链路
- FR-4: 系统必须明确待补安全能力边界，不把未冻结能力写成正式契约
- FR-5: 系统必须统一使用标准错误结构 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`

## 5. Business Rules

- `OIDC begin / token / refresh / logout` 属于会话获取与维护链路；租户级 SSO **读写**见 **租户管理** 模块，不等于 LDAP 配置中心已经冻结
- `API Key` 属于程序化接入凭证链路，不等于审计、风控、导出能力已经冻结
- 总入口只负责汇总、导航和边界说明，不重写 `API Key 管理` 子页完整细节
- 总入口出现安全异常时，必须让用户知道问题落在“登录会话”还是“接入凭证”

## 6. Non-Goals

- 不在本轮扩写租户成员、用户、角色与权限的正式接口页
- 不在本轮扩写 LDAP / SSO 配置中心
- 不在本轮扩写 API Key 审计、恢复、轮换、风险分析、合规导出
- 不在本轮扩写密钥管理、国密加解密、网络安全策略独立资源页

## 7. Design Considerations

- 首屏要先让用户看懂“当前有哪些已冻结能力”，再进入具体子页
- 会话能力和 API Key 能力要分组展示，避免混淆
- 总入口保留子页入口和边界，不复制 `API Key` 的完整 CRUD 契约
- 待补能力使用说明态，不生成看似可点击的伪入口

## 8. Technical Considerations

- 权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 路径归属：`Core / Auth`
- 正式前缀：`/api/v1`
- 会话链路：`/auth/oidc/begin`、`/auth/token`、`/auth/refresh`、`/auth/logout`
- 凭证链路：`/auth/api-keys*`
- 页面不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`

## 9. Success Metrics

- 用户能在 30 秒内识别安全总入口当前已承接的能力范围
- 用户能区分“登录会话问题”和“程序化接入问题”的处理入口
- 文档中关于安全域的路径、schema、返回码与 `v1.yaml` 不再冲突
- 文档不再把用户、角色、SSO 配置等未冻结能力误写成正式接口

## 10. Open Questions

- 后续是否需要把 `LDAP / SSO` 配置中心拆成独立模块，而不是继续留在概览边界中
- 后续是否需要在总入口增加“长期未使用 API Key”或“会话异常”的摘要指标
- 后续是否需要把 API Key 审计与风险分析收口为独立安全模块

## 11. Backfill Dependencies

- 如未来扩展 `用户 / 角色 / LDAP / SSO`，必须先冻结正式路径与 schema
- 如未来扩展 `API Key` 审计、恢复、轮换或导出，必须先写回权威源
- 如未来扩展密钥管理与网络安全策略，必须先明确其资源归属与页面边界
