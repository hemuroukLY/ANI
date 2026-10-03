# ANI Gateway OpenAPI 鉴权策略改造方案 V2

> 文档状态：修订设计方案，尚未实施
>
> 版本：V2
>
> 修订日期：2026-08-20
>
> 基于：`gateway-openapi-authorization-policy-migration.md`
>
> 首个纯鉴权试点：`GET /api/v1/admin/quota-meta`
>
> 第二试点：`GET /api/v1/metering/usage/platform`，必须等待平台计量查询运行时闭环完成
>
> 目标：以 OpenAPI 为产品路由鉴权策略来源，消除 Gateway 根据 URL 和 HTTP 方法猜测权限的行为，同时保证现有 tenant、platform、API Key、service、sandbox 和 dev/local 功能无回归。
>
> 本文是设计方案，不表示代码、测试、数据库权限或生产验证已经完成。

---

## 一、结论

本方案在完成本文定义的兼容阶段和验收矩阵后可实施，并满足以下结果：

1. 现有功能在 operation 显式迁移前继续走原有鉴权链路，运行行为不变。
2. 平台主体的规范语义为 `credential_domain=platform` 且 `tenant_id` 为空；不再把零 UUID 当作真实租户。
3. 迁移期保留旧 `ValidateToken` / `CheckPermission` RPC 的原有字段和值语义，新增 V2 RPC 承载新 Principal 和 boundary，避免滚动升级破坏旧 Gateway。
4. OpenAPI operation 使用结构化 `x-ani-authz` 声明资源、动作、数据边界和允许的主体类型。
5. 构建期生成 Gateway 静态策略，运行时不解析 OpenAPI YAML。
6. 已迁移 operation 的 deny/error 不自动回退 legacy；未迁移 operation 只通过显式 legacy 清单继续运行。
7. sandbox token 的实例绑定、路径限制和本地验签继续强制执行，不能被普通 tenant RBAC 替代。
8. Core、Services 和 OpenAI-compatible proxy 分域迁移；全部产品路由完成迁移前不删除全局 legacy 鉴权。

本方案不会把以下概念混为一体：

- `security`：OpenAPI 接受的凭证方案；
- `principal_kind`：用户、服务身份、API Key 或 sandbox token；
- `credential_domain`：凭证代表的主体域，取值 `tenant|platform|sandbox`；
- `boundary`：operation 访问的数据边界，取值 `own|tenant|platform`；
- `resource + action`：具体业务权限；
- handler/store/RLS：对象所有权、租户过滤和数据面隔离的最终执行层。

`x-ani-authz` 是 ANI 私有 OpenAPI 扩展，不是行业标准字段。通用部分是认证与授权分离、显式策略、构建期校验、默认拒绝和可审计迁移。

---

## 二、兼容性与现有功能影响

### 2.1 最终影响结论

兼容机制按本文顺序实施后：

- tenant 用户现有 API 行为不变；
- tenant API Key 行为不变，仍不能访问 platform boundary；
- platform 管理员可以在没有 tenant ID 的情况下合法访问 platform operation；
- platform 管理员不会自动获得 tenant/own operation 的访问资格；
- sandbox token 继续只能访问 token 绑定实例的允许子资源；
- 公开登录、健康检查和 branding 行为不变；
- `ANI_AUTH_MODE=dev` 的 local/dev 行为保留，但不能在 production profile 启用；
- 未迁移 operation 的 resource/action 推断结果在迁移期保持原样。

这些结论只有在第十一节兼容矩阵全部通过后才成立，不能仅凭 Proto 可编译或单元测试通过宣称“无影响”。

### 2.2 为什么平台管理员可以没有 tenant ID

平台主体访问的是平台或跨租户数据边界，不属于任何单一租户。规范 Principal 示例：

```json
{
  "principal_kind": "user",
  "credential_scheme": "bearer",
  "credential_domain": "platform",
  "tenant_id": "",
  "subject_id": "<platform-user-id>",
  "roles": ["platform-admin"]
}
```

固定规则：

- `credential_domain=platform` 时 `tenant_id` 必须为空；
- `credential_domain=tenant` 时 `tenant_id` 必须是非零合法 UUID；
- `credential_domain=sandbox` 时必须携带合法 tenant ID，并携带已验签的 sandbox instance claim；
- platform handler 使用明确的平台查询/事务边界，不注入伪造 `types.TenantContext`；
- 平台访问指定租户资源时，目标 `tenant_id` 来自 path/query，并由 platform permission 和平台数据访问边界校验，不等于 Principal 自身 tenant ID。

### 2.3 为什么不能直接修改旧 ValidateToken

当前旧链路依赖 platform token 返回零 UUID：

1. Gateway 旧 `withTenantContextStrict` 要求 tenant ID 可解析为 UUID；
2. 旧 `CheckPermission` 无条件要求 `tenant_id` 非空；
3. platform-admin 依赖角色直通继续执行；
4. RateLimit、Audit、Idempotency 和部分 handler 仍读取旧 request-context key。

如果 auth-service 先把旧 `ValidateToken.TenantContext.tenant_id` 改为空，旧 Gateway 和未迁移 platform operation 会在滚动升级期间返回 401/403。

因此 V2 采用以下兼容策略：

- 旧 `ValidateToken`：迁移期保持原返回语义，包括 platform 零 UUID；
- 旧 `CheckPermission`：迁移期保持原请求和判定语义；
- 新 `ValidatePrincipal`：返回规范 Principal，platform tenant ID 为空；
- 新 `CheckPermissionV2`：显式接收 credential domain、required boundary、principal kind、resource 和 action；
- 只有已迁移 operation 使用 V2 RPC；
- 全量迁移、观察窗口和旧 Gateway 清零后，才能单独评审旧 RPC 的废弃计划。

这是新增 V2 RPC 的必要理由：解决当前零停机迁移问题，不是为未来假设提前引入抽象。

### 2.4 禁止的兼容方式

以下方式禁止使用：

- 新鉴权返回 deny/error 后自动再试旧鉴权；
- 直接把空 platform tenant ID 注入 `types.TenantContext`；
- 把零 UUID继续当作规范平台身份保存到业务上下文；
- 根据角色名猜测 credential domain；
- 根据路径前缀猜测新 operation 的 boundary；
- 在 auth-service 未全量支持 V2 RPC 前打开 pilot 开关；
- 为绕过迁移问题让 platform-admin 自动进入 tenant/own boundary。

---

## 三、当前代码事实与问题

### 3.1 当前执行链路

当前全局顺序：

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

当前存在三套互相可能不一致的描述：

| 来源 | 当前用途 | 问题 |
|---|---|---|
| OpenAPI `security` / `x-ani-rbac-scope` | 契约、SDK、文档和兼容基线 | Gateway 未消费精确权限语义 |
| `scopeAllowedForPath` | tenant/platform/sandbox 路由隔离 | 依赖 URL 前缀和手工列表 |
| `inferPermission` | 推断 resource/action | 依赖 HTTP 方法和路径首段 |

### 3.2 平台 token 现状

当前 JWT 内部语义已经允许 `scope=platform` 且无真实 tenant ID，但 `ValidateToken` 使用 `claims.TenantID.String()`，把 `uuid.Nil` 输出为零 UUID字符串。旧链路因而能继续工作，但身份语义不准确。

V2 必须在新 Principal 中规范化为空值，同时保留旧 RPC 的迁移兼容输出。

### 3.3 principal kind 现状

当前 JWT payload 没有可信的 `principal_kind` 字段，也没有在 ANI access token validator 中完成 service-token audience 约束。不能仅通过角色名把用户 JWT 判断为 service JWT。

当前可信来源只能是：

- API Key：通过 auth-service API Key 校验分支确定；
- sandbox：通过 Gateway 本地 sandbox token 格式和签名确定；
- 现有登录/刷新 JWT：迁移期安全默认 `user`；
- service JWT：必须有 ANI 签名的显式 `principal_kind=service` 和预期 audience，未满足时不得视为 service。

### 3.4 权限数据现状

数据库 migration 已定义 `{resource, actions, scope}` 权限结构，但当前 `CheckPermission` 没有读取该结构，只执行角色直通和字符串 scope 匹配。

因此以下能力不能在没有新 evaluator 的情况下宣称完成：

- platform 非管理员精确权限；
- permission scope 与 operation boundary 的严格匹配；
- platform-admin/tenant-admin 通配权限仅在对应 boundary 生效；
- auditor 跨 boundary 的精确控制。

### 3.5 platform metering 现状

`GET /metering/usage/platform` 已存在于 Core OpenAPI，但 Gateway 当前未注册该 route，且只装配 `LocalMeteringService`。因此它不能作为第一批“仅鉴权改造”的 200 E2E 试点。

V2 将首个试点改为当前已注册的 `GET /admin/quota-meta`。platform metering 保留为第二试点，但必须先完成 route、平台查询 port/adapter、跨租户数据边界和真实非超级用户验证。

### 3.6 全局 Gateway 范围

同一中间件链覆盖：

- Core `/api/v1/*`；
- Services `/api/v1/svc/*`；
- OpenAI-compatible `/v1/chat/completions` 和 `/v1/inference/stream`；
- health/no-route 等基础设施路径。

因此只迁移 Core OpenAPI 后不能删除全局 legacy 逻辑。Services 和 proxy 必须各自有明确策略来源或继续留在显式 legacy 清单。

---

## 四、目标鉴权模型

### 4.1 OpenAPI 声明

首个试点 `GET /admin/quota-meta`：

```yaml
/admin/quota-meta:
  get:
    operationId: listQuotaMeta
    security:
      - BearerAuth: []

    # Core v1 兼容字段，保持原值。
    x-ani-rbac-scope: "scope:quota:read"

    x-ani-authz:
      version: v1
      resource: quota
      action: read
      boundary: platform
      principal_kinds: [user]
```

第二试点在数据面前置完成后使用：

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

字段定义：

| 字段 | 必填 | 说明 |
|---|---:|---|
| `version` | 是 | 当前固定 `v1`，未知版本生成失败 |
| `resource` | 是 | 稳定业务资源名，不从 URL 推断 |
| `action` | 是 | 精确业务动作，不从 HTTP 方法推断 |
| `boundary` | 是 | 数据边界：`own|tenant|platform` |
| `principal_kinds` | 是 | `user|service|api_key|sandbox` 的非空集合 |

`security` 决定凭证方案，`principal_kinds` 限制凭证代表的主体类型。两者不能互相替代。

### 4.2 凭证域与数据边界分离

| 概念 | 取值 | 用途 |
|---|---|---|
| credential scheme | bearer、api_key、sandbox_bearer | 凭证传输和校验方式 |
| principal kind | user、service、api_key、sandbox | 主体类型 |
| credential domain | tenant、platform、sandbox | 已认证主体所在域 |
| required boundary | own、tenant、platform | operation 要访问的数据范围 |

固定兼容规则：

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

这只是主体域前置检查。sandbox 仍需额外的实例 claim 限制；`own` 仍需 handler/store 校验对象所有权；tenant 仍需 tenant filter/RLS；platform 仍需平台事务和数据库角色边界。

### 4.3 权限 scope 匹配

数据库 permission 与 required boundary 的匹配规则固定如下：

| required boundary | 允许的 permission scope |
|---|---|
| `platform` | 仅 `platform` |
| `tenant` | 仅 `tenant` |
| `own` | `own` 或更宽的同租户 `tenant` |

附加规则：

- `platform-admin` 的 `*/*/platform` 通配只在 platform boundary 生效；
- `tenant-admin` 的 `*/*/tenant` 通配只在 tenant/own boundary 生效；
- platform-admin 不自动作为 tenant Principal；
- tenant-admin、tenant API Key 和 sandbox 不得进入 platform boundary；
- `operation_id` 只用于审计与诊断，不作为权限 key；
- resource/action `*` 只按现有内置角色和受控权限数据允许，不由 Gateway自行扩展。

### 4.4 sandbox 凭证约束

sandbox token 是受签名和实例范围约束的 capability credential，不进入 auth-service 的通用角色授权。每个 sandbox 请求必须执行：

1. 本地验证签名、过期时间和 token 类型；
2. 验证 operation policy 允许 `principal_kind=sandbox`；
3. 验证 required boundary 为 `own` 或明确允许的 tenant 子资源；
4. 验证 path 中 `instance_id` 与 token claim 完全相等；
5. 验证 method/path 属于 token 允许的 sandbox 子资源集合；
6. 注入真实 tenant data context；
7. handler/store 继续验证实例属于该租户并按跨租户 404 处理。

以上检查全部通过即构成 sandbox 的本地授权决定，不再调用 `CheckPermissionV2`。实现形态为 Gateway 内部 `CredentialConstraintEvaluator`；V1 只提供 sandbox evaluator，不引入通用插件框架。generated policy 仍必须显式允许 sandbox，并声明其 resource/action/boundary，保证路由能力边界可审计、可覆盖检查。

---

## 五、生成策略与路由解析

### 5.1 构建期生成

```text
api/openapi/v1.yaml
  -> scripts/generate_gateway_authz.py
  -> services/ani-gateway/internal/authz/zz_generated_core_policies.go
```

Services 后续独立生成：

```text
api/openapi/services/v1.yaml
  -> scripts/generate_gateway_authz.py
  -> services/ani-gateway/internal/authz/zz_generated_services_policies.go
```

生成物必须表示 OpenAPI Security Requirement 的 OR/AND 结构，不能用一个扁平 schemes 数组丢失语义：

```go
type SecurityRequirement struct {
    AllOf []SecurityScheme
}

type Policy struct {
    OperationID         string
    Method              string
    PathTemplate        string
    Public              bool
    SecurityAlternatives []SecurityRequirement // alternatives 之间 OR；AllOf 内 AND
    Resource            string
    Action              string
    Boundary            Boundary
    PrincipalKinds      []PrincipalKind
}
```

当前 ANI 主要使用单 scheme OR；生成器仍必须正确保留 OpenAPI 结构，或对不支持的多 scheme AND 显式生成失败，不能静默放宽。

### 5.2 路由解析

运行时使用 Hertz 已匹配的 `RequestContext.FullPath()`，不重新实现一套 URL matcher：

1. Hertz 完成 route match；
2. `ResolveAuthzPolicy` 读取 method + `FullPath()`；
3. 将 Hertz `:param` 规范化为 OpenAPI `{param}`；
4. 查找生成策略或显式 legacy entry；
5. 已注册产品路由两者都不存在时 fail closed。

固定行为：

- `FullPath()==""` 的未知路由交给 NoRoute 返回 404；
- 已注册路由缺策略且不在 legacy 清单，返回 503 `AUTHORIZATION_POLICY_UNAVAILABLE`；
- method/path/operationId 重复或冲突在构建期失败；
- 不用原始 path 自行匹配参数模板，避免与 Hertz 路由优先级不一致。

### 5.3 legacy 清单

legacy entry 必须包含：

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

- 当前未迁移路由可一次性生成初始清单；
- 新增产品 route 不得加入 legacy，必须直接声明 `x-ani-authz`；
- legacy 数量只能减少；
- 已迁移 operation 的 deny/error 绝不按请求级自动回退；
- 紧急回滚只能通过受控配置/部署把整个 operation 切回已知 legacy policy，并留下审计记录。

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

迁移期每个中间件根据 operation 的 policy source 进入 generated 或 legacy 分支，但只认证一次，不能先失败再尝试另一套授权。

### 6.2 Principal

```go
type Principal struct {
    Kind             PrincipalKind
    CredentialScheme SecurityScheme
    CredentialDomain CredentialDomain
    TenantID         string
    SubjectID        string
    CredentialID     string
    Roles            []string
    SandboxClaims    *SandboxClaims
}
```

字段规则：

- user JWT：`SubjectID=user_id`；
- service JWT：`SubjectID=service subject`；
- API Key：`CredentialID=key_id`，SubjectID 可为绑定 user ID；
- sandbox：SubjectID 为固定 sandbox actor，SandboxClaims 必填；
- platform：TenantID 必须为空；
- tenant/sandbox：TenantID 必须是合法非零 UUID。

### 6.3 认证

认证固定流程：

1. public policy 直接进入后续公开路由处理；
2. 同时出现 Bearer 和 API Key 时拒绝，避免凭证优先级歧义；
3. sandbox bearer 由 Gateway 本地验签；
4. 普通 Bearer/API Key 调用 `ValidatePrincipal`；
5. Principal 通过严格结构校验后写入 request context；
6. tenant/sandbox 注入 tenant data context；platform 不注入 tenant data context；
7. legacy operation 额外生成兼容 request-context view，保持旧 handler/middleware 读取行为。

`ANI_AUTH_MODE=dev` 继续生成 tenant/user Principal，并仅在允许的 dev/local profile 生效；production profile 检测到 dev mode 必须拒绝启动或 readiness 失败。

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

        // sandbox token 本身是实例范围 capability；本地约束通过即授权，
        // 不把 Gateway 私有 claim 伪装成 auth-service PrincipalContext。
        if principal.Kind == PrincipalSandbox {
            c.Next(ctx)
            return
        }

        result, err := authClient.CheckPermissionV2(ctx, AuthorizationRequest{
            Principal:         principal,
            RequiredResource:  policy.Resource,
            RequiredAction:    policy.Action,
            RequiredBoundary:  policy.Boundary,
            OperationID:       policy.OperationID,
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

错误语义：

- 无凭证或凭证无效：401；
- 主体类型、domain、boundary、constraint 或 permission 不匹配：403；
- auth-service 超时/不可用：503，仍然 fail closed；
- 注册路由缺策略：503；
- 日志和响应不得暴露 Token、API Key、完整 claims、SQL 或内部 gRPC 错误。

### 6.5 横切中间件适配

#### Rate limit

不能再以空 tenant ID 表示“无需限流”。限流主体键：

| Principal | 限流键主体部分 |
|---|---|
| tenant user | `tenant:<tenant_id>:user:<subject_id>` |
| platform user | `platform:user:<subject_id>` |
| API Key | `tenant:<tenant_id>:api_key:<credential_id>` |
| service | `<domain>:service:<subject_id>` |
| sandbox | `tenant:<tenant_id>:sandbox:<instance_id>` |

#### Idempotency

幂等缓存键使用 credential domain + tenant/subject identity + operation，而不是只依赖 tenant ID。已有 API 的 key/body/fingerprint 语义保持不变。

#### Audit

至少记录：

- request_id；
- operation_id；
- policy_source；
- principal_kind；
- credential_domain；
- tenant_id（允许空）；
- subject_id/credential_id 的安全标识；
- required boundary；
- decision 和规范化 reason；
- HTTP status 和 duration。

拒绝和 auth-service error 也必须被审计；不得只在成功进入 handler 后记录。

---

## 七、内部鉴权契约与可信身份

### 7.1 保留旧 RPC

迁移期不改变：

```proto
rpc ValidateToken(ValidateTokenRequest) returns (common.v1.TenantContext);
rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
```

旧生成物和调用方继续兼容。任何废弃必须在全量迁移完成后独立评审。

### 7.2 新增 V2 RPC

示意契约：

```proto
rpc ValidatePrincipal(ValidatePrincipalRequest) returns (PrincipalContext);
rpc CheckPermissionV2(AuthorizationRequest) returns (AuthorizationDecision);

message PrincipalContext {
  string principal_kind = 1;    // user|service|api_key
  string credential_scheme = 2; // bearer|api_key
  string credential_domain = 3; // tenant|platform
  string tenant_id = 4;         // platform 必须为空
  string subject_id = 5;
  string credential_id = 6;
  repeated string roles = 7;
}

message AuthorizationRequest {
  PrincipalContext principal = 1;
  string resource = 2;
  string action = 3;
  string required_boundary = 4; // own|tenant|platform
  string operation_id = 5;      // audit/debug only
}

message AuthorizationDecision {
  bool allowed = 1;
  string reason_code = 2;
}
```

sandbox Principal 不通过该 RPC 返回，也不调用 `CheckPermissionV2`。Gateway 根据 generated policy、本地验签、实例/path/method capability 约束作出授权决定，再由 tenant data context、handler/store 和 RLS 完成数据面隔离。

### 7.3 JWT principal kind

新增 ANI access token claim：

```json
{
  "principal_kind": "user|service",
  "aud": "ani-core"
}
```

兼容规则：

- 现有 ANI 登录/刷新 token 缺少 principal_kind 时按 `user` 处理；
- 缺少 principal_kind 的 token 永远不能访问 service-only operation；
- service 必须显式 `principal_kind=service` 且 audience 包含 `ani-core`；
- principal kind 只能来自签名 claim 或已验证凭证类型，不能来自请求 header 或角色名；
- API Key 永远由数据库校验结果确定为 `api_key`；
- platform/tenant domain 由受信 claim/凭证记录确定，并与 tenant ID 结构交叉验证。

### 7.4 权限 evaluator

`CheckPermissionV2` 必须消费现有角色权限数据结构，而不是只做角色名直通：

1. 解析/获取 Principal 的有效角色权限；
2. 匹配 resource exact 或受控 wildcard；
3. 匹配 action exact 或受控 wildcard；
4. 按 4.3 规则匹配 permission scope 与 required boundary；
5. platform-admin/tenant-admin 只在各自 scope 生效；
6. 返回稳定 reason code，不返回数据库或内部错误文本。

实现复用现有 roles/user_roles 和 permissions JSONB，不新增表、不引入 OPA、不新增独立鉴权服务。可以增加 auth-service 内部窄接口和 Postgres repository 以便测试，但不得绕过现有数据库边界。

生产当前没有精确 platform readonly 角色时，试点的真实正例只声明 platform-admin；精确权限通过隔离 fixture/integration role 验证，不得把测试 fixture 描述成已交付产品角色。

---

## 八、试点与迁移范围

### 8.1 首个试点：GET /admin/quota-meta

选择理由：

- OpenAPI 已声明 operationId 和 `x-ani-rbac-scope`；
- Gateway 已注册 route；
- handler 已通过 `QuotaAdminService` port 调用；
- 路径已由旧 `scopeAllowedForPath` 识别为 platform；
- GET 只读，无幂等写入语义；
- 不依赖尚未完成的 platform metering 查询 adapter。

试点只切换该 operation 的认证/授权策略，不改变 QuotaAdminService、handler response 或数据库查询语义。

预期矩阵：

| 调用方 | 预期 |
|---|---:|
| 无凭证 | 401 |
| tenant user | 403 |
| tenant API Key | 403 |
| platform service JWT | 403 |
| platform user 无 quota/read/platform 权限 | 403 |
| platform-admin，无 tenant ID | 200 |
| auth-service 不可用 | 503 |

### 8.2 第二试点：GET /metering/usage/platform

进入条件：

- Gateway route 已注册；
- platform query 的 port 语义已明确，不用空 TenantID 暗示跨租户；
- PG query adapter 已装配；
- tenant_id 可选筛选执行二次 platform 权限/数据边界校验；
- group_by、时间边界和 UTC 语义与 OpenAPI 一致；
- 非超级用户平台数据库角色、跨租户读取和 RLS 行为已验证；
- handler/adapter E2E 在旧 platform 鉴权下先通过。

满足后才能把该 operation 切换到 generated policy。鉴权改造不能被用来掩盖 route/adapter 不存在的问题。

### 8.3 全量迁移域

| 域 | 策略来源 | 删除 legacy 的条件 |
|---|---|---|
| Core `/api/v1/*` | `api/openapi/v1.yaml` | Core 注册产品路由全部覆盖 |
| Services `/api/v1/svc/*` | `api/openapi/services/v1.yaml` | Services security 默认和 route baseline 收敛 |
| OpenAI-compatible `/v1/*` | 独立正式契约或明确归属的 Services 契约 | 两条产品路由都有显式策略 |
| health/readiness/no-route | 最小基础设施策略 | 与产品 API 策略分离 |

Services OpenAPI 在生成策略前必须先补明确的全局或 operation security；缺少 security 不能被自动解释为“公开”。

---

## 九、零停机实施顺序

### AUTHZ-POLICY-A：策略格式和生成门禁

交付：

- `x-ani-authz` schema/validator；
- Core policy generator；
- OR/AND security 结构；
- 生成物漂移门禁；
- Hertz/OpenAPI path template 规范化测试；
- 显式 legacy inventory。

运行模式保持 `off`，不改变请求行为。

### AUTHZ-COMPAT-B0：Gateway Principal 兼容层

交付：

- Gateway 内部 Principal 类型；
- 从旧 TenantContext 规范化 Principal；
- platform 零 UUID 在内部规范化为空，但 legacy 调用保留旧 wire view；
- RateLimit、Idempotency、Audit 支持 Principal；
- tenant/platform/sandbox/dev 旧行为回归；
- `GATEWAY_AUTHZ_POLICY_MODE=off|pilot|full`，默认 `off`。

仍只调用旧 RPC，不切换 operation。

### AUTHZ-CONTRACT-B1：新增 V2 RPC 和可信 principal kind

交付：

- 新 `ValidatePrincipal` / `CheckPermissionV2` Proto；
- 生成物；
- JWT principal_kind/audience 签发与校验；
- 缺 claim 的现有用户 token 兼容；
- API Key Principal；
- permission evaluator；
- auth-service 同时支持旧/V2 RPC。

部署 auth-service，但 Gateway mode 仍为 `off`。必须等待所有 auth-service 实例都支持 V2，再进入 pilot。

### AUTHZ-PILOT-C：切换 quota-meta

交付：

- `listQuotaMeta` 的 `x-ani-authz`；
- generated policy；
- Gateway pilot allowlist 只包含该 operation；
- 正反 Principal/permission 矩阵；
- 无 tenant ID platform-admin E2E；
- legacy route 回归；
- canary metrics 和受控回滚步骤。

pilot operation 的 deny/error 不自动回退。出现问题时通过配置/部署整体关闭 pilot，而不是单请求 fallback。

### AUTHZ-METERING-D：计量数据面前置与第二试点

先完成第 8.2 节进入条件，再迁移 `getPlatformMeteringUsage`。数据面批次与鉴权切换批次可独立评审，不能在一个 PR 中混合未验证的 PG 查询和全局鉴权删除。

### AUTHZ-MIGRATION-E：Core 全量迁移

- 逐资源域补齐 `x-ani-authz`；
- 已有 action 保持原语义，不顺便批量改名；
- legacy count 单调下降；
- Core operation/route/policy 一致性成为必过 CI；
- 完成观察窗口后删除 Core route 的路径/方法推断。

### AUTHZ-SERVICES-F：Services 和 proxy 迁移

- Services OpenAPI security 语义先收口；
- 生成 Services policies；
- 处理 accepted route baseline；
- 为 `/v1/*` proxy 建立正式策略来源；
- 通过 `make validate-services` 和跨层回归。

### AUTHZ-LEGACY-G：删除全局旧逻辑

仅在以下条件全部满足后执行：

- Core、Services、proxy 的已注册产品路由全部有 generated policy；
- legacy 指标在规定观察窗口持续为零；
- tenant/platform/API Key/service/sandbox/dev/public 回归全部通过；
- 旧 Gateway 实例为零；
- 旧 RPC 调用指标为零；
- 生成物漂移和 route coverage 是必过 CI；
- 已准备并演练受控回滚。

随后才可删除：

- `scopeAllowedForPath`；
- `inferPermission`；
- 产品 API 的 `isPublicPath` 硬编码；
- 旧 principal-kind/service-only 运行时分支；
- 旧 RPC，需独立兼容性评审后处理。

Core v1 已进入兼容基线的扩展字段继续保留为契约元数据，未来 v2 再清理。

---

## 十、预计改动范围

### 10.1 策略和兼容基础

| 模块 | 主要改动 |
|---|---|
| `api/openapi/v1.yaml` | 试点 operation 新增 `security` / `x-ani-authz`，保留旧 scope |
| `api/proto/auth/v1/auth_service.proto` | 新增 V2 RPC/messages，不改旧 tag/语义 |
| Proto 生成物 | 按现有流程更新 |
| auth-service JWT | principal_kind/audience 与旧 token 兼容 |
| auth-service permission | 消费既有 permissions JSONB |
| `scripts/` | schema、generator、route coverage、drift gate |
| Gateway `internal/authz/` | Policy、registry、Principal、legacy inventory、constraint evaluator |
| Gateway middleware | 双路径认证授权和横切中间件适配 |
| tests | 组件版本矩阵、主体矩阵、route/policy、E2E |

### 10.2 metering 第二试点额外范围

由独立数据面方案确认，至少包含：

- platform route/handler；
- 明确的 platform query port 语义；
- PG adapter 和生产装配；
- RLS/角色权限；
- tenant filter、时间和 group_by 语义；
- 真实非超级用户 E2E。

### 10.3 明确不做

- 不一次性迁移全部 Gateway route；
- 不在试点中删除 legacy；
- 不改已冻结 `x-ani-rbac-scope` 值；
- 不批量重命名既有 Core v1 action；
- 不引入 OPA、新鉴权服务或运行时 YAML 解析；
- 不把对象所有权/RLS 上移到 Gateway；
- 不让 platform-admin 自动成为 tenant Principal；
- 不把测试 fixture 的精确角色描述为已交付产品角色。

---

## 十一、测试与验收标准

### 11.1 组件滚动升级矩阵

| Gateway | auth-service | 模式 | 预期 |
|---|---|---|---|
| old | old | legacy | 现有行为 |
| new B0 | old | off | 现有行为 |
| old | new B1 | legacy | 现有行为，旧 RPC 未变 |
| new B0 | new B1 | off | 现有行为 |
| new pilot | new B1 全量实例 | pilot | 仅试点 operation 使用 V2 |
| new full | new B1+ | full | 仅在全量门禁后允许 |

不允许 `new pilot` 在仍有不支持 V2 的 auth-service 实例时启用。

### 11.2 现有功能回归矩阵

至少覆盖：

- tenant 账密登录、refresh、普通 Core GET 和副作用写；
- platform 账密登录、refresh；
- platform-admin 无真实 tenant ID 调用既有 quota/tenant/plan 管理 API；
- tenant 用户和 tenant-admin 无法访问 platform operation；
- platform-admin 无法直接进入 tenant/own operation；
- API Key 正常访问允许的 tenant operation，拒绝 platform operation；
- service-only operation 拒绝普通用户 JWT；
- sandbox token 正常访问绑定实例子资源；
- sandbox token 跨实例、跨路径、过期和篡改均拒绝；
- public health/login/branding；
- dev/local profile；
- auth-service 不可用时 503 且无 legacy bypass；
- platform 请求仍受限流并被审计。

### 11.3 策略生成

- 缺字段、未知 version/boundary/principal kind 生成失败；
- security 继承和 `security: []` 正确；
- OR/AND 结构不丢失；
- operationId/method/path 唯一；
- Hertz `:param` 与 OpenAPI `{param}` 一致；
- 已迁移路由有且仅有一条 generated policy；
- 未迁移路由有显式 legacy entry；
- 新 route 不得加入 legacy；
- 生成物无漂移。

### 11.4 auth-service

- platform V2 Principal tenant ID 为空；
- legacy platform TenantContext 迁移期保持旧语义；
- tenant Principal tenant ID 必须合法非零；
- 缺 principal_kind 的现有 JWT 只能作为 user；
- service JWT 必须有签名 kind 和 audience；
- platform-admin wildcard 仅 platform allow；
- tenant-admin wildcard 仅 tenant/own allow；
- 精确 permission fixture 按 resource/action/scope 生效；
- DB/内部错误返回稳定 error，不泄漏详情。

### 11.5 Gateway

- `FullPath()` 未匹配返回正常 404；
- 注册 route 缺 policy/legacy 时 503 fail closed；
- migrated deny/error 不触发 legacy；
- 同时提交 Bearer/API Key 被拒绝；
- Principal 结构校验严格；
- sandbox credential constraint 始终执行；
- platform 空 tenant 不绕过 rate limit；
- denied/error 请求进入 audit/metrics；
- policy mode 默认为 off。

### 11.6 试点 E2E

quota-meta pilot 必须证明：

```text
platform-admin JWT（tid 缺失）
  -> ValidatePrincipal: platform/user/tenant_id=""
  -> generated quota/read/platform policy
  -> CheckPermissionV2 allow
  -> QuotaAdminService.ListQuotaMeta
  -> HTTP 200
```

同时证明 tenant user、tenant API Key、platform service 和无权限 platform user 返回预期 401/403，auth-service 故障返回 503。

### 11.7 可观测性和退出条件

低基数字段/指标：

- `operation_id`；
- `policy_source=generated|legacy`；
- `principal_kind`；
- `credential_domain`；
- `required_boundary`；
- `decision=allow|deny|error`；
- `reason_code`；
- `auth_rpc_version=legacy|v2`。

退出 legacy 的证据是 route coverage 完整、旧 RPC/legacy 指标归零并经过观察窗口，不是只看单元测试。

---

## 十二、验证命令

每批至少运行与变更相关的 focused tests，并最终运行：

```bash
cd repo

go test -count=1 ./services/ani-gateway/internal/middleware
go test -count=1 ./services/ani-gateway/internal/authz
go test -count=1 ./services/ani-gateway/internal/router
go test -count=1 ./services/auth-service/internal/service

python scripts/validate_auth_gateway_contract.py
python scripts/validate_core_api_compatibility.py
python scripts/validate_services_route_contract_test.py
python scripts/validate_services_route_contract.py --root .

make test
make validate-architecture
make validate-services        # 触及 Services 策略或 mixed Gateway 时
make validate-doc-entrypoints # 更新入口文档时
git diff --check
```

新增并接入 Makefile/CI 的目标名称在实现批次冻结，但必须覆盖：

- authz schema validation；
- policy generation drift；
- route-policy coverage；
- legacy count non-increase；
- component compatibility matrix。

若本地缺少 `make`、真实 PG 或生产 auth-service 环境，必须明确报告未运行范围，不得把 focused unit tests 外推为 runtime ready。

---

## 十三、受控发布与回滚

### 13.1 发布顺序

1. 部署 Gateway B0，mode=off；
2. 验证现有功能回归；
3. 部署同时支持 legacy/V2 的 auth-service B1；
4. 确认所有 auth-service 实例暴露 V2 readiness/capability；
5. 部署含 pilot policy 的 Gateway，mode 仍为 off；
6. 小流量/canary 开启 pilot；
7. 观察 allow/deny/error、legacy、rate-limit 和 audit 指标；
8. 扩大 pilot；
9. 满足退出条件后逐域迁移。

### 13.2 回滚原则

- 自动 fallback 禁止；
- auth-service V2 故障时请求返回 503；
- pilot 回滚通过配置/部署整体切回已审核 legacy policy；
- 回滚动作记录 operation、时间、原因、owner 和恢复条件；
- 不通过重新引入零 UUID业务上下文或放宽 boundary 完成回滚。

---

## 十四、关键代码索引

| 文件 | 当前职责 / 后续影响 |
|---|---|
| `api/openapi/v1.yaml` | Core REST 契约与 Core authz policy 来源 |
| `api/openapi/services/v1.yaml` | Services REST 契约与后续 Services policy 来源 |
| `api/core-v1-compatibility-baseline.yaml` | 冻结旧 operationId/scope 等兼容字段 |
| `services/ani-gateway/internal/middleware/chain.go` | 当前全局中间件顺序 |
| `services/ani-gateway/internal/middleware/auth.go` | 当前认证、公开路径、scope/path 和 tenant context 混合逻辑 |
| `services/ani-gateway/internal/middleware/rbac.go` | 当前 inferPermission 和旧 CheckPermission |
| `services/ani-gateway/internal/middleware/ratelimit.go` | 当前空 tenant ID 跳过限流 |
| `services/ani-gateway/internal/middleware/idempotency.go` | 当前 scope/tenant 幂等键 |
| `services/ani-gateway/internal/middleware/audit.go` | 当前 tenant-only audit entry |
| `services/ani-gateway/internal/router/quota_resources.go` | 首个 pilot route/handler |
| `services/ani-gateway/internal/router/metering_resources.go` | 当前未注册 platform metering route |
| `services/ani-gateway/internal/router/router.go` | Core/Services/proxy 共同路由装配 |
| `services/auth-service/internal/service/auth_service.go` | 旧 ValidateToken/CheckPermission |
| `services/auth-service/internal/service/jwt.go` | 当前 JWT claims/validator |
| `services/auth-service/internal/service/token_issuer.go` | 当前 user/platform token 签发 |
| `api/proto/common/v1/common.proto` | 旧 TenantContext，迁移期保持兼容 |
| `api/proto/auth/v1/auth_service.proto` | 旧 RPC + 新 V2 RPC |
| `deploy/migrations/20260502_003_permissions_schema.sql` | 既有 `{resource, actions, scope}` 权限模型 |
| `scripts/validate_core_api_compatibility.py` | Core v1 兼容性门禁 |
| `scripts/validate_services_route_contract.py` | Services OpenAPI/Gateway route baseline |

---

## 十五、最终完成定义

只有同时满足以下条件，才可标记 Gateway OpenAPI authorization migration 完成：

1. Core、Services 和 proxy 的所有已注册产品路由都有明确 generated policy；
2. Principal kind/domain/boundary 在 Gateway、auth-service、Proto、JWT 和审计中语义一致；
3. platform 主体以空 tenant ID 工作，且现有平台管理功能无回归；
4. tenant、API Key、service 和 sandbox 的正反矩阵通过；
5. sandbox 实例绑定和数据面 tenant/RLS 边界未削弱；
6. permission evaluator 实际消费 `{resource, actions, scope}`；
7. legacy 和旧 RPC 调用指标在观察窗口归零；
8. route coverage、生成漂移和兼容矩阵成为 CI 硬门禁；
9. 受控回滚已演练；
10. 相关 Feature batch 文档和项目进度按仓库规则闭环。

在这些条件之前，只能描述为某个 operation 的 authz pilot/local verification，不得声明 Gateway 全量鉴权迁移完成或 production ready。
