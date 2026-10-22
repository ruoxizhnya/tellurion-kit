# tellurion-kit

> module: `github.com/ruoxizhnya/tellurion-kit` ｜ 状态：**S1 落码已完成（QT-05-B，2026-10-22：13 包全部有实现）／S2 已发版 `v0.1.0`（QT-07 ③-S2-a，2026-10-23，tag 在 `main` 上）**

横切工具层：logging · errors · id · decimal · statistics · httpclient · httpserver · auth/authz · observability · clock · sandbox · **db**（仅 connection/migrate/DSN）。**零业务概念**。

## 目录

| 路径 | 内容 |
|---|---|
| `logging/` `errors/` `id/` `decimal/` `statistics/` `httpclient/` `httpserver/` `auth/` `authz/` `observability/` `clock/` `sandbox/{staticcheck,runner}/` `db/` | **13 个包（QT-05-B 已落地）**：逐包「收什么 / 不收什么」在 [`docs/SCOPE.md`](docs/SCOPE.md) §3，对外形状在 [`docs/api-schema.md`](docs/api-schema.md) §1 |
| `docs/` | 本仓文档：台账 / 边界 / 决策 / 归档 —— 从这里进：`docs/README.md` |
| `docs/SCOPE.md` | **13 包逐包准入 / 禁入清单** + 业务词表 + 依赖方向 + 实测偏差登记（KT-01） |
| `CHANGELOG.md` | 版本账本：每个发版一节，breaking 逐条写明「旧 → 新」 |
| `tools/check_scope/` | 包边界护栏：业务词 / 越界 import / 未登记包 → 红（`go run ./tools/check_scope`） |
| `tools/check_release_matrix/` | 版本纪律护栏：规则表 / 兼容矩阵 ↔ CHANGELOG ↔ tag 不一致即红（`go run ./tools/check_release_matrix`） |
| `AGENTS.md` | 本仓入口（构建、边界、禁止、依赖版本、拉取前提） |
| `doc.go` | 根包（`package kit`），只有文档注释：它让本仓可被 `go list ./...` 识别，代码一律住上面那些子包 |

## 快速开始

```bash
go build ./...
go vet ./...
go run ./tools/check_scope          # 包边界 ↔ docs/SCOPE.md 逐包清单（CI 同款，KT-01 护栏）
go run ./tools/check_release_matrix # 版本规则表 + 兼容矩阵 ↔ CHANGELOG ↔ tag（KT-02 护栏）
go test ./...
```

本仓**单独 checkout 就能跑通**上面五条。需要别的仓源码才过 → 边界错了，不是加 `replace` 能解决的问题。

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
| `v0.1` | 已发版（2026-10-23，首个版本）；13 包落码（QT-05-B） | 尚无 require（QT-05 contracts 腿落码后接入） | 尚无 require（13 包已迁入本仓；core 侧改 `require` 归 QT-05 主单 ③-③） | 尚无 require（同上） |
| `main`（未发版） | 开发基线，**不做兼容承诺** | 已发版 `v0.1.0` | 已发版 `v0.1.0` | 已发版 `v0.1.0` |
<!-- matrix:end -->

首个版本 `v0.1.0` 已于 **2026-10-23** 打在 `main` 上并 push（`QT-07 ③-S2-a`）。发版前的 13 包落码与两条 breaking（`decimal.Div` 舍入口径、`httpclient` 4xx 返回）都记在 [`CHANGELOG.md`](CHANGELOG.md) 的 `v0.1.0` 节；使用方从 `main` 切到 `v0.1.0` 不需要改代码（内容相同）。

## 发版流程

```bash
# 1. 先在 CHANGELOG.md 把本次改动记进对应版本节（标题格式 ## [vX.Y.Z] - YYYY-MM-DD）
# 2. 再更新 README 兼容矩阵：新的 vX.Y 序列加一行，同序列的补丁版改那一行
# 3. 包有增删改名、或某包准入边界变了 → 同步 docs/SCOPE.md §3 那一行
go run ./tools/check_scope && go run ./tools/check_release_matrix   # 两条都绿才继续
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
- 本仓台账：`docs/TASKS.md`（唯一；跨仓任务只留指针）
- 边界与依赖方向：`docs/ARCHITECTURE.md`
- 版本账本：`CHANGELOG.md`（与兼容矩阵、tag 由 CI 三方对拍）
- 决策记录：`docs/adr/`

## 许可

MIT（见 `LICENSE`）。
