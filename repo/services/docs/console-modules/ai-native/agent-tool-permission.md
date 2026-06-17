# AI 原生 — 工具权限控制

## 页面定位

控制 Agent 可调用的 **工具/MCP 能力** 的租户级、角色级与 Agent Profile 级授权页，与 `tenant/role-permission-edit.md` 的 **API scope RBAC** 互补、不替代。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md`。

## 文档管理规则

- 本文是 **工具权限控制** 的主维护源
- **规划权威源**（合入前）：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 不得把草案路径写成**已实现**
- PRD/SPEC：`tasks/modules/prd/console/ai-native/prd-console-agent-tool-permission.md`、`tasks/modules/spec/console/ai-native/spec-console-agent-tool-permission.md`
- TASK：`TASK-SVC-018` 子项（Tool Permission）

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/agent/tools` | `listAgentTools` | 3a |
| GET | `/api/v1/svc/agent/tool-permissions` | `getAgentToolPermissions` | 3a |
| PUT | `/api/v1/svc/agent/tool-permissions` | `updateAgentToolPermissions` | 3a |
| GET | `/api/v1/svc/agent/tool-permissions/effective` | `getEffectiveAgentToolPermissions` | 3b |

Schema（草案）：`AgentToolDefinition`、`AgentToolGrant`、`AgentToolPermissionsSnapshot`、`AgentToolPermissionsUpdateRequest`。

RBAC（草案）：`scope:agent-tools:read`、`scope:agent-tools:write`。

<!-- TODO-YAML: 合入 services/v1.yaml 后删除「规划」标记 -->

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 读权限 | `scope:agent-tools:read` | `403 FORBIDDEN` |
| 写权限 | `scope:agent-tools:write` | `403 FORBIDDEN` |
| PUT 幂等键 | `idempotency_key` 必填 | `400` |
| grants 合法性 | `tool_id` 存在于目录 | `422`（建议 `TOOL_NOT_FOUND`；待 YAML PreconditionFailed） |
| subject 引用 | scope=role 时 `subject_id` 为有效 role_id | `422`（建议 `SUBJECT_NOT_FOUND`） |

**deny 优先于 allow** — 运行时与 Console 预览须一致。

## 页面职责

- **工具目录** Tab：`listAgentTools` — 展示 builtin/MCP 工具、风险等级
- **授权矩阵** Tab：`getAgentToolPermissions` — allow/deny × scope × subject
- 保存：PUT 全量 grants（类似 `kb-permissions.md` 矩阵保存）
- **有效权限预览**（3b）：选择角色/Profile 后调用 `effective` 端点
- 跳转：`mcp-tool-market.md`（上架）、`agent-session.md`（会话内 tool 调用）

## 页面结构

```text
工具权限控制
├── 工具目录（只读 list）
├── 授权矩阵
│   ├── 租户级默认
│   ├── 按角色覆盖
│   └── 按 Agent Profile 覆盖
├── 冲突提示（deny 覆盖 allow）
└── 保存（PUT + idempotency_key）
```

## 操作可用性矩阵

| 操作 | 只读用户 | Agent 管理员 | YAML 合入后 |
|---|---|---|---|
| 查看工具目录 | 可用 | 可用 | GET tools |
| 查看授权 | 可用 | 可用 | GET tool-permissions |
| 更新授权 | 不可用 | 可用 | PUT |
| 有效权限预览 | 可用 | 可用 | GET effective（3b） |
| 编辑租户 RBAC scope | 不可用 | 跳转 | `role-permission-edit.md` |

## 接口冻结规则

> **当前为规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md` §4。

### `GET /api/v1/svc/agent/tools`（规划）

- 成功：`200 + AgentToolDefinitionListResponse`
- 错误：`401`、`403`

### `GET /api/v1/svc/agent/tool-permissions`（规划）

- 成功：`200 + AgentToolPermissionsSnapshot`
- 错误：`401`、`403`

### `PUT /api/v1/svc/agent/tool-permissions`（规划）

- 成功：`200 + AgentToolPermissionsSnapshot`
- 错误：`400`、`401`、`403`、`409`

## 待补边界

- MCP 市场上架/订阅 — `mcp-tool-market.md`（目录数据来源）
- 单工具 rate limit — 非本模块
- 与 Gateway OpenAI tools 声明 — `openai-compatible-api.md` 边界
- 项目级 scope — 产品待定；当前草案 enum：tenant | role | agent_profile

## 相关模块

- `tenant/role-permission-edit.md` — API scope，非 tool grant
- `ai-native/agent-session.md` — 运行时消费授权
- `ai-native/agent-audit.md` — `policy_denied` 事件
- `ai-native/mcp-tool-market.md` — 工具来源

## 验收标准

- [ ] 草案与主维护源 operation 一致
- [ ] 不与 tenant roles PUT 混写契约
- [ ] deny/allow 规则在 UI 有说明
- [ ] 合入 YAML 后切换为冻结口径
- [ ] PRD/SPEC 已同步
