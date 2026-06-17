# PRD: Console kb-video-intelligence

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/knowledge/kb-video-intelligence.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-video-intelligence-draft.md`

## 目标

- 定义视频 **ingest → 转写/章节/OCR** 并写入 KB 的产品边界
- 输出可评审的 Services OpenAPI 草案（ingest + list/get）
- 为 TASK-SVC-018 知识库智能子项提供合入 YAML 的前置材料

## 用户故事

- 作为编辑者，我希望上传培训视频并自动得到章节与转写
- 作为读者，我希望按章节浏览视频摘要并跳转 KB 文档
- 作为开发者，我希望长视频必须 202 异步，object upload 与 ingest 分步

## 范围

- POST videos/ingest、GET videos list/detail；transcript 分片（3b）
- Console 视频 ingest、章节时间轴、转写展示
- Core 只读：objects/upload

## 非目标

- 视频播放器 CDN 契约（Console 前端实现）
- 合入 ANI-main

## 成功标准

- [ ] ingest 与 object-storage-upload 分工清晰
- [ ] 合入 YAML 后详文切换为冻结口径
