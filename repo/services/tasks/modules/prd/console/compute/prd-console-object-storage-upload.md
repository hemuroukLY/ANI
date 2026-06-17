# PRD: Console 对象存储上传下载

> Revised: 2026-06-17  
> 详文：`docs/console-modules/compute/storage/object-storage-upload.md`

## 1. Overview

桶 CRUD 子集 + upload（202 AsyncTask）+ download 预签名。

## 2. Goals

- 区分元数据 CRUD 与 upload
- upload 422 桶不存在

## 3. User Stories

US-001 桶列表创建；US-002 上传；US-003 下载链接。

## 4. References

- 详文、SPEC、module-delivery-workflow §2.9
