# PRD: Console agent-audit

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/ai-native/agent-audit.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-audit-draft.md`

## 目标

- 定义 **Agent 行为审计**与平台 `audit-log`、实例 `security-events` 的分工
- 输出只读 list API 草案（全局 + session scoped）
- 强制 list 响应 PII 脱敏（summary 级）

## 用户故事

- 作为租户管理员，我希望按会话/工具/时间检索 Agent 工具调用与策略命中
- 作为合规人员，我希望审计事件 append-only，且不含完整 prompt 明文
- 作为排障人员，我希望从 Session 详情直接进入该会话审计 Tab

## 范围

- GET audit-events（筛选 + 分页）
- GET sessions/{id}/audit-events
- Console 时间线 UI

## 非目标

- 平台级 Core audit list（见 `security/audit-log.md`）
- 审计导出（3b / BOSS，默认不开放）
- Console 写/删审计事件

## 成功标准

- [ ] 草案通过安全 + Services 评审
- [ ] redact 规则写入 YAML description 或实现说明
- [ ] 与 agent-session / tool-permission 事件类型对齐
