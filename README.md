# tellurion-kit

> 5 分钟了解本服务：[SERVICE.md](SERVICE.md)

> module: `github.com/ruoxizhnya/tellurion-kit` ｜ 状态：**S1 落码已完成（QT-05-B，2026-10-22：13 包全部有实现）／已发版 `v0.1.0`（QT-07 ③-S2-a，2026-10-23）／已发版 `v0.1.1`（KT-08，2026-10-25：新增 `auth/httpadapter` 子包）／已发版 `v0.1.2`（KT-09，2026-10-25：`auth.IssueServiceToken`）／已发版 `v0.1.3`（KT-10，2026-10-25：`statistics.LogReturns` / `AnnualizeVolatility`、`httpserver.IsLoopbackRemoteAddr`）／已发版 `v0.1.4`（KT-11，2026-10-25）／已发版 `v0.2.0`（KT-13，2026-11-10：`logging.New` 可注入 logger）／已发版 `v0.3.0`（K4 前置，2026-11-10：`httpserver.FailCode`）／已发版 `v0.4.0`（KIT-DEP-SECURITY-BUMP，2026-11-12：pgx v5.11.0 / x/net v0.60.0 / x/text v0.42.0，`go` 指令抬到 1.26.11）／已发版 `v0.5.0`（KIT-REQUESTID-STDLIB + KIT-SURFACE-FREEZE，2026-11-12：`observability.RequestIDHandler` stdlib 中间件 + 导出面冻结护栏），tag 均在 `main` 上**

横切工具层：logging · errors · id · decimal · statistics · httpclient · httpserver · auth/authz · observability · clock · sandbox · **db**（仅 connection/migrate/DSN）。**零业务概念**。

## 目录

| 路径 | 内容 |
|---|---|
| `logging/` `errors/` `id/` `decimal/` `statistics/` `httpclient/` `httpserver/` `auth/` `auth/httpadapter/` `authz/` `observability/` `clock/` `sandbox/{staticcheck,runner}/` `db/` | **13 个包（QT-05-B 已落地）+ `auth/httpadapter` 子包（KT-08，2026-10-25）**：逐包「收什么 / 不收什么」在 [`docs/SCOPE.md`](docs/SCOPE.md) §3，对外形状在 [`docs/api-schema.md`](docs/api-schema.md) §1 |
| `docs/` | 本仓文档：台账 / 边界 / 决策 / 归档 —— 从这里进：`docs/README.md` |
| `docs/SCOPE.md` | **13 包逐包准入 / 禁入清单** + 业务词表 + 依赖方向 + 实测偏差登记（KT-01） |
| `docs/api-surface.txt` | **导出面冻结基线**：每个导出符号一行形状指纹（签名 / 字段 / 方法集 / 常量字面量），`KIT-SURFACE-FREEZE` |
| `CHANGELOG.md` | 版本账本：每个发版一节，breaking 逐条写明「旧 → 新」 |
| `tools/check_scope/` | 包边界护栏：业务词 / 越界 import / 未登记包 → 红（`go run ./tools/check_scope`） |
| `tools/check_release_matrix/` | 版本纪律护栏：规则表 / 兼容矩阵 ↔ CHANGELOG ↔ tag 不一致即红（`go run ./tools/check_release_matrix`） |
| `tools/check_surface_freeze/` | 导出面护栏：未登记新增 / 删符号 / 改签名 / 改字段 → 红（`go run ./tools/check_surface_freeze`；有意变更先改账再 `-update`） |
| `AGENTS.md` | 本仓入口（构建、边界、禁止、依赖版本、拉取前提） |
| `doc.go` | 根包（`package kit`），只有文档注释：它让本仓可被 `go list ./...` 识别，代码一律住上面那些子包 |

## 快速开始

```bash
go build ./...
go vet ./...
go run ./tools/check_scope          # 包边界 ↔ docs/SCOPE.md 逐包清单（CI 同款，KT-01 护栏）
go run ./tools/check_release_matrix # 版本规则表 + 兼容矩阵 ↔ CHANGELOG ↔ tag（KT-02 护栏）
go run ./tools/check_surface_freeze # 导出面 ↔ docs/api-surface.txt 逐行指纹（KIT-SURFACE-FREEZE 护栏）
go test ./...
```

本仓**单独 checkout 就能跑通**上面六条。需要别的仓源码才过 → 边界错了，不是加 `replace` 能解决的问题。

`db` 的三条真服务腿挂在 build tag 上，默认不跑（`go test ./...` 不含它们）：

```bash
TEST_DATABASE_DSN="postgres://user:pass@localhost:5432/postgres?sslmode=disable" \
TEST_REDIS_URL="redis://127.0.0.1:6379/0" \
  go test ./db/... -tags integration -count=1 -v
```

未设这两个环境变量时逐条 `t.Skip` 并说明原因（不静默假绿）；迁移那条腿在**总是回滚的事务**里建 TEMP 表，跑完零残留。

## 边界

**可以**：被任意仓 `require`（它是依赖图的底座）；纯技术能力：日志、错误包装、十进制、统计、HTTP 客户端/服务端、鉴权、可观测、时钟、沙箱、DB 连接

**不可以**：任何业务概念：backtest / strategy / order / factor / market / portfolio …；任何表结构（`kit/db` 只做 connection / migrate / DSN 解析）；反向依赖 contracts 或任何应用仓

> 上面两行是**一句话口径**；逐包「收什么 / 不收什么」的完整清单在 [`docs/SCOPE.md`](docs/SCOPE.md)，它是 CI 核对的登记表，不是补充说明。

## 版本规则

规则口径的唯一源头是 `tellurion:docs/TASKS.md` **QT-19 ③**（跨仓决策）；本节是它落在本仓的执行形态，改口径先去那张单，不要在这里另起一套。

本仓在依赖图的**最底层**（`kit ← contracts(含 domain) ← core`，`agent → kit + contracts`），所以一次 breaking 会**同时**打到 contracts、core、agent 三方——版本号不是装饰，是使用方要不要跟着改代码的信号。

<!-- rules:begin -->
| 位 | 何时 bump | 例子（符号名指本仓已落地的导出面；逐包边界见 `docs/SCOPE.md` §3） |
|---|---|---|
| `MAJOR` | **删符号**、**改签名**、**改语义**（默认值、精度与舍入口径、错误分类都算语义） | 删 `Decimal` 的一个方法；改 `Div` 的舍入口径或 `DefaultPoolOptions` 的默认值；`errors.Wrap` 的分类变松；**移动 / 删除 / 改名一个包** |
| `MINOR` | **新增向后兼容的能力**：新包、新类型、新导出符号、struct 新增字段 | `statistics` 新增一个分位函数；`errors` 新增一个 `ErrCode*`；`db` 新增一个 DSN 校验函数 |
| `PATCH` | **文档 / 注释 / 无行为修正** | `SCOPE.md` 措辞订正、CHANGELOG 补条目、README 加前提说明 |
<!-- rules:end -->

**0.x 阶段（现在到 1.0 前）**：允许 breaking，但**必须 bump MINOR + 写 CHANGELOG 条目**，条目里写明「破坏性：<旧 → 新>」。
**1.0 之后**才启用「`deprecated` 的旧接口至少保留 1 个 MINOR 周期」——0.x 阶段没有这条义务，但也没有「随便删」的豁免：删一样要留账。

**判定补充（三个容易误判成 MINOR 的坑）**：

- 给已发布的 **interface 加方法 = breaking**（1.0 之后是 MAJOR；0.x 期间按上面 bump MINOR，但 CHANGELOG 条目必须写成「破坏性」）：所有实现方立刻编不过。本仓的 interface 形状（如 `clock` 的时钟接口）会被 core 与 agent 同时依赖，命中面比一般库更大。
- 给 struct **加字段**通常是 MINOR，但**只要使用方存在无字段名的位置化字面量**（`T{a, b}`）就是 MAJOR —— 不确定就当 breaking 记。
- **包是本仓的对外单位之一**：移动、删除、改名一个包 = breaking，哪怕包内符号一个没变——使用方断的是 import 路径。所以 §「发版流程」第 3 步要把 `docs/SCOPE.md` §3 那一行一起改（`go run ./tools/check_scope` 当场红，不是提醒）。

**不参与 semver 的东西**：内容寻址标识 `strategy_ref` / `artifact_hash` / `run_spec_hash` 自带版本语义（同一个 hash = 同一份内容，变了就是新对象），**不映射到版本号**，也不因它们的变化而 bump。这些标识本身也不住本仓（形状归 contracts）。

**tag 规范**：一律 annotated tag（`git tag -a`），message 首行 `vX.Y.Z: <类型>`，类型取 `feat` / `fix` / `breaking`；tag 必须打在 `main` 上并 push —— 不 push 的 tag 无法 `go get`，等于没发版。

## 兼容矩阵

每一行 = 一个**版本序列** `vX.Y`（同序列的补丁版共用一行，口径来自 QT-19 ③-③），写清该序列同时被哪些上游版本验证过。
使用方在本仓是**三方**而不是两方：contracts 只用 `decimal` / `errors`，core 与 agent 直接 `require` 本仓（QT-05 ③-⑤：`require` 不传递，agent 必须自己 require），所以三列分开记。
**每次发版更新这张表**：新的 `vX.Y` 加一行，同序列的 `vX.Y.Z` 更新那一行 —— 漏更新 → `go run ./tools/check_release_matrix` 当场红（不是提醒）；矩阵里留一个没有账的序列同样红。
「三个仓的 require 版本是否都落在矩阵声明的范围内」是跨仓核对，归 `tellurion:docs/TASKS.md` QT-09 的 CI 步骤（本仓 checkout 不到别的仓，脚本只核对**仓内**三件套的一致性）。

<!-- matrix:begin -->
| kit 版本 | 状态 | 支持的 contracts（`tellurion-contracts`） | 支持的 core（`tellurion`） | 支持的 agent（`tellurion-agent`） |
|---|---|---|---|---|
| `v0.5` | 已发版（`v0.5.0` 2026-11-12 新增 stdlib `net/http` 版 `observability.RequestIDHandler`；同版并入导出面冻结护栏 `KIT-SURFACE-FREEZE`；只增不改） | 尚无 require | 尚未验证升 pin（本版 2026-11-12 刚发版，消费仓跟升归各自台账） | 尚未验证升 pin（同上） |
| `v0.4` | 已发版（`v0.4.0` 2026-11-12 底座依赖安全 bump：pgx `v5.11.0` / x/net `v0.60.0` / x/text `v0.42.0`；`go` 指令 1.26.11；只升依赖、导出面无改动） | 尚无 require | 尚未验证升 pin（本版 2026-11-12 刚发版，消费仓跟升归各自台账） | 尚未验证升 pin（同上） |
| `v0.3` | 已发版（`v0.3.0` 2026-11-10 新增 `httpserver.FailCode` 调用点指定错误码；只增不改） | 尚无 require | 已 require `v0.3.0`（2026-11-10 升 pin 实测 `go.mod`） | 已 require `v0.2.0`（2026-11-10 升 pin 实测 `go.mod`） |
| `v0.2` | 已发版（`v0.2.0` 2026-11-10 新增 `logging.New(level, format, out)` 可注入 logger；旧 `logging.Logger` / `Init` 保留为 `Deprecated` 薄壳，下一 MAJOR 删除） | 尚无 require | 曾被 `v0.2.0` 使用（2026-11-10），已随 `v0.3.0` 前滚 | 已 require `v0.2.0`（2026-11-10 升 pin 实测 `go.mod`） |
| `v0.1` | 已发版（`v0.1.0` 2026-10-23 首个版本；`v0.1.1` 2026-10-25 新增 `auth/httpadapter` 子包；`v0.1.2` 2026-10-25 新增 `auth.IssueServiceToken`；`v0.1.3` 2026-10-25 新增 `statistics.LogReturns` / `AnnualizeVolatility` 与 `httpserver.IsLoopbackRemoteAddr`；`v0.1.4` 2026-10-25 修 `service` token 的 subject 校验）；13 包 + `auth/httpadapter`（QT-05-B / KT-08 / KT-09 / KT-10 / KT-11） | 尚无 require（QT-05 contracts 腿落码后接入） | 已 require `v0.1.1`（QT-73 一期切 `errors` / `decimal` / `clock` / `logging` / `httpclient` 五包，2026-10-25 实测 `go.mod`） | 已 require `v0.1.0`（QT-17 ④ 停写切换后 `pkg/aistore` / `cmd/ai` 走 `kit/db`；2026-10-24 实测 `go.mod`） |
| `main`（未发版） | 开发基线，**不做兼容承诺** | 已发版 `v0.1.0` | 已发版 `v0.1.0` | 已发版 `v0.1.0` |
<!-- matrix:end -->

首个版本 `v0.1.0` 已于 **2026-10-23** 打在 `main` 上并 push（`QT-07 ③-S2-a`）。发版前的 13 包落码与两条 breaking（`decimal.Div` 舍入口径、`httpclient` 4xx 返回）都记在 [`CHANGELOG.md`](CHANGELOG.md) 的 `v0.1.0` 节；使用方从 `main` 切到 `v0.1.0` 不需要改代码（内容相同）。

`v0.1.1` 已于 **2026-10-25** 发版（`KT-08`）：只增不改——新增 `auth/httpadapter` 子包（net/http 适配 + 启动闸 + `Policy`），并把 `KT-06` / `KT-07` 两条 PATCH 修复一并发版。使用方从 `v0.1.0` 升到 `v0.1.1` 不需要改代码。

`v0.1.2` 已于 **2026-10-25** 发版（`KT-09`）：只增不改——`auth.TokenService.IssueServiceToken`（服务账号 access token，无 refresh 半区）+ `ErrInvalidServiceToken`。使用方从 `v0.1.1` 升到 `v0.1.2` 不需要改代码。

`v0.1.3` 已于 **2026-10-25** 发版（`KT-10`）：只增不改——`statistics.LogReturns` / `statistics.AnnualizeVolatility`（导出签名与 core 逐字一致）与 `httpserver.IsLoopbackRemoteAddr`，供 core 二期切包。`LoopbackOnly` gin 中间件仍留 core。使用方从 `v0.1.2` 升到 `v0.1.3` 不需要改代码。

`v0.1.4`（`KT-11`，裁决 ⑯/§五 ⑬ 方向 3）：修「签得出、验不过」——`auth/httpadapter.Authenticate` 与 `auth.Middleware` 对 **`role=service` 的非数字 subject** 不再判 401，用户 token 的 `int64` 校验原样保留。**这一版含行为变更（对一类输入从 401 变为放行）**，使用方升 pin 不需要改代码，但升之前服务间链路是 401、之后才通。risk / data / analytics 三只 pin 由各自仓记账。

`v0.2.0` 已于 **2026-11-10** 发版（`KT-13`）：只增不改——`logging.New(level, format string, out io.Writer) zerolog.Logger` 返回独立实例（writer / 级别 / 过滤都随实例），包级 `Logger` 与 `Init` 原样保留为 `Deprecated` 薄壳（下一 MAJOR 删除）。使用方从 `v0.1.4` 升到 `v0.2.0` **不需要改代码**；想摆脱全局的组件按 CHANGELOG 的迁移指引把 `logging.New(...)` 的返回值显式传下去即可。core / agent / analytics / data / orchestration / risk / strategy 七仓同日升 pin 并实测（见 `docs/TASKS.md` KT-13 施工记录）。

`main` 上的 `KT-13`（2026-11-07 落码）已随 **`v0.2.0`（2026-11-10）** 发版，本段保留为历史：当时未打 tag，现已由上面的 `v0.2.0` 段落接替，兼容矩阵已新增 `v0.2` 行。

`v0.4.0` 已于 **2026-11-12** 发版（`KIT-DEP-SECURITY-BUMP`）：只升依赖、**导出面逐字未动**——pgx `v5.7.3 → v5.11.0`、x/net `v0.55.0 → v0.60.0`、x/text `v0.38.0 → v0.42.0`，`govulncheck ./...` EXIT=0。**唯一需要消费方注意的变化是工具链底线**：x/net v0.60.0 要求 `go >= 1.26.0`，且 stdlib 的安全修复只到 1.26.11，所以本仓 `go` 指令从 `1.25.0` 抬到 `1.26.11`——go.mod 仍写 `go 1.25.x` 的消费方升 pin 后需要 `GOTOOLCHAIN=auto`（默认即 auto）或把本机 / CI 工具链升到 ≥ `1.26.11`，代码不用改。已知残留的 3 条 module 级 advisory（本仓不调用）逐条记在 CHANGELOG 的 `v0.4.0` 节。

`v0.5.0` 已于 **2026-11-12** 发版（`KIT-REQUESTID-STDLIB`）：新增 `observability.RequestIDHandler(next http.Handler) http.Handler`，给用 stdlib `net/http` mux 的服务仓补上 request_id 中间件——入站有 `X-Request-ID` 则透传、无则生成，token 进 request context（`observability.RequestIDFromContext` 取回），出站用同一 context 经 `HTTPTransport` / `WithRequestIDHeader` 传播。**只增不改**：gin 版 `RequestIDMiddleware` 一行未动。同版并入 `KIT-SURFACE-FREEZE` 导出面冻结护栏（`docs/api-surface.txt` + CI 步骤）。

## 发版流程

```bash
# 1. 先在 CHANGELOG.md 把本次改动记进对应版本节（标题格式 ## [vX.Y.Z] - YYYY-MM-DD）
# 2. 再更新 README 兼容矩阵：新的 vX.Y 序列加一行，同序列的补丁版改那一行
# 3. 包有增删改名、或某包准入边界变了 → 同步 docs/SCOPE.md §3 那一行
# 4. 导出面动了（新增 / 改名 / 改签名 / 加删字段）→ 先登记 SCOPE §3 + CHANGELOG，再跑
#    go run ./tools/check_surface_freeze -update 并提交 docs/api-surface.txt 的 diff
go run ./tools/check_scope && go run ./tools/check_surface_freeze && go run ./tools/check_release_matrix   # 三条都绿才继续
git tag -a v0.1.0 -m "v0.1.0: feat <一句话>"
git push origin main && git push origin v0.1.0
```

使用方升级只需 `go get github.com/ruoxizhnya/tellurion-kit@vX.Y.Z` + `go mod tidy`；工具层变更**不走「顺手改一下」**，走这张表 + CHANGELOG 的账。

## 拉取前提（私有仓）

本仓一旦被设为 **private**，`go get` 默认走公共 proxy，**必然失败**——这是 `tellurion:docs/design/qt-decoupling/repo-inventory.md` §1.1/§5 登记的已知风险（现状 12 仓全为 public，转私有前必须先配好下面两条）：

1. 环境与 CI 设 `GOPRIVATE=github.com/ruoxizhnya/*`（让 module 直连 git，不经 proxy/checksum 服务）；
2. 配 git 凭证（HTTPS PAT 或 SSH key）——`GOPRIVATE` 只让 module 直连 git，**认证是另一件事**：只配前者不配后者，`go get` 会在 `git ls-remote https://github.com/ruoxizhnya/tellurion-kit` 那一步 403。

> 本仓是底座，转私有时**先红的是 core 与 agent 的 CI**（它们 `require` 本仓的伪版本或 tag），不是本仓自己——`repo-inventory.md` §5 那句「contracts / kit 变私有后 core 的 CI 与镜像构建会同时红」就是这个意思。

## 文档

- 文档导航：`docs/README.md`
- **包边界清单**：`docs/SCOPE.md`（13 包逐包准入/禁入 + 业务词表；`go run ./tools/check_scope` 核对）
- **导出面基线**：`docs/api-surface.txt`（每符号一行形状指纹；`go run ./tools/check_surface_freeze` 核对）
- 本仓台账：`docs/TASKS.md`（唯一；跨仓任务只留指针）
- 边界与依赖方向：`docs/ARCHITECTURE.md`
- 版本账本：`CHANGELOG.md`（与兼容矩阵、tag 由 CI 三方对拍）
- 决策记录：`docs/adr/`

## 许可

MIT（见 `LICENSE`）。
