# SPRINT7-KICKOFF-A — Sprint 7 Core-only 代码开发入口切换

完成日期：2026-06-04
对应 Sprint：Sprint 7（2026-08-16 ~ 09-01）
验证结果：`make validate-doc-entrypoints` 与 `git diff --check` 均为 EXIT 0

## 实现了什么

完成 Sprint 7 kickoff 文档切换，将当前仓库执行范围明确收窄为 ANI Core：

1. `CORE-INSTALLER-A`：ani-installer 最小可验证闭环。
2. `CORE-OFFLINE-A`：Core 离线包 manifest、镜像/Helm chart/脚本清单与校验器。
3. `CORE-CLI-A`：复用 Core OpenAPI/Core SDK 的 `ani` Core CLI 最小资源覆盖。
4. `CORE-REGRESSION-A`：Sprint 5 real path 与 Sprint 6 Core local profile 的回归门禁入口。

本次只切换入口和执行边界，不新增 API、代码功能或 provider 能力；不代表 installer、离线包、CLI、真实 provider 或 production ready 已完成。

## 关键文件改动

| 文件 | 修改说明 |
|---|---|
| `repo/CURRENT-SPRINT.md` | 当前 Sprint 切换为 Sprint 7 Core-only 代码开发进行中，并列出 Core-only 执行队列 |
| `ANI-06-开发计划.md` | Section 零和 Sprint 7 章节切换为 Core-only，剥离 RAG/Console/Services 为本仓库当前任务 |
| `ANI-DOCS-INDEX.md` | 当前结论切换为 Sprint 7 Core-only 开发阶段 |
| `repo/development-records/README.md` | 新增 Sprint 7 kickoff 归档条目 |

## 完工标准

- [x] 不触碰 `repo/services/model-service/`、`repo/services/kb-service/`、`repo/ai/`、`repo/operators/`、`repo/frontends/`
- [x] 不新增 `M1-REAL-LAB-*` guard
- [x] 不把 local profile、installer、离线包、CLI 或真实 provider 标记为 production ready
- [x] Sprint 7 当前入口只指向 ANI Core 范围

## 备注

Sprint 7 后续 Feature batch 必须继续遵守 API-first、ports/adapters、TDD、最小实现和文档闭环。凡涉及真实底座组件或安装交付能力的批次，必须明确当前是 contract、local validation 还是 real-provider/live evidence。
