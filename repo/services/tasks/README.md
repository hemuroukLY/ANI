# tasks — 任务与辅助材料

实现任务、阶段验收、PRD/SPEC 辅助文档。**模块详文**在 `docs/console-modules/`。

## 目录

| 目录 | 内容 | 何时打开 |
|---|---|---|
| [execution/](execution/) | `CORE/SERVICES-TEAM-TASKS`、依赖图、handler 指南 | 排期、开 PR、验 handler |
| [phase2/](phase2/) | Phase 2 handler 文档验收 | Core/Services Phase 2 |
| [phase3/](phase3/) | **整域联评**、分域验收记录 | Phase 3 YAML 评审 |
| [modules/prd/](modules/prd/) | 86 × PRD | 产品边界辅助 |
| [modules/spec/](modules/spec/) | 86 × SPEC | 技术边界辅助 |
| [support/](support/) | 算力/存储/Services 长文辅助 | 运维、异常处理 |

## 常用入口

- Phase 3 联评：[`phase3/PHASE3-JOINT-REVIEW-AGENDA.md`](phase3/PHASE3-JOINT-REVIEW-AGENDA.md)
- Services 总任务：[`execution/SERVICES-TEAM-TASKS.md`](execution/SERVICES-TEAM-TASKS.md)
- Core 总任务：[`execution/CORE-TEAM-TASKS.md`](execution/CORE-TEAM-TASKS.md)

## 命名约定

- PRD：`modules/prd/prd-console-{module-slug}.md`
- SPEC：`modules/spec/spec-console-{module-slug}.md`
- 与详文 `{module-slug}` 与 `docs/console-modules/**/{module}.md` 对应

## 新增文件放哪

| 类型 | 路径 |
|---|---|
| 新模块 PRD/SPEC | `modules/prd/`、`modules/spec/` |
| 阶段验收 | `phase3/acceptance/` 或 `phase2/` |
| TEAM 任务更新 | `execution/` |
