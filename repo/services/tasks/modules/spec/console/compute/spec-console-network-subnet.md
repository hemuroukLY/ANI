# SPEC: Console 网络子网

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/networks/subnets` | `listNetworkSubnets` | `200 + NetworkSubnetListResponse` | `401`,`403` |
| POST | `/api/v1/networks/subnets` | `createNetworkSubnet` | `201 + NetworkSubnet` | `400`,`401`,`403`,`404` |
| GET | `/api/v1/networks/subnets/{subnet_id}` | `getNetworkSubnet` | `200 + NetworkSubnet` | `401`,`403`,`404` |
| DELETE | `/api/v1/networks/subnets/{subnet_id}` | `deleteNetworkSubnet` | `200 + NetworkSubnet` | `401`,`403`,`404` |

## 3. TODO-YAML

- 子网 IP 分配 list

## 4. References

- `docs/console-modules/compute/network/subnet.md`
