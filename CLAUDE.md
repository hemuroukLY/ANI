# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目背景

广州常青云科技有限公司旗下战略级新产品 **KuberCloud ANI**（KuberCloud AI-Native Infrastructure，中文名：AI专有云）的规划文档集，当前产品定义版本 **V7**。

**必读顺序（任何参与本项目的 Claude 实例，在开始任何任务前必须先读）：**

```
ANI-00  产品战略与开发哲学   ← 最高优先级，所有决策的出发点
ANI-01  客户画像与场景分析
ANI-02  产品功能设计
ANI-03  产品路线图
ANI-04  技术栈设计
ANI-05  系统架构设计
ANI-06  开发计划            ← 详细开发点拆解与月度里程碑
ANI-07  部署工程设计        ← Installer、多部署目标、Region/AZ、计量、白牌化
ANI-08  安全架构设计        ← 平台自身安全（第一层）+ 安全服务能力（第二层）
ANI-09  数据模型设计        ← PostgreSQL 全量表结构 / Milvus Collection / Redis Key / MinIO Bucket
ANI-10  GPT审查提示词集     ← 8组结构化审查提示词，用于 GPT 5.5 跨模型设计审查
ANI-11  代码实现规范        ← Repository接口、三条核心流程时序、Bootstrap模式、Go规范、前端路由
ANI-12  版本管理策略        ← SemVer、发布门禁、Git tag、升级包、兼容性边界
ANI-13  开源组件松耦合适配器架构 ← MinIO/Milvus/NATS/Redis/Harbor 等默认组件的 ports/adapters 强制边界
```

**代码仓库位置：** `repo/`（相对于本文档目录）

```
repo/
├── services/ani-gateway/      # Go，统一 Web Server（Hertz）
├── ai/rag-engine/             # Python，RAG 引擎
├── operators/inference-operator/ # Go，K8s Operator
├── frontends/{console,boss}/  # TypeScript，React 18 + TDesign
├── cli/ani/                   # Go，ani CLI
├── installer/ani-installer/   # Go，bubbletea TUI 安装程序
├── api/openapi/v1.yaml        # OpenAPI 3.1 Spec（先于代码）
├── deploy/docker/             # docker-compose 本地开发环境
├── Makefile                   # make deps / make dev-gateway / make test
└── .env.example               # 环境变量模板，cp 为 .env
```

## 开发阶段命名强制约定

为避免产品计划阶段与 AI 代码生成批次混淆，所有 Claude/Codex/GPT 实例必须遵守：

1. `ANI-06` 中的 `模块 1/2/3...` 是产品开发计划的唯一模块编号来源。
2. 代码生成批次不得再使用 `Stage 3A/3B/3C` 这类容易误解为模块 3 的名称。
3. 代码生成批次必须使用可回溯命名：`M{模块号}.{小节号}-{主题}-{批次}`，例如 `M2.1-TASK-A`。
4. 当前项目进度并行推进 `ANI-06 / 模块 1：基础设施底座` 与 `ANI-06 / 模块 2：ANI Gateway`：
   - `M1-INFRA-A` 已完成：基础设施代码化基线。
   - `M2.1-TASK-A` 已完成：最小 `task-service` 查询接口。
   - `M2.1-TASK-B` 已完成：transactional outbox + NATS publisher。
   - `M2.1-TASK-C` 已完成：worker mutation RPCs with tenant-safe writes。
   - `M2.2-AUTH-A` 已完成：最小 `auth-service` JWT validation + RBAC foundation。
   - `M1-INFRA-B` 已完成：component install profiles and infrastructure dependency contracts。
   - `M2.2-AUTH-B` 已完成：Gateway Auth/RBAC middleware wired to auth-service gRPC。
   - `M1-INFRA-C` 已完成：KubeOVN tenant VPC and NetworkPolicy isolation templates。
   - `M2.2-AUTH-C` 已完成：RLS-safe API Key lifecycle and validation。
   - `ARCH-ADAPTER-A / M1-ARCH-A` 已完成：开源组件松耦合适配器架构设计。
   - `ARCH-ADAPTER-B` 已完成：`pkg/ports` 与 `pkg/adapters` 能力接口骨架。
   - `ARCH-ADAPTER-GUARD-A` 已完成：组件 SDK 直接导入扫描与 allowlist 护栏。
   - `ARCH-ADAPTER-C` 已完成第一批迁移：auth JWT blocklist 使用 `CacheStore`，task outbox publish 使用 `MessageBus`。
   - `M2.2-AUTH-D` 已完成第一批：JWT `RevokeToken` 写入 `CacheStore` blocklist。
   - `M1-INFRA-D` 已完成：cluster preflight validation profile。
   - `ARCH-ADAPTER-C-2` 已完成：pgx/metadata 依赖按 bounded direct 分类，清理 auth-service 非必要 pgx 泄漏。
   - `M2.2-AUTH-E` 已完成：JWT blocklist 增加 PostgreSQL 持久化兜底与 CacheStore 快路径。
   - `M2.2-AUTH-F` 已完成：refresh token 持久化校验与 RS256 AccessToken 签发。
   - `M2.2-AUTH-G` 已完成：OIDC login begin/callback RPC 边界、state 缓存与授权 URL 构造。
   - `M2.2-AUTH-H` 已完成：OIDC code exchange、静态 RS256 ID token verifier、用户/角色映射与 refresh token 签发。
   - `M2.2-AUTH-I` 已完成：OIDC JWKS discovery / `kid` 公钥选择，静态公钥保留为离线 fallback。
   - `M2.2-AUTH-J` 已完成：OIDC group 到 ANI role 的显式映射与默认最小权限策略。
   - `M1-INFRA-E` 已完成：GPU scheduling baseline，覆盖 GPU 节点标签契约、Volcano Queue、HAMi/DCGM 契约与预检模板。
   - `M2.2-AUTH-K` 已完成：OIDC login -> refresh token -> ValidateToken 集成测试剖面。
   - `M1-INFRA-F` 已完成：GPU scheduling preflight/e2e hardening，新增可执行 Kubernetes Job、RBAC、严格检查开关和离线契约校验。
   - `M1-GPU-A` 已完成：异构 GPU 发现与调度契约，覆盖 NVIDIA/Huawei/Hygon、多型号、内核/驱动/运行时兼容和调度决策 port。
   - `M1-RUNTIME-A` 已完成：Workload Runtime / Instance 抽象，覆盖 VM、普通容器、GPU 容器、推理、Notebook、Agent Sandbox 和 Batch Job。
   - `M1-INSTANCE-A` 已完成：核心实例对象、全生命周期、网络平面与存储附件预置契约，明确 VM/Pod 可共享租户 VPC，同时保留 foundation mesh/storage/management 平面。
   - `M1-INSTANCE-B` 已完成：实例规划器 `PlanningRuntime` 最小实现，提供创建前网络/存储/GPU/生命周期校验与计划态记录。
   - `M1-INSTANCE-C` 已完成：Kubernetes/KubeVirt provider dry-run renderer，规划后输出可审查的 VM/Deployment/Job manifest，不直接创建集群资源。
   - `M1-INSTANCE-D` 已完成：本地 admission guardrail，审查 dry-run manifest 的类型、租户/实例标签、网络平面注解、hostNetwork 和 privileged 风险。
   - `M1-INSTANCE-E` 已完成：实例计划/渲染/准入结果持久化与审计，新增 `instance_plan_audits` RLS 表和 `WorkloadPlanAuditStore`。
   - `M1-INSTANCE-F` 已完成：provider dry-run executor 边界，新增 `WorkloadProviderDryRun` 与本地 provider/kind/apiVersion 校验，不创建资源。
   - `M1-INSTANCE-G` 已完成：provider apply/create 执行门控，新增 `WorkloadProviderApply` 与默认关闭的本地执行开关，强制校验用户、租户、权限证明、审计、admission 和 dry-run 证据。
   - `M1-INSTANCE-H` 已完成：实例状态回写/生命周期 reconcile 契约，新增 `WorkloadStatusReconciler`，强制 provider observation 与 apply/audit/resource refs 关联。
   - `M1-INSTANCE-I` 已完成：provider status reader 与实例创建编排 API，新增 `WorkloadProviderStatusReader` 和 `WorkloadInstanceOrchestrator`，业务层通过统一编排端口创建实例。
   - `M1-INSTANCE-J` 已完成：实例持久化/查询 API 契约，新增 `WorkloadInstanceStore`、`workload_instances` RLS 表和 orchestrator 状态写入。
   - `M1-INSTANCE-K` 已完成：Kubernetes/KubeVirt provider adapter 边界，新增 `KubernetesProviderAdapter` 和 `KubernetesProviderClient`，覆盖 server-side dry-run、受控 apply 和状态 observation。
   - `M1-INSTANCE-L` 已完成：实例服务 API 层，新增 `WorkloadInstanceService` 和 `LocalInstanceService`，对 VM、普通容器、GPU 容器提供 Create/Get/List 业务入口。
   - `M1-INSTANCE-M` 已完成：实例生命周期与可视化运维 API，补齐 Start/Stop/Restart/Resize/Delete 和 logs/events/metrics/terminal/exec ops 边界。
   - `M1-E2E-A` 已完成：M1 端到端集成剖面，覆盖 VM、普通容器、GPU 容器 create/lifecycle/query/ops 合同链路。
   - `M1-INSTANCE-N` 已完成：Kubernetes provider 执行剖面，覆盖 `KubernetesProviderClient` server-side dry-run、受控 apply、observe 与 orchestrator 集成。
   - `M1-INSTANCE-O` 已完成：adapter-owned `KubernetesRESTClient`，用标准库 HTTP 实现 dryRun=All、server-side apply 和 Deployment/Job/KubeVirt VM observe。
   - `M1-INSTANCE-P` 已完成：bootstrap/config provider wiring，支持 `WORKLOAD_PROVIDER=kubernetes_rest` 接入 `KubernetesRESTClient`，默认 local 且 apply 关闭。
   - `M1-INSTANCE-Q` 已完成：Kubernetes lifecycle execution，新增 `WorkloadInstanceLifecycleExecutor` 与 `KubernetesLifecycleExecutor`，覆盖 start/stop/restart/resize/delete provider 执行边界。
   - `M1-INSTANCE-R` 已完成：Kubernetes visual ops execution，新增 `KubernetesInstanceOps`，覆盖 logs/events/metrics/terminal/exec provider 执行边界。
   - `M1-E2E-B` 已完成：M1 real provider integration regression profile，统一覆盖 Kubernetes REST provider create/observe/lifecycle/ops 链路。
   - `DEMO-INSTANCE-CONSOLE-A` 已完成：阶段性实例 Demo Console，提前展示 VM、普通容器、GPU 容器 create/lifecycle/ops 体验；该层只允许调用 `WorkloadInstanceService`，不得绕过 M1 核心契约。
   - `M1-INSTANCE-S` 已完成：VM console/VNC/serial remote ops session 边界，支持 KubeVirt 与主流云/虚拟化 console 协议映射，不允许业务层直接接触 provider console API。
   - `DEMO-INSTANCE-WORKSPACE-UI-A` 已完成：实例 Demo 重构为生产控制台候选设计，覆盖 VM、普通容器、GPU 容器的创建、生命周期、运维与独立控制台页面。
   - `2026-05-12-demo-handoff` 已记录：当前 Demo 暂停点、启动步骤、mock 边界、演示验证命令；mock 只允许作为展示层，不代表对应生产 API 已完成。
   - 下一步规划已完成：`repo/development-records/2026-05-11-next-development-plan.md`。
   - 下一步建议：如需汇报环境真实创建资源，先补 `M1-DEMO-SMOKE-A`；否则开始衔接 `M3-MODEL-A`。
5. `Stage 3A/3B/3C` 只允许作为历史旧名出现，并必须注明“不代表模块 3”。
6. 任何进度更新必须同步写入 `repo/development-records/README.md`，并在对应设计文件中标注产品计划映射。

## 版本管理强制约定

所有发布、tag、镜像、Helm Chart、升级包和 Codex Cloud 交付分支必须遵守 `ANI-12-版本管理策略.md`：

1. ANI 采用 SemVer 2.0：`vMAJOR.MINOR.PATCH[-pre.N]`。
2. 首个正式版本为 `v1.0.0`，目标日期 `2026-09-30`，对应 Phase 1 POC 交付就绪。
3. 当前仍处于 `v0.x` 开发期，不得标记为 `v1.0.0` 或 `rc`。
4. 产品阶段、开发模块、代码生成批次与版本号互相独立：`模块 3`、`M2.1-TASK-C` 都不是 SemVer。
5. 任何 API/Proto/DB/CRD/Helm/安全模型的破坏性变更必须按 `ANI-12` 判断 MAJOR/MINOR/PATCH。

**三条核心原则（来自 ANI-00）：**
1. 产品完全从零构建，底层最大化复用成熟开源组件，ANI 的价值在于"编排"与"封装"
2. 全程利用 AI 大模型辅助编码加速开发，架构设计必须对 AI 辅助友好（Spec-First、强类型、单一职责）
3. 这是生产级平台，不是原型或玩具——稳定性、可扩展性、可观测性、安全性有明确的量化标准

**开源组件松耦合强制原则：**
- 除 Kubernetes API 外，业务服务不得直接依赖 MinIO、Milvus、NATS JetStream、Redis、Harbor、CloudNativePG 等开源组件 SDK。
- 业务服务必须依赖 ANI 自定义能力接口，默认组件只能出现在 adapter、bootstrap wiring、deploy profile 和集成测试中。
- `make test` 会执行 `make validate-architecture`，新增直接组件 SDK 导入必须先进入显式 allowlist 并标注迁移批次。
- 架构优先级是可用性/稳定性优先，其次性能可控性与扩展性；核心组件可采用 `bounded_direct`，但必须限定模块边界并记录理由。
- VM、普通容器、GPU 容器、推理、Notebook、Agent Sandbox、Batch Job 都必须先经过 `WorkloadRuntime` 能力抽象；模型部署不得直接把 Pod/Deployment/KubeVirt 细节写入业务流程。
- 异构 GPU 发现、分类和调度必须先经过 `GPUInventory` 能力抽象；厂商差异只能出现在 adapter、deploy profile、preflight 或 bounded runtime module。
- 所有实例必须声明网络平面和存储附件：租户业务互通走 `tenant_vpc`，平台控制/服务互联走 `foundation_mesh`，存储走 `storage`，运维入口走 `management`，公网暴露走 `public_ingress`。
- Provider renderer 只能输出 dry-run/review manifest；真实 apply/create 必须在后续受控 adapter 中接入 server-side dry-run、审计和权限检查后才允许。
- 所有 provider manifest 必须先通过 `WorkloadAdmission`；缺少 dry-run 标记、租户/实例标签、网络平面注解，或包含 hostNetwork/privileged 风险时必须拒绝。
- 实例计划、渲染 manifest、准入结果必须通过 `WorkloadPlanAuditStore` 持久化到租户 RLS 审计表，未审计不得进入真实 provider dry-run/apply。
- Provider dry-run 必须通过 `WorkloadProviderDryRun`；本地实现只校验 provider/kind/apiVersion 映射，未来 Kubernetes/KubeVirt 实现必须使用 server-side dry-run，仍不得绕过审计。
- Provider apply/create 必须通过 `WorkloadProviderApply`；默认执行开关关闭，真实 adapter 必须校验用户、租户、权限证明、审计 ID、admission、dry-run 结果和操作白名单，禁止业务服务直接 apply provider 资源。
- Provider 状态回写必须通过 `WorkloadStatusReconciler`；业务服务不得直接轮询 Kubernetes/KubeVirt/客户云状态 API，provider observation 必须携带 tenant、instance、audit 和 resource refs 关联证据。
- Provider status reader 必须通过 `WorkloadProviderStatusReader` 封装，业务服务创建实例必须优先使用 `WorkloadInstanceOrchestrator`，不得手动串联 renderer/dry-run/apply/status/reconcile 细节。
- 实例查询和恢复必须通过 `WorkloadInstanceStore` 或其上层服务，不能依赖 `PlanningRuntime` 内存状态；持久化记录必须包含 tenant、instance、kind、state、audit ID、provider ID 和 resource refs。
- Kubernetes/KubeVirt 真实执行必须通过 `KubernetesProviderAdapter` 及其内部 `KubernetesProviderClient`；server-side dry-run 必须使用 `dryRun=All`，apply 默认关闭，业务服务不得导入 provider SDK。
- 真实 KubernetesProviderClient 接入前必须先通过 `M1-INSTANCE-N` 执行剖面；后续只能替换 adapter-owned client 实现，不得改变 orchestrator、audit、admission、dry-run、apply、status reconcile 链路。
- 当前真实 KubernetesProviderClient 第一版为 adapter-owned `KubernetesRESTClient`，只使用标准库 HTTP；后续替换为 client-go/KubeVirt client 时仍必须留在 adapter-owned package。
- bootstrap 默认使用 local provider；只有显式设置 `WORKLOAD_PROVIDER=kubernetes_rest` 且提供 `KUBERNETES_API_HOST` 时才接入 `KubernetesRESTClient`，`WORKLOAD_PROVIDER_APPLY_ENABLED` 默认关闭。
- 已存在实例的真实生命周期操作必须通过 `WorkloadInstanceLifecycleExecutor`；`WORKLOAD_LIFECYCLE_PROVIDER=kubernetes_rest` 且 `WORKLOAD_LIFECYCLE_APPLY_ENABLED=true` 时才允许调用 Kubernetes/KubeVirt 生命周期 API。
- 实例可视化运维真实执行必须通过 `WorkloadInstanceOps` adapter；`WORKLOAD_OPS_PROVIDER=kubernetes_rest` 且 `WORKLOAD_OPS_ENABLED=true` 时才允许调用 Kubernetes logs/events/metrics/exec API。
- M1 真实 provider 路径必须通过 `M1-E2E-B` 回归剖面，统一覆盖 create、observe、lifecycle 和 ops adapter 链路。
- VM、普通容器、GPU 容器的业务创建/查询入口必须通过 `WorkloadInstanceService` 或其上层 gRPC/HTTP 服务，不得绕过实例服务 API 直接调用 provider adapter。
- VM、普通容器、GPU 容器生命周期操作必须通过 `WorkloadInstanceService`；可视化运维操作必须通过 `WorkloadInstanceOps`，业务服务不得直接调用 Kubernetes logs/events/metrics/exec 或 KubeVirt console/VNC API。
- VM 远程控制台必须建模为 `WorkloadInstanceOps` 会话；VNC、serial console、SPICE/RDP、OpenStack/VMware/公有云 console URL 都只能由 adapter 映射，业务层只接收 session metadata、connect URL、协议和过期时间。
- Demo/mock 页面只允许用于演示信息架构和前端体验；不得反向改变 M1/M2/M3 产品计划、端口/适配器边界或后端完成度判断。VM、普通容器、GPU 容器真实行为必须继续通过 `WorkloadInstanceService` 与 `WorkloadInstanceOps`。
- 任何违反该原则的实现必须停下来先补架构设计或 adapter 边界。

---

## Karpathy 四条开发原则

> 来源：Andrej Karpathy 关于 LLM 辅助编程的核心建议，整理自 [forrestchang/andrej-karpathy-skills](https://github.com/forrestchang/andrej-karpathy-skills)

### 原则一：先思考，再编码
**不要假设。不要掩饰困惑。要揭示取舍。**

- 如果需求有歧义，明确说出来并询问，而不是悄悄选一种猜测实现
- 存在多种合理方案时，列出并说明各方案的取舍，由人决策
- 面对复杂问题，先陈述理解再动手
- 遇到不确定的地方，停下来问，而不是带着疑惑继续

### 原则二：用能解决问题的最小代码
**拒绝一切带有猜想的实现。**

- 不实现没被要求的功能，哪怕"感觉以后用得到"
- 不为一次性代码创建抽象层
- 不加"灵活性""可配置性"等未被要求的扩展点
- 200 行能写成 50 行的，重写

### 原则三：只触碰你必须改动的部分
**只清理你自己制造的脏。**

- 不顺手"优化"任务范围之外的代码、注释或格式
- 不重构没坏的东西
- 保持现有代码风格，即使你有不同偏好
- 发现死代码，提出来，不要自作主张删除

### 原则四：定义成功标准，循环迭代直到验证通过
**把任务转化为可验证的目标。**

- 每个任务开始前明确"什么状态算完成"
- 多步骤任务先列出简短计划和验证步骤
- 优先给 Claude 成功标准而非操作指令：不是"修复这个 bug"，而是"写一个能复现 bug 的测试，再修复它"

---
