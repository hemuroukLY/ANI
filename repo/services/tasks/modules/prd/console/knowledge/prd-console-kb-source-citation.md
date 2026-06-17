# PRD: Console kb-source-citation

详文：`docs/console-modules/knowledge/kb-source-citation.md` · Phase 2 YAML 已声明项文档收口（2026-06-17）。

## 目标

- 为 **Services / KnowledgeBases / 来源引用** 提供可维护的产品边界说明
- 明确 `GET .../citations` 已冻结、引用导出等待补

## 用户故事（规划）

- 作为知识库读者，我希望查看问答召回的来源引用列表，以便追溯答案依据
- 作为编辑者，我希望从引用列表跳转文档详情，以便核对原文

## 范围

- 引用列表只读页、与 kb-qa-flow citations 口径一致 — 见主维护源
- 不修改 ANI-main OpenAPI（由 Services 团队按 TASK-SVC-014 推进）

## 非目标

- 不自造 Citation 扩展 schema
- 引用单条 GET、导出、向量 debug 视图
