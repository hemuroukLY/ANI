# Console 辅助维护文档：存储模块

## 1. 适用范围

- `docs/console-modules/compute/storage-management.md`
- `docs/console-modules/compute/storage/block-storage.md`
- `docs/console-modules/compute/storage/object-storage.md`
- `docs/console-modules/compute/storage/file-storage.md`
- `docs/console-modules/compute/storage/vector-storage.md`

## 2. 日志规范

| 模块 | 必记事件 | 必记字段 | 脱敏要求 |
|---|---|---|---|
| 存储管理 | `storage_overview_loaded`、`storage_category_opened` | `module`、`category`、`request_id?`、`result` | 不记录无关下游资源细节 |
| 块存储 | `volume_list_loaded`、`volume_created`、`volume_deleted` | `module`、`action`、`request_id`、`volume_id?`、`result` | 不记录开发态调试字段完整正文 |
| 对象存储 | `object_list_loaded`、`object_metadata_created`、`object_deleted` | `module`、`action`、`request_id`、`object_id?`、`bucket`、`result` | 不记录对象内容、上传凭据、未脱敏 key 明细 |
| 文件存储 | `filesystem_list_loaded`、`filesystem_created`、`filesystem_deleted` | `module`、`action`、`request_id`、`filesystem_id?`、`result` | 不记录挂载凭据和未冻结挂载明细 |
| 向量存储 | `vector_store_list_loaded`、`vector_store_created`、`vector_search_submitted`、`vector_store_deleted` | `module`、`action`、`request_id`、`vector_store_id?`、`result` | 不记录原始查询向量、未脱敏命中正文 |

## 3. 异常处理方案

| 模块 | 异常场景 | 前端处理 | 运维检查 |
|---|---|---|---|
| 存储管理 | 子资源区局部失败 | 仅该资源区报错，其他资源区继续可用 | 检查对应子资源路径与总入口映射 |
| 块存储 | 删除失败或状态阻塞 | 保留当前卷上下文，提示是否存在引用或状态限制 | 检查卷状态和引用关系 |
| 对象存储 | 元数据创建失败或详情不存在 | 提示 `request_id` 与对象标识，不把失败文案写成上传失败 | 检查对象元数据路径和字段映射 |
| 文件存储 | endpoint 缺失或删除失败 | 以说明态展示 endpoint，删除失败时保留详情上下文 | 检查文件系统状态与挂载边界是否被误用 |
| 向量存储 | 搜索失败或删除失败 | 仅搜索区报错，不阻断列表与详情；删除失败保留结果反馈 | 检查向量搜索请求、资源状态和依赖关系 |

## 4. 运维操作指南

| 模块 | 日常检查 | 常见处置 | 变更注意 |
|---|---|---|---|
| 存储管理 | 检查四类存储入口、统计与跳转是否一致 | 异常优先定位到子资源页，不在总入口里扩写细节 | 新增能力先下沉到子模块，再回填总入口 |
| 块存储 | 检查卷列表、详情、创建、删除四条主链路 | 删除失败时先核对引用关系，不提前引入 `409` 以外口径 | 快照、挂载需等待正式契约 |
| 对象存储 | 检查元数据列表、详情、创建、删除是否一致 | 任何“上传失败”类反馈都需先确认是否误用了元数据页 | 桶、桶策略、上传下载不能直接塞回本页 |
| 文件存储 | 检查列表、详情、创建、删除与 endpoint 显示 | endpoint 缺失时保持说明态，不误导用户已可挂载 | 挂载目标与访问策略待正式契约后再补 |
| 向量存储 | 检查列表、详情、搜索、删除链路是否独立可用 | 搜索异常优先检查请求和结果映射，不暗示存在写入修复 | 写入、导入、索引维护必须等待正式契约 |

## 5. 版本变更记录

| 模块 | 版本 | 日期 | 变更内容 |
|---|---|---|---|
| 存储管理 | `v1.2` | `2026-06-15` | 清洗 PRD/SPEC 的假设式与开发式残留，统一到现有 Core 冻结路径与待补边界口径 |
| 块存储 | `v1.2` | `2026-06-16` | 清洗 PRD/SPEC 的假设式与开发式残留，统一到 `volumes*` 的正式冻结边界 |
| 对象存储 | `v1.2` | `2026-06-15` | 清洗 PRD/SPEC 的假设式与开发式残留，持续强调当前仅承接对象元数据能力 |
| 文件存储 | `v1.2` | `2026-06-16` | 复核主文档与 HTML 摘要，确认 `filesystems*` 冻结边界与挂载目标待补口径一致 |
| 向量存储 | `v1.2` | `2026-06-16` | 复核主文档与 HTML 摘要，确认 `vector-stores*` 与 `search` 冻结边界一致 |
