# 品牌配置（白牌化）

## 页面定位

`品牌配置` 是 BOSS **交付与安装** 域下 **平台白牌化** 专页（ANI-07 §8 Phase 2）：动态配置 `platform_name`、Logo、主题色、ICP 等；Console/BOSS 首屏通过 **`GET /api/v1/branding`** 加载（`security: []`，无需认证）。

**BOSS 管理写操作** — **TODO-YAML P2**（PATCH/PUT 未在 v1.yaml 声明）。**不得** 伪造写 path 为已冻结。

一级权威：ANI-07 §8 + `v1.yaml` `/branding` GET。Console **只读消费** branding；**无** 租户自助改品牌。

## 文档管理规则

- 主维护源：本文
- PRD/SPEC：prd-boss-brand-config / spec-boss-brand-config
- OpenAPI 对照：`v1.yaml` `/branding`
- 流程：boss-full-depth-checklist

## Core 层要求

- **已声明**：`GET /api/v1/branding` — **无 operationId** — `security: []`
- **待补**：`PATCH /api/v1/branding` 或等价 — **TODO-YAML P2** — BOSS 写入口
- 存储：PostgreSQL `platform_branding` + MinIO `ani-branding`（§8.2）
- 写操作须 `idempotency_key`（合入后）
- 422 仅 YAML 已声明 operation
- **禁止** `/api/v1/boss/branding`

## 页面职责

- BOSS 表单编辑 §8.1 配置项（Logo 上传、色值、ICP）
- 预览 Console/BOSS 首屏效果（读 GET branding）
- 文档化 CSS 变量 `--color-primary` 机制（§8.2）
- 明确 Phase 2 — 非 v1.0.0 阻塞（§11）
- GET 字段与 UI 映射

## 页面结构

```text
品牌配置
├── 预览区（GET branding）
├── 编辑表单（PATCH 待 YAML P2）
├── Logo 上传（light/dark/favicon）
├── 色值与 ICP
└── Phase 2 / ANI-07 §8 说明
```

## 数据来源与分层约束

| 层 | 路径 | 用法 |
|---|---|---|
| Core GET | `/api/v1/branding` | **只读已声明** |
| Core PATCH | **TODO-YAML P2** | BOSS 写 |
| MinIO | ani-branding bucket | Logo URL |
| ANI-07 §8 | 字段语义 | 权威 |

### 关键边界

- GET **无** operationId — 正文写「未声明 operationId」
- OpenAPI 已声明 ≠ handler 已实现
- Console **只读** GET — 不在 Console 改品牌

## 页面区块与数据来源映射

| 区块 | 来源 |
|---|---|
| 预览 | GET branding |
| 平台名/色 | §8.1 |
| Logo | logo_*_url |
| 保存 | PATCH **P2** |

## BOSS 与 Console 分工

| 维度 | BOSS 品牌配置 | Console |
|---|---|---|
| 编辑 | BOSS 管理员 PATCH **P2** | 无 |
| 读取 | GET public | GET public |
| Logo 存储 | MinIO 上传 **P2** | 展示 URL |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/branding` | **未声明** | `security: []`；200 内联 schema |
| PATCH | `/api/v1/branding` | **不存在** | **TODO-YAML P2** |

### GET 响应字段（v1.yaml 已声明）

| 字段 | 类型 | 说明 |
|---|---|---|
| `platform_name` | string | 产品名 |
| `logo_light_url` | uri | 亮色 Logo |
| `logo_dark_url` | uri | 深色 Logo |
| `favicon_url` | uri | favicon |
| `primary_color` | string | 主色 HEX |
| `secondary_color` | string | 辅助色 |
| `icp_number` | string | ICP 备案号 |

### §8.1 扩展（产品 · PATCH 目标）

| 字段 | 说明 |
|---|---|
| `login_bg_image` | 登录背景 **P2** |

## 字段级定义

### 返回字段（GET · YAML）

| 字段 | 来源 | 说明 |
|---|---|---|
| `platform_name` | inline schema | 替换默认名 |
| `logo_light_url` | uri | PNG/SVG ≤2MB |
| `logo_dark_url` | uri | PNG/SVG ≤2MB |
| `favicon_url` | uri | ICO/PNG 32×32 |
| `primary_color` | string | 例 `#1677FF` |
| `secondary_color` | string | HEX |
| `icp_number` | string | 页脚 |

### 写入字段（PATCH · **TODO-YAML P2**）

| 字段 | 必填 | 说明 |
|---|---|---|
| `idempotency_key` | POST/PATCH | UUID |
| `platform_name` | 可选 | string |
| `primary_color` | 可选 | `#RRGGBB` |
| `secondary_color` | 可选 | HEX |
| `icp_number` | 可选 | string |
| `logo_light` | 可选 | multipart **待 YAML** |
| `logo_dark` | 可选 | multipart |
| `favicon` | 可选 | multipart |

### 展示字段（UI）

| 字段 | 说明 |
|---|---|
| `color_preview_swatch` | 色块 |
| `logo_preview` | img src=url |
| `default_fallback` | 未配置时用 ANI 默认 |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| GET 成功 | 预览区渲染 |
| GET 失败 | 默认品牌 + 黄条 |
| PATCH 未就绪 | 表单只读 +「P2 待 YAML」 |
| 403 写 | 无 BOSS 品牌写权限 |
| 非法 HEX | 前端拦截 **合入后 400** |
| Logo 超 2MB | 拒绝上传 |

## 字段口径与单位

| 字段 | 格式 |
|---|---|
| `primary_color` | `#` + 6 hex |
| Logo | PNG/SVG ≤ 2MB |
| favicon | 32×32 推荐 |
| URL | https uri（MinIO signed） |

## 状态与能力口径

| 能力 | 状态 |
|---|---|
| GET branding | YAML 已声明 |
| PATCH branding | **TODO-YAML P2** |
| Logo upload | **P2** |
| login_bg | **P2** §8.1 |

### 编辑态

| 态 | UI |
|---|---|
| readonly_p2 | 灰表单 |
| dirty | 未保存提示 |
| saved | 绿 toast **P2** |

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| GET | 无 auth |
| PATCH RBAC | 403 **P2** |
| idempotency_key | 400 **P2** |
| MinIO bucket | 5xx 运维 **P2** |

## 操作可用性矩阵

| 操作 | 只读 | 平台管理员 |
|---|---|---|
| GET 预览 | ✅ | ✅ |
| PATCH 保存 | ⏳ P2 | ⏳ P2 |
| Logo 上传 | ⏳ P2 | ⏳ P2 |
| 重置默认 | ⏳ P2 | ⏳ P2 |

## 删除前置校验与当前契约边界

品牌 **无** DELETE 资源；重置为默认 — **TODO-YAML P2**（PATCH 空值或 dedicated reset）。

**N/A** DELETE path。

## 接口冻结规则

#### `GET /api/v1/branding`

- **无** operationId
- `security: []`
- 200：内联 properties（见上表）
- **不得** 添加未声明字段为已冻结

#### `PATCH /api/v1/branding`（**TODO-YAML P2**）

<!-- TODO-YAML: PATCH /api/v1/branding + operationId + RBAC scope -->

- 须平台 branding 写 scope
- 须 `idempotency_key`
- Logo multipart 或 presigned upload **待设计**
- **合入前不得写已实现**

## 使用规则

- **不得** 写 PATCH 已实现
- Console **仅** GET — 不在租户路径写品牌
- GET 失败时前端 **fallback** 默认 — 不阻塞登录
- Phase 2 — 与 v1.0 交付范围区分（§11）
- 禁止 `/api/v1/boss/*`

## 待补能力边界

- PATCH + upload — **TODO-YAML P2**
- `login_bg_image` — P2
- RBAC scope — 待 Core
- 优先级：**P2**

## 响应示例

### GET /api/v1/branding 200（YAML 目标）

```json
{
  "platform_name": "Acme AI Platform",
  "logo_light_url": "https://minio.example/ani-branding/logo-light.png",
  "logo_dark_url": "https://minio.example/ani-branding/logo-dark.png",
  "favicon_url": "https://minio.example/ani-branding/favicon.ico",
  "primary_color": "#1677FF",
  "secondary_color": "#52C41A",
  "icp_number": "粤ICP备XXXXXXXX号"
}
```

## 错误示例

### PATCH 无权限（**TODO-YAML P2** 目标）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-brand-403-001"
}
```

### 非法主色（**TODO-YAML P2**）

```json
{
  "code": "BAD_REQUEST",
  "message": "primary_color must match #RRGGBB",
  "request_id": "req-boss-brand-400-001"
}
```

## 相关模块

- [`acceptance-manual.md`](acceptance-manual.md)（Phase2 验收 ACC-P2-02）
- ANI-07 §8

## 回填验收标准

- [x] 满配 22 章
- [x] GET 字段与 v1.yaml 一致；PATCH **TODO-YAML P2**
- [x] 响应 GET 示例 + 400/403
- [x] PRD/SPEC/HTML synced
