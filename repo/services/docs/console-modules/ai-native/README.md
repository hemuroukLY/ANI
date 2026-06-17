# AI 原生模块索引

Console **AI 原生**域子模块主维护源。整域 OpenAPI **待冻结**（`TODO-YAML`）；**×7 模块详文 + OpenAPI 详化草案均已就绪**（2026-06-17），供产品与 Core/Services 评审用。

## 模块清单

| 模块 | 详文 | OpenAPI | 阶段 |
|---|---|---|---|
| Sandbox 安全沙箱 | [ai-native-sandbox-security.md](./ai-native-sandbox-security.md) | **详化草案** | [openapi-phase3-ai-native-sandbox-security-draft.md](../openapi-drafts/phase3/openapi-phase3-ai-native-sandbox-security-draft.md) |
| Agent 会话 | [agent-session.md](./agent-session.md) | **详化草案** | [openapi-phase3-agent-session-draft.md](../openapi-drafts/phase3/openapi-phase3-agent-session-draft.md) |
| 工具权限控制 | [agent-tool-permission.md](./agent-tool-permission.md) | **详化草案** | [openapi-phase3-agent-tool-permission-draft.md](../openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md) |
| Agent 行为审计 | [agent-audit.md](./agent-audit.md) | **详化草案** | [openapi-phase3-agent-audit-draft.md](../openapi-drafts/phase3/openapi-phase3-agent-audit-draft.md) |
| Prompt 防护 | [prompt-injection-guard.md](./prompt-injection-guard.md) | **详化草案** | [openapi-phase3-prompt-injection-guard-draft.md](../openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md) |
| MCP 工具市场 | [mcp-tool-market.md](./mcp-tool-market.md) | **详化草案** | [openapi-phase3-mcp-tool-market-draft.md](../openapi-drafts/phase3/openapi-phase3-mcp-tool-market-draft.md) |
| Agent 编排 | [agent-orchestration.md](./agent-orchestration.md) | **详化草案** | [openapi-phase3-agent-orchestration-draft.md](../openapi-drafts/phase3/openapi-phase3-agent-orchestration-draft.md) |

## 维护规则

- 路径前缀：`/api/v1/svc/agent/*`（**评审前禁止写为冻结事实**）
- 算力 Sandbox 实例：`Core /api/v1/instances?kind=sandbox`；本域 **sandbox-profiles** 为 Agent 运行时安全配置
- 与 `knowledge/kb-qa-flow.md` 的 RAG 会话边界见各模块「待补边界」
- YAML 草案：上表 **×7 均已详化**；合入 ANI-main 前详文「接口冻结规则」保持**规划**标记

## 评审建议顺序

1. Session → Tool Permission → Audit → Prompt Guard（主流程）
2. MCP Market → Sandbox Security（安装与隔离）
3. Orchestration（依赖 Session + AsyncTask）

## 相关

- Backlog：`../governance/console-undefined-features-backlog.md` P2-01～P2-07
- 整域索引：`../openapi-drafts/phase3/openapi-phase3-domain-draft.md` §1
- TASK：`tasks/execution/SERVICES-TEAM-TASKS.md` — TASK-SVC-018
