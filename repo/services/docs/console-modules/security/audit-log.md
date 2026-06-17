# 审计日志

## 页面定位

租户侧 **平台操作审计** 检索页（谁、何时、对何资源、何种动作）；Console 只读，合规导出见 BOSS/`compliance.md`。

## 文档管理规则

- 本文是 **审计日志** 的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 不得把规划路径或 handler stub 写成已实现
- PRD/SPEC 同步规则见 `docs/console-modules/governance/module-delivery-workflow.md` §3.5

## Core 层要求

v1.yaml Tag 组含 `Audit`，但 **paths 段无** `/audit*` 或 `listAuditEvents`。

可部分只读参考（非审计 list 替代）：
| GET | `/api/v1/observability/query` | PromQL 代理 — 不可当作审计事件 API |

<!-- TODO-YAML: GET /audit/events 或 /audit/logs -->

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 模块读/写权限 | 对应 RBAC scope | `403 FORBIDDEN` |

写操作 POST/PUT 的 `idempotency_key` 与 422 口径 — **待 YAML 冻结后按 ../governance/module-delivery-workflow.md §2.10 补充**。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员/运维 | 说明 |
|---|---|---|---|
| 检索审计事件 | 不可用 | 管理员 | TODO-YAML |
| 导出 | 不可用 | 管理员 | BOSS 边界 |

## 页面职责

- 占位 UI + 明确 YAML/OpenAPI 缺口（若适用）
- 跳转关联模块（见「相关模块」）
- 不把 BOSS/平台运营能力写入 Console 冻结契约

## 接口冻结规则

<!-- 上表为当前 YAML 已声明子集；模块完整能力仍待 TODO-YAML -->

## 待补边界

- Audit list API — **TODO-YAML**（Tag 已存在）
- 与 Agent 审计、API Key 审计分工 — Key 审计见 `tenant/api-key-audit.md`、`../openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md`

## 相关模块

- `tenant/api-key-audit.md`
- `ai-native/agent-audit.md`
- `security/compliance.md`

## 验收标准

- [ ] 路径与 OpenAPI 权威源一致（或明确 TODO-YAML / N/A）
- [ ] 正文不把 handler stub 写成已实现
- [ ] 含创建前置条件与操作可用性矩阵
- [ ] PRD/SPEC/HTML 摘要已同步
