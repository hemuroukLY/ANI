# PRD: Console kb-meeting-intelligence

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/knowledge/kb-meeting-intelligence.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-meeting-intelligence-draft.md`

## 目标

- 定义会议 **ingest → 纪要/行动项 → 可选 KB 文档** 的产品边界
- 输出可评审的 Services OpenAPI 草案（ingest + list/get）
- 为 TASK-SVC-018 知识库智能子项提供合入 YAML 的前置材料

## 用户故事

- 作为编辑者，我希望上传会议录音或转写并自动生成纪要
- 作为团队成员，我希望在 Console 查看行动项并跳转关联 KB 文档
- 作为开发者，我希望 3a 支持 file/manual 源，Teams/Zoom 不阻塞首版

## 范围

- POST meetings/ingest、GET meetings list/detail；DELETE（3b）
- Console 会议列表、ingest 表单、纪要/行动项 UI
- RBAC：knowledge-bases read/write

## 非目标

- Teams/Zoom OAuth 完整集成（3b + integration-third-party）
- 合入 ANI-main

## 成功标准

- [ ] ingest 202 + AsyncTask；与 doc-intelligence 管线边界清晰
- [ ] 合入 YAML 后详文切换为冻结口径
