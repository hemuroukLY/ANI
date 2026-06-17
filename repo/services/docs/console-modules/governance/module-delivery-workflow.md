# Console 模块补充执行流程

本文档用于固化后续 `Console` 功能模块的标准补充流程，便于团队内部复用、交接和分享。

目标不是一次性把所有内容都写进 HTML，而是形成一套可重复执行、可一次性交付的规则：先按 `Core` 标准识别边界，再产出 `PRD`、`SPEC`，最后把模块定义收口到模块详文，并同步回填 HTML 摘要。

本流程在原有版本基础上新增了“零回合交付”要求：**在向用户汇报“已完成”之前，必须先完成权威源对照、跨材料一致性检查、反自造契约检查和最终阻塞项清零。**

## 1. 适用范围

- 适用于 `Console` 侧功能模块的新增、补充、收口和对齐工作
- 适用于需要同时满足 `Core` 规范、前端页面定义、接口对齐要求的模块
- 当前已按本流程执行并完成收口的模块范围，以 `docs/console-modules/governance/module-completion-matrix.md` 为准
- 如后续继续新增模块，先补模块产物，再更新完成矩阵，不在本文件重复维护模块清单

## 2. 总体原则

### 2.1 Core 优先

- 先判断模块能力归属 `Core` 还是 `Services`
- `Core` 资源类接口必须使用 `/api/v1/*`
- `Services` 业务类接口必须使用 `/api/v1/svc/*`
- 页面可以做聚合展示，但不能把 `Core` 资源错误写进 `Services` 契约

### 2.2 正文优先

- `HTML` 只保留摘要、边界、入口和导航关系
- `docs/console-modules/...` 下的模块详文是主维护源
- `tasks/prd-*.md`、`tasks/spec-*.md` 是辅助性材料
- 若正文与辅助材料冲突，以模块详文为准，并同步修正辅助材料

### 2.3 一次只收口一个模块

- 每次只处理一个明确模块
- 当前模块未达到“可直接对齐 Core”的标准前，不进入下一个模块
- 每个模块结束前都要做一次严格复审

### 2.4 零回合交付

- 不允许把“初稿完成”汇报成“已达到 Core 标准”
- 不允许把“方向正确”“基本一致”“只差一点”视为可交付
- 在向用户输出“已达到 Core 层要求”前，必须先通过本流程中的全部强制门禁
- 若内部复审仍发现阻塞项，必须**先自行修正**，而不是把找差异的工作留给用户

### 2.5 权威源优先于经验推断

- `ANI-main/repo/api/openapi/v1.yaml` 是当前一级权威源
- `ANI-main/repo/api/core-v1-compatibility-baseline.yaml` 是兼容性校对源
- 如 OpenAPI 与历史认知、旧 HTML、旧草稿冲突，一律以权威源为准
- 任何路径、schema、字段、返回码、错误码、鉴权方式，都必须能在权威源中找到对应证据，或明确标注为“待补能力”，不能模糊处理

### 2.6 禁止自造契约

- 不允许自造不存在于权威源的 schema 名称
- 不允许自造不存在于权威源的字段、状态、动作、错误码
- 不允许把“推断出的未来设计”写成“当前已冻结事实”
- 不允许把“前端可能需要的展示能力”写成“当前 Core 已支撑的字段能力”
- 不允许把 `idempotency_key` 这类请求控制字段误写成用户可见展示字段，除非权威源明确如此定义

### 2.7 最终确认范围必须明确

- 向用户输出“已达到 Core 层要求”时，必须明确说明**确认范围**
- 必须明确区分以下两类结论，不能混写：
  - 单模块结论：例如“`Console / 存储管理 / 对象存储` 已达到 Core 层要求”
  - 全量结论：例如“当前所有可用模块都已重新串审并达到 Core 层要求”
- 如果只复审了当前模块，禁止使用会让用户误以为“整批内容都已重新确认”的表述
- 如果用户问“现在写的内容都符合吗”，必须先判断其指代范围；若当前只完成单模块复审，应明确回答“就当前模块而言，已符合”

### 2.8 统一实例架构（GPU 容器 / Sandbox）

- GPU 容器实例必须使用 Core `/api/v1/instances*` + `kind=gpu_container`
- Sandbox 实例必须使用 Core `/api/v1/instances*` + `kind=sandbox`
- 不得在新文档或新代码中引用 `services/v1.yaml` 的 `/gpu-containers*`、`/sandboxes*`（已 deprecated）
- 模块详文必须包含**状态×操作可用性矩阵**与**创建前置条件**（含 `422 PRECONDITION_FAILED` 语义）

### 2.10 前置条件错误码口径

- `422` + `PreconditionFailed` 仅在对应 operation 的 YAML 已声明时，才可写成**已冻结返回码**
- YAML `PreconditionFailed.description` **已举例**的 `code`（如 Core 的 `IMAGE_NOT_FOUND`、`INSTANCE_STATE_INVALID`；Services 的 `MODEL_NOT_READY`）可写「YAML 已举例 `XXX`」
- 其余产品语义只允许写「具体 `code` 待 Core/Services 冻结；建议语义：…」，**禁止**用 `` `code: XXX` `` 定稿句式冒充冻结契约
- 若 YAML 尚未为该 operation 声明 `422`，须显式标注「**当前 YAML 未声明 `422`**」，不得省略

### 2.9 对象存储专项规则

- `对象存储` 当前只承接 `StorageObject` 对象元数据资源
- `创建对象元数据` 不等于 `上传文件`
- `bucket` 当前只作为 `StorageObject` 的字段，不等于已冻结的独立桶资源
- `对象上传 / 对象下载 / 桶资源 / 桶策略` 当前只能写成待补能力，不能写成已冻结接口
- 如果未来要扩充桶策略，正文和 `SPEC` 在当前阶段只能写“桶策略相关能力”或“桶策略待补”，不能自造过于具体的未来路径

## 3. 标准产出物

每个模块默认产出以下 4 类文件：

### 3.1 HTML 摘要入口

- `prototypes/ani-services-prototype-console.html`
- 必要时同步 `prototypes/ani-services-prototype.html`

作用：

- 保留模块名称、定位、边界、摘要说明、详文入口
- 不承载超长字段定义和完整接口契约

### 3.2 模块详文

路径示例：

- `docs/console-modules/home/platform-overview.md`
- `docs/console-modules/compute/gpu-management.md`

作用：

- 作为该模块唯一主维护源
- 用于沉淀字段定义、分层边界、接口冻结规则、响应示例、验收标准

### 3.3 PRD 辅助文档

路径示例：

- `tasks/modules/prd/console/compute/prd-console-gpu-management.md`

作用：

- 记录模块目标、用户故事、范围、非目标、产品验收口径

### 3.4 SPEC 辅助文档

路径示例：

- `tasks/modules/spec/console/compute/spec-console-gpu-management.md`

作用：

- 记录接口约束、字段模型、错误结构、示例、依赖项

### 3.5 PRD / SPEC 与主维护源同步规则

- `PRD / SPEC` **不得**与主维护源（`docs/console-modules/**`）在路径、返回码、待补边界上冲突
- 主维护源变更后，须回看同名 `tasks/modules/prd/prd-console-*.md`、`tasks/modules/spec/spec-console-*.md`；口径冲突时**以主维护源 + OpenAPI 为准**
- 前置条件与 `422` 的 `code` 写法必须遵守 §2.10（TASK 清单中的实现验收可写目标语义，但须标注「YAML 已举例」或「待冻结」）
- 租户管理、SSO 读写等 Services 能力不得在安全概览 PRD/SPEC 中继续写成「全部待补」

## 4. 标准执行步骤

### 第 0 步：锁定权威源

开始任何模块前，先读取并锁定以下内容：

1. `ANI-main/repo/api/openapi/v1.yaml`
2. `ANI-main/repo/api/core-v1-compatibility-baseline.yaml`
3. 上一级模块主维护源（如果当前是子模块）
4. 当前 HTML 摘要入口

输出要求：

- 列出该模块对应的真实冻结路径
- 列出该模块对应的真实 schema
- 列出真实成功返回码和错误返回码
- 列出当前**没有**冻结、但容易被误写进去的待补能力

### 第 1 步：识别模块归属

先回答 3 个问题：

1. 这是 `Console` 侧模块还是 `BOSS` 侧模块
2. 这个模块的底层资源归 `Core` 还是 `Services`
3. 该页面是资源管理页、业务聚合页，还是首页摘要页

输出要求：

- 明确模块定位
- 明确 `Core / Services` 边界
- 明确是否允许聚合展示

### 第 2 步：套用 Core 标准基线

在开始写模块内容前，先统一套用以下基线：

- 不允许继续使用旧 `/api/v1/console/*` 路径
- 不默认要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- 写操作如为 `POST` 或有副作用动作，默认要求 `idempotency_key`
- 错误结构统一为：

```json
{
  "code": "UPPER_SNAKE",
  "message": "error message",
  "request_id": "req-xxx"
}
```

- 写接口必须明确成功状态码，不允许含糊描述
- 如果动作是异步的，必须明确 `202 Accepted + AsyncTask` 或等价规范

同时必须生成一份“冻结事实表”，至少包含：

- `Frozen Paths`
- `Frozen Schemas`
- `Frozen Response Codes`
- `Non-Frozen Capabilities`
- `Known Risky Assumptions`

这张表不是最终交付物，但必须作为后续 `PRD / SPEC / 主维护源 / HTML` 的统一事实底板。

### 第 3 步：生成 PRD

使用 `prd` 技能生成该模块的产品文档。

PRD 至少应包含：

- 模块目标
- 页面定位
- 用户角色
- 用户故事
- 范围与非目标
- 页面区块或关键能力
- 验收标准
- 开放问题
- 回填前置依赖

### 第 4 步：生成 SPEC

使用 `prd-to-spec` 技能生成技术约束文档。

SPEC 至少应包含：

- 技术范围
- 分层边界
- 数据模型或字段模型
- 接口列表
- 每个接口的路径、方法、用途
- 写接口的幂等要求
- 成功响应示例
- 错误响应示例
- 回填 `Core v1.yaml` 或 `Services v1.yaml` 前置依赖

新增强制要求：

- `SPEC` 中每一条“已冻结路径 / 已冻结 schema / 已冻结返回码”都必须能回指到权威源
- `SPEC` 中可以描述待补能力，但不能给待补能力编造 schema 名称
- `SPEC` 中如果提到页面展示能力，必须先确认该展示能力有现成字段或现成接口支撑；若没有，只能降级为风险提示、占位说明或待补项
- `SPEC` 中如果描述对象存储待补能力，必须明确“当前仅承接对象元数据”；不能把上传、下载、桶资源、桶策略写成既有冻结动作

### 第 5 步：收口为模块详文

将最终稳定口径写入 `docs/console-modules/...` 下的模块详文。

模块详文建议固定包含以下章节：

- 页面定位
- 文档管理规则
- Core 层要求
- 页面职责
- 页面结构
- 数据来源与分层约束
- 字段级定义
- 字段展示规则
- 字段口径与单位
- 状态与能力口径
- 前置条件
- 操作可用性矩阵
- 删除前置校验与当前契约边界
- 接口冻结规则
- 响应示例
- 错误示例
- 回填前置依赖
- 待确认项
- 回填验收标准

要求：

- 模块详文必须比 `PRD`、`SPEC` 更稳定
- 如辅助材料与正文不一致，必须同步修正辅助材料
- 最终要达到“可直接拿去对齐 Core”的程度
- 主维护源中不能保留只在 `SPEC` 出现、但正文未解释的字段处理差异
- 主维护源中必须明确哪些字段是“展示字段”，哪些只是“请求约束 / 后端返回字段 / 开发态辅助字段”
- 对象存储类页面必须明确区分“对象元数据展示/创建”与“对象内容上传/下载”
- 若某字段当前只是已冻结 schema 字段，不代表同名资源域已成为独立冻结能力；例如 `bucket` 作为对象字段，不等于已有桶资源契约

### 第 6 步：回填 HTML 摘要

在 `prototypes/ani-services-prototype-console.html` 中只保留：

- 模块标题
- 模块目标
- 页面定位摘要
- 核心区块摘要
- 文档拆分与维护规则

必要时同步总览文件 `prototypes/ani-services-prototype.html`。

要求：

- HTML 不再堆叠超长字段定义和完整接口契约
- HTML 摘要口径必须与模块详文一致
- HTML 不能把已经清理掉的旧接口路径、旧错误码、旧边界重新带回来

### 第 7 步：严格复审

模块完成后，必须按以下清单复审：

- 是否明确了 `Core / Services` 归属
- 路径前缀是否正确
- 是否误要求前端显式传 `tenant_id / X-Tenant-Id`
- 写操作是否补了 `idempotency_key`
- 是否明确 `operationId / summary / security / responses`
- 是否统一使用标准错误结构
- 是否混入 `BOSS` 平台运营语义
- 主维护源、`PRD`、`SPEC`、HTML 是否已经一致
- 是否还存在会影响 `v1.yaml` 回填的开放契约
- 是否存在任何“权威源里没有、但文档自己写出来”的 schema / 字段 / 动作 / 错误码
- 是否把请求控制字段误写成用户展示字段
- 是否把待补能力写得比当前 Core 现状更具体
- 是否把页面能力写得超过现有 schema/接口可支撑范围
- 若为对象存储，是否把“对象元数据创建”误写成“文件上传”
- 若为对象存储，是否把 `bucket` 字段误写成独立桶资源能力
- 对最终答复是否明确声明了“本次确认范围”

复审必须分 4 轮执行，不能只做一遍表面检查：

1. **权威源对照**
   - 逐项核对路径、schema、字段、返回码、security
2. **跨材料一致性检查**
   - 交叉比对 `PRD / SPEC / 主维护源 / HTML`
3. **反自造契约检查**
   - 检查是否编造 schema、字段、状态、错误码、动作
4. **页面能力降噪检查**
   - 检查是否把“风险提示 / 待补说明 / 跳转入口”误写成“当前已支撑的展示能力”
5. **最终答复范围检查**
   - 检查最终给用户的结论是否明确说明了“当前模块”还是“全量已做内容”

新增“强制阻塞项”如下，只要命中任意一条，就**禁止**向用户汇报“已达到 Core 标准”：

- 文档中存在权威源里没有的 schema 名称
- 文档中存在权威源里没有的返回码，却写成已冻结契约
- 文档中把待补能力写成正式接口表
- 文档中把请求字段写成用户展示字段
- 文档中把无字段支撑的页面能力写成既有能力
- `PRD / SPEC / 主维护源 / HTML` 中任意一份仍保留旧路径或旧边界
- 最终答复没有明确本次确认范围，却直接输出“都已符合”之类的总结性表述

只有在所有阻塞项清零后，才允许输出“已达到 Core 层要求”。

结论必须明确输出为二选一：

- 已达到 `Core` 层要求
- 尚未达到 `Core` 层要求，并列出阻塞项

最终答复时必须补 1 句范围说明，格式建议如下：

- `就当前模块 <模块名> 而言，已达到 Core 层要求`
- `本次结论仅覆盖 <模块名>，不自动外推到所有可用模块`
- 若做的是全量串审，才允许写：`当前可用模块已重新串审，均达到 Core 层要求`

## 5. 文件管理规则

### 5.1 推荐目录

- `docs/console-modules/home/`
- `docs/console-modules/compute/`
- `docs/console-modules/inference/`
- `docs/console-modules/knowledge/`
- `docs/console-modules/alerts/`
- `docs/console-modules/tenant/`

### 5.2 命名建议

- 模块详文：`<module>.md`
- PRD：`prd-console-<module>.md`
- SPEC：`spec-console-<module>.md`

示例：

- `gpu-management.md`
- `prd-console-gpu-management.md`
- `spec-console-gpu-management.md`

### 5.3 修改顺序

固定顺序如下：

1. 先改模块详文
2. 再看 `PRD`、`SPEC` 是否需要同步
3. 最后回填 HTML 摘要

新增收尾要求：

4. 最后重新用权威源回查一次正文与 HTML 摘要
5. 只有最终复审通过后，才向用户交付结论

## 6. 可直接复用的模块模板

后续任一模块都可以按以下顺序直接执行：

1. 明确模块归属与边界
2. 锁定权威源并抽取冻结事实表
3. 按 `Core` 标准建立检查基线
4. 生成 `PRD`
5. 生成 `SPEC`
6. 收口到模块详文
7. 回填 HTML 摘要
8. 做最终 Core 合规复审
9. 仅在阻塞项清零后再向用户交付
10. 交付时明确声明“本次确认范围”

## 7. 当前推荐顺序

当前主摘要层收口已完成；若后续继续维护、新增模块或复查旧模块，建议统一按以下顺序工作：

1. 先查看 `docs/console-modules/governance/final-handoff-summary.md`
2. 再查看 `docs/console-modules/governance/console-document-status-board.md`
3. 再查看 `docs/console-modules/governance/governance-closeout-checklist.md`
4. 再进入对应模块主维护文档核对正式路径、schema、返回码和待补边界
5. 最后再回看 `prototypes/ani-services-prototype-console.html` 是否需要同步摘要

## 8. 分享说明

若将本文档分享给协作者，建议同时说明以下 3 点：

- `HTML` 不是完整契约正文，只是摘要入口
- 模块详文才是主维护源
- 是否达到 `Core` 标准，必须以最终复审结论为准，不能只看 `PRD` 或 `SPEC`
- 用户不需要承担“帮忙找最后几个契约问题”的职责；流程本身必须先把这些问题清理掉
- 向用户做最终确认时，必须明确“当前模块通过”还是“全量串审通过”，不能省略范围
