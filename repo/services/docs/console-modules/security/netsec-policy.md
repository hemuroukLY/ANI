# 网络安全策略

## 页面定位

租户级 **网络安全策略**（微隔离、出口控制、服务网格策略）配置页，与 `compute/network/security-group.md` 安全组 **分资源**。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-netsec-policy-draft.md`。

## 文档管理规则

- 本文是 **网络安全策略** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-netsec-policy-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/v1.yaml`（Core Networks）
- PRD/SPEC：`tasks/modules/prd/console/security/prd-console-netsec-policy.md`、`tasks/modules/spec/console/security/spec-console-netsec-policy.md`
- TASK：`TASK-CORE-015` 子项 / Phase 3 §4

## Core 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/networks/policies` | `listNetworkSecurityPolicies` | 3a |
| POST | `/api/v1/networks/policies` | `createNetworkSecurityPolicy` | 3a |
| GET | `/api/v1/networks/policies/{policy_id}` | `getNetworkSecurityPolicy` | 3a |
| PUT | `.../{policy_id}` | `updateNetworkSecurityPolicy` | 3a |
| DELETE | `.../{policy_id}` | `deleteNetworkSecurityPolicy` | 3b |

Schema（草案）：`NetworkSecurityPolicy`、`NetworkSecurityPolicyRule`、`CreateNetworkSecurityPolicyRequest`。

RBAC（草案）：`scope:networks:read|write|delete`（与 SG 同族）。

## Services 层要求

无 Services 路径。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 读/写权限 | networks scope | `403` |
| POST 幂等键 | 必填 | `400` |
| scope=vpc | 提供有效 `vpc_id` | `422`（建议语义：`VPC_NOT_FOUND`；**待 YAML 声明后定稿**） |
| rules 合法 | CIDR/端口格式 | `422`（建议语义：`INVALID_POLICY_RULE`；**待 YAML 声明后定稿**） |

## 页面职责

- 策略列表/详情/创建/编辑（policy_type、rules）
- 与 SG 分 Tab 或分导航；UI 说明叠加生效顺序
- 跳转 **VPC**、**安全组**

## 操作可用性矩阵

| 操作 | 只读用户 | 网络管理员 | 说明 |
|---|---|---|---|
| 列表/详情 | 可用 | 可用 | GET |
| 创建 | 不可用 | 可用 | POST |
| 编辑/禁用 | 不可用 | 可用 | PUT |
| 删除 | 不可用 | 可用 | DELETE（3b） |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-netsec-policy-draft.md` §4。

### `POST /api/v1/networks/policies`（规划）

- 成功：`201 + NetworkSecurityPolicy`
- 错误：`400`、`422`

## 待补边界

- CNI / ServiceMesh provider 映射 — 实现说明
- 与安全组合并 vs 独立 — **已选独立** `/networks/policies`
- PUT rules 全量替换 vs patch — 评审固定

## 相关模块

- `compute/network/security-group.md`
- `compute/network-management.md`

## 验收标准

- [ ] 与 security-groups path 不混用
- [ ] POST 必填 idempotency_key
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
