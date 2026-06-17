# 虚拟机部署

## 页面定位

`虚拟机部署` 是 BOSS **交付与安装** 域下 **虚拟机（KVM/VMware/Hyper-V）** 路径专页：客户在虚拟化平台新建 VM 后，仍通过 `ani-installer new` 安装 **上游原生 Kubernetes** + ANI 平台（ANI-07 §1 第二行）。

与裸机差异：底层为虚拟化磁盘/vSAN，网络/存储选项相同（KubeOVN + Rook-Ceph 或 vSAN）。**执行面仍为 CLI**，非 REST。

权威源：ANI-07 §1、§2、§10.3、§12。OpenAPI：**TODO-YAML: N/A**。Console **无** 对等页。

## 文档管理规则

- 本文为主维护源
- PRD/SPEC：`prd-boss-vm-deploy.md` / `spec-boss-vm-deploy.md`
- 冲突：ANI-07 + 本文优先
- 流程：module-delivery-workflow + boss-full-depth-checklist
- 对照：[`baremetal-deploy.md`](baremetal-deploy.md)（物理机）、[`ani-installer.md`](ani-installer.md)

## Core 层要求

- VM 部署 REST — **TODO-YAML: N/A**
- 禁止 `/api/v1/boss/*`
- 虚拟化层 **不在** installer 管理范围 — 客户 IT 先建 VM
- Gateway 错误结构统一

## 页面职责

- VM 场景 **额外前置**：vCPU/RAM/磁盘模板、虚拟网卡、时间同步（NTP 透传）
- 文档化 `preflight` → `new` 与裸机相同向导（§2.2）
- 存储：Rook on vSAN/虚拟磁盘 vs 客户 StorageClass 说明
- SOP 与验收链 acceptance-manual
- 云主机 ECS 场景备注（ANI-07 §1 第三行 — 参数类似，网络或 CNI-Hybrid）

## 页面结构

```text
虚拟机部署
├── 虚拟化前置（VM 规格 / 网卡 / 磁盘）
├── 与裸机差异对照表
├── preflight → new 步骤
├── 存储与网络（KubeOVN / vSAN）
└── SOP + 验收
```

## 数据来源与分层约束

| 层 | 来源 | 用法 |
|---|---|---|
| ANI-07 §1 | 部署矩阵 VM 行 | 权威 |
| CLI | preflight / new | 执行 |
| 虚拟化 | 客户 KVM/VMware/Hyper-V | **外部** |
| Core REST | — | N/A |

### 关键边界

- installer **不** 创建 VM — 只装 K8s+ANI
- attach-k8s **不适用** 本页（已有 K8s 见 k8s-attach）
- GPU 直通/ vGPU 须 preflight 可见 — hardware-check

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| VM 规格建议 | ANI-07 §2.3 + 产品 | hardware-check |
| 向导 | §2.2 | ani-installer |
| 存储 | §1 VM 行 | offline-package |
| 证书/IP | Step 5–6 | ip-https / internal-ca |

## BOSS 与 Console 分工

| 维度 | BOSS VM 部署 | Console |
|---|---|---|
| VM 创建 | 客户 IT | 无 |
| ani-installer | 交付工程师 CLI | 无 |
| 文档 | ✅ | 无 |

## 当前冻结事实

| REST | **TODO-YAML: N/A** |
| CLI | `preflight`, `new`（同裸机） |
| 验收 §12 | VM 路径同裸机时长标准 |

### 页面字段

| 字段 | 说明 |
|---|---|
| `hypervisor` | kvm / vmware / hyper-v |
| `vcpu` | 建议 ≥ 8 管理 / ≥ 16 工作 |
| `ram_gb` | 同 §2.3 内存阈值 |
| `disk_type` | thin / thick / vSAN |
| `virtio_net` | 推荐驱动 |
| `gpu_passthrough` | 可选 |

## 字段级定义

### VM 前置字段

| 字段 | 必填 | 说明 |
|---|---|---|
| `hypervisor` | 建议 | 影响磁盘/网卡说明 |
| `vcpu` | 建议 | 低于阈值警告 |
| `ram_gb` | 必填检查 | preflight 内存项 |
| `system_disk_gb` | 必填 | ≥100 可用 |
| `data_disk_gb` | 建议 | Rook 用 |
| `mgmt_network` | 必填 | 管理网段 |

### 向导配置字段（Step 5–8 · 展示）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `mgmt_ip` | Step 5 | 必填 | 集群管理 IP；VM 常为虚拟网卡 IP |
| `pod_cidr` | Step 5 | 必填 | Pod 网段；与 hypervisor 虚拟网络规划一致 |
| `service_cidr` | Step 5 | 必填 | Service 网段；不得与 VM 管理网重叠 |
| `access_host` | Step 6 | 必填 | 域名或纯 IP；纯 IP 须链 ip-https |
| `storage_mode` | Step 7 | 必填 | `rook_raw`（VM 数据盘）/ `existing_sc` |
| `edition` | Step 8 | 必填 | `basic` / `full`；影响组件清单 |
| `cert_mode` | Step 6 | 建议 | `auto_ca` / `import`；Dex 须 HTTPS |
| `offline_bundle_path` | Step 4 | 条件 | 离线部署时必填；指向 `images/` 根 |

### VM 与裸机差异字段

| 字段 | VM 特有 | 说明 |
|---|---|---|
| `virtio_net` | 是 | 推荐 VirtIO 网卡驱动 |
| `gpu_passthrough` | 是 | PCI 直通时链 gpu-driver-install |
| `thin_provision` | 是 | thin 磁盘可能影响 Rook 性能 — 黄警告 |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| VM 内存低于工作节点 32GB | 黄警告 |
| 无 VirtIO | 黄警告（性能） |
| GPU 直通未识别 | 链 gpu-driver-install |
| 403 | BOSS 无权限 |

## 字段口径与单位

| 字段 | 格式 |
|---|---|
| `ram_gb` | GB（十进制 UI） |
| `vcpu` | 整数核 |
| `disk_type` | enum 文案 |
| NTP skew | 秒；≤5 |

## 状态与能力口径

Preflight / SOP 阶段同 baremetal-deploy。

## 创建前置条件

| 依赖 | 失败 |
|---|---|
| VM 已创建并开机 | 无法 SSH |
| 离线包 | CLI 错误 |
| preflight PASS | new 阻止 |
| BOSS 读权限 | 403 |

## 操作可用性矩阵

| 操作 | BOSS | CLI 现场 |
|---|---|---|
| 文档/SOP | ✅ | ✅ |
| preflight/new | ❌ | ✅ |
| REST deploy | ❌ N/A | ❌ |

## 删除前置校验与当前契约边界

**N/A**。卸载 CLI only。

## 接口冻结规则

REST **N/A**；CLI ANI-07 §2.1。

## 使用规则

- VM 仍用 `new` 非 attach
- 虚拟化准备 **客户负责**
- 不得写 VM provision REST
- 裸机文案复用时须改 hypervisor 上下文

## 待补能力边界

- BOSS VM 清单导入 — P2
- REST — **TODO-YAML: N/A** P1

## 响应示例

### preflight on VM（CLI）

```text
Preflight Report — PASS (VM on VMware ESXi 8.0)
Guest: Ubuntu 22.04 · 16 vCPU · 64 GB RAM · 200 GB system disk
VirtIO SCSI · vmxnet3 · NTP skew: 1.2s
No GPU passthrough — CPU-only mode available
Ready: ani-installer new
```

## 错误示例

### 403 BOSS

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-vm-403-001"
}
```

### VM 磁盘不足

```text
BLOCK: system disk free 62 GB < required 100 GB (ANI-07 §2.3)
Expand VM disk or add volume before ani-installer new
Exit code: 1
```

## 相关模块

- [`baremetal-deploy.md`](baremetal-deploy.md) · [`ani-installer.md`](ani-installer.md) · [`hardware-check.md`](hardware-check.md)
- [`offline-package.md`](offline-package.md) · [`acceptance-manual.md`](acceptance-manual.md)

## 回填验收标准

- [x] 满配章节齐全
- [x] ANI-07 VM 行 + REST N/A
- [x] 展示/口径/状态 + 示例
- [x] PRD/SPEC/HTML synced
