# PRD: Console API Key 管理

## 1. Introduction/Overview

`Console / 安全与身份 / API Key` 是租户侧凭证管理页面，用于查看、创建和吊销当前权限范围内的 API Key。该页面直接承接 `Core / Auth` 的正式能力，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`。

当前已确认的正式能力：

- `GET /api/v1/auth/api-keys`
- `POST /api/v1/auth/api-keys`
- `DELETE /api/v1/auth/api-keys/{key_id}`

当前仍属于待补边界的能力：

- API Key 更新
- API Key 恢复或启停切换
- API Key 明文再次查看
- API Key 使用审计明细
- API Key 风险分析与合规导出

## 2. Goals

- 让租户用户在一个页面内完成 API Key 的查看、创建和吊销闭环
- 保证页面定义、路径、返回码和字段与 `Core / Auth` 权威源一致
- 明确 API Key 管理与安全审计、合规、风控模块的边界
- 保持 `Console` 租户视角，不混入平台侧运营或审计口径

## 3. User Stories

### US-001: 查看 API Key 列表

作为租户管理员或凭证管理员，我希望查看当前可见的 API Key 列表，以便确认哪些凭证仍然有效、哪些即将过期。

**Acceptance Criteria**

- [ ] 页面通过 `GET /api/v1/auth/api-keys` 读取当前租户可见的 API Key 列表
- [ ] 列表至少展示 `name`、`key_prefix`、`scopes`、`rate_limit_rpm`、`created_at`、`expires_at`、`last_used_at`、`is_active`
- [ ] 页面支持可选 `user_id` 过滤，但不要求前端显式传 `tenant_id`
- [ ] 列表不回显完整明文 `key_value`

### US-002: 创建 API Key

作为租户管理员或凭证管理员，我希望创建新的 API Key，以便供自动化任务、SDK、CLI 或外部客户端调用 ANI 能力。

**Acceptance Criteria**

- [ ] 页面通过 `POST /api/v1/auth/api-keys` 提交创建请求
- [ ] 创建表单承接 `name`、`user_id?`、`scopes`、`rate_limit_rpm?`、`expires_at?`
- [ ] `scopes` 至少 1 项，并遵循 `Core` 既定 pattern
- [ ] 创建成功按 `201 + CreateAPIKeyResponse` 处理
- [ ] 页面必须强调 `key_value` 仅在创建成功时返回一次，并要求用户立即保存

### US-003: 吊销 API Key

作为租户管理员或凭证管理员，我希望吊销不再需要的 API Key，以便及时回收访问权限。

**Acceptance Criteria**

- [ ] 页面通过 `DELETE /api/v1/auth/api-keys/{key_id}` 提交吊销请求
- [ ] 成功响应按 `200 + {status: revoked}` 处理
- [ ] 吊销动作必须有二次确认和明确影响提示
- [ ] 吊销成功后页面返回列表并刷新当前 key 状态

### US-004: 保持模块边界清晰

作为模块维护人，我希望 API Key 页面只沉淀当前 `Core / Auth` 已冻结的能力，以便后续与安全概览、审计或合规模块拆分维护时没有歧义。

**Acceptance Criteria**

- [ ] 页面正文明确 API Key 管理属于 `Core / Auth`
- [ ] 页面正文不把审计日志、风险分析、合规导出写成当前正式契约
- [ ] 页面正文不把“重新查看明文 key”写成已冻结能力
- [ ] 页面正文不继续使用旧 `/api/v1/console/*` 路径

## 4. Functional Requirements

- FR-1: 系统必须提供 API Key 列表读取能力，并仅展示 `APIKeyInfo` 可回指字段
- FR-2: 系统必须支持 API Key 创建，并严格承接 `CreateAPIKeyRequest`
- FR-3: 系统必须支持 API Key 吊销，并按 `{status: revoked}` 反馈结果
- FR-4: 系统必须强调 `key_value` 仅在创建响应中返回一次，后续不再回显
- FR-5: 系统必须统一使用标准错误结构 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- FR-6: 系统必须按角色权限控制查看、创建、吊销入口
- FR-7: 系统必须把未冻结能力保留为边界说明，不能伪装成正式接口

## 5. Business Rules

- API Key 管理属于 `Console` 租户侧身份管理，不属于 `BOSS` 平台审计
- `user_id` 为空时，按当前认证用户处理
- `scopes` 至少 1 项，且必须满足 `^scope:[a-z0-9_-]+:(\*|[a-z0-9_-]+)$`
- `rate_limit_rpm` 的正式范围为 `1 ~ 10000`，默认值为 `60`
- 吊销成功后，该 key 应被视为不可继续使用
- 只读用户可查看有权范围内的列表，但不显示创建和吊销入口

## 6. Non-Goals

- 不在本轮扩写 API Key 编辑或更新能力
- 不在本轮扩写 API Key 恢复、轮换或启停切换能力
- 不在本轮扩写 API Key 明文再次查看能力
- 不在本轮扩写 API Key 使用审计明细页
- 不在本轮扩写全平台 API Key 风险分析和合规导出

## 7. Design Considerations

- 列表页优先突出名称、前缀、状态、到期时间和最近使用时间
- 创建弹窗必须把“一次性明文返回”做显著提示
- 吊销入口必须带危险动作确认，不允许静默执行
- 页面反馈必须保留 `request_id` 以支持排障

## 8. Technical Considerations

- 权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 路径归属：`Core / Auth`
- 正式前缀：`/api/v1`
- 列表接口：`GET /api/v1/auth/api-keys`
- 创建接口：`POST /api/v1/auth/api-keys`
- 吊销接口：`DELETE /api/v1/auth/api-keys/{key_id}`
- 创建成功：`201 + CreateAPIKeyResponse`
- 吊销成功：`200 + {status: revoked}`
- 页面不要求前端显式传 `tenant_id`

## 9. Success Metrics

- 用户能在一个闭环中完成“查看列表 -> 创建 -> 保存明文 -> 返回列表复核”
- 用户能在列表中快速识别仍有效、已吊销或即将到期的 API Key
- 模块文档与 `v1.yaml` 的路径、schema、返回码和字段保持一致
- 文档中不再出现 `[Assumption]`、旧 `console` 路径或未冻结能力伪装

## 10. Open Questions

- 后续是否需要提供“轮换 API Key”而不是“吊销后重建”的正式能力
- 后续是否需要把 API Key 使用审计拆成 `Console` 独立页
- 后续是否需要在安全域增加“临近过期”或“长期未使用”的运营提示规则

## 11. Backfill Dependencies

- 如未来扩展 API Key 轮换或恢复能力，必须先冻结正式路径与 schema
- 如未来扩展 API Key 审计、风险分析或合规导出，必须先明确其模块归属
- 如未来扩展 API Key 明文再次查看能力，必须先完成专门的安全评审与权威源冻结
