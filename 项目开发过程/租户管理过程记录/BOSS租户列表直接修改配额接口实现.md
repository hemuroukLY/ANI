# BOSS 租户列表「直接修改配额」接口实现设计

> 日期：2026-09-23　状态：待实施
> 关联代码基线：pr 分支（8bf1a1db 之后）

---

## 1. 背景与需求

当前 BOSS 租户列表（`/api/v1/svc/tenants*`）修改租户配额只有一条路：**配额变更审批流**——
`POST quota-requests`（提交 pending）→ `POST quota-requests/:reqId/approve`（审批后才 Upsert 到 Core）。

平台运营在紧急扩容等场景需要**跳过审批、直接改配额**。Core admin 面虽有
`PUT /api/v1/admin/tenants/{tenant_id}/quota`（quota_resources.go），但那是服务间控制面凭证域，BOSS 前端不能用。

**需求**：在 BOSS svc 面新增 `PUT /api/v1/svc/tenants/:tenantId/quota`，平台管理员直接改配额，立即生效，写审计。

### 非目标

- 不改动现有审批流（三条路径并存：审批流 / 套餐绑定同步 / 本接口直改）
- 不做维度删除（DELETE 仍走 Core admin 面）
- 不新增 DB 表 / migration

---

## 2. 方案选型

| 方案 | 说明 | 结论 |
|------|------|------|
| A. gateway svc handler 直转 Core admin 配额端点 | 薄透传，不经 tenant-service | ❌ 绕过租户存在性/状态校验与统一审计，BOSS 错误语义（TENANT_NOT_FOUND 等）丢失 |
| **B. svc handler → tenant-service gRPC 新方法 → Core UpsertQuota** | 与 BindPlanQuota / ReviewQuotaChangeRequest 同构，复用 `validateEnabledQuotaResourceTypes` + `applyTenantQuotaItems`，统一审计 | ✅ **采纳** |

选 B 的核心理由：租户校验（存在性、disabled 拒绝）、维度启用校验、Core 下发、审计全部复用 tenant-service 既有函数，编排层零新逻辑；Core 下发统一走 `UpsertQuota`（PUT upsert 端点，有则改无则建，单次原子）。

---

## 3. API 契约

### 3.1 端点

```
PUT /api/v1/svc/tenants/{tenantId}/quota
```

与既有 `GET /tenants/:tenantId/quota`（tenant_list_resources.go:48）同路径不同方法，路由无冲突。

### 3.2 请求

```json
{
  "items": [
    { "resource_type": "gpu_count", "total": 8 },
    { "resource_type": "cpu_core",  "total": 32 }
  ],
  "idempotency_key": "uuid，必填（可由 X-Idempotency-Key 头提供）"
}
```

- `items` ≥ 1 项，`resource_type` 批内不可重复（→ 422）
- `total` 为具体整数（≥ 0；Core 禁止 NULL，不支持 null 表示"回落默认"）
- 语义为 **upsert**：维度已有配额行则改 total，没有则新建（等价 Core `PUT quota/upsert`）

### 3.3 响应

```json
{
  "tenant_id": "…",
  "items": [
    { "resource_type": "gpu_count", "total": 8, "used": 2, "reserved": 0, "tightened": false }
  ]
}
```

- `tightened=true`：Core 因 `total < used+reserved` 自动收紧为当前占用值（**不算错误**，前端应提示）

### 3.4 错误码（复用 tenant_list 错误表，无新增码）

| HTTP | code | 触发 |
|------|------|------|
| 400 | VALIDATION_FAILED | tenantId 非 UUID / items 空 / total < 0 |
| 404 | TENANT_NOT_FOUND | Core GetTenant 不存在 |
| 409 | TENANT_STATE_INVALID | 租户 disabled（frozen 允许改，与 BindPlanQuota 对齐） |
| 422 | QUOTA_RESOURCE_NOT_REGISTERED | 维度未在 quota-meta 注册或已停用 |
| 409 | IDEMPOTENCY_KEY_REUSED | 幂等键复用（全局中间件） |
| 502 | GRPC_CLIENT_UNAVAILABLE / STORE_UNAVAILABLE | tenant-service 不可达 |
| 502/500 | （Core 错误透传） | Core 配额 API 失败 |

### 3.5 鉴权

与同组 svc 端点完全一致（`/api/v1/svc/*` 属 Services 过渡面，不走 Core policy registry，见 policy.go `IsCorePolicyPath` 豁免）。services/v1.yaml 中照抄 `reviewQuotaChangeRequest` 的 `x-ani-authz`（`resource: tenant, action: update, boundary: platform, principal_kinds: [user]`），即 **platform-admin / platform-ops** 可调用；platform-readonly 仅 GET。

---

## 4. 各层实现

### 4.1 OpenAPI 契约（api/openapi/services/v1.yaml）

在 `/tenants/{tenantId}/quota`（4280 行）的 `get` 之后追加 `put`：

```yaml
    put:
      operationId: updateTenantQuotaDirect
      summary: 直接修改租户配额（需 platform-admin / platform-ops，跳过审批流立即生效）
      description: |
        upsert 语义：维度已有配额行则改 total，无则新建。
        先 ListQuotaMeta 校验维度启用，再 Core PUT quota/upsert 单次原子下发。
        total < used+reserved 时 Core 自动收紧并在响应标 tightened=true。
      tags: [Tenants]
      security: [{ BearerAuth: [] }, { ApiKeyAuth: [] }]
      x-ani-authz:
        version: v1
        resource: tenant
        action: update
        boundary: platform
        principal_kinds: [user]
      # parameters / requestBody / responses 对齐 3.2/3.4；
      # 新增 schema UpdateTenantQuotaDirectRequest（items: [{resource_type, total:int64}] + idempotency_key）
      # 响应复用/扩展 TenantQuotaViewResponse 增加 reserved、tightened 字段
```

改后跑 `make gen-api gen-api-docs`（以及 svc 面契约校验脚本，若 CI 有 drift 校验则一并跑）。

### 4.2 proto（api/proto/tenant/v1/tenant_list_service.proto + tenant_plan.proto）

`tenant_list_service.proto` 追加消息（放在 ReviewQuotaChangeRequestRequest 之后）：

```proto
message UpdateTenantQuotaDirectRequest {
  string tenant_id = 1;
  repeated TenantQuotaTotalInput items = 2;   // min 1; unique resource_type per batch
  string idempotency_key = 3;                  // required UUID
}

message TenantQuotaTotalInput {
  string resource_type = 1;
  int64 total = 2;
}

message UpdateTenantQuotaDirectResponse {
  repeated TenantQuotaWriteItem items = 1;
}

message TenantQuotaWriteItem {
  string resource_type = 1;
  int64  total = 2;        // 生效后的值（tightened 时为收紧后的值）
  int64  used = 3;
  int64  reserved = 4;
  bool   tightened = 5;
}
```

`tenant_plan.proto` 的 `service TenantService` 追加：

```proto
  // UpdateTenantQuotaDirect writes tenant quota directly (BOSS fast path, no approval).
  rpc UpdateTenantQuotaDirect(UpdateTenantQuotaDirectRequest) returns (UpdateTenantQuotaDirectResponse);
```

生成：`make gen-proto`（产物 `pkg/generated/pb/tenant/v1`，含 grpc client/server 桩）。

### 4.3 tenant-service（编排 + 校验 + 审计）

`internal/service/tenant_service.go` 追加方法（样板 = BindPlanQuota 步骤 1-3 + ReviewQuotaChangeRequest 步骤 4）：

```go
// UpdateTenantQuotaDirect 直接修改租户配额（BOSS 快速通道，跳过审批流）。
// 校验链与 BindPlanQuota 对齐：维度启用 → 租户存在 → 非 disabled；下发一律 UpsertQuota。
func (s *TenantService) UpdateTenantQuotaDirect(ctx context.Context, req *tenantv1.UpdateTenantQuotaDirectRequest) (*tenantv1.UpdateTenantQuotaDirectResponse, error) {
	const action = "tenant.quota_update_direct"

	// 步骤 1：入参校验（tenant_id UUID / items≥1 / resource_type 批内去重 / total≥0）
	//         → VALIDATION_FAILED / QUOTA_CHANGE_REQUEST_INVALID 不复用，批内重复用 VALIDATION_FAILED 即可（新接口自定义，语义直白）

	// 步骤 2：维度启用校验（复用 validateEnabledQuotaResourceTypes，tenant_plan_service.go:816）
	//         → QUOTA_RESOURCE_NOT_REGISTERED

	// 步骤 3：经 Core 租户 API GetTenant（复用 s.tenants）
	//         不存在 → TENANT_NOT_FOUND；disabled → TENANT_STATE_INVALID（frozen 放行）

	// 步骤 4：下发（复用 applyTenantQuotaItems → core.UpsertQuota，单次原子）
	//         失败 → 审计 failure + mapStoreError 透传（Core 5xx→STORE/CORE 错误；TENANT_NOT_FOUND 原样）

	// 步骤 5：成功审计（details: items 下发值 + tightened 维度列表 + updated_by），
	//         审计失败只 Warn；返回 UpdateTenantQuotaDirectResponse（total 用 Core 返回的生效值）
}
```

要点：
- **不写 `quota_change_requests` 表**：本接口不产生审批记录（见 §5 语义决策）
- **不做异步补偿**：同步接口，失败即返回错误由前端重试（与审批流"状态已改不回滚"场景不同，这里没有任何本地状态需要先落库）
- `total` 直接取 `CoreQuotaResult.Total`（tightened 时是 Core 生效值），前端展示与 Core 一致

失败/成功审计均走既有 `writeAuditFailure` / `writeAuditSuccess`（resource=`tenant`，action=`tenant.quota_update_direct`）。

### 4.4 gateway svc handler（tenant_list_resources.go）

路由注册（48 行 `GET quota` 下一行）：

```go
svc.PUT("/tenants/:tenantId/quota", api.updateTenantQuotaDirect)
```

handler 样板 = `submitQuotaChangeRequest`（443 行起）：

```go
func (api *tenantListAPI) updateTenantQuotaDirect(ctx context.Context, c *app.RequestContext) {
	if api.tenants == nil { /* 502 GRPC_CLIENT_UNAVAILABLE */ }
	var body struct {
		IdempotencyKey string `json:"idempotency_key"`
		Items          []struct {
			ResourceType string `json:"resource_type"`
			Total        int64  `json:"total"`
		} `json:"items"`
	}
	// BindJSON → idem 兜底 idempotencyHeader(c) → 转 tenantv1 请求
	// → api.tenants.UpdateTenantQuotaDirect(callCtx, …) → mapTenantListError 透传
	// → 200 { tenant_id, items: [{resource_type,total,used,reserved,tightened}] }
}
```

- `mapTenantListError` 现有映射已覆盖全部新错误码（TENANT_STATE_INVALID=409、QUOTA_RESOURCE_NOT_REGISTERED=422 均在 `tenantListBusinessCodeByHTTP` 表中），**错误表零改动**
- `tenant_list_resources_test.go` 的 `fakeTenantListGRPC` 需补 `UpdateTenantQuotaDirect` 实现（编译必需）

### 4.5 Core（零改动）

复用 `PUT /api/v1/admin/tenants/{id}/quota/upsert`（quota_resources.go:33，`UpsertTenantQuota`），含 tightened 收紧语义与维度注册校验兜底。

---

## 5. 语义决策记录

1. **与 approved 锁定的交互（重要）**：审批流改过的维度会永久跳过"套餐绑定/改限额"的批量同步（`GetApprovedQuotaChanges` 过滤）。本接口**不写 `quota_change_requests`**，即：直接改过的维度**不会**获得该锁定，下次改套餐限额仍会被套餐值覆盖。
   - 理由：直改定位是运营快速通道，与审批流语义保持区分；若需"改完不再被套餐覆盖"，走审批流。
   - 若产品后续要求锁定，演进方案：直改成功后 INSERT 一条 `status=approved`、`requested_by=操作人` 的合成行（需评审，勿默认实现）。
2. **frozen 租户允许改配额**：与 BindPlanQuota 对齐（冻结只限制创建资源，不限制运营调整额度）；disabled 拒绝。
3. **tightened 是成功不是冲突**：Core 对 `total < used+reserved` 自动收紧而非报错，接口原样透出，由前端提示"已收紧至 N"。
4. **幂等**：全局 Idempotency 中间件按 `identity + method + path + key` 去重；Upsert 本身幂等，重复提交无副作用，无需业务层去重表。
5. **审计独立 action**（`tenant.quota_update_direct`），与审批流 `quota_change_request.approve` 区分，便于合规审计追溯"哪些改动没走审批"。

---

## 6. 实施清单

| # | 文件 | 改动 |
|---|------|------|
| 1 | `api/openapi/services/v1.yaml` | `/tenants/{tenantId}/quota` 追加 `put` + 请求/响应 schema |
| 2 | `api/proto/tenant/v1/tenant_list_service.proto` | 4 个新消息 |
| 3 | `api/proto/tenant/v1/tenant_plan.proto` | TenantService 追加 rpc |
| 4 | `make gen-proto` | 生成 `pkg/generated/pb/tenant/v1` |
| 5 | `services/tenant-service/internal/service/tenant_service.go` | `UpdateTenantQuotaDirect` 方法 |
| 6 | `services/tenant-service/internal/service/tenant_test.go` | 单测（§7.1） |
| 7 | `services/ani-gateway/internal/router/tenant_list_resources.go` | 路由 + handler |
| 8 | `services/ani-gateway/internal/router/tenant_list_resources_test.go` | fake 补方法 + handler 单测（§7.2） |
| 9 | `make gen-api gen-api-docs`（如有 drift 校验一并跑） | 文档/SDK 同步 |

无 migration、无 Core 改动、无新依赖。

---

## 7. 测试计划

### 7.1 tenant-service 单测（fake Core 客户端）

1. 成功 upsert：2 维度下发，fake 断言 `UpsertQuota` 调用 1 次、items 正确，审计 success 含 tightened 列表
2. 维度未注册 → `QUOTA_RESOURCE_NOT_REGISTERED`，不调 Upsert，审计 failure
3. 租户不存在 → `TENANT_NOT_FOUND`
4. 租户 disabled → `TENANT_STATE_INVALID`；frozen → 放行
5. items 空 / resource_type 重复 / tenant_id 非 UUID → `VALIDATION_FAILED`
6. Core 下发失败 → 错误透传 + 审计 failure（无本地状态需回滚）
7. tightened 回填：fake 返回 tightened + 收紧后 total，响应 total 用生效值

### 7.2 gateway handler 单测

1. 路由注册 + 200 透传（fake gRPC 返回含 tightened）
2. body 非法 / items 空 → 400 VALIDATION_FAILED
3. fake 返回 FailedPrecondition("TENANT_STATE_INVALID: …") → 409
4. fake 返回 FailedPrecondition("QUOTA_RESOURCE_NOT_REGISTERED: …") → 422
5. `api.tenants == nil` → 502 GRPC_CLIENT_UNAVAILABLE

### 7.3 E2E（本地三服务或 test3）

1. 建租户 → GET quota 基线 → PUT quota 改 gpu_count → GET 确认生效
2. PUT 不存在的维度名 → 422
3. PUT total < used 的维度 → 200 且 `tightened=true`、total=占用值
4. 审批流并存验证：直改后再走审批提交/批准，两互不干扰
5. 套餐覆盖验证（语义确认非 bug）：直改维度后改套餐限额 → 该维度被套餐值覆盖（§5.1）
6. disabled 租户 PUT → 409；审计日志出现 `tenant.quota_update_direct`

---

## 8. 风险与工作量

- **风险**：直改与套餐同步并存，运营可能困惑"改了又被套餐覆盖"（§5.1 已明示语义；前端可在改前展示当前绑定套餐的默认值作参考）。
- **工作量**：约 2 个新单测文件段 + 5 处既有文件小改，均为同构样板复制，无架构风险。
