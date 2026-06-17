# PRD: Console 知识库问答流程

> Revised: 2026-06-17  
> 详文：`docs/console-modules/knowledge/kb-qa-flow.md`

## 1. Overview

知识库问答完整 UI：POST query、GET SSE stream、sessions、citations。

## 2. Goals

- 非流式/流式双模式
- `session_id` 延续
- 引用来源展示

## 3. User Stories

### US-001: 提问

验收：POST 带 `question` + `idempotency_key`。

### US-002: 流式

验收：GET `/query/stream`；方法必须为 GET。

### US-003: 历史会话

验收：`listKnowledgeBaseSessions` 分页。

## 4. Functional Requirements

- FR-1: 流式 GET 与 YAML 一致
- FR-2: query 未声明 422 不得写定稿 code

## 5. Non-Goals

- 知识库 CRUD、权限编辑

## 6. References

- 详文、SPEC
