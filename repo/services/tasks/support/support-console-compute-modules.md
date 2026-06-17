# Console 辅助维护文档：算力、实例与网络模块

## 1. 适用范围

- `docs/console-modules/compute/gpu-management.md`
- `docs/console-modules/compute/vm-management.md`
- `docs/console-modules/compute/container-instance-management.md`
- `docs/console-modules/compute/gpu-container-instance-management.md`
- `docs/console-modules/compute/sandbox-instance-management.md`
- `docs/console-modules/compute/k8s-cluster-management.md`
- `docs/console-modules/compute/network-management.md`

## 2. 日志规范

| 模块 | 必记事件 | 必记字段 | 脱敏要求 |
|---|---|---|---|
| GPU 算力管理 | `gpu_overview_loaded`、`gpu_capability_pending_viewed`、`gpu_anomaly_drilldown_opened` | `module`、`action`、`request_id?`、`target_id?`、`result` | 不记录内部凭据和未脱敏节点敏感信息 |
| 云主机 VM | `vm_list_loaded`、`vm_created`、`vm_lifecycle_submitted`、`vm_console_requested` | `module`、`action`、`request_id`、`instance_id?`、`operation_id?`、`result` | 不记录 Console/VNC 长期凭据和 SSH 私钥 |
| 容器实例 | `container_list_loaded`、`container_created`、`container_lifecycle_submitted` | `module`、`action`、`request_id`、`instance_id?`、`operation_id?`、`result` | 不记录镜像私密拉取凭据 |
| GPU 容器实例 | `gpu_container_list_loaded`、`gpu_container_created`、`gpu_container_lifecycle_submitted` | `module`、`action`、`request_id`、`instance_id?`、`gpu_summary?`、`result` | 不记录底层调度密钥和私密镜像凭据 |
| Sandbox 实例 | `sandbox_created`、`sandbox_detail_loaded`、`sandbox_lifecycle_submitted` | `module`、`action`、`request_id`、`instance_id?`、`operation_id?`、`result` | 不记录会话密钥与敏感沙箱内容 |
| K8s 集群 | `cluster_list_loaded`、`cluster_created`、`cluster_upgraded`、`kubeconfig_requested` | `module`、`action`、`request_id`、`cluster_id?`、`task_id?`、`result` | 不记录 kubeconfig 明文、不记录代理凭证 |
| 网络管理 | `network_resource_loaded`、`vpc_created`、`subnet_created`、`security_group_created`、`load_balancer_created` | `module`、`action`、`request_id`、`resource_type`、`resource_id?`、`result` | 不记录无关租户网络明细和内部凭据 |

## 3. 异常处理方案

| 模块 | 异常场景 | 前端处理 | 运维检查 |
|---|---|---|---|
| GPU 算力管理 | 监控缺失或能力待补 | 基础资源继续展示，待补能力明确标注不可直接调用 | 检查监控回传、权威源冻结状态和占用对象关系 |
| 云主机 VM | 生命周期冲突、Console 会话失败 | 保留当前实例上下文，提示冲突或重新申请会话 | 检查实例状态、终端保护和操作历史 |
| 容器实例 | 发布摘要异常或动作冲突 | 列表和详情分区处理错误，保留操作历史入口 | 检查统一实例状态、操作记录和字段映射 |
| GPU 容器实例 | GPU 摘要缺失或动作失败 | 保留实例详情并标记 GPU 摘要暂不可用 | 检查统一实例与 GPU 资源摘要是否一致 |
| Sandbox 实例 | 创建失败、详情不存在、动作冲突 | 保留最近操作或创建上下文，不扩写独立列表 | 检查实例是否已创建、`sandbox` 摘要和操作记录 |
| K8s 集群 | 集群创建/升级/节点池任务失败 | 明确标记进行中/失败态，保留回看入口 | 检查长任务结果、集群状态和节点池变更记录 |
| 网络管理 | 删除冲突、详情不存在、局部资源区失败 | 仅对应资源区报错，删除前给出依赖说明 | 检查下游引用关系和资源状态回写 |

## 4. 运维操作指南

| 模块 | 日常检查 | 常见处置 | 变更注意 |
|---|---|---|---|
| GPU 算力管理 | 检查总览、节点、设备、占用分布定义是否一致 | 先定位占用对象，再判断是否需要跳转到关联实例模块处理 | 新增 GPU 观测或动作能力前先冻结正式契约 |
| 云主机 VM | 检查列表、详情、动作、操作历史是否串通 | 冲突优先查看操作历史，不直接修改页面状态文案 | 新增快照、卷、网络能力时回到对应资源域补契约 |
| 容器实例 | 检查列表、详情、发布摘要、操作历史联动 | 发布异常先看统一实例与操作记录，不创建新子页兜底 | logs/exec/dashboard 等能力需待 Core 扩充后再接入 |
| GPU 容器实例 | 检查 `kind=gpu_container` 列表与 GPU 摘要一致性 | GPU 摘要缺失时先排查返回字段，不伪造默认值 | 不得回退到旧 `/gpu-containers*` 路径 |
| Sandbox 实例 | 检查列表、创建、详情、动作、最近操作是否可回流 | 列表应来自统一实例 `kind=sandbox` 查询；不扩写独立 `/sandboxes*` 路径 | 不得把模板、安全事件等待补能力提前写进页面 |
| K8s 集群 | 检查集群、节点池、升级、kubeconfig 入口可用性 | 长任务失败时优先保留上下文与错误追踪标识 | 工作负载、Helm、事件等能力需单独冻结后再扩页 |
| 网络管理 | 检查四类资源列表与依赖提示是否一致 | 删除失败时先定位引用关系，不先修改错误码口径 | 路由、规则子资源、绑定解绑必须等待正式契约 |

## 5. 版本变更记录

| 模块 | 版本 | 日期 | 变更内容 |
|---|---|---|---|
| GPU 算力管理 | `v1.2` | `2026-06-15` | 按权威源移除伪造 GPU 独立接口，重整为页面定义与待补边界口径 |
| 云主机 VM | `v1.2` | `2026-06-15` | 清洗 PRD/SPEC 的假设式与开发式残留，统一到现有 Core 冻结路径口径 |
| 容器实例 | `v1.2` | `2026-06-16` | 清洗 PRD/SPEC 的假设式与开发式残留，统一到 `kind=container` 的正式冻结边界 |
| GPU 容器实例 | `v1.3` | `2026-06-16` | Core 统一实例为准；`/gpu-containers*` 在 services/v1.yaml 已 deprecated |
| Sandbox 实例 | `v1.4` | `2026-06-16` | Core 统一实例为准；`listInstances.kind` 枚举已含 sandbox；`/sandboxes*` 在 services/v1.yaml 已 deprecated |
| K8s 集群 | `v1.2` | `2026-06-16` | 清洗 PRD 的假设式残留，保持独立 `K8sClusters` 资源组口径 |
| 网络管理 | `v1.2` | `2026-06-16` | 清洗 PRD/SPEC 的假设式与开发式残留，统一到 `networks/*` 的正式冻结边界 |
