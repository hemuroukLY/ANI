# PRD: Console 推理服务调用测试

> Revised: 2026-06-17  
> 详文：`docs/console-modules/inference/inference-call-test.md`

## 1. Overview

`Console / 推理服务 / 调用测试` 在 Console 内对单个推理服务发送测试 prompt，验证 endpoint 可用。Services `POST .../test`。

## 2. 权威源

- `services/v1.yaml`：`testInferenceService`
- 非 Gateway `/v1/chat/completions`

## 3. Goals

- 表单提交 `prompt` + `idempotency_key`
- 展示 `output`、`latency_ms`
- 非 running 状态明确 422

## 4. User Stories

### US-001: 提交测试

验收：仅 `running` 可提交；返回 200 展示结果。

### US-002: 错误展示

验收：404/403/422 展示 `code` + `request_id`。

## 5. Functional Requirements

- FR-1: 路径与 `InferenceServiceTestRequest/Response` 一致
- FR-2: 写操作必填 `idempotency_key`
- FR-3: 422 口径见 module-delivery-workflow §2.10

## 6. Non-Goals

- 流式测试、测试历史持久化

## 7. References

- 详文、SPEC
