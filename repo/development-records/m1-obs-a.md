# M1-OBS-A — 可观测性 API local profile

完成日期：2026-06-03
对应 Sprint：Sprint 6（2026-08-01 ~ 08-15）
验证结果：`make test`、`make validate-architecture`、`make validate-core-api-compatibility`、`make validate-doc-entrypoints`、`make validate-sdk-beta`、`make validate-sdk-mock-smoke`、`git diff --check` 均为 EXIT 0

## 实现了什么

新增 Core 可观测性 API：`GET /observability/query` 作为 PromQL 代理查询入口，`/observability/alert-rules` 提供告警规则 CRUD local profile。当前实现只记录 Core 可观测性产品意图和本地规则状态，`dev_profile.real_provider=false`，不代表真实 Prometheus/Alertmanager provider 或 production ready。

## 关键文件改动

| 文件 | 新增/修改 | 说明 |
|---|---|---|
| `api/openapi/v1.yaml` | 修改 | 新增 observability query、alert-rules CRUD paths 与 schema |
| `pkg/ports/observability.go` | 新增 | 新增 `ObservabilityService` port 与 query/alert rule DTO |
| `pkg/adapters/runtime/local_observability_service.go` | 新增 | local profile PromQL 查询空结果与告警规则 CRUD/idempotency 状态机 |
| `services/ani-gateway/internal/router/observability.go` | 新增 | Gateway `/observability/*` 路由和响应映射 |
| `services/ani-gateway/internal/router/router.go` | 修改 | 注册 observability 路由 |
| `sdks/core/*` | 修改 | 由 `scripts/gen_sdk_alpha.py` 重新生成 Core SDK metadata/schema/operation 常量 |

## 完工标准达成

- [x] OpenAPI 契约先行，并通过 Core API v1 兼容性校验
- [x] 新增 `ObservabilityService` port，未暴露 Prometheus 地址或 SDK 对象
- [x] 新增 local adapter，覆盖 PromQL query 与告警规则 CRUD/idempotency
- [x] router 层和 adapter 层单元测试分开覆盖
- [x] 完整批次门禁通过

## 备注

本批次不新增 REAL-K8S-LAB guard，不触碰 Services 冻结目录，不宣称 Prometheus/Alertmanager real-provider 已接入。真实监控后端代理、规则下发和告警状态回读需后续 provider/live gate 证明。
