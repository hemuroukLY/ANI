# ANI Gateway OpenAPI 鉴权策略改造方案 V4

> 文档状态：综合修订设计方案，尚未实施
>
> 版本：V4
>
> 修订日期：2026-08-20
>
> 综合来源：桌面 V2、仓库 V3，以及 2026-08-20 关于现有 gRPC 信任边界和 ANI-08 零信任规划的审查结论
>
> 首个纯鉴权试点：`GET /api/v1/admin/quota-meta`
>
> 第二试点：`GET /api/v1/metering/usage/platform`，必须等待平台计量查询运行时闭环完成
>
> 目标：以 OpenAPI 为产品路由鉴权策略来源，先完成可验证的功能迁移，保证旧接口和现有功能兼容；随后在不改变产品鉴权契约的前提下，平滑升级到 ANI-08 定义的内部零信任架构。

本文是设计方案，不表示代码、测试、数据库权限、NetworkPolicy、mTLS 或生产验证已经完成。当前状态仍以 `ANI-DOCS-INDEX.md`、`ANI-06-开发计划.md` Section 零和 `repo/CURRENT-SPRINT.md` 为准。

---

## 一、结论

本方案以 V2 的 OpenAPI 鉴权迁移为主线，并吸收 V3 中必要的安全边界修正。

按本文顺序实施后，应达到以下结果：

1. operation 显式迁移前继续走原有鉴权链路，现有 tenant、platform、API Key、service、sandbox 和 dev/local 行为保持不变。
2. 平台主体的规范语义为 `credential_domain=platform` 且 `tenant_id` 为空，不再把零 UUID 当作真实租户。
3. 旧 `ValidateToken`、`CheckPermission` 及其字段和值语义在迁移期保持兼容。
4. 新增 `ValidatePrincipal` / `CheckPermissionV2`，但 V2 不接收或信任 Gateway 提交的 `roles`。
5. auth-service 根据可信的 `subject_id` / `credential_id` 重新读取角色、权限或 API Key scopes。
6. OpenAPI operation 使用结构化 `x-ani-authz` 声明资源、动作、数据边界和允许的主体类型。
7. 构建期生成 Gateway 静态策略，运行时不解析 OpenAPI YAML，也不再根据 URL/HTTP method 猜测新 operation 权限。
8. 已迁移 operation 的 deny/error 不自动回退 legacy；未迁移 operation 仅通过显式 legacy 清单继续运行。
9. sandbox token 继续使用 Gateway 本地签名和 instance/path/method capability 约束，不进入通用角色授权。
10. 功能 MVP 可以先在本地/隔离环境完成；共享集群前必须完成最小安全基线；多租户/生产发布前必须完成 mTLS 和调用方身份验证。
11. Core、Services 和 OpenAI-compatible proxy 分域迁移；全部产品路由完成迁移前不删除全局 legacy 逻辑。

本方案明确区分：

- `security`：OpenAPI 接受的凭证方案；
- `principal_kind`：用户、服务身份、API Key 或 sandbox；
- `credential_domain`：凭证主体域，取值 `tenant|platform|sandbox`；
- `boundary`：operation 数据边界，取值 `own|tenant|platform`；
- `resource + action`：产品权限；
- Caller Identity：哪个内部服务正在调用 auth-service；
- handler/store/RLS：对象所有权、租户过滤和数据面隔离。

`x-ani-authz` 是 ANI 私有 OpenAPI 扩展，不是行业标准字段。通用实践是认证与授权分离、显式策略、服务端权威权限、默认拒绝、纵深防御和可审计迁移。

---

## 二、兼容性与现有功能影响

### 2.1 最终影响结论

兼容机制按本文顺序实施后：

- tenant 用户现有 API 行为不变；
- tenant API Key 行为不变，仍不能访问 platform boundary；
- platform 管理员可以在没有 tenant ID 的情况下访问被明确授权的 platform operation；
- platform 管理员不会自动获得 tenant/own operation 访问资格；
- sandbox token 继续只能访问绑定实例的允许子资源；
- 登录、refresh、health、branding 等公开/基础路径行为不变；
- `ANI_AUTH_MODE=dev` 保留于 dev/local profile，production profile 禁止启用；
- 未迁移 operation 的 legacy resource/action 结果保持原样；
- transport 从明文升级 mTLS 时，旧 Proto 仍可继续运行，不要求同时删除旧 RPC。

这些结论只有在兼容矩阵和回归矩阵全部通过后才成立。

### 2.2 平台管理员为什么可以没有 tenant ID

平台主体访问平台或跨租户边界，不属于任何单一租户。规范 Principal：

```json
{
  "principal_kind": "user",
  "credential_scheme": "bearer",
  "credential_domain": "platform",
  "tenant_id": "",
  "subject_id": "<platform-user-id>"
}
```

固定规则：

- `credential_domain=platform` 时 `tenant_id` 必须为空；
- `credential_domain=tenant` 时 `tenant_id` 必须是非零 UUID；
- `credential_domain=sandbox` 时必须携带合法 tenant ID 和已验签 instance claim；
- platform handler 不注入伪造的 `types.TenantContext`；
- 平台访问指定租户资源时，目标 tenant ID 来自 path/query，并单独经过 platform permission 和数据边界校验。

### 2.3 为什么不能直接修改旧 ValidateToken

旧链路依赖 platform token 返回零 UUID：

1. 旧 Gateway 的 `withTenantContextStrict` 要求 tenant ID 可解析为 UUID；
2. 旧 `CheckPermission` 要求 `tenant_id` 非空；
3. platform-admin 当前依赖角色直通；
4. RateLimit、Audit、Idempotency 和 request context 仍读取旧字段。

若 auth-service 先把旧 `ValidateToken.TenantContext.tenant_id` 改为空，旧 Gateway 和未迁移 platform operation 会在滚动升级期间出现 401/403。

因此：

- 旧 `ValidateToken` 保持当前返回语义，包括 platform 零 UUID；
- 旧 `CheckPermission` 保持 wire contract 和 legacy 判定语义；
- 新 `ValidatePrincipal` 返回规范 Principal，platform tenant ID 为空；
- 新 `CheckPermissionV2` 使用规范 identity/boundary，并从服务端权威数据读取权限；
- 只有明确迁移的 operation 使用 V2；
- 旧 Gateway、旧 RPC 和 legacy policy 分别观察、分别退役。

### 2.4 禁止的兼容方式

禁止：

- V2 deny/error 后自动调用旧授权；
- 把空 platform tenant ID 注入普通 tenant data context；
- 把零 UUID 继续当作新 Principal 的规范身份；
- 根据角色名猜 credential domain/principal kind；
- 根据 URL 前缀猜新 operation boundary；
- 在 auth-service 未全量支持 V2 时开启 pilot；
- 让 platform-admin 自动进入 tenant/own boundary；
- mTLS 配置失败后静默回退明文；
- 将临时 caller token 描述成最终零信任能力。

---

## 三、当前代码事实与风险边界

### 3.1 当前执行链路

```go
h.Use(
    RequestID(),
    AuthWithClient(authClient),
    RBACWithClient(authClient),
    RateLimit(store),
    Idempotency(store),
    Audit(),
)
```

当前存在三套可能不一致的权限描述：

| 来源 | 当前用途 | 问题 |
|---|---|---|
| OpenAPI `security` / `x-ani-rbac-scope` | 契约、SDK、文档和兼容基线 | Gateway 未消费精确权限语义 |
| `scopeAllowedForPath` | tenant/platform/sandbox 隔离 | 依赖路径和手工清单 |
| `inferPermission` | resource/action 推断 | 依赖 HTTP method 和路径首段 |

### 3.2 平台 token 与 principal kind

- 当前 JWT 允许 `scope=platform` 且无真实 tenant ID；旧 `ValidateToken` 会把 `uuid.Nil` 输出为零 UUID。
- 当前用户 JWT 没有完整可信的 `principal_kind` 兼容闭环。
- API Key kind 只能来自 API Key 校验分支。
- sandbox kind 只能来自 Gateway 本地 token 格式和签名校验。
- 缺少 `principal_kind` 的现有登录/刷新 JWT 迁移期按 `user` 处理。
- service 必须有 ANI 签名的 `principal_kind=service` 和预期 audience，不能根据角色名推断。

### 3.3 权限数据

数据库 migration 已定义 `{resource, actions, scope}`，但当前 `CheckPermission` 没有读取 `roles.permissions`，而是信任请求中的 `roles` 并执行角色直通/字符串 scope 匹配。

因此新 V2 必须满足：

- Gateway 不提交 roles；
- auth-service 按 subject/credential 查询权限；
- permission scope 与 required boundary 严格匹配；
- platform-admin/tenant-admin wildcard 只在各自 boundary 生效；
- auditor/精确权限使用同一 evaluator。

### 3.4 当前内部 gRPC 信任边界

当前事实：

- Gateway 使用 `insecure.NewCredentials()`；
- auth-service gRPC Server 没有 TLS credentials；
- interceptor 只有 logging/recovery，没有调用方认证；
- gRPC reflection 总是开启；
- production-shaped manifest 没有证明 auth-service 只允许 Gateway 访问；
- `CreateAPIKey` 等敏感 RPC 与授权 RPC 暴露在同一服务。

正确风险表述：

- 任意 Pod 伪造 roles 调用 `CheckPermission` 本身只会得到 true/false，不会自动获得 Gateway HTTP 权限；
- 但任意可达 auth-service 的工作负载可以调用授权和敏感 RPC；
- `CreateAPIKey` 接受 tenant/user/scopes 并返回原始 key，可能形成直接创建凭证后访问 Gateway 的路径；
- 这是现有代码的安全债务，不是 V2/V4 新增的问题。

### 3.5 ANI-08 状态

ANI-08 的“所有内部 gRPC 使用 mTLS”是目标安全架构。其附录仍将 TLS 1.3 全链路标记为 `Planned`。

因此：

- 功能 MVP 可以先行；
- 共享集群不能完全无隔离；
- production/multi-tenant 前必须完成 mTLS、调用方身份和逐 RPC allowlist；
- JWT/RBAC 功能 verified 不等于内部零信任 verified。

### 3.6 platform metering 与 Gateway 全局范围

`GET /metering/usage/platform` 已在 Core OpenAPI 中存在，但 Gateway 当前未注册该 route，且只装配 `LocalMeteringService`，不能作为第一批纯鉴权试点。

同一 Gateway 中间件还覆盖 Core、Services、OpenAI-compatible proxy 和基础设施路径。因此只迁移 Core OpenAPI 后不能删除全局 legacy。

---

## 四、目标鉴权模型

### 4.1 OpenAPI 声明

首个试点：

```yaml
/admin/quota-meta:
  get:
    operationId: listQuotaMeta
    security:
      - BearerAuth: []
    x-ani-rbac-scope: "scope:quota:read"
    x-ani-authz:
      version: v1
      resource: quota
      action: read
      boundary: platform
      principal_kinds: [user]
```

第二试点在数据面前置完成后：

```yaml
/metering/usage/platform:
  get:
    operationId: getPlatformMeteringUsage
    security:
      - BearerAuth: []
    x-ani-rbac-scope: "scope:metering:platform:read"
    x-ani-authz:
      version: v1
      resource: metering
      action: read
      boundary: platform
      principal_kinds: [user]
```

字段：

| 字段 | 必填 | 说明 |
|---|---:|---|
| `version` | 是 | 当前仅允许 `v1` |
| `resource` | 是 | 稳定业务资源，不从 URL 推断 |
| `action` | 是 | 精确业务动作，不从 HTTP method 推断 |
| `boundary` | 是 | `own|tenant|platform` |
| `principal_kinds` | 是 | `user|service|api_key|sandbox` 非空集合 |

OpenAPI Security Requirement 的 OR/AND 结构必须保留；不支持的组合应在生成期失败，不能静默放宽。

### 4.2 credential domain 与 boundary

| 概念 | 取值 | 用途 |
|---|---|---|
| credential scheme | bearer、api_key、sandbox_bearer | 凭证传输/验证方式 |
| principal kind | user、service、api_key、sandbox | 产品主体类型 |
| credential domain | tenant、platform、sandbox | 主体所在域 |
| required boundary | own、tenant、platform | operation 数据范围 |

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

该检查不能替代 handler/store/RLS。

### 4.3 permission scope

| required boundary | 允许的 permission scope |
|---|---|
| `platform` | 仅 `platform` |
| `tenant` | 仅 `tenant` |
| `own` | `own` 或同租户更宽的 `tenant` |

附加规则：

- platform-admin wildcard 仅 platform；
- tenant-admin wildcard 仅 tenant/own；
- platform-admin 不自动作为 tenant Principal；
- tenant API Key、tenant-admin、sandbox 不得进入 platform；
- `operation_id` 仅用于审计/诊断；
- wildcard 只能来自受控权限数据。

### 4.4 sandbox capability

sandbox token 不进入 auth-service 通用角色授权：

1. Gateway 本地验签、过期和 token 类型；
2. generated policy 允许 sandbox；
3. boundary 为 own 或受控 tenant 子资源；
4. path instance ID 等于 token claim；
5. method/path 在 capability 中；
6. 注入真实 tenant data context；
7. handler/store/RLS 继续隔离。

全部通过即构成本地授权决定，不调用 `CheckPermissionV2`。

---

## 五、策略生成与路由解析

### 5.1 构建期生成

```text
api/openapi/v1.yaml
  -> scripts/generate_gateway_authz.py
  -> services/ani-gateway/internal/authz/zz_generated_core_policies.go

api/openapi/services/v1.yaml
  -> scripts/generate_gateway_authz.py
  -> services/ani-gateway/internal/authz/zz_generated_services_policies.go
```

```go
type SecurityRequirement struct {
    AllOf []SecurityScheme
}

type Policy struct {
    OperationID          string
    Method               string
    PathTemplate         string
    Public               bool
    SecurityAlternatives []SecurityRequirement
    Resource             string
    Action               string
    Boundary             Boundary
    PrincipalKinds       []PrincipalKind
}
```

生成失败条件：

- 缺字段或未知 version/boundary/principal kind；
- operationId/method/path 冲突；
- Security Requirement 语义无法表示；
- 已注册产品 route 同时缺 generated policy 和 legacy entry；
- 生成物漂移。

### 5.2 路由解析

运行时使用 Hertz `RequestContext.FullPath()`：

1. Hertz 完成 route match；
2. 读取 method + FullPath；
3. 将 `:param` 规范化为 OpenAPI `{param}`；
4. 查 generated/public/legacy policy；
5. 已注册产品 route 无策略时 503 fail closed。

未知路由继续由 NoRoute 返回 404，不对 raw path 再实现授权 matcher。

### 5.3 legacy 清单

```go
type LegacyPolicy struct {
    Method       string
    PathTemplate string
    Owner        string
    Reason       string
    MigrateBy    string
}
```

约束：

- 当前未迁移路由可生成初始 inventory；
- 新 route 不得加入 legacy；
- legacy 数量只能下降；
- 单请求 deny/error 不 fallback；
- 紧急回滚只能通过配置/部署整体切回已审核 legacy policy，并记录审计。

---

## 六、Gateway 新链路

### 6.1 中间件顺序

```go
h.Use(
    RequestID(),
    ResolveAuthzPolicy(registry),
    AuthenticatePrincipal(authClient),
    AuthorizePrincipal(authClient),
    RateLimitByPrincipal(store),
    IdempotencyByPrincipal(store),
    AuditPrincipalDecision(),
)
```

每个请求按 policy source 进入 generated/public/legacy 分支，不能先失败再尝试另一套授权。

### 6.2 Principal

```go
type Principal struct {
    Kind             PrincipalKind
    CredentialScheme SecurityScheme
    CredentialDomain CredentialDomain
    TenantID         string
    SubjectID        string
    CredentialID     string

    // 仅旧 ValidateToken/legacy policy 使用；不得发送给 CheckPermissionV2。
    LegacyRoles      []string
    SandboxClaims    *SandboxClaims
}
```

字段规则：

- user：SubjectID 为 user ID；
- service：SubjectID 为签名 service subject；
- API Key：CredentialID 为 key ID；
- sandbox：SandboxClaims 必填；
- platform：TenantID 为空；
- tenant/sandbox：TenantID 为合法非零 UUID；
- LegacyRoles 仅 legacy 路径可读取。

### 6.3 认证

1. public policy 直接进入公开处理；
2. 同时出现 Bearer/API Key 时拒绝；
3. sandbox bearer 由 Gateway 本地验签；
4. generated operation 的 Bearer/API Key 调用 `ValidatePrincipal`；
5. 严格校验 Principal 结构；
6. tenant/sandbox 注入 tenant data context；platform 不注入；
7. legacy operation 继续使用旧 `ValidateToken`，生成旧 request-context view。

### 6.4 授权

```go
func AuthorizePrincipal(authClient AuthClient) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        resolved := MustResolvedPolicy(c)
        if resolved.Public {
            c.Next(ctx)
            return
        }
        if resolved.Source == PolicySourceLegacy {
            authorizeByLegacyPolicy(ctx, c, authClient, resolved.Legacy)
            return
        }

        policy := resolved.Generated
        principal := MustPrincipal(c)
        if !policy.AllowsCredential(principal) ||
            !policy.AllowsPrincipalKind(principal.Kind) ||
            !domainAllowsBoundary(principal, policy.Boundary) {
            respond403(c, "principal not allowed by operation policy")
            return
        }
        if err := evaluateCredentialConstraints(c, principal, policy); err != nil {
            respond403(c, "credential constraints not satisfied")
            return
        }
        if principal.Kind == PrincipalSandbox {
            c.Next(ctx)
            return
        }

        result, err := authClient.CheckPermissionV2(ctx, AuthorizationRequest{
            Principal:        principal.WithoutLegacyRoles(),
            Credential:       MustRawCredential(c),
            CredentialScheme: principal.CredentialScheme,
            Resource:         policy.Resource,
            Action:           policy.Action,
            RequiredBoundary: policy.Boundary,
            OperationID:      policy.OperationID,
        })
        if err != nil {
            respond503(c, "AUTHORIZATION_UNAVAILABLE", "authorization service unavailable")
            return
        }
        if !result.Allowed {
            respond403(c, normalizedDenyReason(result))
            return
        }
        c.Next(ctx)
    }
}
```

错误语义：无效凭证 401；策略/权限不匹配 403；auth-service 不可用 503；注册路由缺 policy 503。

### 6.5 横切中间件

Principal key：

| Principal | identity key |
|---|---|
| tenant user | `tenant:<tenant_id>:user:<subject_id>` |
| platform user | `platform:user:<subject_id>` |
| API Key | `tenant:<tenant_id>:api_key:<credential_id>` |
| service | `<domain>:service:<subject_id>` |
| sandbox | `tenant:<tenant_id>:sandbox:<instance_id>` |

平台空 tenant ID 不得绕过 rate limit。Idempotency 保留现有 key/body/fingerprint 语义，仅替换主体命名空间。Audit 记录成功、拒绝和 auth-service error，但不记录原始凭证。

---

## 七、内部鉴权契约与权威权限

### 7.1 保留旧 RPC

```proto
rpc ValidateToken(ValidateTokenRequest) returns (common.v1.TenantContext);
rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
```

迁移期：

- 不改旧 tag/类型/wire 语义；
- 旧 Gateway 继续工作；
- platform 零 UUID legacy 行为保留；
- 旧 roles 仅被旧 CheckPermission 消费；
- 增加 legacy RPC 调用指标；
- 废弃必须独立评审。

### 7.2 新增 V2 RPC

```proto
rpc ValidatePrincipal(ValidatePrincipalRequest) returns (PrincipalContext);
rpc CheckPermissionV2(AuthorizationRequest) returns (AuthorizationDecision);

message PrincipalContext {
  string principal_kind = 1;      // user|service|api_key
  string credential_scheme = 2;   // bearer|api_key
  string credential_domain = 3;   // tenant|platform
  string tenant_id = 4;           // platform 必须为空
  string subject_id = 5;
  string credential_id = 6;
}

message AuthorizationRequest {
  PrincipalContext principal = 1;
  string resource = 2;
  string action = 3;
  string required_boundary = 4;
  string operation_id = 5;
  string credential = 6;          // raw JWT/API key；只用于服务端重验，禁止记录
  string credential_scheme = 7;   // bearer|api_key
}

message AuthorizationDecision {
  bool allowed = 1;
  string reason_code = 2;
}
```

V2 明确没有 `roles` 字段。`CheckPermissionV2` 必须重新验证 credential，并确认验证得到的 kind/domain/tenant/subject/credential ID 与请求中的 PrincipalContext 一致；不一致直接拒绝。这样无需新增 principal handle，也不会只凭 Gateway 自报 subject/credential 作出授权。

权限查询和最终 decision 必须使用 credential 重验得到的服务端 Principal；请求中的 PrincipalContext 只用于一致性检查和 Gateway/auth-service 语义对齐。

L0 只允许在本地/隔离环境使用。L1 起还必须确认调用者是受控 Gateway，L2 再由 mTLS peer identity 提供最终调用方保证。无论处于哪一层，roles/scopes 都由 auth-service 重新读取，不能由 Gateway 提交。

重复的 JWT 签名验证可以在 auth-service 内部使用有界短 TTL cache 优化，但不能改变吊销、过期和身份一致性语义。原始 credential 不得进入日志、指标、trace attribute 或错误文本。

### 7.3 服务端权威 evaluator

user：

```text
users -> user_roles -> roles -> roles.permissions
```

同时约束 tenant/domain、用户状态和 role scope。

API Key：按 credential ID/hash 重新读取 tenant、scopes、expires_at、revoked_at、instance binding 和 credential version。

service：重新验证 ANI 签名、issuer、audience、principal kind、service scope、expiry/JTI。当前没有 service registry 时不为未来假设新增表。

固定规则：

- 不相信 Gateway roles；
- 不只按 subject 查询而忽略 tenant/domain；
- DB/内部错误返回稳定 reason/error，不泄漏 SQL；
- 短 TTL cache 必须支持吊销/版本失效；
- permission exact/wildcard/scope 使用同一 evaluator。

### 7.4 caller identity

Caller Identity 表示内部 RPC 调用者，与产品 Principal 分离：

```text
Caller Identity: ani-gateway 调用了 CheckPermissionV2
Principal: platform user X 正在访问 quota/read/platform
```

功能 MVP 阶段 Caller Identity 尚未达到零信任；最小安全基线使用 NetworkPolicy + 临时调用凭证；最终由 mTLS peer certificate/SAN/SPIFFE identity 提供。

逐 RPC allowlist：

| RPC | 最终允许 caller |
|---|---|
| ValidatePrincipal / CheckPermissionV2 | ani-gateway |
| Create/List/RevokeAPIKey | ani-gateway |
| Login/Refresh/RevokeToken | ani-gateway |
| IssueServiceToken | 明确 allowlist 的 inference-service |

---

## 八、试点与迁移范围

### 8.1 首个试点：GET /admin/quota-meta

选择理由：OpenAPI 和 Gateway route 已存在；handler 通过 `QuotaAdminService`；只读；明确 platform boundary；不依赖缺失的 metering adapter。

| 调用方 | 预期 |
|---|---:|
| 无凭证 | 401 |
| tenant user/admin | 403 |
| tenant API Key | 403 |
| platform service JWT | 403 |
| 无 quota/read/platform 权限的 platform user | 403 |
| platform-admin，无 tenant ID | 200 |
| auth-service 不可用 | 503 |

### 8.2 第二试点：GET /metering/usage/platform

进入条件：

- Gateway route 已注册；
- platform query port 语义明确；
- PG adapter 和生产装配完成；
- tenant_id 筛选经过二次平台权限/边界校验；
- group_by、UTC/时间边界与 OpenAPI 一致；
- 非超级用户数据库角色和 RLS 行为已验证；
- 旧 platform 鉴权下的数据面 E2E 先通过。

### 8.3 全量迁移域

| 域 | 策略来源 | 删除 legacy 条件 |
|---|---|---|
| Core `/api/v1/*` | Core OpenAPI | 注册产品路由全部覆盖 |
| Services `/api/v1/svc/*` | Services OpenAPI | security 和 route baseline 收敛 |
| OpenAI-compatible `/v1/*` | 正式独立/Services 契约 | 产品路由显式策略 |
| health/readiness/no-route | 基础设施策略 | 与产品 API 分离 |

Services 缺少 security 不能自动解释为 public。

---

## 九、当前实施顺序与后续演进

当前批准实施范围是 OpenAPI 鉴权 Functional MVP。若需要进入共享 Kubernetes 开发/测试集群，再追加 L1 Minimum Security Baseline。

L2 Zero Trust、`AUTHZ-ZEROTRUST-Z1/Z2` 和相关 live gate 仅用于证明本方案保留平滑演进路径，属于后续独立安全方案，当前不实施，也不阻塞 local/隔离环境的功能 MVP。

### 9.1 三个成熟度级别

| 级别 | 交付内容 | 允许环境 | 状态声明 |
|---|---|---|---|
| L0 Functional MVP | OpenAPI policy、Principal、V2 权威授权、单路由 pilot | local/隔离 CI | local/logic verified |
| L1 Minimum Security Baseline | auth-service 网络隔离、敏感 RPC 临时调用保护 | 受控共享开发/测试集群 | security baseline only |
| L2 Zero Trust | mTLS、服务身份、逐 RPC ACL、轮换/live gate | 多租户/生产候选 | 门禁通过后评审 Verified |

L0 可以先实现；L1 是进入共享集群的前置；L2 是 production/multi-tenant 的硬门禁。

### AUTHZ-POLICY-A：策略格式和生成门禁

- `x-ani-authz` validator；
- Core policy generator；
- OR/AND security；
- generated drift；
- FullPath/OpenAPI template；
- legacy inventory；
- mode 保持 off。

### AUTHZ-COMPAT-B0：Gateway Principal 兼容层

- Principal 类型；
- 旧 TenantContext -> legacy Principal view；
- platform 零 UUID 新内部规范为空，legacy wire 保持；
- Principal-aware rate limit/idempotency/audit；
- `GATEWAY_AUTHZ_POLICY_MODE=off|pilot|full`，默认 off；
- 仍只走旧 RPC。

### AUTHZ-CONTRACT-B1：V2 权威授权

- additive Proto；
- ValidatePrincipal；
- CheckPermissionV2，不含 roles；
- JWT principal_kind/audience；
- API Key credential ID；
- permission evaluator；
- auth-service 同时支持旧/V2；
- Gateway mode 仍为 off。

### AUTHZ-PILOT-C：quota-meta 功能试点

- operation policy；
- pilot allowlist 只包含 listQuotaMeta；
- platform 无 tenant ID E2E；
- 正反权限矩阵；
- deny/error no-fallback；
- local/隔离 CI 先完成。

### AUTHZ-BASELINE-S0：共享集群最小安全基线

在 pilot 进入共享 Kubernetes 集群前完成：

1. auth-service 保持 ClusterIP-only；
2. auth-service ingress default deny；
3. 只允许 Gateway 访问通用 auth RPC；
4. inference-service 仅通过受控路径访问 IssueServiceToken；
5. production-like profile 关闭 reflection；
6. 新 Gateway 支持通过 K8s Secret 注入临时 caller credential；
7. 旧 Gateway 清零后，敏感 RPC 强制 caller credential；
8. 未授权 Pod/错误 caller 调用 CreateAPIKey 和 CheckPermissionV2 必须失败；
9. 临时控制有 owner、删除条件和期限。

NetworkPolicy 无法区分 RPC method，必须与应用层 caller check 组合。临时 caller credential 在无 TLS 时不能抵抗节点级/特权 Pod 窃听，不能用于生产零信任声明。

### AUTHZ-METERING-D：数据面与第二试点

先完成 8.2 数据面进入条件，再迁移 getPlatformMeteringUsage。数据面 PR 与鉴权切换应独立评审。

### AUTHZ-MIGRATION-E：Core 全量功能迁移

- 逐资源补齐 x-ani-authz；
- 不顺便批量改 action；
- legacy count 单调下降；
- route/policy CI；
- 可在 L1 受控测试集群验证，但不能声明生产 ready。

### AUTHZ-SERVICES-F：Services 和 proxy

- Services security 先收口；
- 生成 Services policy；
- 处理 route baseline；
- proxy 有正式策略来源；
- 通过 `make validate-services`。

### AUTHZ-ZEROTRUST-Z1：mTLS 与调用方身份（后续规划，当前不实施）

- cert-manager CA/Issuer/Certificate；
- Gateway client certificate；
- auth-service server certificate；
- TLS 1.3/mTLS；
- peer identity interceptor；
- 每个 RPC caller allowlist；
- certificate rotation 和 negative tests；
- NetworkPolicy 保留为纵深防御。

### AUTHZ-ZEROTRUST-Z2：安全传输切换（后续规划，当前不实施）

transport 迁移必须与 Proto 迁移解耦。推荐两种实现之一，在 Z1 批次冻结：

1. 原生 gRPC 双 endpoint：旧明文 endpoint 保留给旧 Gateway，新 mTLS endpoint 给新 Gateway；
2. service mesh permissive -> strict mTLS：先兼容旧流量，再按 workload identity 收紧。

具体端口和配置名不在本设计提前固定，但顺序必须是：

```text
auth-service 先兼容 legacy + secure
  -> 验证旧 Gateway
  -> 新 Gateway 切 secure
  -> 观察 caller/listener/RPC 指标
  -> 旧 Gateway 和明文调用归零
  -> 关闭明文入口
  -> 删除临时 caller credential
```

禁止请求级 silent insecure fallback。

### AUTHZ-LEGACY-G：删除旧业务逻辑

只有以下条件全部满足后执行：

- Core、Services、proxy 注册产品路由全部 generated/public；
- legacy policy 指标归零并经过观察窗口；
- 旧 Gateway 实例归零；
- 明文 transport 调用归零；
- 旧 RPC 调用归零；
- tenant/platform/API Key/service/sandbox/dev/public 回归通过；
- route coverage/drift/兼容矩阵为 CI 硬门禁；
- 回滚演练通过。

关闭明文 transport 不等于必须立即删除旧 Proto。旧 RPC 可以先运行在 mTLS 上，再独立退役。

---

## 十、兼容矩阵

### 10.1 业务 RPC 版本

| Gateway | auth-service | mode | 预期 |
|---|---|---|---|
| old | old | legacy | 当前行为 |
| new B0 | old | off | 当前行为 |
| old | new B1 | legacy | 旧 RPC 行为不变 |
| new B0 | new B1 | off | 当前行为 |
| new pilot | new B1 全量实例 | pilot | 仅试点 V2 |
| new full | new B1+ | full | 仅在全量功能门禁后允许 |

### 10.2 transport 版本

| caller | auth-service transport | 结果 |
|---|---|---|
| 旧 Gateway | legacy endpoint/permissive mesh | 迁移期继续工作 |
| 新 Gateway | legacy endpoint | L0/L1 受控环境临时允许 |
| 新 Gateway | mTLS endpoint/strict identity | 最终目标 |
| 非 Gateway workload | 任意敏感 RPC | L1 起拒绝；L2 由网络+证书双重拒绝 |
| mTLS 配置错误 | secure endpoint | fail closed，不回退明文 |

### 10.3 Principal

| 主体 | tenant ID | domain | 允许 boundary |
|---|---|---|---|
| tenant user | 非零 UUID | tenant | own/tenant，依 permission |
| platform user | 空 | platform | platform，依 permission |
| API Key | 非零 UUID | tenant | own/tenant，依 scope |
| service | 按可信 claim | tenant/platform | 仅明确 service policy |
| sandbox | 非零 UUID | sandbox | capability 允许的 own/tenant |

---

## 十一、测试与验收

### 11.1 现有功能回归

- tenant/platform 登录与 refresh；
- 普通 Core GET 和副作用写；
- platform-admin 无 tenant ID 调用 quota/tenant/plan 管理 API；
- tenant user/admin 拒绝 platform operation；
- platform-admin 拒绝 tenant/own operation；
- API Key tenant 正向和 platform 反向；
- service-only 拒绝普通 JWT；
- sandbox 绑定实例正向及跨实例/路径/篡改反向；
- public health/login/branding；
- dev/local；
- auth-service 故障 503 且无 fallback；
- platform 请求受限流并审计。

### 11.2 策略生成和 Gateway

- schema 正反测试；
- security 继承/`security: []`；
- OR/AND 不丢失；
- operationId/method/path 唯一；
- FullPath 参数模板；
- migrated route 唯一 generated policy；
- legacy route 显式 inventory；
- 新 route 不进入 legacy；
- policy mode 默认 off；
- 注册 route 缺策略 503；
- unknown route 404；
- deny/error no-fallback；
- platform 空 tenant 不绕过横切中间件。

### 11.3 auth-service

- legacy platform TenantContext 保持零 UUID；
- V2 platform Principal tenant ID 为空；
- tenant Principal tenant ID 非零；
- 缺 principal_kind 的旧 JWT 仅作为 user；
- service JWT kind/audience；
- V2 Proto 无 roles；
- CheckPermissionV2 重验 credential，并拒绝 Principal/credential identity 不一致；
- user 权限从 DB 读取；
- API Key scope/revoked/expired 从 DB 读取；
- Gateway 自报 roles 对 V2 无影响；
- permission exact/wildcard/scope；
- DB error 不泄漏详情。

### 11.4 L1 最小安全基线

- auth-service ClusterIP-only；
- 非 Gateway Pod 无法连接或调用敏感 RPC；
- 错误 caller credential 被拒绝；
- 直接 CreateAPIKey 被拒绝；
- production-like reflection 关闭；
- 旧 Gateway rolling compatibility；
- NetworkPolicy rollback 不回到全开放。

### 11.5 后续零信任阶段预留验收（当前不执行）

本节不属于当前 Functional MVP 的验收范围。只有后续单独批准 `AUTHZ-ZEROTRUST-Z1/Z2` 后，才把以下项目转为实际门禁：

- plaintext 访问 secure endpoint 失败；
- 无证书/错误 CA/错误 SAN/过期证书失败；
- Gateway identity 只能调用 allowlisted RPC；
- inference-service 只能调用 IssueServiceToken；
- certificate rotation 无中断；
- legacy + secure 并存阶段旧/新 Gateway 均正常；
- 关闭明文后无 fallback；
- 跨 namespace/未授权 Pod live negative test。

### 11.6 quota-meta E2E

```text
platform-admin JWT（tid 缺失）
  -> ValidatePrincipal: platform/user/tenant_id=""
  -> generated quota/read/platform policy
  -> auth-service reloads permissions
  -> CheckPermissionV2 allow
  -> QuotaAdminService.ListQuotaMeta
  -> HTTP 200
```

### 11.7 可观测性

指标至少包括：

- operation_id；
- policy_source；
- principal_kind/domain；
- required_boundary；
- decision/reason；
- auth_rpc_version；
- caller_identity/trust_mode；
- transport=legacy|secure；
- legacy policy/RPC count。

禁止记录 JWT、API Key、caller credential 和 private key。

---

## 十二、验证命令

```bash
cd repo

go test -count=1 ./services/ani-gateway/internal/middleware
go test -count=1 ./services/ani-gateway/internal/authz
go test -count=1 ./services/ani-gateway/internal/router
go test -count=1 ./services/auth-service/internal/service
go test -count=1 ./pkg/bootstrap

python scripts/validate_auth_gateway_contract.py
python scripts/validate_core_api_compatibility.py
python scripts/validate_services_route_contract_test.py
python scripts/validate_services_route_contract.py --root .

make test
make validate-architecture
make validate-services
git diff --check
```

当前实现新增并接入 CI：authz schema、policy drift、route-policy coverage、legacy count non-increase 和 RPC compatibility。L1 network/caller negative test 仅在进入共享集群前增加；L2 mTLS live gate 留到后续零信任批次。

缺少 make、真实 PG 或真实集群时必须报告未运行范围，不得外推 runtime/production ready。

---

## 十三、发布与回滚

### 13.1 功能发布

1. Gateway B0，mode=off；
2. 验证旧功能；
3. auth-service B1 同时支持旧/V2；
4. 所有 auth-service 实例 V2 ready；
5. Gateway 带 pilot policy，mode 仍 off；
6. local/隔离 CI 开 pilot；
7. L1 完成后进入共享测试 canary；
8. 观察 allow/deny/error/legacy/audit；
9. 逐域迁移。

### 13.2 transport 发布

1. auth-service 先支持 legacy + secure；
2. 旧 Gateway 回归；
3. 新 Gateway 切 secure；
4. 观察明文/caller/RPC 指标；
5. 旧 Gateway 和明文调用归零；
6. 关闭明文；
7. 删除临时 caller credential。

### 13.3 回滚

- pilot 通过配置/部署整体切回已审核 legacy policy；
- 不做请求级 fallback；
- V2 不可用返回 503；
- L1 NetworkPolicy 回滚仍须保持仅 Gateway 可达；
- secure transport 回滚只能在受控迁移窗口切回受 NetworkPolicy 限制的 legacy endpoint；
- 不允许永久 silent insecure fallback；
- 回滚记录 operation、时间、原因、owner 和恢复条件。

---

## 十四、预计改动范围

### 14.1 功能主线

| 模块 | 主要改动 |
|---|---|
| Core/Services OpenAPI | x-ani-authz，保留冻结字段 |
| auth Proto | additive V2，不改旧 tag |
| auth-service JWT/API Key | Principal kind、credential ID |
| auth-service permission | DB 权威 evaluator |
| scripts | schema/generator/coverage/drift |
| Gateway internal/authz | Policy/Principal/legacy/constraints |
| Gateway middleware | generated/legacy 双路径和横切适配 |
| tests | 版本、主体、权限、路由、E2E |

### 14.2 L1/L2 安全演进

- auth-service NetworkPolicy；
- ClusterIP guard；
- production reflection switch；
- 临时 caller credential wiring；
- caller interceptor/逐 RPC ACL；
- cert-manager Issuer/Certificate；
- Gateway client/auth-service server TLS；
- transport migration/live evidence。

### 14.3 明确不做

- 不一次性迁移全部路由；
- 不在 pilot 删除 legacy；
- 不批量改 Core v1 action/scope；
- 不引入 OPA、新鉴权服务或运行时 YAML；
- 不把对象所有权/RLS 上移 Gateway；
- 不把 caller token 当最终零信任；
- 不在本设计提前冻结 mTLS 端口或 service mesh 实现；
- 不把 local/logic 或 L1 baseline 描述为 production-ready。

---

## 十五、关键代码索引

| 文件 | 当前职责 / 后续影响 |
|---|---|
| `api/openapi/v1.yaml` | Core policy 来源 |
| `api/openapi/services/v1.yaml` | Services policy 来源 |
| `api/core-v1-compatibility-baseline.yaml` | 旧 operation/scope 兼容 |
| `services/ani-gateway/internal/middleware/chain.go` | 全局中间件顺序 |
| `services/ani-gateway/internal/middleware/auth_client.go` | 当前 insecure auth client |
| `services/ani-gateway/internal/middleware/auth.go` | 当前认证/context/scope 混合 |
| `services/ani-gateway/internal/middleware/rbac.go` | inferPermission/旧 CheckPermission |
| `services/ani-gateway/internal/middleware/ratelimit.go` | 空 tenant 限流问题 |
| `services/ani-gateway/internal/router/quota_resources.go` | 首个试点 |
| `services/ani-gateway/internal/router/metering_resources.go` | platform metering route 缺口 |
| `services/auth-service/internal/service/auth_service.go` | 旧 auth RPC 和敏感 RPC |
| `services/auth-service/internal/service/api_keys.go` | API Key 权威数据和敏感创建路径 |
| `services/auth-service/internal/service/jwt.go` | JWT validator |
| `pkg/bootstrap/server.go` | 当前无 TLS/caller interceptor 的 gRPC server |
| `api/proto/auth/v1/auth_service.proto` | 旧 RPC + additive V2 |
| `deploy/migrations/20260502_003_permissions_schema.sql` | 权限 JSONB |
| `ANI-08-安全架构设计.md` | 零信任目标和状态来源 |

---

## 十六、完成定义

### 16.1 Functional MVP 完成

1. quota-meta 在 local/隔离 CI 通过；
2. platform 无 tenant ID；
3. V2 无 roles，服务端读取权威权限；
4. tenant/platform/API Key/service/sandbox 正反矩阵通过；
5. 未迁移路由行为不变；
6. deny/error no-fallback；
7. 只能声明 local/logic verified。

### 16.2 Gateway OpenAPI 鉴权功能迁移完成

1. Core、Services、proxy 注册产品路由全部 generated/public；
2. Principal kind/domain/boundary 全链路一致；
3. platform 和其他主体无回归；
4. permission evaluator 实际消费数据库权限；
5. route coverage、drift、compatibility 成为 CI 门禁；
6. legacy policy/业务 RPC 调用在观察窗口归零；
7. 回滚演练和 Feature batch 文档闭环；
8. 若仅在 L1 环境完成，仍不得声明 production-ready。

### 16.3 后续 Zero Trust 完成定义（当前不验收）

本节是未来演进完成条件，不属于当前方案实施完成条件：

1. Gateway/auth-service 使用 mTLS；
2. auth-service 从 peer identity 确认 caller；
3. 敏感 RPC caller allowlist 生效；
4. NetworkPolicy 持续生效；
5. 明文入口和临时 caller credential 删除；
6. 错误身份、跨 namespace、轮换 live gate 通过；
7. ANI-08 状态与代码、部署和 evidence 一致；
8. 完成独立 production readiness 评审。

在对应条件满足前，只能使用对应成熟度描述，不能把功能完成、最小隔离和零信任完成混为一体。
