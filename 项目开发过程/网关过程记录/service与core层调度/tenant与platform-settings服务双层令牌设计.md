# tenant-service / platform-settings-service 双层令牌设计（服务访问 Core 层）

> 记录日期：2026-09-11
> 涉及服务：`tenant-service`、`platform-settings-service`、`auth-service`、`ani-gateway`
> 参照模板：`inference-service` 的双层凭证设计（`internal/runtime/coresdk/minter.go`）

---

## 1. 背景与动机

tenant-service 与 platform-settings-service 的业务实现是"薄 gRPC 层"：handler 收到网关转发来的 gRPC 请求后，通过 Core SDK（HTTP）调用 ani-gateway 的 `/admin/*` 平台管理面 API 完成实际读写。

此前两服务使用**静态 JWT**（K8s Secret `inference-c21-core-token`，TTL 3600s）作为 `CORE_API_TOKEN` 兜底凭证。问题：

- token 过期后服务→Core 全部调用 502（`core tenant api unavailable`），只能手动重签；
- 静态 token 泄漏面大、无法按服务撤销、无短时失效能力。

inference-service 已落地"双层凭证"：**动态 mint 优先 + 静态 token 兜底**。本次将同样的双层设计复制到 tenant-service 与 platform-settings-service（用户指令："tenant-service和platform-service也采用双层设计，硬编码的话就加入这两个服务"——即 mint 白名单硬编码扩展）。

---

## 2. 双层设计总览

```
                     ┌──────────────────────────── 第一层：动态 mint（优先） ───────────────────────────┐
                     │                                                                              │
gRPC handler          │  coreRequest(ctx, ...)                                                       │
    │                │    └─ applyAuthToken(ctx, headers)                                            │
    │                │         └─ coreMinter.Token(ctx)          ← 缓存单 token，剩余寿命 < 30s 刷新    │
    │                │              └─ gRPC IssueServiceToken → auth-service（TTL 300s 短时 JWT）     │
    │                │         └─ headers["Authorization"] = "Bearer <minted JWT>"（覆盖 SDK Token）  │
    │                │                                                                              │
    │                └──────────────────────────── 第二层：静态兜底（回退） ──────────────────────────┘
    │                     coreMinter == nil（AUTH_SERVICE_GRPC_ADDR / AUTH_SERVICE_MINT_SECRET 未配）
    │                     → 沿用 SDK 构造时的静态 CORE_API_TOKEN
    ▼
ani-gateway /admin/*（认证 → 授权 → Core handler）
```

决策要点：

| 决策 | 内容 |
|---|---|
| 切换条件 | `AUTH_SERVICE_GRPC_ADDR` + `AUTH_SERVICE_MINT_SECRET` **两 env 均非空**才启用 minter；任一为空保持静态兜底（部署上可灰度回退） |
| 缓存粒度 | **单 token 缓存**（inference-service 是 per-tenant map）。两服务的调用面全部是平台边界（无 tenant_id 上下文），无需按租户隔离 |
| TTL | 300s（mint 请求 `TtlSeconds`），剩余寿命 < 30s（`mintRefreshWindow`）提前刷新，避免边界过期 |
| 凭证注入方式 | anisdk.Client 的 `Token` 是**值字段**，构造后无法动态更新；但 SDK `Request()` 中 `options.Headers` 在 Token 之后应用（client.go L1072-1077），故 per-request `Headers["Authorization"]` 可覆盖 |
| mint 失败语义 | mint 失败 = `ports.ErrCoreUnavailable`（`CORE_UNAVAILABLE: mint core token: ...`），**不回退**静态 token——静态 token 已过期时静默回退会掩盖故障，宁 fail fast |

---

## 3. auth-service 侧改动（硬编码白名单扩展）

### 3.1 mint 调用方白名单

[auth_service.go](../../../repo/services/auth-service/internal/service/auth_service.go)：

```go
// IssueServiceToken 允许的内部调用方集合
allowedMintCallers = "inference-service,tenant-service,platform-settings-service"

// 访问 /admin/*（V2 policy 仅允许 PrincipalUser）的调用方集合：
// 签发"服务扮演 platform user"形态 token，而非 service 主体形态
platformAdminMintCallers = "tenant-service,platform-settings-service"
```

- `AUTH_SERVICE_MINT_CREDENTIALS` env（`name:secret,...` 多条目格式）经 `parseMintCredentials` 解析为 `s.mintSecrets[name]`；
- `mintAllowed` 用 `subtle.ConstantTimeCompare` 校验调用方 secret（防时序探测）。

### 3.2 Token 形态：为什么不是 service 形态

`IssueServiceToken` 有两种签发形态：

| 维度 | service 形态（inference-service 用） | **服务扮演 platform user 形态（本两服务用）** |
|---|---|---|
| `principal_kind` | `service` | `user` |
| `credential_domain` | tenant / platform | `platform` |
| `aud` | `ani-core` | （无） |
| `roles` | `["service"]` | `["platform-admin"]` |
| `scope`（legacy） | 第一个 permission | `platform` |
| `sub` | caller 服务名 | caller 服务名 |
| `uid` | magic UUID | magic UUID（`00000000-0000-0000-0000-0000000000aa`） |
| `tenant_id` | 可带 | 必须为空（带 tenant_id 直接 `InvalidArgument`） |
| 网关可达面 | 被 rbac.go 限制在 `/api/v1/platform-workloads*` | 走平台管理员授权链，可达 `/admin/*` |

原因：网关对 `/admin/*` 的 **V2 policy**（`zz_generated_core_policies.go`）`PrincipalKinds` 只允许 `user`，service 主体进不去授权链；而 **legacy policy** 面直接要求 `roles=[platform-admin]` 放行。所以为这两个服务签"服务扮演 platform user"形态——`roles`/`scope` 满足 legacy 链，`principal_kind=user` 满足 V2 链。

[token_issuer.go](../../../repo/services/auth-service/internal/service/token_issuer.go) 关键分支：

```go
if isPlatformAdminMintCaller(req.GetCallerService()) {
    if domain != "platform" || strings.TrimSpace(req.GetTenantId()) != "" {
        return jwtPayload{}, status.Error(codes.InvalidArgument,
            "platform admin caller must use platform credential domain without tenant")
    }
    return jwtPayload{
        PrincipalKind:    "user",
        TenantID:         "",
        Subject:          req.GetCallerService(),           // V2 权威服务身份
        UserID:           serviceActorUserID.String(),      // magic UUID，避免服务名泄漏到 legacy wire 契约
        CredentialDomain: "platform",
        Permissions:      permissions,                      // scope:tenants:* 等权威权限
        Scope:            "platform",                        // legacy 投影
        Roles:            []string{"platform-admin"},        // legacy 投影
    }, nil
}
```

### 3.3 V2 授权链缺口修复（关键）

**问题**：网关 V2 链 `CheckPermissionV2` 调 auth-service `permissionStore.Allows`；对 `user` 主体会查库：

```sql
SELECT r.permissions FROM users u JOIN user_roles ur ... WHERE u.id = $1
```

magic UUID（`00000000-...-0000aa`）在 users 表中**没有对应行**（users 种子只有 platform-admin 角色，无该用户）→ 空权限 → `PERMISSION_DENIED`。动态 mint 的 token 会被 V2 链拒绝。

**修复**（[permissions.go](../../../repo/services/auth-service/internal/service/permissions.go)）：

```go
func (s *permissionStore) Allows(...) (bool, error) {
    switch {
    case principal.Kind == "user" && isServiceActorPrincipal(principal):
        // 服务扮演的 platform user：SubjectID 是 magic UUID，users 表无对应行；
        // 权威权限来自签名 token 的 permissions，与 service 主体同信任模型，不查库。
        permissions, err = permissionsFromScopes(principal.Permissions, principal.Domain)
    case principal.Kind == "user":
        permissions, err = s.userPermissions(ctx, principal)   // 真人 user 才查库
    case principal.Kind == "api_key", principal.Kind == "service":
        permissions, err = permissionsFromScopes(principal.Permissions, principal.Domain)
    ...
}

func isServiceActorPrincipal(principal principalRecord) bool {
    return principal.SubjectID == serviceActorUserID.String() && principal.Domain == "platform"
}
```

判定条件是 `SubjectID == magicUUID && Domain == "platform"`：真人 platform user 的 SubjectID 是真实 UUID，不会误命中。

### 3.4 测试（service_token_test.go）

新增三个用例（fixture 补三服务 mintSecrets）：
- `TestIssueServiceTokenPlatformAdminCaller`：tenant-service / platform-settings-service 调用 → 断言 `Kind=user`、`Domain=platform`、`TenantID=""`、`sub=caller`、`uid=magic UUID`、`Legacy.Scope="platform"`、`Roles=[platform-admin]`、`Permissions` 内容；
- `TestIssueServiceTokenPlatformAdminCallerRejectsTenant`：带 tenant_id → `InvalidArgument`；
- `TestIssueServiceTokenServiceFormUnaffectedForInference`：回归——inference-service 仍签 service 形态（`Kind=service`、有 tenant）。

---

## 4. 两个服务侧实现

两服务实现结构完全一致，仅 mint caller 与 permissions 不同。

### 4.1 minter.go（新建）

[tenant-service minter.go](../../../repo/services/tenant-service/internal/repo/adapters/core/minter.go) / [platform-settings-service minter.go](../../../repo/services/platform-settings-service/internal/repo/adapters/core/minter.go)

```go
const (
    mintCaller        = "tenant-service"                    // 或 "platform-settings-service"
    mintTTLSeconds    = 300
    mintRefreshWindow = 30 * time.Second
)

// tenant-service：
//   /admin/tenants* GET/POST/PUT 与 /admin/quota-meta 是 V2 policy（resource=tenants/quota，需覆盖）
//   其余全部是 legacy policy（roles 放行，不读 permissions）
mintPermissions = "scope:tenants:*,scope:quota:read"

// platform-settings-service：
//   /admin/platform-users* 目前全部 legacy policy；permissions 冗余携带 users 域
//   平台边界权限，供 V2 审计链与后续 policy 迁移使用
mintPermissions = "scope:users:*"
```

`Minter` 结构：`client`（gRPC AuthServiceClient）+ `secret` + `sync.Mutex` 保护的 `cached/expiresAt`。核心方法：

```go
// Token 返回平台管理面 JWT；剩余寿命低于 mintRefreshWindow 时提前刷新
func (m *Minter) Token(ctx context.Context) (string, error) {
    now := m.now()
    m.mu.Lock()
    if m.cached != "" && now.Add(mintRefreshWindow).Before(m.expiresAt) {
        token := m.cached
        m.mu.Unlock()
        return token, nil
    }
    m.mu.Unlock()

    issued, err := m.client.IssueServiceToken(ctx, &authv1.IssueServiceTokenRequest{
        CallerService:    mintCaller,
        CallerSecret:      m.secret,
        CredentialDomain:  "platform",              // 平台管理面：不带 tenant
        Permissions:       strings.Split(mintPermissions, ","),
        TtlSeconds:        mintTTLSeconds,
    })
    ...
    m.mu.Lock()
    m.cached = token
    m.expiresAt = now.Add(ttl)
    m.mu.Unlock()
    return token, nil
}
```

`DialMinter(addr, secret)`：`grpc.NewClient` + `insecure.NewCredentials()`（集群内明文，mTLS 由后续迭代）。

### 4.2 sdk_client.go（凭证注入层）

[tenant-service sdk_client.go](../../../repo/services/tenant-service/internal/repo/adapters/core/sdk_client.go)：

```go
// 全局可选 minter；nil = 静态兜底。main.go 启动时注入。
var coreMinter *Minter

func SetupMinter(addr, secret string) (*Minter, error) {
    // 双 env 任一为空 → 返回 nil（保持静态兜底），非空 → DialMinter 并赋值全局
}

// 双层注入：minter 可用时 per-request 覆盖 Authorization
func applyAuthToken(ctx context.Context, headers map[string]string) (map[string]string, error) {
    if coreMinter == nil {
        return headers, nil                       // 第二层：静态兜底
    }
    token, err := coreMinter.Token(ctx)
    if err != nil {
        return nil, fmt.Errorf("%w: mint core token: %v", ports.ErrCoreUnavailable, err)
    }
    if headers == nil { headers = map[string]string{} }
    headers["Authorization"] = "Bearer " + token  // 第一层：动态 mint
    return headers, nil
}

// 所有 Core SDK 调用的统一入口
func coreRequest(ctx context.Context, sdk anisdk.Client, method, path string, opts anisdk.RequestOptions) (any, error) {
    headers, err := applyAuthToken(ctx, opts.Headers)
    if err != nil { return nil, err }
    opts.Headers = headers
    opts.Context = ctx
    return sdk.Request(method, path, opts)
}
```

覆盖原理：anisdk `Request()` 先用 `client.Token` 设 Authorization，**之后**应用 `options.Headers`（[client.go L1072-1077](../../../repo/sdks/core/go/anisdk/client.go)），故 per-request `Authorization` 覆盖 SDK 静态 Token 字段。

### 4.3 调用点接入

所有 `c.sdk.Request(...)` 批量替换为 `coreRequest(ctx, c.sdk, ...)`：

| 服务 | 客户端文件 | 调用点数 |
|---|---|---|
| tenant-service | tenant_svc_client.go / tenant_admin_svc_client.go / tenant_plan_svc_client.go / quota_svc_client.go | 30 处 |
| platform-settings-service | platform_user_client.go | 9 处 |

### 4.4 main.go 接线

[tenant-service main.go](../../../repo/services/tenant-service/main.go) / [platform-settings-service main.go](../../../repo/services/platform-settings-service/main.go)：

```go
func main() {
    cfg := config.Load()
    deps := bootstrap.MustConnect(cfg)          // DATABASE_URL → Postgres
    defer deps.Close()

    // 双层凭证：双 env 都配置时启用 minter，否则静态兜底
    if _, err := core.SetupMinter(
        os.Getenv("AUTH_SERVICE_GRPC_ADDR"),
        os.Getenv("AUTH_SERVICE_MINT_SECRET"),
    ); err != nil {
        log.Fatalf("setup core minter: %v", err)
    }
    ...  // 构造 client / service / Register gRPC
}
```

---

## 5. 网关授权链如何接受该 token

网关对 `/admin/*` 存在**双授权链**（由 policy 来源二分）：

```
请求 Authorization: Bearer <服务扮演 platform user JWT>
    │
    ├─ V2 链（zz_generated_core_policies.go 静态注册的路由）
    │    认证：验签 → principal{kind=user, domain=platform, sub=tenant-service, uid=magicUUID, permissions=[scope:tenants:*,...]}
    │    授权：CheckPermissionV2 → auth-service Allows()
    │          → isServiceActorPrincipal 命中 → permissionsFromScopes(token permissions) → 按 policy resource/action 匹配
    │
    └─ legacy 链（其余 /admin/* 路由）
         roles=[platform-admin] 直接放行（不读 permissions）
```

V2 policy 路由（需 token permissions 覆盖）与 legacy 路由的划分（调用面视角）：

- tenant-service：`/admin/tenants*` GET/POST/PUT、`/admin/quota-meta` → V2（Resource=tenants/quota，Action=read/create/update）→ 需 `scope:tenants:*`、`scope:quota:read` 覆盖；其余 → legacy
- platform-settings-service：`/admin/platform-users*` 目前全部 legacy（roles 放行）；permissions 冗余携带 `scope:users:*` 供后续迁移

---

## 6. K8s 部署接线

### 6.1 Secret（`inference-mint` 扩展为三服务格式）

| 键 | 值 | 消费方 |
|---|---|---|
| `credentials` | `inference-service:<sec1>,tenant-service:<sec2>,platform-settings-service:<sec3>` | auth-service `AUTH_SERVICE_MINT_CREDENTIALS` |
| `caller_secret` | `<sec1>` | inference-service `AUTH_SERVICE_MINT_SECRET` |
| `tenant_secret` | `<sec2>` | tenant-service `AUTH_SERVICE_MINT_SECRET` |
| `platform_settings_secret` | `<sec3>` | platform-settings-service `AUTH_SERVICE_MINT_SECRET` |

> Secret 更新只触发引用它的 Deployment 需显式重启（`kubectl rollout restart`）才生效——env 的 secretKeyRef 在 pod 启动时注入。

### 6.2 两个服务的 Deployment env

| env | 值 | 作用 |
|---|---|---|
| `AUTH_SERVICE_GRPC_ADDR` | `ani-auth-service.<ns>.svc.cluster.local:9101` | minter 拨号地址 |
| `AUTH_SERVICE_MINT_SECRET` | secretKeyRef `inference-mint/<各自键>` | mint 调用方凭证 |
| `CORE_API_BASE_URL` | `http://ani-gateway.<ns>.svc.cluster.local:8080/api/v1` | SDK base（兜底层共用） |
| `CORE_API_TOKEN` | secretKeyRef `inference-c21-core-token/token` | 静态兜底（已过期占位） |
| `DATABASE_URL` | secretKeyRef `ani-services-runtime/database_url` | 各自业务库 |

网关侧补 `TENANT_SERVICE_ADDR` / `PLATFORM_SETTINGS_SERVICE_ADDR`（缺省回落 `127.0.0.1:9105/9106` → `GRPC_CLIENT_UNAVAILABLE connection refused`）。

### 6.3 镜像

三镜像 tag `dev-20260911-dual`（构建机 192.168.18.35 构建，auth-service Dockerfile 补 `GOPROXY=https://goproxy.cn,direct`）：
- `docker.changqingyun.cn/ani/ani-auth-service:dev-20260911-dual`
- `docker.changqingyun.cn/ani/tenant-service:dev-20260911-dual`
- `docker.changqingyun.cn/ani/platform-settings-service:dev-20260911-dual`

部署环境：**ani-system**（NodePort 30080）与 **ani-test2**（NodePort 30083，含独立 Postgres/网关/auth-service，克隆 Deployment 时需将跨 namespace 地址改为本环境地址，否则 mint forbidden）。

---

## 7. 端到端验证

验证链路（两环境均通过）：

```
curl POST /api/v1/auth/platform/password/login（root 登录 → platform user JWT）
  → GET /api/v1/svc/tenants        → 网关 → tenant-service gRPC → coreRequest
      → minter → auth-service IssueServiceToken（300s，服务扮演形态）
      → Core /admin/tenants（网关 V2/legacy 链放行）→ 200
  → GET /api/v1/svc/platform-admins → 网关 → platform-settings-service gRPC → 同上 → 200
```

关键证据：静态 `CORE_API_TOKEN`（exp=1787294854）在验证时**已过期数周**，接口仍 200 → 证明实际走的是 minter 动态签发链路，且网关 V2 链接受了"服务扮演 platform user"形态。

---

## 8. 踩坑记录

| 坑 | 现象 | 解法 |
|---|---|---|
| V2 授权链缺口 | mint token 被 `CheckPermissionV2` 拒（magic UUID 不在 users 表） | `Allows` 加 `isServiceActorPrincipal` 分支，从 token permissions 读权威权限 |
| SDK Token 字段不可动态更新 | `anisdk.Client.Token` 是值字段 | per-request `Headers["Authorization"]` 覆盖（Headers 在 Token 后应用） |
| 构建机残留旧文件 | `AuditLog` 重复声明编译失败 | 上传前先删远程目标目录（scp 不清理） |
| 本地 Docker 构建网络 | GOPROXY 回落 proxy.golang.org 不可达 | auth Dockerfile 显式 `GOPROXY=goproxy.cn`，或走构建机 |
| 跨 namespace 克隆 Deployment | test2 服务仍指向 ani-system 的 auth/gateway → mint `PermissionDenied forbidden` | 改 `AUTH_SERVICE_GRPC_ADDR`/`CORE_API_BASE_URL` 为本环境 svc 地址 |
| 克隆 Service 带 clusterIP | `failed to allocate IP: already allocated` | 剥离 `clusterIP/clusterIPs` 后 apply |

---

## 9. 关键文件索引

| 层 | 文件 |
|---|---|
| auth 白名单/形态 | `repo/services/auth-service/internal/service/auth_service.go`（allowedMintCallers / platformAdminMintCallers） |
| auth 签发 | `repo/services/auth-service/internal/service/token_issuer.go`（resolveIssueServiceTokenClaims 平台管理面分支） |
| auth V2 权限 | `repo/services/auth-service/internal/service/permissions.go`（Allows 服务扮演分支 + isServiceActorPrincipal） |
| auth 测试 | `repo/services/auth-service/internal/service/service_token_test.go`（3 个新用例） |
| tenant minter | `repo/services/tenant-service/internal/repo/adapters/core/minter.go` |
| tenant 注入 | `repo/services/tenant-service/internal/repo/adapters/core/sdk_client.go`（SetupMinter/applyAuthToken/coreRequest） |
| tenant 接线 | `repo/services/tenant-service/main.go` |
| settings minter | `repo/services/platform-settings-service/internal/repo/adapters/core/minter.go` |
| settings 注入 | `repo/services/platform-settings-service/internal/repo/adapters/core/sdk_client.go` |
| settings 接线 | `repo/services/platform-settings-service/main.go` |
| 网关 policy | `repo/services/ani-gateway/internal/authz/zz_generated_core_policies.go`（V2/legacy 二分） |
| 网关 V2 链 | `repo/services/ani-gateway/internal/middleware/generated_authz.go` |
| SDK 覆盖点 | `repo/sdks/core/go/anisdk/client.go` L1072-1077（Headers 在 Token 后应用） |
| 设计模板 | `repo/services/inference-service/internal/runtime/coresdk/minter.go`（per-tenant map 版） |
