# M1-REGISTRY-A — 镜像仓库 API local profile

完成日期：2026-06-03
对应 Sprint：Sprint 6（2026-08-01 ~ 08-15）
验证结果：`make test`、`make validate-architecture`、`make validate-core-api-compatibility`、`make validate-doc-entrypoints`、`make validate-sdk-beta`、`make validate-sdk-mock-smoke`、`git diff --check` 均为 EXIT 0

## 实现了什么

新增 Core 镜像仓库 API local profile：`GET/POST /registry/projects`、`/registry/projects/{project}/repositories`、`/registry/projects/{project}/repositories/{repository}/artifacts`、`/registry/projects/{project}/repositories/{repository}/permissions`、`/registry/projects/{project}/pull-secret`、`/registry/projects/{project}/scan-report` 与 `/registry/images/scan-result` 覆盖项目创建/浏览、仓库/artifact 浏览、仓库权限设置、pull secret 引用创建和扫描结果/汇总查询。当前实现只返回确定性的本地 registry 元数据和 local scan 结果，`dev_profile.real_provider=false`，不代表真实 Harbor、Trivy provider 或 production ready。

## 关键文件改动

| 文件 | 新增/修改 | 说明 |
|---|---|---|
| `api/openapi/v1.yaml` | 修改 | 新增 registry projects/repositories/artifacts/permissions/pull-secret/scan-report/scan-result paths 与 schema |
| `pkg/ports/image_registry.go` | 修改 | 扩展 `ImageRegistry` port，覆盖项目、仓库、artifact、权限、pull secret 和扫描结果/汇总 DTO |
| `pkg/adapters/registry/local_image_registry.go` | 新增 | local profile 镜像仓库 adapter，支持租户项目、seeded repository/artifact、权限幂等、pull secret 引用和 local scan |
| `pkg/adapters/registry/not_configured.go` | 修改 | 补齐扩展后的 `ImageRegistry` not-configured adapter |
| `services/ani-gateway/internal/router/registry_resources.go` | 新增 | Gateway `/registry/*` 路由和响应映射 |
| `services/ani-gateway/internal/router/stubs.go` | 修改 | 移除旧 Harbor proxy stub，避免与 registry local profile 路由重复注册 |
| `sdks/core/*` | 修改 | 由 `scripts/gen_sdk_alpha.py` 重新生成 Core SDK metadata/schema/operation 常量 |

## 完工标准达成

- [x] OpenAPI 契约先行，并通过 Core API v1 兼容性校验
- [x] 复用并扩展 `ImageRegistry` port，未暴露 Harbor/Trivy SDK 对象或 provider 地址
- [x] 新增 local adapter，覆盖项目/仓库/artifact 浏览、权限幂等、pull secret 引用、scan result 和 project scan report
- [x] router 层和 adapter 层单元测试分开覆盖
- [x] 完整批次门禁通过

## 备注

本批次不新增 REAL-K8S-LAB guard，不触碰 Services 冻结目录，不宣称 Harbor、Trivy real-provider 或 production ready 已完成。真实 Harbor API 对接、真实镜像推拉凭证、扫描报告回读和 provider/live gate 需后续批次证明。
