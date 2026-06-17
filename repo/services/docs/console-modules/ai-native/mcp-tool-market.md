# AI 原生 — MCP 工具市场

## 页面定位

租户浏览、订阅与安装 **MCP 工具包** 的市场页；安装后工具 ID 进入 `agent-tool-permission.md` 的授权矩阵。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-mcp-tool-market-draft.md`。

## 文档管理规则

- 本文是 **MCP 工具市场** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-mcp-tool-market-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 不得把草案路径写成**已实现**
- PRD/SPEC：`tasks/modules/prd/console/ai-native/prd-console-mcp-tool-market.md`、`tasks/modules/spec/console/ai-native/spec-console-mcp-tool-market.md`
- TASK：`TASK-SVC-018`

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/agent/mcp/tools` | `listMCPMarketToolPackages` | 3a |
| GET | `/api/v1/svc/agent/mcp/tools/{package_id}` | `getMCPMarketToolPackage` | 3a |
| POST | `/api/v1/svc/agent/mcp/tools/{package_id}/install` | `installMCPMarketToolPackage` | 3a |
| POST | `/api/v1/svc/agent/mcp/tools/{package_id}/uninstall` | `uninstallMCPMarketToolPackage` | 3b |

Schema（草案）：`MCPMarketToolPackage`、`MCPMarketInstallRequest`、`MCPMarketInstallation`。

RBAC（草案）：`scope:agent-mcp:read`、`scope:agent-mcp:write`。

<!-- TODO-YAML: 合入后删除「规划」标记 -->

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 读/写权限 | `agent-mcp:read` / `:write` | `403 FORBIDDEN` |
| install 幂等键 | POST 携带 `idempotency_key` | `400` |
| 包存在 | 有效 `package_id` | `404 NOT_FOUND` |
| 重复安装 | 未已 active | `422`（建议 `PACKAGE_ALREADY_INSTALLED`） |

## 页面职责

- 市场 catalog：包名、版本、publisher、risk_level、installed 标记
- 包详情：含展开 tool 列表预览
- 安装/卸载（3b）：异步时 `202 + AsyncTask`
- 安装完成后跳转 **工具权限** 配置 allow/deny

## 操作可用性矩阵

| 操作 | 只读用户 | Agent 管理员 | YAML 合入后 |
|---|---|---|---|
| 浏览/详情 | 可用 | 可用 | GET |
| 安装 | 不可用 | 可用 | POST install |
| 卸载 | 不可用 | 可用 | POST uninstall（3b） |
| 授权工具 | 跳转 | 跳转 | tool-permissions |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-mcp-tool-market-draft.md` §4。

### `GET /api/v1/svc/agent/mcp/tools`（规划）

- 成功：`200 + MCPMarketToolPackageListResponse`
- 错误：`401`、`403`

### `POST .../install`（规划）

- 成功：`202 + MCPMarketInstallation` 或 `201`
- 错误：`400`、`404`、`422`

## 待补边界

- 计费/metering 按包 — 待 BOSS/Services
- 社区包审核 — BOSS 运营，Console 只读 catalog
- catalog 与 `GET /agent/tools` 分工 — 市场 vs 已安装集合

## 相关模块

- `ai-native/agent-tool-permission.md`
- `integration/integration-third-party.md`

## 验收标准

- [ ] 安装后 tool_ids 与 listAgentTools 一致
- [ ] 不与 `/integrations` 混 path
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
