# SPEC: Console 网络安全组

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/networks/security-groups` | `listNetworkSecurityGroups` | `200 + NetworkSecurityGroupListResponse` | `401`,`403` |
| POST | `/api/v1/networks/security-groups` | `createNetworkSecurityGroup` | `201 + NetworkSecurityGroup` | `400`,`401`,`403` |
| GET | `/api/v1/networks/security-groups/{security_group_id}` | `getNetworkSecurityGroup` | `200 + NetworkSecurityGroup` | `401`,`403`,`404` |
| DELETE | `/api/v1/networks/security-groups/{security_group_id}` | `deleteNetworkSecurityGroup` | `200 + NetworkSecurityGroup` | `401`,`403`,`404` |

## 3. TODO-YAML

- rules 子资源 CRUD；实例绑定

## 4. References

- `docs/console-modules/compute/network/security-group.md`
