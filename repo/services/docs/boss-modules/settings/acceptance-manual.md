# 交付验收手册

## 页面定位

`交付验收手册` 是 BOSS **交付与安装** 域下 **验收测试与 SOP** 专页：对应 ANI-07 §10.2《验收测试手册》、§12 验收要点、§10.3 标准交付 SOP（T+1 签署验收单）。

BOSS 提供 **检查清单**、通过标准、与现场测试结果记录入口（文档/表单 — REST **N/A**）。

Console **无** 对等页。

## 文档管理规则

- 主维护源：本文
- PRD/SPEC：prd-boss-acceptance-manual / spec-boss-acceptance-manual
- 链：baremetal/vm/k8s-attach、internal-ca、ip-https、gpu-driver-install

## Core 层要求

- 验收 REST — **TODO-YAML: N/A**
- 计量验收 §12 可 **只读参考** `metering` API — **非** 本页主路径
- 禁止 boss 验收 path 自造

## 页面职责

- §12 测试场景表：裸机/new、attach、纯 IP、GPU、计量
- §10.3 SOP 时间线与《验收确认单》
- §10.2 交付文档清单（PDF/Excel/MP4）
- CA 信任配置步骤（T+1 上午 — §9.3）
- Phase2 项标注（Patch、白牌、AZ）

## 页面结构

```text
交付验收
├── §12 验收场景表
├── 逐步操作清单（可勾选）
├── SOP T+1 流程
├── 交付文档下载索引
└── 签署与移交
```

## 数据来源与分层约束

| 层 | 来源 |
|---|---|
| ANI-07 §10、§12 | 权威 |
| CLI/status | 健康确认 |
| BOSS | 清单 UI |

## 页面区块与数据来源映射

| 场景 | 验证方法 | 标准 | 来源 |
|---|---|---|---|
| 裸机安装 | new | ≤2h 无错 | §12 |
| attach | attach-k8s | 平台启动 | §12 |
| 纯 IP | 浏览器 | 无警告 | §12 |
| GPU | nvidia-smi | 驱动正确 | §12 |
| 计量 | 推理10次 | Token 准确 | §12 |

## BOSS 与 Console 分工

| BOSS | Console |
|---|---|
| 验收清单、签署流程 | 功能试用 |

## 当前冻结事实

| REST | **TODO-YAML: N/A** |
| Phase2 测试 | Patch/白牌/AZ — 非 v1.0 阻塞 |

###  checklist 字段

| 字段 | 说明 |
|---|---|
| `test_id` | 场景 ID |
| `title` | 场景名 |
| `method` | 验证方法 |
| `pass_criteria` | 通过标准 |
| `result` | pass/fail/skip |
| `evidence` | 截图/日志 |
| `tester` | 工程师 |
| `tested_at` | 时间 |

## 字段级定义

### §12 行

| test_id | title | pass_criteria |
|---|---|---|
| ACC-01 | 裸机全新安装 | ≤2h 无报错 |
| ACC-02 | 复用 K8s | attach 后正常 |
| ACC-03 | 纯 IP | HTTPS 无警告 |
| ACC-04 | GPU 驱动 | smi 正确 |
| ACC-05 | 计量 | Token 准确 |
| ACC-P2-01 | Patch 升级 | Phase2 |
| ACC-P2-02 | 白牌化 | Phase2 |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| pass | 绿勾 |
| fail | 红 + 阻塞签署 |
| Phase2 | 灰「非 v1.0 门禁」 |
| 缺 CA 信任 | ACC-03 失败提示 |

## 字段口径与单位

| 字段 | 格式 |
|---|---|
| `tested_at` | ISO 8601 |
| 安装时长 | 小时 |
| Token 数 | 整数 |

## 状态与能力口径

| result | 可签署 |
|---|---|
| all_pass | 是 |
| has_fail | 否 |
| partial_p2_skip | v1.0 可 |

## 创建前置条件

| 依赖 | 失败 |
|---|---|
| 安装完成 Step 11 | 无可测 |
| BOSS 读 | 403 |
| CA 分发（ACC-03） | 纯 IP 失败 |

## 操作可用性矩阵

| 操作 | 只读 | 交付工程师 |
|---|---|---|
| 查看清单 | ✅ | ✅ |
| 勾选结果 | ✅ P1 | ✅ |
| REST 提交 | ❌ N/A | ❌ |
| 签署 PDF | 线下 | 线下 |

## 删除前置校验与当前契约边界

**N/A**

## 接口冻结规则

REST **N/A**；§12 场景 ANI-07 冻结。

## 使用规则

- v1.0 **不** 阻塞 Phase2 项
- ACC-03 **依赖** internal-ca 客户信任
- 不得写验收 API 已实现
- SOP 顺序：安装 → CA → 验收 → 签署

## 待补能力边界

- BOSS 验收记录持久化 — P1
- PDF 导出 — P2

## 响应示例

### 验收记录（BOSS 页面目标 · 本地表单）

```json
{
  "site": "customer-a-site-01",
  "tests": [
    {
      "test_id": "ACC-01",
      "result": "pass",
      "evidence": "install log 2026-06-17, duration 1h42m",
      "tested_at": "2026-06-17T11:00:00Z"
    },
    {
      "test_id": "ACC-03",
      "result": "pass",
      "evidence": "Chrome secure lock on https://10.20.30.40/",
      "tested_at": "2026-06-17T11:30:00Z"
    }
  ],
  "signoff_ready": true
}
```

## 错误示例

### 403

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-acc-403-001"
}
```

### ACC-03 失败 — 未信任 CA

```text
FAIL ACC-03: Browser shows NET::ERR_CERT_AUTHORITY_INVALID
Action: Install ani-root-ca.crt per internal-ca guide, re-test
Blocks sign-off until pass
Exit code: 1
```

### 验收项未填（表单校验 · 400 类）

```json
{
  "code": "BAD_REQUEST",
  "message": "result and evidence required for ACC-01 before sign-off",
  "request_id": "req-boss-acc-400-001"
}
```

## 相关模块

- [`ani-installer.md`](ani-installer.md) · [`internal-ca.md`](internal-ca.md) · [`ip-https.md`](ip-https.md)
- [`baremetal-deploy.md`](baremetal-deploy.md) · [`k8s-attach.md`](k8s-attach.md) · [`gpu-driver-install.md`](gpu-driver-install.md)

## 回填验收标准

- [x] 满配 22 章
- [x] §10/§12 完整
- [x] PRD/SPEC/HTML 与本文同步
