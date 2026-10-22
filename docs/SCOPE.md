---
status: active
last-verified: 2026-10-22
verified-by: QT-05-B 落地（§5 九条偏差逐条裁决并核对导出面；13 包全部有实现，check_scope 绿）
---

# SCOPE — kit 的逐包准入 / 禁入清单

> 本仓是依赖图的**底座**：`kit ← contracts(含 domain) ← core`，`agent → kit + contracts`。
> 底座一旦被塞进业务概念，它就不再是工具层，而是**新的耦合中心**——所有仓都得跟着它改。
> 所以这份文件不是「参考说明」，是**评审检查项 + CI 护栏的登记表**：`tools/check_scope` 按 §2 的词表与 §4 的依赖面核对，越界即红。
>
> 上游口径（改判要去改这些地方，不要只改本文件）：
> `tellurion:docs/design/qt-decoupling/target-module-topology.md` 的 kit 包边界表（含依赖方向订正 banner）·
> `tellurion:docs/TASKS.md` QT-05 ③-②（12 包 + `db` 的迁移范围）· 本仓 `docs/ARCHITECTURE.md` §3（简表）·
> `tellurion:docs/TASKS.md` QT-05-B（本仓落码执行单）。

## 1. 准入判据（三条同时成立才能进 kit）

| # | 判据 | 谁定的 |
|---|---|---|
| 1 | **零业务概念**——包名、类型名、函数名、变量名里不出现 §2 词表里的任何词 | `target-module-topology.md`「kit 的准入判据」第 1 条 |
| 2 | **被 ≥2 个部署单元使用**（core 与 agent 都要用）——只在单个仓里用的东西留在那个仓 | 同上第 2 条 |
| 3 | **不反向依赖**——kit 不 import contracts、core、agent 的任何包；只依赖第三方与标准库 | 依赖不变量 `kit ← contracts(含 domain) ← core`（同一文档 §50 的 2026-10-20 订正 banner） |

> ⚠️ **「零依赖 ≠ 该进 kit」是这条边界上最容易犯的错**：`fees`（A 股费率）与 `marketimpact`（冲击模型）在 core 里就是零依赖包，
> 但它们是**业务规则**不是工具，裁决是留 core。判据 1 优先于「搬起来方便」。

## 2. 机械核对的词表

`tools/check_scope` 扫的是**声明出来的标识符名**（包级 func / type / const / var / 方法名 / 结构体字段名）与**非标准库 import 路径的段**。
切词规则：标识符按驼峰与下划线切段，**段与词表整词相等**才算命中；import 路径按 `/` 切段同理。

| 词 | 出处 | 为什么算业务 |
|---|---|---|
| `backtest` | 本仓 `AGENTS.md` §3 禁入清单 | 回测是 core 的仪器动作 |
| `strategy` | 同上 | 策略是 core 的运行时对象 |
| `order` | 同上 | 订单是交易域 |
| `factor` | 同上 | 因子是研究域 |
| `market` | 同上 | 行情是数据域 |
| `portfolio` | 同上 | 组合是结算域 |
| `screen` | QT-05-B ③（core 专用类型） | 选股 |
| `sector` | 同上 | 板块 |
| `supply_chain` | 同上 | 产业链 |
| `macro` | 同上 | 宏观序列 |
| `mining` | 同上 | 挖掘（实验域） |
| `ohlcv` | **本单补充** | `market` 的同族数据形状，出现即行情域 |
| `bar` | **本单补充** | K 线 bar 是策略热路径形状（`sandbox/wasm` 的 `BarContext` 就是它） |
| `signal` | **本单补充** | 交易信号（`EmittedSignal` / `hostEmitSignal` 就是它）。标准库 `os/signal` 按「非标准库路径才扫」豁免，不会误红 |

**不扫的东西（覆盖面限制，如实登记）**：

- **注释与字符串字面量**：注释里写「供回测循环使用」是正当的（要说清调用方是谁），扫注释会把文档本身变成违规物。
  代价是注释可以藏业务语义 —— 靠 §3 的逐包禁入清单 + 人工评审兜。
- **`*_test.go`**：库对外的形状由非测试文件决定，测试夹具借用真实形状（如造一个 `OrderID` 来验序列化）不构成对外承诺。
  代价是测试里可以藏形状 —— 同上，靠评审。
- **`tools/` 下的 `package main`**：本仓自查工具自己必须写出这些词才能拦住它们（本文件的词表就在 `main.go` 里）。

## 3. 逐包准入 / 禁入清单（13 包）

<!-- scope:begin -->
| 包 | 准入（收什么） | 禁入（不许收什么） | kit 内允许依赖 | core 现状出处 |
|---|---|---|---|---|
| `logging` | logger 形状、级别、context 注入、敏感值脱敏（API key / 账号 / 通用串） | 业务埋点指标（「回测耗时」「LLM token 数」这类计数属于 observability 的使用方各自定义）；任何 core 类型进字段 | — | `pkg/logging`（`logger.go` `MaskAPIKey` 等 8 个导出，零内部依赖 ✅） |
| `errors` | 错误包装（`Wrap`/`Wrapf`）、分类码（internal / invalid input / not found / conflict / unavailable / timeout）、`errors.Is` 友好形状 | **业务错误码语义**（`ErrFillTimingMismatch` 这类具体判定属 contracts/core）；把某类业务失败降格成通用码 | — | `pkg/errors`（9 个 `ErrCode*`，stdlib-only ✅；`ErrCodeDataQuality` 是边缘名，见 §5 偏差 #6） |
| `id` | 与域无关的 ID 类型与生成器（UUIDv7、job id、订阅 id） | **交易对象 ID**（`OrderID` 这类：它的存在就意味着订单域进了底座）；带业务前缀的编码规则 | — | `pkg/id/order.go`（`OrderID`/`JobID`/`SubscriptionID`）——`OrderID` 命中词表，见 §5 偏差 #1 |
| `decimal` | 定点数运算、舍入模式、精度与 scale、JSON 编解码 | **费率模型 / 价格档位 / 单位换算规则**（`fees` 留 core）；任何「按 A 股规则取整」的语义 | — | `pkg/decimal`（`Decimal` + 12 个构造/运算，stdlib-only ✅，最干净的候选） |
| `statistics` | 纯数学：均值 / 方差 / 标准差 / 相关 / 分位 / 回归 / 排名 / z 分数 | **金融口径**（年化因子、夏普、回撤定义、收益率序列的业务含义）；参数名或文档里携带策略语义 | — | `pkg/statistics`（5 文件；`LogReturns` / `AnnualizeVolatility` 是量化口径，见 §5 偏差 #2） |
| `httpclient` | 客户端形状：请求 / 响应包装、超时与重试的通用机制、JSON 解码、4xx → `*errors.AppError` 的映射 | 具体上游的 endpoint 表（tushare / eastmoney 的接口名与参数是数据服务的事）；鉴权凭证的持有与下发 | logging · errors | `pkg/httpclient/client.go`（→ `pkg/logging` / `pkg/errors`，两包同仓后是**同仓内部 import**，允许 ✅；`→ errors` 由 KT-04 增加） |
| `httpserver` | gin 骨架：CORS、错误中间件、AppError → status 映射、运行模式配置 | **路由清单与业务端点**（`/api/backtest` 这种路径表属 core / agent）；把某个 handler 的语义写进中间件 | errors | `pkg/httpserver`（cors / errors / ginmode，→ `pkg/errors` 同仓内部 ✅） |
| `auth` | 认证机制本身：token 签发 / 解析 / 刷新、bcrypt 成本、claims 形状、中间件 | **表结构 SQL**（`users` / `audit_logs` 的建表与读写）；角色语义与业务权限的绑定关系；任何直连业务表 | — | `pkg/auth`（含 `pgxpool` 与 10 处 `users`/`audit_logs` SQL）——**今天不是 kit 形状**，见 §5 偏差 #3 |
| `authz` | 授权裁决的**机制**：策略条目、角色匹配、按方法是否 mutating 分流、中间件 | **具体路由 → 角色的映射表**（那是各服务的资产）；业务动作名（「下单」「选股」）作为规则条件 | auth | `pkg/authz/policy.go`（→ `pkg/auth` 同仓内部 ✅；注释里是 core 路由清单，见 §5 偏差 #4） |
| `observability` | metrics 注册表、request-id 传递、HTTP transport 埋点、handler 暴露 | **业务指标名**（`backtest_duration_seconds` 这类由使用方在自己仓里注册）；把某条业务链路的口径固化成通用 metric | — | `pkg/observability/metrics.go`（`BacktestDuration` / `ObserveBacktest(strategyName, universe)` / `LLMTokens` / `NTCall*`）——命中词表，见 §5 偏差 #5 |
| `clock` | 时间原语：`Now` / `Since`、虚拟时钟与推进、单调性守卫、时钟模式 | **回测语义**（把「回测确定性」写进类型名或错误文案的形状）；交易日历与交易日推进（那是 core 的仪器域） | — | `pkg/clock/interfaces.go`（`Mode` / `ModeBacktest` / `ErrTimeRewind`）——`ModeBacktest` 命中词表，见 §5 偏差 #7 |
| `sandbox` | `staticcheck`（编译期静态门）+ `runner`（进程隔离、rlimit、超时与 OOM 判定） | **WASM 策略宿主**：`wasm` 子包（`BarContext` / `EmittedSignal` / `hostGetBar` / `OnBarSession`）——它逐 bar 喂数据，必须在引擎进程内（ADR-027 §2.6.1，裁决 ❌ 留 core） | — | `pkg/sandbox/{staticcheck,runner}` ✅ 可搬；`pkg/sandbox/wasm` 4 个文件 import `tellurion/pkg/backtest/contracts`，**不许进本仓** |
| `db` | 连接（pgx pool / redis client）、DSN 解析与占位符校验、**迁移执行机制**（按传入的语句序列跑）、迁移记录表之外的 plumbing | **任何表结构 / DDL 文本 / 列名**（K 裁决：`migrationStatements()` 与所有 `CREATE TABLE` 的唯一真相留在 core `pkg/storage/postgres.go`）；业务表的读写函数；把「迁移」做成带业务语义的框架 | — | **本仓新建**，源出 `pkg/storage/dsn.go`（`BuildDSN` + 三个 DSN 哨兵错误，stdlib-only ✅）与 `postgres.go` 的连接半区（`NewPostgresStoreConnected` / `connectPostgresStore` / `Ping` / `Close`）；⚠️ 「迁移执行」今天在本仓**没有可搬的 runner**，见 §5 偏差 #8 |
<!-- scope:end -->

## 4. 依赖方向与允许的 import 面

```
标准库 · 第三方库
        ↓
      kit（本仓，13 包）
        ↑
contracts（允许：只用 decimal / errors）← core / agent / 各服务（都可以直接 require kit）
```

- **kit 内部**：允许，方向由 §3 的「kit 内允许依赖」列钉住（`httpclient → logging / errors`、`httpserver → errors`、`authz → auth`）。除此之外的包间边是新耦合，加了要在 ADR 里留账。
- **kit → 仓外**：只允许标准库与第三方。**禁止** `tellurion/...`、`.../tellurion-contracts`、`.../tellurion-agent` 与任何应用仓路径。
  不变量是 `kit ← contracts`：contracts 可以依赖 kit 的 `decimal`/`errors`，**反向一条都不许有**（本仓 CI 里没有 core 源码，所以这条红向在编译期就成立，护栏再补一道静态的）。
- 新增第三方依赖前，先在本仓 `AGENTS.md` §4 登记，再改 `go.mod`。

## 5. 实测偏差登记（回源码核对，不是纸面作业）

KT-01 的清单是**目标边界**；下表是 core 今天的形状与之不符之处。九条在 QT-05-B 落码时逐条做了裁决，
最后一列记的是**裁决结果与落地形状**（原来的「对 QT-05-B 的后果」已兑现，保留偏差与证据两列作为历史依据 ——
它们说明「为什么本仓这个包与 core 同名包形状不同」，这是 core 切依赖时最容易踩的地方）。

| # | 偏差 | 实测证据 | 裁决与落地（QT-05-B，2026-10-22） |
|---|---|---|---|
| 1 | `pkg/id` 的文件名与导出符号就是交易域 | `pkg/id/order.go`：`OrderID` / `JobID` / `SubscriptionID` | **已裁决（QT-05-B）**：本仓 `id` 只导出 `JobID` / `SubscriptionID`；`OrderID` 留 core（后续随 contracts 的 domain 值对象走）。不靠改名糊。⚠️ UUIDv7 生成器（`newV7`）**未导出**，所以 §3 那行「通用 UUIDv7」今天只是内部机制——要对外提供得按 MINOR 加一个导出函数（已在 `api-schema.md` §5 登记为"今天没有"）。 |
| 2 | `pkg/statistics` 有量化口径函数 | `LogReturns`、`AnnualizeVolatility`（`pkg/statistics/basic.go` 一族） | **已裁决（QT-05-B）**：本仓 `statistics` 只收纯数学（均值/方差/相关/分位/回归/排名/z 分数）；`LogReturns` / `AnnualizeVolatility` 留 core。 |
| 3 | `pkg/auth` 带表结构与 SQL | `pkg/auth/auth.go:169,204,237,311,320,360,372,389,494,510` 打 `users` / `audit_logs`；`pgxpool` 是构造入参 | **已裁决（QT-05-B）**：SQL 留 core。本仓 `auth` 只收 token 签发/解析/刷新、bcrypt、claims、中间件；不带 pgxpool、不命名任何表。角色是不透明字符串（角色词表与 `CanTrade` 语义留 core）；`Config.PublicPaths` 由调用方传入；刷新经 `IdentityLoader` 端口回读身份。 |
| 4 | `pkg/authz` 的 `Policy` 是按 core 路由清单写的 | `pkg/authz/policy.go` 注释里的「策略 CRUD / 回测 / 下单 / 选股 / 风控」条目 | **已裁决（QT-05-B）**：本仓 `authz` 只收机制——`Policy` 是 route key → `Rule`（Public / AnyAuthenticated / Roles 三档，`NewPolicy` 校验互斥），角色匹配是 any-of；具体路由→角色的表由各服务维护。未登记路由 fail-closed。 |
| 5 | `pkg/observability` 里已经注册了业务指标 | `metrics.go:41` `BacktestDuration`（`backtest_duration_seconds`）、`:132` `ObserveBacktest(strategyName, universe string, ...)`、`LLMTokens`、`NTCallCount/Duration` | **已裁决（QT-05-B）**：本仓 `observability` 只收通用 registry、request-id 中间件/transport、`/metrics` handler；`BacktestDuration` / `LLMTokens` / `NTCall*` 留 core，由使用方注册。 |
| 6 | `pkg/errors` 有 `ErrCodeDataQuality` | `pkg/errors/errors.go` 的 9 个 `ErrCode*` 之一 | **已裁决（QT-05-B）**：保留 `ErrCodeDataQuality`。它是**分类码**而不是业务规则（判据 1 管的是业务概念）；data-service 会用同一个分类，改名反而制造第二个词。 |
| 7 | `pkg/clock` 的模式名与错误文案带着回测语义 | `interfaces.go`：`Mode`、`ModeBacktest`、`ErrTimeRewind`（"…破坏回测确定性"） | **已裁决（QT-05-B）**：`ModeBacktest` → **`ModeDeterministic`**；`ModeLive` 原名保留（它指的是墙钟，不在词表内）。**iota 顺序照搬不变**（`ModeDeterministic=0`、`ModeLive=1`，`clock.go:22` 注明是冻结契约），所以 core 切到本仓时只有名字要改、数值语义不动。错误文案与注释同步中性化（实测 `clock/clock.go` 里「回测」「实盘」「backtest」零命中）。 |
| 8 | `db` 的「迁移执行」今天没有可搬的 runner | core `migration_manager.go` 已于 2026-09-18 删除（零调用方，见 core `migrations/README.md`）；唯一的迁移路径是 `pkg/storage/postgres.go` 的 `migrate()` + `migrationStatements()`（`:150`，约 30 条内联 DDL）；`NewPostgresStore`（`:33`）调 `store.migrate()`，**不能搬** | **已落地（QT-05-B）**：`db` = `BuildDSN` + `Connect` / `ConnectRedis` + `Migrate(ctx, Execer, statements)`；零 DDL、零迁移记录表，空白语句跳过、遇错即停并报位次。 |
| 9 | 四处内部 import 不是障碍 | `httpclient → logging`、`httpserver → errors`、`authz → auth`（以上 QT-05-B）、`httpclient → errors`（KT-04，见 ADR 0002）；测试侧 `pkg/auth/{bootstrap,middleware}_test.go → pkg/testenv`、`pkg/clock/*_test.go` 自 import | **已落地（QT-05-B / KT-04）**：生产侧四条是同仓内部边；搬 auth 时测试改为自带 `TestMain`（gin 模式集中设置），不依赖 core 的 `testenv`。 |

## 6. 护栏

```bash
go run ./tools/check_scope                 # 词表 + 依赖方向 + SCOPE.md 覆盖面
go run ./tools/check_scope -scope docs/SCOPE.md   # 指定清单文件（默认就是这个）
```

红向（`main.go` 顶部注释同步维护）：

1. **业务词命中** —— 非测试文件里的声明标识符或仓外 import 段命中 §2 词表。
2. **越界依赖** —— import 路径命中 `tellurion/`（core）、`tellurion-contracts`、`tellurion-agent` 或其它应用仓前缀。
3. **清单没覆盖** —— 本仓出现了 §3 十三包之外、也不在任何登记包的子树里的 Go 包目录（新增一个没登记的包即红）。
4. **结构红** —— `docs/SCOPE.md` 里找不到 `<!-- scope:begin/end -->`、表头分隔行被删、或 §3 十三包少任何一行（**「把清单删了/缩了来让检查通过」这条路堵死**）。

夹具（`tools/check_scope/testdata/ok/`）是一份齐的清单 + 干净包，并自带两条正腿：同仓内部 import（`sandbox/runner → kit/logging`）
与标准库 `os/signal` —— 都不许被误红。每个红向由用例在临时副本上「改坏 → 断言红」，
其中撞的形状就是 §5 的真实泄漏形状：`OrderID`、`ObserveBacktest`、`BarCount` 字段、`EmitSignal` 接口方法、`import "tellurion/pkg/storage"`。
误红防线也各有断言（`Chain` / `LinkedChain` / `Factorial` 不许命中词表），防止词表日后被「顺手加个词」扩成拦路石。

## 7. 与其他文档的互锁

| 文档 | 分工 |
|---|---|
| `docs/ARCHITECTURE.md` §3 | 只有 5 行的**简表**（定位级），逐包全表在本文件；两处若有出入，以本文件 + CI 为准 |
| `docs/db-schema.md` | 本仓**不拥有任何表**；`db` 包只做连接/DSN/迁移机制（K 裁决），表结构的真相对应 core 的 `pkg/storage/postgres.go` |
| `tellurion-contracts:docs/API.md` | **符号级**导出面清单归 contracts（CT-01）；本文件只管「包该不该收这个东西」，不重复列符号 |
| `README.md` §版本规则/兼容矩阵 + `CHANGELOG.md` | **版本纪律三件套**（KT-02 已落地）：bump 表、按 `vX.Y` 序列建行的矩阵、发版账本，由 `tools/check_release_matrix` 在 CI 与 tag 三方对拍；本文件只管「包该不该收这个东西」，不重复写版本口径 |
| `docs/adr/0001-零业务概念.md` | 「为什么零业务」的决策；本文件是它的执行清单 |
