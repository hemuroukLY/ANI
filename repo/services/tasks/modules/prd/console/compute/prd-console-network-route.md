# PRD: Console 网络路由

> Revised: 2026-06-17  
> 详文：`docs/console-modules/compute/network/route.md`

## 1. Overview

路由 list + create；无 GET/DELETE by route_id。

## 2. Goals

- vpc_id 筛选
- create 409 冲突

## 3. User Stories

US-001 按 VPC 列表；US-002 创建路由。

## 4. Non-Goals

- 单条 GET/DELETE（TODO-YAML）

## 5. References

- 详文、SPEC
