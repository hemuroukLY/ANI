# PRD: Console 网络子网

> Revised: 2026-06-17  
> 详文：`docs/console-modules/compute/network/subnet.md`

## 1. Overview

子网 CRUD；POST 必填 vpc_id、name、idempotency_key。

## 2. Goals

- 对齐 NetworkSubnet schema
- 创建时 vpc_id 404 处理

## 3. User Stories

US-001 列表/详情；US-002 创建；US-003 删除。

## 4. Functional Requirements

- FR-1: 不自造 IP 分配 list API

## 5. Non-Goals

- 子网 PATCH、IP 池

## 6. References

- 详文、SPEC
