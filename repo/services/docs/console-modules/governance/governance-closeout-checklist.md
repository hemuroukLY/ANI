# Console 文档治理收口清单

## 目的

本文用于在当前阶段给出一份可直接执行的 `Console` 文档治理结论，明确：

- 哪些模块已经达到可继续维护的文档基线
- 哪些内容只是导航占位或说明态，不能算作已冻结能力
- 哪些旧材料不再作为当前维护基线

本文只讨论产品文档治理，不讨论前端代码实现。

## 适用范围

- `prototypes/ani-services-prototype-console.html`
- `docs/console-modules/**/*.md`
- `tasks/support/support-console-*.md`
- `docs/console-modules/governance/console-document-status-board.md`

不包含：

- `ANI-main` 内文档更新
- `ANI-main` 内前端实现是否已完全落地当前命名与导航

## 使用规则

- `可用` 表示文档层 `HTML / PRD / SPEC / 主维护文档 / 辅助维护文档` 已完成一致性收口，可继续作为维护基线
- `可用` 不等于相关前端实现、路由名称或运行时页面已经全部落地
- `导航占位` 表示当前只保留入口名称或说明态，不能被引用为已冻结产品能力
- `旧版参考` 表示材料仍可保留历史上下文，但不能继续作为当前口径来源

## 当前可用模块

以下模块当前均可按正式维护基线继续使用：

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

## 当前已冻结的高频结论

- **GPU 容器 / Sandbox** 统一走 Core `/api/v1/instances*`（`kind=gpu_container` / `kind=sandbox`）；`services/v1.yaml` 中旧路径已 deprecated

- `平台概览` 是聚合页，没有独立冻结接口、独立 schema 或独立返回码
- `API Key 管理` 的正式边界仍只有列表、创建、吊销，且 `key_value` 只在创建成功时返回一次
- `安全与身份概览` 只收口 `OIDC begin / token / refresh / logout` 与 `API Key 管理` 总入口
- `租户用量报表` 的正式查询边界只有 `GET /api/v1/metering/usage`
- `开放与集成总入口` 当前只稳定收口 `REST API`、`OIDC`、`API Key 管理`、`ani CLI` 源码入口
- `模型中心` 属于 `Services / Models`，导入模型仍是 `202 + AsyncTask`
- `推理服务` 属于 `Services / InferenceServices`，部署与删除都按 `202` 语义处理
- `知识库管理` 属于 `Services / KnowledgeBases`，文档上传是 `202`，流式问答方法是 `GET /query/stream`

## 当前仅为导航占位的内容

以下内容目前不能算作已冻结能力，只能作为导航占位或说明态存在：

- `OpenAI 兼容 API（待补）`
- `Go SDK（待补）`
- `Python SDK（待补）`
- `TypeScript Client（待补）`
- `Java SDK（待补）`
- `Webhook（待补）`
- `企业微信 / 钉钉 Bot（待补）`
- `第三方业务系统集成（待补）`

额外说明：

- `ani CLI（入口说明）` 当前只确认 `ANI-main/repo/cli/ani/main.go` 这一仓内源码入口
- 上述 CLI 入口不等于已经存在独立 `Console` 下载页、配置页或交付页

## 不计入已冻结能力的页面

以下页面即使在 HTML 中存在，也不能计入“当前可用模块”：

- `开放与集成` 域下的待补接入页
- `安全与身份` 域下尚未冻结的配置中心、审计、合规、密钥管理类页
- 任何只写了“占位页面 / 后续补充 / Phase 2”的页面

判定原则：

- 若正文只说明“后续补详细页面结构、接口、状态机和错误码”，则当前仍属于占位态
- 若正文无法回指权威源中的正式路径、schema 或返回码，则当前不能算作已冻结能力

## 当前不再作为维护基线的材料

- `prototypes/ani-services-prototype.html`

当前处理方式：

- 该文件仅保留为 `旧版参考`
- 可用于回溯历史表达，但不能继续作为当前命名、回跳关系和冻结边界的判断依据
- 当前基线应统一以 `prototypes/ani-services-prototype-console.html` 和 `docs/console-modules/**/*.md` 为准

## 当前唯一总控入口

后续继续维护时，应优先查看以下文件：

- `docs/console-modules/governance/final-handoff-summary.md`
- `docs/console-modules/governance/console-document-status-board.md`
- `docs/console-modules/governance/governance-closeout-checklist.md`

建议顺序：

1. 先看最终交付版总览说明，快速理解当前全局结论
2. 再看状态板判断模块当前是否为 `可用 / 待收口 / 待重写`
3. 再看本文判断该模块是否包含“仅占位、不计入冻结能力”的子页
4. 最后再进入对应主维护文档核对正式路径、schema 和边界

## 后续维护规则

- 不再把“文件存在”直接等同于“模块已完成”
- 不再把导航占位页统计进“已冻结能力”
- 不再引用 `prototypes/ani-services-prototype.html` 作为当前口径来源
- 后续若新增或修正文档，先改状态板，再改主维护文档和 HTML 摘要
- 后续若某个待补页真的进入冻结范围，必须先补权威源依据，再移出本文的“导航占位”名单
