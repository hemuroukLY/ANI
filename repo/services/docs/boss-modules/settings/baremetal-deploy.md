# 裸机全新部署

## 页面定位

`裸机全新部署` 是 BOSS **交付与安装** 域下 **裸物理服务器** 路径专页：客户自购服务器、无任何基础设施时，经 `ani-installer new` 在上游原生 Kubernetes 上完成 ANI 平台安装（ANI-07 §1 第一行、§2.2 Step 4「全新原生 K8s」）。

**本页 ≠ 安装执行 API。** 实际操作在目标裸机上执行 CLI；BOSS 提供参数口径、前置清单、进度文档化与验收跳转。

一级权威源：`ANI-main/ANI-07-部署工程设计.md` §1、§2.2–§2.3、§10.3、§12。OpenAPI：**TODO-YAML: N/A**。

Console **无** 对等页。

## 文档管理规则

- 本文是 `裸机全新部署` 的主维护源
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-baremetal-deploy.md`](../../tasks/modules/prd/boss/settings/prd-boss-baremetal-deploy.md) 与 [`spec-boss-baremetal-deploy.md`](../../tasks/modules/spec/boss/settings/spec-boss-baremetal-deploy.md) 为辅助材料
- 冲突时以 ANI-07 + 本文为准，回写 PRD/SPEC
- 流程：[`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)
- 枢纽：[`ani-installer.md`](ani-installer.md)

## Core 层要求

- 裸机部署 **无** Core REST 契约 — **TODO-YAML: N/A**
- **禁止** `/api/v1/boss/*` 自造 path
- 网络/存储口径：KubeOVN + Rook-Ceph/本地 SSD（ANI-07 §1）
- 错误经 Gateway 时：`code` / `message` / `request_id`
- 底座：**上游原生 Kubernetes**，非发行版绑定（ANI-07 §1 关键规则）

## 页面职责

- 文档化裸机路径 **前置条件**（OS/内核/内存/磁盘/NTP/端口 — ANI-07 §2.3）
- 链向 preflight、硬件检测、GPU 驱动、内部 CA、纯 IP HTTPS、离线包
- 展示 ANI-07 §2.2 **new 向导** 在裸机场景的 Step 4–10 要点（网络 CIDR、存储裸盘、组件选择）
- 提供标准 SOP 时间线（T-3 ~ T+1 — §10.3）与验收跳转 [`acceptance-manual.md`](acceptance-manual.md)
- 区分 VM 路径 — [`vm-deploy.md`](vm-deploy.md)

## 页面结构

```text
裸机全新部署
├── 场景说明（ANI-07 部署矩阵 · 裸物理）
├── 前置清单（§2.3 表格 + 现场勘察 Excel）
├── 执行步骤（preflight → new · 11 步高亮）
├── 网络/存储/证书要点（链 ip-https / internal-ca）
└── SOP 时间线与验收入口
```

## 数据来源与分层约束

| 层 | 来源 | 本页用法 |
|---|---|---|
| 产品 | ANI-07 §1、§2、§10、§12 | **一级权威** |
| CLI | `ani-installer preflight` / `new` | 执行面 |
| Core | — | **N/A** |
| 离线包 | ANI-07 + offline-package | 介质分发 |

### 关键边界

- 裸机 **必须** 走 `new`，**不是** `attach-k8s`
- 数据盘裸盘警告：无盘则无 Rook-Ceph（§2.3）
- GPU 场景链 gpu-driver-install；纯 CPU 可跳过 Step 3

## 页面区块与数据来源映射

| 区块 | 来源 | 说明 | 跳转 |
|---|---|---|---|
| 前置清单 | §2.3 | 阻止/警告项 | hardware-check |
| 向导步骤 | §2.2 | Step 1–11 | ani-installer |
| 网络存储 | §2.2 Step 5–7 | CIDR/裸盘/SC | ip-https / internal-ca |
| SOP | §10.3 | T-3~T+1 | acceptance-manual |
| 离线包 | §10.2 | U 盘分发 | offline-package |

## BOSS 与 Console 分工

| 维度 | BOSS 裸机部署 | Console |
|---|---|---|
| 执行 | 现场 CLI | 无 |
| 文档/SOP | ✅ | 无 |
| 租户资源 | 不涉及 | 安装后使用 |

## 当前冻结事实

| REST | 状态 |
|---|---|
| 裸机 deploy API | **TODO-YAML: N/A** |

| CLI | 说明 |
|---|---|
| `ani-installer preflight` | 裸机首要步骤 |
| `ani-installer new` | 裸机主路径 |

| 验收 | ANI-07 §12 |
|---|---|
| 裸机全新安装 | ≤ 2h，无报错 |

### 页面字段（文档/清单 · 非 API）

| 字段 | 说明 |
|---|---|
| `os_type` | Ubuntu 22.04 / UOS / 麒麟 |
| `kernel_version` | ≥ 5.15 |
| `cpu_arch` | x86_64 / ARM64 |
| `memory_gb` | 管理节点 ≥16，工作 ≥32 |
| `system_disk_free_gb` | ≥ 100 |
| `raw_data_disks` | 裸盘块数/容量 |
| `ntp_skew_sec` | ≤ 5 |
| `ports_free` | 6443/2379/2380/10250 |
| `gpu_detected` | 有/无 → 驱动链 |

## 字段级定义

### 前置检查字段（ANI-07 §2.3）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `os_type` | preflight | 建议 | 不满足警告 |
| `kernel_version` | preflight | 必填 | <5.15 阻止 |
| `cpu_arch` | preflight | 必填 | x86_64/ARM64 |
| `memory_gb` | preflight | 建议 | 低于阈值警告 |
| `system_disk_free_gb` | preflight | 必填 | <100 阻止 |
| `raw_data_disks` | preflight | 建议 | 0 块警告 |
| `ntp_skew_sec` | preflight | 必填 | >5 阻止 |
| `ports_free` | preflight | 必填 | 占用阻止 |

### 向导配置字段（Step 5–8 · 展示）

| 字段 | 说明 |
|---|---|
| `mgmt_ip` | 集群管理 IP |
| `pod_cidr` | Pod 网段 |
| `service_cidr` | Service 网段 |
| `access_host` | 域名或纯 IP |
| `storage_mode` | Rook 裸盘 / 已有 SC |
| `edition` | 基础版 / 完整版 |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 内核/架构失败 | 红色阻止条 + preflight 命令 |
| 无数据裸盘 | 黄色警告 + Rook 说明 |
| 有 GPU 未装驱动 | 蓝色提示 → gpu-driver-install |
| 纯 IP 访问 | 链 ip-https + internal-ca |
| 403 | 无 BOSS 读权限 |

## 字段口径与单位

| 字段 | 口径 | 格式 |
|---|---|---|
| `kernel_version` | ≥ 5.15 | semver |
| `memory_gb` | GiB 十进制展示 | GB |
| `system_disk_free_gb` | 可用空间 | GB |
| `ntp_skew_sec` | 节点间时差 | 秒；≤5 |
| `raw_data_disks` | 未挂载裸盘 | 块数 + GB |
| 安装时长 | §12 通过标准 | ≤ 2 小时 |

## 状态与能力口径

### Preflight 结果

| 值 | UI |
|---|---|
| PASS | 绿 |
| WARN | 黄（可继续） |
| BLOCK | 红（不可 new） |

### 部署阶段（SOP）

| 阶段 | 含义 |
|---|---|
| T-3 | 前置清单确认 |
| T-1 | 离线包到场 |
| T+0 AM | preflight |
| T+0 PM | new 执行 |
| T+1 | 验收 + CA 信任 |

## 创建前置条件

| 依赖 | 要求 | 失败 |
|---|---|---|
| BOSS 读权限 | 交付角色 | 403 |
| 离线包 | offline-package 就位 | CLI 失败 |
| preflight PASS | 无 BLOCK 项 | new 拒绝 |
| 工程师现场 | SSH/带外 | 不可执行 |

## 操作可用性矩阵

| 操作 | 只读 | 交付工程师 | 平台管理员 |
|---|---|---|---|
| 查看 SOP/清单 | ✅ | ✅ | ✅ |
| preflight/new | ❌ BOSS | ✅ CLI | ✅ CLI |
| BOSS 触发部署 | ❌ | ❌ | ❌ |

## 删除前置校验与当前契约边界

**N/A**（无 REST DELETE）。卸载见 `ani-installer uninstall` — [`ani-installer.md`](ani-installer.md)。

## 接口冻结规则

- REST：**TODO-YAML: N/A** — 不得写已实现 deploy API
- CLI：ANI-07 §2.1 冻结命令语义

## 使用规则

- 裸机 **只用** `new`，不用 `attach-k8s`
- 必须先 `preflight`（§10.3）
- 纯 IP 须 Step 6 证书 — internal-ca / ip-https
- 禁止伪造 BOSS deploy REST
- VM 场景见 vm-deploy — 不混用文案

## 待补能力边界

- BOSS 部署进度实时上报 — **待产品** P1
- REST orchestration — **TODO-YAML: N/A**
- 优先级：**P1**

## 响应示例

### preflight 裸机 PASS（CLI 节选）

```text
Preflight Report — PASS (bare metal)
OS: Ubuntu 22.04.4 LTS · kernel 5.15.0-91 · x86_64
Memory: 256 GB · System disk free: 512 GB
Raw disks: 4 x 3.84 TB NVMe (unmounted)
NTP skew: 0.8s · Required ports: free
Recommendation: ani-installer new (bare metal, greenfield K8s)
```

### new 完成（CLI Step 11）

```text
Installation complete.
Console: https://10.20.30.40/
BOSS:    https://10.20.30.40/boss/
Initial admin: admin / ******** (change on first login)
Download root CA: /var/lib/ani/certs/ani-root-ca.crt
Next: distribute CA (see BOSS → internal-ca) · run acceptance tests
```

## 错误示例

### BOSS 403

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-bm-403-001"
}
```

### preflight BLOCK — 内核不足

```text
BLOCK: kernel 5.10.0 — minimum 5.15 required (ANI-07 §2.3)
Installation blocked. Upgrade kernel or use supported OS image.
Exit code: 1
```

### new 失败 — 端口占用

```text
Step 10 FAILED: port 6443 already in use (pid 8821)
Rollback initiated for step kubernetes-bootstrap
Exit code: 1
```

## 相关模块

- [`ani-installer.md`](ani-installer.md) · [`hardware-check.md`](hardware-check.md) · [`offline-package.md`](offline-package.md)
- [`gpu-driver-install.md`](gpu-driver-install.md) · [`internal-ca.md`](internal-ca.md) · [`ip-https.md`](ip-https.md)
- [`acceptance-manual.md`](acceptance-manual.md) · [`vm-deploy.md`](vm-deploy.md)（对照）

## 回填验收标准

- [x] 满配 22 章齐全
- [x] ANI-07 §1/§2/§10/§12 引用；REST N/A
- [x] 字段展示/口径/状态 + CLI 示例 + 403/操作失败
- [x] 删除 N/A
- [x] PRD/SPEC/HTML synced
