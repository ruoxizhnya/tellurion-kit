# CHANGELOG — tellurion-kit

> 本仓的版本账本。**每个发版一节**，标题格式固定 `## [vX.Y.Z] - YYYY-MM-DD`（`tools/check_release_matrix` 按它和 README 兼容矩阵、`git tag` 三方对拍）。
> 未发版的改动记在 `## [Unreleased]` 下，发版时把那一节的内容搬进新版本节。
> 规则口径见 [`README.md`](README.md) 的「版本规则」（源头：`tellurion:docs/TASKS.md` QT-19 ③）。
> breaking 的条目必须写成「破坏性：<旧 → 新>」，光写「改了签名」不算留账。
> ⚠️ 本仓在最底层：一条 breaking 会同时打到 contracts / core / agent，所以「谁被迫改代码」要写在条目里。

## [Unreleased]

尚未发版（**零 tag**；13 包已随 QT-05-B 落码，可 import 但无兼容承诺）。首个版本将是 `v0.1.0`，由 `tellurion:docs/TASKS.md` QT-07 ③-S2-a 在 S2 打；在此之前使用方只能 pin `main`。

### 变更（breaking：导出面语义 —— `decimal.Div` 的舍入口径）

- **KT-03 ②**（2026-10-23）：`decimal.Div` 的舍入口径改为 **half away from zero（四舍五入）**。
  **破坏性**：余数恰好为半个单位时，结果由「向零截断」变为「向远离零进位」——`Div(1, 2, 0)` `0 → 1`、`Div(-1, 2, 0)` `0 → -1`、`Div(1, 8, 2)` `0.12 → 0.13`、`Div(1, 20, 1)` `0.0 → 0.1`。
  **裁决理由**：① 本包 `Round` 已是 half-away-from-zero，`Div` 与它必须同口径，否则「先除再舍」与「先舍再除」没有成文解释；② 向零截断在金额上系统性少算（正数偏小、负数偏大）；③ banker's rounding 虽无偏，但改它要连带改 `Round`，且与会计直觉相反，属另一张单。
  **使用方影响**：受影响面 = 直接调 `Div` 且商恰好落在半个单位的算式（本仓与 core 今天**零生产调用点**，core `pkg/decimal` 仍标「待采用」）。本仓尚未发版，按 README 的 0.x 规则记 breaking、不 bump 版本号；`v0.1.0` 之前 pin `main` 的使用方需自查。判据：`decimal/div_rounding_test.go` 的 `TestDiv_Rounding`（同表覆盖正负两侧、除数正负、奇偶除数）。
  **core 侧处置**：本次**不改** core `pkg/decimal` 的同一份代码（零生产调用、且 core 切本仓依赖属 QT-05 主单）；处置结论已写入 `tellurion:docs/TASKS.md` 的跨仓索引，同号指针为 KT-03。

### 变更（MINOR 级：新增导出面 —— 13 包落码）

- **QT-05-B**（2026-10-22）：S1 工具层落地。13 个包全部有实现，`go build ./...` / `go vet ./...` / `go test ./...` / `go run ./tools/check_scope` / `go run ./tools/check_release_matrix` 全绿（真退出码 0），单独 checkout 即可自证（R1/R2）。
  - **新增导出面（按包）**：`logging`（logger + 脱敏）· `errors`（`AppError` + 9 个 `ErrCode*` + 包装）· `id`（`JobID` / `SubscriptionID`）· `decimal`（定点运算）· `statistics`（纯数学：均值/方差/相关/分位/回归/排名/z 分数）· `httpclient`（超时 + 重试 + JSON 解码）· `httpserver`（`AppError` → status 映射、CORS、gin 模式）· `auth`（token 签发/解析/刷新、bcrypt、claims、中间件）· `authz`（route key → `Rule` + 执法中间件，未登记变更端点 fail-closed）· `observability`（私有 registry + 两个通用族 + `RegisterCollectors` + `/metrics` handler + request-id 贯穿）· `clock`（`Clock` / `VirtualClock` / `LiveClock`，`ModeDeterministic` / `ModeLive`）· `sandbox/staticcheck` + `sandbox/runner` · **`db`（新建）**：`BuildDSN` + 三个 DSN 哨兵错误、`Connect` / `ConnectRedis`、`Migrate(ctx, Execer, statements)`。
  - **与 core 同名包的形状差异**（切依赖时会咬人的地方，逐条在 `docs/SCOPE.md` §5 有裁决与证据）：`id` 不带 `OrderID`；`statistics` 不带 `LogReturns` / `AnnualizeVolatility`；`auth` **不带任何 SQL 与 pgxpool**（身份回读经 `IdentityLoader` 端口，`PublicPaths` 由调用方传入）；`authz` 的 `Policy` 是机制而非 core 的路由表；`observability` 不带 `BacktestDuration` / `LLMTokens` / `NTCall*`（使用方自建 collector 挂进来）；`clock` 的 `ModeBacktest` 改名 `ModeDeterministic`（**iota 顺序不变**，数值语义不动）；`sandbox/wasm` **不在本仓**（WASM 策略宿主留 core）。
  - **`db` 零表结构**（K 裁决）：不含一条 DDL、不含迁移记录表；DDL 文本的唯一真相仍在 core `pkg/storage/postgres.go`。
  - **新登记的第三方依赖**：`golang.org/x/sys v0.46.0` 由间接转直接（`sandbox/runner` 的 Windows Job Object）；`github.com/redis/go-redis/v9 v9.7.3` 已登记进 `AGENTS.md` §4。版本与 core 当前 `go.mod` 对齐，无大版本分叉。
  - **不在本次范围**：没有 tag（S2 = QT-07）；没有 `replace`；没有 `go.work`；core / agent 改 `require` 本仓属 QT-05 主单 ③-③（core 侧）；`auth` 的审计中间件与 `users` / `audit_logs` SQL 留 core。
  - **顺带登记两条既有缺陷（照搬自 core，本次刻意不改）**：`decimal.Div` 的舍入口径与源码注释不符（注释写 banker's rounding，实现是恰好一半向零截断）→ **KT-03**；`httpclient` 带 body 的请求在第二次及以后的尝试发出**空 body**，且 4xx（429 除外）按成功返回 → **KT-04**（两条都已实测取证）。改口径属 MAJOR 语义变更，不在「不改任何业务逻辑」的 QT-05-B 里顺手做。

### 变更（PATCH 级：文档与仓内自查工具，不改导出面）

- **KT-03 ①**（2026-10-23）：`decimal.Div` 的舍入口径取证。新增表驱动测试 `decimal/div_rounding_test.go` 钉住当前行为（恰好一半的正负两侧、除数正负、奇偶除数、`resultScale` 大于/小于被除数 scale），并把源码注释与 `docs/api-schema.md` §1 改成与实现同一句话（当前：余数严格超过半个单位才向远离零进位，恰好一半向零截断）。导出行为零变化。
  取证时另发现两条既有缺陷：`resultScale` 小于被除数 scale 时**先截断被除数再除**（`2.5 / 1` 取整得 `2`），以及 `shift > 0` 溢出时**静默少乘 10 的幂**（错量级）——按「不顺手修未列出语义」拆为 **KT-06**。

- **KT-02**（2026-10-22）：README 补「版本规则表 + 兼容矩阵 + 发版流程 + 拉取前提（`GOPRIVATE`）」；新建本 `CHANGELOG.md`；新建护栏 `tools/check_release_matrix`（矩阵 ↔ CHANGELOG ↔ tag 三件套不一致即红）并接进 CI。导出面零变化 —— 当时本仓没有符号（同日 QT-05-B 落码后已有 13 包，见上节）。
- **KT-01**（2026-10-22）：新建 `docs/SCOPE.md`（13 包逐包准入/禁入 + 业务词表 + 依赖方向 + 九条实测偏差登记）与 `tools/check_scope`（业务词 / 越界 import / 未登记包即红），接进 CI 步骤「包边界核对」。导出面零变化。
