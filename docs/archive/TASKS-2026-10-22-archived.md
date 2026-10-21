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
