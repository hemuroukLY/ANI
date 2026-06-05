# M1-METER-A — 计量 API local profile

完成日期：2026-06-03
对应 Sprint：Sprint 6（2026-08-01 ~ 08-15）
验证结果：`make test`、`make validate-architecture`、`make validate-core-api-compatibility`、`make validate-doc-entrypoints`、`make validate-sdk-beta`、`make validate-sdk-mock-smoke`、`git diff --check` 均为 EXIT 0

## 实现了什么

新增 Core 计量 API local profile：`GET /metering/usage` 返回租户用量聚合结果和 local profile 标记，`POST /metering/token-usage` 支持 Services/控制面上报模型 Token 用量并按 `idempotency_key` 去重。当前实现只记录本地计量事件并聚合 token 输入、输出和总量，`dev_profile.real_provider=false`，不代表真实计量后端、账单系统或 production ready。

## 关键文件改动

| 文件 | 新增/修改 | 说明 |
|---|---|---|
| `api/openapi/v1.yaml` | 修改 | 补齐 metering usage 响应字段并新增 token usage 上报 path/schema |
| `pkg/ports/metering.go` | 新增 | 新增 `MeteringService` port、usage query DTO 与 token usage report DTO |
| `pkg/adapters/runtime/local_metering_service.go` | 新增 | local profile token usage 上报、幂等去重、租户内聚合和 dev profile 标记 |
| `services/ani-gateway/internal/router/metering_resources.go` | 新增 | Gateway `/metering/usage` 与 `/metering/token-usage` 路由和响应映射 |
| `services/ani-gateway/internal/router/stubs.go` | 修改 | 移除旧 metering stub，避免与真实 local profile 路由重复注册 |
| `sdks/core/*` | 修改 | 由 `scripts/gen_sdk_alpha.py` 重新生成 Core SDK metadata/schema/operation 常量 |

## 完工标准达成

- [x] OpenAPI 契约先行，并保持 Core API v1 兼容性基线通过
- [x] 新增 `MeteringService` port，未暴露计量后端 SDK、账单系统对象或外部 provider 地址
- [x] 新增 local adapter，覆盖 token usage 上报、幂等重放和 usage 聚合
- [x] router 层和 adapter 层单元测试分开覆盖
- [x] 完整批次门禁通过

## 备注

本批次不新增 REAL-K8S-LAB guard，不触碰 Services 冻结目录，不宣称真实 metering backend、billing backend 或 production ready 已完成。实例 CPU/GPU/内存真实用量采集、计费报表、账单出账和 provider/live gate 需后续批次证明。
