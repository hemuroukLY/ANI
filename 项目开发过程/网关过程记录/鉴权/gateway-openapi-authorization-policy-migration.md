# ANI Gateway OpenAPI 鉴权策略改造方案

> 文档状态：设计方案，尚未实施
>
> 创建日期：2026-08-20
>
> 首个试点接口：`GET /api/v1/metering/usage/platform`
>
> 目标：以 OpenAPI 为路由鉴权策略的唯一真实来源，Gateway 不再根据 URL 和 HTTP 方法猜测权限。
>
> 本文是独立设计文档，不表示相关代码、测试或生产验证已经完成。

---

## 一、结论

本次改造采用以下方案：

1. 在 OpenAPI operation 上新增结构化的 `x-ani-authz`，统一声明资源、动作、授权边界和允许的主体类型。
2. 构建期从 OpenAPI 生成 Gateway 静态鉴权策略，运行时不解析 YAML。
3. Gateway 将现有 `Auth + RBAC` 拆分为“解析策略、认证主体、执行授权”三个明确阶段。
4. 内部鉴权 Proto 追加 `boundary`、`operation_id`、`principal_kind`，并让平台主体合法地不携带 `tenant_id`。
5. 首先迁移 `GET /metering/usage/platform`；未迁移接口短期走显式 legacy fallback。
6. 全量迁移完成后删除 `scopeAllowedForPath`、`inferPermission` 和产品 API 的硬编码公开路径表。

不采用 `x-ani-principal-scope`。`scope` 在当前系统中已经同时指 JWT 主体域、API Key 权限和 RBAC 权限字符串，继续复用会放大歧义。新模型使用：

- `security`：凭证方案，例如 Bearer JWT、API Key；
- `principal_kinds`：主体类型，例如用户、服务身份、API Key；
- `boundary`：授权和数据边界，例如 own、tenant、platform；
- `resource + action`：具体业务权限。

---

## 二、当前代码及问题

### 2.1 当前执行链路

当前 Gateway 中间件顺序为：

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

`AuthWithClient` 同时承担两种职责：

1. 调用 auth-service 校验 Bearer Token 或 API Key；
2. 调用 `scopeAllowedForPath(path, scope)` 判断 tenant/platform/sandbox 主体能否访问该 URL。

`RBACWithClient` 再调用 `inferPermission(method, path)` 推断资源与动作，并请求 auth-service：

```go
resource, action := inferPermission(string(c.Method()), string(c.Path()))
resp, err := authClient.CheckPermission(ctx, &authv1.CheckPermissionRequest{
    TenantId: tenantID,
    UserId:   getStringValue(c, "user_id"),
    Roles:    getStringSliceValue(c, "roles"),
    Resource: resource,
    Action:   action,
})
```

当前推断规则主要是：

```go
// URL 中 v1 后的第一个 path segment 被当作 resource。
// GET/HEAD -> get，POST -> create，PUT/PATCH -> update，DELETE -> delete。
```

因此，当前代码实际存在三套可能互相不一致的鉴权描述：

| 来源 | 当前用途 | 问题 |
|---|---|---|
| OpenAPI `security` / `x-ani-rbac-scope` | 契约、SDK、文档、兼容基线 | Gateway 没有真正消费其权限语义 |
| `scopeAllowedForPath` | tenant/platform/sandbox 路由隔离 | 依赖 URL 前缀和人工白名单 |
| `inferPermission` | 生成 resource/action | 依赖 HTTP 方法和路径首段推断 |

### 2.2 计量平台接口的具体错误

OpenAPI 已声明：

```yaml
/metering/usage/platform:
  get:
    operationId: getPlatformMeteringUsage
    x-ani-rbac-scope: "scope:metering:platform:read"
```

但当前请求经过 Gateway 时会发生：

```text
GET /api/v1/metering/usage/platform
  -> Auth.scopeAllowedForPath
  -> 不匹配 /auth/platform/*、/platform/*、/admin/*
  -> 被归类为 tenant 路由
  -> platform token 在进入 RBAC 前返回 403
```

即使只在 `scopeAllowedForPath` 中追加该路径，使请求能够进入 RBAC，后续仍然是：

```text
inferPermission(GET, /api/v1/metering/usage/platform)
  -> resource = metering
  -> action = get
  -> CheckPermission(metering, get)
```

这与 OpenAPI 声明的 `metering + read + platform` 不一致。当前 `platform-admin` 因为角色直通可能得到允许，但系统并没有验证 OpenAPI 声明的精确权限；较小权限的平台角色也无法正确工作。

### 2.3 平台主体的 tenant_id 语义不正确

平台 JWT 的语义是 `scope=platform` 且没有 tenant ID。当前 `ValidateToken` 却使用 `claims.TenantID.String()`，会把 `uuid.Nil` 序列化为零 UUID 字符串。结果是：

- Gateway 和 auth-service 看到了一个形式有效、语义无效的 tenant ID；
- `CheckPermission` 的 `tenant_id required` 没有真正验证平台/租户边界；
- 平台权限是否生效依赖实现偶然行为，而不是显式契约。

目标实现必须让 platform 主体的 `tenant_id` 为空成为合法且明确的状态，并只在 tenant/own 边界强制 tenant ID。

### 2.4 现状与目标对照

| 维度 | 当前代码 | 目标代码 |
|---|---|---|
| 路由策略来源 | Gateway 手写路径规则 | OpenAPI `security + x-ani-authz` 生成 |
| 公开接口 | `isPublicPath` 硬编码 | OpenAPI `security: []`，基础设施健康路径除外 |
| 主体域判断 | `scopeAllowedForPath` 猜测 | policy boundary 与 Principal 显式比较 |
| 权限资源 | 从 URL 首段推断 | `x-ani-authz.resource` |
| 权限动作 | 从 HTTP 方法推断 | `x-ani-authz.action` |
| 主体类型 | 角色名或认证分支隐式表达 | `principal_kind` 明确传递 |
| 平台 tenant_id | 零 UUID 偶然通过 | 空值是 platform 主体的合法契约 |
| 未知路由策略 | 继续按默认规则推断 | fail closed |
| 迁移回退 | 无 | 仅未迁移显式清单；新策略 deny/error 不回退 |

---

## 三、目标鉴权模型

### 3.1 OpenAPI 声明

`GET /metering/usage/platform` 的目标声明如下：

```yaml
/metering/usage/platform:
  get:
    operationId: getPlatformMeteringUsage
    security:
      - BearerAuth: []

    # Core v1 兼容字段：已经进入兼容性基线，暂时保留原值。
    x-ani-rbac-scope: "scope:metering:platform:read"

    # 新的唯一运行时授权策略来源。
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
| `version` | 是 | 当前固定为 `v1`，用于后续演进和生成器拒绝未知结构 |
| `resource` | 是 | 稳定的业务资源名，不从 URL 推断 |
| `action` | 是 | 业务动作，不从 HTTP 方法推断 |
| `boundary` | 是 | `own`、`tenant` 或 `platform` |
| `principal_kinds` | 是 | `user`、`service`、`api_key`、`sandbox` 的非空数组 |

`security` 与 `x-ani-authz` 的职责不同：

- `security` 决定接受哪类凭证，例如此接口只接受 Bearer JWT；
- `principal_kinds` 进一步限制凭证代表的主体，例如 Bearer JWT 中只允许平台用户，不允许服务身份；
- `resource/action/boundary` 决定主体通过认证后还必须具备什么权限。

### 3.2 为什么使用结构化对象

不再为新接口使用 `scope:<resource>:<boundary>:<action>` 四段字符串，原因是：

1. `boundary` 与 `action` 是不同维度，编码到一个字符串后需要每个消费者自行解析；
2. 当前数据库权限模型本来就是 `{resource, actions, scope}` 三个字段；
3. 结构化字段可以被 OpenAPI 校验器、代码生成器和 CI 独立检查；
4. 后续增加主体类型不需要再次发明五段或六段字符串。

`x-ani-authz` 是 ANI 的 OpenAPI 扩展，不是 OpenAPI 行业标准字段。它符合 OpenAPI 的 Specification Extension 机制，但只有 ANI 自己的生成器和校验器理解其语义。

### 3.3 权限命名规范

#### Resource

- 使用小写 kebab-case，例如 `metering`、`vector-stores`、`gpu-inventory`；
- 表示稳定业务资源，不包含 `/api/v1` 等传输层信息；
- 不用 `admin-`、`platform-` 前缀表达授权边界；边界放在 `boundary`；
- 同一业务资源的 tenant/platform 接口使用同一个 resource。

示例：

```yaml
# 租户用量
resource: metering
action: read
boundary: tenant

# 平台跨租户用量
resource: metering
action: read
boundary: platform
```

#### Action

新接口优先使用以下命名：

| 类型 | 推荐动作 |
|---|---|
| 标准资源操作 | `read`、`create`、`update`、`delete` |
| 领域动作 | `search`、`invoke`、`execute`、`console`、`bind`、`proxy` |

新接口不再引入 `get`、`list`、`write`、`manage` 这类粒度不一致或语义过宽的动作。已有 Core v1 操作不在本批次批量改名；旧动作在迁移期原样兼容，后续通过独立权限迁移处理，不能在生成器中做不透明的全局映射。

#### Boundary

| 值 | 含义 | tenant_id 要求 |
|---|---|---|
| `own` | 仅主体拥有的资源 | 必须存在；对象所有权还需 handler/store 强制校验 |
| `tenant` | 当前租户范围 | 必须存在 |
| `platform` | 平台或跨租户范围 | 必须为空 |

Gateway 负责验证主体是否有进入该边界的资格。`own` 的对象所有权、tenant 资源过滤和数据库 RLS 仍由 handler/service/store 强制执行，不能仅依赖 Gateway。

#### Principal kind

| 值 | 含义 |
|---|---|
| `user` | 用户 JWT，包括 tenant 用户和 platform 用户 |
| `service` | 短期服务身份 JWT |
| `api_key` | API Key |
| `sandbox` | Sandbox 短期令牌 |

现有 `x-ani-principal-kind`、`x-ani-service-only` 在 Core v1 中先保留。某 operation 存在 `x-ani-authz` 时，运行时只使用 `x-ani-authz.principal_kinds`；旧字段只作为兼容元数据，CI 校验它们不得与新策略冲突。

---

## 四、生成策略与 Gateway 新链路

### 4.1 构建期生成

新增生成器：

```text
api/openapi/v1.yaml
  -> scripts/generate_gateway_authz.py
  -> services/ani-gateway/internal/authz/zz_generated_policies.go
```

运行时不加载或解析 OpenAPI YAML。生成代码包含已经解析好的 operation 策略：

```go
type Policy struct {
    OperationID    string
    Method         string
    PathTemplate   string
    Public         bool
    SecuritySchemes []SecurityScheme
    Resource       string
    Action         string
    Boundary       Boundary
    PrincipalKinds []PrincipalKind
}
```

生成器必须解析 OpenAPI 全局和 operation 级 `security` 的继承/覆盖关系：

- `security: []` 表示产品 API 公开；
- operation 未声明 `security` 时继承全局 security；
- 非公开 operation 必须声明合法的 `x-ani-authz`，或在迁移期进入显式 legacy 清单；
- 未知 `version`、boundary、principal kind 或空 resource/action 直接生成失败。

健康检查等不属于 OpenAPI 产品契约的基础设施路径可以保留独立的最小 allowlist。产品 API 是否公开必须来自 OpenAPI。

### 4.2 新中间件顺序

目标顺序：

```go
h.Use(
    RequestID(),
    ResolveAuthzPolicy(generatedRegistry),
    AuthenticateWithClient(authClient),
    AuthorizeWithClient(authClient),
    RateLimit(store),
    Idempotency(store),
    Audit(),
)
```

职责边界：

| 中间件 | 只负责 |
|---|---|
| `ResolveAuthzPolicy` | 用 method + 已注册 path template 找到生成策略 |
| `Authenticate` | 校验凭证并构造 Principal；失败返回 401 |
| `Authorize` | 校验 security scheme、principal kind、boundary、resource/action；失败返回 403 |

### 4.3 Principal

Gateway 内部统一主体结构：

```go
type Principal struct {
    Kind               PrincipalKind
    CredentialScheme   SecurityScheme
    CredentialBoundary Boundary
    TenantID           string
    UserID             string
    Roles              []string
}
```

现有 Proto 的 `TenantContext.scope` 在兼容期继续传输 `tenant|platform`，但进入 Gateway 后立即规范化为 `CredentialBoundary`。不再在业务代码中继续传播含义不清的 `scope` 名称。

`Authenticate` 的伪代码：

```go
func Authenticate(authClient AuthClient) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        policy := MustPolicy(c)
        if policy.Public {
            c.Next(ctx)
            return
        }

        principal, err := authenticateCredential(ctx, c, authClient)
        if err != nil {
            respond401(c)
            return
        }

        // 仅 tenant/own 主体注入租户数据库上下文。
        // platform 主体保持 TenantID 为空，平台 handler 使用平台事务边界。
        ctx, err = attachPrincipalContext(ctx, principal)
        if err != nil {
            respond401(c)
            return
        }
        SetPrincipal(c, principal)
        c.Next(ctx)
    }
}
```

### 4.4 Authorize

```go
func Authorize(authClient AuthClient) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        policy := MustPolicy(c)
        if policy.Public {
            c.Next(ctx)
            return
        }

        principal := MustPrincipal(c)
        if !policy.AllowsScheme(principal.CredentialScheme) ||
            !policy.AllowsPrincipalKind(principal.Kind) ||
            !boundaryCompatible(principal, policy.Boundary) {
            respond403(c, "principal not allowed by operation policy")
            return
        }

        result, err := authClient.CheckPermission(ctx, &authv1.CheckPermissionRequest{
            TenantId:     principal.TenantID,
            UserId:       principal.UserID,
            Roles:        principal.Roles,
            Resource:     policy.Resource,
            Action:       policy.Action,
            Boundary:     string(policy.Boundary),
            OperationId:  policy.OperationID,
            PrincipalKind: string(principal.Kind),
        })
        if err != nil || !result.GetAllowed() {
            respond403(c, denyReason(result, err))
            return
        }
        c.Next(ctx)
    }
}
```

`boundaryCompatible` 的固定语义是：

```go
func boundaryCompatible(principal Principal, required Boundary) bool {
    switch required {
    case BoundaryOwn, BoundaryTenant:
        return principal.CredentialBoundary == BoundaryTenant && principal.TenantID != ""
    case BoundaryPlatform:
        return principal.CredentialBoundary == BoundaryPlatform && principal.TenantID == ""
    default:
        return false
    }
}
```

关键规则：

- 不再调用 `scopeAllowedForPath`；
- 不再调用 `inferPermission`；
- 已迁移 operation 被新策略拒绝或执行出错时，绝不能回退到旧逻辑；
- 找不到策略时 fail closed，并记录 operation/method/path，不允许自动猜测。

---

## 五、内部鉴权契约与 auth-service

### 5.1 Proto 兼容扩展

Proto 只追加字段，不修改已有 tag：

```proto
// api/proto/common/v1/common.proto
message TenantContext {
  string tenant_id = 1;
  string user_id = 2;
  repeated string roles = 3;
  string scope = 4;          // legacy wire name: tenant|platform
  string principal_kind = 5; // user|service|api_key
}

// api/proto/auth/v1/auth_service.proto
message CheckPermissionRequest {
  string tenant_id = 1;
  string user_id = 2;
  repeated string roles = 3;
  string resource = 4;
  string action = 5;
  string boundary = 6;       // own|tenant|platform
  string operation_id = 7;   // audit/debug only, not a permission key
  string principal_kind = 8;
}
```

`ValidateToken` 必须返回真实主体类型：

- 普通/平台登录 JWT -> `user`；
- 服务身份 JWT -> `service`；
- API Key -> `api_key`；
- sandbox token 仍由 Gateway 本地校验并构造 `sandbox` Principal。

### 5.2 auth-service 授权逻辑

目标伪代码：

```go
func CheckPermission(req *CheckPermissionRequest) Decision {
    if !validResourceAction(req.Resource, req.Action) ||
        !validBoundary(req.Boundary) ||
        !validPrincipalKind(req.PrincipalKind) {
        return Deny("invalid authorization request")
    }

    switch req.Boundary {
    case "platform":
        if req.TenantId != "" {
            return Deny("platform principal must not carry tenant_id")
        }
        // platform-admin 的通配能力只在 platform 边界生效。
        return evaluatePlatformPermission(req)

    case "tenant", "own":
        if req.TenantId == "" {
            return Deny("tenant_id required")
        }
        // tenant-admin 的通配能力只在 tenant/own 边界生效。
        return evaluateTenantPermission(req)

    default:
        return Deny("unknown boundary")
    }
}
```

授权匹配最终对齐数据库既有权限结构：

```json
{
  "resource": "metering",
  "actions": ["read"],
  "scope": "platform"
}
```

以下越权必须被明确拒绝：

- `tenant-admin` 不能因为角色直通访问 platform boundary；
- `platform-admin` 不能自动作为 tenant 主体进入 tenant/own 路由；
- tenant API Key 即使角色列表中混入同名字符串也不能进入 platform boundary；
- auditor 的只读能力不能跨越其被授予的 boundary。

`operation_id` 只用于审计、日志和诊断，不能作为权限判断依据，避免 operationId 重命名改变授权语义。

---

## 六、计量接口迁移后的完整请求示例

平台管理员调用：

```http
GET /api/v1/metering/usage/platform?start_time=...&end_time=...
Authorization: Bearer <platform-user-jwt>
```

处理过程：

```text
1. ResolveAuthzPolicy
   -> operationId=getPlatformMeteringUsage
   -> scheme=BearerAuth
   -> resource=metering
   -> action=read
   -> boundary=platform
   -> principal_kinds=[user]

2. Authenticate
   -> auth-service ValidateToken
   -> Principal{
        Kind: user,
        CredentialScheme: bearer,
        CredentialBoundary: platform,
        TenantID: "",
        Roles: [platform-admin],
      }

3. Authorize 本地前置检查
   -> BearerAuth: 通过
   -> principal kind user: 通过
   -> credential boundary platform: 通过

4. auth-service CheckPermission
   -> resource=metering
   -> action=read
   -> boundary=platform
   -> platform-admin 或对应平台权限: 通过

5. handler
   -> 使用平台查询事务读取跨租户计量数据
   -> 返回 200
```

对照拒绝结果：

| 调用方 | 结果 | 原因 |
|---|---:|---|
| 无凭证 | 401 | 未认证 |
| tenant 用户 JWT | 403 | credential boundary 不匹配 |
| tenant API Key | 403 | security scheme、principal kind 和 boundary 均不满足 |
| platform service JWT | 403 | principal kind 不满足该 operation |
| platform 用户但无 `metering/read/platform` 权限 | 403 | 精确权限不满足 |
| platform 管理员 | 200 | 主体、边界和权限均满足 |

---

## 七、兼容与迁移方案

### 阶段 A：建立策略格式和生成门禁

新增：

- `x-ani-authz` 格式校验器及单元测试；
- OpenAPI -> Go 静态策略生成器；
- 生成物漂移检查；
- method/path template/operationId 与 Gateway 路由的一致性检查。

此阶段不改变现有接口运行行为。

### 阶段 B：迁移平台计量接口

1. 为 `GET /metering/usage/platform` 添加 `security` 和 `x-ani-authz`；
2. 保留已冻结的 `x-ani-rbac-scope: scope:metering:platform:read`；
3. 追加 Proto 字段并重新生成代码；
4. 实现 Principal、策略解析和新 Authorize 路径；
5. 修改 auth-service 的 boundary 授权逻辑；
6. 只让该 operation 使用新策略，其余接口保持旧行为。

迁移选择必须由生成策略中是否存在该 operation 决定：

```go
policy, migrated := generatedRegistry.Resolve(method, path)
if migrated {
    authorizeByGeneratedPolicy(policy) // deny/error 不得 legacy fallback
    return
}
authorizeByExplicitLegacyPolicy(method, path)
```

禁止采用“新鉴权失败后再试旧鉴权”的双重判定方式，否则新策略的拒绝可能被旧逻辑绕过。

### 阶段 C：全量迁移受保护 operation

- 为全部受保护 Core operations 添加 `x-ani-authz`；
- 已有 `x-ani-rbac-scope`、`x-ani-principal-kind`、`x-ani-service-only` 暂时保留；
- 对每个旧 action 保持原授权语义，不在同一批次顺便改名；
- CI 要求每个已注册产品路由恰好对应一个 OpenAPI operation 和一条生成策略；
- legacy fallback 数量必须只减不增。

### 阶段 D：删除旧逻辑

满足以下条件后删除：

- 所有受保护路由已由生成策略覆盖；
- 线上 legacy fallback 指标持续为零；
- tenant/platform/API key/service/sandbox 回归矩阵通过；
- 路由覆盖和生成物漂移已成为必过 CI。

随后删除：

- `scopeAllowedForPath`；
- `inferPermission`；
- 产品 API 的 `isPublicPath` 硬编码项；
- `x-ani-principal-kind` / `x-ani-service-only` 的运行时读取逻辑。

Core v1 中已经进入兼容基线的旧扩展字段可以继续保留为文档兼容信息；在未来 v2 契约中再统一清理。

---

## 八、预计改动范围

### 8.1 计量试点最小范围

| 模块 | 主要改动 |
|---|---|
| `api/openapi/v1.yaml` | 为平台计量 operation 增加 `security` 和 `x-ani-authz` |
| `api/proto/common/v1/common.proto` | 追加 `principal_kind` |
| `api/proto/auth/v1/auth_service.proto` | 追加 boundary、operation_id、principal_kind |
| Proto 生成物 | 按现有生成流程更新 |
| `scripts/` | 新增策略校验、生成、漂移检查 |
| Gateway `internal/authz/` | Policy 类型、生成 registry、path template 匹配 |
| Gateway middleware | 引入 Principal 和新 Authorize 路径 |
| auth-service | 返回 principal kind，按 boundary 校验权限 |
| 测试 | 生成器、middleware、auth-service、平台计量集成测试 |

这不是单行修改，但试点可以控制在鉴权契约和中间件边界内，不要求一次迁移全部 211 个 `x-ani-rbac-scope` operation。全量迁移的主要工作是 OpenAPI 策略补录和机械性测试，适合后续独立批次完成。

### 8.2 明确不在试点中完成

- 不批量重命名已有 Core v1 权限 action；
- 不删除兼容性基线中的 `x-ani-rbac-scope`；
- 不引入 OPA 或新的鉴权服务；
- 不把对象所有权和数据库 RLS 上移到 Gateway；
- 不在运行时读取 OpenAPI YAML；
- 不把新策略失败自动降级到 legacy 逻辑。

---

## 九、测试与验收标准

### 9.1 策略生成和契约

- `x-ani-authz` 缺少任一必填字段时生成失败；
- 未知 version、boundary、principal kind 时生成失败；
- 公开 operation 不生成资源权限；
- 受保护且已迁移的 operation 必须生成唯一策略；
- `GET /metering/usage/platform` 生成 `metering/read/platform/user`；
- Core v1 旧 `x-ani-rbac-scope` 保持不变，兼容性门禁通过；
- 生成物无漂移，Gateway 路由与 OpenAPI operation 一一对应。

### 9.2 Gateway

- 缺少/无效凭证返回 401；
- 合法凭证但 scheme、principal kind、boundary 或 permission 不匹配返回 403；
- platform Principal 保持空 tenant ID，不注入伪造的零 UUID tenant 上下文；
- tenant/own Principal 缺少有效 tenant ID 时返回 401；
- 已迁移 operation 的授权拒绝不会触发 legacy fallback；
- 未迁移 operation 在显式 legacy 清单内保持现有行为；
- 未注册、未声明或策略缺失的受保护 route fail closed。

### 9.3 auth-service 权限矩阵

| Principal | Boundary | 权限 | 预期 |
|---|---|---|---:|
| platform user / platform-admin | platform | metering/read | allow |
| platform user / readonly 精确权限 | platform | metering/read | allow |
| platform user / 无权限 | platform | metering/read | deny |
| tenant user / tenant-admin | platform | metering/read | deny |
| platform user / platform-admin | tenant | metering/read | deny |
| tenant API Key | platform | metering/read | deny |
| tenant user / tenant-admin | tenant | metering/read | allow |
| tenant user / 缺 tenant_id | tenant | metering/read | deny |

### 9.4 可观测性

至少记录以下低基数字段或指标：

- `operation_id`；
- `policy_source=generated|legacy`；
- `principal_kind`；
- `boundary`；
- `decision=allow|deny|error`；
- 规范化的 deny reason。

日志不得记录原始 Token、API Key 或完整敏感 claims。全量迁移完成的直接证据是 legacy 策略计数归零，而不是只看单元测试通过。

---

## 十、实施顺序建议

建议拆为四个可独立评审的批次：

1. **AUTHZ-POLICY-A**：定义 `x-ani-authz`、校验器、生成器和生成物，不切换运行时。
2. **AUTHZ-CONTRACT-B**：追加 Proto 字段，auth-service 返回规范化 Principal 并支持 boundary，但保持旧调用兼容。
3. **AUTHZ-METERING-C**：只切换 `getPlatformMeteringUsage`，完成正反权限矩阵和集成测试。
4. **AUTHZ-MIGRATION-D**：逐域补齐全部 operation，legacy 指标归零后删除路径/方法推断逻辑。

每个批次都必须先修改 OpenAPI 或内部 Proto 契约，再修改实现和测试；不得在同一个计量功能 PR 中直接删除全局 legacy 鉴权。

---

## 十一、关键代码索引

| 文件 | 当前职责 / 后续影响 |
|---|---|
| `api/openapi/v1.yaml` | Core REST 契约和新鉴权策略唯一来源 |
| `api/core-v1-compatibility-baseline.yaml` | 已冻结旧 `x-ani-rbac-scope`，试点不可静默改值 |
| `services/ani-gateway/internal/middleware/chain.go` | 当前 Auth/RBAC 中间件装配顺序 |
| `services/ani-gateway/internal/middleware/auth.go` | 当前认证、公开路径和 `scopeAllowedForPath` 混合逻辑 |
| `services/ani-gateway/internal/middleware/rbac.go` | 当前 `inferPermission` 和 CheckPermission 调用 |
| `api/proto/common/v1/common.proto` | 当前认证结果缺少 principal kind |
| `api/proto/auth/v1/auth_service.proto` | 当前权限请求缺少 boundary 和 operation ID |
| `services/auth-service/internal/service/auth_service.go` | 当前 tenant_id 前置校验、角色直通和 scope 字符串匹配 |
| `deploy/migrations/20260502_003_permissions_schema.sql` | 既有 `{resource, actions, scope}` 权限数据模型 |
| `scripts/validate_core_api_compatibility.py` | 当前会阻止修改已冻结的 `x-ani-rbac-scope` |

---

## 十二、规范定位与参考

- OpenAPI 3.1 允许使用 `x-` Specification Extensions，因此 `x-ani-authz` 是合法扩展，但不是跨产品通用字段：<https://spec.openapis.org/oas/v3.1.0.html#specification-extensions>
- OpenAPI Security Requirement 描述 operation 使用的安全方案；标准 scope 列表主要与 OAuth2/OpenID Connect security scheme 配合，不能替 ANI 表达 tenant/platform 数据边界：<https://spec.openapis.org/oas/v3.1.0.html#security-requirement-object>
- OAuth 2.0 的 scope 字符串由授权服务器定义，RFC 不规定 ANI 应采用何种 resource/action/boundary 命名：<https://www.rfc-editor.org/rfc/rfc6749.html#section-3.3>

因此，本方案的通用部分是“认证方案与授权策略分离、显式策略、默认拒绝、构建期校验”；`x-ani-authz` 的具体字段是 ANI 内部规范，需要由生成器、Gateway、auth-service 和 CI 共同保证一致性。
