# PRD: Console k8s-workloads

详文：`docs/console-modules/compute/k8s-workloads.md` · Phase 2 YAML 已声明项文档收口（2026-06-17）。

## 目标

- 为 **Core / K8s 集群 / 工作负载摘要** 子页提供可维护的产品边界说明
- 明确 workloads 列表只读、非 Helm/应用市场完整能力

## 用户故事（规划）

- 作为集群管理员，我希望在集群详情查看工作负载摘要，以便快速了解负载概况
- 作为开发者，我希望从文档得知 Pod 事件/Helm 未在 YAML 冻结

## 范围

- 工作负载摘要只读页 — 见主维护源
- 不修改 ANI-main OpenAPI（由 Core 团队按 TASK-CORE-012 推进）

## 非目标

- 不自造工作负载写 API
- Pod 事件、Helm、应用市场
