# PRD: Console 网络安全组

> Revised: 2026-06-17  
> 详文：`docs/console-modules/compute/network/security-group.md`

## 1. Overview

安全组 CRUD；rules 内嵌于 NetworkSecurityGroup，无独立 rules 子路径。

## 2. Goals

- 规则摘要展示
- POST idempotency_key

## 3. User Stories

US-001 列表详情；US-002 创建含 rules；US-003 删除。

## 4. Non-Goals

- rules 子资源 CRUD（TODO-YAML）

## 5. References

- 详文、SPEC
