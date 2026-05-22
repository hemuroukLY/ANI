# M1-ENCRYPT-A · Encryption Core API dev profile

- 完成日期：2026-05-21；2026-05-22 补齐 seal/unseal-token local profile；2026-05-23 由 M1-ENCRYPT-B 补齐 rotate/revoke local profile
- 状态：✅

实现了 `/api/v1/encryption/keys` 的 create/get/list/delete 路由、local service、租户隔离与幂等创建；2026-05-22 继续补齐 `POST /api/v1/encryption/seal` 和 `POST /api/v1/encryption/unseal-token` 的 local dev profile，返回模拟 sealed object URI、unseal token、过期时间和 `dev_profile` 标记；并补充了 OpenAPI 合同、SDK 生成和网关单元测试。

## 真实边界

本批次仍是 Core dev/local profile 切片，不包含：

- 真实 KMS/SM4 provider 集成。
- 真实 KMS/SM4 provider 下的密钥生命周期操作或 Secret 绑定。
