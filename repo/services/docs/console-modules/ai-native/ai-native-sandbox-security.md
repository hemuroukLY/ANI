# AI 原生 — Sandbox 安全沙箱

## 页面定位

AI 原生域下的 **Agent 运行时安全沙箱** 配置与策略页，与算力域 `kind=sandbox` 实例（`sandbox-instance-management.md`）不同：本页面向 Agent 代码执行隔离、egress 策略与安全摘要，不替代 Core 实例 CRUD。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-ai-native-sandbox-security-draft.md`。

## 文档管理规则

- 本文是 **AI 原生 Sandbox 安全** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-ai-native-sandbox-security-draft.md`
- Core 只读引用：`instances?kind=sandbox`、`security-events`（见 `sandbox-templates.md`）
- **禁止** Services `/sandboxes*`（deprecated）
- PRD/SPEC：`tasks/modules/prd/console/ai-native/prd-console-ai-native-sandbox-security.md`、`tasks/modules/spec/console/ai-native/spec-console-ai-native-sandbox-security.md`
- TASK：`TASK-SVC-018`

## Core 层要求（只读引用）

| 能力 | 路径 | 说明 |
|---|---|---|
| Sandbox 实例列表 | `GET /api/v1/instances?kind=sandbox` | 关联展示，非本模块写 API |
| 安全事件 | `GET /api/v1/instances/{id}/security-events` | 算力域事件；Console 跳转 |

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/agent/sandbox-profiles` | `listAgentSandboxProfiles` | 3a |
| GET | `/api/v1/svc/agent/sandbox-profiles/{profile_id}` | `getAgentSandboxProfile` | 3a |
| PUT | `/api/v1/svc/agent/sandbox-profiles/{profile_id}` | `updateAgentSandboxProfile` | 3a |
| GET | `/api/v1/svc/agent/sandbox-security/summary` | `getAgentSandboxSecuritySummary` | 3b |

Schema（草案）：`AgentSandboxProfile`、`AgentSandboxEgressRule`、`AgentSandboxProfileUpdateRequest`。

RBAC（草案）：`scope:agent-sandbox:read`、`scope:agent-sandbox:write`。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 读/写权限 | agent-sandbox scope | `403` |
| PUT 幂等键 | 必填 | `400` |
| profile 存在 | 有效 id | `404` |
| egress 规则合法 | CIDR/域名格式 | `422`（建议 `INVALID_EGRESS_RULE`） |

default profile 不可 DELETE。

## 页面职责

- Sandbox Profile 列表/编辑：isolation_level（standard/strict/locked）、资源上限、egress 白名单
- `allow_tool_ids` 与 tool-permission **交集**生效（UI 说明）
- 安全摘要（3b）：blocked egress、violations 计数
- 跳转：算力 Sandbox 实例、security-events、Prompt Guard

## 操作可用性矩阵

| 操作 | 只读用户 | 安全管理员 | 说明 |
|---|---|---|---|
| 查看 profile | 可用 | 可用 | GET |
| 编辑 profile | 不可用 | 可用 | PUT |
| 安全摘要 | 可用 | 可用 | GET summary（3b） |
| Core 实例 CRUD | 跳转 | 跳转 | sandbox-instance-management |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-ai-native-sandbox-security-draft.md` §4。

### `PUT /api/v1/svc/agent/sandbox-profiles/{profile_id}`（规划）

- 成功：`200 + AgentSandboxProfile`
- 错误：`400`、`404`、`422`

## 待补边界

- profile → Core 实例映射 — Services 实现说明
- agent-audit 扩展 `sandbox_egress_denied` — 评审时与 audit 草案同步
- BOSS 全局 baseline profile 下发 — Console 租户覆盖边界

## 相关模块

- `compute/sandbox-instance-management.md`
- `compute/sandbox-templates.md`
- `ai-native/prompt-injection-guard.md`
- `ai-native/agent-tool-permission.md`

## 验收标准

- [ ] 无 Services sandboxes 路径
- [ ] 与 Core instances 边界清晰
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
