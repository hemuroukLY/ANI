# PRD: Console 租户管理

## 模块目标

帮助租户管理员管理成员、查看角色、配置 SSO 与 Webhook。

## 范围

- 成员列表、邀请、移除
- 角色列表（只读）
- SSO 配置读写
- Webhook 列表与创建

## 非目标

- API Key 管理（见 `api-key-management.md`，Core Auth）
- 平台级租户运营（BOSS）
- 角色权限编辑（待补）

## 验收标准

- [ ] 路径使用 `/api/v1/svc/tenant/*`
- [ ] 写操作含 `idempotency_key`
- [ ] 前置条件与操作矩阵与主维护源一致
- [ ] 不与 deprecated `/sandboxes*`、`/gpu-containers*` 混写
