# PRD: Console 知识库

## 1. Introduction/Overview

`Console / 知识库与 AI 应用 / 知识库管理` 是租户侧知识管理页面，用于管理知识库主对象、知识库文档和问答能力。

本模块属于 `Services` 层，权威源为 `ANI-main/repo/api/openapi/services/v1.yaml` 中的 `KnowledgeBases` 路径组。正式访问前缀为 `/api/v1/svc`，不允许把知识库资源写入 `Core /api/v1/*`。

当前已确认的正式能力：

- `GET /api/v1/svc/knowledge-bases`
- `POST /api/v1/svc/knowledge-bases`
- `GET /api/v1/svc/knowledge-bases/{kb_id}`
- `DELETE /api/v1/svc/knowledge-bases/{kb_id}`
- `GET /api/v1/svc/knowledge-bases/{kb_id}/documents`
- `POST /api/v1/svc/knowledge-bases/{kb_id}/documents`
- `DELETE /api/v1/svc/knowledge-bases/{kb_id}/documents/{doc_id}`
- `POST /api/v1/svc/knowledge-bases/{kb_id}/query`
- `GET /api/v1/svc/knowledge-bases/{kb_id}/query/stream`

当前仍属于待补边界的能力：

- 向量化与索引独立资源
- 来源引用
- 对话历史
- 知识库权限
- 文档智能 / 会议智能 / 视频智能

## 2. Goals

- 让租户用户管理知识库主对象
- 支持上传与删除知识库文档
- 支持知识库问答与 SSE 流式问答
- 明确向量化、来源引用、对话历史、权限目前不作为已冻结独立资源

## 3. User Stories

### US-001: 查看知识库列表与详情

作为知识管理员，我希望查看知识库列表与详情，以便确认知识库基础信息和当前状态。

**Acceptance Criteria**

- [ ] 页面通过 `GET /api/v1/svc/knowledge-bases` 查询知识库列表
- [ ] 页面通过 `GET /api/v1/svc/knowledge-bases/{kb_id}` 查询知识库详情
- [ ] 页面至少展示 `name`、`status`、`doc_count`、`created_at`

### US-002: 创建与删除知识库

作为知识管理员，我希望创建和删除知识库，以便管理不同知识域。

**Acceptance Criteria**

- [ ] 页面通过 `POST /api/v1/svc/knowledge-bases` 创建知识库
- [ ] 创建请求至少包含 `idempotency_key`、`name`
- [ ] 成功返回按 `201 + KnowledgeBase` 处理
- [ ] 页面通过 `DELETE /api/v1/svc/knowledge-bases/{kb_id}` 删除知识库
- [ ] 删除成功按 `204` 处理
- [ ] 删除属于危险动作，页面必须有明确确认与回流反馈

### US-003: 上传与删除知识库文档

作为内容编辑成员，我希望查看、上传和删除知识库文档，以便维护知识库内容。

**Acceptance Criteria**

- [ ] 页面通过 `GET /api/v1/svc/knowledge-bases/{kb_id}/documents` 查看文档列表
- [ ] 页面通过 `POST /api/v1/svc/knowledge-bases/{kb_id}/documents` 上传文档
- [ ] 上传成功按 `202 + KBDocument` 处理，表示解析任务进行中
- [ ] 页面通过 `DELETE /api/v1/svc/knowledge-bases/{kb_id}/documents/{doc_id}` 删除文档
- [ ] 删除成功按 `204` 处理
- [ ] 页面不把文档解析、索引构建写成独立正式资源

### US-004: 知识库问答

作为内容编辑成员，我希望在知识库上下文内执行同步或流式问答，以便快速验证知识效果。

**Acceptance Criteria**

- [ ] 页面通过 `POST /api/v1/svc/knowledge-bases/{kb_id}/query` 执行普通问答
- [ ] 页面通过 `GET /api/v1/svc/knowledge-bases/{kb_id}/query/stream` 执行 SSE 流式问答
- [ ] 同步问答返回 `200 + KBQueryResponse`
- [ ] 流式问答返回 `200 + text/event-stream`
- [ ] 页面不把来源引用、对话历史、权限写成当前已冻结独立资源

### US-005: 保持边界清晰

作为模块维护人，我希望文档明确哪些能力已经冻结、哪些仍待补，避免把规划项写成正式契约。

**Acceptance Criteria**

- [ ] 页面正文明确知识库属于 `Services`
- [ ] 页面正文不把知识库对象写成 `Core` 资源
- [ ] 页面正文不把向量化、来源引用、对话历史、权限写成当前已冻结能力
- [ ] 页面正文使用真实路径方法和成功返回码

## 4. Functional Requirements

- FR-1: 系统必须支持知识库列表、详情、创建和删除
- FR-2: 系统必须支持知识库文档列表、上传和删除
- FR-3: 系统必须支持同步问答和 SSE 流式问答
- FR-4: 系统必须展示知识库、文档和问答的关键状态
- FR-5: 系统必须支持只读用户查看知识库与文档列表，但不显示写操作和问答提交入口
- FR-6: 系统必须明确待补边界，不把未冻结能力写成正式接口

## 5. Business Rules

- 知识库创建成功 `201` 表示主对象已创建
- 文档上传成功 `202` 表示文档已上传且解析任务进行中，不等于文档已经完成索引
- 文档删除和知识库删除成功返回 `204`
- 文档状态当前只承接权威源中已确认的枚举：`uploaded`、`parsing`、`indexed`、`failed`、`deleted`
- 知识库状态当前只承接 `active`、`rebuilding`、`deleted`
- 问答区和文档区允许独立失败，但失败不应阻断整个知识库详情

## 6. Non-Goals

- 不在本轮扩写独立向量化索引资源
- 不在本轮扩写来源引用、对话历史、权限独立资源
- 不在本轮实现文档智能、会议智能、视频智能正式路径
- 不在本轮定义问答会话管理独立资源

## 7. Design Considerations

- 列表页突出“知识库状态 + 文档数量”
- 详情页突出“基础信息 + 文档区 + 问答区”
- 文档上传必须保留解析中反馈，不能误导为同步完成
- 流式问答必须明确区分开始、增量输出和结束

## 8. Technical Considerations

- 权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 路径归属：`Services / KnowledgeBases`
- 服务前缀：`/api/v1/svc`
- 页面不要求前端显式传 `tenant_id`
- 创建知识库与同步问答需要承接 `idempotency_key`
- 文档上传使用 `multipart/form-data`
- 错误响应复用 `ErrorResponse`

## 9. Success Metrics

- 用户能完成“创建知识库 -> 上传文档 -> 执行问答 -> 查看结果”的闭环
- 文档上传与问答都具备明确的过程态和失败态
- 模块文档与 `services/v1.yaml` 的路径、方法、返回码、schema 名称保持一致
- 文档中不再出现错误的删除文档路径或错误的流式问答方法

## 10. Open Questions

- 文档上传后是否需要未来补充任务历史或解析详情页
- 流式问答的结束事件和错误事件是否需要未来统一前端事件协议
- `doc_count` 是否需要未来区分“已索引文档数”和“总文档数”

## 11. 回填前置依赖

- 如未来补充向量化、来源引用、对话历史、权限，必须先冻结正式路径与 schema
- 如未来补充智能增强能力，必须先明确其资源归属是否仍在 `KnowledgeBases`
- 如未来补充问答会话资源，必须先明确其与 `session_id` 的正式关系
