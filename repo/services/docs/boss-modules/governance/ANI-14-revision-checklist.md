# BOSS × ANI-14 对齐清单

> **用途**：跟踪 BOSS 文档链与 [`ANI-main/ANI-14-API对齐与开发工作流.md`](../../ANI-main/ANI-14-API对齐与开发工作流.md) 的对齐状态。  
> **Console 对照**：[`docs/console-modules/governance/ANI-14-revision-checklist.md`](../../console-modules/governance/ANI-14-revision-checklist.md)

## 文档链（Phase 1 产出）

| 产出物 | 路径 | 状态（2026-06-17） |
|---|---|---|
| BOSS 模块详文 | `docs/boss-modules/{域}/{module}.md`（排除 `governance/` 与各域 README） | ✅ 54 份模块叶子详文 |
| BOSS PRD | `tasks/modules/prd/prd-boss-*.md` | ✅ 54 套 |
| BOSS SPEC | `tasks/modules/spec/spec-boss-*.md` | ✅ 54 套 |
| BOSS HTML 摘要 | `prototypes/ani-services-prototype-boss.html` | ✅ **54 模块**已回填（2026-06-17） |
| BOSS 状态板 | `boss-document-status-board.md` | ✅ |

## Phase 进度

| Phase | BOSS 目标 | 当前状态 |
|---|---|---|
| **0** GAP | BOSS vs OpenAPI 差异报告 | ✅ 8 域摘要 + [boss-phase0-gap-index.md](./boss-phase0-gap-index.md)（2026-06-17） |
| **1** 文档对齐 | 详文 + PRD + SPEC + Core 层要求 | ✅ **54 模块**满配（Core）+ 已同步 + 阻塞项清零 |
| **2** YAML 扩充 | BOSS 平台级 TODO-YAML → OpenAPI | ⏳ 待产品/架构确认平台 API 归属 |
| **3** TASK | CORE/Services 任务清单追加 BOSS handler | ⏳ 依赖 Phase 2 |
| **4** Handler | 按 TASK 实现 | ⏳ |

## P0 约束（与 Console 相同）

- Core `/api/v1/*` · Services `/api/v1/svc/*`
- 统一实例 `kind=gpu_container|sandbox` 走 Core `/instances*`
- 禁止自造 schema / 路径 / 返回码
- BOSS 额外：**平台 RBAC**；租户级 YAML 路径不得直接冒充 BOSS 平台 API

## 下一步

1. Phase 2：平台 RBAC + 平台 API YAML 批次（见 [boss-phase0-gap-index.md](./boss-phase0-gap-index.md) Section D）
2. Phase 3 TASK → Phase 4 Handler
