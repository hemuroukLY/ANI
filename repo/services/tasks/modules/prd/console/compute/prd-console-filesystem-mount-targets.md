# PRD: Console filesystem-mount-targets

详文：`docs/console-modules/compute/storage/filesystem-mount-targets.md` · Phase 2 YAML 已声明项文档收口（2026-06-17）。

## 目标

- 为 **Core / Filesystems / 挂载目标** 子页提供可维护的产品边界说明
- 明确 `GET .../mount-targets` 已冻结、mount/unmount 写操作未声明

## 用户故事（规划）

- 作为存储管理员，我希望在文件系统详情查看挂载目标列表，以便了解关联实例
- 作为开发者，我希望从文档得知本页只读、不写 mount/unmount API

## 范围

- 挂载目标只读 Tab、跳转关联实例 — 见主维护源
- 不修改 ANI-main OpenAPI（由 Core 团队按 TASK-CORE-010 推进）

## 非目标

- 不自造 mount/unmount 写接口
- 活跃挂载实时拓扑、NFS/POSIX 策略细化
