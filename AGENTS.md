# AGENTS — tellurion-kit

> 本文件是**本仓**的唯一入口。判据：一个 agent 只 checkout 本仓，就能知道「怎么构建、边界在哪、依赖哪版契约、活在哪、禁止什么」，**不需要读别的仓的源码**。

## 1. 这是什么仓

横切工具层：logging · errors · id · decimal · statistics · httpclient · httpserver · auth/authz · observability · clock · sandbox · **db**（仅 connection/migrate/DSN）。**零业务概念**。

**状态**：S1 落码已完成（QT-05-B，2026-10-22：13 包全部有实现）／已发版 `v0.1.0`（2026-10-23，QT-07 ③-S2-a）／已发版 `v0.1.1`（2026-10-25，KT-08：新增 `auth/httpadapter` 子包）／已发版 `v0.1.2`（2026-10-25，KT-09：新增 `auth.IssueServiceToken` 服务账号签发）／已发版 `v0.1.3`（2026-10-25，KT-10：`statistics.LogReturns` / `AnnualizeVolatility`、`httpserver.IsLoopbackRemoteAddr`）／已发版 **`v0.1.4`**（2026-10-25，KT-11：`auth/httpadapter.Authenticate` 与 `auth.Middleware` 对 `role=service` 的**非数字 subject** 不再判 401——这一版含行为变更，risk / data / analytics 已升 pin；三件套齐：tag `v0.1.4` 已在**远端**（`git ls-remote --tags origin` 复核）、`CHANGELOG.md` 的 `v0.1.4` 节、`README.md` 兼容矩阵 `v0.1` 行）／已发版 `v0.2.0`（2026-11-10，KT-13：新增 `logging.New(level, format, out)` 可注入实例，writer / 级别随实例；包级 `Logger` / `Init` 保留为 `Deprecated` 薄壳，下一 MAJOR 删除）／已发版 **`v0.3.0`**（2026-11-10，K4 前置 / 裁决 ⑲：新增 `httpserver.FailCode(c, status, code, message)` 由调用点指定错误码，供 core 读面复用 analytics 三态 `code`；只增不改。2026-11-11 复核：tag `v0.3.0` 在远端，core 已 require `v0.3.0`，agent 仍 `v0.2.0`）／已发版 **`v0.4.0`**（2026-11-12，KIT-DEP-SECURITY-BUMP：pgx `v5.11.0` / x/net `v0.60.0` / x/text `v0.42.0`，`go` 指令 1.25.0 → 1.26.11；导出面未动，`govulncheck ./...` EXIT=0）／已发版 **`v0.5.0`**（2026-11-12，KIT-REQUESTID-STDLIB：新增 `observability.RequestIDHandler(next http.Handler) http.Handler`，stdlib mux 的 request_id 中间件；gin 版未动、只增不改。同版并入 KIT-SURFACE-FREEZE 导出面冻结护栏）

## 2. 构建 / 测试（必须单仓可跑）

```bash
go build ./...
go vet ./...
go test ./...
go test ./... -race -count=1        # 竞态门禁（CI 同款，KT-12）；覆盖面 = 全仓所有包，需 CGO_ENABLED=1
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out  # 覆盖率门（CI 同款，KIT-COVERAGE-FLOOR，阈值 81%）
go run ./tools/check_scope          # 包边界：业务词 / 越界 import / SCOPE.md 未登记包 即红（CI 同款，KT-01 护栏）
go run ./tools/check_release_matrix # 版本规则表 + 兼容矩阵 ↔ CHANGELOG ↔ tag（CI 同款，KT-02 护栏）
```

- **覆盖率门（KIT-COVERAGE-FLOOR）**：CI 的 Test 步改为带 `-coverprofile` 的覆盖率门，阈值 **81%**（语句加权，`go tool cover -func` 的 `total`）。来源 = 2026-11-13 本机全量实测 `86.4%` 减 5 点，先测后设、不猜数（跨仓口径见 `tellurion:docs/TASKS.md` → `X-COVERAGE-FLOOR`）。低于阈值 CI `exit 1`；**不许**用空壳测试抬覆盖。破坏验证基线：移除 `tools/check_surface_freeze/main_test.go` ⇒ 实测掉到 81% 以下门红，还原 ⇒ `86.4%` 绿。

- 前三条是 R1/R2 的判据：**单独 checkout 就能绿**；后三条是本仓的随仓护栏（CI 同款步骤）。
- **`-race` 的覆盖面与本机限制（KT-12）**：CI 步骤名写作「Test -race（覆盖面：全仓所有包）」，**没有**缩小到 `auth` / `httpserver` / `db` 三包——缩小覆盖面会让读者以为其余包已过竞态检测。`-race` 依赖 cgo，**Windows 无 gcc 时本机跑不动**（`go: -race requires cgo`），这条腿的结论**一律以 CI 输出为准**，不拿本地绿当证据。
- 若某一步需要别的仓的源码目录才过，那是**边界画错了**——停下来修边界，**不要**加 `replace ../xxx` 绕过。
- 本机同时改多个仓时可用 `go.work` 联调，但**它不入库**（`.gitignore` 已挡）。
- **私有化前提**：本仓一旦被设为 private，构建与全部 CI 环境必须配 `GOPRIVATE=github.com/ruoxizhnya/*` **加** git 凭证（PAT / SSH），否则 `go get` 走公共 proxy 必然拉不到（风险注见 `tellurion:docs/design/qt-decoupling/repo-inventory.md` §1.1/§5；现状 12 仓全 public）。本仓在依赖图底座，转私有先红的是 contracts / core / agent 的 CI，不是本仓。细则与发版流程在 `README.md`「拉取前提」「发版流程」。

## 3. 边界（硬约束）

**允许**

- 被任意仓 `require`（它是依赖图的底座）
- 纯技术能力：日志、错误包装、十进制、统计、HTTP 客户端/服务端、鉴权、可观测、时钟、沙箱、DB 连接

**禁止**

- 任何业务概念：backtest / strategy / order / factor / market / portfolio …
- 任何表结构（`kit/db` 只做 connection / migrate / DSN 解析）
- 反向依赖 contracts 或任何应用仓

下面三条对所有仓都成立：

- **禁止** `replace` 指向仓外路径（单仓就不可构建了）。
- **禁止**提交 `go.work` / `go.work.sum`（本地产物；已在 `.gitignore`）。
- 新增第三方依赖前，先在第 4 节登记，再改 `go.mod`；没登记就加依赖 = 绕过评审。

## 4. 依赖的契约版本

| 依赖 | 版本 | 说明 |
|---|---|---|
| `（无内部依赖）` | — | 本仓是依赖图底层；不 `require` 任何内部仓 |

第三方依赖（QT-05-B 迁入时登记；新增先改这张表，再改 `go.mod`）：

| 依赖 | 版本 | 用途 | 引入包 |
|---|---|---|---|
| `github.com/rs/zerolog` | v1.33.0 | 结构化日志 | `logging` / `httpclient` / `httpserver` |
| `github.com/google/uuid` | v1.6.0 | UUIDv7 / 通用 ID | `id` / `observability` |
| `github.com/gin-gonic/gin` | v1.10.0 | HTTP 服务端骨架 | `httpserver` / `observability` / `auth` / `authz` |
| `github.com/spf13/viper` | v1.18.2 | 配置读取（CORS / gin mode） | `httpserver` |
| `github.com/golang-jwt/jwt/v5` | v5.3.1 | token 签发 / 解析 | `auth` |
| `golang.org/x/crypto` | v0.57.0 | bcrypt | `auth` |
| `github.com/jackc/pgx/v5` | v5.11.0 | PG 连接 / 迁移执行机制 | `db` |
| `github.com/redis/go-redis/v9` | v9.7.3 | Redis 连接 | `db` |
| `github.com/prometheus/client_golang` | v1.23.2 | metrics 注册表 | `observability` |
| `golang.org/x/sys` | v0.48.0 | Windows Job Object（进程限额） | `sandbox/runner`（仅 `_windows.go`） |
| `golang.org/x/net` | v0.60.0 | **间接**（HTTP/2 与 IDNA，经 gin / pgx）——安全 bump 主体，见 CHANGELOG `v0.4.0` | 无直接 import |
| `golang.org/x/text` | v0.42.0 | **间接**（编码转换，经 gin / viper / pgx）——安全 bump 主体 | 无直接 import |
| `github.com/stretchr/testify` | v1.11.1 | **仅测试**断言 | 各包 `_test.go` |

> 版本与 core 当前 `go.mod` 对齐，避免迁移期同一个库在两个仓里出现两个大版本。
> **`go` 指令 = 1.26.11**（`v0.4.0` 起）：x/net v0.60.0 要求 `go >= 1.26.0`，而 stdlib 的安全修复要到 1.26.11 才齐；
> 消费方 go.mod 仍是 `go 1.25.x` 时靠 `GOTOOLCHAIN=auto`（默认）自动取 1.26.11，或显式升工具链。升级顺序、残留 advisory 与重开条件见 CHANGELOG `v0.4.0`。
> `github.com/tetratelabs/wazero` 属 core 的 WASM host（`sandbox/wasm` 不搬），**不进本仓**。

> 跨仓依赖**只认版本号**（伪版本或 tag），不认源码路径。契约变更走 `tellurion-contracts` 的发布，不走"顺手改一下"。

## 5. 活在哪、文档在哪

- **本仓台账**：`docs/TASKS.md`（唯一）。**没有单号别开工。**
- 文档导航：`docs/README.md`；边界与依赖方向：`docs/ARCHITECTURE.md`；决策：`docs/adr/`。
- 跨仓背景（只读，不在本仓改）：`tellurion:docs/design/qt-decoupling/repo-topology-and-naming.md`、`tellurion:docs/TASKS.md`。

> 跨仓引用一律写 `repo:path` 纯文本，**不做可点击相对链接**——独立 checkout 时 `../../` 必断，R2 会永远红。

## 6. 文档真相源（改功能前先看这里）

本仓的「功能边界」由三份文件钉住，**它们是唯一真相源，不是参考读物**：

| 维度 | 文件 |
|---|---|
| 数据形状（我拥有哪些表、谁能写） | `docs/db-schema.md` |
| 接口形状（对外暴露什么、错误码、幂等） | `docs/api-schema.md` |
| 流向与顺序（数据怎么走、哪里 fail-closed） | `docs/dataflow.md` |

**本仓还有第四份**：`docs/SCOPE.md` —— **13 包逐包准入/禁入清单** + 业务词表 + 依赖方向 + 实测偏差登记。
它与别的仓的三份不同，因为它是**机器核对**的：出现一个没登记的包、或某个声明标识符撞上词表，`go run ./tools/check_scope` 当场红（KT-01）。
§3 那句「禁止任何业务概念」是口径，`docs/SCOPE.md` §2 那张词表才是可执行边界：**判定撞车时以 SCOPE 与 CI 为准**，改词表要同时改 SCOPE §2 的出处列。

**第五份是版本纪律三件套**：`README.md` 的「版本规则 / 兼容矩阵」+ `CHANGELOG.md` + `git tag`。
同样是机器核对的：发版只写 CHANGELOG 不更新矩阵（或反过来、或 tag 打在没账的版本上），`go run ./tools/check_release_matrix` 当场红（KT-02）。
本仓的 breaking 面比别仓更宽（底座一次改动同时打到 contracts / core / agent），所以「移动 / 删除 / 改名一个包」在本仓**记作 breaking**，见 README「判定补充」。

**第六份是导出面冻结基线**：`docs/api-surface.txt`（`KIT-SURFACE-FREEZE`）—— 每个导出符号一行形状指纹
（函数签名 / 结构体字段 / 接口方法集 / 具名类型底层类型 / 常量字面量），`go run ./tools/check_surface_freeze` 逐行比对。
它和第四份分工不同：`check_scope` 核**包边界**（不该有的包、不该出现的业务词），这一份核**已有导出面的形状有没有动**（改名、改签名、加删字段、改底层类型）。
**红只有一条合法出路**：确认这是有意的 API 变更 → 先在 `CHANGELOG.md` 落账（新增**包**或改包边界另需登记 `docs/SCOPE.md` §3）→ 跑 `go run ./tools/check_surface_freeze -update` 把符号写进 `docs/api-surface.txt` → 把基线的 diff 随代码一起提交评审。符号级登记**只**落在 `docs/api-surface.txt`，SCOPE §3 表是包级准入边界、不列符号。
**不许**为了让 CI 变绿直接 `-update` 而不改账；基线缺失 / 空基线 / 一个库包都没扫到，护栏自身即红（拒绝空绿）。

**导出面登记口径（不照搬 contracts 的 `docs/API.md` 三件套）**：本仓的符号级登记源是 `docs/api-surface.txt`（`tools/check_surface_freeze` 逐行比对形状指纹），发布矩阵由 `tools/check_release_matrix` 承担（`README.md` 兼容矩阵 ↔ `CHANGELOG.md` ↔ tag）；本仓**不设** `docs/API.md`，也**不引入** contracts 的 `tools/check_api_surface` / `check_contract_shape`。理由三条：

- **不造第二份登记源**：`api-surface.txt` 由 `check_surface_freeze -update` 从代码机械生成并逐行冻结；再手写一份符号清单就是第二个真相源，两份迟早漂，且漂的时候只有一头会红。
- **`check_contract_shape` 的判据不适用于本仓**：它禁导出 `*Config` / `*Options` / `*Impl` 形状名、禁任何非测试文件 import 数据库/驱动包——而 `auth.Config`、`db` 包直连 `pgx` 正是本仓的本职能力，照搬会把自己的核心包判红。
- **人读层已由 `docs/api-schema.md` 承担**：它记「每包承诺什么形状、不提供什么」，与 `docs/SCOPE.md` §3（包边界）、`docs/api-surface.txt`（符号形状）三层分工，不缺 contracts `API.md` 那一层。

**变更触发是硬约束**（完整表见 `docs/README.md` §3）：改表 → 改 `db-schema.md`；改接口 → 改 `api-schema.md`；改链路 → 改 `dataflow.md`；**新增 / 改名 / 删除包，或改某包准入边界 → 改 `docs/SCOPE.md` 那一行**；**新增 / 改名 / 改签名 / 加删导出字段 → 改 `docs/api-surface.txt`（`-update`）**；**发版 → 改兼容矩阵那一行 + `CHANGELOG.md` 那一节**；改分层 → 改 `ARCHITECTURE.md` + 一条 ADR；删档/迁移/换工具 → 一条 ODR。PR 里代码动了而对应文档没动，应被打回。
