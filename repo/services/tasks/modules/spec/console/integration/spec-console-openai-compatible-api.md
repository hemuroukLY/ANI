# SPEC: Console OpenAI 兼容 API

> Revised: 2026-06-17

## 1. Summary

Gateway 文档模块；无 repo OpenAPI 冻结。

## 2. Documented Paths（非 YAML 冻结）

| Method | Path | 说明 |
|---|---|---|
| POST | `/v1/chat/completions` | inferenceProxy stub |
| GET | `/v1/inference/stream` | 流式（对照 gateway 代码） |

## 3. Related Frozen（鉴权/资源）

- Core Auth Bearer/ApiKey
- Services InferenceService.endpoint_url

## 4. TODO-YAML

- Gateway OpenAPI 正式 YAML

## 5. References

- `docs/console-modules/integration/openai-compatible-api.md`
