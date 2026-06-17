# 硬件检测（Preflight）

## 页面定位

`硬件检测` 是 BOSS **交付与安装** 域下 **`ani-installer preflight`** 与 **new 向导 Step 1–2** 专页：OS/内核/CPU/内存/磁盘/网卡/NTP/端口/GPU 枚举（ANI-07 §2.3、§3.1）。

BOSS 展示 preflight **报告摘要**与检查项口径；**执行** 在节点上 `ani-installer preflight` 或 new Step 1–2。

权威源：ANI-07 §2.3、§3.1。REST：**TODO-YAML: N/A**。Console **无** 对等页。

## 文档管理规则

- 主维护源：本文
- PRD/SPEC：prd-boss-hardware-check / spec-boss-hardware-check
- 链：[`ani-installer.md`](ani-installer.md)、[`gpu-driver-install.md`](gpu-driver-install.md)

## Core 层要求

- 硬件检测 REST — **N/A**
- 禁止 `/api/v1/boss/*`
- GPU 检测逻辑 §3.1 优先级：nvidia-smi → npu-smi → hy-smi → rocm-smi

## 页面职责

- 完整 §2.3 检查表（阻止/警告）
- §3.1 CPU/GPU/存储/网络检测字段说明
- preflight 报告字段与 UI 展示规则
- 失败项 → 修复指引（不阻止 vs 阻止）
- GPU 检出 → gpu-driver-install

## 页面结构

```text
硬件检测
├── §2.3 前置检查表
├── §3.1 检测维度（CPU/GPU/磁盘/网卡）
├── preflight CLI 用法
├── 报告摘要展示区
└── 与 new Step 1–2 对应关系
```

## 数据来源与分层约束

| 层 | 来源 |
|---|---|
| ANI-07 §2.3 | 阈值 |
| ANI-07 §3.1 | 检测算法 |
| CLI | preflight 输出 |

## 页面区块与数据来源映射

| 区块 | 来源 |
|---|---|
| 检查表 | §2.3 |
| GPU 检测链 | §3.1 |
| NUMA/IB | §3.1 网络 |
| 报告 | CLI |

## BOSS 与 Console 分工

| BOSS | Console |
|---|---|
| 文档 + 报告摘要 | 无 |

## 当前冻结事实

| REST | N/A |
| CLI | `ani-installer preflight` |

### 检测输出字段

| 字段 | 说明 |
|---|---|
| `os_type` | OS 识别 |
| `kernel_version` | ≥5.15 |
| `cpu_arch` | x86_64/ARM64 |
| `cpu_brand` | Intel/AMD/鲲鹏等 |
| `numa_nodes` | NUMA 拓扑 |
| `memory_total_gb` | 内存 |
| `system_disk_free_gb` | 系统盘 |
| `raw_disks[]` | 裸盘列表 |
| `ntp_skew_sec` | NTP |
| `ports_status` | 6443 等 |
| `gpu_devices[]` | vendor/model/driver |
| `nics[]` | 速率/MTU |

## 字段级定义

### §2.3 检查项

| 检查项 | 最低要求 | 不满足 |
|---|---|---|
| OS | Ubuntu22.04/UOS/麒麟 | 警告 |
| 内核 | ≥5.15 | **阻止** |
| 架构 | x86_64/ARM64 | **阻止** |
| 内存 | 管理≥16G 工作≥32G | 警告 |
| 系统盘 | ≥100GB | **阻止** |
| 数据裸盘 | ≥1 块 | 警告 |
| NTP | ≤5s | **阻止** |
| 端口 | 未占用 | **阻止** |
| SELinux | permissive 建议 | 警告 |

### GPU 检测字段（§3.1）

| 字段 | 说明 |
|---|---|
| `vendor` | nvidia/ascend/hygon/amd/none |
| `model` | 型号 |
| `driver_installed` | bool |
| `vram_gb` | 显存 |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| BLOCK | 红 + 不可 new |
| WARN | 黄 + 可继续说明 |
| 无 GPU | 灰「纯 CPU 模式」 |
| GPU 无驱动 | 蓝 → gpu-driver-install |
| MTU<1500 | 黄 Jumbo 建议 |

## 字段口径与单位

| 字段 | 单位 |
|---|---|
| `memory_total_gb` | GB |
| `ntp_skew_sec` | 秒 |
| `vram_gb` | GB |
| `raw_disks[].size_gb` | GB |
| `nics[].speed_gbps` | Gbps |

## 状态与能力口径

### 单项检查结果

| 值 | UI |
|---|---|
| PASS | 绿 |
| WARN | 黄 |
| BLOCK | 红 |

### 报告总评

| 值 | 可 new |
|---|---|
| PASS | 是 |
| PASS_WITH_WARNINGS | 是 |
| FAIL | 否 |

## 创建前置条件

| 依赖 | 失败 |
|---|---|
| 节点 SSH | 无法 preflight |
| BOSS 读 | 403 |

## 操作可用性矩阵

| 操作 | BOSS | CLI |
|---|---|---|
| 查看检查表 | ✅ | ✅ |
| 运行 preflight | ❌ | ✅ |
| 粘贴报告 | ✅ P1 | ✅ |

## 删除前置校验与当前契约边界

**N/A**

## 接口冻结规则

REST **N/A**；检测项 ANI-07 §2.3/§3.1 冻结。

## 使用规则

- new **前** 须 preflight（§10.3）
- 不得伪造 hardware REST
- GPU 驱动安装见专页 — preflight 只检测

## 待补能力边界

- BOSS 报告持久化 — P1
- REST — N/A

## 响应示例

### preflight 完整报告（CLI 节选）

```text
Preflight Report — PASS_WITH_WARNINGS
──────────────────────────────────────
[PASS] kernel 5.15.0-91 · x86_64
[WARN] memory 24 GB < recommended 32 GB for worker role
[PASS] system disk free 180 GB
[WARN] no raw data disks — Rook-Ceph unavailable
[PASS] NTP skew 2.1s · ports free
GPU Detection:
  NVIDIA A100-SXM4-80GB x2 · driver: not installed
  → Step 3 / gpu-driver-install required for inference
Network: eth0 25Gbps MTU 1500 · no RDMA NIC detected
Overall: ready for ani-installer new with warnings acknowledged
```

## 错误示例

### 403

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-hw-403-001"
}
```

### BLOCK — NTP

```text
BLOCK: NTP skew 12.4s > maximum 5s (ANI-07 §2.3)
Fix chrony/ntp on all nodes before installation
Exit code: 1
```

## 相关模块

- [`ani-installer.md`](ani-installer.md) · [`gpu-driver-install.md`](gpu-driver-install.md)
- [`baremetal-deploy.md`](baremetal-deploy.md) · [`vm-deploy.md`](vm-deploy.md)

## 回填验收标准

- [x] 满配 22 章
- [x] §2.3/§3.1 完整
- [x] PRD/SPEC/HTML synced
