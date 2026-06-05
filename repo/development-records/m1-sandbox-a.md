# M1-SANDBOX-A — Kata Containers 安全沙箱实例类型

完成日期：2026-06-03
对应 Sprint：Sprint 6（2026-08-01 ~ 08-15）
验证结果：`make test`、`make validate-architecture`、`make validate-core-api-compatibility`、`make validate-doc-entrypoints`、`make validate-sdk-beta`、`make validate-sdk-mock-smoke`、`git diff --check` 均为 EXIT 0

## 实现了什么

新增 Core `/instances` 的 `sandbox` 实例类型契约、`sandbox_config` 请求字段和响应摘要，补齐 `SandboxRuntime` port 与 local profile adapter。当前实现只表达 ANI sandbox/Kata 产品意图，并以 `dev_profile.real_provider=false` 标记本地状态机，不代表真实 Kata Containers provider 或生产 RuntimeClass 已验证。

## 关键文件改动

| 文件 | 新增/修改 | 说明 |
|---|---|---|
| `api/openapi/v1.yaml` | 修改 | `CreateInstanceRequest`/`InstanceRecord` 新增 `sandbox`、`instance_type` 兼容别名、`sandbox_config` 和 `SandboxInstanceStatus` |
| `pkg/ports/sandbox_runtime.go` | 新增 | 新增 `SandboxRuntime` port，只表达 runtime class、会话超时、出口策略和 local dev profile 意图 |
| `pkg/ports/workload_runtime.go` | 修改 | 新增 `WorkloadKindSandbox`、`WorkloadSpec.Sandbox` 和 `WorkloadInstanceRecord.Sandbox` |
| `pkg/adapters/runtime/local_sandbox_runtime.go` | 新增 | local sandbox profile，返回 pending/running 状态机和 non-real dev profile 标记 |
| `pkg/adapters/runtime/instance_service.go` | 修改 | `WorkloadInstanceService` 支持 sandbox create，写入统一 instance store 和 operation timeline |
| `services/ani-gateway/internal/router/demo_instances.go` | 修改 | `/instances` 请求映射 `sandbox_config`，响应返回 sandbox 摘要 |
| `sdks/core/*` | 修改 | 由 `scripts/gen_sdk_alpha.py` 重新生成 Core SDK metadata/schema 常量 |

## 完工标准达成

- [x] OpenAPI 契约先行，并通过 Core API v1 兼容性校验
- [x] 新增 `SandboxRuntime` port，未封装 Kubernetes/Kata SDK
- [x] 新增 local profile adapter，返回一致状态机和 `dev_profile.real_provider=false`
- [x] router 层和 adapter 层单元测试分开覆盖
- [x] 完整批次门禁通过

## 备注

本批次不新增 REAL-K8S-LAB guard，不触碰 Services 冻结目录，不宣称 sandbox 已 real-provider 或 production ready。真实 Kata RuntimeClass/provider 验证需后续 live gate 或真实 provider 批次证明。
