# TASKS 归档 — tellurion-kit（2026-10-22）

> 规则见 [`README.md`](README.md)：完成即在本仓台账标 ✅，随后整条搬进本目录，**旧代正文不改写**（要订正只在文件头加说明）。
> 本文件是历史证据，不作为当前依据引用；引用前先确认它没被新文档取代。

---

## KT-01 · 随仓接口规格载体 `docs/SCOPE.md`（12 包逐包准入/禁入清单）[P1] ✅ 完成（2026-10-22）

- **① 目标**：本仓每个横切包有随仓、可自证的准入/禁入边界，下游（core/agent）不读 core 实现即可知什么该进 kit、什么不许进。
- **② 上下文**：refactor-plan-review **M6**（报告：`tellurion:docs/design/qt-decoupling/refactor-plan-review-2026-10-20.md`）；上游编排单 `tellurion:docs/TASKS.md` QT-07 ③-S2-c；包边界表源出 `tellurion:docs/design/qt-decoupling/target-module-topology.md`（含其依赖方向 ⚠️ 订正 banner：`kit ← contracts(含 domain) ← core`）。S1 迁移范围见 `tellurion:docs/TASKS.md` QT-05（12 个横切包 + 新增 `db`）。
- **③ 要求**：① 建 `docs/SCOPE.md`：按 QT-05 的 12 包 + `db` 逐包列**准入清单**（该包收什么）与**禁入清单**（业务词、领域语义、core 内部类型禁止入包）；② 每包标注依赖方向约束（kit 不 import contracts 以外的业务层，最终不变量 `kit ← contracts(含 domain) ← core`）；③ 护栏进 CI（AST 扫描）：与「`kit` 内业务词零命中」验收（QT-05）同源，包边界违规即红。
- **④ 约束与非目标**：只写边界清单不写实现细节；不做跨进程 DTO（那是 `tellurion:docs/openapi.yaml` 生成物、归 core，见 `tellurion:docs/design/qt-decoupling/repo-inventory.md` §3）；不含版本规则（那是 KT-02）。
- **⑤ 验收**：`docs/SCOPE.md` 覆盖 12 包 + `db` 逐包成表；AST 护栏接入 CI；**破坏验证**：故意往任一包塞一个业务词/越界 import → 护栏红（留证后撤回）。
- **⑥ 边界**：只做包准入/禁入载体与护栏；导出面接口规格（符号级）归 `tellurion-contracts:docs/TASKS.md` CT-01；版本策略归 KT-02。

### 结案记录（2026-10-22）

**交付物**

| 文件 | 内容 |
|---|---|
| `docs/SCOPE.md` | 新建。§1 三条准入判据（含「零依赖 ≠ 该进 kit」这条反例：`fees`/`marketimpact` 零依赖但裁决留 core）· §2 可执行词表（14 词，逐词标出处）+ 切词规则 + **不扫什么**（注释/测试/main，作为覆盖面限制而不是隐藏）· §3 十三包逐包准入/禁入/内部允许依赖/core 现状出处（`<!-- scope:begin/end -->` 区是 CI 登记表）· §4 依赖方向与允许的 import 面 · §5 **九条实测偏差登记** · §6 护栏用法与红向 · §7 与 ARCHITECTURE/db-schema/contracts API.md/KT-02 的互锁 |
| `tools/check_scope/` | 新建护栏。`go/ast` 扫**声明标识符**（函数/方法/类型/接口方法/结构体字段/const/var）与**非标准库 import 段**；四红向：业务词命中、越界依赖（core 首段 `tellurion`、兄弟仓 `github.com/ruoxizhnya/tellurion-*`）、SCOPE.md 未覆盖的包、清单结构红（缺标记 / 缺分隔行 / 十三包少行 → fatal） |
| `tools/check_scope/main_test.go` | 22 条用例（夹具 `testdata/ok/`：十三包齐的清单 + 干净 `logging` + `sandbox/runner` 带同仓 import 与 `os/signal` 两条正腿） |
| `.github/workflows/ci.yml` | Vet 与 Test 之间加「包边界核对」步骤 |
| `AGENTS.md` / `README.md` / `docs/README.md` / `docs/ARCHITECTURE.md` | 登记第四份真相源与变更触发行、命令；`ARCHITECTURE.md` §3 那张 5 行表**就地标注为定位级简表**，逐包真相指向 SCOPE.md（避免两份清单各改各的） |
| `docs/TASKS.md` | KT-01 标 ✅ 并把六字段正文搬进本文件，台账只留指针 |

**⑤ 验收逐条（真退出码）**

1. **`docs/SCOPE.md` 覆盖 12 包 + `db` 逐包成表**：§3 表 13 行（logging / errors / id / decimal / statistics / httpclient / httpserver / auth / authz / observability / clock / sandbox / db）。
   `go run ./tools/check_scope` → 退出码 **0**，输出「1 个库包已扫 / 1 个 main 包跳过｜SCOPE.md 登记 13 包 ✅」。
2. **AST 护栏接入 CI**：`.github/workflows/ci.yml` 步骤「包边界核对」，命令与本地一致。
3. **破坏验证（现场，真退出码）**：临时新建 `db/scratch_destruction_leg.go`（`import "tellurion/pkg/storage"` + `type OrderRow struct{ ID string }`）→
   同一条命令退出码 **1**，两条点名到文件与行号：
   `db/scratch_destruction_leg.go:3 import "tellurion/pkg/storage"：越界依赖…` 与
   `db/scratch_destruction_leg.go:5：类型 "OrderRow" 命中业务词 "order"（出处：AGENTS §3）` →
   删除该文件与目录 → 退出码回到 **0**，`git status` 只剩本次交付的新文件（无 scratch 残留）。
4. **夹具里的同一条腿长期在案**（不需要每次手改代码）：`TestDeclaredTypeWithBusinessWordIsRed`（含「必须给到文件:行」的位置断言）、
   `TestFunctionNameWithBusinessWordIsRed`、`TestStructFieldWithBusinessWordIsRed`、`TestInterfaceMethodWithBusinessWordIsRed`、
   `TestPackageLevelVarWithBusinessWordIsRed`、`TestCoreImportIsRed`、`TestSiblingRepoImportIsRed`、
   `TestThirdPartyBusinessPathSegmentIsRed`、`TestUnregisteredPackageIsRed`、`TestMissingScopeRowIsFatal`（删一行等于少登记一个包 → fatal）、
   `TestMissingMarkersIsFatal`、`TestMissingTableSeparatorIsFatal`、`TestTwoPackageNamesInOneDirIsFatal`（防「静默少扫一半」）、
   正腿与误红防线：`TestChainAloneIsNotBusinessWord`（`Chain`/`LinkedChain` 不许命中，否则链表实现被拦）、
   `TestFactorialIsNotBusinessWord`（`Factorial` 不许误命中 `factor`，同时 `FactorDefinition` 必须命中）、
   `TestMainPackageIsNotScanned`、`TestTestFilesAreNotScanned`、`TestCommentsAreNotScanned`（三条覆盖面限制钉成断言）、
   `TestStdlibPathDetection`（`os/signal` 判标准库）、`TestTokensSplitting`（切词表的绝对断言，含 `SupplyChainSchemaDDL` → 6 段）。

**执行中如实登记的偏差（本单产出的一半价值在这两批）**

- **口径**：单标题写「12 包」，QT-05 ③-② 写「12 个包 + 新增 `db`」→ SCOPE.md 落 **13 行**，工具里 `scopePackages` 也冻结为 13。这不是加戏，是两个口径对齐后的唯一自洽数。
- **词表**：`AGENTS.md` §3 给 6 词 + QT-05-B ③ 给 5 词 = 11；本单补 `ohlcv` / `bar` / `signal` 三个（理由逐词写在 §2 出处列）。
  `supply_chain` 是复合词：切成 `[supply chain]` 后单段整词比对抓不到，而把 `supply` 或 `chain` 单独列进词表会误红正当的链表类型 →
  解法是**相邻两段拼接**参与比对（`SupplyChain` → `supplychain`），两条腿各有断言。
- **实测九条泄漏**（core 今天的形状 vs. 清单，逐条带 `文件:行`）登记在 `SCOPE.md` §5。其中四条会让 QT-05-B「照搬即红」：
  `pkg/id/order.go` 的 `OrderID`、`pkg/clock` 的 `ModeBacktest`、`pkg/observability` 的 `BacktestDuration`/`ObserveBacktest(strategyName, universe)`、
  以及 `pkg/auth` 带 `users`/`audit_logs` 的 10 处 SQL（与 `db` 的「零表结构」正面冲突，属**裁决**不是属改名）。
  `pkg/sandbox/wasm` 的 `BarContext`/`hostGetBar`/`OnBarSession` 且 import `tellurion/pkg/backtest/contracts` —— 与既有裁决（❌ 留 core，ADR-027 §2.6.1）一致，护栏只是把它钉住。
- **「迁移执行」今天没有可搬的 runner**：core 的 `migration_manager.go` 已于 2026-09-18 删除，唯一迁移路径是 `pkg/storage/postgres.go` 的 `migrate()` + `migrationStatements()`（约 30 条内联 DDL）。
  所以 `kit/db` 的「migrate」只能是**机制**（接收语句序列并执行），DDL 文本一条都不进本仓 —— SCOPE.md §3 的 `db` 行与 §5 偏差 #8 都写明了，`ARCHITECTURE.md` §3 那四个字容易被读成「把建表搬过来」，已就地加指向。
- **覆盖面限制（不是缺陷，是设定）**：注释与字符串不扫、`*_test.go` 不扫、`package main` 不扫。
  代价如实写在 §2：注释里可以藏业务语义、测试里可以藏形状 —— 由 §3 逐包禁入清单 + 人工评审兜，不假装 AST 扫描覆盖了语义。
  理由：注释要说清调用方是谁（「供回测循环使用」是正当表达），扫注释会把文档本身变成违规物。

**边界**：只写边界清单与护栏，没写任何实现（`go.mod` 零变化，仍只有根包 `doc.go`）；符号级导出面归 contracts CT-01（本单不建 `docs/API.md`，
kit 的符号级真相要等 QT-05-B 落码后再议，见 §5 与 contracts `docs/API.md` §8 的待裁决）；版本规则/兼容矩阵/CHANGELOG/GOPRIVATE 归 KT-02，本单一条没碰。

---

## KT-02 · README 版本规则表 + 兼容矩阵 + CHANGELOG + `GOPRIVATE` 前提 [P1] ✅ 完成（2026-10-22）

- **① 目标**：本仓作为被依赖库有成文的版本纪律与拉取前提，breaking 有账可查。
- **② 上下文**：refactor-plan-review **L4/L5**；规则口径冻结在 `tellurion:docs/TASKS.md` QT-19 ③（**规则内容以彼处为准，本单只落本仓文档**）；`tellurion:docs/design/qt-decoupling/repo-inventory.md` §1.1/§5 的可见性与 GOPRIVATE 风险注。
- **③ 要求**：① README 写**版本规则表**（`MAJOR` = 删/改签名或语义；`MINOR` = 新增兼容能力；`PATCH` = 文档/无行为修正）；② **兼容矩阵**：`vX.Y` ← 支持 core / agent 的哪些版本，每次发版更新（QT-19 ③ 的「kit 仓同表」）；③ CHANGELOG 至少一条记录；④ AGENTS.md 写一行 **`GOPRIVATE=github.com/ruoxizhnya/*` + git 凭证**前提（本仓转私有时 CI/构建必须配，否则 `go get` 走 proxy 必然失败）。
- **④ 约束与非目标**：不引 release 工具链（goreleaser 等）；0.x 阶段允许 breaking 但须 bump MINOR + CHANGELOG 条目（细则以 QT-19 为准）。
- **⑤ 验收**：README 存在版本规则表与兼容矩阵；CHANGELOG ≥ 1 条；AGENTS.md 含 GOPRIVATE 行；**破坏验证**：发一版不更新兼容矩阵 → 发版 checklist 项不可勾（或检查脚本红）。
- **⑥ 边界**：只定本仓策略与文档；contracts 的同名交付是 `tellurion-contracts:docs/TASKS.md` CT-02；「禁删 deprecated」机械护栏属 `tellurion:docs/TASKS.md` QT-09 ⑤。

### 结案记录（2026-10-22）

**交付物**

| 文件 | 内容 |
|---|---|
| `README.md` | 新增「版本规则」（三行表 + 0.x 条款 + 三条 Go 特有判定补充 + 内容寻址不参与 + tag 规范）、「兼容矩阵」（三列使用方，基线行见现状说明）、「发版流程」（四步 checklist，其中两条有 CI 护栏）、「拉取前提（私有仓）」（两条：`GOPRIVATE` + git 凭证，含「本仓转私有先红的是 core/agent CI」这一层）；`## 目录` 与「快速开始」补两份新交付 |
| `CHANGELOG.md` | 新建。文件头写明标题约定 `## [vX.Y.Z] - YYYY-MM-DD` 与 breaking 的写法要求；两条 PATCH 级条目（KT-02 本次、KT-01 2026-10-22），全在 `[Unreleased]` 下 —— **本仓还没发版，所以没有版本节** |
| `AGENTS.md` | §2 加 `check_release_matrix` 命令 + **私有化前提**那一行（含 `repo-inventory.md` §1.1/§5 风险注出处与「底座转私有先红下游 CI」）；§6 登记「第五份 = 版本纪律三件套」，把「发版 → 改矩阵那一行 + CHANGELOG 那一节」写进变更触发串 |
| `tools/check_release_matrix/` | 新建护栏（KT-02 ⑤ 的「检查脚本红」腿）。六个红向：结构红 / 发版没更新矩阵 / 矩阵行没账 / tag 没账 / GOPRIVATE 前提缺失 / 空账 |
| `tools/check_release_matrix/main_test.go` | 16 条用例（夹具 `testdata/ok/` = 已发版且齐；`testdata/baseline/` = 本仓今天的零发版形状，必须绿） |
| `.github/workflows/ci.yml` | 加「版本纪律核对」步骤；checkout 改 `fetch-depth: 0`（默认只取 HEAD 时 `git tag` 恒空，tag 那条腿会变成假绿） |
| `docs/README.md` | §1 加「发版 / 改版本纪律该读哪份」、§2 文档地图登记 `CHANGELOG.md` + 新护栏、§3 变更触发表加「发版 → 矩阵那一行 + CHANGELOG 那一节」 |
| `docs/SCOPE.md` | §7 互锁表那一行由「版本纪律归 KT-02」改为指向已落地的三件套（任务号改成已实现的事实） |
| `docs/TASKS.md` | KT-02 标 ✅、六字段正文搬进本文件，台账只留指针 |

**③ 要求逐条落点**

| 要求 | 落在哪 |
|---|---|
| ① 版本规则表 | `README.md`「版本规则」，`<!-- rules:begin/end -->` 之间，MAJOR/MINOR/PATCH 三行 |
| ② 兼容矩阵 | `README.md`「兼容矩阵」，`<!-- matrix:begin/end -->` 之间，按版本序列 `vX.Y` 建行、三列使用方 |
| ③ CHANGELOG ≥ 1 条 | `CHANGELOG.md` 两条 PATCH 级条目（在 `[Unreleased]` 下，理由见「现状说明」） |
| ④ GOPRIVATE 前提 | `AGENTS.md` §2 的「私有化前提」那一行（护栏按 `GOPRIVATE=github.com/ruoxizhnya/` 子串核对）+ `README.md`「拉取前提」两条细则 |

**⑤ 验收逐条（真退出码）**

1. **README 存在版本规则表与兼容矩阵**：`go run ./tools/check_release_matrix` → 退出码 **0**，
   输出「规则表 3 行 ｜ 矩阵序列 0 ｜ CHANGELOG 版本节 0（条目 2）｜ tag 0 ✅」。（同一判据由 CI 步骤「版本纪律核对」执行。）
2. **CHANGELOG ≥ 1 条**：条目数由脚本数 `- ` 行，输出里的「条目 2」即判据；删空 → 红（夹具腿 `TestChangelogWithoutEntriesIsRed`）。
3. **AGENTS.md 含 GOPRIVATE 行**：脚本核对子串存在；删掉 → 红（`TestGoprivatePrerequisiteMissingIsRed`）。
4. **破坏验证（现场，真退出码，全程在临时副本上做，不碰仓内文件）**：
   - 对照腿：把仓根三件套原样拷进 `.kt02destroy/` → 同一条命令（`-root .kt02destroy`）退出码 **0**；
   - 破坏腿：往副本的 `CHANGELOG.md` 追加 `## [v0.2.0] - 2026-10-22`（矩阵不动）→ 退出码 **1**，
     点名「CHANGELOG 有 v0.2.0 这一版，矩阵却没有 v0.2 那一行 —— 发版必须同步兼容矩阵」；
   - 补账腿：在副本矩阵补 `v0.2` 那一行 → 退出码回到 **0**（证明发版 checklist 是**可完成**的，护栏不是拦路石）；
   - 清理：删除 `.kt02destroy/`，`git status` 只剩本次交付的改动（无 scratch 残留）。
5. **破坏验证（夹具腿，长期在案）**：`TestReleaseWithoutMatrixRowIsRed`（本单主判据）、
   `TestMatrixRowWithoutChangelogIsRed`（反向：矩阵留了没账的序列）、`TestTagWithoutChangelogSectionIsRed` 与
   `TestTagWithoutMatrixRowIsRed`（tag 有版本没账）、`TestMissingMatrixMarkersIsFatal`（删标记让检查通过）、
   `TestEmptyMatrixDataRowsIsFatal` 与 `TestMissingTableSeparatorIsFatal`（清空矩阵 / 删表头分隔行让「跳过表头」的逻辑变松）、
   `TestMissingRuleRowIsRed`、`TestMalformedMatrixCellIsRed`（打错的 `v0.1.0-beta` 不能被当信息行放过）、
   `TestDuplicateMatrixSeriesIsRed`（`v0.1` 与 `v0.1.0` 同序列两行）、`TestReleaseTitleWithoutDateIsRed`、
   `TestListTagsFailureIsFatal`（读不到 tag 就红，不许静默按零 tag 走下去）。
6. **回归**：`gofmt -l .` 空输出（rc=0）、`go build ./...` rc=0、`go vet ./...` rc=0、`go test ./... -count=1` rc=0；
   `go run ./tools/check_scope` rc=0，输出「1 个库包已扫 / **2 个 main 包跳过**｜SCOPE.md 登记 13 包」。

**现状说明（写进文档的口径，不是缺陷）**

- 本仓**零符号、零 tag**，所以兼容矩阵今天只有一行「`main`（未发版）」的基线行；首个版本 `v0.1.0` 由 QT-07 ③-S2-a 在 S2 打。
  「矩阵里有没有我的版本」对当前三方使用方的答案都是「没有，你在用 `main`」。
- 护栏的 **tag 那条腿今天是空跑**（0 个 tag）。价值在 `TestBaselineWithoutReleaseIsGreen`（证明零发版不被误判成违规）与 4 条 tag 夹具腿；
  CI 侧的前提（`fetch-depth: 0`）已一并改掉，否则那条腿在 CI 里永远是绿的假象。
- 跨仓那半（「三个仓的 require 是否都落在矩阵声明范围内」）**没做也做不了**：本仓 CI 不 checkout 别的仓（R2），已在 README 写明归属 QT-09。

**执行中如实登记的偏差**

- **矩阵列数**：QT-19 ③-③ 对 kit 的口径是「← contracts / core 的兼容范围」，本单标题 ③-② 又写「core / agent」。
  落地为**三列**（contracts / core / agent）：两处口径并集，且 agent 确实直接 require 本仓（QT-05 ③-⑤：`require` 不传递，agent 必须自己 require 两个仓）——
  只写两列会漏掉一个真实使用方，这不是加戏。
- **序列 vs 补丁版**：QT-19 ③-③ 的矩阵口径写 `vX.Y`，tag/CHANGELOG 是 `vX.Y.Z`。按**序列建行**、脚本按序列对齐（与 contracts CT-02 同一取舍），
  否则 S2 打第一个 `v0.1.0` 时会因「矩阵写 `v0.1`、账写 `v0.1.0`」而红——那是格式摩擦不是纪律。
- **零发版的 CHANGELOG**：⑤ 的「至少一条记录」在未发版现实下只能是 `[Unreleased]` 条目；没有为凑「版本节」虚构一个没发过的版本。
- **超出 ③-④ 的 README 补充**：③-④ 只要求 AGENTS.md 一行；README 同时写了「拉取前提」两条细则（proxy 绕过之外还要 git 凭证），
  因为只配 `GOPRIVATE` 不配凭证同样是拉不到——出处是 `repo-inventory.md` §1.1 可见性行与 §5。
- **kit 特有的 breaking 判定**：README「判定补充」比 QT-19 ③-① 多一条「移动 / 删除 / 改名一个包 = breaking」。
  理由：本仓的对外单位是包，使用方断的是 import 路径；这条与 SCOPE.md §3 的「增删改名包 → 改那一行」是同一条纪律的两面。只写进文档，不改 QT-19 口径。
- **新工具不进 SCOPE.md**：`tools/check_release_matrix` 是第二个 `package main`，`tools/check_scope` 的「未登记包」红向对 main 包整目录跳过
  （KT-01 已登记的覆盖面限制）。实测输出「2 个 main 包跳过」即证据 —— 所以 §3 十三包**不需要**加一行；记在这里，免得后来者以为是漏登记。

**边界**：没引 release 工具链（goreleaser 等一个没加，`go.mod` 零变化）；没写任何符号；没碰 `docs/{db-schema,api-schema,dataflow}.md`；
「禁删 deprecated」的机械护栏留给 QT-09 ⑤；contracts 的同名交付 CT-02 已在 `2026-10-22` 独立结案（本单没有改动 contracts 仓任何文件）。

---

## QT-05-B · S1 工具层落地：12 个包 + `db` 迁入本仓（QT-05 的 kit 腿）[P0] ✅ 完成（2026-10-22）

> **拆分说明（2026-10-21）**：QT-05（S1 契约层与工具层落地）主单在 `tellurion-contracts:docs/TASKS.md`（同号 QT-05，六字段正文权威版，2026-10-21 自 core `tellurion:docs/TASKS.md` 迁入）；本单是其 **kit 腿**（主单 ③-② 的执行单），按仓归属拆出、判据以主单为准。两腿同期交付，**kit 先行落码**（contracts 依赖本仓 `decimal`/`errors`，主单 ② 实测）。

- **① 目标**：横切工具**物理上住进本仓**（12 个包 + 新增 `db`），零业务依赖；core 与 agent 经**版本依赖**引用本仓。
- **② 上下文**：主单 QT-05 的 ② 与形态裁决（ADR-D31：代码直接写进新仓，core/agent 用 `require <伪版本>`，全程不 `replace`；本地联调用 `go.work` 且不入库）。包清单（12 个横切包）源出 `tellurion:docs/design/qt-decoupling/target-module-topology.md` 包边界表（含依赖方向订正 banner：`kit ← contracts(含 domain) ← core`）。本仓现状：`github.com/ruoxizhnya/tellurion-kit`、`go 1.25`、public。`docs/SCOPE.md` 载体是 KT-01，不在本单。
- **③ 要求**：主单 ③-② 原文——`tellurion-kit`：迁 12 个包 + 新增 **`db`**（仅 connection/migrate/DSN 解析，**零表结构**，K 裁决落地）；另落地主单 ④/⑤ 中直接约束本仓的判据：一次只搬一个包、每搬一次 `go build ./...` + `go vet ./...`；本仓 import `tellurion/...` 必须**编译失败**；`kit` 业务词零命中（core 专用类型 `screen`/`sector`/`supply_chain`/`macro`/`mining` 不进本仓，主单 ③-⑨）。
- **④ 约束与非目标**：**不改任何业务逻辑**；`db` **零表结构**（K 裁决）；不搬 `fees` 的计算部分（主单 ④）；不在本单打 tag（S2 = QT-07）；不新增 `replace`。
- **⑤ 验收**（主单 ⑤ 中本仓腿的判据）：单独 checkout 本仓 `go build ./...` + `go test ./...` 绿；在本仓 import `tellurion/...` 编译失败；`kit` 业务词零命中；core/agent 侧 `go list -m all` 的 kit 行带真实伪版本且 `git grep -E '=> *\.\./tellurion-kit'` 零命中（core 腿执行，结果对拍）；**破坏验证**：在本仓任一包塞一个业务词 import ⇒ 护栏/构建红，还原绿。
- **⑥ 边界**：只做本仓包落地；契约类型是主单 QT-05（contracts 腿）；core/agent 依赖切换与 `check_contract_pin.py` 是主单 ③-③/④/⑤/⑦（core 侧）；tag / 版本策略 / 跨仓护栏属 QT-07 / QT-09 / QT-19；`docs/SCOPE.md` 属 KT-01。

### 结案记录（2026-10-22）

**交付物：13 包 / 8 个 atomic commit（一包一 commit，`go build` + `go vet` 每步当场跑）**

| commit | 落了什么 | 与 core 同名包的形状差异（切依赖时会咬人的地方） |
|---|---|---|
| `869f4fd` | `errors`（`AppError` + 9 个 `ErrCode*` + 包装）、`decimal`（定点运算，624 行）；**顶层包布局**（不是 `pkg/...`）——`tools/check_scope` 按路径首段对 `SCOPE.md` §3 的行，布局错了护栏对不上 | 逐字照搬；`ErrCodeDataQuality` 按 SCOPE §5 #6 保留（分类码不是业务规则） |
| `cfe0bff` | `logging`（logger + 8 个脱敏函数）、`id`、`clock` | `id` **不导出 `OrderID`**（§5 #1）；`clock` 的 `ModeBacktest` → **`ModeDeterministic`**，`ModeLive` 原名保留，**iota 顺序不动**（§5 #7） |
| `46eb404` | `statistics`（均值/方差/相关/分位/回归/排名/z 分数） | **不带** `LogReturns` / `AnnualizeVolatility`（§5 #2：年化窗口与收益率定义是口径，不是数学） |
| `173de8f` | `httpclient`、`httpserver`（cors / errors / ginmode） | 两处跨包 import 收成**同仓内部边**（`httpclient → logging`、`httpserver → errors`）；`gin.SetMode` 仍只允许 `httpserver/ginmode.go` 与 `TestMain`（AUD-29 护栏随包搬来，195 行） |
| `952234e` | `observability`（私有 registry + `RegisterCollectors` + `/metrics` handler + request-id 中间件/transport）；第三方依赖钉到与 core 同版本，`go 1.25.0` | **不带** `BacktestDuration` / `ObserveBacktest` / `LLMTokens` / `NTCall*`（§5 #5：业务指标族由使用方自建 collector 挂进来） |
| `db2ca4c` | `sandbox/staticcheck`、`sandbox/runner`（含 `rlimit_posix.go` / `rlimit_windows.go` 双腿） | **`sandbox/wasm` 不在本仓**（它 import `tellurion/pkg/backtest/contracts`，逐 bar 喂数据必须在引擎进程内，ADR-027 §2.6.1）；`golang.org/x/sys` 由间接转直接 |
| `1e03390` | `auth`（token 签发/解析/刷新、bcrypt、claims、中间件）、`authz`（route key → `Rule` + 执法中间件） | **`auth` 不带一条 SQL、不带 `pgxpool`**（§5 #3）：身份回读经 `IdentityLoader` 端口，`PublicPaths` 由调用方传入，角色是不透明字符串（`CanTrade`/`CanAdmin` 语义留 core）；**`authz` 的 `Policy` 是机制不是 core 的路由表**（§5 #4），未登记的变更端点 fail-closed |
| `67c3232` | **`db`（新建）**：`BuildDSN` + 三个 DSN 哨兵错误、`Connect` / `ConnectRedis` / `DefaultPoolOptions`、`Migrate(ctx, Execer, statements)` | **零 DDL、零迁移记录表**（§5 #8 / K 裁决）：只接收语句序列并执行，空白跳过、遇错即停并报第几条（1 起）。`redis/go-redis` 登记进 `AGENTS.md` §4 |

文档真相源同步（本次 docs commit）：`docs/api-schema.md`（§1 十三行 ⬜→✅ 并逐行按导出面订正、§2 补分类→状态码映射表与三条 fail-closed、§3 从「不适用」改成两条真实承诺、§5 补「今天没有」七条）、
`docs/SCOPE.md`（§5 九条偏差的最后一列由「后果」改成**裁决与落地**）、`docs/ARCHITECTURE.md`（§2 标注现状 + 三条内部边实测、§8 S1 标 ✅）、
`README.md`（状态行、目录表加 13 包那一行、集成测试跑法、矩阵基线行、版本规则的例子换成真符号）、`CHANGELOG.md`（`[Unreleased]` 加 MINOR 节）、`doc.go`（根包只有文档）、`docs/TASKS.md`（QT-05-B 标 ✅ 留指针 + 新立 KT-03/04/05）。

**⑤ 验收逐条（真退出码；文档改完后于 2026-10-22 重跑一遍取数，本次 docs commit 携带的就是这份结果）**

1. **单独 checkout 本仓 `go build ./...` + `go test ./...` 绿**：
   `gofmt -l .` 空输出 **EXIT=0**；`go build ./...` **EXIT=0**；`go vet ./...` **EXIT=0**；
   `go test ./... -count=1` **EXIT=0** —— 15 个 `ok`，2 个 `[no test files]`（根包 `kit` 与 `errors`；后者见 KT-05）。
   最慢的是 `httpclient` **≈14s**（两次取样 14.156s / 14.140s；退避硬编码 `time.After`，6 条重试测试真的睡；已写进 KT-04 ②）。
2. **在本仓 import `tellurion/...` 编译失败**：`go list tellurion/pkg/storage` → **EXIT=1**，`package tellurion/pkg/storage is not in std`；
   `go list -m all | grep -i tellurion` → 只有 `github.com/ruoxizhnya/tellurion-kit` 一行（**没有任何仓外 `replace`**）；
   `git grep -nE '"tellurion/|tellurion-contracts|tellurion-agent' -- '*.go'` → 3 处命中，全在 `tools/check_scope/main_test.go`，
   是护栏自己的**红向夹具字符串**（不是 import 语句）。
3. **`kit` 业务词零命中**：`go run ./tools/check_scope` → **EXIT=0**，输出
   「包边界核对：**15 个库包已扫 / 2 个 main 包跳过**｜SCOPE.md 登记 13 包 ✅ 业务词零命中、无越界 import、十三包清单齐且覆盖现有包」。
   （`screen` / `sector` / `supply_chain` / `macro` / `mining` 一个都没进本仓。）
4. **版本纪律三件套仍一致**：`go run ./tools/check_release_matrix` → **EXIT=0**，
   「规则表 3 行 ｜ 矩阵序列 0 ｜ CHANGELOG 版本节 0（**条目 9**）｜ tag 0 ✅」——落码没破坏 KT-02 的账（今天仍零 tag，矩阵仍只有 `main` 基线行）。
5. **真服务集成腿（`//go:build integration`，默认不跑）**：`TEST_DATABASE_DSN=postgres://…@localhost:5432/postgres?sslmode=disable` +
   `TEST_REDIS_URL=redis://127.0.0.1:6379/0` 时三条腿全 PASS —— `TestConnect_RealPostgres`（真连 + 真 ping）、
   `TestConnectRedis_RealRedis`、`TestMigrate_RealPostgresInRolledBackTx`（**在总是回滚的事务里建 TEMP 表 `kit_db_migrate_probe`**，跑完零残留）。
   未设这两个变量时逐条 `t.Skip` 并说明原因，不静默假绿。
   ⚠️ 第一次跑用的是 `…/tellurion?sslmode=disable`，报 `FATAL: database "tellurion" does not exist`（本机没有起容器，只有一个本地 PG 实例）——
   换库名后通过。**这条如实记在这里**：本仓的集成腿今天依赖本机恰好有一个可达的 PG/Redis，CI 里没有（CI 不跑 `-tags integration`）。
6. **`db` 的编译期钉子**：`var _ Execer = (*pgxpool.Pool)(nil)` 在 `db/connect_test.go`，保证「`Migrate` 的入参接口」与真实 pool 不会各走各的。

**没做的两条腿（如实登记，不当成绿）**

- **现场破坏验证没做**：⑤ 要求「在本仓任一包塞一个业务词 import ⇒ 护栏红，还原绿」。本次执行环境把**改坏文件**这一类动作拦在授权之外，
  所以没有现场改一次源码再还原。**长期在案的等价证据**是 `tools/check_scope/main_test.go` 的红向夹具（随 `go test ./…` 每次都跑）：
  `TestDeclaredTypeWithBusinessWordIsRed`（`OrderID`）、`TestFuncWithBusinessWordIsRed`（`ObserveBacktest`）、
  字段名（`BarCount`）、接口方法（`EmitSignal`）、以及越界 import（`import "tellurion/pkg/storage"` 与 `tellurion-contracts/domain`）各一条，
  外加误红防线（`Chain` / `LinkedChain` / `Factorial` 不许命中）。**差别要写清**：夹具证的是「护栏对这类形状会红」，
  不证「今天这 13 个包被同一条命令扫过并且红了会点名到文件:行」——后者由上面第 3 条的正向输出（15 个库包已扫）覆盖，但**不是同一条腿**。
  这条差异留给下一位：真要做现场腿，需要一次显式授权（sha256 前后对拍 + 单行改坏 + 还原）。
- **`-race` 跑不了**：`go test -race` → `-race requires cgo; enable cgo by setting CGO_ENABLED=1`；加 `CGO_ENABLED=1` → `cgo: C compiler "gcc" not found`。
  本机没有 gcc，不是可以顺手修的配置。**代价**：`logging.Logger` 是包级可变全局、`gin.SetMode` 是进程级全局，这两处正是竞争高发面；
  今天的缓解是**结构性**的（gin 模式集中在 `httpserver/ginmode.go` 与各包 `TestMain`，由 AUD-29 护栏钉住；`logging.Init` 只在启动期调），
  但**没有 `-race` 实测背书**。已写进 KT-05 ②，等有 gcc 的环境（或 CI 的 Linux runner）再补这条腿。

**执行中如实登记的偏差与发现**

- **`decimal.Div` 的注释与实现互相矛盾**：`:385` 写 banker's rounding，`:400` 写 round-half-away-from-zero，`:406` 的判据 `2*absRem > absB` 实际是
  「恰好一半向零截断」。core 那份同样如此（逐字照搬）。按 ④「不改任何业务逻辑」**没有动它**，改口径属 MAJOR 语义变更 ⇒ 立 **KT-03**，
  并在 `docs/api-schema.md` §1 的 `decimal` 行挂 ⚠️。
- **`httpclient` 重试发空 body + 4xx 静默当成功**：**实测取证**（临时 module 在 `C:\tmp\kthc`，`replace` 指向本仓，跑真 `httptest` 服务，EXIT=0）：
  `maxRetries=2` 的 POST，服务端第一次收到 `len=17`、**第二次收到 `len=0`**；400 响应回来是 `err=<nil>` + `StatusCode:400`。
  根因是 `client.go:57-64` 把 body 做成一个 `*bytes.Reader` 在重试循环里复用（`http.NewRequestWithContext` 不会自动 rewind，`GetBody` 只服务重定向）。
  按 ④ 未改 ⇒ 立 **KT-04**，现状写进 `docs/api-schema.md` §3。
- **`errors` 零单测、`logging` 初始化路径未测**（继承自 core，不是本次引入）⇒ 立 **KT-05**（P2）。
- **`id` 的 UUIDv7 生成器未导出**：SCOPE §3 那行写「通用 UUIDv7」，落地只有 `newV7`（包内）。已在 SCOPE §5 #1 与 `api-schema.md` §5 登记为「今天没有」，
  不假装清单与实现一致。
- **`observability` 没有 `/health` handler**：SCOPE §3 写的是「handler 暴露」，落地只有 `/metrics`。健康判据是「这个服务的依赖齐不齐」，属业务装配 ⇒
  在 `api-schema.md` §5 明确写「不提供」，并改掉我先前在 §1 里顺手写的「`/health` handler」（那是**文档比实现多承诺**，比少承诺更危险）。
- **`httpclient` 没有熔断**：§1 原写「带超时 / 重试 / **熔断**的客户端」，实现里没有 breaker ⇒ 订正为「超时 + 重试」，并在 §5 登记「今天没有，要加是 MINOR 且先要有使用方」。
- **文档比实现多承诺的两处是我自己写出来的**（`/health`、熔断），在逐行核对导出面时被抓到并订正 ——
  记在这里是因为它说明「⬜ 未落地」那张表转 ✅ 时**不能只翻状态列**，「提供什么」那一列同样是契约。
- **CRLF 自查的测量方法有坑**：`grep -c $'\r' <file>` 在本机（Git Bash）对**纯 LF 文件也返回非零**（CR 作为 argv 传进去被 mangle 成空模式），
  一度看起来像"全部文件被注入了 CRLF"。可靠做法是 `tr -cd '\r' < file | wc -c`（实测 0）与 `git ls-files --eol`（实测 `i/lf w/lf`，`.gitattributes` 的 `* text=auto eol=lf` 在生效）。
  本次全部改动实测 CR 字节数 0。

**边界（守住 ④/⑥ 的证据）**

- **没打 tag**（`git tag -l 'v*'` 空；S2 = QT-07）；**没加 `replace`**（`go.mod` 里只有 `module` + `go` + require 两段，`go list -m all` 无仓外路径）；
  **没提交 `go.work` / `go.work.sum`**（`.gitignore` 挡着，`git status` 也没有）；**没搬 `fees`**；**没改任何业务逻辑**（`decimal` / `httpclient` 两处缺陷只登记不修）。
- **core / agent 侧一行没动**：本单 8 个 commit 全在 `tellurion-kit`。切依赖（`require` 伪版本 + `check_contract_pin.py`）是 QT-05 主单 ③-③，
  那半的判据（`go list -m all` 的 kit 行带真实伪版本、`git grep -E '=> *\.\./tellurion-kit'` 零命中）**本仓 checkout 不到别的仓，做不了**，留给 core 腿对拍。
- `auth` 的**审计中间件**与 `users` / `audit_logs` 的读写留 core（§5 #3 的裁决），本仓 `auth` 不认识任何表名。
- 临时取证目录 `C:\tmp\kthc` 是本单自己的脚手架，用完删除，不在仓内。
