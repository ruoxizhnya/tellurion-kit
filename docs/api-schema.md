---
status: active
last-verified: 2026-10-25
verified-by: KT-09 / KT-10 落地（`auth.IssueServiceToken`、`statistics.LogReturns` / `AnnualizeVolatility`、`httpserver.IsLoopbackRemoteAddr` 已按导出面核对）
---

# API Schema — tellurion-kit

> **契约地位**：本仓的"API"是 **Go 包级 API**（导出函数 / 类型）。消费方式是 `require` + 版本号。
> 本文件是**本仓对外接口的唯一真相源**。代码与它不一致时，先判断谁错：契约错 → 改本文 + 立 ADR；实现错 → 改代码。
> 跨仓消费方（core / agent / ui / 其他服务）按本文件编码，**不读本仓实现**。
>
> **粒度分工**：本表记的是**每包承诺什么形状**（收什么、不收什么、哪些没提供）。
> 符号级清单不写在这里——那是代码本身 + [`SCOPE.md`](SCOPE.md) §3 的逐包准入/禁入表（`tools/check_scope` 按它核对包边界）。
> 版本承诺（哪个符号属哪个版本、breaking 怎么记）在 `README.md` 的兼容矩阵与 `CHANGELOG.md`。

## 1. 接口面

| 包 | 提供什么 | 典型用途 | 实现状态 |
|---|---|---|---|
| `logging` | 结构化日志、级别、字段约定、**敏感字段脱敏**（`masking.go`） | 全仓统一日志形状 | ✅ 已落地（QT-05-B，2026-10-22）；**`Init` 只能在启动期调一次**——它同时写包级 `Logger` 与 zerolog 全局级别，测试不得并发调（KT-05 ③ 裁决）；`level` 参数真正决定全局级别（KT-07 修复，2026-10-23；未知值回落 info） |
| `errors` | 分类错误 `AppError`（9 个 `ErrCode*`）+ 包装（保留 `%w`）+ 堆栈 | 与 contracts 的三类哨兵配合 | ✅ 已落地（同上） |
| `id` | 与域无关的 ID：`JobID` / `SubscriptionID` | job_id / 订阅 id | ✅ 已落地；**`OrderID` 不在本仓**（交易域，SCOPE §5 #1）；UUIDv7 生成器是包内机制，**未导出通用入口** |
| `decimal` | 定点小数（scale ≤ 18、`Round(scale)`、`Div(a, b, resultScale)`） | 价格与金额（禁止 float 做钱） | ✅ 已落地；`Div` 的舍入口径：**四舍五入（half away from zero）**——余数超过**或恰好等于**半个单位都向远离零进位，与 `Round` 同口径（KT-03 ② 裁决，2026-10-23）。`resultScale` 小于被除数 scale 时定标走**分母侧**（被除数低位数字活到舍入步，`2.5 / 1` 在 scale 0 答 `3`）；任何定标步越界返回 `ErrOverflow`，不静默截断也不静默跳乘（KT-06，2026-10-23）。按 README，舍入口径属 MAJOR 语义 |
| `statistics` | 均值 / 方差 / 标准差 / 相关 / 分位 / 回归 / 排名 / z 分数 + **`LogReturns` / `AnnualizeVolatility`**（公式型，参数与周期数全由调用方显式传入） | 指标计算 | ✅ 已落地；**只收纯数学与公式**——参数显式、无内置年化窗口或收益约定；夏普 / 回撤等复合口径留 core（SCOPE §5 #2，KT-10 / 裁-18 改判） |
| `httpclient` | 带超时 / 重试的客户端 + 统一 JSON 解码 | 服务间调用 | ✅ 已落地；**熔断未提供**（今天没有 breaker，别按"有熔断"编码）；**4xx 返回 `*errors.AppError`**，`err == nil` 只代表 2xx/3xx（口径见 §3） |
| `httpserver` | 统一错误响应（`AppError` + 状态码↔错误码映射）、CORS、gin 模式、**回环来源判据 `IsLoopbackRemoteAddr`**（只看 `RemoteAddr`，不信 `X-Forwarded-For`） | 各服务 HTTP 面；insecure 模式的逐请求来源判定 | ✅ 已落地（同上）；KT-10 只迁判据，**`LoopbackOnly` gin 中间件与 403 文案仍留 core**（应用层装配口径） |
| `auth` | token 签发 / 解析 / 刷新、**服务账号 access token 签发（`IssueServiceToken(subject, role, ttl)`）**、bcrypt、claims、gin 中间件 | API 保护（认证） | ✅ 已落地；**只收机制**——`users` / `audit_logs` 的 SQL 与角色语义留 core，身份回读经 `IdentityLoader` 端口（SCOPE §5 #3）；服务账号签发见 KT-09（2026-10-25） |
| `auth/httpadapter` | net/http 适配（源码不 import gin；依赖闭包因父包 `auth` 的 gin 中间件仍含 gin，见 ADR 0003）：`Authenticate`（Bearer → `Parse` → 拒 `kind != access`，失败 `ErrUnauthenticated` → 401）、`Authorize`（角色 any-of，失败 `ErrForbidden` → 403）、`Guard` 中间件、`ClaimsFromContext`、`DecideStartup` 启动闸、`Policy`（route key → Public / AnyAuthenticated / Roles，**未声明变更路由 403**） | risk / analytics / data 等 net/http 服务的认证与鉴权 | ✅ 已落地（KT-08，2026-10-25）；route key 取 `http.Request.Pattern`，须按路由挂载 `Guard`；路由 → 角色的策略表仍由各服务维护 |
| `authz` | route key → `Rule`（Public / AnyAuthenticated / Roles）+ 执法中间件 | API 保护（鉴权） | ✅ 已落地；**路由→角色的表由各服务维护**，未登记的变更端点 fail-closed（SCOPE §5 #4） |
| `observability` | 私有 Prometheus registry（`Metrics`）+ 两个通用族（`http_client_requests_total` / `cache_hit_ratio`）+ `RegisterCollectors` 扩展点 + `/metrics` handler + request-id 中间件与 `HTTPTransport` 贯穿 | `/metrics`、跨服务追一个请求 | ✅ 已落地；**业务指标族不在本仓**——使用方自建 collector 经 `RegisterCollectors` 挂进来（SCOPE §5 #5）；**`/health` 不在本仓**（见 §5） |
| `clock` | 可注入时钟（`ModeDeterministic` = VirtualClock / `ModeLive` = 墙钟） | 让"纯函数"可测试 | ✅ 已落地；模式名中性化（原 core 的 `ModeBacktest` → `ModeDeterministic`，SCOPE §5 #7）；**取值顺序是冻结契约**（`ModeDeterministic=0`、`ModeLive=1`，禁重排） |
| `sandbox` | `staticcheck`（源码静态检查）+ `runner`（进程隔离与 rlimit / Job Object） | 策略沙箱 | ✅ 已落地；**`wasm` 子包留 core**（那是 WASM host，不是工具） |
| `db` | `BuildDSN` + `Connect` / `ConnectRedis` + `Migrate(ctx, Execer, statements)` | 各仓自己的 store | ✅ 已落地；**零 DDL、零迁移记录表**（SCOPE §5 #8） |

## 2. 错误码与 fail-closed 语义

- 统一用 `errors` 包包装，保留 `%w` 链（可 `errors.Is` / `As`）。
  ⚠️ `AppError.Is` 比的是**分类码**而不是实例身份：`errors.Is(err, errors.New(ErrCodeNotFound, "…"))` 为真，只要两者的 `Code` 相同。
  要按具体原因分支时用 `errors.As(err, &appErr)` 读 `appErr.Code` / `Cause`，不要指望哨兵变量的指针相等。
- 对外错误必须能映射到 HTTP 码（`httpserver` 提供映射表）；**不许把内部错误字符串直接吐给客户端**
  —— `httpserver.Error` / `Wrap` 在 5xx 一律换成通用文案（`genericMessage`），真实原因只进日志（里面可能有连接串 / SQL / 内部路径）。
- 分类 → 状态码由 `httpserver.StatusForAppError` 定，判据只有一条「这件事该谁去修」：

  | `errors` 分类 | HTTP | 含义 |
  |---|---|---|
  | `ErrCodeInvalidInput` | 400 | 请求本身不合法，调用方改请求 |
  | `ErrCodeNotFound` | 404 | 点名的资源不存在 |
  | `ErrCodeDataQuality` / `ErrCodeConflict` | 409 | 请求没错，但当前数据/资源状态不满足前置 |
  | `ErrCodePermission` | 403 | 角色不够 |
  | `ErrCodeRateLimit` | 429 | 限流 |
  | `ErrCodeTimeout` | 504 | 上游超时 |
  | `ErrCodeUnavailable` | 503 | 服务不可用 |
  | 其余（含**不是** `AppError` 的错误） | 500 | 保守默认：未知即服务端问题 |

  ⚠️ 一处已知不对称（`httpserver/errors.go` 内注释在案）：`codeForStatus` 没有 504 分支，
  所以 TIMEOUT → 504 的响应体里 `code` 落到默认的 `"internal"`。修它要动 `code` 的取值集合 = 对外契约变更，须单独裁决。
- 重试只由 `httpclient` 决定，调用方不自己写重试循环。
- **鉴权 fail-closed**：`authz.Middleware` 遇到「变更方法（POST/PUT/PATCH/DELETE）但策略表里没有这一条」时回 **403 并点名该路由**，
  不静默放行；读方法（GET/HEAD/OPTIONS/TRACE）不经策略表。鉴权整体关闭（`Authenticator.Enabled() == false`）时中间件 no-op —— 那是 dev/CI 的 open-access 模式，不是漏洞。
- **net/http 鉴权 fail-closed（`auth/httpadapter`）**：`Authenticate` / `Authorize` 的失败哨兵是 `ErrUnauthenticated`（→401）与 `ErrForbidden`（→403）。
  `Guard` 在鉴权开启时：`Policy` 声明的 Public 路由放行；**未声明的变更方法路由一律 403 并点名 route key**；其余请求必须带有效 access token（401），
  Roles 档再做 any-of 角色判定（403）。未声明的读方法也要求有效 token（fail-closed），Public 读端点须显式登记。鉴权关闭时 `Guard` no-op，与 gin 侧同口径。
- **服务账号 token（`auth.IssueServiceToken`）**：只签 `kind=access`，**不产生 refresh token**，因此不经 `/api/auth/refresh` 续期；
  空 subject / 空 role / `ttl <= 0` 一律 `ErrInvalidServiceToken`；角色是调用方传入的**不透明字符串**（本仓不建角色枚举）。
  轮换周期、分发（env / secret 注入）与吊销策略归调用方（core QT-71）。
  ⚠️ 消费者侧注意：`auth.Middleware` 与 `auth/httpadapter.Authenticate` 仍要求 **subject 可解析为 int64**，
  因此 `sub=service:<repo>` 的服务 token 会被这两条守卫以 401 拒绝；接入服务间调用前须先裁决并放宽该校验（本版未改，见 `docs/TASKS.md` KT-09 的跨仓前置）。
- **DB 配置 fail-closed**：`db.BuildDSN` 对 `${...}` 占位符、空用户名、空密码**当场报错**，不把「连不上」留到第一次查询。
- **迁移 fail-closed**：`db.Migrate` 遇错即停并报出是第几条语句（1 起），不吞错、不续跑；空白语句跳过。

## 3. 幂等性与重试

业务幂等不适用（幂等键属于各服务的 API 契约）。工具层只有一条与幂等相关的承诺：

- `db.Migrate` **不提供**幂等保证，也**不记迁移账**（没有版本表、没有 down、没有执行记录）。
  重复跑同一批语句的结果完全由语句自己决定 —— core 的做法是每条都写成 `CREATE TABLE IF NOT EXISTS` 这类幂等 DDL。
  使用方若传入非幂等语句（`CREATE TABLE` 无 `IF NOT EXISTS`、`INSERT` 无冲突处理），第二次运行就会失败，这是**预期行为**而不是缺陷。
- `httpclient` 的重试**不看方法**：传输错误、读响应体失败、5xx、429 都会重试，POST 也在内 —— 它不保证上游端点幂等，
  有副作用的调用要不要用带重试的 `Client` 由调用方决定（`New(baseURL, timeout, 0)` 可关掉重试；这是唯一开关，**没有**按方法自动降级）。
- `httpclient` 的请求体在**每次尝试都完整重发**（已 marshal 的 `[]byte` 每轮新建 reader），重试不会变成空 body。
- **4xx（429 除外）不是成功**：`Do` 返回 `nil` + `*errors.AppError`，错误码按状态映射 ——
  400 / 422 / 其余未映射 4xx → `INVALID_INPUT`，401 / 403 → `PERMISSION`，404 → `NOT_FOUND`，408 → `TIMEOUT`，409 → `CONFLICT`。
  `message` 以 `unexpected HTTP status <code>` 开头并附上游响应体（截断 200 字），`op` 是 `"<METHOD> <Path>"`。
  调用方按 `if err != nil` 处理即可；`err == nil` 当且仅当 2xx / 3xx（429 在重试耗尽后仍返回「max retries exceeded」错误）。

## 4. 版本与兼容

- 加函数 / 加包 → MINOR
- 改签名 / 改语义 → MAJOR（0.x 期间 bump MINOR 并写 CHANGELOG）
- 删导出符号 → MAJOR（1.0 后先 deprecated 一个 MINOR 周期）

## 5. 不提供什么（边界）

- 不提供业务概念（哪怕是"通用"的 Strategy 接口）
- 不提供表结构 / 迁移 SQL（`db` 只有**执行机制**；DDL 文本的唯一真相在 core `pkg/storage/postgres.go`）
- 不提供 DI 容器 / 应用框架（那是 core 的装配层，QT-50 后外移）
- 不提供测试 mock 框架

**QT-05-B 落地时逐条核实过的"今天没有"**（写在这里是为了让使用方不去猜，缺什么按 MINOR 加）：

| 没有的东西 | 归属 / 理由 |
|---|---|
| `/health` handler | 各服务自己写：健康判据是「这个服务的依赖齐不齐」，属业务装配，不属工具层 |
| 熔断（circuit breaker） | `httpclient` 只有超时 + 重试；要加是 MINOR，但**先要有使用方**（今天没有） |
| 通用 UUIDv7 生成器 | `id` 只导出 `JobID` / `SubscriptionID`；`newV7` 是包内机制。加新 ID 族 = MINOR 加一个函数 |
| `users` / `audit_logs` 的读写与建表 | 留 core（SCOPE §5 #3）；本仓 `auth` 经 `IdentityLoader` 端口回读身份 |
| 具体路由 → 角色的策略表 | 各服务的资产（SCOPE §5 #4）；本仓 `authz` 只提供 `Rule` 机制与执法中间件 |
| `sandbox/wasm`（WASM 策略宿主） | 留 core：它逐 bar 喂数据，必须在引擎进程内（ADR-027 §2.6.1） |
| `OrderID` / `BacktestDuration` 一族 | 分别留 core（SCOPE §5 #1 / #5）：交易域 ID 与业务指标不是工具。`LogReturns` / `AnnualizeVolatility` 已由 KT-10 / 裁-18 迁入本仓（SCOPE §5 #2 改判） |
