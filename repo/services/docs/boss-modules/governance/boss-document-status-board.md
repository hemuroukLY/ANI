# BOSS 文档总控状态表

## 说明

本表是 BOSS 产品文档的唯一总控视图，对齐 ANI-14 Phase 1 产出约定。

每个模块必须同时具备：**主维护文档 + PRD + SPEC**。

- `可用（Phase 1 结构）`：方案 A 章节门禁（文档管理规则、Core 层要求、冻结事实、前置条件、操作矩阵、接口冻结规则）
- `满配（Core）`：通过 [`boss-full-depth-checklist.md`](boss-full-depth-checklist.md) 全部章节 + 阻塞项清零
- `待收口`：基线已生成，尚未达到 Phase 1 结构门禁

## 模块状态总表

| 模块 | 详文 | PRD | SPEC | HTML | 优先级 | 主维护源 |
|---|---|---|---|---|---|---|
| 运营总览 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `overview/operations-overview.md` |
| 资源池与容量态势 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `overview/resource-capacity-trend.md` |
| GPU 资源池态势 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `overview/gpu-pool-trend.md` |
| AI 服务运营态势 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `overview/inference-ops-trend.md` |
| 知识库运营态势 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `overview/kb-ops-trend.md` |
| 平台告警与待处理 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `overview/platform-alerts-pending.md` |
| 租户列表 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `tenant/tenant-list.md` |
| 配额策略 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `tenant/tenant-quota-policy.md` |
| 租户管理员 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `tenant/tenant-admin.md` |
| 租户计费与用量 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `tenant/tenant-usage-billing.md` |
| 平台资源池总览 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `ops/platform-resource-pool.md` |
| GPU 资源池管理 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `ops/gpu-pool-management.md` |
| 节点状态 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `ops/node-status.md` |
| 存储基础设施 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `ops/storage-infrastructure.md` |
| 网络基础设施 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `ops/network-infrastructure.md` |
| 镜像仓库 / 项目配额 | 满配（Core） | ✅ | ✅ | ✅ | P2 | `ops/registry-project-quota.md` |
| 镜像仓库 / 漏洞扫描 | 满配（Core） | ✅ | ✅ | ✅ | P2 | `ops/registry-vulnerability-scan.md` |
| 镜像仓库 / 垃圾回收 | 满配（Core） | ✅ | ✅ | ✅ | P2 | `ops/registry-garbage-collection.md` |
| 平台健康 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `health/platform-health.md` |
| GPU 监控 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `health/gpu-monitoring.md` |
| 推理监控 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `health/inference-monitoring.md` |
| 知识库监控 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `health/kb-monitoring.md` |
| 日志 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `health/platform-logs.md` |
| Trace | 满配（Core） | ✅ | ✅ | ✅ | P1 | `health/platform-trace.md` |
| 告警规则 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `health/alert-rules.md` |
| 运维 Skills | 满配（Core） | ✅ | ✅ | ✅ | P2 | `health/maint-skills.md` |
| 任务历史 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `health/job-history.md` |
| 故障处理 | 满配（Core） | ✅ | ✅ | ✅ | P2 | `health/incident-handling.md` |
| 平台 GPU-Hours | 满配（Core） | ✅ | ✅ | ✅ | P1 | `metering/platform-gpu-hours.md` |
| 平台 CPU-Hours | 满配（Core） | ✅ | ✅ | ✅ | P1 | `metering/platform-cpu-hours.md` |
| 平台 Memory-GBHours | 满配（Core） | ✅ | ✅ | ✅ | P1 | `metering/platform-memory-gbhours.md` |
| 平台 Storage-GBDays | 满配（Core） | ✅ | ✅ | ✅ | P1 | `metering/platform-storage-gbdays.md` |
| 平台 Input Tokens | 满配（Core） | ✅ | ✅ | ✅ | P1 | `metering/platform-input-tokens.md` |
| 平台 Output Tokens | 满配（Core） | ✅ | ✅ | ✅ | P1 | `metering/platform-output-tokens.md` |
| 平台 KB Queries | 满配（Core） | ✅ | ✅ | ✅ | P1 | `metering/platform-kb-queries.md` |
| 平台审计日志 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `audit/platform-audit-log.md` |
| API Key 审计 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `audit/platform-apikey-audit.md` |
| 推理调用审计 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `audit/platform-inference-audit.md` |
| 合规导出与取证 | 满配（Core） | ✅ | ✅ | ✅ | P2 | `audit/compliance-export.md` |
| ani-installer | 满配（Core） | ✅ | ✅ | ✅ | P1 | `settings/ani-installer.md` |
| 裸机部署 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `settings/baremetal-deploy.md` |
| VM 部署 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `settings/vm-deploy.md` |
| 已有 K8s 接入 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `settings/k8s-attach.md` |
| 硬件检测 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `settings/hardware-check.md` |
| GPU 驱动安装 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `settings/gpu-driver-install.md` |
| 内部 CA | 满配（Core） | ✅ | ✅ | ✅ | P1 | `settings/internal-ca.md` |
| 纯 IP HTTPS | 满配（Core） | ✅ | ✅ | ✅ | P1 | `settings/ip-https.md` |
| 离线安装包 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `settings/offline-package.md` |
| 交付验收手册 | 满配（Core） | ✅ | ✅ | ✅ | P1 | `settings/acceptance-manual.md` |
| 品牌配置 | 满配（Core） | ✅ | ✅ | ✅ | P2 | `settings/brand-config.md` |
| 升级与备份 | 满配（Core） | ✅ | ✅ | ✅ | P2 | `settings/upgrade-backup.md` |
| 运维 Webhook | 满配（Core） | ✅ | ✅ | ✅ | P1 | `integration/ops-webhook.md` |
| 企业通知集成 | 满配（Core） | ✅ | ✅ | ✅ | P2 | `integration/enterprise-notification.md` |
| 运营系统对接 | 满配（Core） | ✅ | ✅ | ✅ | P2 | `integration/ops-system-integration.md` |

## 权威源对齐说明

- Core：`ANI-main/repo/api/openapi/v1.yaml`
- Services：`ANI-main/repo/api/openapi/services/v1.yaml`
- 交付类：TODO-YAML: N/A
- BOSS 平台级 API：多数 **TODO-YAML**（租户级已声明路径仅作只读参考）
