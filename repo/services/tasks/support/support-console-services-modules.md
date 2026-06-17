# Console 辅助维护文档：Services 业务模块

## 1. 适用范围

- `docs/console-modules/inference/model-center.md`
- `docs/console-modules/inference/inference-service.md`
- `docs/console-modules/knowledge/knowledge-base.md`

本文件用于补充上述模块的日志规范、异常处理方案、版本变更记录和运维操作指南。

## 2. 日志规范

| 模块 | 必记事件 | 必记字段 | 脱敏要求 |
|---|---|---|---|
| 模型中心 | `model_list_loaded`、`model_created`、`model_import_submitted`、`model_version_created`、`model_deleted` | `module`、`action`、`request_id`、`model_id?`、`task_id?`、`result`、`duration_ms?` | 不记录模型内容、密钥、上传源凭据 |
| 推理服务 | `inference_service_list_loaded`、`inference_service_created`、`inference_service_deleted`、`inference_service_detail_loaded` | `module`、`action`、`request_id`、`service_id?`、`result`、`duration_ms?` | 不记录访问令牌、未脱敏 endpoint 凭证 |
| 知识库管理 | `kb_list_loaded`、`kb_created`、`kb_document_uploaded`、`kb_query_submitted`、`kb_query_stream_started`、`kb_deleted` | `module`、`action`、`request_id`、`kb_id?`、`document_id?`、`stream`、`result` | 不记录文档原文、问答原始敏感内容、上传源凭据 |

统一要求：

- 所有前端操作日志都要带 `module` 与 `action`
- 所有失败日志都要记录 `request_id`
- 异步或流式场景必须补充 `task_id`、`operation_id` 或等价追踪标识
- 日志只记录定位问题所需的最小字段，不额外复制业务正文

## 3. 异常处理方案

| 模块 | 异常场景 | 前端处理 | 运维检查 |
|---|---|---|---|
| 模型中心 | 列表或详情加载失败 | 保留页面壳与当前筛选，展示错误文案和重试入口 | 检查 `Services / Models` 路径可用性与最近错误请求 |
| 模型中心 | 导入任务失败 | 在任务反馈区展示失败结果与 `request_id`，不假定导入成功 | 检查异步任务状态、导入源可达性和模型记录是否已落盘 |
| 推理服务 | 部署请求失败或详情异常 | 保留服务列表，详情区显示失败信息，不伪造运行中状态 | 检查 `InferenceServices` 主资源状态、最近创建请求和删除动作结果 |
| 推理服务 | 删除冲突或删除后未刷新 | 维持当前详情与列表上下文，展示删除失败或处理中反馈 | 检查资源是否仍处于后端处理中，确认结果是否已同步回列表 |
| 知识库管理 | 文档上传失败 | 仅文档区报错，不阻断知识库详情与问答区 | 检查上传请求、文档记录状态与文件源可达性 |
| 知识库管理 | 同步/流式问答失败 | 结束当前问答态并展示失败反馈，不继续拼接不完整答案 | 检查查询请求、知识库主对象状态和流式连接中断原因 |

统一要求：

- 失败提示统一遵循 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}` 口径
- 不把待补能力的失败解释成已冻结能力异常
- 局部失败只影响对应区块，不扩大成整页不可用

## 4. 运维操作指南

| 模块 | 日常检查 | 常见处置 | 变更注意 |
|---|---|---|---|
| 模型中心 | 检查模型列表是否可读、最近导入任务是否有失败积压 | 导入失败时先核对任务结果，再决定是否允许用户重试 | 版本能力扩展前先更新主维护文档与 SPEC，不直接补新子页 |
| 推理服务 | 检查服务列表、详情和最近部署/删除结果是否一致 | 部署异常时优先确认服务主对象状态，不把待补观测能力当成现有排障入口 | 若未来补 OpenAI 或访问策略，必须先冻结独立契约再扩页 |
| 知识库管理 | 检查知识库详情、文档区、问答区是否都能独立工作 | 文档失败与问答失败要分开排查，避免混成统一“知识库故障” | 向量化、引用、历史、权限等增强能力必须先补正式边界再接入页面 |

## 5. 版本变更记录

| 模块 | 版本 | 日期 | 变更内容 |
|---|---|---|---|
| 模型中心 | `v1.2` | `2026-06-16` | 复核 `Services / Models` 正式边界，确认 `svc/models*`、`import=202 + AsyncTask` 与 HTML 摘要一致 |
| 推理服务 | `v1.2` | `2026-06-16` | 复核 `Services / InferenceServices` 正式边界，确认 `svc/inference-services*`、部署/删除 `202` 语义与 HTML 摘要一致 |
| 知识库管理 | `v1.2` | `2026-06-16` | 复核 `Services / KnowledgeBases` 正式边界，确认 `documents*`、`query*`、上传 `202` 与流式问答方法口径一致 |
