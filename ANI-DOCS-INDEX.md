# KuberCloud ANI · 文档导航与一致性矩阵

> 最后更新：2026-05-23
> 目的：让人类开发者和 AI 工具在 5 分钟内判断当前开发阶段、文档职责、下一步入口和闭环规则。

---

## 当前结论

```text
当前阶段：Phase 1 / Sprint 5 收敛中
当前不是 Phase 2：Phase 2 指 2026-10 以后延期能力
当前优先级：继续补齐 Sprint 5 真实 provider 主链路，并通过 REAL-K8S-LAB-A / make validate-real-k8s-profile 强制并行引入 K8s/Kube-OVN/KubeVirt/vCluster 真实底座验证环境
刚完成：Sprint 5 K8s CRUD+kubeconfig+proxy、Encryption keys+seal/unseal-token+rotate+revoke、Secret CRUD+bindings local profile，以及 WorkloadReconcileController bootstrap opt-in 运行剖面（2026-05-23）
下一步入口：repo/CURRENT-SPRINT.md（继续 Sprint 5，不能直接切换 Sprint 6）
```

本地真实代码显示：Sprint 4 API/SDK/Mock/Docs 收尾批次已归档；Sprint 5 目前完成可验证的 Core dev/local profile 主链路切片：`/api/v1/k8s-clusters` create/get/list/delete + kubeconfig + proxy，`/api/v1/encryption/keys` create/get/list/delete + seal + unseal-token + rotate + revoke，以及 `/api/v1/secrets` CRUD + bindings；`M1-RECONCILE-A` 已完成 controller adapter/capability 和默认关闭的 bootstrap opt-in 后台运行剖面。完整 vCluster 生命周期、真实 vCluster API 转发、controller leader election/指标/退避、真实 KMS/SM4 provider 和真实 K8s Secret 注入尚未由代码证明完成。

因此当前入口仍停留在 Sprint 5 收敛与后续切片，不进入 Sprint 6。从 Sprint 5 起，K8s、Kube-OVN、KubeVirt、vCluster、KMS/SM4、K8s Secret 注入等真实底座组件必须并行建设验证环境；local profile 只能证明 API/SDK/状态机/调用边界，不能证明真实组件已经跑通。`REAL-K8S-LAB-A` 是当前真实底座验证批次，默认通过 `make validate-real-k8s-profile` 校验门禁定义，三台云 VM 就绪后用同一脚本的 live 模式做真实验证。后续文档更新必须以真实代码、OpenAPI 契约、测试和真实环境验证记录共同落地为准。

`CORE-DEV-PROFILE-A` 是 Core dev/local profile 与 Services 业务 mock 的稳定边界名称：Core 可以提供合同兼容的本地开发剖面，但不得承载 Services 业务 mock。

`SPEC-SPLIT-A` 已完成：`/models`、`/inference-services`、`/knowledge-bases` 只保留在 Services API，Core API 和 Core SDK 不再承载这些业务路径。

`SPEC-CORE-BETA` 已完成首个切片：`repo/api/core-beta-readiness.yaml` 和 `make validate-core-beta` 用于持续校验 Core P0 path/schema、分页、幂等、状态机、RBAC scope 和 Core/Services 关联边界。

`SPEC-COMPAT-A` 已完成首个切片：`repo/api/core-v1-compatibility-baseline.yaml` 和 `make validate-core-api-compatibility` 用于持续保护 Core API v1 的 path/method/operationId/参数/响应/schema 字段，允许新增可选能力但阻止破坏性变更。

`SDK-BETA-A` 已完成首个切片：四语言 SDK 已生成 `idempotency_key` helper，并通过 `make validate-sdk-beta` 持续校验。

`SDK-BETA-B` 已完成首个切片：四语言 SDK 已生成 cursor 分页 helper，并在 SDK metadata 中标出支持 `limit/cursor` 的 Core 列表操作。

`SDK-BETA-C` 已完成首个切片：四语言 SDK 已生成统一 API error helper，并在 SDK metadata 中标出 API 契约声明的标准错误码。

`SDK-BETA-D` 已完成首个切片：四语言 SDK 已生成 basic example，覆盖 client 初始化、幂等、cursor 分页和 API error helper 的组合用法。

`SDK-MOCK-SMOKE-A` 已完成首个切片：Core Python SDK 已提供标准库 HTTP `request()` 能力，并通过 `make validate-sdk-mock-smoke` 调用由 API 契约驱动的 Core Mock Server。

`SDK-MOCK-SMOKE-B` 已完成首个切片：Core TypeScript SDK 已提供基于 `fetch` 的 `request()` 能力，并通过 `make validate-sdk-mock-smoke` 调用同一个 Core Mock Server。

`SDK-MOCK-SMOKE-C` 已完成首个切片：Core Go SDK 已提供基于 `net/http` 的 `Request()` 能力，并通过 `make validate-sdk-mock-smoke` 调用同一个 Core Mock Server。

`SDK-MOCK-SMOKE-D` 已完成首个切片：Core Java SDK 已提供基于 `java.net.http.HttpClient` 的 `request()` 能力；有 JDK 时调用同一个 Core Mock Server，无 JDK 时执行 source smoke。

`MOCK-A` 已完成首个切片：Core Mock Server 由 `repo/api/openapi/v1.yaml` 驱动，`make validate-mock-a` 校验全量 Core path 可 mock。

`DOC-API-A` 已完成首个切片：Core/Services 静态 API 文档由 API 契约生成到 `repo/docs/api/`，`make validate-doc-api` 校验 operation 和 schema 覆盖。

`SPRINT4-CLOSURE-A` 已完成首个切片：`make validate-sprint4-closure` 统一校验 Sprint 4 API/SDK/Mock/Docs/Records 关联性闭环。

---

## 唯一真实来源矩阵

| 问题 | 先看哪里 | 说明 |
|---|---|---|
| 当前做什么 | `repo/CURRENT-SPRINT.md` | 当前 Sprint 的执行入口，状态、任务、验收命令以它为准 |
| 全局开发节奏 | `ANI-06-开发计划.md` | Sprint 计划、Services 解锁门禁、延期项以它为准 |
| 产品功能边界 | `ANI-02-产品功能设计.md` | Core/Services 分层、v1.0.0 P0 能力边界以它为准 |
| 系统架构图和模块边界 | `ANI-05-系统架构设计.md` | Core/Services、API/SDK、ports/adapters、local profile/real provider 的结构图以它为准 |
| 路线图阶段 | `ANI-03-产品路线图.md` | Phase 1/2/3 与版本号关系以它为准 |
| 工程约定和 AI 工作规则 | `CLAUDE.md` | AI/人类开发前必须先读；只维护稳定规则和入口，不维护批次流水账 |
| API 契约 | `repo/api/openapi/v1.yaml` | Core OpenAPI REST API 与 Core/Services 跨层控制面契约的唯一真实来源 |
| API Beta 准备矩阵 | `repo/api/core-beta-readiness.yaml` | Core P0 Beta 审查、兼容性和自动校验矩阵 |
| API 兼容性基线 | `repo/api/core-v1-compatibility-baseline.yaml` | Core API v1 已交付 path/schema 的防破坏基线 |
| Services API 契约 | `repo/api/openapi/services/v1.yaml` | Services 层业务 API 契约 |
| 已完成批次 | `repo/development-records/README.md` | 历史完成记录索引，不作为当前任务清单 |
| 单批次细节 | `repo/development-records/*.md` | 追溯实现、验证和关键文件时再读 |
| 审查提示词模板 | `ANI-10-GPT审查提示词集.md` | 只作为审查问题模板；内置示例不得作为当前事实来源 |

---

## 推荐阅读路径

### 人类开发者

1. `ANI-DOCS-INDEX.md`
2. `CLAUDE.md` 的 5 分钟快速上手
3. `repo/CURRENT-SPRINT.md`
4. `ANI-06-开发计划.md` Section 零和当前 Sprint
5. `ANI-05-系统架构设计.md`
6. `repo/api/openapi/v1.yaml` + `repo/api/openapi/services/v1.yaml` + 相关代码入口

### AI 编码工具

1. 必须先读 `CLAUDE.md`
2. 再读 `repo/CURRENT-SPRINT.md`
3. 开发前检查 `ANI-06-开发计划.md` Section 零
4. 涉及架构边界时检查 `ANI-05-系统架构设计.md`
5. 涉及接口时先改 `repo/api/openapi/v1.yaml` 或 `repo/api/openapi/services/v1.yaml`
6. 完成后按 `CLAUDE.md` 的进度更新规约闭环

---

## 当前开发门禁

| 日期 | 门禁 | 当前影响 |
|---|---|---|
| 2026-05-31 | P0 依赖矩阵冻结 | 已完成历史批次归档，后续只按当前 Sprint 补缺口 |
| 2026-06-10 | Core API Alpha Freeze | 已完成 instances 等核心路径冻结；新增能力必须保持兼容性 |
| 2026-06-20 | SDK Alpha | 四语言 Core/Services SDK 已可生成，并由 SDK Beta/Mock smoke 持续校验 |
| 2026-06-30 | Core Dev Profile Ready | Core dev/local profile 边界已建立；Sprint 5 继续补真实 provider |
| 2026-07-31 | Core Real Path Beta | 当前关键门禁：K8s/Kube-OVN/KubeVirt/vCluster 真实底座验证和真实 provider 主链路 |
| 2026-09-30 | v1.0.0 Final Delivery | ANI Core v1.0.0 + ANI Services P0 |

---

## 文档维护规则

1. 当前阶段变更时，必须同步 `ANI-DOCS-INDEX.md`、`ANI-06-开发计划.md` 和 `repo/CURRENT-SPRINT.md`。
2. 批次完成时，必须新增或更新 `repo/development-records/{批次名}.md`，并更新 `repo/development-records/README.md`。
3. 历史归档文档允许保留当时日期和上下文，不反向改写为当前态。
4. 若 `CLAUDE.md` 与其它文档冲突，以 `CLAUDE.md` 的工程规则为准；若是进度状态冲突，以 `ANI-06-开发计划.md` Section 零和 `repo/CURRENT-SPRINT.md` 为准。
5. `CLAUDE.md` 只保留稳定强制规则、读取顺序、架构边界、提交门禁和 Karpathy 四条开发原则；禁止写入单批次完成清单、API path 长列表、文件级变更清单和每日开发流水账。
6. 动态进度只维护在 `repo/CURRENT-SPRINT.md`、`ANI-06-开发计划.md` Section 零和 `repo/development-records/*.md`；入口文档只保留当前状态、下一步和链接。
7. 更换 AI 模型或工具时，必须先重新读取本文件、`CLAUDE.md` 和 `repo/CURRENT-SPRINT.md`，不得依赖上一个会话的记忆。
8. 修改文档入口后必须运行 `make validate-doc-entrypoints`，确认 `CLAUDE.md` 没有重新承担动态进度记录职责。
