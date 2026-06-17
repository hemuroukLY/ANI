# Console 文档治理一页式同步简报

## 适用场景

本文用于当前阶段的内部同步、转发和交接，目标是用一页内容快速说明：

- 现在到底哪些模块已经可以继续维护
- 哪些页面只是导航占位，不能当成正式能力
- 哪些旧材料不能再作为当前基线
- 后续继续维护时应该先看哪些文件

本文只讨论产品文档治理，不讨论前端代码实现，不包含 `ANI-main` 内文档更新。

## 一句话结论

当前 `Console` 文档体系已经完成主摘要层收口，正式可继续维护的模块共有 `20` 个；`开放与集成` 下部分子页仍只是导航占位，`prototypes/ani-services-prototype.html` 已降级为旧版参考，不再作为当前判断依据。

## 当前可继续维护的 20 个模块

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

## 当前必须记住的冻结边界

- `平台概览` 是聚合页，没有独立冻结接口、schema 或返回码
- `API Key 管理` 当前只有列表、创建、吊销三条正式链路，`key_value` 仅在创建成功时返回一次
- `安全与身份概览` 只收口 `OIDC begin / token / refresh / logout` 与 `API Key 管理` 总入口
- `租户用量报表` 的正式查询边界只有 `GET /api/v1/metering/usage`
- `开放与集成总入口` 当前只稳定收口 `REST API`、`OIDC`、`API Key 管理`、`ani CLI` 源码入口
- `模型中心` 属于 `Services / Models`，导入模型是 `202 + AsyncTask`
- `推理服务` 属于 `Services / InferenceServices`，部署与删除按 `202` 语义处理
- `知识库管理` 属于 `Services / KnowledgeBases`，文档上传是 `202`，流式问答方法是 `GET /query/stream`

## 当前不能当正式能力的内容

以下内容当前只算导航占位或说明态，不能统计进“已冻结能力”：

- `OpenAI 兼容 API（待补）`
- `Go SDK（待补）`
- `Python SDK（待补）`
- `TypeScript Client（待补）`
- `Java SDK（待补）`
- `Webhook（待补）`
- `企业微信 / 钉钉 Bot（待补）`
- `第三方业务系统集成（待补）`

补充说明：

- `ani CLI（入口说明）` 当前只确认仓内源码入口 `ANI-main/repo/cli/ani/main.go`
- 这不等于已经存在独立的 `Console` 下载页、配置页或交付页

## 当前不能再作为基线的材料

- `prototypes/ani-services-prototype.html`

当前处理原则：

- 只保留为 `旧版参考`
- 可以回看历史表达
- 不能继续作为当前命名、首页回跳、冻结边界和模块关系的依据

## 后续查看顺序

后续继续维护时，统一按下面顺序查看：

1. [final-handoff-summary.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/final-handoff-summary.md)
2. [console-document-status-board.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/console-document-status-board.md)
3. [governance-closeout-checklist.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/governance-closeout-checklist.md)
4. 对应模块主维护文档
5. 对应领域辅助维护文档
6. `prototypes/ani-services-prototype-console.html`

## 后续维护规则

- 不再把“文件存在”直接等同于“模块已完成”
- 不再把导航占位页统计进“已冻结能力”
- 不再引用 `prototypes/ani-services-prototype.html` 作为当前口径来源
- 后续若修正文档，先改状态板，再改主维护文档和 HTML 摘要
- 后续若某个待补页真的进入冻结范围，必须先补权威源依据，再从占位名单移出

## 推荐转发方式

若需要给未参与历史对话的人同步当前状态，建议最少转发以下 4 份材料：

- [final-sync-brief.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/final-sync-brief.md)
- [final-handoff-summary.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/final-handoff-summary.md)
- [console-document-status-board.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/console-document-status-board.md)
- [governance-closeout-checklist.md](file:///Users/viccao/Desktop/service-ani/docs/console-modules/governance/governance-closeout-checklist.md)
