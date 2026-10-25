---
status: active
last-verified: 2026-10-25
verified-by: KT-08 结案归档（auth/httpadapter 落码、表驱动用例、破坏验证、v0.1.1 发版账）
---

# TASKS — tellurion-kit 未完成项

> 本仓**唯一**任务台账。跨仓任务只在 §二 留指针，**不镜像正文**（双写必然漂移）。
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `✅ 完成（附日期）`
> 每张单六字段（① 目标 ② 上下文 ③ 要求 ④ 约束与非目标 ⑤ 验收 ⑥ 边界），判据见 `tellurion:AGENTS.md` §8.5。
> 完成即标 ✅（附日期）并移入 `docs/archive/TASKS-YYYY-MM-DD-archived.md`；本文件只留未完成项。

---

## 一、本仓任务

> **本仓已不是占位仓**（2026-10-22 起）：QT-05-B 落码完成，13 包全部有实现；KT-01/KT-02 两张随仓文档交付在此之前已结案。
> 现在台账上的是**落码过程中实测出来的既有缺陷**（KT-06 是 KT-03 取证时新发现；KT-07 是 KT-05 补测时取证。KT-03 / KT-04 / KT-05 已结案）。
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

### KT-05 · `errors` 零单测、`logging` 的初始化路径未测（含一个包级可变全局）[P2] ✅ 完成（2026-10-23，已归档）

> 六字段正文与结案记录（`errors` 表驱动覆盖面、`logging` 的 `parseLevel` / `newWriter` 提取与断言、两包覆盖率 100.0%、③ 全局处置裁决、
> 两条破坏验证的真退出码、补测时取证的 `Init` 强制 Debug 缺陷拆为 KT-07）在
> [`archive/TASKS-2026-10-23-archived.md`](archive/TASKS-2026-10-23-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对使用方的直接后果**：`logging.Init` 只能在启动期调一次（它同时写包级 `Logger` 与 zerolog 全局级别），口径钉在 `docs/api-schema.md` §1；
> `errors` 的「`Is` 按分类码相等、`Wrap` 保留 `%w` 链、`ToMap` 只出 `code` / `message` / `op`」从今天起由本仓 CI 直接咬住。导出面零变化。

### KT-06 · `decimal.Div` 定标路径的静默精度/量纲损失（KT-03 取证时新发现）[P2] ✅ 完成（2026-10-23）

> **结案**：① `shift < 0` 改走分母侧（`denominator × 10^(-shift)`，被除数低位活到舍入步）；② `shift > 0` 的 `break` 改 `ErrOverflow`，且越界守卫从「乘前 > 1<<62」收紧为「乘前 > `math.MaxInt64/10`」（旧守卫乘完仍可越 MaxInt64）；两条定标循环共用 `maxScaled10`。③ 测试：`TestDiv_ResultScaleBelowDividendScale` 订正钉缺陷的期望（2→3）+ 补 2.5/8 与负数例；新增 `TestDiv_OverflowFailsClosed` 五腿（上界/上界内/负上界/旧 break 形状/分母侧）。④ core `pkg/decimal` 已同步同一份 `Div`（其原副本停在「文档写 banker's、实现 tie 向零」的旧分叉，一并对齐 KT-03 口径）并补同款测试，指针留 `tellurion:docs/TASKS.md`。**破坏验证（真退出码）**：突变一（分母侧改回 `numerator /= 10` 预截断）→ `TestDiv_ResultScaleBelowDividendScale` 2 子例红（exit=1）；突变二（`shift > 0` 改回旧 `break`）→ `TestDiv_OverflowFailsClosed` 3 子例红（exit=1）；各还原后 `go test ./decimal/` exit=0。CHANGELOG `[Unreleased]` PATCH 节记账；`docs/api-schema.md` §1 `decimal` 行同步。

- **① 目标**：`Div` 在**任意** `resultScale` 下都给出「精确商 + 一次舍入」的结果，或在无法精确时显式返回 `ErrOverflow`，不再静默给出错量级/多次截断的值。
- **② 上下文**：KT-03 ① 取证时实测（`decimal/div_rounding_test.go` 的 `TestDiv_ResultScaleBelowDividendScale`）：
  `Div(New(25,1), New(1,0), 0)`（2.5 / 1）→ **2**，正确应是 3——`shift < 0` 分支在 `decimal.go` 把被除数逐位 `/10` 截断后再除，信息在舍入前就被丢掉；
  另一条：`shift > 0` 分支在 `numerator` 越过 `1<<62` 时直接 `break`（不再乘 10 的幂），于是商**静默小了 10^k 倍**（量纲错误，不是精度损失）。
  两条都在 KT-03 之前就存在，KT-03 只裁决「恰好一半怎么走」，按「不顺手修未列出语义」的规则拆到本单。
- **③ 要求**：① 消除预截断：`shift < 0` 时要么在分母侧乘 10 的幂、要么先约分再除、要么走大整数，三条路任选其一，**溢出处必须 `ErrOverflow` 而不是截断**；② `shift > 0` 的溢出 `break` 同样改成 fail-closed（返回 `ErrOverflow`），禁止静默少乘；③ 补表驱动测试：`resultScale` 小于/大于被除数 scale 两侧、`a.value` 接近 `1<<62` 的溢出边界；④ 同步 core `pkg/decimal` 的同一份代码处置，在 `tellurion:docs/TASKS.md` 留同号指针。
- **④ 约束与非目标**：**不改** KT-03 已裁决的 half-away-from-zero 口径；不引 `rounding-mode` 配置项；不动 `Mul` / `Round` / `FromString`；不新增导出符号（`ErrOverflow` 已导出）。
- **⑤ 验收**：`go test ./decimal/... -count=1` 绿；`Div(New(25,1), New(1,0), 0)` = 3、`Div(New(25,1), New(8,0), 0)` = 0（0.3125）；**破坏验证**：把任一路径改回预截断 → 对应测试必须红。
- **⑥ 边界**：只做 `Div` 的定标与溢出路径；舍入口径是 KT-03，`httpclient` 是 KT-04。

### KT-07 · `logging.Init` 的 `level` 参数被无条件覆盖为 Debug（KT-05 补测时取证）[P2] ✅ 完成（2026-10-23）

> **结案**：① 删除无条件 `SetGlobalLevel(zerolog.DebugLevel)` 覆盖；② `parseLevel` 映射与未知值回落 info 不动；③ `TestInit_LevelIsAlwaysForcedToDebug` 改写为 `TestInit_LevelArgumentControlsGlobalLevel`（表驱动：debug/info/warn/error/未知→info，断言 `GlobalLevel()` 且对拍 `parseLevel`）。④ CHANGELOG `[Unreleased]` PATCH 节记账（破坏性：生产可观测面——日志从恒 Debug 变按配置）；`docs/api-schema.md` §1 `logging` 行同步。**破坏验证（真退出码）**：把删掉的强制 Debug 加回 → 改写后测试 4 子例红（exit=1；debug 例本就等于 Debug 不红），还原后 `go test ./logging/` exit=0。全仓 build/vet/test + `check_scope` + `check_release_matrix` 全 0。

- **① 目标**：`logging.Init(level, format)` 的 `level` 参数真正决定全局日志级别，不再被启动期硬编码的 Debug 覆盖。
- **② 上下文**：实测（2026-10-23，`logging/logger_test.go` 的 `TestInit_LevelIsAlwaysForcedToDebug`）：`logging/logger.go` 先
  `zerolog.SetGlobalLevel(parseLevel(level))`，紧接着**无条件**再 `zerolog.SetGlobalLevel(zerolog.DebugLevel)`（源码注释 "Force debug level for testing"）。
  于是 `Init("error", "json")` 之后 `zerolog.GlobalLevel()` 仍是 `DebugLevel` —— 生产环境会把 Debug 日志全打出来，`level` 配置完全失效。
  该缺陷照搬自 core `pkg/logging`；KT-05 ④ 明确「不改现有行为」，所以 KT-05 只钉住现状并登记在本单。
- **③ 要求**：① 删除那次无条件覆盖，让 `level` 生效；② 保留 `parseLevel` 的既有映射与「未知值回落 info」口径；
  ③ 把 `TestInit_LevelIsAlwaysForcedToDebug` 改写为「`Init(level, …)` 后 `zerolog.GlobalLevel() == parseLevel(level)`」的表驱动断言
  （该测试当前钉的是缺陷行为，**必须随修复一起改**）；④ CHANGELOG PATCH 节记账（导出行为变化：日志级别从「恒为 Debug」变为「按配置」）。
- **④ 约束与非目标**：不改 `Init` 的签名与 `Logger` 包级全局形状（改成可注入是另一张单，KT-05 ③ 已裁决「文档钉死」）；
  不动脱敏、`WithContext`、writer 选择；不引 zerolog 全局钩子改造。
- **⑤ 验收**：`go test ./logging/... -count=1` 绿；改写后的测试在 `Init("error", …)` 时 `GlobalLevel() == ErrorLevel`、`Init("info", …)` 时为 `InfoLevel`；
  **破坏验证**：把删掉的 `SetGlobalLevel(DebugLevel)` 加回去 → 改写后的测试必须红。全仓 `go build ./...` / `go vet ./...` / `go test ./... -count=1` / 两个 `tools/` 护栏绿。
- **⑥ 边界**：只做 `Init` 的级别生效路径；`Logger` 全局的可注入化与 `errors` 包不在本单。

### KT-08 · `auth` 的 net/http 适配 + 启动闸（裁-15 的机制腿）[P1] ✅ 完成（2026-10-25，已归档）

> 六字段正文与结案记录（`auth/httpadapter` 的导出面、表驱动 9 + 14 例、破坏验证真退出码与 SHA256 还原凭据、
> `go build` / `go vet` / `go test` / `check_scope` / `check_release_matrix` 的真退出码、`v0.1.1` 发版账）在
> [`archive/TASKS-2026-10-25-archived.md`](archive/TASKS-2026-10-25-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对三条服务腿（risk / analytics / data）的直接后果**：鉴权机制从「各写一份」变为复用 `auth/httpadapter.Guard`；
> 各服务的「端点 × 角色」策略表仍留各自仓（`Policy` 是机制）。`Guard` 必须逐路由挂载（`r.Pattern` 取 route key），
> 漏挂的端点不受保护——服务侧守卫测试应用 `Policy.Entries()` 与路由表双向对拍；未声明的变更路由运行时 403 是 backstop，不是替代测试。

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
| S2.5 | ✅ `auth` 的 net/http 适配 + 启动闸，打 `v0.1.1`（裁-15）（2026-10-25 完成，正文见归档 KT-08） | `KT-08` | risk / analytics / data 三腿开工前 |

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
