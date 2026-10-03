# 开物接入 ANI Gateway 后端实现方案

> 版本：2026-09-22
> 范围：ANI Console / ANI BOSS 点击“开物”后的后端入口、鉴权、K8s Service/Secret 读取，以及 ANI Gateway 到开物 ClusterIP Service 的内部代理。
> 结论：可以实现。前端只拿到 ANI Gateway 自己的代理入口，不接触开物 `webToken`，也不接触 Kaiwu `ClusterIP`。

> 2026-09-22 最新代码复核结论：Kaiwu 是产品侧集成能力，API 应归入 ANI Services 层，而不是 Core API。本文已按最新代码将入口调整为 `/api/v1/svc/integrations/kaiwu/*`，OpenAPI 契约改为 `api/openapi/services/v1.yaml`，并把 Kubernetes 读取收敛到 adapter/port 边界，避免 Gateway handler 直接操作 K8s REST Client。

> 2026-09-28 决议（当前生效实现）：子路径反向代理方案已废弃并**从代码中删除**，改为**开物独占一个 origin，ANI Gateway 只负责鉴权 + 跳转**。原因、契约、部署与验证见下方 **第 0 节**；第 6 章与 7.2 节保留的只是历史决策记录。

## 0. 2026-09-28 决议：开物独占 origin（当前实现）

### 0.1 为什么放弃子路径反代

开物页面由 `dsh-web-app 0.1.2-rc.1` 渲染：页面带 `<base href="/">`，客户端把 API 基址硬编码为 `/api`，静态资源走 `/assets/*`、`/plugins/*` 等**根绝对路径**。挂到 `/kaiwu/console` 子路径后，浏览器发出的 `/assets/*`、`/plugins/*`、`/api/*` 全部落到 ANI Gateway 自己的根路径上（404 或 ANI 的 `/api/v1`），表现就是“开物页什么都没有”。该版本没有 base path 部署选项，子路径方案在开物侧无解。

因此当前实现为：开物通过 NodePort 或独立域名独占一个 origin，ANI Gateway 只做鉴权并下发开物 origin 的绝对入口地址。

### 0.2 入口 API 契约

授权通过后，`GET /api/v1/svc/integrations/kaiwu/console/entry` 与 `GET /api/v1/svc/integrations/kaiwu/boss/entry` 返回：

```json
{
  "client": "console",
  "entry_url": "http://10.10.1.66:30088/?token=<DSH webToken>",
  "expires_in": 300
}
```

1. `entry_url` 是**开物自身 origin 的绝对地址**，前端只需 `window.location.assign(entry_url)`，无需任何子路径改写；
2. DSH 收到 `?token=` 后立即 `303` 跳转到 `/`、下发自己的 `dsh-auth-*` Cookie，并返回 `referrer-policy: no-referrer`，token 不会停留在地址栏；
3. 该模式下 Gateway 不签发 `ani_kaiwu_bootstrap_*` / `ani_kaiwu_session_*` Cookie；
4. `expires_in` 来自 `KAIWU_ENTRY_TOKEN_TTL`，默认 120 秒；
5. 鉴权语义不变（仍校验 bearer 用户、scope、租户名/ID/状态、平台角色），运行时读取失败仍按 503 失败关闭。

### 0.3 K8s 侧要求

1. `kaiwu/kaiwu-console`、`kaiwu/kaiwu-boss` 保持 `ClusterIP`：Gateway 读取 ClusterIP 与 `webToken` 用，adapter 只接受 `ClusterIP` 类型。
2. 新增浏览器入口 Service（同一 selector，NodePort 或由独立域名 Ingress 承担）：

```text
kaiwu-console-public   NodePort   3080 -> 30088
kaiwu-boss-public      NodePort   3080 -> 30089
```

   清单见 `kaiwu-praxis/deploy/k8s/kaiwu-console.yaml` 与 `kaiwu-boss.yaml`。
3. `kaiwu-console` / `kaiwu-boss` 的 `KAIWU_TRUSTED_HOST` 各自只写自己 origin 的 authority（`10.10.1.66:30088` / `10.10.1.66:30089`），不要写 ANI 网关 authority，改完重启 Pod。
4. 使用独立域名时把 NodePort 换成 Ingress，并把 `KAIWU_*_PUBLIC_URL` 改成 `https://kaiwu-console.example.com` 这类 **origin**（不能带子路径）；HTTPS 由 Ingress 终结，`entry_url` 必须是 `https://`。
5. 读取 `webToken` 走开物侧提供的 `kaiwu/kaiwu-web-token-reader`，它只绑定 `ani-system/ani-gateway`；其他命名空间接入需自行追加绑定（见第 8 章）。

### 0.4 Gateway 环境变量

```text
KAIWU_CONSOLE_PUBLIC_URL=http://10.10.1.66:30088
KAIWU_BOSS_PUBLIC_URL=http://10.10.1.66:30089
KAIWU_ENTRY_TOKEN_TTL=5m
```

取值规则：只接受 `scheme://host[:port]`（允许末尾斜杠），带子路径、查询串、凭据或非法 scheme 时网关启动即失败；`KAIWU_*_PUBLIC_URL` 为空表示该客户端回退到子路径代理入口。

### 0.5 安全注意

1. `entry_url` 携带 DSH `webToken`，等同一次性启动凭据：不得写入日志、工单、截图或前端埋点，前端只做即时跳转。
2. 换来的 DSH Cookie 有效期 30 天且绑定 origin authority；建议对开物 NodePort/域名做网络层限制（内网 ACL、Ingress 鉴权），避免绕过 ANI 直接访问。
3. 若要求 token 完全不出现在浏览器，需要开物侧支持“服务端换一次性 code”，或让开物与 ANI 同 registrable domain 共享 Cookie Domain，属于后续可插拔增强项。

### 0.6 联调验证（2026-09-28，10.10.1.66）

```bash
# 1) 鉴权后取独占 origin 入口
curl -H "Authorization: Bearer <ANI token>" \
  http://10.10.1.66:30080/api/v1/svc/integrations/kaiwu/console/entry
# {"client":"console","entry_url":"http://10.10.1.66:30088/?token=...","expires_in":300}

# 2) 跟跳转：303 + dsh-auth Cookie
curl -i -c /tmp/kw.jar "http://10.10.1.66:30088/?token=<webToken>"
# 3) 带 Cookie 打开页面：200，且 <base href="/">
curl -o /tmp/kw.html -b /tmp/kw.jar http://10.10.1.66:30088/
# 4) 根路径静态资源：200
curl -o /dev/null -w '%{http_code}\n' -b /tmp/kw.jar http://10.10.1.66:30088/assets/index-*.js
```

实测结果：30080、30083 两个网关均返回绝对入口地址（`expires_in=300`）；`?token=` 返回 `303` 并下发 `dsh-auth-*`；`/` 返回 `200`（约 25 KB HTML，`<base href="/">`）；`/assets/index-*.js` 返回 `200`；`kaiwu-boss-public`（30089）同样 `303`/`200`。

## 1. 目标与边界

### 1.1 目标

1. ANI Console 点击“开物”后，只有 `tenant-a` 且认证上下文中的租户 ID 与数据库中 `tenant-a` 的 ID 一致，才允许进入开物 Console。
2. ANI BOSS 点击“开物”后，只有当前平台用户角色为 `platform-admin`，才允许进入开物 BOSS。
3. ANI Gateway 从 K8s 中读取对应 Service 和 Secret：
   - Console：`kaiwu/kaiwu-console`、`kaiwu/kaiwu-console-web-token`
   - BOSS：`kaiwu/kaiwu-boss`、`kaiwu/kaiwu-boss-web-token`
4. 前端只接收 ANI Gateway 自己的入口 URL，例如：

```text
/kaiwu/console
/kaiwu/boss
```

5. ANI Gateway 在集群内部通过 Kaiwu `ClusterIP Service` 转发 HTTP、SSE、WebSocket 和静态资源请求。
6. 不向浏览器返回或暴露 Kaiwu `ClusterIP`、DSH `webToken`、K8s Secret 原文或 ServiceAccount token。

### 1.2 非目标

本文不改变开物内部会话模型，不新增开物前端页面，也不实现开物内部租户隔离。本文只解决 ANI 侧入口授权、K8s Service/Secret 读取、ANI Gateway 代理入口和集群内 ClusterIP 转发。

当前 `kaiwu-console` 和 `kaiwu-boss` 是共享 DSH 实例。ANI Gateway Cookie 只做入口授权，不会把 DSH 内部的会话列表、工作区或文件目录按 ANI 用户拆分。如果仍要求“每个 ANI 用户在同一个开物部署内拥有独立会话”，必须另行实现 DSH 侧的 tenant/session 隔离插件，或恢复每用户实例；不能只靠本方案中的 URL、Cookie 和反向代理实现。

## 2. 现状与代码依据

### 2.1 认证上下文

ANI Gateway 已在认证中间件中写入 `tenant_id`、`user_id`、`roles` 和 `scope`。相关代码：

- `repo/services/ani-gateway/internal/middleware/auth.go`
- `repo/services/ani-gateway/internal/middleware/request_id.go`

可用读取函数：

```go
tenantID := middleware.GetTenantID(c)
userID := middleware.GetUserID(c)
```

实现时必须从认证上下文读取，不允许信任前端传入的 `tenant_id`、`user_id` 或 `client`。

### 2.2 Console 租户校验

Gateway 已注入 `ports.TenantService`，其中提供：

```go
GetTenant(ctx context.Context, tenantID string) (Tenant, error)
```

相关代码：`repo/pkg/ports/tenant.go`。

Console 入口应使用当前认证租户 ID 查询租户，并校验：

- 租户存在；
- `tenant.ID == middleware.GetTenantID(c)`；
- `tenant.Name == "tenant-a"`；
- `tenant.Status == "active"`。

当前主库中 `tenant-a` 的租户 ID 为 `00000000-0000-0000-0000-000000000001`。建议不要在代码中硬编码该 UUID，而是通过部署环境变量固定：

```text
KAIWU_CONSOLE_TENANT_NAME=tenant-a
KAIWU_CONSOLE_TENANT_ID=00000000-0000-0000-0000-000000000001
```

这样既满足“租户名称必须是 tenant-a”，也满足“租户 ID 必须一致”。如果未来环境重建导致 `tenant-a` 的 UUID 变化，只需修改环境变量，不需要改代码。

### 2.3 BOSS 平台角色校验

Gateway 已注入 `ports.PlatformUserAdminStore`，其中可查询平台账号：

```go
Get(ctx context.Context, userID uuid.UUID) (PlatformUserAdmin, error)
```

相关代码：`repo/pkg/ports/platform_user_admin.go`。

BOSS 入口应使用当前认证用户 ID 查询平台账号，并校验：

- 用户存在；
- `Status == "active"`；
- `Role == "platform-admin"`。

`platform-ops`、`platform-readonly` 均不允许访问开物 BOSS。

### 2.4 K8s 运行时读取边界

ANI Gateway 已有 `runtimeadapter.KubernetesRESTClient`，支持原始 K8s REST 调用：

```go
Do(ctx context.Context, method, endpoint, contentType string, body []byte) ([]byte, int, error)
```

相关代码：`repo/pkg/adapters/runtime/kubernetes_rest_client.go`。Gateway 在生产 Pod 内运行时，可通过 ServiceAccount token 和集群 CA 访问 K8s API。

最新代码中该客户端来自 GPU inventory / instance runtime 初始化；当 `GPU_INVENTORY_PROVIDER` 不是 `kubernetes_rest` 且 instance runtime 未提供客户端时，`RegisterOptions.KubernetesRESTClient` 可能为 `nil`。Kaiwu 不应复用这个隐式条件，应新增独立的 Kaiwu runtime reader：

```text
pkg/adapters/runtime/kubernetes_kaiwu.go
```

约束：

1. `pkg/adapters/runtime` 内可以读取 K8s Service/Secret；
2. `services/ani-gateway/internal/router/kaiwu_resources.go` 不得直接 import Kubernetes REST client；
3. Handler 只依赖注入的 Kaiwu runtime reader；
4. Reader 返回内部目标地址和已解码的 `webToken`，但不把 token 写入日志或 JSON。

## 3. API 设计

### 3.1 为什么拆成两个入口

Console 和 BOSS 的授权边界不同：Console 是 tenant boundary，BOSS 是 platform boundary。Services OpenAPI 的 `x-ani-authz` 对一个 operation 也只能声明一个 `boundary`。如果设计成单个混合入口，例如 `GET /api/v1/svc/integrations/kaiwu/entry?client=console|boss`，就无法在契约上精确表达 tenant 与 platform 两种边界，只能把细粒度判断全部放到 handler，降低契约清晰度。

Kaiwu 是产品侧第三方集成能力，按仓库最新架构规则应归入 ANI Services 层，而不是 Core API。因此推荐两个 Services API：

```text
GET /api/v1/svc/integrations/kaiwu/console/entry
GET /api/v1/svc/integrations/kaiwu/boss/entry
```

前端可以封装成一个函数：

```ts
getKaiwuEntry(client: "console" | "boss")
```

函数内部根据 `client` 调用对应路由。后端通过路由本身判断 Console 还是 BOSS，不信任查询参数、请求体或自定义 Header。

### 3.2 Services 鉴权语义

最新代码中，`/api/v1/svc/*` 路径在 Gateway 中仍走 legacy Auth + `CheckPermission`，尚未自动执行 Services OpenAPI 的 `x-ani-authz` 生成 policy。因此不能只依赖 `x-ani-authz` 声明，handler 必须继续做权威业务判断：

1. Console：要求 Bearer 用户、`scope == "tenant"`、当前租户为 `tenant-a` 且租户 ID 一致；
2. BOSS：要求 Bearer 用户、`scope == "platform"`、平台账号为 active `platform-admin`；
3. Services OpenAPI contract 要求 authenticated operation 同时声明 `BearerAuth` 与 `ApiKeyAuth` alternatives；本入口虽然按契约声明两者，但产品语义只允许 Bearer 用户，API Key 请求应返回 403 `KAIWU_USER_CREDENTIAL_REQUIRED`；
4. 需要给 middleware 增加只读的 credential scheme accessor，或复用已有 Principal 上下文，禁止 handler 重新解析 Authorization Header；
5. `x-ani-authz` 仍必须声明正确 boundary，作为契约与后续 V2 切换的稳定语义。

### 3.3 Console Entry

```text
GET /api/v1/svc/integrations/kaiwu/console/entry
Authorization: Bearer <ANI access_token>
```

OpenAPI authz 建议：

```yaml
security:
  - BearerAuth: []
  - ApiKeyAuth: []
x-ani-authz:
  version: v1
  resource: integration
  action: read
  boundary: tenant
  principal_kinds: [user]
```

成功响应：

```json
{
  "client": "console",
  "entry_url": "/kaiwu/console",
  "expires_in": 120
}
```

同时由 Gateway 设置一个短期、HttpOnly、SameSite 的 Console 代理 bootstrap Cookie。

失败响应：

| 状态码 | 业务码 | 场景 |
| --- | --- | --- |
| 401 | `UNAUTHORIZED` | 未登录或 token 无效 |
| 403 | `KAIWU_CONSOLE_TENANT_NOT_ALLOWED` | 当前租户不是 `tenant-a`，或租户 ID 不一致 |
| 403 | `KAIWU_TENANT_NOT_ACTIVE` | 租户不是 active |
| 503 | `KAIWU_BACKEND_UNAVAILABLE` | Kaiwu Service 或 Secret 不存在、不可读、数据不完整 |

### 3.4 BOSS Entry

```text
GET /api/v1/svc/integrations/kaiwu/boss/entry
Authorization: Bearer <ANI access_token>
```

OpenAPI authz 建议：

```yaml
security:
  - BearerAuth: []
  - ApiKeyAuth: []
x-ani-authz:
  version: v1
  resource: integration
  action: read
  boundary: platform
  principal_kinds: [user]
```

成功响应：

```json
{
  "client": "boss",
  "entry_url": "/kaiwu/boss",
  "expires_in": 120
}
```

同时设置 BOSS 代理 bootstrap Cookie。

失败响应：

| 状态码 | 业务码 | 场景 |
| --- | --- | --- |
| 401 | `UNAUTHORIZED` | 未登录或 token 无效 |
| 403 | `KAIWU_BOSS_ROLE_REQUIRED` | 当前平台用户不是 `platform-admin` |
| 403 | `KAIWU_PLATFORM_USER_NOT_ACTIVE` | 平台用户不是 active |
| 503 | `KAIWU_BACKEND_UNAVAILABLE` | Kaiwu Service 或 Secret 不存在、不可读、数据不完整 |

### 3.5 前端调用方式

ANI Console 调用 `/api/v1/svc/integrations/kaiwu/console/entry`，ANI BOSS 调用 `/api/v1/svc/integrations/kaiwu/boss/entry`。收到响应后：

```ts
window.location.assign(response.entry_url);
```

或新窗口打开：

```ts
window.open(response.entry_url, "_blank", "noopener");
```

URL 中不需要携带开物 token，也不需要携带 `user_id`。

当前 `ani-console/src/api/request.ts` 与 `ani-boss-console/src/api/request.ts` 已经为 Axios 设置 `withCredentials: true`，这对本方案是必要条件。部署时还必须满足：

1. entry API 与 `/kaiwu/*` 代理入口使用同一个浏览器站点；最简单做法是都通过 ANI Gateway 同一 origin 提供；
2. 如果二者跨 origin，需要配置 CORS、`Access-Control-Allow-Credentials`，且不能使用通配符 `*`；
3. `SameSite=Lax` 只保证同站顶层导航携带 Cookie，不解决跨 origin CORS 的凭证问题。

## 4. 鉴权流程

### 4.1 Console 流程

```text
前端 Console
    |
    | Bearer ANI token
    v
GET /api/v1/svc/integrations/kaiwu/console/entry
    |
    |- Auth middleware 校验 Bearer 用户
    |- Handler 校验 credential scheme 必须是 bearer
    |- Handler 校验 scope == tenant
    |- 契约声明 boundary=tenant；当前 legacy Services 链路由 handler 兜底执行
    |- 读取 middleware.GetTenantID(c)
    |- TenantService.GetTenant(tenantID)
    |- 校验 tenant.ID == 当前 tenantID
    |- 校验 tenant.Name == tenant-a
    |- 校验 tenant.Status == active
    |- 校验环境变量 KAIWU_CONSOLE_TENANT_ID == 当前 tenantID
    |- 读取 K8s Service kaiwu-console
    |- 读取 K8s Secret kaiwu-console-web-token
    |- 校验 Service ClusterIP 和 3080 端口存在
    |- 校验 Secret webToken 非空
    |- 签发短期 ANI Gateway bootstrap Cookie
    v
返回 entry_url=/kaiwu/console
```

### 4.2 BOSS 流程

```text
前端 BOSS
    |
    | Bearer ANI token
    v
GET /api/v1/svc/integrations/kaiwu/boss/entry
    |
    |- Auth middleware 校验 Bearer 用户
    |- Handler 校验 credential scheme 必须是 bearer
    |- Handler 校验 scope == platform
    |- 契约声明 boundary=platform；当前 legacy Services 链路由 handler 兜底执行
    |- 读取 middleware.GetUserID(c)
    |- PlatformUserAdminStore.Get(userID)
    |- 校验 Status == active
    |- 校验 Role == platform-admin
    |- 读取 K8s Service kaiwu-boss
    |- 读取 K8s Secret kaiwu-boss-web-token
    |- 校验 Service ClusterIP 和 3080 端口存在
    |- 校验 Secret webToken 非空
    |- 签发短期 ANI Gateway bootstrap Cookie
    v
返回 entry_url=/kaiwu/boss
```

### 4.3 禁止事项

1. 不允许前端传 `tenant_id`。
2. 不允许前端传 `user_id`。
3. 不允许只根据前端传来的 `client` 判断权限。
4. 不允许 API Key 调用 Kaiwu entry。
5. 不允许把 `ClusterIP` 放入 JSON 响应。
6. 不允许把 DSH token 写入日志、工单、埋点或持久化存储；按第 0 节决议，token 只允许作为入口地址中的一次性启动参数交给浏览器即时跳转。

## 5. K8s Service 与 Secret 读取

### 5.1 固定资源名

```text
Namespace: kaiwu

Console Service: kaiwu-console
Console Secret:  kaiwu-console-web-token

BOSS Service:     kaiwu-boss
BOSS Secret:      kaiwu-boss-web-token
```

Service 端口：

```text
name: http
port: 3080
targetPort: http
type: ClusterIP
```

### 5.2 K8s API 请求

Console：

```text
GET /api/v1/namespaces/kaiwu/services/kaiwu-console
GET /api/v1/namespaces/kaiwu/secrets/kaiwu-console-web-token
```

BOSS：

```text
GET /api/v1/namespaces/kaiwu/services/kaiwu-boss
GET /api/v1/namespaces/kaiwu/secrets/kaiwu-boss-web-token
```

这些 K8s API 请求只能在 `pkg/adapters/runtime/kubernetes_kaiwu.go` 中通过现有 `KubernetesRESTClient.Do` 发起。`services/ani-gateway/internal/router/kaiwu_resources.go` 不直接拼装 K8s API URL，也不直接 import Kubernetes REST client。

### 5.3 Service 解析

1. 读取 `spec.clusterIP`；
2. 读取 `spec.ports`；
3. 找到 `name=http` 或 `port=3080` 的端口；
4. 拒绝 `clusterIP == None` 或空值；
5. 构造内部目标 `http://<clusterIP>:3080`。

`ClusterIP` 只在 Gateway 进程内使用，不返回给前端。

### 5.4 Secret 解析

Secret 数据结构包含 `data.webToken`、`data.podName` 和 `data.instanceId`。K8s Secret 的 `data` 字段是 base64 编码。Gateway 需要：

1. 读取 `data.webToken`；
2. base64 decode；
3. trim；
4. 校验非空；
5. 只保存在请求处理内存中；
6. 不写入日志；
7. 不写入 JSON 响应。

`podName` 和 `instanceId` 可用于诊断，但本期不作为授权依据。

### 5.5 缓存策略

Service 元数据可以缓存 5–10 秒，用于减少 Service GET 压力；DSH `webToken` 完全不做缓存。Gateway 每次 entry 授权成功后、以及每次 proxy 需要执行 DSH token 交换前，都必须重新执行 K8s Secret GET，实时读取当前 `data.webToken` 并 base64 decode。这样 Kaiwu Pod 重启、Secret 更新或手动轮换后，Gateway 不会继续使用旧 token。

## 6. ANI Gateway 代理设计（2026-09-28 已废弃，代码已删除）

> 本章描述的 `/kaiwu/*` 子路径反向代理、代理 Cookie（`ani_kaiwu_bootstrap_*` / `ani_kaiwu_session_*`）与 DSH token 交换已在 2026-09-28 随“开物独占 origin”决议整体删除：`internal/router/kaiwu_proxy.go`、`kaiwu_cookies.go`、`kaiwu_proxy_runtime.go`、`kaiwu_cookie_runtime.go` 及对应测试文件均已移除，`RegisterOptions` 也不再接受 `KaiwuCookieSigner` / `KaiwuProxyHTTPClient`。以下内容仅作为历史决策记录保留，实施口径以第 0 节为准。

### 6.1 入口 URL

```text
/kaiwu/console
/kaiwu/console/*

/kaiwu/boss
/kaiwu/boss/*
```

这些 URL 是 ANI Gateway 自己的代理入口，不是 Kaiwu Service 的 NodePort。

### 6.2 为什么需要 Gateway 代理 Cookie

浏览器从 ANI 前端跳转到 `/kaiwu/console` 时，不会自动带上 `Authorization: Bearer ...` Header，因为这是新的顶层导航请求。因此 entry API 需要先给浏览器设置一个 ANI Gateway 自己的 HttpOnly Cookie：

```text
1. 前端带 Bearer 调用 entry API
2. Gateway 完成 Console/BOSS 鉴权
3. Gateway 设置短期 HttpOnly bootstrap Cookie
4. Gateway 返回 entry_url
5. 前端跳转 entry_url
6. 浏览器自动携带 bootstrap Cookie
7. Gateway 校验 Cookie
8. Gateway 读取 DSH webToken Secret
9. Gateway 用 webToken 向 DSH 完成首次 token 交换
10. Gateway 下发 DSH 会话 Cookie，并转发页面
```

这个 Cookie 是 ANI Gateway 签发的，不是 DSH 原始 token。

### 6.3 Cookie 设计

Console bootstrap Cookie：

```text
Name: ani_kaiwu_bootstrap_console
Path: /kaiwu/console
HttpOnly: true
SameSite: Lax
Secure: true  # HTTPS 入口必须开启；当前 HTTP 测试环境必须关闭，上线 HTTPS 后开启
Max-Age: 120
```

BOSS bootstrap Cookie：

```text
Name: ani_kaiwu_bootstrap_boss
Path: /kaiwu/boss
HttpOnly: true
SameSite: Lax
Secure: true  # HTTPS 入口必须开启；当前 HTTP 测试环境必须关闭，上线 HTTPS 后开启
Max-Age: 120
```

代理会话 Cookie 使用 `ani_kaiwu_session_console` 和 `ani_kaiwu_session_boss`。Cookie claims 使用 HMAC-SHA256 签名，建议包含：

```json
{
  "typ": "kaiwu_proxy",
  "client": "console",
  "scope": "tenant",
  "tenant_id": "00000000-0000-0000-0000-000000000001",
  "user_id": "用户ID",
  "issued_at": 1760000000,
  "expires_at": 1760000300,
  "jti": "随机数"
}
```

要求：

1. `client` 必须与当前代理路径一致；
2. Console Cookie 只能访问 `/kaiwu/console`；
3. BOSS Cookie 只能访问 `/kaiwu/boss`；
4. Console Cookie 必须绑定 `scope=tenant`、`tenant_id` 和 `user_id`；
5. BOSS Cookie 必须绑定 `scope=platform` 和 `user_id`，不得伪造或要求 `tenant_id`；
6. 签名密钥来自环境变量，不能使用进程随机默认值；
7. 多副本 Gateway 必须共享同一签名密钥。

以上 Cookie 机制已删除，不再需要 `KAIWU_PROXY_SIGNING_SECRET` 等代理签名变量（见第 0 节）。

### 6.4 DSH token 交换

当用户第一次打开 `/kaiwu/console` 或 `/kaiwu/boss` 时：

1. 校验 ANI Gateway bootstrap Cookie；
2. 每次都实时 GET K8s Secret，读取 DSH `webToken`，不使用任何进程内缓存、内存缓存或本地缓存；
3. 构造内部请求：

```text
GET http://<kaiwu-cluster-ip>:3080/?token=<webToken>
Host: <浏览器实际访问的 ANI Gateway authority>
```

4. HTTP client 禁止自动跟随 redirect；
5. 读取 DSH 返回的 `Set-Cookie`；
6. 将 DSH Cookie 下发给浏览器；
7. 设置 ANI Gateway session Cookie；
8. 302 跳转到不带 token 的代理路径；
9. 后续请求由 Gateway 代理到 Kaiwu Service。

注意：

- `Host` 必须是浏览器实际访问的 ANI Gateway authority；
- DSH launch token 在当前实现中是进程生命周期内固定的 token，不是一次性 token；同一个 Secret `webToken` 可以为多个 ANI 用户分别完成首次 token 交换；
- Kaiwu Pod 重启后新的启动流程会更新 Secret，Gateway 不能缓存 `webToken`，否则会把旧 token 发给新的 DSH 进程；
- token 交换请求必须访问 DSH 内部根路径 `/`，不能访问 `/kaiwu/console/` 或 `/kaiwu/boss/`；`/kaiwu/*` 只是 ANI Gateway 的外部代理前缀；
- DSH token 交换成功后返回 `303`，`Location` 通常是 `/`，并下发 DSH Cookie。Gateway 不能跟随该重定向，也不能把 `Location: /` 原样返回给浏览器，必须改写为 `/kaiwu/console` 或 `/kaiwu/boss`；
- 浏览器最终 URL、entry API JSON 和日志中都不能出现 `webToken`；
- Kaiwu Pod 的 `KAIWU_TRUSTED_HOST` 必须配置为同一个 authority；
- 否则 DSH browser-trust fence 会拒绝请求。

### 6.5 路径重写

外部路径 `/kaiwu/console/assets/foo.js` 应转发为内部路径 `/assets/foo.js`。规则：

1. 去掉 `/kaiwu/console` 或 `/kaiwu/boss` 前缀；
2. 保留 query；
3. 保留 WebSocket 子协议 Header；
4. 重写 `Location` 响应头；
5. 重写或隔离 `Set-Cookie`；
6. 检查 DSH 是否生成根路径静态资源引用；如不支持 base path，需要改为子域名入口或为 DSH 增加可配置 base path。

### 6.6 Header 处理

转发到 DSH 时必须设置：

```text
Host: <浏览器实际访问的 ANI Gateway authority>
X-Forwarded-Host: <浏览器实际访问的 ANI Gateway authority>
X-Forwarded-Proto: http 或 https
X-Forwarded-For: <客户端 IP>
```

必须保留 WebSocket Header：`Upgrade`、`Connection`、`Sec-WebSocket-Key`、`Sec-WebSocket-Version`、`Sec-WebSocket-Protocol`、`Sec-WebSocket-Extensions`。

必须删除或不转发 `Authorization`，以及 Cookie 中属于 ANI Gateway 的签名 Cookie。这里不转发 ANI `Authorization` 不是丢失认证，因为代理认证已由 ANI Gateway Cookie 完成。

### 6.7 Cookie 隔离

Console 和 BOSS 都使用 DSH，如果它们部署在同一个浏览器 origin 下，DSH Cookie 名可能冲突。必须二选一：

1. 推荐：Gateway 对 DSH Cookie 增加 client 前缀，例如 `kaiwu_console_<DSH cookie name>` 和 `kaiwu_boss_<DSH cookie name>`。出站时去掉前缀，DSH 看到原始 Cookie 名；入站时加回前缀，浏览器侧两个实例互不覆盖。
2. 或者：为 Console 和 BOSS 使用不同子域名。

第一版建议实现方案 1，保持单域名、单 Gateway 入口。

使用方案 1 时，Gateway 必须完整改写 DSH 返回的 `Set-Cookie`：

1. Cookie name 加 client 前缀，例如 `kaiwu_console_<DSH cookie name>`、`kaiwu_boss_<DSH cookie name>`；
2. `Path=/` 改写为 `Path=/kaiwu/console` 或 `Path=/kaiwu/boss`；
3. 保留必要的 `HttpOnly` 和 `SameSite` 属性；生产 HTTPS 入口必须保留或补齐 `Secure`；
4. 浏览器请求进入 Gateway 后，只把带前缀的 DSH Cookie 还原成原始名称转发给对应 DSH 实例；
5. 不转发 ANI Gateway 的 bootstrap/session Cookie 给 DSH；
6. DSH Cookie 是实例级认证，不包含 ANI user ID、tenant ID 或 session owner，不能作为用户隔离依据；
7. 每一个 `/kaiwu/*` 请求都必须继续校验 ANI Gateway 签名 Cookie，并确认其中的 user/tenant/client 与入口授权一致。

不能只依赖 DSH Cookie 判定访问者身份。DSH Cookie 只证明“该浏览器曾被 Gateway 引导到这个 DSH 实例”，不能证明“当前 ANI 用户仍有权访问该实例”。

### 6.8 流式与 WebSocket

代理层必须支持 HTTP request/response、SSE / streaming response、WebSocket upgrade、大文件下载和上传，并逐字节 flush，避免 SSE 被网关缓冲。

如果 Hertz 路由层无法稳定支持 WebSocket hijack / upgrade，可以在 `ani-gateway` 进程内增加一个独立的 Go `net/http` reverse proxy listener，并由同一个 ANI Gateway Service / Ingress 将 `/kaiwu/*` 路由到该 listener。对外仍然是 ANI Gateway 的统一入口。

## 7. 代码落点

建议新增：

```text
repo/services/ani-gateway/internal/router/kaiwu_resources.go
repo/services/ani-gateway/internal/router/kaiwu_resources_test.go
repo/services/ani-gateway/internal/router/kaiwu_proxy.go
repo/services/ani-gateway/internal/router/kaiwu_proxy_test.go
repo/pkg/adapters/runtime/kubernetes_kaiwu.go
repo/pkg/adapters/runtime/kubernetes_kaiwu_test.go
```

### 7.1 API Handler

`kaiwu_resources.go` 负责：

1. 在 Services route group 注册 `svc.GET("/integrations/kaiwu/console/entry", api.consoleEntry)` 和 `svc.GET("/integrations/kaiwu/boss/entry", api.bossEntry)`；
2. Console 授权；
3. BOSS 授权；
4. 调用注入的 Kaiwu runtime reader 读取 K8s Service/Secret；
5. 按独占 origin 基址拼接 `entry_url`（带一次性启动 token），不签发任何 Cookie。

Handler 文件不得 import `pkg/adapters/runtime`，不得持有 `KubernetesRESTClient`。也不建议为该过渡能力新增 `pkg/ports/kaiwu.go` 后再让 Services handler import；应在 `kaiwu_resources.go` 中定义本地 `KaiwuRuntimeReader` 接口，并通过 `RegisterOptions.KaiwuRuntimeReader` 注入实现。reader 的实现和 K8s REST 调用全部留在 adapter 层。

该改动触碰 Gateway shared/mixed handler 与 Core adapter，需要按 CODEOWNERS 申请 Core/Services 共同 review。

### 7.2 Proxy Handler（已删除）

`kaiwu_proxy.go`（子路径反代、Cookie 校验、DSH token 交换、路径重写、ClusterIP 转发、SSE/WebSocket 透传）已随独占 origin 决议整体删除，见第 0 节。Gateway 现在只有 `kaiwu_resources.go` 一组入口 handler 与 `pkg/adapters/runtime/kubernetes_kaiwu.go` 运行时读取器。

### 7.3 OpenAPI 与生成物

修改 `repo/api/openapi/services/v1.yaml`，在 `/integrations/kaiwu/console/entry` 与 `/integrations/kaiwu/boss/entry` 下新增两个 operation：

```text
getKaiwuConsoleEntry
getKaiwuBossEntry
```

并添加 response schema：

```yaml
KaiwuEntryResponse:
  type: object
  required: [client, entry_url, expires_in]
  properties:
    client:
      type: string
      enum: [console, boss]
    entry_url:
      type: string
      description: ANI Gateway 自己的代理入口，不包含 Kaiwu webToken 或 ClusterIP
    expires_in:
      type: integer
      description: bootstrap Cookie 有效秒数
```

Services 层不修改 `repo/api/openapi/v1.yaml`，也不生成 Core authz policy。需要同步维护：

```text
repo/sdks/services/
repo/docs/api/
```

最小验证命令：

```bash
cd repo
make validate-services
git diff --check
```

其中 Services route contract 会校验 OpenAPI 路径与 Gateway 注册路径一致；新增路径不得登记为 `services-route-baseline.yaml` exception。

### 7.4 与现有中间件的关系

入口 API `/api/v1/svc/integrations/kaiwu/console/entry` 和 `/api/v1/svc/integrations/kaiwu/boss/entry` 继续走现有 Auth、Authz、RateLimit、Idempotency、Audit 中间件。

代理路径 `/kaiwu/console` 和 `/kaiwu/boss` 是浏览器顶层导航，不会携带 Bearer Header，因此需要：

1. 新增专用的 `isKaiwuProxyPath(path)`，并让 Auth/RBAC/RateLimit 对 `/kaiwu/console`、`/kaiwu/boss` 及其子路径跳过全局 Bearer 认证；
2. 在 proxy handler 内部强制校验签名 Cookie；
3. 对无 Cookie 或无效 Cookie 请求返回 401；
4. proxy handler 内实现自定义限流；
5. proxy handler 内实现自定义审计，不记录 Cookie 和 query token。

最新中间件顺序中还包含 `Idempotency`。代理路径上的 DSH POST/PUT/PATCH/DELETE 不属于 ANI JSON API 幂等操作，如果请求意外携带 `Idempotency-Key`，会被全局幂等中间件缓存或阻断。因此 `Idempotency` 也必须对 `isKaiwuProxyPath(path)` 直接放行，避免破坏流式响应和 WebSocket。

“public path”只表示跳过全局 Bearer 中间件，绝不表示公开访问。代理 handler 必须 fail closed。

## 8. K8s 访问权限

> 2026-09-28 更新：本节按“最小权限绑定”口径重写。开物侧现在提供
> `deploy/k8s/kaiwu-web-token-rbac.yaml`，子路径代理时代“沿用既有权限”的说法已作废。

开物侧提供 `kaiwu/kaiwu-web-token-reader`：只对 `kaiwu-console-web-token`、`kaiwu-boss-web-token` 授予 `get`，并把 RoleBinding **只绑定 `ani-system/ani-gateway`**。入口 API 的最小权限集合就是“按名字读这两个 Secret”，多绑定命名空间等于扩大可读秘密范围，因此：

1. 网关部署在 `ani-system` 时无需额外授权；
2. 网关部署在其他命名空间（测试用 `ani-test2`、`ly-test` 等）时，需在该命名空间自行追加等价 RoleBinding，不应依赖集群级兜底角色；
3. Gateway 不需要 `list` / `watch`：适配器每次调用都按名字读单个 Service 和 Secret，不做缓存。

## 9. 部署调整

### 9.1 Kaiwu Deployment

> 2026-09-28 更新：本节按“开物独占 origin”口径重写，子路径代理时代的要求已作废，详见第 0 节。

Kaiwu Console 和 BOSS 的服务保持 `ClusterIP`（`3080`），供 ANI Gateway 读取运行时；浏览器入口由额外的 NodePort/域名 Service 提供（`kaiwu-console-public: 30088`、`kaiwu-boss-public: 30089`）。每个 Deployment 的 `KAIWU_TRUSTED_HOST` 只写自己的 origin authority。示例：

```text
# 只能写 host:port，不写 scheme
# kaiwu-console
KAIWU_TRUSTED_HOST=10.10.1.66:30088
# kaiwu-boss
KAIWU_TRUSTED_HOST=10.10.1.66:30089
```

要求：

1. 每个 Deployment 只写「浏览器实际访问自己」的 authority：`kaiwu-console` 写 `10.10.1.66:30088`，`kaiwu-boss` 写 `10.10.1.66:30089`，两者不再要求一致；
2. 不要写 ANI 网关 authority（`10.10.1.66:30080` / `30083` / `30093` 等）：浏览器不会带着这些 Host 访问开物，写进去只会放宽 `/api/` 的 browser-trust fence；
3. 同一实例存在多个入口（NodePort + 独立域名）时全部列入，空格分隔；
4. 如果使用 HTTPS 和非默认端口，仍需包含端口；
5. 修改后需要重启 Kaiwu Pod，让 DSH 使用新的 trusted host；
6. 浏览器入口 Service 用固定 `nodePort`（`kaiwu-console-public: 30088`、`kaiwu-boss-public: 30089`），避免 origin 漂移导致 `entry_url` 与 trusted host 不一致。

### 9.2 NetworkPolicy

如果集群启用默认拒绝 NetworkPolicy，需要允许 ANI Gateway Pod 访问 `kaiwu/kaiwu-console Pod:3080` 和 `kaiwu/kaiwu-boss Pod:3080`，同时确保 Gateway 能访问 Kubernetes API 和 DNS。

### 9.3 Gateway 环境变量

建议：

```text
KAIWU_NAMESPACE=kaiwu

KAIWU_CONSOLE_SERVICE=kaiwu-console
KAIWU_CONSOLE_SECRET=kaiwu-console-web-token
KAIWU_CONSOLE_PUBLIC_URL=http://10.10.1.66:30088

KAIWU_BOSS_SERVICE=kaiwu-boss
KAIWU_BOSS_SECRET=kaiwu-boss-web-token
KAIWU_BOSS_PUBLIC_URL=http://10.10.1.66:30089

KAIWU_ENTRY_TOKEN_TTL=5m
```

代理 Cookie 相关的 `KAIWU_PROXY_SIGNING_SECRET` / `KAIWU_PROXY_BOOTSTRAP_TTL` / `KAIWU_PROXY_SESSION_TTL` 已随子路径代理删除，现有环境可从网关 Deployment 中移除（留存的 `ani-kaiwu-runtime` Secret 也不再被读取）。

## 10. 测试计划

### 10.1 单元测试

Console 必须覆盖：

1. `tenant-a` 且 tenant ID 一致返回 200；
2. 其他租户返回 403；
3. tenant ID 与 `tenant-a` 不一致返回 403；
4. tenant frozen/disabled 返回 403；
5. Service 缺失返回 503；
6. Secret 缺失返回 503。

BOSS 必须覆盖：

1. `platform-admin` 且 active 返回 200；
2. `platform-ops` 返回 403；
3. `platform-readonly` 返回 403；
4. disabled 用户返回 403；
5. Service/Secret 缺失返回 503。

响应安全必须覆盖：

1. JSON 中不出现 `webToken`；
2. JSON 中不出现 `cluster_ip`；
3. `entry_url` 不包含 token；
4. Cookie 是 HttpOnly；
5. Console Cookie 不能访问 BOSS 路径；
6. BOSS Cookie 不能访问 Console 路径。

最新代码额外必须覆盖：

1. API Key 调用 Console/BOSS entry 返回 403 `KAIWU_USER_CREDENTIAL_REQUIRED`；
2. Console 请求 `scope != "tenant"` 返回 403；
3. BOSS 请求 `scope != "platform"` 返回 403；
4. `PlatformUserAdminStore` 未注入时返回 503，而不是 panic；
5. Kaiwu runtime reader 返回空 ClusterIP、空 webToken、非 3080 端口时返回 503；
6. Handler 不直接 import Kubernetes REST client 或 `pkg/adapters/runtime`；
7. `make validate-services` 中 Services route contract 不产生新 exception。

### 10.2 Proxy 测试（已删除）

> 子路径代理测试（`kaiwu_proxy_test.go`、`kaiwu_cookies_test.go`）已随代理代码删除。当前开物测试覆盖 `kaiwu_resources_test.go`（鉴权矩阵、独占 origin 契约、未配置 origin 时 503、运行时校验）与 `kaiwu_runtime_test.go`（环境变量解析与非法值拒绝）。

使用 fake DSH server 验证：

1. token 交换；
2. DSH `Set-Cookie` 转发；
3. Cookie 前缀隔离；
4. `/kaiwu/console/assets/x` 转发为 `/assets/x`；
5. `Location` 重写；
6. WebSocket upgrade；
7. SSE flush；
8. 无 Gateway Cookie 返回 401；
9. Cookie 过期返回 401；
10. Cookie client 不匹配返回 403。
11. 同一个 DSH `webToken` 可以为两个不同 ANI 用户分别完成首次 token 交换；
12. DSH 返回 `303 Location: /` 时，浏览器收到的是 `/kaiwu/console` 或 `/kaiwu/boss`；
13. DSH 返回 `Path=/` Cookie 时，浏览器收到 `Path=/kaiwu/console` 或 `Path=/kaiwu/boss`；
14. 只有 DSH Cookie、没有 ANI Gateway 签名 Cookie 时返回 401；
15. 只有 ANI Gateway 签名 Cookie、没有 DSH Cookie 时会重新完成 token exchange，而不是直接公开 DSH 页面；
16. token exchange 过程的 URL query 不进入访问日志和错误日志。
17. 连续两次 token exchange 都会各自调用 K8s Secret GET，第二次能读到更新后的 `webToken`，验证 Gateway 没有缓存旧值；
18. Kaiwu Pod 重启并更新 Secret 后，Gateway 使用新 `webToken` 完成 token exchange，而不是复用重启前的旧值。

### 10.3 联调验证

Console 允许场景：

```text
tenant-a active 用户
    -> GET /api/v1/svc/integrations/kaiwu/console/entry
    -> 200，entry_url=/kaiwu/console
    -> 浏览器打开 /kaiwu/console
    -> DSH 页面 200
```

Console 拒绝场景：非 `tenant-a` 用户调用 entry 返回 403。

BOSS 允许场景：

```text
platform-admin active 用户
    -> GET /api/v1/svc/integrations/kaiwu/boss/entry
    -> 200，entry_url=/kaiwu/boss
    -> 浏览器打开 /kaiwu/boss
    -> DSH 页面 200
```

BOSS 拒绝场景：`platform-ops`、`platform-readonly` 或 tenant 用户调用 entry 返回 403。

安全验证：

```text
entry 响应中无 webToken
entry 响应中无 ClusterIP
浏览器 URL 中无 webToken
Gateway 日志中无 webToken
Kaiwu Service 仍为 ClusterIP
外部无法直接访问 Kaiwu NodePort
```

## 11. 验收命令

### 11.1 确认 Kaiwu Service

```bash
kubectl -n kaiwu get svc kaiwu-console kaiwu-boss \
  -o custom-columns=NAME:.metadata.name,TYPE:.spec.type,CLUSTER_IP:.spec.clusterIP,PORT:.spec.ports[0].port
```

预期 `TYPE=ClusterIP` 且 `PORT=3080`。

### 11.2 确认 Secret 存在

```bash
kubectl -n kaiwu get secret kaiwu-console-web-token kaiwu-boss-web-token
```

不要在输出中展示 Secret 内容。

### 11.3 从 Gateway Pod 内部连通性

```bash
kubectl -n ani-system exec deploy/ani-gateway -- \
  sh -c 'wget -qO- --timeout=3 http://kaiwu-console.kaiwu.svc.cluster.local:3080/ >/dev/null && echo console-ok'

kubectl -n ani-system exec deploy/ani-gateway -- \
  sh -c 'wget -qO- --timeout=3 http://kaiwu-boss.kaiwu.svc.cluster.local:3080/ >/dev/null && echo boss-ok'
```

如果使用 ClusterIP 直连，也可以用 Service 中解析到的 ClusterIP 测试，但不要把 IP 输出到前端。

## 12. 实施顺序

1. 在 `api/openapi/services/v1.yaml` 中新增两个 entry API；
2. 在 `kaiwu_resources.go` 中定义本地 `KaiwuRuntimeReader` 接口，并通过 `RegisterOptions` 注入；
3. 在 `pkg/adapters/runtime/kubernetes_kaiwu.go` 中实现 K8s Service/Secret 只读 adapter；
4. 实现 `kaiwu_resources.go`，只通过注入 reader 获取内部目标和 webToken；
5. 实现开物独占 origin 入口：配置 `KAIWU_*_PUBLIC_URL` 后返回绝对入口地址（代理 Cookie 与子路径反代已删除）；
7. 增加单元测试、adapter 测试和 proxy 测试；
8. 运行 `make validate-services`，确认 Services route contract、SDK/API docs 和架构门禁通过；
9. 修改 Kaiwu 两个 Deployment 的 `KAIWU_TRUSTED_HOST`（各自只写自身 origin 的 authority）；
10. 重启 Kaiwu Console/BOSS；
11. 构建、推送并部署 ANI Gateway；
12. 按 Console/BOSS 正反向用例联调；
13. 验证响应和日志不泄露 token 与 ClusterIP。

## 13. 安全清单

- [ ] Console 只允许 `tenant-a`；
- [ ] Console tenant ID 与环境变量和数据库一致；
- [ ] BOSS 只允许 active `platform-admin`；
- [ ] 不信任前端传入身份；
- [ ] API 契约维护在 `api/openapi/services/v1.yaml`；
- [ ] Gateway handler 不直接读取 Kubernetes API；
- [ ] API Key 不能调用 Kaiwu entry；
- [ ] entry API 不返回 DSH token；
- [ ] Gateway 不缓存 DSH token，每次 token exchange 前实时读取 K8s Secret；
- [ ] entry API 不返回 ClusterIP；
- [ ] entry API 不下发任何 ANI 代理 Cookie；
- [ ] `entry_url` 是开物 origin 的绝对地址，含一次性启动 token；
- [ ] `KAIWU_*_PUBLIC_URL` 未配置时 entry 返回 503，不降级为公开入口；
- [ ] Gateway 日志不记录 DSH token 与 `entry_url`；
- [ ] Kaiwu Service 保持 ClusterIP；
- [ ] 浏览器入口只经 `kaiwu-console-public` / `kaiwu-boss-public`，并按内网 ACL 或 Ingress 鉴权限制；
- [ ] 每个开物实例的 `KAIWU_TRUSTED_HOST` 只包含自身 origin 的 authority，不含 ANI 网关 authority。
