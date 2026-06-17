# 离线安装包

## 页面定位

`离线安装包` 是 BOSS **交付与安装** 域下 **离线介质** 专页：ANI 安装包与 `.anipatch` 结构（ANI-07 §7.1、§3.2 驱动随包、§10.2–§10.3 U 盘分发）。

含 `images/` OCI tar、`charts/`、`drivers/`、`migrations/`（patch）、`signature.sig` 等。**制作与校验** 在构建机/CLI；BOSS 文档化布局与版本 manifest。

REST：**TODO-YAML: N/A**。Console **无** 对等页。

## 文档管理规则

- 主维护源：本文
- PRD/SPEC：prd-boss-offline-package / spec-boss-offline-package
- 链：ani-installer、gpu-driver-install、upgrade-backup

## Core 层要求

- 离线包 upload REST — **Phase 2** patch 上传见 upgrade-backup（§7.3 BOSS 界面）
- v1.0 安装包 — **无** OpenAPI path
- 禁止 `/api/v1/boss/offline*`
- patch 验签 RSA-PSS（§7.1）

## 页面职责

- 文档化 **安装包** 目录：`images/`、`charts/`、`drivers/`、installer 二进制
- 文档化 **.anipatch** 结构（§7.1 manifest 字段）
- SOP T-1 U 盘/局域网分发（§10.3）
- 链 gpu-driver-install（drivers/）
- 链 upgrade-backup（patch 升级）

## 页面结构

```text
离线安装包
├── 安装包目录结构
├── .anipatch 结构（§7.1）
├── manifest.yaml 字段
├── 验签与 from_versions
├── 分发 SOP
└── 与 upgrade 关系
```

## 数据来源与分层约束

| 层 | 来源 |
|---|---|
| ANI-07 §7.1 | patch 格式 |
| ANI-07 §3.2 | drivers/ |
| §10.3 | 分发 |
| CLI | new/upgrade 消费包 |

## 页面区块与数据来源映射

| 区块 | 来源 |
|---|---|
| images/ | §7.1 + 安装包 |
| charts/ | Helm |
| drivers/ | §3.2 |
| manifest | §7.1 yaml |
| 验签 | signature.sig |

## BOSS 与 Console 分工

| BOSS | Console |
|---|---|
| 包版本文档、Phase2 上传 | 无 |

## 当前冻结事实

| 安装包 REST | N/A |
| patch 上传 BOSS UI | Phase 2 §7.3 |
| manifest 字段 | §7.1 冻结 |

### 安装包字段

| 路径 | 说明 |
|---|---|
| `ani-installer` | 二进制 |
| `images/*.tar` | OCI 镜像 |
| `charts/*.tgz` | Helm |
| `drivers/` | GPU/NPU/DCU 离线 |
| `crds/` | ANI CRD |

### .anipatch manifest（§7.1）

| 字段 | 说明 |
|---|---|
| `version` | 目标版本 |
| `from_versions` | 兼容源 |
| `min_supported_version` | 最低 |
| `rollback_supported` | bool |
| `components[]` | name/version |
| `migrations[]` | SQL 列表 |

## 字段级定义

### manifest.yaml 必填

| 字段 | 必填 |
|---|---|
| `version` | 是 |
| `from_versions` | 是 |
| `components` | 是 |
| `signature.sig` | 是（验签） |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 包版本 | 展示 manifest version |
| 验签失败 | 红 — 禁止 upgrade |
| 缺 images/ | 红 — new 不可 |
| Phase2 上传 | 链 upgrade-backup |

## 字段口径与单位

| 字段 | 格式 |
|---|---|
| `version` | vX.Y.Z |
| 镜像 tar | OCI layout |
| 包大小 | GB（UI） |

## 状态与能力口径

| package_status | 含义 |
|---|---|
| verified | 验签通过 |
| invalid_sig | 拒绝 |
| missing_component | 缺目录 |

## 创建前置条件

| 依赖 | 失败 |
|---|---|
| 离线包到场 | new 失败 |
| from_versions 匹配 | upgrade 拒绝 |
| BOSS 读 | 403 |

## 操作可用性矩阵

| 操作 | BOSS | 现场 |
|---|---|---|
| 文档/清单 | ✅ | ✅ |
| 复制 U 盘 | ❌ | ✅ |
| new/upgrade 消费 | CLI | CLI |
| patch 上传 UI | P2 | P2 |

## 删除前置校验与当前契约边界

**N/A**（包删除为文件系统操作，非 REST）

## 接口冻结规则

- 安装包 REST — **N/A**
- patch upload — Phase 2 upgrade-backup
- manifest schema — §7.1

## 使用规则

- installer **零外网** — 所有镜像/驱动在包内
- upgrade 须验签 + from_versions
- 不得伪造 offline upload API（Phase2 前）
- 与 upgrade-backup 同步 patch 口径

## 待补能力边界

- BOSS patch 拖拽上传 — Phase 2 P2
- 包 build CI — M5 工具链

## 响应示例

### manifest.yaml（§7.1 示例）

```yaml
version: v1.1.0
from_versions:
  - v1.0.x
min_supported_version: v1.0.0
rollback_supported: true
components:
  - name: ani-gateway
    version: v1.1.0
  - name: task-service
    version: v1.1.0
migrations:
  - 20260601_002_add_regions.sql
```

### 安装包目录（文档）

```text
ani-offline-v1.0.0/
├── ani-installer
├── images/
│   ├── ani-gateway-v1.0.0.tar
│   └── postgres-15.tar
├── charts/
│   └── ani-platform-v1.0.0.tgz
├── drivers/
│   ├── nvidia/535.104.05/
│   └── ascend/7.0.0/
└── crds/
    └── ani-crds.yaml
```

## 错误示例

### 403 BOSS

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-offline-403-001"
}
```

### upgrade 验签失败

```text
ERROR: signature verification failed for patch-v1.1.0.anipatch
RSA-PSS invalid — refuse upgrade (ANI-07 §7.1)
Exit code: 1
```

### 缺 images/

```text
ERROR: offline bundle not found at ./images/
Copy full offline package from USB (ANI-07 §10.3)
Exit code: 1
```

## 相关模块

- [`ani-installer.md`](ani-installer.md) · [`gpu-driver-install.md`](gpu-driver-install.md)
- [`upgrade-backup.md`](upgrade-backup.md) · [`baremetal-deploy.md`](baremetal-deploy.md)

## 回填验收标准

- [x] 满配 22 章
- [x] §7.1 + 安装包 images/drivers
- [x] PRD/SPEC/HTML synced
