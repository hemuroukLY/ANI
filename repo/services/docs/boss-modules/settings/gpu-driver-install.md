# GPU/NPU/DCU 驱动安装

## 页面定位

`GPU/NPU/DCU 驱动安装` 是 BOSS **交付与安装** 域下 **离线驱动包** 专页：对应 `ani-installer new` **Step 3** 与 ANI-07 §3.2（NVIDIA / 昇腾 / 海光 DCU；离线 rpm/deb，无外网）。

BOSS 文档化驱动映射表、随包组件与验证命令；**安装** 在 installer TUI/CLI 或离线包 `drivers/` 完成。

权威源：ANI-07 §3.2、§3.1 GPU 检测。REST：**TODO-YAML: N/A**。Console **无** 对等页。

## 文档管理规则

- 主维护源：本文
- PRD/SPEC：prd-boss-gpu-driver-install / spec-boss-gpu-driver-install
- 链：hardware-check、offline-package、ani-installer

## Core 层要求

- 驱动安装 REST — **N/A**
- 禁止 boss 路径自造
- K8s 层 Operator 生命周期 §3.2 末 — **安装后** nvidia-gpu-operator 等

## 页面职责

- 三厂商离线安装步骤与随包组件表（§3.2）
- 驱动版本映射 YAML 示例字段说明
- 验证：nvidia-smi / npu-smi / hy-smi
- 无 GPU → 纯 CPU 模式说明
- 链 offline-package images/ 与 drivers/

## 页面结构

```text
GPU 驱动安装
├── 检测 → 驱动映射（§3.2 YAML）
├── 厂商安装步骤表
├── 离线包 drivers/ 布局
├── 验证命令
└── K8s Operator 接管说明
```

## 数据来源与分层约束

| 层 | 来源 |
|---|---|
| ANI-07 §3.2 | 策略 |
| 离线包 | drivers/ rpm/deb |
| CLI | new Step 3 |
| Core gpu-inventory | **只读参考** P2 监控 |

## 页面区块与数据来源映射

| 区块 | 来源 |
|---|---|
| 映射表 | §3.2 yaml |
| NVIDIA 步骤 | §3.2 表 |
| 昇腾 CANN | §3.2 |
| 海光 DCU | §3.2 |
| 验证 | CLI |

## BOSS 与 Console 分工

| BOSS | Console |
|---|---|
| 交付文档 + 状态 | 租户用 GPU 资源 |

## 当前冻结事实

| REST | N/A |
| CLI | new Step 3 |

### 映射字段（installer 内置）

| 字段 | 说明 |
|---|---|
| `gpu_model` | A100/H100/910B 等 |
| `driver_version` | 535.x 等 |
| `cuda_version` | NVIDIA |
| `cann_version` | 昇腾 |

### 安装结果字段

| 字段 | 说明 |
|---|---|
| `vendor` | nvidia/ascend/hygon/skipped |
| `install_status` | success/failed/skipped |
| `verify_command` | smi 命令 |
| `verify_output` | 摘要 |

## 字段级定义

| 厂商 | 随包组件 | 验证 |
|---|---|---|
| NVIDIA | Driver + Container Toolkit | nvidia-smi |
| 昇腾 | Ascend-CANN + NPU 驱动 | npu-smi |
| 海光 | DCU 驱动包 | hy-smi |
| 无 GPU | — | skipped |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 已装正确版本 | 绿 + 跳过 |
| 需安装 | 蓝 Step 3 |
| 失败 | 红 + 日志路径 |
| 纯 CPU | 灰说明 |

## 字段口径与单位

| 字段 | 格式 |
|---|---|
| `driver_version` | vendor  semver |
| `cuda_version` | x.y |
| `vram_gb` | GB |

## 状态与能力口径

| install_status | UI |
|---|---|
| success | 绿 |
| failed | 红 |
| skipped | 灰 |
| pending | 蓝 |

## 创建前置条件

| 依赖 | 失败 |
|---|---|
| hardware-check GPU 检出 | 可 skip |
| 离线 drivers/ | CLI 失败 |
| BOSS 读 | 403 |

## 操作可用性矩阵

| 操作 | BOSS | CLI |
|---|---|---|
| 文档/映射 | ✅ | ✅ |
| Step 3 安装 | ❌ | ✅ |
| REST | ❌ N/A | ❌ |

## 删除前置校验与当前契约边界

**N/A**（驱动卸载非 BOSS REST；重装走 upgrade/offline）

## 接口冻结规则

REST **N/A**；映射表随 patch 更新（§3.2）。

## 使用规则

- **离线 only** — 禁止 installer 拉外网
- 已安装且版本匹配则跳过
- 失败不自动继续推理服务 — 验收 §12 GPU 项
- 不得写 driver install REST

## 待补能力边界

- BOSS 驱动版本 dashboard — P2
- gpu-inventory 对齐 — P2

## 响应示例

### Step 3 成功 NVIDIA（CLI）

```text
Step 3 Driver Install — SUCCESS
GPU: NVIDIA A100 x8
Installed: driver 535.104.05 · CUDA 12.2 · nvidia-container-toolkit
Verify: nvidia-smi — 8 devices visible
Next: Step 4 deployment mode
```

### 无 GPU 跳过

```text
Step 3 Driver Install — SKIPPED
No GPU/NPU/DCU detected (ANI-07 §3.1)
Platform will run CPU-only mode (KB Q&A available, inference rate-limited)
```

## 错误示例

### 403 BOSS

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-gpu-403-001"
}
```

### 驱动包缺失

```text
ERROR: offline package missing drivers/nvidia/535.104.05/
Expected under offline bundle drivers/ (ANI-07 §3.2)
Exit code: 1
```

## 相关模块

- [`hardware-check.md`](hardware-check.md) · [`offline-package.md`](offline-package.md)
- [`ani-installer.md`](ani-installer.md) · [`acceptance-manual.md`](acceptance-manual.md)

## 回填验收标准

- [x] 满配 22 章
- [x] §3.2 三厂商 + 离线约束
- [x] PRD/SPEC/HTML synced
