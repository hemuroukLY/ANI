# 租户管理员权限变更方案：排除移交所有者

> **状态：** 待评审
> **创建日期：** 2026-08-21
> **关联：** [migration 003](../../repo/deploy/migrations/20260502_003_permissions_schema.sql) / [migration 20260821_001](../../repo/deploy/migrations/20260821_001_tenant_admin_invitation.sql) / [SPEC §5.1.4](../../repo/services/tasks/modules/spec/boss/tenant/spec-new-boss-tenant-admin.md)
> **范围：** 仅数据库层（迁移文件 + roles seed），不含 CheckPermission / RBAC 中间件改动

---

## 1. 背景与目标

### 1.1 问题

当前 tenant-admin / tenant-owner 的 permissions seed：

```json
[{"resource":"*","actions":["*"],"scope":"tenant"}]
```

`actions:["*"]` 通配符自动覆盖所有 action，**无法从数据层排除任何操作**。SPEC §5.1.4 要求移交所有者（transfer-ownership）仅 platform-admin / platform-ops 可发起，tenant-admin 不可执行，但当前 seed 无法实现这一排除。

### 1.2 目标

- **tenant-owner**：列出全部 20 个 resource，actions 保持 `*`（全权限，不变）
- **tenant-admin**：列出全部 20 个 resource，其中 19 个 actions=`*`；`tenants` resource 按业务能力拆分 actions，排除 `transfer-ownership`
- 未来新增自定义 action，tenant-admin 的 `tenants` 行默认无权限，人工审核后写入

### 1.3 action 设计原则

| 分类 | action 命名 | 说明 |
|------|------------|------|
| 普通 CRUD | create / read / list / update / delete | 标准资源操作，不拆分 |
| 状态机转换 | disable / enable / freeze / unfreeze | 业务状态变更，独立 action |
| 安全敏感操作 | reset-password / change-role / rotate-key | 安全相关，独立 action |
| 所有权和授权变化 | transfer-ownership / assign / revoke | 权限流转，独立 action |
| 业务动作 | invite / resend-invitation | 业务语义，独立 action |

**"业务动作语义相同的接口共享一个 action"**：多个端点表达同一种业务语义时映射到同一个 action，不重复。例：`POST /instances/{id}/start` 和 `POST /instances/{id}/power-on` 共享 `start`。反例：`invite`（邀请）和 `transfer-ownership`（移交）语义不同，必须用不同 action。

**普通 CRUD 不拆分**：非业务能力操作（create/read/list/update/delete）保持标准动词，不为每个端点单独定义 action。

### 1.4 排除原理

`transfer-ownership` 是所有权变更类业务 action，不是标准 CRUD。tenant-admin 的 `tenants` resource 只列标准 CRUD + 其他业务 action，`transfer-ownership` 不在其中 → 默认拒绝。

**代码依赖**：当前 `inferPermission`（[rbac.go:82-106](../../repo/services/ani-gateway/internal/middleware/rbac.go)）把所有 POST 映射成 `action=create`，无法区分业务 action。需对业务路径返回对应 action。**不在本次 DB 变更范围，但需同步实施。**

### 1.5 不含范围

本方案**仅修改数据库层**（迁移文件 + roles seed），不涉及：

- `auth-service` 的 `CheckPermission` 代码
- `ani-gateway` 的 `inferPermission` / RBAC 中间件
- pilot-v2 的 `permissionAllows` / `requiredBoundary` evaluator

---

## 2. permissions 变更

### 2.1 有效 resource 清单（20 个）

来自 [migration 003 第 13-16 行](../../repo/deploy/migrations/20260502_003_permissions_schema.sql)：

```
instances | networks | volumes | filesystems | objects |
vector-stores | k8s-clusters | baremetal-hosts | gpu-inventory |
dpu-inventory | registry | encryption | secrets | metering |
observability | notifications | audit | users | tenants | api-keys
```

### 2.2 tenants resource 的 action 拆分

`tenants` 是业务能力资源，按操作语义拆分 action（非标准 CRUD）：

| 分类 | action | 对应端点 | tenant-admin |
|------|--------|----------|:---:|
| 普通 CRUD | read | GET .../admins/{userId}、GET .../role、GET .../changeable-roles、GET .../audit-logs | ✅ |
| 普通 CRUD | list | GET /tenant-admins | ✅ |
| 普通 CRUD | delete | DELETE .../admins/{userId} | ✅ |
| 状态机转换 | disable | POST .../disable | ✅ |
| 状态机转换 | enable | POST .../enable | ✅ |
| 安全敏感 | reset-password | POST .../reset-password | ✅ |
| 安全敏感 | change-role | PUT .../role | ✅ |
| 业务动作 | invite | POST .../invite | ✅ |
| 业务动作 | resend-invitation | POST .../invitation/resend | ✅ |
| 所有权变更 | **transfer-ownership** | POST .../transfer-ownership | **❌ 排除** |

> 注：`create`/`update` 不列入 tenants —— 邀请用 `invite`（不是 create），改角色用 `change-role`（不是 update）。

### 2.3 tenant-owner permissions（全权限，不变）

列出全部 20 个 resource，actions 均为 `*`：

```json
[
  {"resource":"instances",      "actions":["*"],"scope":"tenant"},
  {"resource":"networks",       "actions":["*"],"scope":"tenant"},
  {"resource":"volumes",        "actions":["*"],"scope":"tenant"},
  {"resource":"filesystems",    "actions":["*"],"scope":"tenant"},
  {"resource":"objects",        "actions":["*"],"scope":"tenant"},
  {"resource":"vector-stores",  "actions":["*"],"scope":"tenant"},
  {"resource":"k8s-clusters",   "actions":["*"],"scope":"tenant"},
  {"resource":"baremetal-hosts","actions":["*"],"scope":"tenant"},
  {"resource":"gpu-inventory",  "actions":["*"],"scope":"tenant"},
  {"resource":"dpu-inventory",  "actions":["*"],"scope":"tenant"},
  {"resource":"registry",       "actions":["*"],"scope":"tenant"},
  {"resource":"encryption",     "actions":["*"],"scope":"tenant"},
  {"resource":"secrets",        "actions":["*"],"scope":"tenant"},
  {"resource":"metering",       "actions":["*"],"scope":"tenant"},
  {"resource":"observability", "actions":["*"],"scope":"tenant"},
  {"resource":"notifications", "actions":["*"],"scope":"tenant"},
  {"resource":"audit",          "actions":["*"],"scope":"tenant"},
  {"resource":"users",          "actions":["*"],"scope":"tenant"},
  {"resource":"tenants",        "actions":["*"],"scope":"tenant"},
  {"resource":"api-keys",       "actions":["*"],"scope":"tenant"}
]
```

> tenant-owner 的 `tenants` actions=`*`，包含 `transfer-ownership`。tenant-owner 是租户所有者，可发起移交。

### 2.4 tenant-admin permissions（tenants 拆分，排除 transfer-ownership）

列出全部 20 个 resource。其中 19 个 actions=`*`（普通 CRUD 资源不拆分），`tenants` 按业务能力拆分：

```json
[
  {"resource":"instances",      "actions":["*"],"scope":"tenant"},
  {"resource":"networks",       "actions":["*"],"scope":"tenant"},
  {"resource":"volumes",        "actions":["*"],"scope":"tenant"},
  {"resource":"filesystems",    "actions":["*"],"scope":"tenant"},
  {"resource":"objects",        "actions":["*"],"scope":"tenant"},
  {"resource":"vector-stores",  "actions":["*"],"scope":"tenant"},
  {"resource":"k8s-clusters",   "actions":["*"],"scope":"tenant"},
  {"resource":"baremetal-hosts","actions":["*"],"scope":"tenant"},
  {"resource":"gpu-inventory",  "actions":["*"],"scope":"tenant"},
  {"resource":"dpu-inventory",  "actions":["*"],"scope":"tenant"},
  {"resource":"registry",       "actions":["*"],"scope":"tenant"},
  {"resource":"encryption",     "actions":["*"],"scope":"tenant"},
  {"resource":"secrets",        "actions":["*"],"scope":"tenant"},
  {"resource":"metering",       "actions":["*"],"scope":"tenant"},
  {"resource":"observability", "actions":["*"],"scope":"tenant"},
  {"resource":"notifications", "actions":["*"],"scope":"tenant"},
  {"resource":"audit",          "actions":["*"],"scope":"tenant"},
  {"resource":"users",          "actions":["*"],"scope":"tenant"},
  {"resource":"tenants",        "actions":["read","list","delete","disable","enable","reset-password","change-role","invite","resend-invitation"],"scope":"tenant"},
  {"resource":"api-keys",       "actions":["*"],"scope":"tenant"}
]
```

### 2.5 其他角色不变

| 角色 | permissions | 说明 |
|------|-------------|------|
| platform-admin | `[{"resource":"*","actions":["*"],"scope":"platform"}]` | 超级管理员，`*` 通配覆盖 transfer-ownership |
| user | 显式枚举（5 个 resource，scope=own/tenant） | 已是显式，无需改动 |
| auditor | `[{"resource":"*","actions":["read","list"],"scope":"tenant"}]` | 只读，无需改动 |

---

## 3. 迁移文件变更

### 3.1 变更策略

当前项目处于开发阶段，迁移文件尚未上生产，直接修改 seed。若已上生产，改用增量迁移（见 §3.4）。

### 3.2 修改 migration 003（tenant-admin seed）

文件：[20260502_003_permissions_schema.sql](../../repo/deploy/migrations/20260502_003_permissions_schema.sql) 第 48-52 行。

**修改前**：

```sql
    (
        '00000000-0000-0000-0000-000000000002',
        NULL, 'tenant-admin',
        '[{"resource":"*","actions":["*"],"scope":"tenant"}]'
    ),
```

**修改后**：

```sql
    (
        '00000000-0000-0000-0000-000000000002',
        NULL, 'tenant-admin',
        '[
            {"resource":"instances",      "actions":["*"],"scope":"tenant"},
            {"resource":"networks",       "actions":["*"],"scope":"tenant"},
            {"resource":"volumes",        "actions":["*"],"scope":"tenant"},
            {"resource":"filesystems",    "actions":["*"],"scope":"tenant"},
            {"resource":"objects",        "actions":["*"],"scope":"tenant"},
            {"resource":"vector-stores",  "actions":["*"],"scope":"tenant"},
            {"resource":"k8s-clusters",   "actions":["*"],"scope":"tenant"},
            {"resource":"baremetal-hosts","actions":["*"],"scope":"tenant"},
            {"resource":"gpu-inventory",  "actions":["*"],"scope":"tenant"},
            {"resource":"dpu-inventory",  "actions":["*"],"scope":"tenant"},
            {"resource":"registry",       "actions":["*"],"scope":"tenant"},
            {"resource":"encryption",     "actions":["*"],"scope":"tenant"},
            {"resource":"secrets",        "actions":["*"],"scope":"tenant"},
            {"resource":"metering",       "actions":["*"],"scope":"tenant"},
            {"resource":"observability", "actions":["*"],"scope":"tenant"},
            {"resource":"notifications", "actions":["*"],"scope":"tenant"},
            {"resource":"audit",          "actions":["*"],"scope":"tenant"},
            {"resource":"users",          "actions":["*"],"scope":"tenant"},
            {"resource":"tenants",        "actions":["read","list","delete","disable","enable","reset-password","change-role","invite","resend-invitation"],"scope":"tenant"},
            {"resource":"api-keys",       "actions":["*"],"scope":"tenant"}
        ]'
    ),
```

### 3.3 修改 migration 20260821_001（tenant-owner seed）

文件：[20260821_001_tenant_admin_invitation.sql](../../repo/deploy/migrations/20260821_001_tenant_admin_invitation.sql) 第 81-86 行。

**修改前**：

```sql
INSERT INTO roles (id, tenant_id, name, permissions) VALUES
    ('00000000-0000-0000-0000-000000000005',
     NULL, 'tenant-owner',
     '[{"resource":"*","actions":["*"],"scope":"tenant"}]')
ON CONFLICT (tenant_id, name) DO UPDATE
    SET permissions = EXCLUDED.permissions;
```

**修改后**：

```sql
INSERT INTO roles (id, tenant_id, name, permissions) VALUES
    ('00000000-0000-0000-0000-000000000005',
     NULL, 'tenant-owner',
     '[
         {"resource":"instances",      "actions":["*"],"scope":"tenant"},
         {"resource":"networks",       "actions":["*"],"scope":"tenant"},
         {"resource":"volumes",        "actions":["*"],"scope":"tenant"},
         {"resource":"filesystems",    "actions":["*"],"scope":"tenant"},
         {"resource":"objects",        "actions":["*"],"scope":"tenant"},
         {"resource":"vector-stores",  "actions":["*"],"scope":"tenant"},
         {"resource":"k8s-clusters",   "actions":["*"],"scope":"tenant"},
         {"resource":"baremetal-hosts","actions":["*"],"scope":"tenant"},
         {"resource":"gpu-inventory",  "actions":["*"],"scope":"tenant"},
         {"resource":"dpu-inventory",  "actions":["*"],"scope":"tenant"},
         {"resource":"registry",       "actions":["*"],"scope":"tenant"},
         {"resource":"encryption",     "actions":["*"],"scope":"tenant"},
         {"resource":"secrets",        "actions":["*"],"scope":"tenant"},
         {"resource":"metering",       "actions":["*"],"scope":"tenant"},
         {"resource":"observability", "actions":["*"],"scope":"tenant"},
         {"resource":"notifications", "actions":["*"],"scope":"tenant"},
         {"resource":"audit",          "actions":["*"],"scope":"tenant"},
         {"resource":"users",          "actions":["*"],"scope":"tenant"},
         {"resource":"tenants",        "actions":["*"],"scope":"tenant"},
         {"resource":"api-keys",       "actions":["*"],"scope":"tenant"}
     ]')
ON CONFLICT (tenant_id, name) DO UPDATE
    SET permissions = EXCLUDED.permissions;
```

### 3.4 增量迁移方案（若已上生产）

新建 `20260821_002_tenant_admin_permissions_exclude_transfer.sql`：

```sql
-- ANI Platform · Migration 20260821_002
-- Description: tenant-owner / tenant-admin 权限从通配 */*/tenant 改为显式 resource 枚举；
--              tenant-admin 的 tenants resource 按业务能力拆分 action，排除 transfer-ownership
-- Depends on: 20260502_003_permissions_schema.sql, 20260821_001_tenant_admin_invitation.sql
-- Rationale:
--   1. 显式列出全部 20 个 resource，不再用 resource:"*" 通配
--   2. tenant-owner：全部 resource actions=*（全权限，不变）
--   3. tenant-admin：19 个 resource actions=*（普通 CRUD 不拆分）；
--      tenants 按业务能力拆分为 read/list/delete/disable/enable/reset-password/change-role/invite/resend-invitation，
--      排除 transfer-ownership（所有权变更，仅 platform-admin/ops 可发起）

BEGIN;

-- tenant-owner：显式列出全部 resource，actions=*
UPDATE roles SET permissions = '[
    {"resource":"instances",      "actions":["*"],"scope":"tenant"},
    {"resource":"networks",       "actions":["*"],"scope":"tenant"},
    {"resource":"volumes",        "actions":["*"],"scope":"tenant"},
    {"resource":"filesystems",    "actions":["*"],"scope":"tenant"},
    {"resource":"objects",        "actions":["*"],"scope":"tenant"},
    {"resource":"vector-stores",  "actions":["*"],"scope":"tenant"},
    {"resource":"k8s-clusters",   "actions":["*"],"scope":"tenant"},
    {"resource":"baremetal-hosts","actions":["*"],"scope":"tenant"},
    {"resource":"gpu-inventory",  "actions":["*"],"scope":"tenant"},
    {"resource":"dpu-inventory",  "actions":["*"],"scope":"tenant"},
    {"resource":"registry",       "actions":["*"],"scope":"tenant"},
    {"resource":"encryption",     "actions":["*"],"scope":"tenant"},
    {"resource":"secrets",        "actions":["*"],"scope":"tenant"},
    {"resource":"metering",       "actions":["*"],"scope":"tenant"},
    {"resource":"observability", "actions":["*"],"scope":"tenant"},
    {"resource":"notifications", "actions":["*"],"scope":"tenant"},
    {"resource":"audit",          "actions":["*"],"scope":"tenant"},
    {"resource":"users",          "actions":["*"],"scope":"tenant"},
    {"resource":"tenants",        "actions":["*"],"scope":"tenant"},
    {"resource":"api-keys",       "actions":["*"],"scope":"tenant"}
]' WHERE id = '00000000-0000-0000-0000-000000000005' AND name = 'tenant-owner';

-- tenant-admin：19 个 resource actions=*；tenants 按业务能力拆分，排除 transfer-ownership
UPDATE roles SET permissions = '[
    {"resource":"instances",      "actions":["*"],"scope":"tenant"},
    {"resource":"networks",       "actions":["*"],"scope":"tenant"},
    {"resource":"volumes",        "actions":["*"],"scope":"tenant"},
    {"resource":"filesystems",    "actions":["*"],"scope":"tenant"},
    {"resource":"objects",        "actions":["*"],"scope":"tenant"},
    {"resource":"vector-stores",  "actions":["*"],"scope":"tenant"},
    {"resource":"k8s-clusters",   "actions":["*"],"scope":"tenant"},
    {"resource":"baremetal-hosts","actions":["*"],"scope":"tenant"},
    {"resource":"gpu-inventory",  "actions":["*"],"scope":"tenant"},
    {"resource":"dpu-inventory",  "actions":["*"],"scope":"tenant"},
    {"resource":"registry",       "actions":["*"],"scope":"tenant"},
    {"resource":"encryption",     "actions":["*"],"scope":"tenant"},
    {"resource":"secrets",        "actions":["*"],"scope":"tenant"},
    {"resource":"metering",       "actions":["*"],"scope":"tenant"},
    {"resource":"observability", "actions":["*"],"scope":"tenant"},
    {"resource":"notifications", "actions":["*"],"scope":"tenant"},
    {"resource":"audit",          "actions":["*"],"scope":"tenant"},
    {"resource":"users",          "actions":["*"],"scope":"tenant"},
    {"resource":"tenants",        "actions":["read","list","delete","disable","enable","reset-password","change-role","invite","resend-invitation"],"scope":"tenant"},
    {"resource":"api-keys",       "actions":["*"],"scope":"tenant"}
]' WHERE id = '00000000-0000-0000-0000-000000000002' AND name = 'tenant-admin';

COMMIT;

-- ===========================================================================
-- Rollback
-- ===========================================================================
-- UPDATE roles SET permissions = '[{"resource":"*","actions":["*"],"scope":"tenant"}]'
--   WHERE id = '00000000-0000-0000-0000-000000000002' AND name = 'tenant-admin';
-- UPDATE roles SET permissions = '[{"resource":"*","actions":["*"],"scope":"tenant"}]'
--   WHERE id = '00000000-0000-0000-0000-000000000005' AND name = 'tenant-owner';
```

---

## 4. CHECK 约束兼容性

[migration 003 第 22-39 行](../../repo/deploy/migrations/20260502_003_permissions_schema.sql) 的 `roles_permissions_schema` CHECK 约束要求：

- `permissions` 必须是数组 ✅
- 每项必须是 object ✅
- 含 `resource`(string) ✅
- 含 `actions`(非空数组) ✅
- 可选 `scope` ∈ {`tenant`, `own`, `platform`} ✅

本方案的 JSON 完全符合约束，无需修改 CHECK 约束。

---

## 5. 排除效果验证

### 5.1 移交所有者（transfer-ownership）—— tenant-admin 排除

- 端点：`POST /api/v1/svc/tenants/{tenantId}/transfer-ownership`
- `inferPermission`（更新后）→ `resource=tenants, action=transfer-ownership`
- tenant-admin 的 `tenants` 条目：`{"actions":["read","list","delete","disable","enable","reset-password","change-role","invite","resend-invitation"]}`
- `actionMatch`：`"transfer-ownership" ∈ [...]` → **false**
- → **拒绝** ✅

### 5.2 移交所有者 —— tenant-owner 允许

- tenant-owner 的 `tenants` 条目：`{"actions":["*"]}`
- `actionMatch`：`"transfer-ownership" ∈ ["*"]` → **true**
- → **允许** ✅

### 5.3 tenant-admin 其余操作 —— 保留

| 操作 | 端点 | 推断 action | tenants actions 含该 action | 结果 |
|------|------|------------|---------------------------|------|
| 邀请管理员 | POST .../invite | invite | ✅ | 允许 |
| 重发邀请 | POST .../invitation/resend | resend-invitation | ✅ | 允许 |
| 跨租户列表 | GET /tenant-admins | read | ✅ | 允许 |
| 管理员详情 | GET .../admins/{userId} | read | ✅ | 允许 |
| 软删除管理员 | DELETE .../admins/{userId} | delete | ✅ | 允许 |
| 查角色权限 | GET .../role | read | ✅ | 允许 |
| 修改角色 | PUT .../role | change-role | ✅ | 允许 |
| 可变更角色 | GET .../changeable-roles | read | ✅ | 允许 |
| 重置密码 | POST .../reset-password | reset-password | ✅ | 允许 |
| 禁用管理员 | POST .../disable | disable | ✅ | 允许 |
| 启用管理员 | POST .../enable | enable | ✅ | 允许 |
| 操作历史 | GET .../audit-logs | read | ✅ | 允许 |

### 5.4 tenant-admin 其余 19 个 resource —— 全权限

instances / networks / volumes / filesystems / objects / vector-stores / k8s-clusters / baremetal-hosts / gpu-inventory / dpu-inventory / registry / encryption / secrets / metering / observability / notifications / audit / users / api-keys 均为 `actions:["*"]`，普通 CRUD 资源不拆分，全权限保留。

---

## 6. 代码依赖（不在本次范围，但需同步实施）

### 6.1 inferPermission 识别业务 action

[rbac.go:82-106](../../repo/services/ani-gateway/internal/middleware/rbac.go) 需对业务路径返回对应 action 而非 `create`/`update`：

```go
func inferPermission(method, path string) (string, string) {
    // ... 现有 resource 推断逻辑 ...

    // 业务 action 识别（按路径段映射）
    switch {
    case strings.Contains(path, "/transfer-ownership"):
        return resource, "transfer-ownership"
    case strings.Contains(path, "/reset-password"):
        return resource, "reset-password"
    case strings.Contains(path, "/disable"):
        return resource, "disable"
    case strings.Contains(path, "/enable"):
        return resource, "enable"
    case strings.HasSuffix(path, "/role") && method == http.MethodPut:
        return resource, "change-role"
    case strings.Contains(path, "/invite"):
        return resource, "invite"
    case strings.Contains(path, "/invitation/resend"):
        return resource, "resend-invitation"
    }

    // 标准 CRUD 映射（不变）
    switch method {
    case http.MethodPost:
        return resource, "create"
    // ...
    }
}
```

### 6.2 CheckPermission 改为读取 permissions JSONB

[auth_service.go:201-227](../../repo/services/auth-service/internal/service/auth_service.go) 当前是角色名直通，不读 permissions JSONB。需改为查询 DB 加载 permissions 并判断 `action ∈ actions`，本方案的 seed 数据才能生效。

### 6.3 GET → read 命名统一

`inferPermission` 把 GET 映射为 `get`，但 seed 用 `read`。需在代码层统一，否则 GET 请求会因 `get ∉ actions`（含 `read` 不含 `get`）被拒绝。

---

## 7. 未来扩展

### 7.1 新增标准 action（非 tenants resource）

例如给 instances 新增 `snapshot`：tenant-admin 的 instances 是 `actions:["*"]`，自动覆盖，无需改 seed。

### 7.2 新增需排除的业务 action（tenants resource）

例如新增 `merge-tenant`（合并租户），要排除 tenant-admin：

1. `inferPermission` 对 `/merge-tenant` 路径返回 `action=merge-tenant`
2. tenant-admin 的 tenants actions 数组不加 `merge-tenant` → 自动排除
3. tenant-owner 的 `*` 和 platform-admin 的 `*` 覆盖

### 7.3 新增 resource

需在 tenant-owner / tenant-admin 的 permissions JSON 中追加对应 resource 条目，否则新 resource 默认无权限（因为不再用 `resource:"*"` 通配）。

### 7.4 新增业务端点（tenants resource）

例如新增 `POST .../freeze`（冻结租户）：

1. `inferPermission` 对 `/freeze` 路径返回 `action=freeze`
2. 评估 tenant-admin 是否需要 → 若需要，在 tenants actions 数组追加 `"freeze"`
3. 若不需要排除，不加即默认拒绝

---

## 8. 注意事项

1. **本方案仅改数据库 seed**：`CheckPermission` 代码当前仍为角色名直通，不读 permissions JSONB。要让 seed 数据生效，需后续代码改动（见 §6）。

2. **tenant-owner 全权限**：tenant-owner 的 `tenants` actions=`*`，包含 `transfer-ownership`。tenant-owner 是租户所有者，可发起移交。

3. **tenant-admin 仅排除 transfer-ownership**：tenant-admin 的 `tenants` resource 按业务能力拆分为 9 个 action（read/list/delete/disable/enable/reset-password/change-role/invite/resend-invitation），仅 transfer-ownership 被排除。其余 19 个 resource 全权限 `*`（普通 CRUD 不拆分）。

4. **显式 resource 枚举**：不再用 `resource:"*"` 通配，全部 20 个 resource 逐条列出。新增 resource 时需手动追加到 seed，否则默认无权限。

5. **platform-admin 不受影响**：保留 `*/*/platform`，超级管理员仍可执行所有操作。

6. **action 设计原则**：普通 CRUD 用标准动词（create/read/list/update/delete），业务能力按语义独立命名（invite/disable/enable/reset-password/change-role/transfer-ownership 等）。语义相同的端点共享一个 action，不为路径别名单独定义。
