# 纯 IP + HTTPS（Dex SAN）

## 页面定位

`纯 IP + HTTPS` 是 BOSS **交付与安装** 域下 **无公网域名** 部署专页：访问地址为纯 IP 时，仍须 **HTTPS** 与正确 **IP SAN**（ANI-07 §9、§2.2 Step 5–6）。

**Dex OIDC issuer 必须 HTTPS**（§9.2 关键约束）— 与 [`internal-ca.md`](internal-ca.md) 强关联。

BOSS 文档化组件适配表 §9.2、`ani-installer add-san --ip`；执行在 installer/cert-manager。

REST：**TODO-YAML: N/A**。Console **无** 对等页。

## 文档管理规则

- 主维护源：本文
- PRD/SPEC：prd-boss-ip-https / spec-boss-ip-https
- 交叉：internal-ca、acceptance-manual

## Core 层要求

- IP/HTTPS 配置 REST — **N/A**
- Gateway/Dex TLS — 集群配置 + 内部 CA
- 禁止 boss 自造 path

## 页面职责

- §9.2 组件无域名适配清单（cert-manager/Gateway/Harbor/Dex/…）
- Step 5 `access_host` 纯 IP 口径
- Dex `issuer: https://<IP>:<port>` + IP SAN 证书
- `add-san --ip` 追加节点 IP
- 验收 §12「纯 IP 访问」标准

## 页面结构

```text
纯 IP + HTTPS
├── 约束说明（Dex HTTPS）
├── §9.2 组件适配表
├── Step 5–6 与 CA 关系
├── add-san 用法
└── 验收与 CA 信任
```

## 数据来源与分层约束

| 层 | 来源 |
|---|---|
| ANI-07 §9 | 权威 |
| internal-ca | 证书 |
| CLI | add-san, Step 5–6 |

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| Dex 约束 | §9.2 | internal-ca |
| 组件表 | §9.2 | — |
| 访问 URL | Step 11 | acceptance |
| SAN 追加 | §2.1 add-san | internal-ca |

## BOSS 与 Console 分工

| BOSS | Console |
|---|---|
| 文档 + 安全引导 | `https://<IP>` 登录 |

## 当前冻结事实

| REST | N/A |
| CLI | add-san, new Step 5–6 |
| 验收 §12 | 无证书警告 |

### 配置字段

| 字段 | 说明 |
|---|---|
| `access_ip` | 管理/访问 IP |
| `dex_issuer_url` | https://IP:port |
| `gateway_url` | https://IP/ |
| `san_ips[]` | 证书 IP 列表 |
| `use_domain` | false |

## 字段级定义

### §9.2 组件行

| 组件 | 适配 | 注意 |
|---|---|---|
| cert-manager | ipAddresses | 原生 |
| ANI Gateway | IP SAN TLS | ✅ |
| Harbor | hostname=IP | ✅ |
| Dex | IP SAN + HTTPS issuer | ⚠️ OIDC |
| Grafana | 集群内 DNS | 无需外部 |
| vLLM/Milvus | 内部 | HTTP/gRPC |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 纯 IP 模式 | 显示 access_ip |
| Dex 非 HTTPS | 红阻止说明 |
| 浏览器警告 | 链 CA 下载 |
| 多节点 IP | add-san 列表 |

## 字段口径与单位

| 字段 | 格式 |
|---|---|
| `access_ip` | IPv4/IPv6 |
| `dex_issuer_url` | https URI |
| 端口 | 1–65535 |

## 状态与能力口径

| https_ready | 含义 |
|---|---|
| true | CA 信任后无警告 |
| ca_pending | 待客户导入根 CA |
| misconfigured | Dex/Gateway SAN 缺失 |

## 创建前置条件

| 依赖 | 失败 |
|---|---|
| Step 6 CA | 无 TLS |
| IP 规划 | 网络不可达 |
| BOSS 读 | 403 |

## 操作可用性矩阵

| 操作 | BOSS | CLI |
|---|---|---|
| 文档 | ✅ | ✅ |
| add-san | ❌ | ✅ |
| REST | ❌ N/A | ❌ |

## 删除前置校验与当前契约边界

**N/A**

## 接口冻结规则

REST **N/A**；Dex issuer HTTPS — ANI-07 §9.2 冻结。

## 使用规则

- **禁止** HTTP 对外 Dex issuer
- 纯 IP **仍须** internal CA 或合法导入 cert
- 验收前客户 IT 须信任根 CA
- 与 internal-ca 同步维护

## 待补能力边界

- BOSS SAN 列表只读 — P2
- REST — N/A

## 响应示例

### Step 11 纯 IP 输出（CLI）

```text
Installation complete (pure IP mode).
Console: https://10.20.30.40/
Dex issuer: https://10.20.30.40:5556/
Import ani-root-ca.crt on client machines to avoid browser warnings
Test: open Console URL — expect secure lock after CA trust
```

## 错误示例

### 403 BOSS

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-ip-403-001"
}
```

### Dex HTTP 误配（运维失败）

```text
ERROR: Dex configured with http://10.20.30.40:5556/ — OIDC requires HTTPS (ANI-07 §9.2)
Fix: re-run certificate step or ani-installer add-san --ip 10.20.30.40
```

## 相关模块

- [`internal-ca.md`](internal-ca.md) · [`ani-installer.md`](ani-installer.md)
- [`acceptance-manual.md`](acceptance-manual.md) · [`baremetal-deploy.md`](baremetal-deploy.md)

## 回填验收标准

- [x] 满配 22 章
- [x] §9.2 Dex 约束 + internal-ca 链接
- [x] PRD/SPEC/HTML synced
