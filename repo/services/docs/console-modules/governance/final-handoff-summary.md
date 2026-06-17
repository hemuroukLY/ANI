# Console 产品文档最终交付版总览说明

## 说明

本文是当前阶段的最终交付版总览说明，用于在不回看历史对话的前提下，直接回答以下问题：

- 现在哪些模块已经稳定
- 哪些内容只是导航占位，不能按正式能力理解
- 哪些旧材料不应继续作为当前基线
- 后续继续维护时应该先看什么、再看什么

本文只讨论产品文档治理，不讨论前端代码实现或 `ANI-main` 内前端落地状态。

## 一句话结论

当前 `Console` 文档体系已完成 **86 套模块详文**（21 基线 + 65 backlog）Core 标准收口；**OpenAPI/handler** 仍按 backlog 与 YAML 草案推进。GPU 容器与 Sandbox 以 Core 统一实例 `/api/v1/instances*` 为准。

## 已稳定模块清单

**基线 21 模块** + **P0/P1/P2 共 65 子模块**均已具备详文 + PRD/SPEC。完整索引见 `module-completion-matrix.md` 与 `README.md`。

基线模块（2026-06-16 起）：

- `平台概览`
- `GPU 算力管理`
- `云主机 VM`
- `容器实例`
- `GPU 容器实例`
- `Sandbox 实例`
- `K8s 集群`
- `网络管理`
- `存储管理`
- `块存储`
- `对象存储`
- `文件存储`
- `向量存储`
- `API Key 管理`
- `安全与身份概览`
- `租户用量报表`
- `开放与集成总入口`
- `模型中心`
- `推理服务`
- `知识库管理`
- `租户管理`

## 开发任务清单（Phase 3）

- `tasks/execution/CORE-TEAM-TASKS.md`
- `tasks/execution/SERVICES-TEAM-TASKS.md`
- `tasks/execution/TASK-DEPENDENCY-MAP.md`
- Schema 跟踪：`docs/console-modules/governance/schema-completion-tracker.md`

## 当前正式冻结的关键口径

- `平台概览`：聚合页，无独立冻结接口、schema、operationId 或返回码
- `API Key 管理`：正式能力只有列表、创建、吊销；`key_value` 仅在创建成功时返回一次
- `安全与身份概览`：只收口 `OIDC begin / token / refresh / logout` 与 `API Key 管理` 总入口
- `租户用量报表`：正式查询边界只有 `GET /api/v1/metering/usage`
- `开放与集成总入口`：当前只稳定收口 `REST API`、`OIDC`、`API Key 管理`、`ani CLI` 源码入口
- `模型中心`：属于 `Services / Models`，导入模型是 `202 + AsyncTask`
- `推理服务`：属于 `Services / InferenceServices`，部署与删除都是 `202` 语义
- `知识库管理`：属于 `Services / KnowledgeBases`，文档上传 `202`，流式问答方法是 `GET /query/stream`

## 仅占位能力清单（已更新 2026-06-17）

以下条目**已有详文**，但 OpenAPI 仍为 TODO-YAML 或 N/A（交付页）：

- AI 原生整域（×7）— 见 `ai-native/README.md`
- 知识库智能（×3）、模型加密/推荐、多数安全合规项
- SDK 四件套 + CLI 下载页 — `TODO-YAML: N/A`

以下已从「待补导航」升级为有详文（handler 仍可能 stub）：

- ~~`Webhook（待补）`~~ → `integration-webhook-overview.md`
- ~~`Go/Python/TS/Java SDK（待补）`~~ → 各 `integration-*-sdk.md`

额外说明：

- `ani CLI（入口说明）` 当前只确认仓内源码入口 `ANI-main/repo/cli/ani/main.go`
- 这不等于已经存在独立的 `Console` 下载页、配置页或交付页

## 不应再误判为正式能力的页面

以下页面即使在 HTML 中存在，也不能统计进“当前可用模块”：

- `开放与集成` 域下的待补接入页
- `安全与身份` 域下未冻结的配置中心、审计、合规、密钥管理类页
- 任意正文只写了“占位页面 / 后续补充 / Phase 2”的页面

判定标准：

- 不能回指权威源正式路径、schema、返回码的，不算已冻结能力
- 只承认“后续补充”而没有冻结依据的，不算正式页面

## 旧版参考材料清单

以下材料当前不再作为维护基线：

- `prototypes/ani-services-prototype.html`

当前角色：

- 仅保留为 `旧版参考`
- 可用于回溯历史表达
- 不可继续作为当前命名、首页回跳、冻结边界、模块关系的判断依据

## 当前应使用的唯一查看顺序

后续继续维护时，统一按下面顺序查看：

1. [console-document-status-board.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/console-document-status-board.md)
2. [governance-closeout-checklist.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/governance-closeout-checklist.md)
3. 对应模块主维护文档
4. 对应领域辅助维护文档
5. `prototypes/ani-services-prototype-console.html`

不要反过来从旧 HTML 或零散页面标题倒推当前结论。

## 后续维护建议

- 后续若继续改文档，先更新状态板，再更新主维护文档和 HTML 摘要
- 后续若某个待补页真正进入冻结范围，必须先补权威源依据，再从“仅占位能力清单”中移出
- 后续若要汇报当前状态，优先使用 `可用 / 待收口 / 待重写`，不要只说“已完成”
- 后续若发生新漂移，先判断是 `命名漂移`、`首页回跳漂移` 还是 `待补能力被写成正式能力`

## 交接建议

如果后续由其他人继续维护，建议把以下 3 份文件作为最小交接包：

- [final-handoff-summary.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/final-handoff-summary.md)
- [console-document-status-board.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/console-document-status-board.md)
- [governance-closeout-checklist.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/governance-closeout-checklist.md)

如果只是做内部同步或需要快速转发当前结论，可额外附上：

- [final-sync-brief.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/final-sync-brief.md)
