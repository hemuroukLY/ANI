# PRD: Console 网络 VPC

> Revised: 2026-06-17  
> 详文：`docs/console-modules/compute/network/vpc.md`

## 1. Overview

Console 网络管理 VPC 子页：列表、详情、创建、删除。Core `/api/v1/networks/vpcs*`。

## 2. Goals

- CRUD 对齐 Core Networks 组
- POST 带 `idempotency_key`
- RBAC scope:networks:*

## 3. User Stories

### US-001: VPC 列表与详情

验收：list/get 200；404 单条不存在。

### US-002: 创建 VPC

验收：201；400/403；name + idempotency_key 必填。

### US-003: 删除 VPC

验收：200；404。

## 4. Functional Requirements

- FR-1: 不自造 VPC 扩展 schema
- FR-2: 接口冻结见详文

## 5. Non-Goals

- VPC PATCH、对等连接

## 6. References

- 详文、SPEC、network-management.md
