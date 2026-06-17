# SPEC: Console 网络路由

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/networks/routes` | `listNetworkRoutes` | `200 + NetworkRouteListResponse` | `401`,`403` |
| POST | `/api/v1/networks/routes` | `createNetworkRoute` | `201 + NetworkRoute` | `400`,`401`,`403`,`409` |

Query list: vpc_id, limit, cursor

## 3. TODO-YAML

- GET/DELETE `/networks/routes/{route_id}`

## 4. References

- `docs/console-modules/compute/network/route.md`
