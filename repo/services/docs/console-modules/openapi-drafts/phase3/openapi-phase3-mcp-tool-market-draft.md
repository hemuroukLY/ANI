# Phase 3 — MCP Tool Market OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`ai-native/mcp-tool-market.md`  
> **TASK**：`TASK-SVC-018` 子项  
> **依赖**：安装后工具进入 `agent-tool-permission.md` 的 `listAgentTools` 可见集

---

## 1. 与工具权限 / 第三方集成的边界

| 维度 | MCP 市场（本草案） | 工具权限 | 第三方集成 |
|---|---|---|---|
| 模块 | `mcp-tool-market.md` | `agent-tool-permission.md` | `integration-third-party.md` |
| 路径 | `/agent/mcp/tools*` | `/agent/tool-permissions` | `/integrations` |
| 职责 | **发现、订阅、安装** MCP 包 | allow/deny **已安装**工具 | 通用 CRM/Webhook 集成 |
| 写操作 | install/uninstall | PUT grants | createIntegration |

安装成功后，工具 `id` 形如 `mcp:{package_id}/{tool_name}`，与 `AgentToolDefinition.provider=mcp` 对齐。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [MCPTools]
x-ani-rbac-scope:
  read:  scope:agent-mcp:read
  write: scope:agent-mcp:write
```

---

## 3. Schemas（草案）

### MCPMarketToolPackage

```yaml
MCPMarketToolPackage:
  type: object
  required: [id, name, version, publisher]
  properties:
    id:          { type: string, description: 市场包 ID }
    name:        { type: string }
    version:     { type: string }
    publisher:   { type: string, enum: [platform, partner, community] }
    description: { type: string, nullable: true }
    risk_level:  { type: string, enum: [low, medium, high], nullable: true }
    tool_count:  { type: integer, minimum: 0 }
    installed:   { type: boolean, description: 当前租户是否已安装 }
```

### MCPMarketToolPackageListResponse

```yaml
MCPMarketToolPackageListResponse:
  type: object
  required: [items]
  properties:
    items:       { type: array, items: { $ref: '#/components/schemas/MCPMarketToolPackage' } }
    next_cursor: { type: string, nullable: true }
```

### MCPMarketInstallRequest

```yaml
MCPMarketInstallRequest:
  type: object
  required: [idempotency_key]
  properties:
    idempotency_key: { type: string, format: uuid }
    version:         { type: string, nullable: true, description: 空则安装 latest }
```

### MCPMarketInstallation

```yaml
MCPMarketInstallation:
  type: object
  required: [package_id, status, installed_at]
  properties:
    package_id:   { type: string }
    status:       { type: string, enum: [installing, active, failed, uninstalling] }
    installed_at: { type: string, format: date-time }
    tool_ids:     { type: array, items: { type: string }, description: 展开后的 AgentToolDefinition.id 列表 }
```

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/svc/agent/mcp/tools`

- operationId: `listMCPMarketToolPackages`
- Query: `publisher`、`risk_level`、`installed`（bool）、`limit`、`cursor`
- 成功：`200 + MCPMarketToolPackageListResponse`
- 说明：市场 **catalog**；与 `GET /agent/tools`（已安装+内置）分工 — 市场页用本 op，权限页用 tools

### 4.2 `GET /api/v1/svc/agent/mcp/tools/{package_id}`

- operationId: `getMCPMarketToolPackage`
- 成功：`200 + MCPMarketToolPackage` + 可选 `tools[]` 明细
- 错误：`404`

### 4.3 `POST /api/v1/svc/agent/mcp/tools/{package_id}/install`

- operationId: `installMCPMarketToolPackage`
- Body: `MCPMarketInstallRequest`
- 成功：`202 + MCPMarketInstallation` 或 `201`（同步安装完成时）
- 异步时：`202 + AsyncTask`，轮询 Core `GET /tasks/{id}`
- 422：包已安装 / 版本不存在 — 建议 `PACKAGE_ALREADY_INSTALLED`（待 YAML 举例）

### 4.4 `POST /api/v1/svc/agent/mcp/tools/{package_id}/uninstall`（3b）

- operationId: `uninstallMCPMarketToolPackage`
- Body: `{ idempotency_key }`
- 成功：`202` 或 `200`
- 卸载后相关 `AgentToolGrant` 建议自动 revoke（Services 实现说明，非 Console 契约）

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | Agent 管理员 | 说明 |
|---|---|---|---|
| 浏览市场 | 可用 | 可用 | GET catalog |
| 查看包详情 | 可用 | 可用 | GET package |
| 安装 | 不可用 | 可用 | POST install |
| 卸载 | 不可用 | 可用 | POST uninstall（3b） |
| 配置 allow/deny | 跳转 | 跳转 | tool-permissions |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/agent/mcp/tools?limit=20"

curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'"}' \
  "$BASE/api/v1/svc/agent/mcp/tools/$PACKAGE_ID/install"
```

---

## 7. 评审检查清单

- [ ] 与 `/integrations` 无 path 冲突
- [ ] 安装后 `tool_ids` 与 `listAgentTools` 一致
- [ ] install 必填 `idempotency_key`
- [ ] 合入后更新 `mcp-tool-market.md`

---

## 相关文件

- `docs/console-modules/ai-native/mcp-tool-market.md`
- `docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md`
