# Phase 3 安全/合规扩展 ×4 — 文档验收记录

> **生成日期**：2026-06-17  
> **用途**：Phase 3 §4 四模块 **文档层验收**。  
> **约束**：本阶段**不修改** `ANI-main/**`。  
> **索引**：`docs/console-modules/security/README.md`

---

## 总表

| 模块 | 详文 | 归属 | 草案 | 文档验收 | YAML | Handler |
|---|---|---|---|---|---|---|
| API Key 审计 | `tenant/api-key-audit.md` | Core Auth | `../openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md` | ✅ | ☐ | ☐ |
| 网络安全策略 | `security/netsec-policy.md` | Core Networks | `../openapi-drafts/phase3/openapi-phase3-netsec-policy-draft.md` | ✅ | ☐ | ☐ |
| 合规能力 | `security/compliance.md` | Services 只读 | `../openapi-drafts/phase3/openapi-phase3-compliance-draft.md` | ✅ | ☐ | ☐ |
| 账单导出 | `tenant/billing-export.md` | Services | `../openapi-drafts/phase3/openapi-phase3-billing-export-draft.md` | ✅ | ☐ | ☐ |

**文档验收 ✅ 定义**

- Core 2 项：v1.yaml；Services 2 项：services/v1.yaml
- compliance：**无 write**；billing-export：**202 AsyncTask**，不含 BOSS 结算金额
- netsec：`/networks/policies` ≠ `/networks/security-groups`

---

> **整域联评**：`tasks/phase3/PHASE3-JOINT-REVIEW-AGENDA.md`（场次 C：安全/合规 + 合入计划）

---

## 联评议程（域内摘要 · 完整版见整域议程 §4.4）

1. **审计分工**：api-key-audit vs audit-log vs agent-audit
2. **网络策略**：netsec policy_type 枚举与 SG  UI 分流
3. **BOSS 边界**：compliance 报告生成、billing 发票
4. **合入批次**：Core（auth audit + netsec）与 Services（compliance + billing）可分 PR

---

## curl 预览（合入后）

```bash
curl .../auth/api-keys/$KEY/audit-events?limit=20
curl .../networks/policies?vpc_id=$VPC
curl .../svc/compliance/summary
curl -X POST .../svc/billing/exports -d '{"idempotency_key":"...","export_type":"usage_csv",...}'
```

---

## 签核

- [x] ×4 详文 + PRD/SPEC + 草案（2026-06-17）
- [x] `security/README.md`、`../openapi-drafts/phase3/openapi-phase3-domain-draft.md` §4 更新
- [ ] 架构联评通过 — **整域**：`PHASE3-JOINT-REVIEW-AGENDA.md` §7
- [ ] YAML 合入（v1 + services）
- [ ] Handler 验通（ANI-main）

---

## 相关

- `tasks/TASK-CORE-015.md` — Core 子项
- `tasks/execution/SERVICES-TEAM-TASKS.md` — TASK-SVC-018
- `tasks/phase3/acceptance/PHASE3-MODEL-ENHANCEMENT-ACCEPTANCE-RECORD.md`
