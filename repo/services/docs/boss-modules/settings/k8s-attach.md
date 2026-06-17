# 已有 K8s 集群接入

## 页面定位

`已有 K8s 集群接入` 是 BOSS **交付与安装** 域下 **`ani-installer attach-k8s`** 专页：客户已有 Kubernetes 时，**只部署 ANI 平台层**，不修改 K8s 自身配置（ANI-07 §1 第四行、§4）。

**attach 执行走 CLI**，不是 BOSS REST。集群 **注册到控制面后**，BOSS 可 **只读参考** Core `GET/POST /api/v1/k8s-clusters`（`listK8sClusters` / `createK8sCluster`）展示已纳管集群 — **不得** 将 POST 描述为 attach 向导的替代。

权威源：ANI-07 §4。Console **无** 对等页。

## 文档管理规则

- 本文为主维护源
- PRD/SPEC：`prd-boss-k8s-attach.md` / `spec-boss-k8s-attach.md`
- OpenAPI 对照：`v1.yaml` `/k8s-clusters*`
- 枢纽：[`ani-installer.md`](ani-installer.md)

## Core 层要求

- attach 工作流 REST — **TODO-YAML: N/A**（CLI only）
- Core 已声明（**注册后只读参考**）：
  - `GET /api/v1/k8s-clusters` — `listK8sClusters` — `scope:k8s-clusters:read`
  - `POST /api/v1/k8s-clusters` — `createK8sCluster` — `scope:k8s-clusters:create` — 须 `idempotency_key`
- **禁止** BOSS 用 POST 冒充「一键 attach」
- JWT `tenant_id` 来自 claims；**不信任** body `tenant_id` 越权
- 422 仅 YAML 已声明 operation（如 create 422）

## 页面职责

- 文档化 §4.1 兼容性检查（K8s≥1.28、CNI、SC、Ingress）
- 文档化 §4.2 在已有 K8s 上的 6 步操作清单
- BOSS 展示 **已注册** 集群列表（`listK8sClusters` — handler 实现前标注）
- 深链 baremetal/vm（greenfield 对照）
- attach 完成后链 internal-ca、acceptance-manual

## 页面结构

```text
已有 K8s 接入
├── attach-k8s CLI 流程（§4）
├── 兼容性检查表（§4.1）
├── 平台组件 Helm 清单（§4.2）
├── 已注册集群表格（listK8sClusters · 只读）
└── 与 new 路径对照
```

## 数据来源与分层约束

| 层 | 来源 | 用法 |
|---|---|---|
| CLI | `ani-installer attach-k8s` | **attach 执行** |
| Core | `GET/POST /api/v1/k8s-clusters` | **注册后 list/create 参考** |
| ANI-07 §4 | 兼容性与步骤 | 权威 |

### 关键边界

- attach **不** 安装 K8s — 与 `new` 互斥
- `createK8sCluster` **≠** attach CLI — API 为控制面资源登记
- OpenAPI 已声明 ≠ handler 已实现

## 页面区块与数据来源映射

| 区块 | 来源 | 说明 |
|---|---|---|
| CLI 指引 | §4 | attach 步骤 |
| 兼容检查 | §4.1 | 阻止/警告 |
| 集群表格 | listK8sClusters | id/name/version/state |
| kubeconfig | getK8sClusterKubeconfig | **非 attach**；运维只读 P2 |

## BOSS 与 Console 分工

| 维度 | BOSS attach 页 | Console |
|---|---|---|
| attach 执行 | CLI | 无 |
| 集群 list | 平台 RBAC read | 租户 K8s 若有时为 Services |
| create API | 平台登记 | 非租户自助 attach |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| — | attach REST | **不存在** | CLI only |
| GET | `/api/v1/k8s-clusters` | `listK8sClusters` | YAML 已声明；BOSS 只读展示 |
| POST | `/api/v1/k8s-clusters` | `createK8sCluster` | YAML 已声明；**非 attach 向导** |
| GET | `/api/v1/k8s-clusters/{cluster_id}` | `getK8sCluster` | 详情 |
| DELETE | `/api/v1/k8s-clusters/{cluster_id}` | `deleteK8sCluster` | 删除登记 |

### K8sCluster schema（YAML 已声明）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | string | UUID |
| `tenant_id` | string | JWT 上下文 |
| `name` | string | 集群名 |
| `version` | string | K8s 版本 |
| `state` | enum | provisioning / running / deleting |
| `reason` | string | 状态原因 |
| `created_at` | date-time | — |
| `updated_at` | date-time | — |

### attach 兼容字段（CLI · §4.1）

| 字段 | 说明 |
|---|---|
| `k8s_version` | ≥ 1.28 |
| `cni_type` | KubeOVN 优选 |
| `storageclass_count` | ≥1 |
| `ingress_class` | 可选复用 |

## 字段级定义

### 查询字段（listK8sClusters）

| 字段 | 来源 | 必填 |
|---|---|---|
| `limit` | query | 可选 1–100 |
| `cursor` | query | 可选 |

### 写入字段（createK8sCluster · **非 attach**）

| 字段 | 必填 |
|---|---|
| `idempotency_key` | 是 |
| `name` | 是 |
| `version` | 可选 |

### 展示字段（UI）

| 字段 | 说明 |
|---|---|
| `state_badge` | provisioning 蓝 / running 绿 / deleting 灰 |
| `cni_warning` | 非 KubeOVN 时展示 §4.1 文案 |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| list 未实现 | 「k8s-clusters API 待 handler」 |
| state=provisioning | 行内 spinner |
| CNI 警告 | 黄条 §4.1 原文摘要 |
| attach 说明 |  prominent「请用 CLI attach-k8s」 |
| 403 | 无 read scope |

## 字段口径与单位

| 字段 | 格式 |
|---|---|
| `version` | semver K8s |
| `state` | YAML enum |
| `k8s_version` attach | ≥ 1.28 |

## 状态与能力口径

### K8sCluster.state

| 值 | UI |
|---|---|
| `provisioning` | 蓝色 |
| `running` | 绿色 |
| `deleting` | 灰色 |

### attach 兼容结论

| 值 | 含义 |
|---|---|
| PASS | 可 attach |
| WARN | 用户确认继续 |
| BLOCK | 不可 attach |

## 创建前置条件

| 依赖 | 失败 |
|---|---|
| 有效 kubeconfig | CLI BLOCK |
| K8s ≥ 1.28 | CLI BLOCK |
| BOSS list read | 403 |
| createK8sCluster | idempotency + create scope |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 平台管理员 |
|---|---|---|---|
| 查看 attach 文档 | ✅ | ✅ | ✅ |
| CLI attach-k8s | ❌ BOSS | ✅ | ✅ |
| listK8sClusters | ✅* | ✅* | ✅* |
| createK8sCluster | ❌ | ✅* | ✅* |
| deleteK8sCluster | ❌ | ✅* | ✅* |

\* handler 就绪后；OpenAPI 已声明 ≠ 已实现

## 删除前置校验与当前契约边界

`deleteK8sCluster`（YAML 已声明）预期：

| 校验 | 响应 |
|---|---|
| 集群存在 | 404 |
| delete scope | 403 |
| 集群上仍有 ANI workload | 409/422 待产品 |

attach CLI **无** HTTP DELETE。

## 接口冻结规则

#### `GET /api/v1/k8s-clusters` — `listK8sClusters`

- scope: `scope:k8s-clusters:read`
- 200: `K8sClusterListResponse`
- 401/403

#### `POST /api/v1/k8s-clusters` — `createK8sCluster`

- scope: `scope:k8s-clusters:create`
- body: `idempotency_key`, `name` required
- 201 / 400 / 409 / 422
- **非 attach 替代**

#### attach CLI — ANI-07 §4

- **TODO-YAML: N/A** for attach workflow

## 使用规则

- attach **必须** CLI
- **不得** 写 BOSS「POST = attach」
- greenfield 用 new — 见 baremetal/vm
- list 仅展示登记集群
- 不信任 body tenant_id

## 待补能力边界

- attach 进度 BOSS 展示 — P1 文档
- attach REST — **N/A**
- kubeconfig 下载 UI — P2

## 响应示例

### listK8sClusters 200（YAML 目标）

```json
{
  "items": [
    {
      "id": "kc-attach-001",
      "tenant_id": "platform",
      "name": "customer-existing-k8s",
      "version": "1.29.2",
      "state": "running",
      "reason": "",
      "created_at": "2026-06-15T10:00:00Z",
      "updated_at": "2026-06-15T10:00:00Z"
    }
  ],
  "total": 1,
  "next_cursor": null
}
```

### attach-k8s 成功（CLI · §4.2 摘要）

```text
attach-k8s complete.
Namespace ani-system created · CRDs applied · Helm releases: 5/5 ready
Console: https://192.168.1.50/ · Register cluster in control plane if required
```

## 错误示例

### 403 list

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-k8s-403-001"
}
```

### createK8sCluster 400

```json
{
  "code": "BAD_REQUEST",
  "message": "idempotency_key is required",
  "request_id": "req-boss-k8s-400-001"
}
```

### attach BLOCK — K8s 版本

```text
BLOCK: Kubernetes 1.27.3 < minimum 1.28 (ANI-07 §4.1)
Upgrade cluster before ani-installer attach-k8s
Exit code: 1
```

## 相关模块

- [`ani-installer.md`](ani-installer.md) · [`baremetal-deploy.md`](baremetal-deploy.md) · [`vm-deploy.md`](vm-deploy.md)
- [`internal-ca.md`](internal-ca.md) · [`acceptance-manual.md`](acceptance-manual.md)

## 回填验收标准

- [x] 满配 22 章
- [x] attach CLI vs k8s-clusters API 边界清晰
- [x] list/create YAML 引用 + handler  caveat
- [x] PRD/SPEC/HTML synced
