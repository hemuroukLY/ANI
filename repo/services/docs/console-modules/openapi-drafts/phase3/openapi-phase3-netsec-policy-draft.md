# Phase 3 — 网络安全策略 OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`security/netsec-policy.md`  
> **TASK**：`TASK-CORE-015` 子项 / Phase 3 §4  
> **原则**：与 **安全组**（`/networks/security-groups`）分资源；面向微隔离/出口/网格策略。

---

## 1. 与安全组 / VPC 的边界

| 维度 | 安全组 SG | 网络安全策略（本草案） |
|---|---|---|
| 路径 | `/networks/security-groups` | `/networks/policies` |
| 模型 | ENI 级 L4 规则 | 租户/VPC 级 **策略**（egress、微隔离、mesh） |
| 已冻结 | ✅ YAML 已声明 | 规划 |
| 模块 | `security-group.md` | 本文 |

同一 VPC 可同时挂 SG 与 netsec policy；生效优先级 — Services/CNI 实现说明。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [Networks]
x-ani-rbac-scope:
  read:   scope:networks:read
  write:  scope:networks:write
  delete: scope:networks:delete
```

与 security-groups 共用 networks scope 族（**评审确认 create vs write 拆分**）。

---

## 3. Schemas（草案）

### NetworkSecurityPolicyRule

```yaml
NetworkSecurityPolicyRule:
  type: object
  required: [action, protocol]
  properties:
    action:     { type: string, enum: [allow, deny] }
    protocol:   { type: string, enum: [tcp, udp, icmp, any] }
    port_range: { type: string, nullable: true, description: 如 443 或 8000-9000 }
    cidr:       { type: string, nullable: true }
    host:       { type: string, nullable: true, description: 域名 egress 规则 }
    description: { type: string, nullable: true }
```

### NetworkSecurityPolicy

```yaml
NetworkSecurityPolicy:
  type: object
  required: [id, name, policy_type, status, created_at]
  properties:
    id:          { type: string, format: uuid }
    name:        { type: string }
    description: { type: string, nullable: true }
    policy_type: { type: string, enum: [egress_control, micro_segmentation, service_mesh] }
    scope:       { type: string, enum: [tenant, vpc], default: vpc }
    vpc_id:      { type: string, format: uuid, nullable: true }
    rules:
      type: array
      items: { $ref: '#/components/schemas/NetworkSecurityPolicyRule' }
    status:      { type: string, enum: [active, disabled, pending] }
    created_at:  { type: string, format: date-time }
    updated_at:  { type: string, format: date-time, nullable: true }
```

### CreateNetworkSecurityPolicyRequest

```yaml
CreateNetworkSecurityPolicyRequest:
  type: object
  required: [idempotency_key, name, policy_type]
  properties:
    idempotency_key: { type: string, format: uuid }
    name:            { type: string }
    description:     { type: string, nullable: true }
    policy_type:     { type: string, enum: [egress_control, micro_segmentation, service_mesh] }
    scope:           { type: string, enum: [tenant, vpc] }
    vpc_id:          { type: string, format: uuid, nullable: true }
    rules:           { type: array, items: { $ref: '#/components/schemas/NetworkSecurityPolicyRule' } }
```

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/networks/policies`

- operationId: `listNetworkSecurityPolicies`
- Query: `vpc_id`、`policy_type`、`status`、`limit`、`cursor`
- 成功：`200` + items + `next_cursor`

### 4.2 `POST /api/v1/networks/policies`

- operationId: `createNetworkSecurityPolicy`
- Body: `CreateNetworkSecurityPolicyRequest`
- 成功：`201 + NetworkSecurityPolicy`
- 错误：`400`、`401`、`403`、`422`（vpc 不存在 — `VPC_NOT_FOUND`；规则非法 — `INVALID_POLICY_RULE`）

### 4.3 `GET /api/v1/networks/policies/{policy_id}`

- operationId: `getNetworkSecurityPolicy`
- 成功：`200 + NetworkSecurityPolicy`
- 错误：`404`

### 4.4 `PUT /api/v1/networks/policies/{policy_id}`（3a）

- operationId: `updateNetworkSecurityPolicy`
- Body: `{ idempotency_key, name?, description?, rules?, status? }`
- 成功：`200 + NetworkSecurityPolicy`

### 4.5 `DELETE /api/v1/networks/policies/{policy_id}`（3b）

- operationId: `deleteNetworkSecurityPolicy`
- Body: `{ idempotency_key }`
- 成功：`204`

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 网络管理员 | 说明 |
|---|---|---|---|
| 列表/详情 | 可用 | 可用 | GET |
| 创建 | 不可用 | 可用 | POST |
| 编辑/禁用 | 不可用 | 可用 | PUT |
| 删除 | 不可用 | 可用 | DELETE（3b） |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/networks/policies?vpc_id=$VPC_ID"

curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","name":"egress-default","policy_type":"egress_control","vpc_id":"'"$VPC_ID"'","rules":[{"action":"allow","protocol":"tcp","port_range":"443"}]}' \
  "$BASE/api/v1/networks/policies"
```

---

## 7. 评审检查清单

- [ ] 与 `/networks/security-groups` 无资源混用
- [ ] POST 必填 `idempotency_key`
- [ ] scope=vpc 时 vpc_id 必填 — 422 口径
- [ ] 合入后更新 `netsec-policy.md`

---

## 相关文件

- `docs/console-modules/security/netsec-policy.md`
- `docs/console-modules/compute/network/security-group.md`
- `docs/console-modules/compute/network-management.md`
