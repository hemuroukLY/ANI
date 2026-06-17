# BOSS 安全审计与合规域

> **PRD/SPEC/HTML 同步**：audit 域 **4 模块**已于 2026-06-17 与满配详文对齐。

| 模块 | 详文 | HTML ID | 优先级 |
|---|---|---|:---:|
| 平台审计日志 | [`platform-audit-log.md`](platform-audit-log.md) | `boss.audit.ops` | P1 |
| API Key 审计 | [`platform-apikey-audit.md`](platform-apikey-audit.md) | `boss.audit.apikey` | P1 |
| 推理调用审计 | [`platform-inference-audit.md`](platform-inference-audit.md) | `boss.audit.inference` | P1 |
| 合规导出与取证 | [`compliance-export.md`](compliance-export.md) | `boss.audit.export` | P2 |

Phase 0 GAP 摘要：[`governance/boss-phase0-gap-audit.md`](../governance/boss-phase0-gap-audit.md)

**门禁**：

```bash
python3 scripts/validate-boss-audit-gate.py
```

权威顺序：详文 > SPEC > PRD > HTML 摘要。
