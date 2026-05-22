# M1-SECRETS-A · Secret Core API dev profile

- 完成日期：2026-05-22
- 状态：✅

实现了 `/api/v1/secrets` 的 create/get/list/delete 路由，以及 `POST /api/v1/secrets/{secret_id}/bindings` 的绑定记录；补齐 OpenAPI 合同、ports 接口、local runtime service、租户隔离、幂等创建、SDK 生成和 router 单元测试。

## 真实边界

本批次仍是 Core dev/local profile 切片，不包含：

- 真实 Kubernetes Secret 写入。
- 真实实例环境变量或文件挂载注入。
- Secret 值读取 API；响应只返回元数据和 key 名称，不返回明文。
- 后台 controller 驱动的绑定生效与状态回写。
