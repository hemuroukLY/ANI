# SPEC: Console 文件存储

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-file-storage.md`
> Generated: 2026-06-14 | Scope: `Console / 算力与云资源 / 存储管理 / 文件存储`

## 1. Summary

- Source: `tasks/modules/prd/console/compute/prd-console-file-storage.md`
- User Stories covered: `US-001` ~ `US-004`
- Functional Requirements covered: `FR-1` ~ `FR-4`

## 2. Architecture

- 页面属于 `Console / 算力与云资源 / 存储管理`
- 页面直接消费 `Core /api/v1/filesystems*`
- 页面不定义新的 `Services /api/v1/svc/*` 文件存储资源契约
- VM、容器、GPU 容器等页面中的共享文件摘要和跳转只消费本模块已冻结定义

## 3. Data Model

- `StorageFilesystem`
- `StorageFilesystemListResponse`
- `CreateStorageFilesystemRequest`
- `StorageResourceState`

## 4. API Design

| Method | Path | Description | Response |
|---|---|---|---|
| GET | `/api/v1/filesystems` | 查询文件存储列表 | `200 + StorageFilesystemListResponse` |
| POST | `/api/v1/filesystems` | 创建文件存储 | `201 + StorageFilesystem` |
| GET | `/api/v1/filesystems/{filesystem_id}` | 查询文件存储 | `200 + StorageFilesystem` |
| DELETE | `/api/v1/filesystems/{filesystem_id}` | 删除文件存储 | `200 + StorageFilesystem` |

## 5. Business Logic

- 文件存储页默认展示当前租户可见文件存储资源
- 总览区负责说明资源边界与能力状态
- 创建区只承接现有 `CreateStorageFilesystemRequest`
- 对未冻结能力，只记录为“待补 Core 契约”，不能写成正式接口表

## 6. Error Handling

```json
{
  "code": "UPPER_SNAKE",
  "message": "error message",
  "request_id": "req-xxx"
}
```

## 7. Security

- 全部文件存储接口依赖标准认证
- 后端必须从认证上下文获取租户边界
- 前端不可信任也不显式传 `tenant_id / X-Tenant-Id`

## 8. Testing Strategy

- 校验文件存储列表、详情、创建、删除路径与返回码对齐 `Core v1.yaml`
- 校验创建请求带有 `idempotency_key`
- 校验待补能力不会误显示为已冻结接口
