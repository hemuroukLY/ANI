# BOSS 模块补充执行流程

本文档是 BOSS 侧对 [`ANI-main/ANI-14-API对齐与开发工作流.md`](../../ANI-main/ANI-14-API对齐与开发工作流.md) Phase 1 的执行细则，并继承 [`docs/console-modules/governance/module-delivery-workflow.md`](../../console-modules/governance/module-delivery-workflow.md) 的 Core 优先、正文优先、零回合交付原则。

## 1. 适用范围

- 适用于 `BOSS` 侧功能模块的新增、补充、收口和对齐工作
- 适用于需要同时满足 Core 规范、平台级 RBAC、OpenAPI 对齐要求的模块
- 当前模块范围以 `boss-document-status-board.md` 与 `module-completion-matrix.md` 为准

## 2. 标准产出物（与 Console 对齐）

| 产物 | 路径 | 角色 |
|---|---|---|
| 模块详文（主维护源） | `docs/boss-modules/{域}/{module}.md` | **唯一正式口径** |
| PRD | `tasks/modules/prd/prd-boss-{module}.md` | 产品目标、用户故事、验收 |
| SPEC | `tasks/modules/spec/spec-boss-{module}.md` | 权威源事实、技术边界 |
| HTML 摘要 | `prototypes/ani-services-prototype-boss.html` | 导航与边界入口 |

修改顺序：**详文 → PRD/SPEC 同步 → HTML 摘要 → 状态板更新**

## 3. ANI-14 Phase 1 强制门禁

执行每个 BOSS 模块时必须完成：

1. **锁定权威源**：`v1.yaml`、`services/v1.yaml`、Console 对照模块详文（若有）
2. **Core 层基线**：路径前缀、分层、错误结构、`idempotency_key`、§2.10 的 422 口径
3. **BOSS 增量**：只写平台级差异，不复制 Console 租户契约
4. **冻结事实表**：区分「YAML 已声明（只读参考）」与「BOSS 平台 API 待补」
5. **创建前置条件 + 操作可用性矩阵**（Phase 1 规则 4/5）
6. **四轮回审**：权威源对照 → 跨材料一致性 → 反自造契约 → 页面能力降噪

## 4. BOSS 专属原则

### 平台视角优先

- BOSS 默认 **全平台 / 多租户 / 资源池** 视角
- 平台级 list 必须走后端平台 RBAC；**禁止**信任前端 `tenant_id` 越权
- 租户内 Console 路径（如 `/api/v1/svc/tenant/members`）**不能**直接写成 BOSS 平台租户 CRUD 的冻结契约

### Console 对照而非重复

- 已有 Console 详文的能力，BOSS 只写 **平台级增量**
- 字段口径与 Console 同源模块一致，聚合范围不同

### 交付与安装类

- 权威源为 `ANI-07` 等交付文档，**TODO-YAML: N/A**
- 不得自造 REST 契约

## 5. 禁止自造契约（与 Console 相同）

- 不得把 TODO-YAML 写成已实现
- 不得编造 schema / operationId / 返回码
- 不得把 handler stub 写成已实现
- OpenAPI 已声明的租户级路径 ≠ BOSS 平台级 API 已冻结

## 6. 相关文件

- [boss-document-status-board.md](./boss-document-status-board.md)
- [module-completion-matrix.md](./module-completion-matrix.md)
- [../../console-modules/governance/ANI-14-revision-checklist.md](../../console-modules/governance/ANI-14-revision-checklist.md)
