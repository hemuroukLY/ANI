# PRD: Console compliance

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/security/compliance.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-compliance-draft.md`

## 目标

- 定义 Console **只读合规摘要与报告** 与 BOSS 的边界
- 输出 Services GET summary/reports 草案
- 为 TASK-SVC-018 / Phase 3 §4 提供合入 YAML 前置材料

## 用户故事

- 作为租户管理员，我希望查看控制项 pass/fail 摘要
- 作为合规专员，我希望下载 BOSS 已生成的 ready 报告
- 作为开发者，我希望 Console 无 POST 创建报告 API

## 范围

- GET compliance/summary、reports list/detail
- download_url 只读跳转

## 非目标

- BOSS 报告生成、全量取证
- 合入 ANI-main

## 成功标准

- [ ] 无 write path；BOSS 边界 UI 明示
- [ ] 合入 YAML 后详文切换为冻结口径
