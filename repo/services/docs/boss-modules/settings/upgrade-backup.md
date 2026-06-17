# 升级与备份

## 页面定位

`升级与备份` 是 BOSS **交付与安装** 域下 **`ani-installer upgrade`** 与 **平台备份/ Patch 升级管理** 专页（ANI-07 §7）。

- **CLI**：`ani-installer upgrade` 应用 `.anipatch`（§7.1）
- **BOSS UI**：§7.3 升级管理界面（版本面板、包上传、Dry-Run、SSE 进度）— **Phase 2**
- **备份**：运维手册口径（§10.2《运维手册》）+ 升级前备份 SOP — **无** 独立 REST **N/A**

权威源：ANI-07 §7、§10.2。OpenAPI upgrade：**TODO-YAML: N/A**（`k8s-clusters/{id}/upgrade` 为 **不同域** 租户集群升级，**非** 本页平台 patch）。

Console **无** 对等页。优先级：**P2**。

## 文档管理规则

- 主维护源：本文
- PRD/SPEC：prd-boss-upgrade-backup / spec-boss-upgrade-backup
- 链：offline-package、ani-installer

## Core 层要求

- 平台 patch BOSS API — **Phase 2 TODO-YAML**
- `POST /api/v1/k8s-clusters/{cluster_id}/upgrade` — **租户 K8s 集群升级** — **非** 本页 ANIPatch 升级 — **只读边界说明**
- 禁止 `/api/v1/boss/upgrade*`
- patch 验签 RSA-PSS（§7.1）
- ANIPatch CRD §7.2 — Operator 执行

## 页面职责

- 文档化 `.anipatch` 与 `ani-installer upgrade` CLI
- §7.3 BOSS 目标能力（上传/Dry-Run/SSE/历史）— 标注 Phase 2
- 升级前 **备份清单**：etcd/PostgreSQL/MinIO/配置 — 运维 SOP
- 链 offline-package manifest / from_versions
- 区分 Phase 2 vs v1.0（§11 Patch 为 Phase 2）

## 页面结构

```text
升级与备份
├── 当前版本面板（§7.3 · P2）
├── Patch 上传 / Dry-Run（P2）
├── ani-installer upgrade CLI
├── 备份清单（升级前）
├── Operator 流程（§7.2）
└── 升级历史（P2）
```

## 数据来源与分层约束

| 层 | 来源 | 用法 |
|---|---|---|
| CLI | upgrade | v1 可选 CLI |
| ANI-07 §7 | patch/operator | 权威 |
| BOSS §7.3 | UI | **P2** |
| k8s-clusters upgrade API | v1.yaml | **非本页** |

### 关键边界

- **不得** 把 `k8s-clusters/upgrade` 写成平台 patch 升级
- v1.0 **不** 阻塞 Patch（§12）
- uninstall `--keep-data` 见 ani-installer — 非备份替代

## 页面区块与数据来源映射

| 区块 | 来源 |
|---|---|
| 版本面板 | status + §7.3 P2 |
| 包上传 | §7.3 P2 |
| CLI | §2.1 upgrade |
| 备份 SOP | §10.2 运维手册 |
| Operator | §7.2 CRD |

## BOSS 与 Console 分工

| BOSS | Console |
|---|---|
| 平台升级/备份 | 无 |

## 当前冻结事实

| 方法 | 路径 | 说明 |
|---|---|---|
| — | platform patch upload | **TODO-YAML P2** |
| CLI | `ani-installer upgrade` | §2.1 |
| CRD | `ANIPatch` | §7.2 Operator |

| k8s-clusters upgrade | **非本页** |
|---|---|
| POST `/api/v1/k8s-clusters/{cluster_id}/upgrade` | 租户集群版本升级 |

### §7.3 BOSS 目标字段（P2）

| 字段 | 说明 |
|---|---|
| `current_version` | 平台版本 |
| `component_versions[]` | 各组件 |
| `last_upgrade_at` | 时间 |
| `patch_upload_status` | uploading/ready |
| `dry_run_result` | 变更预览 |
| `upgrade_progress[]` | SSE 组件状态 |
| `history[]` | 版本/操作人/结果 |

### 备份清单字段（SOP · 非 API）

| 项 | 说明 |
|---|---|
| `etcd_snapshot` | 若 installer 部署 etcd |
| `postgres_dump` | PG 逻辑备份 |
| `minio_bucket` | 模型/品牌对象 |
| `helm_values` | ani-system values |
| `certs_backup` | 可选 |

## 字段级定义

### manifest（offline-package）

见 [`offline-package.md`](offline-package.md) `version` / `from_versions` / `rollback_supported`。

### Operator spec（§7.2）

| 字段 | 说明 |
|---|---|
| `version` | 目标 |
| `source` | uploaded |
| `dryRun` | bool |
| `maintenanceWindow` | 可选 |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| Phase 2 未就绪 | 「Patch BOSS UI Phase 2」 |
| 验签失败 | 红 — 禁止升级 |
| from_versions 不匹配 | 红 |
| 备份未完成 | 黄警告 gate |
| 403 | 无升级权限 P2 |

## 字段口径与单位

| 字段 | 格式 |
|---|---|
| `version` | vX.Y.Z |
| SSE 进度 | % 或 step |
| 备份大小 | GB |

## 状态与能力口径

### upgrade_progress（§7.3 P2）

| status | UI |
|---|---|
| pending | 灰 |
| running | 蓝 |
| success | 绿 |
| failed | 红 |
| rolled_back | 橙 |

### 备份检查

| checked | 允许 upgrade |
|---|---|
| all | 是（建议） |
| partial | 需确认 |

## 创建前置条件

| 依赖 | 失败 |
|---|---|
| from_versions 匹配 | CLI/Operator 拒绝 |
| signature 有效 | 验签失败 |
| 备份完成 | SOP 警告 |
| BOSS upgrade RBAC | 403 P2 |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 平台管理员 |
|---|---|---|---|
| 查看文档/备份清单 | ✅ | ✅ | ✅ |
| CLI upgrade | ❌ BOSS | ✅ | ✅ |
| BOSS 上传 patch | ⏳ P2 | ⏳ P2 | ⏳ P2 |
| Dry-Run / 一键升级 | ⏳ P2 | ⏳ P2 | ⏳ P2 |
| 查看历史 | ⏳ P2 | ⏳ P2 | ⏳ P2 |

## 删除前置校验与当前契约边界

**N/A**（无 DELETE patch 资源）。回滚 §7.2 step 9 — Operator **非** REST DELETE。

## 接口冻结规则

### Platform patch BOSS（**TODO-YAML P2**）

<!-- TODO-YAML: upload / dry-run / upgrade / history SSE -->

- **合入前不得写已实现**

### `ani-installer upgrade`（ANI-07 §2.1）

- 消费 `.anipatch` — 见 offline-package
- 验签 + compatibility

### `POST .../k8s-clusters/{id}/upgrade`（**边界 · 非本页**）

- operationId 见 v1.yaml
- **不得** 作为平台 patch 升级契约

## 使用规则

- Phase 2 — 非 v1.0 阻塞
- upgrade **前** 执行备份 SOP
- CLI 与 BOSS UI **并存** Phase 2 后
- 不得混淆 k8s-clusters upgrade
- 禁止伪造 boss upgrade REST

## 待补能力边界

- §7.3 全功能 — **P2**
- ANIPatch Operator — Phase 2
- 自动备份 Job — P2
- REST upload/SSE — **TODO-YAML P2**

## 响应示例

### ani-installer upgrade 成功（CLI）

```text
Applying patch-v1.1.0.anipatch ...
Signature: OK · from_version v1.0.3 compatible
Pre-upgrade script: OK
Upgrading ani-gateway ... OK
Upgrading task-service ... OK
Migrations: 2/2 applied
Post-upgrade: OK
Health check: all components ready
Upgrade complete: v1.0.3 → v1.1.0
```

### BOSS 版本面板（§7.3 目标 · P2）

```json
{
  "current_version": "v1.0.3",
  "components": [
    { "name": "ani-gateway", "version": "v1.0.3" },
    { "name": "ani-boss", "version": "v1.0.3" }
  ],
  "last_upgrade_at": "2026-05-01T02:00:00Z",
  "patch_ready": false
}
```

## 错误示例

### 403 BOSS upgrade（P2 目标）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-upg-403-001"
}
```

### from_versions 不兼容（CLI）

```text
ERROR: current v1.0.3 not in from_versions [v1.1.x] for patch-v1.2.0.anipatch
See ANI-12 for major upgrade path
Exit code: 1
```

### 备份未确认（运维 gate）

```text
WARN: pre-upgrade backup checklist incomplete (postgres_dump missing)
Continue? [y/N] — BOSS P2 will block until backup tasks checked
```

## 相关模块

- [`offline-package.md`](offline-package.md) · [`ani-installer.md`](ani-installer.md)
- [`acceptance-manual.md`](acceptance-manual.md)（ACC-P2-01）

## 回填验收标准

- [x] 满配 22 章
- [x] §7 CLI/Operator/BOSS §7.3 + 备份 SOP
- [x] k8s-clusters upgrade 边界
- [x] P2 标注；PRD/SPEC/HTML synced
