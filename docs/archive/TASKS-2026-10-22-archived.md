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
