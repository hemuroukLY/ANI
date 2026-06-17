# PRD: Console 块存储卷快照

> Revised: 2026-06-17  
> 详文：`docs/console-modules/compute/storage/block-storage-snapshot.md`

## 1. Overview

卷快照 list + create（202）；DELETE 未声明。

## 2. Goals

- idempotency_key on POST
- 202 异步语义

## 3. User Stories

US-001 列表；US-002 创建快照。

## 4. References

- 详文、SPEC、block-storage.md
