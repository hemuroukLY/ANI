# ANI 安装向导（CLI/TUI）

## 页面定位

`ANI 安装向导` 是 BOSS **交付与安装** 域下的 **安装工具总入口** 专页：为交付工程师提供 `ani-installer` 单二进制 CLI/TUI 的命令索引、执行指引、现场状态只读展示，以及到各子流程（裸机/VM/attach/升级/卸载）的导航。

**本页 ≠ 实际安装执行面。** 全流程安装、前置检查、节点加入、K8s 接入、升级与卸载均在 **目标节点或跳板机** 上通过 `ani-installer` 命令完成（见 `ANI-main/ANI-07-部署工程设计.md` §2）；BOSS 页仅 **文档化 + 状态聚合展示**，**不得** 伪造 REST 安装 API。

一级权威源为 `ANI-main/ANI-07-部署工程设计.md` §2（命令与 11 步向导）。OpenAPI 安装工作流：**TODO-YAML: N/A** — Core `v1.yaml` **无** installer 相关 path。

Console **无** 对等页；租户侧不参与平台底座安装。

## 文档管理规则

- 本文是 `ANI 安装向导` 的主维护源；命令口径、向导步骤与验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-ani-installer.md`](../../tasks/modules/prd/boss/settings/prd-boss-ani-installer.md) 与 [`spec-boss-ani-installer.md`](../../tasks/modules/spec/boss/settings/spec-boss-ani-installer.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 ANI-07 §2，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)
- Console 对照：**无直接对应页**

## Core 层要求

- 安装/交付类能力 **不在 OpenAPI** — **TODO-YAML: N/A**
- **禁止** 自造 `/api/v1/boss/*` 或 `/api/v1/installer/*` 冒充已冻结契约
- BOSS 若展示组件健康，可 **只读参考** 现场 `ani-installer status` 输出或后续 Agent 上报 — **非** 当前 YAML 冻结路径
- 页面 RBAC：平台交付工程师 / SRE / 平台管理员；无 scope 时按 BOSS 平台角色
- 统一错误结构（经 Gateway 的 BOSS 读接口）：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- OpenAPI 已声明 ≠ handler 已实现（本页 REST **N/A**）

## 页面职责

- 提供 `ani-installer` **全命令索引**（preflight / new / join / attach-k8s / upgrade / status / uninstall / add-san）及 ANI-07 §2.2 十一向导步骤说明
- 链向专页：[`baremetal-deploy.md`](baremetal-deploy.md)、[`vm-deploy.md`](vm-deploy.md)、[`k8s-attach.md`](k8s-attach.md)、[`hardware-check.md`](hardware-check.md)、[`offline-package.md`](offline-package.md)、[`upgrade-backup.md`](upgrade-backup.md)
- 展示 **最近一次** 现场 preflight/status 摘要（若 Agent 或人工粘贴上报 — **待产品** P1 UI）
- 明确 **BOSS 不触发安装**；操作须 SSH/现场执行 CLI
- 安装完成后深链：[`internal-ca.md`](internal-ca.md)、[`ip-https.md`](ip-https.md)、[`acceptance-manual.md`](acceptance-manual.md)

## 页面结构

- 首屏至少包含：`命令卡片区`、`向导步骤时间线`、`现场状态摘要`、`相关模块导航`

```text
ANI 安装向导
├── 命令卡片（preflight / new / join / attach-k8s / upgrade / status / uninstall）
├── 11 步向导时间线（ANI-07 §2.2 · 只读说明）
├── 现场状态摘要（status 输出粘贴 / Agent — 待产品）
├── 部署路径入口（裸机 / VM / attach-k8s）
└── 交付链路（离线包 → 验收 → CA/HTTPS）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 / 模块 | 本页用法 |
|---|---|---|
| 产品 | ANI-07 §2 | **一级权威**：命令、向导、前置项 |
| CLI | `ani-installer *` | **实际操作面** |
| Core OpenAPI | — | **TODO-YAML: N/A** |
| BOSS UI | 本页 + 子模块详文 | 文档化 + 状态展示 |

### 关键边界

- **不得** 写「BOSS 一键安装 API 已实现」
- **不得** 将 Console 租户操作混入安装向导
- `ani-installer new` 内含 Step 2–6 硬件/驱动/证书 — 专页见 hardware-check、gpu-driver-install、internal-ca、ip-https
- attach 流程见 [`k8s-attach.md`](k8s-attach.md)；Core `k8s-clusters` **仅注册后只读参考**

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 命令卡片 | ANI-07 §2.1 | 命令名 + 一行说明 | 各专页 |
| 向导时间线 | ANI-07 §2.2 | Step 1–11 只读 | hardware-check / gpu-driver-install 等 |
| 状态摘要 | CLI `status` | 组件健康 JSON/文本 | — |
| 部署路径 | ANI-07 §1 矩阵 | 裸机/VM/attach | baremetal / vm / k8s-attach |
| 离线包提示 | ANI-07 + offline-package | U 盘/局域网分发 | offline-package |

## BOSS 与 Console 分工

| 维度 | BOSS ANI 安装向导 | Console |
|---|---|---|
| 用户 | 交付工程师、平台 SRE | 租户管理员、业务用户 |
| 安装执行 | **CLI/TUI 现场** | **无** |
| 状态查看 | BOSS 文档 + 可选 status 摘要 | **无** |
| API | **TODO-YAML: N/A** | **无** |
| 验收 | 链 acceptance-manual | **无** |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| — | installer 相关 REST | **不存在** | **TODO-YAML: N/A** |

| CLI 命令 | 权威源 | 说明 |
|---|---|---|
| `ani-installer preflight` | ANI-07 §2.1 | 仅检查，不安装 |
| `ani-installer new` | ANI-07 §2.1–§2.2 | 全新安装 TUI |
| `ani-installer join` | ANI-07 §2.1 | 节点加入 |
| `ani-installer attach-k8s` | ANI-07 §2.1、§4 | 接入已有 K8s |
| `ani-installer upgrade` | ANI-07 §2.1、§7 | 应用 `.anipatch` |
| `ani-installer status` | ANI-07 §2.1 | 组件健康 |
| `ani-installer uninstall` | ANI-07 §2.1 | 卸载（`--keep-data`） |
| `ani-installer add-san --ip` | ANI-07 §2.1、§9 | 证书 IP SAN |

| 能力 | 状态 |
|---|---|
| BOSS installer REST | **TODO-YAML: N/A** |
| BOSS status Agent 上报 | **待产品** P1 |
| 向导远程触发 | **不支持** |

### 页面目标字段（CLI / 状态摘要 · 非 REST）

| 字段 | 说明 |
|---|---|
| `command` | preflight / new / join / attach-k8s / upgrade / status / uninstall |
| `last_run_at` | 最近一次 CLI 执行时间（人工/Agent） |
| `exit_code` | 0 成功 / 非 0 失败 |
| `summary` | 单行结果摘要 |
| `components[]` | status：组件名 / 健康 / 版本 |
| `wizard_step` | new 向导当前步（1–11） |

## 字段级定义

### 展示字段（BOSS UI · CLI 映射）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `command` | CLI | — | 命令枚举 |
| `last_run_at` | 上报 | 可选 | ISO 8601 |
| `exit_code` | CLI | 可选 | integer |
| `summary` | CLI stdout | 可选 | 单行 |
| `components[].name` | status | 可选 | 如 ani-gateway |
| `components[].health` | status | 可选 | ok / degraded / down |
| `components[].version` | status | 可选 | semver |

### 查询字段（BOSS list 目标 · **待产品**）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `site_id` | query | 可选 | 客户现场 ID |
| `command` | query | 可选 | 筛选命令历史 |
| `limit` | query | 可选 | 分页 |

### 返回字段

本页 **无** Core REST 返回；status 以 CLI 文本/JSON 为准（见 §响应示例）。

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 命令卡片 + 向导时间线可浏览 | — |
| 无 status 上报 | 「暂无现场状态，请在节点执行 `ani-installer status`」 | **不** 伪造组件行 |
| preflight 失败摘要 | 红色告警条 + 链 hardware-check | exit_code ≠ 0 |
| 无权限态 | 403 提示 | 平台 RBAC |
| 安装进行中 | 黄色「现场安装中」+ 禁止 BOSS 触发 | CLI 独占 |
| upgrade | 链 upgrade-backup | Phase 2 能力标注 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `exit_code` | 进程退出码 | integer；0=成功 |
| `last_run_at` | 上报或粘贴时间 | ISO 8601 date-time |
| `components[].health` | ok / degraded / down | enum（CLI 约定） |
| `components[].version` | 组件 semver | string |
| `wizard_step` | 1–11 | integer |
| 安装时长（验收） | ANI-07 §12 | ≤ 2h（裸机 new） |

## 状态与能力口径

### CLI 执行状态（页面展示）

| 值 | 含义 | UI |
|---|---|---|
| `not_run` | 未执行 | 灰色 |
| `running` | 现场执行中（人工标记） | 蓝色 |
| `success` | exit_code=0 | 绿色 |
| `failed` | exit_code≠0 | 红色 |

### 组件健康（status）

| 值 | 含义 | UI |
|---|---|---|
| `ok` | 就绪 | 绿点 |
| `degraded` | 部分异常 | 橙点 |
| `down` | 不可用 | 红点 |

### 状态 × 操作可用性矩阵

| 操作 \ 场景 | BOSS 页 | 现场 CLI |
|---|---|---|
| 查看命令说明 | ✅ | ✅ |
| 执行 preflight/new/… | ❌ | ✅ |
| 粘贴 status 摘要 | ✅ P1 | ✅ |
| REST 触发安装 | ❌ **N/A** | ❌ |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的响应 |
|---|---|---|
| BOSS 读权限 | 交付/SRE/平台管理员 | `403 FORBIDDEN` |
| 离线包就位 | ANI-07 §10.3 T-1 | CLI 报错（见错误示例） |
| preflight 通过（new 前） | ANI-07 §2.3 | CLI 阻止安装 |
| 目标节点 SSH/现场 | 工程师可达 | 操作不可用（非 HTTP） |

## 操作可用性矩阵

| 操作 | 平台只读 | 交付工程师 | 平台管理员 | 说明 |
|---|---|---|---|---|
| 浏览命令与向导 | ✅ | ✅ | ✅ | 本文 |
| 查看 status 摘要 | ✅ **P1** | ✅ | ✅ | 粘贴/Agent |
| 现场执行 CLI | ❌ | ✅ | ✅ | **非 BOSS** |
| BOSS 触发安装 | ❌ | ❌ | ❌ | **N/A** |
| 跳转子模块 | ✅ | ✅ | ✅ | 裸机/VM/attach 等 |

## 删除前置校验与当前契约边界

本页 **无** DELETE REST 资源。**N/A**

- `ani-installer uninstall` 为 **CLI Destructive** 操作；BOSS **不得** 提供一键卸载 API（**TODO-YAML: N/A**）
- 卸载前置：备份（见 [`upgrade-backup.md`](upgrade-backup.md)）、客户确认 — CLI 交互确认

## 接口冻结规则

### Installer REST（**TODO-YAML: N/A**）

- Core `v1.yaml` **无** `/installer*` path
- **合入前不得** 写入「已冻结安装 API」

### CLI（ANI-07 §2.1 · **产品冻结**）

- 命令名与语义以 ANI-07 为准
- BOSS 正文 **仅引用**，不扩展未文档化 flag

### `GET /api/v1/k8s-clusters`（**非本页主路径**）

- attach 完成后 **可选** 只读展示 — 见 [`k8s-attach.md`](k8s-attach.md)
- **不得** 描述为 installer 触发创建的 REST

## 使用规则

- **不得** 伪造 REST 安装 path 或 operationId
- **不得** 在 BOSS 写「远程一键 new」 unless YAML 合入
- 安装执行 **必须** 经 `ani-installer` + 离线包（见 offline-package）
- preflight → new 顺序：**强制**（ANI-07 §10.3）
- attach 与 new **互斥** 路径 — 见 k8s-attach / baremetal / vm
- Console **无** 安装入口 — 不得链到租户页
- OpenAPI 已声明 ≠ handler 已实现（本页 REST N/A）

## 待补能力边界

- BOSS 现场 status Agent 上报 — **待产品** P1
- 命令执行历史 list API — **TODO-YAML: N/A** P1
- 远程 SSH 编排 — **Out of scope**
- 优先级：**P1**（文档与导航）；REST **N/A**

## 响应示例

### `ani-installer status` 成功（CLI · ANI-07 §2.1）

```text
ANI Platform Status (cluster: prod-01)
──────────────────────────────────────
ani-gateway          ok      v1.0.0
ani-auth-service     ok      v1.0.0
ani-model-service    ok      v1.0.0
ani-kb-service       ok      v1.0.0
ani-boss             ok      v1.0.0
dex                  ok      v1.0.0
──────────────────────────────────────
Overall: healthy
```

### `ani-installer preflight` 成功摘要（CLI）

```text
Preflight Report — PASS
OS: Ubuntu 22.04 · kernel 5.15 · x86_64
Memory: 128 GB · System disk free: 420 GB
NTP skew: 0.3s · Ports 6443/2379/2380/10250: free
GPU: NVIDIA A100 x8 (driver not installed — Step 3 will install)
Ready for: ani-installer new
```

## 错误示例

### 无 BOSS 交付页读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-installer-403-001"
}
```

### `ani-installer new` 前置未通过（CLI · 操作失败）

```text
ERROR: preflight not passed — kernel 5.10 < required 5.15
Run: ani-installer preflight
Exit code: 1
```

### 离线包缺失（CLI · 400 类业务失败）

```text
ERROR: offline bundle not found at /opt/ani/offline/images/
Expected layout: see ANI-07 offline package (images/ charts/ drivers/)
Exit code: 1
```

## 相关模块

- [`baremetal-deploy.md`](baremetal-deploy.md) · [`vm-deploy.md`](vm-deploy.md) · [`k8s-attach.md`](k8s-attach.md)
- [`hardware-check.md`](hardware-check.md) · [`gpu-driver-install.md`](gpu-driver-install.md)
- [`internal-ca.md`](internal-ca.md) · [`ip-https.md`](ip-https.md)
- [`offline-package.md`](offline-package.md) · [`acceptance-manual.md`](acceptance-manual.md) · [`upgrade-backup.md`](upgrade-backup.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] 权威源 ANI-07 §2；REST **TODO-YAML: N/A**
- [x] 含字段展示规则、字段口径与单位、状态与能力口径
- [x] 含响应示例（CLI status/preflight）与错误示例（403 + CLI 失败）
- [x] 删除前置校验 N/A 已声明
- [x] PRD/SPEC/HTML 与本文同步
