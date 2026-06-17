# Console 文档总控状态表

## 说明

本表用于在当前阶段作为 `Console` 产品文档的唯一总控视图，解决“文件已经存在，但质量状态不清楚”的问题。

当前阶段只讨论产品文档，不讨论前端代码。

本文将每个模块拆成 3 类状态：

- `可用`：`HTML / PRD / SPEC / 主维护文档 / 辅助维护文档` 均已存在，且主维护文档结构基本完整，可继续作为正式维护基线
- `待收口`：文档体系已存在，但命名、章节深度、HTML 同步口径或跨材料一致性仍需统一
- `待重写`：虽然文件存在，但主维护文档或 `PRD / SPEC` 明显低于当前仓内成熟模块标准，不宜直接作为正式交付材料

**子维度（2026-06-16 起强制汇报）**：

| 子维度 | 含义 |
|---|---|
| 文档结构 | 主维护源章节完整（含前置条件、状态×操作矩阵） |
| 权威源对齐 | 路径/schema/返回码与 OpenAPI 一致；不得把规划路径或建议错误码写成冻结契约 |
| 开发可执行 | 已纳入 TASK 清单，可驱动 handler 实现 |

### 汇总（2026-06-17 P0 文档收口后）

| 维度 | 状态 | 说明 |
|---|---|---|
| 文档文件完整性 | ✅ | 21 模块 PRD / SPEC / 主维护源 / HTML 摘要均存在 |
| 文档结构 | ⚠️ → ✅（P1 2026-06-17） | 21 模块均已统一 `## 创建前置条件`；容器实例已补状态×操作矩阵 |
| 权威源对齐 | ⚠️ | **不可整体标 ✅**；见下表分项；P0 已修 phantom 路径、自造错误码、向量状态枚举、安全域口径 |
| 开发可执行 | ✅（文档层） | Core Handler 指南 + 模块详文已齐；**handler 代码由 Core 团队实现** |

额外说明：

- `HTML` 当前只用于页面摘要、导航关系、模块边界和详细材料入口，不等于完整详文
- `支持文档` 以领域级文档复用，不为每个模块单独复制一份
- 当前模块范围以 `module-completion-matrix.md` 为准
- 未定义 Console 功能见 `console-undefined-features-backlog.md`
- 如需快速判断“哪些算可用、哪些只是占位、哪些旧材料不再作为基线”，统一参考 `docs/console-modules/governance/governance-closeout-checklist.md`
- 前置条件错误码写法统一见 `../governance/module-delivery-workflow.md` §2.10

## 支持文档映射

| 模块域 | 支持文档 |
|---|---|
| 首页 / 租户 / 接入 | `tasks/support/support-console-overview-tenant-integration.md` |
| 算力 / 实例 / 网络 | `tasks/support/support-console-compute-modules.md` |
| 存储 | `tasks/support/support-console-storage-modules.md` |
| Services | `tasks/support/support-console-services-modules.md` |

## 权威源对齐明细

图例：✅ 路径/schema/返回码与 YAML 一致 | ⚠️ 文档可维护但仍有缺口或依赖 stub | ❌ 存在 phantom 路径/自造契约（P0 已修项不再标 ❌）

| 模块 | 权威源对齐 | 备注 |
|---|---|---|
| 平台概览 | ✅ | 聚合页，无独立 API；不自造路径 |
| GPU 算力管理 | ⚠️ | `gpu-inventory` **Phase 2 已声明**（见 `gpu-inventory-ui.md`）；handler stub |
| 云主机 VM | ⚠️ | 实例路径 ✅；部分 422 `code` 为建议语义（P0-2 已改） |
| 容器实例 | ⚠️ → ✅ | 统一实例路径 ✅；P1 已补状态×操作矩阵 |
| GPU 容器实例 | ⚠️ | 统一实例路径 ✅；`GPU_UNAVAILABLE` 等为建议语义 |
| Sandbox 实例 | ⚠️ | 统一实例路径 ✅；`IMAGE_NOT_FOUND` 为 YAML 举例 |
| K8s 集群 | ✅ | 路径 ✅；创建 422 **YAML 已声明**（Core v1.yaml Phase 2 2026-06-17，P0 已修） |
| 网络管理 | ✅ | 路径与 schema 对齐 |
| 存储管理 | ✅ | 总览不自造子路径 |
| 块存储 / 对象存储 | ✅ | 边界清晰 |
| 文件存储 | ✅ | `StorageResourceState` 对齐 |
| 向量存储 | ✅ | `VectorStoreState` 用 `ready`（P0-3 已修）；search 无 422 已标注 |
| API Key 管理 | ✅ | Core Auth 路径对齐 |
| 安全与身份概览 | ✅ | Core Auth 对齐；租户入口分流至租户管理（P0-4 已修） |
| 租户用量报表 | ✅ | `GET /metering/usage` 对齐 |
| 开放与集成总入口 | ✅ | 导航聚合，占位清晰 |
| 租户管理 | ⚠️ | Services YAML 有路径；handler 未实现 |
| 模型中心 | ⚠️ | Services stub |
| 推理服务 | ⚠️ | 路径在 YAML；部分 422 code 为建议语义 |
| 知识库管理 | ⚠️ | Services stub |

## 模块状态总表

| 模块 | HTML 摘要 | PRD | SPEC | 主维护文档 | 支持文档 | 当前状态 | 当前判断 |
|---|---|---|---|---|---|---|---|
| 平台概览 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | 聚合页无独立冻结接口；首页回跳与可用模块入口一致 |
| GPU 算力管理 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | 页面定义 + Phase 2 gpu-inventory；handler stub |
| 云主机 VM | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | 统一实例 `kind=vm`；422 错误码按 §2.10 口径 |
| 容器实例 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | 统一实例 `kind=container` |
| GPU 容器实例 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | 统一实例 `kind=gpu_container`；deprecated Services 路径已拒绝 |
| Sandbox 实例 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | 列表来自 `kind=sandbox` 查询 |
| K8s 集群 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | 独立 `K8sClusters` 资源组；创建 422 **已在 Core YAML 声明**（Phase 2 2026-06-17） |
| 网络管理 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | `networks/*` 正式边界 |
| 存储管理 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | 子模块入口与待补边界 |
| 块存储 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | `volumes*` 正式边界 |
| 对象存储 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | 仅对象元数据，非上传 |
| 文件存储 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | `filesystems*` 正式边界 |
| 向量存储 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | `vector-stores*` + `search`；状态枚举 `ready` |
| API Key 管理 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | `api-keys*` 与一次性 `key_value` |
| 安全与身份概览 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core Auth + 租户管理入口分流 |
| 租户用量报表 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | `GET /metering/usage` 单一查询 |
| 开放与集成总入口 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Webhook/Bot/第三方已有子模块；SDK 仍为占位 |
| 模型中心 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Services `/models*`；handler stub |
| 推理服务 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | `/inference-services*`；部署/删除 `202` |
| 知识库管理 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | `/knowledge-bases*`；上传 `202` |
| 租户管理 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Services `/tenant/*`；与安全域 HTML 已对齐 |
| 告警与待处理事项 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚；handler/TODO-YAML 聚合 API 仍待 |
| 异步任务中心 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚；list tasks TODO-YAML |
| 推理服务调用测试 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚 |
| 推理服务可观测性 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚；events TODO-YAML |
| 知识库问答流程 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚 |
| 首页资源使用概览 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚；无独立 API |
| 首页推理服务状态 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚 |
| 首页 GPU 利用率 | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚 |
| P1 网络子模块（×5） | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚；见 compute/network/*.md |
| P1 算力/存储/租户/推理/搜索（×13） | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚；见 backlog P1 表 |
| P2 YAML 已有项（×13） | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 | Core 标准已加厚；见 backlog P2 表 |
| P2 新模块（×26） | 已有 | 已有 | 已有 | 已有 | 已覆盖 | 可用 / 待 YAML | 详文已补；整域 YAML/handler 仍待 Core/Services |

## 架构决策（2026-06-16）

- **GPU 容器 / Sandbox**：以 Core 统一实例 `/api/v1/instances*` 为准（`kind=gpu_container` / `kind=sandbox`）
- `services/v1.yaml` 中 `/gpu-containers*`、`/sandboxes*` 已全部 `deprecated`，不得在新代码中使用
- `listInstances` 的 `kind` 枚举已含 `sandbox`（与 createInstance 对齐）
- Phase 3 任务清单：`tasks/execution/CORE-TEAM-TASKS.md`、`tasks/execution/SERVICES-TEAM-TASKS.md`、`tasks/execution/TASK-DEPENDENCY-MAP.md`
- Schema 跟踪：`docs/console-modules/governance/schema-completion-tracker.md`

## P0 文档收口记录（2026-06-17）

| 项 | 动作 | 涉及文件 |
|---|---|---|
| P0-1 phantom `gpu-inventory` | 改为「待 Core 冻结」 | `gpu-management.md`、`module-completion-matrix.md`、backlog、HTML |
| P0-2 自造错误码 | 改建议语义 / YAML 已举例 | VM、GPU 容器、Sandbox、K8s、推理等模块详文；`../governance/module-delivery-workflow.md` §2.10 |
| P0-3 向量状态 | `available` → `ready`；去掉未冻结 search 422 | `vector-storage.md` |
| P0-4 安全域口径 | 租户管理从「待补」改为分流入口 | `security-identity-overview.md` |
| P0-5 治理诚实度 | 本表分项对齐，取消「全量 ✅」 | 本文件、`schema-completion-tracker.md` |
| P1 结构统一 | `## 创建前置条件` 全模块；容器状态矩阵；README backlog 链 | 见 §P1 记录 |
| P1 辅助材料 | `tasks/modules/prd/prd-console-*.md`、`tasks/modules/spec/spec-console-*.md`、TASK 清单与主维护源对齐 | §3.5、`spec-console-inference/security/tenant/container/vector/gpu` 等 |
| P0 backlog 模块 | 8 项 P0 详文 + PRD/SPEC + HTML | 2026-06-17 |
| P1 backlog 模块 | 18 项 P1 详文 + PRD/SPEC + HTML | 2026-06-17 |
| Phase 3 TASK 对齐 | Core 13 + Services 16 TASK，映射 Phase 2 全量 ops | 2026-06-17 |
| Phase 4 Core 文档 | `tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` + 模块详文/SPEC 加厚 | 2026-06-17 |
| P2 剩余 26 项文档 | AI 原生、知识智能、模型、安全、SDK/CLI | 2026-06-17 |
| P0/P1/P2（39 项）Core 标准加厚 | 详文 + PRD/SPEC 全章节；matrix 待收口→可用 | 2026-06-17 |
| Phase 2 handler 文档验收 | `PHASE2-HANDLER-ACCEPTANCE-RECORD.md` + Services 指南 | 2026-06-17 |
| P0 YAML 草案 | 标记 **待定** | 2026-06-17 |
| Agent Session Phase 3 详化 | `../openapi-drafts/phase3/openapi-phase3-agent-session-draft.md` | 2026-06-17 |
| Tool Permission / Agent Audit 详化 | `../openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md`、`../openapi-drafts/phase3/openapi-phase3-agent-audit-draft.md` | 2026-06-17 |
| Prompt Guard 详化 | `../openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md` | 2026-06-17 |
| AI 原生 ×7 详化完成 | MCP / Orchestration / Sandbox Security 详文 + PRD/SPEC + 草案 | 2026-06-17 |
| Phase 3 AI 原生文档验收 | `PHASE3-AI-NATIVE-ACCEPTANCE-RECORD.md` | 2026-06-17 |
| 知识库智能 ×3 详化 | doc/meeting/video 详文 + PRD/SPEC + 草案 | 2026-06-17 |
| Phase 3 KB 智能文档验收 | `PHASE3-KB-INTELLIGENCE-ACCEPTANCE-RECORD.md` | 2026-06-17 |
| 模型增强 ×3 详化 | encryption / recommend / usage-stats 详文 + PRD/SPEC + 草案 | 2026-06-17 |
| Phase 3 模型增强文档验收 | `PHASE3-MODEL-ENHANCEMENT-ACCEPTANCE-RECORD.md` | 2026-06-17 |
| 安全/合规扩展 ×4 详化 | api-key-audit / netsec / compliance / billing-export | 2026-06-17 |
| Phase 3 安全合规文档验收 | `PHASE3-SECURITY-COMPLIANCE-ACCEPTANCE-RECORD.md` | 2026-06-17 |
| Phase 3 整域文档层收口 | `../openapi-drafts/phase3/openapi-phase3-domain-draft.md` 全 § 详化 ✅ | 2026-06-17 |
| Phase 3 整域联评议程 | `tasks/phase3/PHASE3-JOINT-REVIEW-AGENDA.md`（2～3 场 · 待召开） | 2026-06-17 |
| **仓库目录规整** | governance / openapi-drafts / tasks 分子目录；根 [`README.md`](../../../README.md) | 2026-06-17 |

## P1 文档收口记录（2026-06-17）

| 项 | 动作 | 涉及文件 |
|---|---|---|
| 容器实例状态矩阵 | 对齐 VM/GPU 容器：`InstanceRecord.state` × lifecycle action | `container-instance-management.md` |
| 章节命名统一 | `## 前置条件` → `## 创建前置条件` | 网络、存储、API Key、推理、模型、知识库等 |
| 聚合页补章节 | 平台概览 / 开放集成 / GPU 算力：声明不承接直接创建 | `platform-overview.md`、`open-integration-overview.md`、`gpu-management.md` |
| 结构修复 | 网络/存储总览：POST 前置条件从「产品验收」段迁回创建章节 | `network-management.md`、`storage-management.md` |
| README | 链到 backlog + 维护规则 | `README.md` |

## P1 辅助材料收口（2026-06-17）

| 项 | 动作 | 涉及文件 |
|---|---|---|
| 推理 SPEC | 错误码改 §2.10 口径 | `spec-console-inference-service.md` |
| 安全 SPEC/PRD | 租户管理从「全待补」改为分流 | `spec-console-security-identity-overview.md`、`prd-console-security-identity-overview.md` |
| 租户 SPEC | 对齐主维护源创建前置条件 | `spec-console-tenant-management.md` |
| 容器/向量/GPU SPEC | 指向主维护源或补状态/路径口径 | `spec-console-container-*`、`vector-storage`、`gpu-management` |
| TASK 清单 | 实现验收用语标注 YAML 已举例/待冻结 | `CORE-TEAM-TASKS.md`、`SERVICES-TEAM-TASKS.md` |
| 治理规则 | PRD/SPEC 同步规则 | `../governance/module-delivery-workflow.md` §3.5 |

## 当前可信结论

- 当前 **65** 项 backlog 模块文档均已存在（P0×8 + P1×18 + P2×39）
- **文档层** 39 项待收口已加厚至 Core 标准；26 项从零创建完成
- **实现层** 仍依赖 OpenAPI/handler 补全（见各模块 TODO-YAML）
- 权威源对齐必须按上表分项查看，**禁止**再使用「20 模块全部 ✅」类表述
- 后续工作的重点：按 `console-undefined-features-backlog.md` 推进时，先确认依赖路径已在 YAML 冻结

## 当前治理状态

### 待重写

- 当前已清空

### 待收口

- 当前已清空（P1 结构统一已于 2026-06-17 完成）

### 可用模块统一复查

- P0 项已完成；P1 结构统一已完成（2026-06-17）
- 后续：按 `console-undefined-features-backlog.md` 推进时，先确认依赖路径已在 YAML 冻结

## 后续执行规则

- 后续不再用“已完成”作为唯一状态词
- 每次汇报必须明确写出：`可用 / 待收口 / 待重写`，以及权威源对齐分项（✅/⚠️）
- 后续新增修改任何模块时，先更新本表，再修改对应模块文档
- 前置条件表必须遵守 `../governance/module-delivery-workflow.md` §2.10
