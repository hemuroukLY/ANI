# 内部 CA 与证书

## 页面定位

`内部 CA 与证书` 是 BOSS **交付与安装** 域下 **内部 PKI** 专页：安装时自动生成 ANI Root/Intermediate CA 与叶证书（ANI-07 §9.1），对应 `ani-installer new` **Step 6** 与 `add-san --ip`。

BOSS 提供 **根 CA 下载**、分发指引（§9.3）、cert-manager 轮换说明；**生成** 在 installer/cert-manager，非 BOSS REST。

与 [`ip-https.md`](ip-https.md) **强关联**（Dex IP SAN / OIDC HTTPS）。

权威源：ANI-07 §9.1、§9.3。REST：**TODO-YAML: N/A**（GET branding 无关）。Console **无** 对等页。

## 文档管理规则

- 主维护源：本文
- PRD/SPEC：prd-boss-internal-ca / spec-boss-internal-ca
- 交叉：ip-https、ani-installer、acceptance-manual

## Core 层要求

- CA 签发 REST — **N/A**（installer + cert-manager）
- 禁止 `/api/v1/boss/ca*`
- 叶证书 90 天自动轮换 — cert-manager（§9.1）
- BOSS 首次登录引导 §9.3 — 本页实现目标

## 页面职责

- 文档化 CA 层级：Root 10y → Intermediate 5y → 叶 90d
- 列出 SAN 覆盖：gateway/harbor/dex/内部通配符
- BOSS **下载 ani-root-ca.crt** 与 OS 分发步骤（§9.3）
- 链 ip-https（Dex issuer HTTPS 约束）
- `ani-installer add-san --ip` 追加 IP SAN 说明

## 页面结构

```text
内部 CA
├── CA 层级图（§9.1）
├── 叶证书与服务映射
├── 下载根 CA + 分发指引（§9.3）
├── cert-manager 自动轮换
└── 关联 ip-https / Step 6
```

## 数据来源与分层约束

| 层 | 来源 |
|---|---|
| ANI-07 §9 | PKI 设计 |
| CLI | new Step 6, add-san |
| cert-manager | 集群内 |
| BOSS UI | 下载按钮 **待产品** P1 |

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| CA 树 | §9.1 | — |
| 下载 | §9.3 | acceptance |
| Dex 约束 | §9.2 | ip-https |
| 导入证书 | Step 6 选项 | — |

## BOSS 与 Console 分工

| BOSS | Console |
|---|---|
| 根 CA 分发、平台安全 | 用户信任 CA 后 HTTPS 访问 |

## 当前冻结事实

| REST CA API | **N/A** |
| BOSS 下载 | **待产品** P1 UI |
| CLI | Step 6 自动 CA / 导入 |

### 证书字段（文档）

| 字段 | 说明 |
|---|---|
| `root_ca_expiry` | 10 年 |
| `intermediate_expiry` | 5 年 |
| `leaf_ttl_days` | 90 |
| `san_hosts[]` | IP/域名 |
| `issuer_type` | selfSigned / imported |

## 字段级定义

### CA 层级

| 实体 | TTL | 用途 |
|---|---|---|
| ANI Root CA | 10y | 信任锚 |
| Intermediate CA | 5y | 签发 |
| ani-gateway | 90d | 入口 TLS |
| dex | 90d | OIDC IP SAN |
| harbor | 90d | 镜像 |
| `*.ani-system.svc` | 90d | 内部 |

### BOSS 展示字段

| 字段 | 说明 |
|---|---|
| `ca_download_url` | BOSS 静态或 signed URL **待产品** |
| `leaf_expires_at` | 最近叶证书到期 |
| `auto_renew` | cert-manager 状态 |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 首次登录 | 引导下载 CA（§9.3） |
| 叶证书 <14d | 黄（应自动续） |
| 纯 IP 部署 | 链 ip-https |
| 导入模式 | 显示 issuer_type=imported |
| 403 | 无平台安全读 |

## 字段口径与单位

| 字段 | 格式 |
|---|---|
| TTL | 天/年 |
| SAN IP | IPv4/IPv6 |
| 文件 | PEM `.crt` |

## 状态与能力口径

| ca_status | 含义 |
|---|---|
| healthy | 自动轮换正常 |
| renew_pending | cert-manager 续期中 |
| import_invalid | 导入证书错误 |

## 创建前置条件

| 依赖 | 失败 |
|---|---|
| 安装 Step 6 完成 | 无 CA 可下载 |
| BOSS 读 | 403 |

## 操作可用性矩阵

| 操作 | 只读 | 安全管理员 |
|---|---|---|
| 查看 CA 文档 | ✅ | ✅ |
| 下载 root CA | ✅ P1 | ✅ |
| CLI 生成/导入 | CLI | CLI |
| REST 签发 | ❌ N/A | ❌ |

## 删除前置校验与当前契约边界

**N/A**（不通过 REST 删 CA）。轮换由 cert-manager 替换叶证书 — **非 DELETE 用户操作**。

## 接口冻结规则

- CA REST — **TODO-YAML: N/A**
- BOSS 下载 endpoint — **待产品** P1
- CLI Step 6 — ANI-07 冻结

## 使用规则

- Dex **必须** HTTPS + IP SAN — 见 ip-https
- 禁止伪造 CA REST
- 客户 IT 须安装根 CA 后再验收浏览器无警告（§12）
- add-san 仅追加 IP，不替换 Root

## 待补能力边界

- BOSS CA 下载 API — P1 产品
- 证书 inventory list — P2

## 响应示例

### 安装完成 CA 提示（CLI Step 11 节选）

```text
Internal CA generated (recommended).
Root CA: /var/lib/ani/certs/ani-root-ca.crt
Intermediate: auto · Leaf TTL: 90 days (cert-manager)
Dex issuer: https://10.20.30.40:5556/ (IP SAN configured)
BOSS: Settings → Internal CA → download root certificate
```

### BOSS 下载成功（页面目标 · 待产品）

```text
Content-Type: application/x-x509-ca-cert
Filename: ani-root-ca.crt
(PEM body)
```

## 错误示例

### 403

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-ca-403-001"
}
```

### 导入证书无效（CLI Step 6）

```text
ERROR: imported certificate missing IP SAN for Dex issuer 10.20.30.40
Use auto CA or add ipAddresses in cert (ANI-07 §9.2)
Exit code: 1
```

## 相关模块

- [`ip-https.md`](ip-https.md) · [`ani-installer.md`](ani-installer.md)
- [`acceptance-manual.md`](acceptance-manual.md) · [`baremetal-deploy.md`](baremetal-deploy.md)

## 回填验收标准

- [x] 满配 22 章
- [x] §9.1/§9.3 + ip-https 交叉链接
- [x] PRD/SPEC/HTML synced
