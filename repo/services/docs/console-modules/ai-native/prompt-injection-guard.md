# AI 原生 — Prompt 注入防护

## 页面定位

配置 **Prompt 注入 / 越狱 / 敏感内容** 检测策略（输入/输出扫描、阻断/告警/仅日志）的 Console 页，面向 **Agent 会话**与**推理 Gateway** 统一防护。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md`。

## 文档管理规则

- 本文是 **Prompt 注入防护** 的主维护源
- **规划权威源**（合入前）：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 不得把草案路径写成**已实现**
- PRD/SPEC：`tasks/modules/prd/console/ai-native/prd-console-prompt-injection-guard.md`、`tasks/modules/spec/console/ai-native/spec-console-prompt-injection-guard.md`
- TASK：`TASK-SVC-018` 子项（Prompt Guard）

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/agent/guard-policies` | `listPromptGuardPolicies` | 3a |
| GET | `/api/v1/svc/agent/guard-policies/{policy_id}` | `getPromptGuardPolicy` | 3a |
| PUT | `/api/v1/svc/agent/guard-policies/{policy_id}` | `updatePromptGuardPolicy` | 3a |
| PUT | `/api/v1/svc/inference-services/{service_id}/guard-binding` | `updateInferenceServiceGuardBinding` | 3b |
| POST | `/api/v1/svc/agent/guard-policies/evaluate` | `evaluatePromptGuard` | 3b |

Schema（草案）：`PromptGuardPolicy`、`PromptGuardRule`、`PromptGuardPolicyUpdateRequest`、`PromptGuardEvaluateRequest/Response`。

RBAC（草案）：`scope:agent-guard:read`、`scope:agent-guard:write`。

<!-- TODO-YAML: 合入后删除「规划」标记 -->

## 与推理限流策略的分工

| 能力 | 模块 | 路径 |
|---|---|---|
| RPM、allowed_scopes | `inference-rate-limit-policy.md` | `PUT .../inference-services/{id}/policies` |
| 注入扫描、block/warn | **本文** | `/agent/guard-policies` |
| 推理服务绑定 guard | **本文（3b）** | `PUT .../guard-binding` |

**禁止**在 `InferenceServicePolicyUpdateRequest` 中增加 scan/action 字段。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 读权限 | `scope:agent-guard:read` | `403 FORBIDDEN` |
| 写权限 | `scope:agent-guard:write` | `403 FORBIDDEN` |
| PUT 幂等键 | `idempotency_key` 必填 | `400` |
| policy 存在 | 有效 `policy_id` | `404 NOT_FOUND` |
| custom 规则 | `pattern` 合法正则 | `422`（建议 `INVALID_GUARD_RULE`） |

默认策略（`is_default=true`）由系统 seed；**不可 DELETE**。

## 页面职责

- **策略列表/详情**：enabled、scan_target（input/output/both）、action（block/warn/log_only）
- **规则矩阵**：builtin 规则开关 + custom 正则（高风险须二次确认）
- **推理绑定**（3b）：服务详情独立 Tab，引用 `guard-binding`（不替代 RPM Tab）
- **策略演练**（3b）：`evaluate` 输入样例，展示 matched_rules（不写 audit）
- 跳转：`agent-audit.md`（policy_hit 时间线）

## 页面结构

```text
Prompt 注入防护
├── 租户默认策略（is_default）
│   ├── 扫描方向 / 处置动作
│   └── 规则列表（builtin + custom）
├── 推理服务绑定（3b · 子入口）
├── 策略演练（3b）
└── 关联审计（跳转 agent-audit）
```

## 操作可用性矩阵

| 操作 | 只读用户 | Agent 安全管理员 | YAML 合入后 |
|---|---|---|---|
| 查看策略 | 可用 | 可用 | GET |
| 编辑策略 | 不可用 | 可用 | PUT |
| 推理 guard 绑定 | 不可用 | 可用 | PUT guard-binding（3b） |
| 策略演练 | 不可用 | 可用 | POST evaluate（3b） |
| 编辑 RPM | 跳转 | 跳转 | inference-rate-limit-policy |

## 接口冻结规则

> **当前为规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md` §4。

### `GET /api/v1/svc/agent/guard-policies`（规划）

- 成功：`200 + PromptGuardPolicyListResponse`
- 错误：`401`、`403`

### `PUT /api/v1/svc/agent/guard-policies/{policy_id}`（规划）

- 成功：`200 + PromptGuardPolicy`
- 错误：`400`、`401`、`403`、`404`、`422`

### `POST /api/v1/svc/agent/guard-policies/evaluate`（规划 · 3b）

- 成功：`200 + PromptGuardEvaluateResponse`
- 不落库、默认不写 audit

## 待补边界

- 检测引擎实现（本地规则 vs 外部 API）— Services 实现说明，非 Console 契约
- 多命名策略 POST/create — 3b 可选；3a 仅 default + list
- 与 `ai-native-sandbox-security.md` 沙箱隔离策略 — 互补，不合并 API
- BOSS 全局规则模板下发 — Console 租户页只读/覆盖边界待产品定稿

## 相关模块

- `ai-native/agent-audit.md` — `policy_hit` 事件
- `ai-native/agent-session.md` — Agent 输入输出扫描挂载点
- `inference/inference-rate-limit-policy.md` — RPM/scope，非注入扫描
- `inference/inference-service.md` — guard-binding 宿主

## 验收标准

- [ ] 与 inference policies PUT 无字段混用
- [ ] default 策略不可删
- [ ] evaluate 不写 audit（默认）
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
