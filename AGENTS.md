# AGENTS — tellurion-kit

> 本文件是**本仓**的唯一入口。判据：一个 agent 只 checkout 本仓，就能知道「怎么构建、边界在哪、依赖哪版契约、活在哪、禁止什么」，**不需要读别的仓的源码**。

## 1. 这是什么仓

横切工具层：logging · errors · id · decimal · statistics · httpclient · httpserver · auth/authz · observability · clock · sandbox · **db**（仅 connection/migrate/DSN）。**零业务概念**。

**状态**：S1 落码已完成（QT-05-B，2026-10-22：13 包全部有实现）／已发版 `v0.1.0`（2026-10-23，QT-07 ③-S2-a）／已发版 `v0.1.1`（2026-10-25，KT-08：新增 `auth/httpadapter` 子包）／已发版 `v0.1.2`（2026-10-25，KT-09：新增 `auth.IssueServiceToken` 服务账号签发）／已发版 `v0.1.3`（2026-10-25，KT-10：`statistics.LogReturns` / `AnnualizeVolatility`、`httpserver.IsLoopbackRemoteAddr`）

## 2. 构建 / 测试（必须单仓可跑）

```bash
go build ./...
go vet ./...
go test ./...
go run ./tools/check_scope          # 包边界：业务词 / 越界 import / SCOPE.md 未登记包 即红（CI 同款，KT-01 护栏）
go run ./tools/check_release_matrix # 版本规则表 + 兼容矩阵 ↔ CHANGELOG ↔ tag（CI 同款，KT-02 护栏）
```

- 前三条是 R1/R2 的判据：**单独 checkout 就能绿**；后两条是本仓的两条随仓护栏（CI 同款步骤）。
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
| `golang.org/x/crypto` | v0.53.0 | bcrypt | `auth` |
| `github.com/jackc/pgx/v5` | v5.7.3 | PG 连接 / 迁移执行机制 | `db` |
| `github.com/redis/go-redis/v9` | v9.7.3 | Redis 连接 | `db` |
| `github.com/prometheus/client_golang` | v1.23.2 | metrics 注册表 | `observability` |
| `golang.org/x/sys` | v0.46.0 | Windows Job Object（进程限额） | `sandbox/runner`（仅 `_windows.go`） |
| `github.com/stretchr/testify` | v1.11.1 | **仅测试**断言 | 各包 `_test.go` |

> 版本与 core 当前 `go.mod` 对齐，避免迁移期同一个库在两个仓里出现两个大版本。
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

**变更触发是硬约束**（完整表见 `docs/README.md` §3）：改表 → 改 `db-schema.md`；改接口 → 改 `api-schema.md`；改链路 → 改 `dataflow.md`；**新增 / 改名 / 删除包，或改某包准入边界 → 改 `docs/SCOPE.md` 那一行**；**发版 → 改兼容矩阵那一行 + `CHANGELOG.md` 那一节**；改分层 → 改 `ARCHITECTURE.md` + 一条 ADR；删档/迁移/换工具 → 一条 ODR。PR 里代码动了而对应文档没动，应被打回。
