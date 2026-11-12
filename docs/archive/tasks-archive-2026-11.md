# TASKS 归档 — tellurion-kit（2026-11）

> **本文件是 `docs/TASKS.md` 的历史归档，进 git、不删除**。归档时间 2026-11-11；收录范围 = 截至 2026-11-10 已收口（✅）的任务原文与当时的施工记录，**逐字搬入、未改写**（原文件 frontmatter 留在新台账；任务标题与正文原样保留）。

> **检索方式**：按任务 ID `grep -n "<ID>" docs/archive/tasks-archive-2026-11.md`；按 commit `grep -n "<sha>" docs/archive/tasks-archive-2026-11.md`；索引表在本节下方，含 ID / 标题 / 收口日期 / commit / 证据链接。

## 检索索引

| ID | 标题 | 收口日期 | commit / 正文所记哈希 | 证据链接 |
|---|---|---|---|---|
| KT-01 | 随仓接口规格载体 `docs/SCOPE.md`（12 包逐包准入/禁入清单） | 2026-10-22 | —（原文未记） | [原文](#KT-01)；[`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md) |
| KT-02 | README 版本规则表 + 兼容矩阵 + CHANGELOG + `GOPRIVATE` 前提 | 2026-10-22 | —（原文未记） | [原文](#KT-02)；[`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md) |
| QT-05-B | S1 工具层落地：12 个包 + `db` 迁入本仓（QT-05 的 kit 腿） | 2026-10-22 | —（原文未记） | [原文](#QT-05-B)；[`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md) |
| KT-03 | `decimal.Div` 的舍入口径与源码注释不一致（裁决口径 + 让注释与实现同一句话） | 2026-10-23 | —（原文未记） | [原文](#KT-03)；[`archive/TASKS-2026-10-23-archived.md`](archive/TASKS-2026-10-23-archived.md) |
| KT-04 | `httpclient` 重试发空 body + 4xx 静默当成功（两条都已实测取证） | 2026-10-23 | —（原文未记） | [原文](#KT-04)；[`archive/TASKS-2026-10-23-archived.md`](archive/TASKS-2026-10-23-archived.md) |
| KT-05 | `errors` 零单测、`logging` 的初始化路径未测（含一个包级可变全局） | 2026-10-23 | —（原文未记） | [原文](#KT-05)；[`archive/TASKS-2026-10-23-archived.md`](archive/TASKS-2026-10-23-archived.md) |
| KT-06 | `decimal.Div` 定标路径的静默精度/量纲损失（KT-03 取证时新发现） | 2026-10-23 | —（原文未记） | [原文](#KT-06) |
| KT-07 | `logging.Init` 的 `level` 参数被无条件覆盖为 Debug（KT-05 补测时取证） | 2026-10-23 | —（原文未记） | [原文](#KT-07) |
| KT-08 | `auth` 的 net/http 适配 + 启动闸（裁-15 的机制腿） | 2026-10-25 | —（原文未记） | [原文](#KT-08)；[`archive/TASKS-2026-10-25-archived.md`](archive/TASKS-2026-10-25-archived.md) |
| KT-09 | `IssueServiceToken`：服务账号 token 签发（裁-17） | 2026-10-25 | —（原文未记） | [原文](#KT-09) |
| KT-10 | 二期切包前置：`statistics` / `httpserver` 的缺口补齐（裁-18） | 2026-10-25 | —（原文未记） | [原文](#KT-10) |
| KT-11 | `service` token 的 subject 被 `int64` 守卫拒成 401（裁决 ⑯ / §五 ⑬ 方向 3） | 2026-10-25 | `11ed852`、`8077c36` | [原文](#KT-11) |
| KT-12 | 把 `-race` 真接进 CI（QT-05-B 登记的两条没做腿之一） | 2026-11-10 | `d9a5c81`、`721faf2` | [原文](#KT-12)；[actions/runs/42788116604](https://github.com/ruoxizhnya/tellurion-kit/actions/runs/42788116604) |
| KT-13 | `logging.Logger` 的包级全局改成可注入（KT-05 ③ / KT-07 ④ 押后的「另一张单」） | 2026-11-10 | `eedeedb`、`65881a6`、`7f35517` | [原文](#KT-13) |
| KIT-DEP-SECURITY-BUMP | 底座依赖安全 bump（pgx GO-2026-5004 / x/net / x/text） | 2026-11-12 | `60bcfba`、`9cc72dc` | [原文](#KIT-DEP-SECURITY-BUMP) |
| KIT-SURFACE-FREEZE | kit 缺导出面冻结护栏（对照 contracts 的 check_surface_freeze） | 2026-11-12 | `6eba77b` | [原文](#KIT-SURFACE-FREEZE) |
| KIT-REQUESTID-STDLIB | RequestID 中间件只有 gin 版，6 服务仓 stdlib mux 装不上 | 2026-11-12 | `c794171` | [原文](#KIT-REQUESTID-STDLIB) |

---

## 原文归档（逐字）

# TASKS — tellurion-kit 未完成项

> 本仓**唯一**任务台账。跨仓任务只在 §二 留指针，**不镜像正文**（双写必然漂移）。
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `✅ 完成（附日期）`
>
> **台账终扫（2026-11-10）**：任务标题带 `⬜` / `🔶` / `⛔` 的 = **0 条**；正文旧标记以标题与最新「收口」块为准。**无开放项**。
> 每张单六字段（① 目标 ② 上下文 ③ 要求 ④ 约束与非目标 ⑤ 验收 ⑥ 边界），判据见 `tellurion:AGENTS.md` §8.5。
> 完成即标 ✅（附日期）并移入 `docs/archive/TASKS-YYYY-MM-DD-archived.md`；本文件只留未完成项。

> **第三轮裁决指针（2026-10-25）**：⑬/⑯ service token 401 = **方向 3 已裁**（按 `claims.Role == "service"` 分流，只豁免 service 角色；用户 token 的 int64 校验原样保留）。已展开成 **`KT-11`**（本文件）并 **✅ 闭合**：CI 六步全绿 → 合 `main`（`8077c36`）→ `v0.1.4` 已发版 → risk/data/analytics 三仓同日升 pin（各自 CI 绿）。**本机无 Go ⇒ 验收证据全部来自 CI**；破坏验证那腿本机跑不了，已如实登记在 KT-11 结案里。正文见 `tellurion:docs/design/qt-decoupling/cross-repo-decision-record-2026-10-25.md` §五。

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

<a id="KT-01"></a>
### KT-01 · 随仓接口规格载体 `docs/SCOPE.md`（12 包逐包准入/禁入清单）[P1] ✅ 完成（2026-10-22，已归档）

> 六字段正文与结案记录（⑤ 验收真退出码、现场破坏验证的还原凭据、22 条夹具腿、词表与覆盖面限制的理由）在
> [`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对 QT-05-B 的直接后果（读 `SCOPE.md` §5 那九条，别照搬）**：清单落的是 **13 行**（12 横切包 + `db`）；
> `pkg/id` 的 `OrderID`、`pkg/clock` 的 `ModeBacktest`、`pkg/observability` 的 `BacktestDuration`/`ObserveBacktest` 撞词表，
> `pkg/auth` 带 `users`/`audit_logs` 的 SQL 撞「零表结构」——这四条**不解决就搬，本仓 CI 第一天就红**；
> 其中 auth 那条是裁决（SQL 留 core 还是 auth 暂不进本仓），不是改名能过关的。

<a id="KT-02"></a>
### KT-02 · README 版本规则表 + 兼容矩阵 + CHANGELOG + `GOPRIVATE` 前提 [P1] ✅ 完成（2026-10-22，已归档）

> 六字段正文与结案记录（③ 逐条落点、⑤ 验收真退出码、现场破坏验证的三条腿与还原凭据、六条口径偏差登记）在
> [`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md)。本文件按规则只留未完成项，这里只留指针。
> **对 QT-05-B / QT-07 的直接后果**：S2 打第一个 tag（`v0.1.0`）那天起，`README.md` 的兼容矩阵必须加 `v0.1` 那一行、
> `CHANGELOG.md` 必须把 `[Unreleased]` 的内容搬进 `## [v0.1.0] - 日期` 节，否则 `go run ./tools/check_release_matrix` 红
> （CI「版本纪律核对」同款）。tag 那条腿今天空跑（0 tag），从第一个 tag 开始才咬人。

<a id="QT-05-B"></a>
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

<a id="KT-03"></a>
### KT-03 · `decimal.Div` 的舍入口径与源码注释不一致（裁决口径 + 让注释与实现同一句话）[P1] ✅ 完成（2026-10-23，已归档）

> 六字段正文与结案记录（① 取证的真退出码、② half-away-from-zero 的裁决理由、破坏验证、core 侧处置）在
> [`archive/TASKS-2026-10-23-archived.md`](archive/TASKS-2026-10-23-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对 KT-06 / QT-07 的直接后果**：口径已冻结为 **half away from zero（四舍五入）**，`docs/api-schema.md` §1 在册；
> KT-06 只动定标与溢出路径，**不得回改口径**；S2 打 `v0.1.0` 时 CHANGELOG 的 breaking 条目已就位（`[Unreleased]` 里第一条）。

<a id="KT-04"></a>
### KT-04 · `httpclient` 重试发空 body + 4xx 静默当成功（两条都已实测取证）[P1] ✅ 完成（2026-10-23，已归档）

> 六字段正文与结案记录（① 复现测试的先红后绿真退出码、③ 4xx → `*AppError` 的裁决与错误码映射、④ 重试不看方法的文档钉死、破坏验证、新增 `httpclient → errors` 边与 ADR 0002、core 侧处置）在
> [`archive/TASKS-2026-10-23-archived.md`](archive/TASKS-2026-10-23-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对 QT-07 / 使用方的直接后果**：4xx 的返回语义已从「`err == nil` 的成功」改为 `nil` + `*errors.AppError`（**breaking**，CHANGELOG `[Unreleased]` 在册），
> `err == nil` 从此只代表 2xx/3xx；重试仍**不看方法**，有副作用的调用须显式 `New(baseURL, timeout, 0)`。口径唯一真相在 `docs/api-schema.md` §3。

<a id="KT-05"></a>
### KT-05 · `errors` 零单测、`logging` 的初始化路径未测（含一个包级可变全局）[P2] ✅ 完成（2026-10-23，已归档）

> 六字段正文与结案记录（`errors` 表驱动覆盖面、`logging` 的 `parseLevel` / `newWriter` 提取与断言、两包覆盖率 100.0%、③ 全局处置裁决、
> 两条破坏验证的真退出码、补测时取证的 `Init` 强制 Debug 缺陷拆为 KT-07）在
> [`archive/TASKS-2026-10-23-archived.md`](archive/TASKS-2026-10-23-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对使用方的直接后果**：`logging.Init` 只能在启动期调一次（它同时写包级 `Logger` 与 zerolog 全局级别），口径钉在 `docs/api-schema.md` §1；
> `errors` 的「`Is` 按分类码相等、`Wrap` 保留 `%w` 链、`ToMap` 只出 `code` / `message` / `op`」从今天起由本仓 CI 直接咬住。导出面零变化。

<a id="KT-06"></a>
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

<a id="KT-07"></a>
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

<a id="KT-08"></a>
### KT-08 · `auth` 的 net/http 适配 + 启动闸（裁-15 的机制腿）[P1] ✅ 完成（2026-10-25，已归档）

> 六字段正文与结案记录（`auth/httpadapter` 的导出面、表驱动 9 + 14 例、破坏验证真退出码与 SHA256 还原凭据、
> `go build` / `go vet` / `go test` / `check_scope` / `check_release_matrix` 的真退出码、`v0.1.1` 发版账）在
> [`archive/TASKS-2026-10-25-archived.md`](archive/TASKS-2026-10-25-archived.md)。本文件按规则只留未完成项，这里只留指针。
>
> **对三条服务腿（risk / analytics / data）的直接后果**：鉴权机制从「各写一份」变为复用 `auth/httpadapter.Guard`；
> 各服务的「端点 × 角色」策略表仍留各自仓（`Policy` 是机制）。`Guard` 必须逐路由挂载（`r.Pattern` 取 route key），
> 漏挂的端点不受保护——服务侧守卫测试应用 `Policy.Entries()` 与路由表双向对拍；未声明的变更路由运行时 403 是 backstop，不是替代测试。

<a id="KT-09"></a>
### KT-09 · `IssueServiceToken`：服务账号 token 签发（裁-17）[P1] ✅ 完成（2026-10-25）

> **结案（2026-10-25）**：新增 `auth/service_token.go` 的 `(*TokenService).IssueServiceToken(subject, role string, ttl time.Duration) (string, error)`
> + 哨兵 `ErrInvalidServiceToken`；只签 `kind=access`（无 refresh 半区，`ParseRefresh` 拒绝同一 token），`Username` 镜像 `subject`（让审计 actor 非空）。
> 判据：`auth/service_token_test.go` 的 `TestIssueServiceToken_RoundTrip`（kind/role/subject/issuer/ttl/no-refresh）、
> `TestIssueServiceToken_Expires`（1ms TTL + 20ms 等待后 `Parse` 报过期）、`TestIssueServiceToken_RejectsBadParams`（6 例：空/空白 subject、空/空白 role、0/负 ttl）、
> `TestIssueServiceToken_DisabledService`（无 secret → `ErrTokensDisabled`）。**破坏验证（真退出码）**：把 `Kind: KindAccess` 改成 `KindRefresh` →
> `TestIssueServiceToken_RoundTrip` 两条断言红（`kind` 不相等 + `ParseRefresh` 不再报错，exit=1）；逐字节还原（`service_token.go` SHA256 `AA4368…837C` 一致）后 exit=0。
> `go test ./auth/... -count=1`、`go run ./tools/check_scope`、`go run ./tools/check_release_matrix` 均退出 0；已发 `v0.1.2`（annotated tag，在 `main` 上）。
>
> ✅ **跨仓前置已落（2026-10-25，`KT-11` / `v0.1.4`）**：`auth.Middleware` 与 `auth/httpadapter.Authenticate` 原要求 subject 可解析为 `int64`，
> 而裁-17 的口径是 `sub=service:<repo>` ⇒ 服务 token 会被这两条守卫 401 拒绝。裁决 ⑯/§五 ⑬ 拍定方向 3 后，
> **仅 `role=service` 免该校验**（用户 token 一字未改），判据是把真 `IssueServiceToken` 签出的 token 送进 `Authenticate`。
> 本单（KT-09）当时按派发只做签发入口、未改既有守卫语义，这条因此挂到 KT-11 才闭上。

- **① 目标**：给运维一个签发 **`service` 角色长期 access token** 的入口，不经过用户登录/刷新链路。
- **② 上下文**：现有 `TokenService.IssueTokens(*User)` 要求 DB 用户（ID/Username），TTL 取 `access_token_ttl`（默认 15m），服务间无法使用。core 侧 `裁-17` 定稿要求 `role=service`、`sub=service:<repo>`、TTL 30d；凭证只经 env 注入，不落库、不进镜像。
- **③ 要求**：① 新增 `IssueServiceToken(subject, role string, ttl time.Duration) (string, error)`（`*TokenService` 方法或包级函数，实现取最小侵入）；② 只签 access token（`kind=access`），**不产生 refresh token**；③ 拒绝空 `subject` / 空 `role` / `ttl<=0`；④ 角色按调用方传入字符串写入 claims，不在 kit 新增枚举；⑤ 文档写明轮换与分发责任在使用方（core QT-71）。
- **④ 约束与非目标**：不改 `IssueTokens` 现有语义；不做吊销列表 / 服务注册表；不把 token 落盘。
- **⑤ 验收**：签出的 token 经 `Parse` 得到 `kind=access` + 传入 role + 传入 subject；TTL 到期后 `Parse` 报过期；空参数报错；`go test ./auth/...` 绿；打 `v0.1.2`。
- **⑥ 边界**：只做签发；分发、部署清单与服务侧角色登记属 core QT-71。

<a id="KT-10"></a>
### KT-10 · 二期切包前置：`statistics` / `httpserver` 的缺口补齐（裁-18）[P2] ✅ 完成（2026-10-25）

> **结案（2026-10-25）**：
> ① `statistics/regression.go` 逐字追加 `LogReturns(prices []float64) []float64` / `AnnualizeVolatility(dailyVol float64, periodsPerYear int) float64`（含私有 `isFinitePositive`），导出签名与 core `pkg/statistics` 逐字一致；`statistics/doc.go` 的 scope note 同步说明「公式型 helper 进来、复合口径仍在外」。
> ② `httpserver/loopback.go` 新增 `IsLoopbackRemoteAddr(remoteAddr string) bool`。**`LoopbackOnly` 不搬**：kit 的依赖白名单**允许** gin，但按派发与本单「只补 core 二期切包缺口」的边界，403 文案与容器装配口径属于应用层策略，core 保留那层 gin 包装。⚠️ **跨仓腿**：core 二期台账需写明「`LoopbackOnly` 留 core、只切 `IsLoopbackRemoteAddr`」——本仓按「只在自己的仓工作」未改 core。
> ③ 判据：`statistics/returns_test.go` 的 `TestLogReturnsParity`（11 例：nil/空/单元素/复合/0 价/负价/NaN/±Inf/混合）、`TestLogReturnsRandomParity`（3 组随机向量，seed=20261025，比 `math.Float64bits`）、`TestAnnualizeVolatilityParity`（25 组 core 黄金 bits，含 `ppy<=0 → 0`、±Inf、负数；NaN 只断言 `IsNaN`）；`httpserver/loopback_test.go` 的 `TestIsLoopbackRemoteAddr`（14 例：回环放行 / 非回环 / 空串 / 不可解析 fail-closed）。
> **破坏验证（真退出码）**：把 `AnnualizeVolatility` 的 `math.Sqrt(float64(periodsPerYear))` 改成 `float64(periodsPerYear)` →
> `TestAnnualizeVolatilityParity` 红（`vol=0.01 ppy=252` 期望 `0x3fc451c23ca192cc` 实得 `0x400428f5c28f5c29`，`ppy=365` 与 `vol=0.2` 两组同红，测试进程非 0 退出）；
> 逐字节还原（`statistics/regression.go` SHA256 `83C515CB472BEBD60B52CAA6B899CD8170DCAD5FF8FC09920AECEC47A117D4C0` 前后一致）后 `go test ./statistics/...` exit=0。
> **验收（真退出码）**：`gofmt -l .` 空；`go build ./...` 0；`go vet ./...` 0；`go test ./... -count=1` 全 ok（`statistics` / `httpserver` 均 ok）；`go run ./tools/check_scope` 0（「业务词零命中、无越界 import、十三包清单齐且覆盖现有包」）；`go run ./tools/check_release_matrix` 0。已发 `v0.1.3`（annotated tag，在 `main` 上，**不与 KT-09 混发**）。

- **① 目标**：让 core 二期能切 `statistics` 与 `httpserver`，先把 kit 缺的函数补齐。
- **② 上下文**：core `statistics` 有 `LogReturns` / `AnnualizeVolatility`，kit 无；core `httpserver` 有 `LoopbackOnly` / `IsLoopbackRemoteAddr`，kit 无。两组是 core 二期切换的唯一硬缺口（导出面已实测；core 侧点数分别 20 / 34）。
- **③ 要求**：① 搬 `LogReturns` / `AnnualizeVolatility` 连同单测（纯函数、无业务配置）；② 评估 `LoopbackOnly`（gin 中间件）与 `IsLoopbackRemoteAddr`：若 kit 的依赖白名单不允许 gin，则只搬 `IsLoopbackRemoteAddr`（`net` 无关），core 保留 gin 包装并在二期台账写明；③ 导出签名与 core 逐字一致；④ 单独打 tag（`v0.1.3` 或二期同批），不与 KT-09 混发。
- **④ 约束与非目标**：不搬 core `metrics.go` 的业务指标；不改 `LoopbackOnly` 不信 XFF 的原则；不改一期五包。
- **⑤ 验收**：kit 版与 core 版同函数在表驱动 + 随机向量上输出逐位一致；core 二期删本地实现后测试绿；`go run ./tools/check_scope` 仍退出 0（无业务词）。
- **⑥ 边界**：只补缺口；core 侧实际切换归 core 二期单，本仓不改 core。

<a id="KT-11"></a>
### KT-11 · `service` token 的 subject 被 `int64` 守卫拒成 401（裁决 ⑯ / §五 ⑬ 方向 3）[P1] ✅ 完成（2026-10-25，已发 `v0.1.4`）

> **结案（2026-10-25）**：
> **代码**（`11ed852`）：`auth/claims.go` 加 `const ServiceRole = "service"`（注释写明不是角色枚举）；`auth/httpadapter/httpadapter.go` 的 `authenticate()` 在 `Kind != access` 之后、`UserIDInt64()` 之前按 `claims.Role == auth.ServiceRole` 分流；`auth/middleware.go` 同口径且**不设 `CtxUserID`**；`httpadapter_test.go` 新增 `TestAuthenticate_ServiceTokenFromIssueServiceToken`（真 `IssueServiceToken("service:risk", …)` 喂进 `Authenticate`，并断非 service + 非数字仍 401）。
> **文档与发版记账**（`8077c36`）：本单六字段 + KT-09 那条「跨仓前置（未落）」改为「已闭合」+ 施工序 S2.7 行 + `CHANGELOG.md` 的 `[v0.1.4]` 节（**该输入类的处置记为破坏性**：原 401 现放行）+ README 兼容矩阵 `v0.1` 行（顺带订正 core 列——原写「尚无 require」，实测 core `go.mod` 已 `require v0.1.1`）+ `docs/api-schema.md` 的 `auth` / `auth/httpadapter` 两行与 §2 服务账号注。
> **验收（全部来自 CI，本机无 Go ⇒ 本地零复跑，这是⑤ 已写明的边界）**：PR #1 rebase 合并 → `main = 8077c36`；PR run 与 main push run `40235174365` **六步全 success**，Test 步含
> `ok github.com/ruoxizhnya/tellurion-kit/auth 0.976s`、`ok github.com/ruoxizhnya/tellurion-kit/auth/httpadapter 0.008s`（新用例就在后者）。
> annotated tag `v0.1.4` 已 push；`https://proxy.golang.org/github.com/ruoxizhnya/tellurion-kit/@v/v0.1.4.info` 返回 `Version=v0.1.4 / Ref=refs/tags/v0.1.4 / Hash=8077c36c006bc497a65a65ad4f322f931b52e0ef`。
> **下游已接**：risk / data / analytics 三仓同日各一笔升 pin PR（各自 CI 全步骤绿 + 合并后 main push run 绿，三条 Build 日志各有 `go: downloading …tellurion-kit v0.1.4`）。
> **破坏验证（如实）**：**本机无法执行**（无 Go，跑不了「改成恒假 → 红」这条）。等价证据是新增用例本身：它把「签出的 token 能被 `Authenticate` 接受」这条断言钉在 CI 里，
> 分流分支一旦被删，`auth/httpadapter` 包必然红——**但这条「红」尚未被实际观测过**，下一个在有 Go 的环境里跑本仓的门禁时应补做（做法：把 `if claims.Role == auth.ServiceRole` 改成恒假，跑 `go test ./auth/httpadapter -run ServiceToken`）。
> **仍未落（不在本仓）**：`JWT_SECRET` 与服务 token 的配发 = core QT-71；「core 真进程 → 服务真进程 → `sub=service:<repo>`」端到端腿一条没跑过。**本单闭的是机制，不是链路。**

> **来由**：KT-09 结案时点名的那条跨仓前置——签发侧允许任意非空 subject（`service:risk`），
> 而 `auth/httpadapter/httpadapter.go:147` 与 `auth/middleware.go` 要求 subject 可解析为 `int64`（`auth/claims.go` 的 `UserIDInt64()`），
> ⇒ 按裁-17 口径签出的服务 token 到消费侧**必然 401**，整条服务间链路走不通。
> 2026-10-25 实测新发现并登记为裁决记录 ⑯，同日拍定方向 3。
>
> **为什么两条既有测试都是绿的**（本单的教训）：kit 的 `service_token_test.go` 只验签发（`NotEmpty` 即过），
> risk 的 `TestPolicy_*` 用 `IssueTokens(1, …)` 的数字 subject ⇒ **两端各自绿、中间那条缝没人测**。本单的判据就是补那条缝。

- **① 目标**：`role=service` 的 access token 能被 `auth/httpadapter.Authenticate` 与 `auth.Middleware` 接受，**且用户 token 的 `int64` 校验不放松**。
- **② 上下文**：`auth/service_token.go`（KT-09）把 subject 原样写进 `sub` 与 `Username`；`auth/claims.go` 的 `UserIDInt64()` 做 `strconv.ParseInt`；
  core 侧同一逻辑已上线并有钉子（`tellurion:pkg/auth/middleware.go:68-81` + `middleware_test.go:164` `TestMiddleware_ServiceTokenNonNumericSubject_OK`），本单是把它移植到 kit 的两条消费守卫。
  裁决正文：`tellurion:docs/design/qt-decoupling/cross-repo-decision-record-2026-10-25.md` §一 ⑯ + §五 ⑬。
- **③ 要求**：① 新增导出常量 `auth.ServiceRole = "service"`，注释写明**不是角色枚举**（KT-09 ③-④ 口径不变）；
  ② `httpadapter.authenticate()` 在 `Kind != access` 之后、`UserIDInt64()` 之前按 `claims.Role == auth.ServiceRole` 分流；
  ③ gin 侧 `Middleware` 同样分流：设 `CtxUsername/CtxRole/CtxClaims`，**不设 `CtxUserID`**（不伪造用户身份，审计侧看到「无用户」而不是用户 0）；
  ④ 判据必须用**真 `IssueServiceToken` 签出的 token** 喂进 `Authenticate`，并断言「非 service 角色 + 非数字 subject 仍 401」；
  ⑤ 记账：`CHANGELOG.md` 的 `v0.1.4` 节 + README 兼容矩阵那一行 + `docs/api-schema.md` 的 `auth` / `auth/httpadapter` 两行与 §2 那条服务账号注记。
- **④ 约束与非目标**：**不改** `IssueServiceToken` 的签名与校验；**不放宽**用户 token 的 `int64` 要求（方向 1 被否）；
  **不**把 subject 强制成数字（方向 2 被否，与「service 主体要可读」及 core `RoleService` 口径冲突）；不建角色枚举、不做吊销；
  不改三个消费仓的 `go.mod`（升 pin 是它们的单）。
- **⑤ 验收**：本仓 CI 六步（`gofmt` / `build` / `vet` / `check_scope` / `check_release_matrix` / `go test ./... -count=1`，`fetch-depth: 0`）全绿；
  新用例 `TestAuthenticate_ServiceTokenFromIssueServiceToken` 通过；破坏验证：把分流分支的 `claims.Role == auth.ServiceRole` 改成恒假 ⇒ 该用例红，还原后绿；
  发版 = `v0.1.4` annotated tag 打在含本节账的 `main` 上并 push。
  ⚠️ **本机无 Go 工具链**（`go` 不在 PATH）⇒ 本地一条都没跑，验收证据只能取 CI 退出码。
- **⑥ 边界**：只做 kit 的两条守卫 + 文档与发版记账；risk / data / analytics 升 pin、orchestration/strategy 首次接 kit 都在别的仓。

---

<a id="KT-12"></a>
### KT-12 · 把 `-race` 真接进 CI（QT-05-B 登记的两条没做腿之一） [P2] ✅（2026-11-10：CI 首跑 `Test -race（覆盖面：全仓所有包）` step 级绿，run `42788116604`）

- **① 目标**：`httpserver` / `auth` / `db` / `observability` 这类并发面在 CI 里被竞态检测器跑过，而不是只被 `-count=1` 跑过。
- **② 上下文**：QT-05-B 结案把两条腿明文记成**没做**——「现场破坏验证与 `-race`」（本文件 `:63-65`，正文在 `archive/TASKS-2026-10-22-archived.md`）；实测 `.github/workflows/ci.yml:49` 今天只有 `go test ./... -count=1`，全仓 CI **无 `-race`**；`AGENTS.md` §2 的三条命令也没带它。⇒ 「并发安全」这句话在本仓目前**没有机械判据**。
- **③ 要求**：① `ci.yml` 增一步 `go test ./... -race -count=1`（Linux runner，gcc 可用）；② 先跑一次全量，把**已经红的**用例逐条登记进 `docs/TASKS.md` 或就地修，**禁止**把用例改名/加 `t.Skip` 来变绿；③ 若一次性全绿不可达，允许先只对 `auth` / `httpserver` / `db` 三包开 `-race` 门禁，并在 CI 步骤名里写明覆盖面（覆盖面必须成文，不许让读者以为全仓已过）；④ `AGENTS.md` §2 同步这条命令与覆盖面注记。
- **④ 约束与非目标**：不改任何被测代码的语义（本单是测试与 CI 配置）；不为跑绿而放宽断言；不引第三方竞态工具；本机 Windows 无 gcc ⇒ `-race` 需 `CGO_ENABLED=1`，**本机跑不动属预期**，结论一律以 CI 输出为准（如实登记，不拿本地绿当证据）。
- **⑤ 验收**：① CI 的 `-race` 步骤**真实执行**（step 级看，不是 job 绿——`skipped` 不算过）；② 破坏验证一条：在 `auth` 或 `httpserver` 里临时加一个无锁并发写 map 的用例 ⇒ `-race` 步骤必须红、`-count=1` 步骤必须保持绿（证明这条腿**抓的是别的东西**，不是重复覆盖）；还原后复绿；③ 覆盖面（全包 or 三包）在 CI 与 `AGENTS.md` 两处口径一致。
- **⑥ 边界**：只动 CI 与本仓测试；core / agent 各自的 `-race` 门禁归它们自己的单。

- **施工记录（2026-11-07）**：
  - **改动面**（2 文件）：`.github/workflows/ci.yml` 在 `Test` 步骤之后新增步骤 **「Test -race（覆盖面：全仓所有包）」** → `go test ./... -race -count=1`；`AGENTS.md` §2 命令块加同一行（带覆盖面与本机限制注记），并把「后两条是本仓的两条随仓护栏」订正为「后三条」。
  - **覆盖面口径**：采用**全仓所有包**，不采用 ③ 允许的三包收窄。理由写进 CI 注释与 `AGENTS.md`：收窄会让读者以为其余包已过竞态检测，那是假账。两处口径一致（⑤-③）。
  - **本地门禁实测**（go1.25.14 windows/amd64，2026-11-07）：`gofmt -l .` 输出空 → exit 0；`go build ./...` → exit 0；`go vet ./...` → exit 0；`go test ./... -count=1` → exit 0（16 包 `ok` / 1 个 `no test files`）；`go run ./tools/check_scope` → exit 0；`go run ./tools/check_release_matrix` → exit 0。
  - **破坏验证（本地可取证半边）**：临时新增 `auth/zz_kt12_race_probe_test.go`（8 goroutine × 5000 次**无锁** `n++`）。`go test ./auth/ -run TestKT12RaceProbeUnsynchronizedCounter -count=1 -v` → `PASS`，exit 0 ⇒ 证实 `-count=1` **抓不到**这类共享写；同一用例加 `-race` → `go: -race requires cgo; enable cgo by setting CGO_ENABLED=1`，exit 2（`go env CGO_ENABLED` = `0`，本机无 gcc）⇒ 「`-race` 下必红」这半边本机原生无法取证；**2026-11-10 已用 Docker（CGO）补齐，见下条**。探针文件已逐字节删除，`git status --porcelain` 为空。
  - **`-race` 红绿腿补齐（2026-11-10，Docker + CGO）**：本机无 C 编译器，改用 `golang:1.25` 镜像挂载本仓、`CGO_ENABLED=1`。① 全量：`go test ./... -race -count=1` → **EXIT=0**，17 个包 `ok`（`auth 7.088s`、`httpclient 16.077s` 等），无 `DATA RACE`。② 破坏探针（8 goroutine × 5000 次无锁 `n++`）：`-count=1` → `ok ... 0.006s`，**NO_RACE_EXIT=0**；同用例 `-race` → `WARNING: DATA RACE`、`--- FAIL ... race detected during execution of test`，**RACE_EXIT=1**。③ 探针文件已删除，`git status --porcelain` 为空。⇒ ⑤② 达成；⑤①（CI step 级首跑）仍缺，故本单维持 🔶。
  - **剩余腿（故记 🔶 非 ✅）**：CI 首跑须 **step 级**确认「Test -race（覆盖面：全仓所有包）」真实执行且绿（`skipped` 不算）；若首跑翻出既有竞态，按 ③ 逐条登记或就地修，**禁止**改名 / `t.Skip` 洗绿。
  - **CI 首跑取证（2026-11-10，主任授权 push）**：`main` 推 `d9a5c81`（6 笔，含 `721faf2` 的 `-race` 步骤）⇒ run **`42788116604`**（[actions/runs/42788116604](https://github.com/ruoxizhnya/tellurion-kit/actions/runs/42788116604)，job `Build & Test` `129734369424`）。**step 级原文**：`Test -race（覆盖面：全仓所有包）` = `success`（非 `skipped`），job 总时长 **2m7s**、结论 `success`；同 job 的 `gofmt` / `Build` / `Vet` / 包边界 / 版本纪律 / `Test` 六步亦全 `success`。⇒ ⑤① 达成；⑤② 已由 2026-11-10 的 Docker 探针腿取证（`-count=1` 绿 / `-race` 红、探针已删）；⑤③ 两处口径同为「全仓所有包」。**本单转 ✅。**

<a id="KT-13"></a>
### KT-13 · `logging.Logger` 的包级全局改成可注入（KT-05 ③ / KT-07 ④ 押后的「另一张单」） [P2] ✅（2026-11-10：`v0.2.0` 已发版并 push tag，七仓升 pin 实测完成）

- **① 目标**：让「取一个 logger」这件事不依赖包级可变全局，使 core/服务能在测试里注入 writer 与 level，而不是各自绕过 kit。
- **② 上下文**：KT-05 ③ 当年裁决「文档钉死」并把可注入化推给后单；KT-07 ④ 原文「不改 `Init` 的签名与 `Logger` **包级全局形状**（改成可注入是**另一张单**）」、KT-07 ⑥ 再写「`Logger` 全局的可注入化与 `errors` 包不在本单」——**那张单此前不存在**（本文件 KT-01~KT-11 无对应条目）。KT-07 已把「level 被无条件覆盖成 Debug」修掉，今天的 `Init(level, format)` 语义正确，本单只动**形状**。
- **③ 要求**：① 先出**决策材料**：注入式 API 的三种候选（返回 `*Logger` / 传 `io.Writer` 参数 / 保留全局但加 `New()` 构造函数）与本仓「移动/改名包即 breaking」的规则如何适用；② 由主任定形状后再改实现；③ 迁移期兼容姿势（保留 `Init` 一个版本周期并标 Deprecated，或一次性 breaking + 下游四仓同批升 pin）在单内定死；④ `docs/SCOPE.md` 的 `logging` 行、`docs/api-schema.md` §1、`README.md` 兼容矩阵、`CHANGELOG.md` 四件套同步；⑤ 下游（core / 六服务）改调用点的单由协调者发，不在本仓。
- **④ 约束与非目标**：**不改** KT-07 已定死的 level 语义（`parseLevel` 映射与未知值回落 info）；不动脱敏、`WithContext`、writer 选择；不引 zerolog 全局钩子改造；**不在材料未拍板前动导出面**（本仓在依赖图底座，一次改动同时打到 contracts / core / agent 的 CI）。
- **⑤ 验收**：① 决策材料落档且主任拍板记录在案（这是本单第一道验收，不是可选项）；② 拍板后：`go build/vet/test ./... -count=1` + `check_scope` + `check_release_matrix` 全绿；③ 注入用例：同一进程内两个 logger 各自写到各自 `io.Writer`，断言互不串（旧形状下这条写不出来，可作为「实现不存在时测试必红」的先写证据）；④ 破坏验证：把注入退回全局 ⇒ 该用例红；⑤ 版本纪律：CHANGELOG 的 MAJOR 节 + 兼容矩阵那一行 + `git tag` 三者一致（`check_release_matrix` 核对）。
- **⑥ 边界**：只做 `logging` 的形状；`errors` 包的同类问题在 KT-05 ⑥ 里被一并押后，若要做请另立单，不与本单混提。

- **施工记录（2026-11-07）**：
  - **① 实测清单（本仓，改动前基线）**：`logging.Logger` 是唯一的包级 logger 全局。**写点 2 处**：`logging/logger.go:12`（声明）、`logging/logger.go:24`（`Init` 赋值）；测试内写 4 处：`logging/logger_test.go:26,146,166,189`。**生产读点 2 处**：`logging/logger.go:29`（`WithContext`）、`httpclient/client.go:38`（`httpclient.New` 构造时快照为默认 logger；`:102,111,122,129,137` 用的是字段不是全局）。**生产调用方**：`logging.Logger` 1 个（`httpclient.New`）；`logging.WithContext` 0 个（仅 3 条测试）；`logging.Init` 本仓 0 个（外部启动期调用）。另记 zerolog 进程级全局：`SetGlobalLevel` 写点 `logging/logger.go:22`（测试 2 处）；`TimeFieldFormat` / `ErrorStackMarshaler` 写点 `logging/logger.go:19-20`。`tools/check_scope/testdata/ok/**` 里的 `logging.Logger` 是夹具自带的同名包，不参与本仓编译。
  - **② 选定形态与理由**：新增 `logging.New(level, format string, out io.Writer) zerolog.Logger`——返回**实例**，sink 用 `out`、级别用 per-logger `.Level(parseLevel(level))`、格式复用 `newWriter`；**不读也不写**包级 `Logger`、**不调** `zerolog.SetGlobalLevel`（`TimeFieldFormat` / `ErrorStackMarshaler` 是与 `Init` 共用的进程级**格式常量**，为保持记录形状一致而同样设置）。保留 `Logger` + `Init` 作为**兼容薄壳**并标 `Deprecated`，写明**下一 MAJOR（v1.0.0）删除**。理由：彻底移除要同时改本仓 `httpclient.New` 与仓外 core/服务的调用点（外部面在本仓不可闭环），且 KT-07 刚修好 `Init` 的级别语义；薄壳 + 新增没有「下游必改项」，破坏面最小且不撒谎。`Init` 函数体**逐字未动** ⇒ 级别过滤 / 字段 / 格式 / sink 行为逐字节不变。
  - **③ 测试**（`logging/logger_test.go`）：`TestNew_WritesToInjectedWriter`（写入传入 writer，且记录形状含 `time` / `caller`）、`TestNew_TwoLoggersWriteToTheirOwnSinks`（同进程两个 logger 各写各的、互不串）、`TestNew_LevelIsScopedToTheInstance`（error / debug 实例各自过滤）、`TestNew_DoesNotMutatePackageLogger`（新入口不写包级全局）、`TestInit_PackageLoggerThinShellStillWorks`（`os.Pipe` 捕获 `Init` 后的 `Logger` 输出，证明薄壳仍工作）。
  - **④ 版本判定**：**MINOR，不是 MAJOR**。依据 README「版本规则」：新增导出符号 = MINOR；`Deprecated` 只是注释级标注，不删符号、不改签名与语义。因此没有 MAJOR 节、本轮不打 tag、兼容矩阵不新增行（未发版）。**迁移指引（下游）**：`logging.Init(level, format)` + `logging.Logger` → `logger := logging.New(level, format, os.Stdout)` 并把 `logger` 显式传给组件；`WithContext` 仍是全局读点，迁移期保留，下一 MAJOR 随薄壳一并处置。落点：`CHANGELOG.md` `[Unreleased]`、`docs/SCOPE.md` §3 `logging` 行、`docs/api-schema.md` §1、`README.md` 兼容矩阵节末段。**下游升级验证：未执行**（本轮不 push、不发版）。
  - **⑤ 破坏验证（本地真退出码，go1.25.14 windows/amd64）**：把 `New` 改成 `_ = level; _ = format; _ = out; return Logger`（忽略传入、退回包级全局）⇒ `go test ./logging/ -run 'TestNew_' -count=1` **4 条全红、exit=1**（首条报 `"" does not contain "\"message\":\"injected\""`）；逐字节还原后 `logging/logger.go` SHA256 前后同为 `28E9ABDDE2706D0C69146CCBCD9173B153E65D4E07DB1FE80DC9E3773C49F25C`，同一命令 **exit=0**。
  - **⑥ 本地门禁实测（同一命令批次）**：`gofmt -l .` 无输出 → **exit 0**；`go build ./...` → **exit 0**；`go vet ./...` → **exit 0**；`go test ./... -count=1` → **exit 0**（16 包 `ok` / 1 个 `no test files`）；`go run ./tools/check_scope` → **exit 0**；`go run ./tools/check_release_matrix` → **exit 0**。
  - **未执行腿（故记 🔶 非 ✅）**：KT-13 ⑤-⑤ 的 `git tag` / `git push`；下游 core / 六服务升 pin 后的编译与行为验证（不在本仓，属协调者发单）；跨仓调用点清单同上。`errors` 包同类问题仍按 ⑥ 另立单。
  - **⑦ 发版与下游升 pin（2026-11-10，主任授权 push + tag）**：
    - **发版**：`eedeedb` 移入 `CHANGELOG [v0.2.0] - 2026-11-10`、README 兼容矩阵新增 `v0.2` 行并补发版段；本地门禁 `gofmt -l .` 空、`go build/vet/test ./... -count=1` 全 0、`check_scope` = 0、`check_release_matrix` = 0（三件套对拍通过）；附注 tag `v0.2.0`（annotated）并 **push**（远端 `refs/tags/v0.2.0 = e3b797c5…`，指向发版提交）。
    - **下游升 pin（七仓，全部 `go get …@v0.2.0`）**：`tellurion`（`65881a6`）、`tellurion-agent`（`7f35517`）、`tellurion-analytics`（`d1be9c9`）、`tellurion-data`（`841ab2d`）、`tellurion-orchestration`（`8acc31b`）、`tellurion-risk`（`9d204d2`）、`tellurion-strategy`（`47193cc`）。**实测**：七仓 `go build ./...` / `go vet ./...` 全 **EXIT=0**；analytics / data / orchestration / risk / strategy 五仓 `go test ./... -count=1` 全 **EXIT=0**；core 全量 `go test ./... -count=1` **EXIT=0**（`pkg/backtest 10.8s` 等全 ok）。
    - **agent 的两条红不是升 pin 引入**：`pkg/ai/client` 的 openapi 对拍红（`/api/v1/trading-calendar`、`/api/v1/listing-windows` 未被 core openapi 登记）。已在 agent `docs/TASKS.md:508` 登记为 **QT-17 ③ 待另立单**的既有红；本轮把 `go.mod/go.sum` stash 回 `v0.1.4` 复跑同包 ⇒ **同样两条红**，升回 `v0.2.0` 后不变，证明与 kit 版本无关。
    - **结论**：⑤-⑤ 三件套一致 + 下游编译/行为验证完成 ⇒ 本单转 ✅；`未执行腿` 里的 tag/push 与下游验证两腿至此结清，`errors` 包仍按 ⑥ 另立单。

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
| S2.6 | ✅ `IssueServiceToken` 服务账号签发，打 `v0.1.2`（裁-17）（2026-10-25 完成，正文见 KT-09） | `KT-09` | core QT-71 开工前 |
| S3（二期） | ✅ 补 `LogReturns` / `AnnualizeVolatility` / `IsLoopbackRemoteAddr`（`LoopbackOnly` 留 core），打 `v0.1.3`（裁-18）（2026-10-25 完成，正文见 KT-10） | `KT-10` | core 二期切包前；不与一期混发 |
| S2.7 | ✅ 修 `service` token 的 subject 被 `int64` 守卫拒成 401（方向 3），打 `v0.1.4`（2026-10-25 完成，正文见 KT-11 结案） | `KT-11` | risk / data / analytics **已同日升 pin**（各仓 CI 绿） |
| S2.8 | ✅ 发版 `v0.2.0`（`logging.New` 可注入）并 push tag；core / agent / analytics / data / orchestration / risk / strategy 七仓同日升 pin 实测（2026-11-10 完成，正文见 KT-13 ⑦） | `KT-13` | KT-12 CI 取证后 |


## 二、跨仓指针（不在本仓，勿双写）

| 单号 | 主题 | 在哪 |
|---|---|---|
| QT-05 | S1 主单（契约层腿 + core 侧依赖切换），本仓腿拆为 QT-05-B | `tellurion-contracts:docs/TASKS.md`（2026-10-21 自 core 迁入；core `tellurion:docs/TASKS.md` 留同号指针） |
| QT-07 | S2：打 `v0.1.0` 并把下游伪版本换成 tag | `tellurion:docs/TASKS.md` |
| QT-17 | agent 数据自持会用到本仓的 `db` | `tellurion-agent:docs/TASKS.md`（2026-10-21 自 core 迁入；core 留同号指针） |
| A 组（三形状包：risk / 编排 / broker DTO） | **本仓零改动、无需立单**：三个形状都不进 kit——`tools/check_scope/main.go:45-60` 的 `businessWords` 含 `order` / `market`，且切词命中驼峰（`OrderRequest` → `[order request]`）；kit 准入判据第 2 条还要求「被 ≥2 个部署单元使用」（`docs/SCOPE.md:18-27`） | 裁决材料在 `tellurion-contracts:docs/design/shape-packages-decision-material.md`（2026-10-26，§3.4 有机械证据） |

---

## 2026-11-12 收口（台账卫生批次）

> 以下三条为 2026-11-12 已收口（✅）任务的原文，自 `docs/TASKS.md` 逐字移入，未改写。

<a id="KIT-DEP-SECURITY-BUMP"></a>
### KIT-DEP-SECURITY-BUMP · 底座依赖安全 bump（pgx GO-2026-5004 / x/net / x/text）[P0] ✅ 2026-11-12

- **① 目标**：把 kit 的 pgx / x/net / x/text 升到安全版并发版，让依赖图底座先闭合漏洞，7 消费仓再跟升（`tellurion:X-DEP-BUMP-CONSUMERS`）。
- **② 上下文**：审查发现 W5-F1（ODR-072）+ 裁决 D1=b（pgx 一步收敛 v5.11.0）。GO-2026-5004（pgx SQLi）fixed v5.9.2，kit `db` 现 pin `pgx v5.7.3`；x/net×4 fixed v0.60.0、x/text×2 需升。kit 是底座，bump 后消费仓跟升才真生效。contracts 是唯一 govulncheck EXIT=0 的仓（零 pgx）。
- **③ 要求**：① pgx 升 v5.11.0（D1=b 一步到位，不逐档）；② x/net→v0.60.0、x/text 升安全版；③ `govulncheck ./...` EXIT=0；④ 发版三件套（CHANGELOG + 兼容矩阵 + tag 推远端）；⑤ AGENTS §4 依赖表同批。
- **④ 约束与非目标**：不手改版本号、不加 replace；不在本单顺手升无关依赖；pgx 大版本收敛的 breaking 面如实登记（底座一次打到 contracts/core/agent）。
- **⑤ 验收**：`govulncheck ./...` EXIT=0（E1）；`go build/vet/test ./... -count=1` 绿；`check_release_matrix` 绿；破坏验证 = pin 回旧版 ⇒ govulncheck 红。
- **⑥ 边界**：kit 依赖与发版；消费仓跟升归 `tellurion:X-DEP-BUMP-CONSUMERS`。
- **owner**：待定
- **重开条件**：无外部前置（A5 已授权立单；发版推远端由主任做）。
- **收口记录（2026-11-12）**：
  - 落地：pgx `v5.7.3 → v5.11.0`、x/net `v0.55.0 → v0.60.0`、x/text `v0.38.0 → v0.42.0`；MVS 连带 x/crypto `v0.53.0 → v0.57.0`、x/sys `v0.46.0 → v0.48.0`、x/sync `v0.21.0 → v0.23.0`。
  - **`go` 指令 1.25.0 → 1.26.11**（x/net v0.60.0 要求 ≥1.26.0；stdlib 29 条 advisory 要到 1.26.11 才齐，故钉补丁版）。版本判定 **MINOR → `v0.4.0`**（依据：工具链底线抬升是消费方可见变化，不是 PATCH 的「无行为修正」；导出面无改动）。
  - 验收（本机真实输出）：`govulncheck ./...` EXIT=0；`go build ./...` / `go vet ./...` / `go test ./... -count=1` 全绿；`go run ./tools/check_scope` 与 `check_release_matrix` 绿。
  - 破坏验证：x/text 降回 v0.38.0（x/net 被 MVS 压到 v0.56.0）⇒ govulncheck EXIT=3（6 条符号级，含 GO-2026-6629 / GO-2026-6617 / GO-2026-5970）；`go.mod` / `go.sum` 逐字节还原（SHA256 一致）后复绿 EXIT=0。单降 pgx 只到 package 级、默认模式 EXIT=0（kit 不调用其漏洞符号），故不作为破坏验证证据。
  - 残留（如实登记）：x/net v0.60.0 两条 module 级（fixed v0.62.0）、x/crypto 一条（GO-2026-5932），govulncheck 判「不调用」；本单版本由裁决指定为 x/net v0.60.0，不擅自再升。
  - 残留的**重开条件**（2026-11-12 裁决：接受现状，不再前推到 x/net v0.62.0；不是「永久不做」）：① 任一消费仓 `go list -m all` 已解析到 x/net ≥ v0.62.0（本仓 pin 成为下游 MVS 下限 / 账实不符）；② 下一次安全轮次（govulncheck 复查 / 季度依赖体检）重新扫出这三条，或任一条从 module 级变为符号级；③ 上游把 GO-2026-5932 的修复发成正式版。触发后一并升到修好版本（x/net ≥ v0.62.0、x/crypto 正式版），判定 MINOR / PATCH 发新 tag 并通知 7 仓二次升 pin；触发前不做任何版本动作。
  - 发版：tag `v0.4.0`（annotated，`main` 上）已推远端；CHANGELOG `v0.4.0` 节 + README 兼容矩阵 `v0.4` 行 + AGENTS §1/§4 同批。

<a id="KIT-SURFACE-FREEZE"></a>
### KIT-SURFACE-FREEZE · kit 缺导出面冻结护栏（对照 contracts 的 check_surface_freeze）[P1] ✅ 2026-11-12

- **① 目标**：给 kit 一个导出面冻结护栏，让「新增/改名/删除导出符号」可机器判定，对齐 contracts 仓已有的 `check_api_surface` + `check_surface_freeze`。
- **② 上下文**：审查发现 W2-F1（ODR-069）。kit `docs/` 无 `API.md`/`api-surface.txt`；`tools/` 只有 `check_scope` + `check_release_matrix`，无导出面冻结。kit 是底座，breaking 面比别仓宽（README「判定补充」：移动/删除/改名包记 breaking），却无机器护栏钉住导出面。
- **③ 要求**：① 建 `docs/api-surface.txt` 基线 + `tools/check_surface_freeze`（对照 contracts 形状）；② 新增导出符号须先登记再更新基线；③ CI 增该护栏步骤；④ AGENTS §6 第四份真相源登记。
- **④ 约束与非目标**：不改现有导出符号（只建护栏）；不照搬 contracts 的业务词表（kit 零业务概念）；护栏破坏验证须真做。
- **⑤ 验收**：`check_surface_freeze` 绿；破坏验证 = 临时加一个未登记导出符号 ⇒ 红，还原复绿。
- **⑥ 边界**：kit 导出面护栏；contracts 侧已有护栏不动。
- **owner**：待定
- **收口记录（2026-11-12）**：
  - 落地：新增 `tools/check_surface_freeze/`（`main.go` + `main_test.go` + `testdata/ok/` 夹具）与 `docs/api-surface.txt`（**17 个库包 / 272 行指纹**；`v0.5.0` 增 `observability.RequestIDHandler` 一行后为 273 行）。指纹格式：`dir/Name(params) results`（函数与导出类型的导出方法，指针接收者记为 `(*T)`）、`dir/Type struct{…}`、`dir/Type interface{…}`、`dir/Type <底层类型>`、`dir/Name <类型|字面量>`；按行确定性排序去重。跳过 `package main`、`*_test.go`、`testdata/`、`vendor/`、`.`/`_` 开头目录。
  - 三个红向：基线缺失 / 空基线 / 一个库包都没扫到（拒绝空绿）；形状漂移先列 `-`（基线有、代码无）再列 `+`（代码有、基线无）。
  - 实现方式：照 contracts 的 `check_surface_freeze` 思路**独立一份**写在 kit 内（go/parser + go/ast，不依赖 go/types，单仓 checkout 可跑），未跨仓引用对方文件。
  - 验收（真实输出）：`go run ./tools/check_surface_freeze` → `✅ 导出面形状与 docs/api-surface.txt 逐行一致`；`go test ./tools/check_surface_freeze -count=1` 9 条用例全 PASS（绿腿 + 新增红 + 删除红 + 改字段红 + 基线缺失红 + 空基线红 + 空扫红 + update 后复绿 + package main 不入基线）。
  - 破坏验证（真做真还原）：① 临时加 `id/surface_probe.go`（`func SurfaceProbe()`）⇒ `+ id/SurfaceProbe()`，EXIT=1；删除后复绿。② 把 `id/id.go` 临时改名为 `id/id.go.off` ⇒ `- id/JobID() string`、`- id/SubscriptionID() string`，EXIT=1；改回后复绿（`git status` 确认无残留）。
  - 接线：CI 在 `check_scope` / `check_release_matrix` 之后新增该步骤；AGENTS §6 登记为「第六份真相源」并写入变更触发；README 目录表 / 快速开始 / 发版流程 / 文档清单同步；`docs/README.md` 文档地图与 §3 变更触发表同步。
  - 说明：本单只建护栏，**未改任何导出符号**，基线即当前导出面的快照。

<a id="KIT-REQUESTID-STDLIB"></a>
### KIT-REQUESTID-STDLIB · RequestID 中间件只有 gin 版，6 服务仓 stdlib mux 装不上 [P1] ✅ 2026-11-12

- **① 目标**：给 kit 一个 **stdlib `net/http` 兼容**的 RequestID 中间件，让 6 个用 stdlib mux 的服务仓能真装上请求 ID 传播。
- **② 上下文**：审查发现 W5-F3（ODR-072）。kit 现 RequestID 中间件是 gin 版，6 服务仓用 stdlib mux 装不上；core `observability/middleware.go:58` 零安装；agent 出站 HTTPTransport 有、入站 RequestIDMiddleware 无。ADR-027:1314-1316 / ADR-029:126-127/:448 / ADR-017:151 的「请求 ID 全链路」断言被证伪（订正归 `tellurion-docs:DOC-ADR-REQUESTID-CORRECTION`）。
- **③ 要求**：① kit 增 stdlib 版 RequestID 中间件（`func(http.Handler) http.Handler`）；② 与 gin 版共用同一 ID 生成/传播口径；③ 发版（只增不改，MINOR）；④ 6 服务仓安装归各仓 + `X-CI-RACE-STATICCHECK` 同批核对。
- **④ 约束与非目标**：不改 gin 版语义；不引入新依赖（用 stdlib + 现有 id 包）；不强制 core/agent 同批安装（各仓自持）。
- **⑤ 验收**：stdlib 中间件单测（入站生成/出站传播/已有 ID 透传）；破坏验证 = 删透传 ⇒ 用例红；`check_scope` 绿（不撞业务词）。
- **⑥ 边界**：kit stdlib 中间件与发版；各仓安装归各仓；ADR 断言订正归 docs 仓。
- **owner**：待定
- **收口记录（2026-11-12）**：
  - 落地：新增 `observability/requestid_http.go` 的 `RequestIDHandler(next http.Handler) http.Handler`（**gin 版 `middleware.go` 一行未改**）。语义与 gin 版逐条对齐：入站 `X-Request-ID` 有则透传、无则 `uuid.NewString()`；token 经既有 `withRequestID` 写进 request context（`RequestIDKey`，取回用 `RequestIDFromContext`）；回镜像响应头；出站复用既有 `HTTPTransport` / `WithRequestIDHeader`（同一 context 即传播）。装配两层幂等（第二次读到第一次的 token）。
  - 依赖：未引入新依赖——UUID 生成走 gin 版已在用的 `github.com/google/uuid`，口径完全一致。
  - 验收（真实输出）：`go test ./observability -count=1` 全绿；新增 4 条用例 PASS——生成（无 header ⇒ 非空且响应头一致）、透传（`upstream-token` 原样）、出站传播（`HTTPTransport` 把 `trace-9` 写到下游 `X-Request-ID`）、64 并发不串号（context 与响应头都等于各自 token）；`go build/vet/test ./... -count=1` 全绿；`check_scope` / `check_release_matrix` / `check_surface_freeze` 绿。
  - 破坏验证：临时把 `rid := r.Header.Get(RequestIDHeader)` 改成 `rid := ""` ⇒ `TestRequestIDHandler_PreservesInboundToken` 红（expected `upstream-token`，actual 为新 UUID）；逐字节还原（`Get-FileHash` 前后一致）后复绿。
  - 发版：`v0.5.0`（MINOR，只增不改）。基线 `docs/api-surface.txt` 增 `observability/RequestIDHandler(http.Handler) http.Handler`，CHANGELOG / 兼容矩阵 `v0.5` / `SCOPE.md` observability 行 / `api-schema.md` 同批更新。
  - 未做（按 ④ 边界）：6 服务仓安装与 `X-CI-RACE-STATICCHECK`、ADR 断言订正分别在各自仓与 docs 仓，本仓不动。
