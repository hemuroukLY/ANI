# ANI-14 修订 Diff 清单

> **用途**：将 `ANI-14-API对齐与开发工作流.md` 与 2026-06-16/17 Console 文档收口、OpenAPI 现状、统一实例架构决策对齐。  
> **权威对照**：`docs/console-modules/governance/console-document-status-board.md`、`../governance/module-delivery-workflow.md` §2.8/§2.10、`tasks/execution/CORE-TEAM-TASKS.md`  
> **最后更新**：2026-06-17

## 修订优先级

| 级别 | 含义 |
|---|---|
| **P0** | 不修订会导致 AI/团队按错误架构继续开发 |
| **P1** | 工作流可跑但产出路径/口径与仓内现状不一致 |
| **P2** | 示例、统计数字、次要路径约定 |

> **2026-06-17**：P0/P1 项已写入根目录 `ANI-14-API对齐与开发工作流.md` v1.1 与 `ANI-main/` 副本；P2 示例 TASK 块仍为 illustrative，随 Phase 2 YAML 落地再精修。

## P0 — 架构与路径（必须改）

### P0-1 §架构约束 第 3 条：GPU 容器 / Sandbox 归属

| | 内容 |
|---|---|
| **旧文** | GPU 容器 → `/api/v1/svc/gpu-containers`；Sandbox → `/api/v1/svc/sandboxes` |
| **新文** | 统一 Core 实例：`POST/GET /api/v1/instances*`，`kind=gpu_container` / `kind=sandbox`；`services/v1.yaml` 中 `/gpu-containers*`、`/sandboxes*` 已 **deprecated**，Phase 2/4 **不得**再新增或实现 |
| **依据** | `services/v1.yaml` info.description；`../governance/module-delivery-workflow.md` §2.8 |
| **状态** | ✅ 根目录 + ANI-main（2026-06-17） |

### P0-2 §架构约束：新增「统一实例」与「错误码」条款

| | 内容 |
|---|---|
| **新增** | §9 统一实例；§10 前置条件错误码（§2.10） |
| **状态** | ✅ |

### P0-3 Phase 3 / 产出物：TASK 清单路径

| | 内容 |
|---|---|
| **新文** | `tasks/execution/CORE-TEAM-TASKS.md` 等（非 `repo/docs/`） |
| **状态** | ✅ |

### P0-4 Phase 3 任务归属：deprecated 路径

| | 内容 |
|---|---|
| **新增** | 跳过 deprecated；GPU/Sandbox → Core instances TASK |
| **状态** | ✅ |

---

## P1 — 工作流输入/输出与 Phase 口径

### P1-1 §0 三个输入源 → 扩展为 Console 文档链

| | 内容 |
|---|---|
| **旧文** | 仅 `v1.yaml` + `services/v1.yaml` + `ani-services.html` |
| **新文** | 增加：`docs/console-modules/**`（主维护源）、`prototypes/ani-services-prototype-console.html`（Console HTML 摘要）；`ani-services.html` 仍为产品全量定义源 |
| **状态** | ✅ |

### P1-2 §0 / §E 统计数字

| | 内容 |
|---|---|
| **旧文** | Core 57+ 路径；Services 21 条全 501 |
| **新文** | 以 `schema-completion-tracker.md` 为准（Core ~88 路径 schema ~67%；Services 有效路径 ~29 handler ~28%）；**禁止写死过时数字** |
| **状态** | ✅ |

### P1-3 Phase 1 规则 4/5

| | 内容 |
|---|---|
| **旧文** | `### 前置条件`；示例 ``422 + code: MODEL_NOT_READY`` 定稿式 |
| **新文** | 章节名 `## 创建前置条件`；Services 模型未 ready 写「YAML 已举例 `MODEL_NOT_READY`」；其余写建议语义（§2.10） |
| **Deliverable** | Phase 1 产出包括 `docs/console-modules/{module}.md` 更新，不仅 HTML |
| **状态** | ✅ |

### P1-4 Phase 2：422 与 gpu-inventory

| | 内容 |
|---|---|
| **新增** | `POST /instances`、`POST .../lifecycle` 已挂 `422 PreconditionFailed`；`POST /k8s-clusters`、`searchVectorStore` 等待补 |
| **新增** | `GET /api/v1/gpu-inventory` Phase 2 **已声明**；详文 `gpu-inventory-ui.md`；handler 见 TASK-CORE-003 |
| **状态** | ✅（ANI-14 已写；YAML 待 Core 补） |

### P1-5 Phase 0 §D 示例

| | 内容 |
|---|---|
| **旧文** | `PATCH .../inference-services/{id}/scale` |
| **新文** | `PATCH /api/v1/svc/inference-services/{service_id}`（`patchInferenceService`，已在 services/v1.yaml） |
| **状态** | ✅ |

### P1-6 §文档产出物位置约定

| 产出物 | 旧位置 | 新位置 / 说明 |
|---|---|---|
| Console 模块详文 | （无） | `docs/console-modules/**/*.md` |
| Console 状态板 | （无） | `docs/console-modules/governance/console-document-status-board.md` |
| GAP REPORT | 仓库根目录 `GAP-REPORT-*.md` | **`GAP-REPORT-2026-06-17.md`** 为当前基线；2026-06-09 版部分结论已 supersede |
| TASK 三件套 | `repo/docs/` | **`tasks/`** |
| Core-ready 复审 | （无） | `docs/console-modules/governance/core-ready-review-checklist.md` |

---

## P2 — 示例与引用修正

### P2-1 Phase 3 示例 TASK-CORE-001

| | 内容 |
|---|---|
| **说明** | `gpu-inventory` 路径已在 Phase 2 声明；`gpu-inventory/{id}/metrics` 仍为 illustrative 示例 |
| **状态** | ⚠️ 示例保留；已加 listGPUInventory 前置说明 |

### P2-2 Phase 4 错误码列表

| | 内容 |
|---|---|
| **说明** | Phase 4 prompt 中 QUOTA_EXCEEDED 等仍为 handler 参考；文档层见 §2.10 |
| **状态** | ⚠️ 待 Phase 4 prompt 精修 |

### P2-3 版本与触发条件

| | 内容 |
|---|---|
| **新文** | v1.1 | 2026-06-17 |
| **状态** | ✅ |

### P2-4 双份 ANI-14 同步

| 路径 | 状态 |
|---|---|
| 根目录 + `ANI-main/` | ✅ |

---

## Phase 进度对照（2026-06-17）

| Phase | ANI-14 目标 | 当前状态 | 下一步 |
|---|---|---|---|
| **0** GAP | 客观差异报告 | ✅ **`GAP-REPORT-2026-06-17.md`** | 已完成 |
| **1** HTML/文档对齐 | 路径+前置条件+矩阵+ALIGNED/TODO-YAML | ✅ Console 21 模块 + HTML 摘要（2026-06-17） | 可进入 Phase 2 |
| **2** YAML 扩充 | TODO-YAML → OpenAPI | ✅ **`docs/console-modules/governance/YAML-EXPANSION-SUMMARY-2026-06-17.md`**（Core +19 / Svc +12） | Phase 4 handler |
| **3** TASK 分解 | 三份 TASK 清单 | ✅ **`tasks/`** 已与 Phase 2 同步（2026-06-17；Core 13 + Services 16 TASK） | Phase 4 按 TASK 实现 |
| **4** Handler | 按 TASK 实现 | ❌ 绝大部分待开始 | 按 `TASK-DEPENDENCY-MAP.md` Sprint 建议 |

## Core-ready 闭环检查（修订 ANI-14 之后）

1. 跑 `core-ready-review-checklist.md` 全模块串审  
2. `schema-completion-tracker.md` 与 YAML PR 同步  
3. TASK 完成打 ✅ → 更新状态板「开发可执行」列  
4. **禁止**在 ANI-14 / Console 文档中写「全量严格 Core-ready」直至上三项完成  

## 相关文件

- [module-delivery-workflow.md](./module-delivery-workflow.md)
- [console-document-status-board.md](./console-document-status-board.md)
- [schema-completion-tracker.md](./schema-completion-tracker.md)
- [../../tasks/execution/CORE-TEAM-TASKS.md](../../tasks/execution/CORE-TEAM-TASKS.md)
- [../../tasks/execution/SERVICES-TEAM-TASKS.md](../../tasks/execution/SERVICES-TEAM-TASKS.md)
