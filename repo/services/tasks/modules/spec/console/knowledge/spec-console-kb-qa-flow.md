# SPEC: Console 知识库问答流程

> Revised: 2026-06-17

## 1. Summary

Services KnowledgeBases 问答交互子流程。

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| POST | `/api/v1/svc/knowledge-bases/{kb_id}/query` | `queryKnowledgeBase` | `200 + KBQueryResponse` | `400`,`401`,`404` |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/query/stream` | `streamQueryKnowledgeBase` | `200` SSE | `400`,`401` |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/sessions` | `listKnowledgeBaseSessions` | `200 + KnowledgeBaseSessionListResponse` | `401`,`403`,`404` |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/citations` | `listKnowledgeBaseCitations` | `200 + KnowledgeBaseCitationListResponse` | `401`,`403`,`404` |

## 3. Page Scope

- 知识库选择、会话侧栏、对话区、citations

## 4. Non-Goals

- 不自造 ChatMessage schema

## 5. References

- `docs/console-modules/knowledge/kb-qa-flow.md`
