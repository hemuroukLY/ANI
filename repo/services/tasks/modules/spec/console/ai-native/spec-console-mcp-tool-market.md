# SPEC: Console mcp-tool-market

> Technical specification — Phase 3 规划详化  
> Source: `tasks/modules/prd/console/ai-native/prd-console-mcp-tool-market.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-mcp-tool-market-draft.md`  
> Revised: 2026-06-17

## 1. Summary

MCP 工具市场 catalog + install/uninstall（规划）；安装后工具进入 `listAgentTools` 可见集。合入 `services/v1.yaml` 前标注**规划**。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Planned Paths

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/svc/agent/mcp/tools` | `listMCPMarketToolPackages` | 200 | `scope:agent-mcp:read` |
| GET | `/api/v1/svc/agent/mcp/tools/{package_id}` | `getMCPMarketToolPackage` | 200 | read |
| POST | `.../install` | `installMCPMarketToolPackage` | 202/201 | write |
| POST | `.../uninstall` | `uninstallMCPMarketToolPackage` | 202 | write（3b） |

### 2.3 Planned Schemas

- `MCPMarketToolPackage`、`MCPMarketToolPackageListResponse`
- `MCPMarketInstallRequest`、`MCPMarketInstallation`

### 2.4 Tool ID 命名

安装成功后：`mcp:{package_id}/{tool_name}`，`AgentToolDefinition.provider=mcp`。

## 3. Page Scope

- 市场 catalog、包详情、安装/卸载、跳转 tool-permissions
- **Non-Goals**：BOSS 审核；`/integrations` 混 path

## 4. 创建前置条件

见主维护源 §创建前置条件。重复安装建议 `422 PACKAGE_ALREADY_INSTALLED`。

## 5. Handler 验收（合入 YAML 后）

```bash
curl .../agent/mcp/tools
curl -X POST .../agent/mcp/tools/$PKG/install -d '{"idempotency_key":"..."}'
curl .../agent/tools   # 安装后可见
```

## 6. 主维护源

- `docs/console-modules/ai-native/mcp-tool-market.md`
- 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-mcp-tool-market-draft.md`

OpenAPI **未合入** ≠ handler 已实现。
