# 安全与合规扩展索引

Console **安全 / 租户扩展** Phase 3 详化模块。Core secrets/encryption 部分已声明；§4 四项 **详化草案** ✅（2026-06-17）。

## Phase 3 详化模块（§4）

| 模块 | 详文 | 归属 | OpenAPI 草案 |
|---|---|---|---|
| API Key 审计 | [../tenant/api-key-audit.md](../tenant/api-key-audit.md) | Core Auth | [openapi-phase3-api-key-audit-draft.md](../openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md) |
| 网络安全策略 | [netsec-policy.md](./netsec-policy.md) | Core Networks | [openapi-phase3-netsec-policy-draft.md](../openapi-drafts/phase3/openapi-phase3-netsec-policy-draft.md) |
| 合规能力 | [compliance.md](./compliance.md) | Services 只读 | [openapi-phase3-compliance-draft.md](../openapi-drafts/phase3/openapi-phase3-compliance-draft.md) |
| 账单导出 | [../tenant/billing-export.md](../tenant/billing-export.md) | Services | [openapi-phase3-billing-export-draft.md](../openapi-drafts/phase3/openapi-phase3-billing-export-draft.md) |

## 已有 Core 能力（非本次详化）

| 模块 | 详文 | OpenAPI |
|---|---|---|
| 密钥管理 | [secrets-management.md](./secrets-management.md) | **部分** `/secrets*` |
| 国密加解密 | [crypto-sm.md](./crypto-sm.md) | **部分** `/encryption/*` |
| 审计日志 | [audit-log.md](./audit-log.md) | Tag 有；list **待补**（非 §4 本次范围） |

## 维护规则

- Core 项合入 `v1.yaml`；Services 项合入 `services/v1.yaml`
- compliance / billing-export：**Console 租户侧**；BOSS 运营不写为 Console write API
- TASK：Core §4 → `TASK-CORE-015`；Services §4 → `TASK-SVC-018`
- 验收：`tasks/phase3/acceptance/PHASE3-SECURITY-COMPLIANCE-ACCEPTANCE-RECORD.md`

## 评审建议顺序

1. api-key-audit + audit-log 事件模型对齐（字段映射）
2. netsec-policy vs security-group 分资源确认
3. compliance（只读）+ billing-export（202 export）与 BOSS 边界

## 相关

- 整域索引：`../openapi-drafts/phase3/openapi-phase3-domain-draft.md` §4
- 安全入口：`../tenant/security-identity-overview.md`
- Backlog：P2-21、P2-24、P2-26、P2-27
