# SPEC: Console 向量存储写入

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| POST | `/api/v1/vector-stores/{vector_store_id}/documents` | `insertVectorStoreDocuments` | `202 + VectorStoreDocumentInsertResponse` | `400`,`401`,`403`,`404`,`422` |
| POST | `/api/v1/vector-stores/{vector_store_id}/search` | `searchVectorStore` | `200 + VectorStoreSearchResponse` | `400`,`401`,`403`,`404`,`422` |

## 3. TODO-YAML

- rebuild / bulk import

## 4. References

- `docs/console-modules/compute/storage/vector-store-write.md`
