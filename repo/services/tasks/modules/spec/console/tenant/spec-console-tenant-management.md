# SPEC: Console 租户管理

> 权威源：`ANI-main/repo/api/openapi/services/v1.yaml` → `/tenant/*`  
> 主维护源：`docs/console-modules/tenant/tenant-management.md`

## Frozen Paths

| Method | Path | operationId | Success |
|---|---|---|---|
| GET | `/tenant/members` | listTenantMembers | 200 |
| POST | `/tenant/members` | inviteTenantMember | 202 |
| DELETE | `/tenant/members/{member_id}` | removeTenantMember | 204 |
| GET | `/tenant/roles` | listTenantRoles | 200 |
| GET | `/tenant/sso` | getTenantSSOConfig | 200 |
| PUT | `/tenant/sso` | updateTenantSSOConfig | 200 |
| GET | `/tenant/webhooks` | listTenantWebhooks | 200 |
| POST | `/tenant/webhooks` | createTenantWebhook | 201 |

完整 URL 前缀：`https://{host}/api/v1/svc`

## 创建前置条件

| 场景 | 依赖项 | 未满足时的 HTTP 响应 |
|---|---|---|
| 邀请成员 | 操作者为租户管理员 | `403 FORBIDDEN` |
| 邀请成员 | 邮箱未被占用 | `409 CONFLICT` |
| 更新 SSO | OIDC Provider 配置合法 | `400 BAD_REQUEST` |
| 创建 Webhook | URL 可达且 secret 已配置 | 产品建议 `422 PRECONDITION_FAILED`（**当前 `createWebhook` YAML 未声明 `422`**；建议语义：前置校验失败） |

## 操作可用性矩阵

见 `docs/console-modules/tenant/tenant-management.md`
