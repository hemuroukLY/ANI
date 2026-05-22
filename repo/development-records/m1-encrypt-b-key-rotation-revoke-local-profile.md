# M1-ENCRYPT-B · Encryption key rotation/revoke local profile

- 完成日期：2026-05-23
- 状态：✅

在 `M1-ENCRYPT-A` 的密钥、seal 和 unseal-token local profile 基础上，补齐密钥轮转和吊销的 Core API local profile：

- `POST /api/v1/encryption/keys/{key_id}/rotate`
- `POST /api/v1/encryption/keys/{key_id}/revoke`

本批次同步更新了 OpenAPI 契约、ports 接口、local runtime service、Gateway 路由、四语言 SDK、静态 API 文档和 router 单元测试。local profile 中，rotate 会把旧 key 标记为 `rotated`，生成新的 `active` key；revoke 会把 key 标记为 `revoked`，并禁止 revoked key 继续 seal 或生成 unseal token。rotate/revoke 都要求 `idempotency_key`，并在 local service 中提供幂等回放。

## 真实边界

本批次仍是 Core dev/local profile 切片，不代表真实 KMS/SM4 provider 已接入。真实 provider 下的密钥生成、轮换、吊销、审计、密钥材料保护和真实加解密验证仍属于 Sprint 5 后续未完成项。

## 验证

```bash
go test ./services/ani-gateway/internal/router ./pkg/adapters/runtime -run TestEncryptionAPIDevProfileAndIdempotency -v
python scripts/validate_yaml.py api/openapi/v1.yaml api/openapi/services/v1.yaml api/core-beta-readiness.yaml api/core-v1-compatibility-baseline.yaml
python scripts/generate_api_docs.py
make gen-core-sdk
```
