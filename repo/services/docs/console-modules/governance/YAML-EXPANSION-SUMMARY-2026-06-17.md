# YAML EXPANSION SUMMARY

生成时间：2026-06-17  
前置：Phase 1 完成（`GAP-REPORT-2026-06-17.md` + `TODO-YAML` 标记）  
脚本：`scripts/phase2-yaml-expansion.py`

## 汇总

| 文件 | Phase 2 前 ops | Phase 2 后 ops | 新增 |
|---|---|---|---|
| Core `v1.yaml` | 92 | 111 | **+19** |
| Services `services/v1.yaml`（active） | 29 | 41 | **+12** |
| Services deprecated | 19 | 19 | 0（未复活） |

## 新增到 Core v1.yaml（19 条）

| 路径 | HTTP | operationId | 对应文档功能 | 复杂度 |
|---|---|---|---|---|
| `/gpu-inventory` | GET | listGPUInventory | GPU 设备清单 | 中 |
| `/gpu-inventory/occupancy` | GET | getGPUOccupancy | GPU 占用分布 | 低 |
| `/sandbox-templates` | GET | listSandboxTemplates | Sandbox 模板列表 | 低 |
| `/instances/{instance_id}/logs` | GET | listInstanceLogs | 容器/GPU 日志 | 中 |
| `/instances/{instance_id}/events` | GET | listInstanceEvents | 容器/GPU 事件 | 中 |
| `/instances/{instance_id}/metrics` | GET | getInstanceMetrics | 容器/GPU 监控 | 中 |
| `/instances/{instance_id}/exec` | POST | createInstanceExecSession | 终端/exec | 高 |
| `/instances/{instance_id}/security-events` | GET | listInstanceSecurityEvents | Sandbox 安全事件 | 中 |
| `/networks/routes` | GET | listNetworkRoutes | 路由表 | 中 |
| `/networks/routes` | POST | createNetworkRoute | 创建路由 | 中 |
| `/volumes/{volume_id}/snapshots` | GET | listVolumeSnapshots | 卷快照列表 | 中 |
| `/volumes/{volume_id}/snapshots` | POST | createVolumeSnapshot | 创建卷快照 | 中 |
| `/buckets` | GET | listStorageBuckets | 对象存储桶 | 中 |
| `/buckets` | POST | createStorageBucket | 创建桶 | 中 |
| `/objects/upload` | POST | uploadStorageObject | 对象上传 | 高 |
| `/objects/{object_id}/download` | GET | downloadStorageObject | 对象下载 | 中 |
| `/filesystems/{filesystem_id}/mount-targets` | GET | listFilesystemMountTargets | 挂载目标 | 中 |
| `/vector-stores/{vector_store_id}/documents` | POST | insertVectorStoreDocuments | 向量写入 | 高 |
| `/k8s-clusters/{cluster_id}/workloads` | GET | listK8sClusterWorkloads | K8s 工作负载摘要 | 中 |

## Core 422 扩充（非新路径）

| 路径 | 变更 |
|---|---|
| `POST /k8s-clusters` | 增 `422 PreconditionFailed` + 400/401/403/409 |
| `POST /vector-stores/{vector_store_id}/search` | 增 `422 PreconditionFailed` |

## 新增到 services/v1.yaml（12 条 active）

| 路径 | HTTP | operationId | 对应文档功能 | 复杂度 |
|---|---|---|---|---|
| `/inference-services/{service_id}/logs` | GET | getInferenceServiceLogs | 推理日志 | 中 |
| `/inference-services/{service_id}/test` | POST | testInferenceService | 推理调用测试 | 中 |
| `/inference-services/{service_id}/policies` | PUT | updateInferenceServicePolicies | 限流/访问策略 | 中 |
| `/knowledge-bases/{kb_id}/citations` | GET | listKnowledgeBaseCitations | 来源引用 | 低 |
| `/knowledge-bases/{kb_id}/sessions` | GET | listKnowledgeBaseSessions | 对话历史 | 中 |
| `/knowledge-bases/{kb_id}/permissions` | PUT | updateKnowledgeBasePermissions | 知识库权限 | 中 |
| `/tenant/members/{member_id}` | GET | getTenantMember | 成员详情 | 低 |
| `/tenant/roles/{role_id}` | PUT | updateTenantRole | 角色权限编辑 | 中 |
| `/tenant/webhooks/{webhook_id}/deliveries` | GET | listWebhookDeliveries | Webhook 投递日志 | 中 |
| `/integrations` | GET | listIntegrations | 第三方集成列表 | 中 |
| `/integrations` | POST | createIntegration | 创建第三方集成 | 中 |
| `/integrations/bots` | POST | createIntegrationBot | 企微/钉钉 Bot | 中 |

## Services schema/422 修补

| 路径 | 变更 |
|---|---|
| `DELETE /inference-services/{service_id}` | 补 `202 + AsyncTask` schema |
| `POST /tenant/webhooks` | 增 `422 PreconditionFailed` |

## 刻意未新增（说明）

| 项 | 原因 |
|---|---|
| `/gpu-containers*`、`/sandboxes*` | 已 deprecated；统一 Core `/instances*` |
| `POST /v1/chat/completions` | Gateway 层，不在 `services/v1.yaml` |
| Sandbox 延长/暂停 lifecycle action | 复用现有 `POST .../lifecycle`；待 Core 扩展 action 枚举，非新路径 |
| `GET /api/v1/alerts` | backlog P0-01；未在 Phase 1 TODO-YAML 主清单 |

## 文档标记

- `docs/console-modules/**` 中 Phase 1 `TODO-YAML` 已批量改为 `ADDED-TO-YAML`（15 个文件）
- `prototypes/ani-services-prototype-console.html` GPU 监控行已更新

## Phase 3 输入

- Core 新增 handler：**19** 个（+ 2 个 422 验通）→ **TASK-CORE-003~012** 等，见 `tasks/execution/CORE-TEAM-TASKS.md`
- Services 新增 handler：**12** 个（+ schema 修补）→ **TASK-SVC-013~016** 等，见 `tasks/execution/SERVICES-TEAM-TASKS.md`
- ~~更新 TASK 清单~~ ✅ 2026-06-17 Phase 3 对齐完成
- ~~P2「YAML 已有」模块详文~~ ✅ 2026-06-17：13 项（见 `console-undefined-features-backlog.md` P2「文档已补」表）

## 相关文件

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`
- `docs/console-modules/governance/schema-completion-tracker.md`
- `GAP-REPORT-2026-06-17.md`
