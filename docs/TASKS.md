---
status: active
last-verified: 2026-10-23
verified-by: KT-04 结案归档（4xx 返回 AppError、重试保留 body；新增 httpclient → errors 边，ADR 0002）
---

# TASKS — tellurion-kit 未完成项

> 本仓**唯一**任务台账。跨仓任务只在 §二 留指针，**不镜像正文**（双写必然漂移）。
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `✅ 完成（附日期）`
> 每张单六字段（① 目标 ② 上下文 ③ 要求 ④ 约束与非目标 ⑤ 验收 ⑥ 边界），判据见 `tellurion:AGENTS.md` §8.5。
> 完成即标 ✅（附日期）并移入 `docs/archive/TASKS-YYYY-MM-DD-archived.md`；本文件只留未完成项。

---

## 一、本仓任务

> **本仓已不是占位仓**（2026-10-22 起）：QT-05-B 落码完成，13 包全部有实现；KT-01/KT-02 两张随仓文档交付在此之前已结案。
> 现在台账上的是**落码过程中实测出来的既有缺陷**（KT-04 / KT-05 照搬自 core、QT-05-B 按 ④「不改任何业务逻辑」刻意未动；KT-06 是 KT-03 取证时新发现）。
>
> ✅ **KT-01 已结案（2026-10-22）**：`docs/SCOPE.md`（13 包逐包准入/禁入 + 词表 + 实测偏差九条）+ `tools/check_scope`（CI 步骤「包边界核对」）落地，
> 整条搬进 [`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md)。
> **对 QT-05-B 的直接后果**：落码时包必须先进 `SCOPE.md` §3 那张表（少一行是 fatal，不是提醒）；§5 那九条偏差里有四条会让照搬即红，先读那张表再动手。
>
> ✅ **KT-02 已结案（2026-10-22）**：README 的版本规则表 + 兼容矩阵 + `CHANGELOG.md` + AGENTS 的 `GOPRIVATE` 前提落地，
> 护栏 `tools/check_release_matrix` 进 CI 步骤「版本纪律核对」，整条搬进同一份归档。
> **对 QT-05-B / QT-07 的直接后果**：S2 打第一个 tag（`v0.1.0`）那天起，README 兼容矩阵必须有 `v0.1` 那一行、
> CHANGELOG 必须有 `## [v0.1.0] - 日期` 那一节，漏一条 `go run ./tools/check_release_matrix` 红（CI 同款）。tag 那条腿今天空跑（本仓 0 tag）。
> 另：`tools/` 下的 `package main` 被 `check_scope` 整目录跳过，新工具**不需要**在 `SCOPE.md` §3 加行（实测「2 个 main 包跳过」）。
>
> 第一张单长什么样（照此写，别自造格式）：
>
> ### XX-01 · <一句话目标> [P?] ⬜
> - **① 目标**：…
> - **② 上下文**：本仓路径/行号 + 冻结契约的路径与节号
> - **③ 要求**：①②③（可执行、可核对）
> - **④ 约束与非目标**：明确不做什么
> - **⑤ 验收**：可复现的命令 + 期望输出；至少一条"破坏验证"（改坏它应当红）
> - **⑥ 边界**：只做什么；相邻的什么归别的单

### KT-01 · 随仓接口规格载体 `docs/SCOPE.md`（12 包逐包准入/禁入清单）[P1] ✅ 完成（2026-10-22，已归档）

> 六字段正文与结案记录（⑤ 验收真退出码、现场破坏验证的还原凭据、22 条夹具腿、词表与覆盖面限制的理由）在
> [`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对 QT-05-B 的直接后果（读 `SCOPE.md` §5 那九条，别照搬）**：清单落的是 **13 行**（12 横切包 + `db`）；
> `pkg/id` 的 `OrderID`、`pkg/clock` 的 `ModeBacktest`、`pkg/observability` 的 `BacktestDuration`/`ObserveBacktest` 撞词表，
> `pkg/auth` 带 `users`/`audit_logs` 的 SQL 撞「零表结构」——这四条**不解决就搬，本仓 CI 第一天就红**；
> 其中 auth 那条是裁决（SQL 留 core 还是 auth 暂不进本仓），不是改名能过关的。

### KT-02 · README 版本规则表 + 兼容矩阵 + CHANGELOG + `GOPRIVATE` 前提 [P1] ✅ 完成（2026-10-22，已归档）

> 六字段正文与结案记录（③ 逐条落点、⑤ 验收真退出码、现场破坏验证的三条腿与还原凭据、六条口径偏差登记）在
> [`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md)。本文件按规则只留未完成项，这里只留指针。
> **对 QT-05-B / QT-07 的直接后果**：S2 打第一个 tag（`v0.1.0`）那天起，`README.md` 的兼容矩阵必须加 `v0.1` 那一行、
> `CHANGELOG.md` 必须把 `[Unreleased]` 的内容搬进 `## [v0.1.0] - 日期` 节，否则 `go run ./tools/check_release_matrix` 红
> （CI「版本纪律核对」同款）。tag 那条腿今天空跑（0 tag），从第一个 tag 开始才咬人。

### QT-05-B · S1 工具层落地：12 个包 + `db` 迁入本仓（QT-05 的 kit 腿）[P0] ✅ 完成（2026-10-22，已归档）

> 六字段正文与结案记录（13 包逐包落点、⑤ 验收真退出码、三条真服务集成腿、SCOPE §5 九条偏差的裁决、
> 两条**没做**的腿——现场破坏验证与 `-race`——的如实登记）在
> [`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对后续单的直接后果**：
> - **QT-05 主单（contracts 腿 + core 侧依赖切换）**：core 切 `require` 本仓时会撞上 7 处形状差异（`id` 无 `OrderID`、`statistics` 无量化口径、
>   `auth` 无 SQL、`authz` 的 `Policy` 是机制、`observability` 无业务指标、`clock` 的 `ModeBacktest` → `ModeDeterministic`、`sandbox/wasm` 不在本仓），
>   逐条在 `docs/SCOPE.md` §5。**这不是"顺手改一下"能过的**：core 侧要么保留自己的实现，要么按裁决改调用点。
> - **QT-07（S2 打 tag）**：`v0.1.0` 那天 README 矩阵要有 `v0.1` 行、CHANGELOG 要把 `[Unreleased]` 搬进版本节（KT-02 的护栏会红）。
>   ⚠️ 打 tag 前应先裁决 **KT-04 / KT-05 / KT-06**：都属导出面语义或底层行为，`v0.1.0` 之后再改就是 breaking（要按 README 规则记账）。
>   KT-03 已于 2026-10-23 结案（`Div` 口径冻结为 half away from zero）；KT-04 同日结案（4xx → `*AppError`、重试保留 body、新增 `httpclient → errors` 边）。

### KT-03 · `decimal.Div` 的舍入口径与源码注释不一致（裁决口径 + 让注释与实现同一句话）[P1] ✅ 完成（2026-10-23，已归档）

> 六字段正文与结案记录（① 取证的真退出码、② half-away-from-zero 的裁决理由、破坏验证、core 侧处置）在
> [`archive/TASKS-2026-10-23-archived.md`](archive/TASKS-2026-10-23-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对 KT-06 / QT-07 的直接后果**：口径已冻结为 **half away from zero（四舍五入）**，`docs/api-schema.md` §1 在册；
> KT-06 只动定标与溢出路径，**不得回改口径**；S2 打 `v0.1.0` 时 CHANGELOG 的 breaking 条目已就位（`[Unreleased]` 里第一条）。

### KT-04 · `httpclient` 重试发空 body + 4xx 静默当成功（两条都已实测取证）[P1] ✅ 完成（2026-10-23，已归档）

> 六字段正文与结案记录（① 复现测试的先红后绿真退出码、③ 4xx → `*AppError` 的裁决与错误码映射、④ 重试不看方法的文档钉死、破坏验证、新增 `httpclient → errors` 边与 ADR 0002、core 侧处置）在
> [`archive/TASKS-2026-10-23-archived.md`](archive/TASKS-2026-10-23-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对 QT-07 / 使用方的直接后果**：4xx 的返回语义已从「`err == nil` 的成功」改为 `nil` + `*errors.AppError`（**breaking**，CHANGELOG `[Unreleased]` 在册），
> `err == nil` 从此只代表 2xx/3xx；重试仍**不看方法**，有副作用的调用须显式 `New(baseURL, timeout, 0)`。口径唯一真相在 `docs/api-schema.md` §3。

### KT-05 · `errors` 零单测、`logging` 的初始化路径未测（含一个包级可变全局）[P2] ⬜

- **① 目标**：本仓最底层的两个包（所有仓都会 `require`）有**自己的**单测，其导出面的语义变化能被本仓 CI 当场逮住，而不是等使用方红。
- **② 上下文**：实测（2026-10-22，`go test ./... -count=1`）：`errors` 包**没有测试文件**（输出 `? …/errors [no test files]`），
  core 的 `pkg/errors` 同样没有——QT-05-B 逐字照搬，所以这个空洞是**继承来的**，不是本次引入的。
  `errors` 今天只被间接覆盖：`httpserver/errors_test.go` 有 12 处用 `apperrors.ErrCode*` 走「分类 → 状态码」那条链
  （`StatusForAppError`），但 `Wrap`/`Wrapf` 的 `%w` 链、`AppError.Is` **按分类码相等**（不是实例身份）、`ToMap`、`captureStack` 都没有直接断言。
  `logging` 只有 `masking_test.go`；`logger.go` 的 `Init(level, format)` 与 `WithContext(map)` 未测，
  且 `logging.Logger` 是**包级可变全局**（`logger.go:12`），`Init` 写它 —— 并行测试里同时 `Init` 就是数据竞争，
  而本机 `-race` 不可用（无 cgo/gcc，见 QT-05-B 结案记录），所以这类问题今天**测不出来**。
- **③ 要求**：① `errors` 补表驱动单测：`Is` 按码相等的正/反例（同码不同 message → 真；不同码 → 假）、`Unwrap` 链能被 `errors.Is`/`As` 穿透、
  `WithOperation`/`WithCause` 的返回值语义（是否改原值）、`ToMap` 的字段集合、`Error()` 文案在有无 `Cause` 两种形状下的断言；
  ② `logging` 补 `Init` 的级别与格式（json/console）解析断言 + `WithContext` 的字段注入断言；
  ③ 就 `Logger` 全局给一个处置：要么在文档里钉死「`Init` 只能在启动期调一次，测试不得并发调」，
  要么把它改成可注入（新导出面 = MINOR，且会影响所有使用方，**要单独裁决不要顺手改**）；
  ④ 若 ③ 选了「文档钉死」，`docs/api-schema.md` §1 的 `logging` 行要写这条约束。
- **④ 约束与非目标**：**不改 `errors` / `logging` 的现有行为**（补测暴露出的行为问题另立单，按 KT-03 的两步法处理）；
  不引第三方断言库之外的新依赖（`testify` 已在 `go.mod`）；不做 `zerolog` 的全局钩子改造。
- **⑤ 验收**：`go test ./errors/... ./logging/... -count=1` 绿且**不再是 `[no test files]`**；
  `go test ./... -count=1 -cover` 里 `errors` 与 `logging` 的覆盖率有非零数字（把数字写进结案记录，不设阈值门槛——设了就成了凑数）；
  **破坏验证**：把 `AppError.Is` 改成 `return false` → 新增测试必须红；把 `MaskMap` 的敏感键判断去掉 → `masking_test.go` 必须红（证明新测与旧测都在咬）。
- **⑥ 边界**：只做这两个包的测试与全局的处置裁决；`decimal` / `httpclient` 的行为缺陷是 KT-03 / KT-04。

### KT-06 · `decimal.Div` 定标路径的静默精度/量纲损失（KT-03 取证时新发现）[P2] ⬜

- **① 目标**：`Div` 在**任意** `resultScale` 下都给出「精确商 + 一次舍入」的结果，或在无法精确时显式返回 `ErrOverflow`，不再静默给出错量级/多次截断的值。
- **② 上下文**：KT-03 ① 取证时实测（`decimal/div_rounding_test.go` 的 `TestDiv_ResultScaleBelowDividendScale`）：
  `Div(New(25,1), New(1,0), 0)`（2.5 / 1）→ **2**，正确应是 3——`shift < 0` 分支在 `decimal.go` 把被除数逐位 `/10` 截断后再除，信息在舍入前就被丢掉；
  另一条：`shift > 0` 分支在 `numerator` 越过 `1<<62` 时直接 `break`（不再乘 10 的幂），于是商**静默小了 10^k 倍**（量纲错误，不是精度损失）。
  两条都在 KT-03 之前就存在，KT-03 只裁决「恰好一半怎么走」，按「不顺手修未列出语义」的规则拆到本单。
- **③ 要求**：① 消除预截断：`shift < 0` 时要么在分母侧乘 10 的幂、要么先约分再除、要么走大整数，三条路任选其一，**溢出处必须 `ErrOverflow` 而不是截断**；② `shift > 0` 的溢出 `break` 同样改成 fail-closed（返回 `ErrOverflow`），禁止静默少乘；③ 补表驱动测试：`resultScale` 小于/大于被除数 scale 两侧、`a.value` 接近 `1<<62` 的溢出边界；④ 同步 core `pkg/decimal` 的同一份代码处置，在 `tellurion:docs/TASKS.md` 留同号指针。
- **④ 约束与非目标**：**不改** KT-03 已裁决的 half-away-from-zero 口径；不引 `rounding-mode` 配置项；不动 `Mul` / `Round` / `FromString`；不新增导出符号（`ErrOverflow` 已导出）。
- **⑤ 验收**：`go test ./decimal/... -count=1` 绿；`Div(New(25,1), New(1,0), 0)` = 3、`Div(New(25,1), New(8,0), 0)` = 0（0.3125）；**破坏验证**：把任一路径改回预截断 → 对应测试必须红。
- **⑥ 边界**：只做 `Div` 的定标与溢出路径；舍入口径是 KT-03，`httpclient` 是 KT-04。

---

## 施工序（候选单，开工前必须展开成六字段）

> 这张表是**施工顺序**，不是任务台账：每条在执行前要展开成六字段单（目标/上下文/要求/约束/验收/边界），并写明冻结契约的路径与节号。
> 顺序有依赖，**不要跳级**：跳序的代价是返工（例：X1a 的协议没冻结就写 aggregator，等于把协议写死在一个服务里）。

| 阶段 | 内容 | 上游单号 | 前置 |
|---|---|---|---|
| S1 | ✅ 迁 12 个横切包 + 新增 `db`（connection/migrate/DSN，零表结构）（2026-10-22 完成，正文见归档 QT-05-B） | `QT-05-B` | W0 全绿 |
| S1 | ✅ `kit` 内业务词零命中（`go run ./tools/check_scope` 退出码 0，实测「业务词零命中」） | `QT-05-B` | 同上 |
| S1.5 | 裁决并修掉两条既有缺陷（`decimal` 舍入口径 / `httpclient` 重试语义）——**排在打 tag 前**，否则 `v0.1.0` 之后改就是 breaking | `KT-03` `KT-04` | QT-05-B 完成 |
| S2 | 打 `v0.1.0`；下游换 tag | `QT-07` | S1 全绿 + S1.5 有裁决（改或不改都要留账） |

## 二、跨仓指针（不在本仓，勿双写）

| 单号 | 主题 | 在哪 |
|---|---|---|
| QT-05 | S1 主单（契约层腿 + core 侧依赖切换），本仓腿拆为 QT-05-B | `tellurion-contracts:docs/TASKS.md`（2026-10-21 自 core 迁入；core `tellurion:docs/TASKS.md` 留同号指针） |
| QT-07 | S2：打 `v0.1.0` 并把下游伪版本换成 tag | `tellurion:docs/TASKS.md` |
| QT-17 | agent 数据自持会用到本仓的 `db` | `tellurion-agent:docs/TASKS.md`（2026-10-21 自 core 迁入；core 留同号指针） |

---

## 维护本文件

- **新增单**：写全六字段；六字段里必须显式引用**冻结契约的路径与节号**（契约仓的包路径、设计文档的节号），不许写"参考相关实现"。
- **派工前先证伪前提**：读码核对工单写的前提；不符就停下回报 + 给拆法，不硬做、不扩范围。
- **跨仓引用写 `repo:path`**（如 `tellurion:docs/TASKS.md`），不做可点击相对链接。
- **顺序纪律**：同一波次内撞同一批文件的单按单号串行，一单一 commit。
