# BOSS 资源池与基础设施域

> **PRD/SPEC/HTML 同步**：ops 域 **8 模块**（P1×5 + P2 镜像仓库×3）已于 2026-06-17 与满配详文对齐。

| 模块 | 详文 | HTML ID | 优先级 |
|---|---|---|:---:|
| 平台资源池总览 | [`platform-resource-pool.md`](platform-resource-pool.md) | `boss.ops.pool` | P1 |
| GPU 资源池管理 | [`gpu-pool-management.md`](gpu-pool-management.md) | `boss.ops.gpu` | P1 |
| 节点状态 | [`node-status.md`](node-status.md) | `boss.ops.nodes` | P1 |
| 存储基础设施 | [`storage-infrastructure.md`](storage-infrastructure.md) | `boss.ops.storage` | P1 |
| 网络基础设施 | [`network-infrastructure.md`](network-infrastructure.md) | `boss.ops.network` | P1 |
| 镜像仓库 / 项目配额 | [`registry-project-quota.md`](registry-project-quota.md) | `boss.registry.quota` | P2 |
| 镜像仓库 / 漏洞扫描 | [`registry-vulnerability-scan.md`](registry-vulnerability-scan.md) | `boss.registry.vuln` | P2 |
| 镜像仓库 / 垃圾回收 | [`registry-garbage-collection.md`](registry-garbage-collection.md) | `boss.registry.gc` | P2 |

权威顺序：详文 > SPEC > PRD > HTML 摘要（`prototypes/ani-services-prototype-boss.html`）。

**P2 门禁（OpenAPI × 四材料）**：

```bash
python3 scripts/validate-boss-registry-p2-gate.py
```

报告输出：`docs/boss-modules/governance/registry-p2-gate-report.json`（`--json` 可选）。
