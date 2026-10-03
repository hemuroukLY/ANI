# ANI Gateway OpenAPI 鉴权策略迁移方案 V3

> 状态：设计方案，尚未实施
>
> 基线：`pr/liangzai006/112` / `5860972141d70ba4d1616cd13a2185457422ef42`
>
> 日期：2026-08-20
>
> 目标：先交付可验证的 OpenAPI 鉴权功能 MVP，保持旧 Gateway 和旧 auth RPC 兼容；随后通过受控阶段平滑升级到 ANI-08 定义的内部零信任架构。

本文是迁移方案，不表示代码、数据库、部署配置或真实集群安全门禁已经完成。当前开发状态仍以 `ANI-DOCS-INDEX.md`、`ANI-06-开发计划.md` Section 零和 `repo/CURRENT-SPRINT.md` 为准。

---

## 一、决策摘要

本方案采用三层交付模型：

| 层级 | 目标 | 允许的部署范围 | 能否声明 production-ready |
|---|---|---|---|
| Functional MVP | 完成 OpenAPI policy、Principal、boundary 和 V2 授权功能 | 本地、单人开发、隔离 CI | 否 |
| Minimum Security Baseline | 限制 auth-service 暴露面，保护敏感 RPC | 受控共享开发/测试集群 | 否 |
| Zero Trust | mTLS、服务身份、逐 RPC caller ACL、证书轮换和 live gate | 多租户/生产候选环境 | 通过全部门禁后才可评审 |

固定结论：

1. 可以先实现功能 MVP，但不能把它部署到无隔离的共享多租户集群。
2. 旧 `ValidateToken`、`CheckPermission` 和现有 API Key RPC 在迁移期保持 wire compatibility。
3. 新 `CheckPermissionV2` 从第一天起不接收 `roles`，避免把错误信任边界固化进新契约。
4. 旧 Gateway 继续使用旧 RPC；新 Gateway 只对显式迁移的 operation 使用 V2。
5. V2 deny/error 不回退旧 `CheckPermission`，未迁移 operation 仅通过显式 legacy policy 继续运行。
6. platform Principal 的规范语义为 `credential_domain=platform` 且 `tenant_id` 为空。
7. 完整 mTLS 可以后置，但共享集群前至少要完成 auth-service NetworkPolicy、敏感 RPC 过渡保护和 production reflection 关闭。
8. production/multi-tenant 发布必须完成 ANI-08 的 mTLS 和调用方身份门禁，不能用共享密钥永久替代。

---

## 二、当前代码事实与风险边界

### 2.1 当前鉴权链路

```text
HTTP client
  -> ANI Gateway Auth middleware
  -> auth-service.ValidateToken
  -> Gateway request context: tenant_id/user_id/roles/scope
  -> ANI Gateway RBAC middleware
  -> auth-service.CheckPermission(tenant_id, user_id, roles, resource, action)
  -> handler
```

当前事实：

- Gateway 使用 `insecure.NewCredentials()` 连接 auth-service；
- auth-service gRPC Server 未配置 TLS credentials；
- Server interceptor 只有 logging 和 recovery，没有调用方认证；
- production-shaped manifest 中 auth-service 是 ClusterIP，但没有 auth-service 专用 ingress allowlist 证据；
- gRPC reflection 当前总是开启；
- `CheckPermission` 直接相信请求里的 `roles`，没有根据 `user_id` 查询数据库；
- `CreateAPIKey` 等敏感 RPC 与授权 RPC 暴露在同一个无调用方认证的服务上。

### 2.2 正确的风险表述

“任意 Pod 伪造 roles 调用 `CheckPermission` 就直接获得 Gateway HTTP 权限”不准确，因为该 RPC 单独只返回授权结果。

但当前暴露面仍然是安全缺口：

1. 任意可达工作负载可以调用 auth-service RPC；
2. `CheckPermission` 可以被当作可伪造的授权 oracle；
3. `CreateAPIKey` 接收调用者提供的 tenant、user 和 scopes，成功后返回原始 API Key；
4. 如果调用者知道有效 tenant UUID，可能形成“直接创建凭证 -> 调用 Gateway”的真实越权路径；
5. `IssueServiceToken` 仅靠静态 caller secret，且当前传输未加密。

### 2.3 ANI-08 的定位

ANI-08 的零信任章节是目标安全架构，不是当前完成声明。该文档附录明确将 TLS 1.3 全链路标记为 `Planned`。

因此本方案的交付状态必须诚实区分：

- JWT/OIDC/API Key/RBAC 功能链路存在，不等于内部零信任完成；
- 局部 NetworkPolicy live evidence 不等于所有 ANI 组件完成默认拒绝；
- Functional MVP 不得标记 security-ready、runtime-ready 或 production-ready。

---

## 三、目标架构中的两个身份

鉴权链路必须区分两个不同主体：

### 3.1 服务调用方身份 CallerIdentity

表示“谁在调用 auth-service RPC”。

```go
type CallerIdentity struct {
    ServiceName string
    TrustMode   string // local | transitional_token | mtls
    CertificateSerial string
}
```

最终来源必须是经过验证的 mTLS peer certificate/SAN，而不是请求消息中的 `caller_service` 字段。

### 3.2 产品访问主体 Principal

表示“谁在访问 ANI 产品 API”。

```go
type Principal struct {
    Kind             PrincipalKind   // user | service | api_key | sandbox
    CredentialScheme CredentialScheme
    CredentialDomain CredentialDomain // tenant | platform | sandbox
    TenantID         string
    SubjectID        string
    CredentialID     string
    PrincipalRef     string
}
```

固定结构约束：

- platform：`TenantID == ""`；
- tenant user/service/API Key：TenantID 必须是非零 UUID；
- sandbox：由 Gateway 本地验签，携带真实 tenant ID 和 instance claim；
- `PrincipalRef` 由 auth-service 签发，短期、不可伪造、不得由 Gateway 自行拼接。

CallerIdentity 与 Principal 不得混用：

```text
mTLS 证明：调用 RPC 的是 ani-gateway
JWT/API Key 证明：产品请求主体是某 user/service/api_key
数据库权限证明：该主体此刻被授予哪些 resource/action/boundary
```

---

## 四、OpenAPI 鉴权策略模型

### 4.1 operation 扩展

受保护 operation 使用结构化扩展：

```yaml
x-ani-authz:
  version: v1
  resource: quota
  action: read
  boundary: platform
  principal_kinds:
    - user
```

字段定义：

| 字段 | 必填 | 语义 |
|---|---|---|
| `version` | 是 | 当前仅允许 `v1` |
| `resource` | 是 | 稳定业务资源，不从 URL 猜测 |
| `action` | 是 | 稳定业务动作，不从 HTTP method 猜测 |
| `boundary` | 是 | `own|tenant|platform` |
| `principal_kinds` | 是 | `user|service|api_key|sandbox` 非空集合 |

OpenAPI `security` 继续决定凭证传输方案；`x-ani-authz` 决定产品主体、权限和数据边界。二者不能互相替代。

### 4.2 boundary 规则

```go
func domainAllowsBoundary(p Principal, required Boundary) bool {
    switch required {
    case BoundaryOwn, BoundaryTenant:
        return (p.CredentialDomain == DomainTenant && p.TenantID != "") ||
            (p.CredentialDomain == DomainSandbox && p.TenantID != "")
    case BoundaryPlatform:
        return p.CredentialDomain == DomainPlatform && p.TenantID == ""
    default:
        return false
    }
}
```

这只是身份域前置检查：

- `own` 仍由 handler/store 校验对象所有权；
- `tenant` 仍由 tenant filter、transaction context 和 RLS 隔离；
- `platform` 仍由平台事务、数据库角色和显式目标 tenant 参数约束；
- sandbox 仍执行本地 instance/path/method capability 检查。

### 4.3 权限数据来源

auth-service 使用现有 `roles.permissions` JSONB 结构：

```json
{
  "resource": "instances",
  "actions": ["read", "list"],
  "scope": "tenant"
}
```

匹配规则：

| required boundary | 允许的 permission scope |
|---|---|
| `platform` | 仅 `platform` |
| `tenant` | 仅 `tenant` |
| `own` | `own` 或同租户更宽的 `tenant` |

通配规则：

- `platform-admin` 的 `*/*/platform` 只在 platform boundary 生效；
- `tenant-admin` 的 `*/*/tenant` 只在 tenant/own boundary 生效；
- platform-admin 不自动变成 tenant Principal；
- tenant-admin/API Key/sandbox 不得进入 platform boundary；
- wildcard 只按受控数据库权限数据解释，Gateway 不自行扩展。

---

## 五、RPC 兼容设计

### 5.1 旧 RPC 冻结

迁移期不破坏以下接口：

```proto
rpc ValidateToken(ValidateTokenRequest) returns (common.v1.TenantContext);
rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
```

兼容要求：

- 旧字段编号、类型和错误语义保持不变；
- 旧 Gateway 继续调用旧 RPC；
- platform 旧链路继续保留当前零 UUID wire behavior；
- 旧 `CheckPermission.roles` 只服务于 legacy Gateway，不复制到 V2；
- 为旧 RPC 增加调用量、caller、listener 和结果指标；
- 旧 RPC 只能在迁移指标归零、旧 Gateway 实例归零后废弃。

### 5.2 新增 ValidatePrincipal

```proto
rpc ValidatePrincipal(ValidatePrincipalRequest) returns (PrincipalContext);

message ValidatePrincipalRequest {
  string credential = 1;          // raw JWT or API key
  string credential_scheme = 2;   // bearer | api_key
}

message PrincipalContext {
  string principal_kind = 1;      // user | service | api_key
  string credential_domain = 2;   // tenant | platform
  string tenant_id = 3;           // platform 必须为空
  string subject_id = 4;
  string credential_id = 5;
  string principal_ref = 6;       // auth-service 签发的短期 opaque handle
  int64 expires_at_unix = 7;
}
```

`PrincipalContext` 不返回用于授权决策的 roles。Gateway 可以记录 kind/domain/subject 进行审计，但不能据此自行授予业务权限。

### 5.3 新增 CheckPermissionV2

```proto
rpc CheckPermissionV2(AuthorizationRequest) returns (AuthorizationDecision);

message AuthorizationRequest {
  string principal_ref = 1;
  string resource = 2;
  string action = 3;
  string required_boundary = 4; // own | tenant | platform
  string operation_id = 5;      // audit/debug only
}

message AuthorizationDecision {
  bool allowed = 1;
  string reason_code = 2;
}
```

禁止字段：

- `roles`；
- 调用者自报的 `caller_service`；
- Gateway 自报的 permission wildcard；
- 用于覆盖 PrincipalRef 中 subject/domain/tenant 的字段。

### 5.4 PrincipalRef

Functional MVP 可以使用 auth-service HMAC 签发的短期 opaque handle，至少绑定：

- principal kind；
- credential domain；
- tenant ID；
- subject ID 或 credential ID；
- credential JTI/hash version；
- issued-at/expiry；
- 随机 nonce；
- schema version。

固定要求：

- 有效期不超过对应 credential 剩余有效期；
- API Key revoked/expired 状态仍需在授权时重新确认；
- role/permission 不写入 handle；
- handle 不能替代服务调用方身份；
- 生产零信任阶段可将 caller identity/audience 纳入 handle binding。

### 5.5 为什么不直接修改旧 CheckPermission

直接修改旧 RPC 会同时造成：

- 旧 Gateway 滚动升级不兼容；
- platform tenant ID 语义变化；
- 旧测试和 mock 全部变化；
- 无法区分 legacy roles 信任和新的服务端权威权限。

新增 V2 是为解决当前兼容问题，而不是提前建设通用授权框架。

---

## 六、服务端权威授权

### 6.1 user Principal

根据 PrincipalRef 中已验证的 subject/domain/tenant 查询：

```text
users
  -> user_roles
  -> roles
  -> roles.permissions
```

查询必须同时约束：

- tenant user：`users.tenant_id = principal.tenant_id`；
- platform user：`users.tenant_id IS NULL` 且 role scope 为 platform；
- user 未禁用/删除；
- role 绑定仍然有效。

不能只按 `subject_id` 查询后接受跨 tenant 角色。

### 6.2 API Key Principal

授权时重新读取或使用可安全失效的短 TTL cache 检查：

- key ID/hash；
- tenant ID；
- scopes；
- expires_at；
- revoked_at；
- instance binding；
- credential version。

API Key 的 scopes 不能通过 Gateway 的 roles 数组传入 V2。

### 6.3 service Principal

服务 JWT 必须重新验证：

- ANI 签名；
- issuer；
- audience；
- principal kind；
- tenant/domain；
- service scope；
- expiry/JTI blocklist。

若后续引入 service registry，再单独评审数据库化；本批次不为未来假设新增表。

### 6.4 sandbox Principal

sandbox token 是 Gateway 本地 capability credential，不调用 `CheckPermissionV2`：

1. 本地验签和过期检查；
2. generated policy 明确允许 sandbox；
3. boundary 仅允许 own/受控 tenant 子资源；
4. instance ID 与 path 完全相等；
5. method/path 在签名 capability 中；
6. 注入真实 tenant data context；
7. handler/store/RLS 继续做数据隔离。

---

## 七、Gateway 策略生成和运行时

### 7.1 构建期生成

新增生成器读取 Core 和 Services OpenAPI：

```text
api/openapi/v1.yaml
api/openapi/services/v1.yaml
  -> validate x-ani-authz
  -> compare registered routes
  -> generate operation policy table
  -> Gateway binary
```

运行时不解析 YAML。

生成失败条件：

- 未知 version/boundary/principal kind；
- resource/action 为空；
- 非 public route 同时缺 generated policy 和 legacy allowlist；
- operationId 重复；
- 注册 route 与 OpenAPI method/path 不一致；
- 生成物漂移。

### 7.2 路由解析

使用 Hertz route template/`FullPath()` 匹配，不对 raw path 编写新的正则授权器。

策略来源必须可观测：

```go
type PolicySource string

const (
    PolicySourceGenerated PolicySource = "generated"
    PolicySourceLegacy    PolicySource = "legacy"
    PolicySourcePublic    PolicySource = "public"
)
```

### 7.3 V2 授权流程

```text
ResolvePolicy
  -> public: pass
  -> legacy: old ValidateToken + old CheckPermission
  -> generated:
       ValidatePrincipal
       validate Principal structure
       credential/principal/boundary pre-check
       sandbox: local capability authorization
       others: CheckPermissionV2(principal_ref, resource, action, boundary)
       allow -> handler
       deny -> 403
       auth-service unavailable -> 503
```

禁止行为：

- V2 deny 后回退 legacy；
- V2 error 后回退 legacy；
- route 未命中 policy 时按 HTTP method/path 猜权限；
- platform 空 tenant ID 绕过 rate limit/audit；
- 将空 tenant ID 注入普通 tenant transaction context。

### 7.4 横切中间件

RateLimit、Idempotency 和 Audit 改用 Principal key：

| Principal | identity key |
|---|---|
| tenant user | `tenant:<tenant_id>:user:<subject_id>` |
| platform user | `platform:user:<subject_id>` |
| API Key | `tenant:<tenant_id>:api_key:<credential_id>` |
| service | `<domain>:service:<subject_id>` |
| sandbox | `tenant:<tenant_id>:sandbox:<instance_id>` |

现有幂等 key/body/fingerprint 语义不变，只调整主体命名空间。

---

## 八、三段式安全升级

### 8.1 Level 0：Functional MVP

目标：证明产品鉴权语义正确，不声称内部零信任完成。

允许环境：

- localhost；
- 单人 Docker Compose；
- 网络隔离的 CI fixture；
- 明确无不受信任工作负载的临时环境。

交付：

- `x-ani-authz` schema/validator/generator；
- Principal 和 platform 空 tenant ID；
- `ValidatePrincipal` / `CheckPermissionV2`；
- 服务端权威 permission evaluator；
- Gateway generated/legacy 双策略；
- quota-meta 试点；
- 完整正反测试矩阵；
- `GATEWAY_AUTHZ_POLICY_MODE=off|pilot|full`，默认 `off`。

限制：

- 只能标记 `local/logic verified`；
- 不运行在共享多租户集群；
- 不标记 ANI-08 zero trust completed；
- 不删除旧 RPC；
- 不使用 `full` 模式进行生产流量。

### 8.2 Level 1：Minimum Security Baseline

目标：允许进入受控共享开发/测试集群，但仍不是生产零信任。

最低交付：

1. auth-service Service 仅为 ClusterIP；
2. auth-service ingress default deny；
3. 只允许 Gateway Pod/namespace 访问 legacy/V2 auth port；
4. `IssueServiceToken` 的调用路径单独允许 inference-service；
5. production-like profile 关闭 reflection；
6. 临时内部调用凭证通过 K8s Secret 注入，不写入仓库或日志；
7. 新 Gateway 对敏感 RPC 发送临时 caller token；
8. auth-service 在旧 Gateway 清零后对敏感 RPC 强制校验 caller token；
9. 非 Gateway Pod 访问 `CreateAPIKey/ListAPIKeys/RevokeAPIKey/CheckPermissionV2` 必须失败；
10. 记录该方案仍未满足 ANI-08 mTLS，禁止外推 production-ready。

临时 caller token 不能单独作为安全边界，必须与 NetworkPolicy 组合；它在明文链路中不能抵抗节点级或特权 Pod 窃听，因此必须有删除期限。

### 8.3 Level 2：Zero Trust

目标：达到 ANI-08 定义的内部通信目标。

交付：

- cert-manager 内部 CA、Issuer/Certificate；
- Gateway client certificate 自动签发和轮换；
- auth-service server certificate 自动签发和轮换；
- TLS 1.3/mTLS；
- auth-service 从 peer certificate SAN/SPIFFE identity 解析 CallerIdentity；
- 每个 RPC 的 caller allowlist；
- `ValidatePrincipal/CheckPermissionV2/API Key RPC` 仅允许 Gateway；
- `IssueServiceToken` 仅允许明确的 inference-service identity；
- caller identity、certificate serial 和 auth decision 进入审计；
- NetworkPolicy 继续作为纵深防御；
- 证书轮换、错误证书、过期证书、错误 caller 和跨 namespace live gate。

达到以上条件前，ANI-08 的 TLS/零信任状态不能标记 `Verified`。

---

## 九、旧接口到 mTLS 的平滑迁移

### 9.1 双监听器过渡

为保持旧 Gateway 可用，transport 迁移采用双监听器，而不是一次性把 9101 改成 mTLS：

| listener | 用途 | 生命周期 |
|---|---|---|
| legacy `:9101` | 旧 Gateway 明文 gRPC | 受 NetworkPolicy 限制，迁移期保留 |
| secure `:9443` | 新 Gateway mTLS gRPC | Zero Trust 阶段启用并成为最终入口 |

固定顺序：

1. 先部署 auth-service 双 listener，旧 9101 行为不变；
2. 验证旧 Gateway 继续工作；
3. 部署支持 mTLS 的新 Gateway，切换到 9443；
4. 观察 listener/RPC/caller 指标；
5. 等待所有旧 Gateway 实例清零；
6. 关闭 9101 Service port 和 server listener；
7. 删除临时 caller token；
8. auth-service 只保留 mTLS listener。

禁止：

- 在旧 Gateway 尚存活时强制 9101 mTLS；
- 根据单次 RPC 失败自动降级到 9101；
- 新 Gateway 配置 mTLS 失败时静默使用 insecure；
- 将双 listener 作为永久状态。

### 9.2 Proto 兼容与 transport 兼容分离

旧 RPC 可以继续运行在 secure listener 上。因此关闭 legacy listener 不等于立即删除旧 Proto：

```text
阶段 1：旧 RPC + legacy listener
阶段 2：旧 RPC/V2 RPC + secure listener
阶段 3：仅 V2 RPC + secure listener
```

这允许先完成 transport zero trust，再独立观察和废弃旧业务 RPC。

### 9.3 配置 fail-closed

建议配置：

```text
AUTH_GRPC_LISTEN_MODE=legacy|dual|mtls
AUTH_GRPC_LEGACY_PORT=9101
AUTH_GRPC_MTLS_PORT=9443
AUTH_GRPC_SERVER_CERT_FILE=...
AUTH_GRPC_SERVER_KEY_FILE=...
AUTH_GRPC_CLIENT_CA_FILE=...

AUTH_SERVICE_TRANSPORT=insecure_local|transitional|mtls
AUTH_SERVICE_ADDR=...
AUTH_SERVICE_SERVER_NAME=...
AUTH_SERVICE_CLIENT_CERT_FILE=...
AUTH_SERVICE_CLIENT_KEY_FILE=...
AUTH_SERVICE_CA_FILE=...
```

约束：

- `insecure_local` 只允许 loopback/dev profile；
- production profile 检测到 `insecure_local` 必须 readiness 失败；
- `mtls` 配置缺证书/CA/server name 时启动失败；
- 禁止 `InsecureSkipVerify`；
- certificate reload/rotation 失败必须有指标和告警。

---

## 十、试点选择

### 10.1 首个试点

```text
GET /api/v1/admin/quota-meta
operationId: listQuotaMeta
```

选择理由：

- OpenAPI 已存在；
- Gateway route 已注册；
- QuotaAdminService 已有运行时注入；
- 只读；
- 明确 platform boundary；
- 能直接验证 platform Principal 无 tenant ID。

建议 policy：

```yaml
security:
  - BearerAuth: []
x-ani-authz:
  version: v1
  resource: quota
  action: read
  boundary: platform
  principal_kinds: [user]
```

验收矩阵：

| 场景 | 预期 |
|---|---|
| platform-admin，无 tenant ID | 200 |
| tenant-admin | 403 |
| tenant user | 403 |
| API Key | 403 |
| service JWT | 403 |
| 无凭证 | 401 |
| V2 deny | 403，不能 fallback |
| auth-service unavailable | 503，不能 fallback |

### 10.2 第二试点

`GET /api/v1/metering/usage/platform` 只有在以下前置完成后迁移：

- Gateway route 注册；
- persisted query adapter；
- platform transaction/DB role/RLS 边界；
- representative data；
- tenant_id 可选筛选的二次 platform 权限校验；
- 分页、时间边界、group_by 语义和测试。

OpenAPI 中存在路径不等于运行时闭环。

---

## 十一、实施批次

### AUTHZ-FUNC-A：策略格式与生成器

- `x-ani-authz` schema；
- Core/Services spec validator；
- operation policy generator；
- registered route coverage；
- generated drift gate；
- 默认 policy mode `off`。

### AUTHZ-FUNC-B：Principal 兼容层

- Gateway Principal 类型；
- 旧 TenantContext -> legacy Principal view；
- platform 零 UUID 在新内部语义规范化为空；
- legacy wire behavior 不变；
- principal-aware rate limit/idempotency/audit；
- tenant/platform/API Key/service/sandbox/dev 回归。

### AUTHZ-FUNC-C：V2 权威授权

- additive Proto；
- `ValidatePrincipal`；
- PrincipalRef；
- `CheckPermissionV2`；
- user role/permission store；
- API Key scope/status store；
- JWT/service evaluator；
- 明确禁止 V2 roles 字段。

### AUTHZ-FUNC-D：quota-meta 本地试点

- operation policy；
- pilot resolver；
- 无 tenant ID platform-admin E2E；
- 正反矩阵；
- deny/error no-fallback；
- 仅 local/isolated CI 标记通过。

### AUTHZ-BASELINE-E：共享测试集群最小保护

- auth-service default-deny/allow NetworkPolicy；
- ClusterIP-only guard；
- production reflection off；
- 临时 caller token；
- 敏感 RPC caller enforcement；
- 未授权 Pod negative test；
- 删除期限和 Zero Trust 迁移 owner。

### AUTHZ-ZEROTRUST-F：mTLS 双 listener

- cert-manager manifests；
- server/client TLS config；
- peer identity interceptor；
- per-RPC caller ACL；
- 9101/9443 dual listener；
- new Gateway 切换；
- certificate rotation tests。

### AUTHZ-ZEROTRUST-G：关闭 legacy transport

- 旧 Gateway 实例清零；
- 9101 调用指标清零；
- 删除 legacy Service port/listener；
- 删除临时 caller token；
- strict mTLS live gate；
- ANI-08 状态更新评审。

### AUTHZ-MIGRATION-H：产品路由全量迁移

- Core operation；
- Services operation；
- OpenAI-compatible proxy；
- sandbox capability routes；
- legacy route count 归零。

### AUTHZ-LEGACY-I：废弃旧业务 RPC

仅在以下条件全部满足后执行：

- 所有已注册产品路由都有 generated/public policy；
- legacy policy 指标在观察窗口为零；
- 旧 Gateway 实例为零；
- 旧 RPC 调用指标为零；
- tenant/platform/API Key/service/sandbox/dev/public 回归通过；
- rollback 演练通过。

---

## 十二、兼容矩阵

### 12.1 Gateway/auth-service 版本

| Gateway | auth-service | 结果 |
|---|---|---|
| 旧 Gateway | 旧 auth-service | 当前行为 |
| 旧 Gateway | 新 auth-service，legacy listener | 旧 RPC 行为不变 |
| 新 Gateway，policy off | 新 auth-service | 全部 legacy，等价当前行为 |
| 新 Gateway，pilot | 新 auth-service | quota-meta 使用 V2，其余 legacy |
| 新 Gateway，mtls | 新 auth-service，dual | secure listener；旧 Gateway 仍走 legacy listener |
| 新 Gateway，full | 新 auth-service，mtls only | 最终目标 |
| 新 Gateway，V2 unavailable | 任意 | 503，禁止自动 fallback |

### 12.2 Principal 行为

| 主体 | tenant_id | credential domain | 允许边界 |
|---|---|---|---|
| tenant user | 非零 UUID | tenant | own/tenant，依 permission |
| platform user | 空 | platform | platform，依 permission |
| API Key | 非零 UUID | tenant | own/tenant，依 scope |
| service | 按签名 claim | tenant/platform | 仅明确 service policy |
| sandbox | 非零 UUID | sandbox | 本地 capability 允许的 own/tenant |

### 12.3 安全成熟度

| 能力 | Functional MVP | Minimum Baseline | Zero Trust |
|---|---:|---:|---:|
| OpenAPI generated policy | 是 | 是 | 是 |
| 服务端权威权限 | 是 | 是 | 是 |
| NetworkPolicy | 非必需，仅隔离环境 | 必需 | 必需 |
| 临时 caller token | 否 | 必需 | 删除 |
| mTLS caller identity | 否 | 否 | 必需 |
| 逐 RPC caller ACL | 逻辑可先实现 | token identity | certificate identity |
| production-ready 候选 | 否 | 否 | 门禁通过后评审 |

---

## 十三、测试与验收

### 13.1 Functional MVP

- OpenAPI extension schema 正反测试；
- generator golden/drift；
- registered route coverage；
- legacy/generated/public resolver；
- platform empty tenant Principal；
- tenant/platform boundary 正反矩阵；
- PrincipalRef 篡改、过期和 credential 失效；
- DB permission exact/wildcard/scope；
- API Key revoked/expired；
- V2 request 不含 roles；
- sandbox 跨实例/跨路径/篡改；
- quota-meta pilot；
- V2 deny/error no-fallback；
- rate limit/audit 不因 platform tenant 为空而跳过。

### 13.2 Minimum Security Baseline

- auth-service 不是 NodePort/LoadBalancer；
- 非 Gateway Pod 无法连接 auth-service；
- 错误 caller token 调用敏感 RPC 被拒绝；
- 未认证直接 CreateAPIKey 被拒绝；
- reflection 在 production-like profile 不可用；
- 旧 Gateway rolling compatibility；
- NetworkPolicy 变更 rollback。

### 13.3 Zero Trust

- plaintext 连接 secure listener 失败；
- 无 client certificate 失败；
- 非受信 CA certificate 失败；
- 错误 SAN/service identity 失败；
- Gateway certificate 只可调用 allowlisted RPC；
- inference-service 只可调用 `IssueServiceToken`；
- 过期/吊销 certificate 失败；
- 证书轮换期间无中断；
- dual listener 期间旧/新 Gateway 并存；
- legacy listener 关闭后无 fallback；
- 跨 namespace/未授权 Pod live negative test。

### 13.4 回归命令

实施时至少执行：

```bash
cd repo

go test -count=1 ./services/ani-gateway/internal/middleware/...
go test -count=1 ./services/ani-gateway/internal/router/...
go test -count=1 ./services/auth-service/internal/service/...
go test -count=1 ./pkg/bootstrap/...

python scripts/validate_auth_gateway_contract.py
python scripts/validate_services_route_contract.py --root .
python scripts/validate_core_api_compatibility.py

make validate-services
make validate-architecture
make test
git diff --check
```

新增 Zero Trust 批次后还必须提供真实集群 live gate，至少证明未授权 Pod、错误证书和证书轮换行为。仅单元测试不能把 ANI-08 TLS 状态改为 `Verified`。

---

## 十四、可观测性

最低指标：

```text
gateway_authz_policy_total{source,operation_id,result}
gateway_authz_legacy_total{operation_id}
auth_rpc_total{listener,method,caller,result}
auth_rpc_unauthenticated_total{listener,method,reason}
auth_principal_ref_total{kind,result}
auth_permission_decision_total{kind,boundary,resource,action,result}
auth_grpc_certificate_expiry_seconds{identity,serial}
```

审计至少包含：

- request ID；
- operation ID；
- caller service identity；
- principal kind/domain/subject/credential ID；
- tenant ID（platform 允许空）；
- resource/action/boundary；
- allow/deny/error；
- policy source；
- auth listener/trust mode；
- certificate serial（mTLS 阶段）。

禁止记录：

- 原始 JWT；
- 原始 API Key；
- PrincipalRef 完整值；
- caller token；
- private key/certificate secret。

---

## 十五、回滚

### 15.1 Functional pilot 回滚

- 将 `GATEWAY_AUTHZ_POLICY_MODE` 从 `pilot` 改回 `off`；
- 回滚 Gateway deployment；
- 不删除 V2 RPC/Proto，避免旧新实例生成物不一致；
- 不对单个 deny/error 请求自动 fallback。

### 15.2 Minimum Baseline 回滚

- NetworkPolicy 回滚必须恢复“仅 Gateway 可达”，不能回滚为全开放；
- caller token enforcement 可在旧 Gateway 尚未清零时短暂关闭，但必须保留 NetworkPolicy；
- 记录临时风险窗口和关闭期限。

### 15.3 mTLS 回滚

- dual listener 阶段可将新 Gateway 暂时切回受 NetworkPolicy 限制的 legacy listener；
- mtls-only 阶段只有在保留经过审批的 rollback listener 时才能回退；
- 任何回退都必须告警并设定自动到期时间；
- 禁止长期保留 silent insecure fallback。

---

## 十六、预计改动范围

### OpenAPI 与生成器

- `api/openapi/v1.yaml`
- `api/openapi/services/v1.yaml`
- 新增 authz extension validator/generator
- generated Gateway policy table
- route coverage/drift gate

### Proto 和生成物

- `api/proto/auth/v1/auth_service.proto`
- `pkg/generated/pb/auth/v1/*`

### Gateway

- `services/ani-gateway/internal/middleware/auth_client.go`
- `services/ani-gateway/internal/middleware/auth.go`
- `services/ani-gateway/internal/middleware/rbac.go`
- `services/ani-gateway/internal/middleware/ratelimit.go`
- `services/ani-gateway/internal/middleware/idempotency.go`
- `services/ani-gateway/internal/middleware/audit.go`
- policy resolver/generated policy package
- TLS/caller-token transport config

### auth-service

- V2 RPC implementation
- PrincipalRef signer/validator
- permission store/evaluator
- API Key authoritative evaluator
- caller identity interceptor
- per-RPC allowlist
- dual listener/TLS config
- production reflection switch

### 部署与门禁

- auth-service/Gateway NetworkPolicy
- ClusterIP guard
- temporary caller token Secret wiring
- cert-manager Issuer/Certificate
- TLS Secret mounts
- shared-cluster baseline test
- zero-trust live gate/evidence

### 文档

实现批次按 Feature batch 规则更新：

- `repo/development-records/{batch}.md`
- `repo/development-records/README.md`
- `repo/CURRENT-SPRINT.md`
- `ANI-06-开发计划.md`
- ANI-08 状态仅在对应门禁完成后更新

---

## 十七、最终完成定义

### Functional MVP 完成

1. quota-meta pilot 在 local/isolated CI 通过；
2. platform 无 tenant ID 正常工作；
3. tenant/platform/API Key/service/sandbox 正反矩阵通过；
4. V2 不接收 roles；
5. auth-service 从权威数据计算权限；
6. 未迁移路由行为不变；
7. 不声明 security/production ready。

### Minimum Security Baseline 完成

1. auth-service ClusterIP-only；
2. default-deny + Gateway allowlist NetworkPolicy 生效；
3. 未授权 Pod 无法调用敏感 RPC；
4. 临时 caller token 生效且有删除期限；
5. production reflection 关闭；
6. 共享测试集群回归通过；
7. 不声明 ANI-08 zero trust Verified。

### Zero Trust 完成

1. Gateway/auth-service 全部使用 mTLS；
2. auth-service 从证书身份确认 caller；
3. 每个敏感 RPC 有 caller allowlist；
4. 旧 plaintext listener 和临时 caller token 已删除；
5. 证书轮换和错误身份 live gate 通过；
6. legacy Gateway/RPC 指标达到废弃条件；
7. ANI-08、当前 Sprint、开发记录与真实代码状态一致；
8. 完成独立 production readiness 评审。

