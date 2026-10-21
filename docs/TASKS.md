---
status: active
last-verified: 2026-10-21
verified-by: 任务台账批量重分布（QT-05-B 拆分落本仓）
---

# TASKS — tellurion-kit 未完成项

> 本仓**唯一**任务台账。跨仓任务只在 §二 留指针，**不镜像正文**（双写必然漂移）。
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `✅ 完成（附日期）`
> 每张单六字段（① 目标 ② 上下文 ③ 要求 ④ 约束与非目标 ⑤ 验收 ⑥ 边界），判据见 `tellurion:AGENTS.md` §8.5。
> 完成即标 ✅（附日期）并移入 `docs/archive/TASKS-YYYY-MM-DD-archived.md`；本文件只留未完成项。

---

## 一、本仓任务

> **占位仓**：实现落地前不写业务代码；KT-01/KT-02 两张是**随仓文档交付**（refactor-plan-review M6/L4/L5 修复传播，2026-10-20 登记），可与 S1 落码并行。**QT-05-B**（S1 落码 kit 腿，2026-10-21 拆分落本仓）在两单之后。
>
> 第一张单长什么样（照此写，别自造格式）：
>
> ### XX-01 · <一句话目标> [P?] ⬜
> - **① 目标**：…
> - **② 上下文**：本仓路径/行号 + 冻结契约的路径与节号
> - **③ 要求**：①②③（可执行、可核对）
> - **④ 约束与非目标**：明确不做什么
> - **⑤ 验收**：可复现的命令 + 期望输出；至少一条"破坏验证"（改坏它应当红）
> - **⑥ 边界**：只做什么；相邻的什么归别的单

### KT-01 · 随仓接口规格载体 `docs/SCOPE.md`（12 包逐包准入/禁入清单）[P1] ⬜

- **① 目标**：本仓每个横切包有随仓、可自证的准入/禁入边界，下游（core/agent）不读 core 实现即可知什么该进 kit、什么不许进。
- **② 上下文**：refactor-plan-review **M6**（报告：`tellurion:docs/design/qt-decoupling/refactor-plan-review-2026-10-20.md`）；上游编排单 `tellurion:docs/TASKS.md` QT-07 ③-S2-c；包边界表源出 `tellurion:docs/design/qt-decoupling/target-module-topology.md`（含其依赖方向 ⚠️ 订正 banner：`kit ← contracts(含 domain) ← core`）。S1 迁移范围见 `tellurion:docs/TASKS.md` QT-05（12 个横切包 + 新增 `db`）。
- **③ 要求**：① 建 `docs/SCOPE.md`：按 QT-05 的 12 包 + `db` 逐包列**准入清单**（该包收什么）与**禁入清单**（业务词、领域语义、core 内部类型禁止入包）；② 每包标注依赖方向约束（kit 不 import contracts 以外的业务层，最终不变量 `kit ← contracts(含 domain) ← core`）；③ 护栏进 CI（AST 扫描）：与「`kit` 内业务词零命中」验收（QT-05）同源，包边界违规即红。
- **④ 约束与非目标**：只写边界清单不写实现细节；不做跨进程 DTO（那是 `tellurion:docs/openapi.yaml` 生成物、归 core，见 `tellurion:docs/design/qt-decoupling/repo-inventory.md` §3）；不含版本规则（那是 KT-02）。
- **⑤ 验收**：`docs/SCOPE.md` 覆盖 12 包 + `db` 逐包成表；AST 护栏接入 CI；**破坏验证**：故意往任一包塞一个业务词/越界 import → 护栏红（留证后撤回）。
- **⑥ 边界**：只做包准入/禁入载体与护栏；导出面接口规格（符号级）归 `tellurion-contracts:docs/TASKS.md` CT-01；版本策略归 KT-02。

### KT-02 · README 版本规则表 + 兼容矩阵 + CHANGELOG + `GOPRIVATE` 前提 [P1] ⬜

- **① 目标**：本仓作为被依赖库有成文的版本纪律与拉取前提，breaking 有账可查。
- **② 上下文**：refactor-plan-review **L4/L5**；规则口径冻结在 `tellurion:docs/TASKS.md` QT-19 ③（**规则内容以彼处为准，本单只落本仓文档**）；`tellurion:docs/design/qt-decoupling/repo-inventory.md` §1.1/§5 的可见性与 GOPRIVATE 风险注。
- **③ 要求**：① README 写**版本规则表**（`MAJOR` = 删/改签名或语义；`MINOR` = 新增兼容能力；`PATCH` = 文档/无行为修正）；② **兼容矩阵**：`vX.Y` ← 支持 core / agent 的哪些版本，每次发版更新（QT-19 ③ 的「kit 仓同表」）；③ CHANGELOG 至少一条记录；④ AGENTS.md 写一行 **`GOPRIVATE=github.com/ruoxizhnya/*` + git 凭证**前提（本仓转私有时 CI/构建必须配，否则 `go get` 走 proxy 必然失败）。
- **④ 约束与非目标**：不引 release 工具链（goreleaser 等）；0.x 阶段允许 breaking 但须 bump MINOR + CHANGELOG 条目（细则以 QT-19 为准）。
- **⑤ 验收**：README 存在版本规则表与兼容矩阵；CHANGELOG ≥ 1 条；AGENTS.md 含 GOPRIVATE 行；**破坏验证**：发一版不更新兼容矩阵 → 发版 checklist 项不可勾（或检查脚本红）。
- **⑥ 边界**：只定本仓策略与文档；contracts 的同名交付是 `tellurion-contracts:docs/TASKS.md` CT-02；「禁删 deprecated」机械护栏属 `tellurion:docs/TASKS.md` QT-09 ⑤。

### QT-05-B · S1 工具层落地：12 个包 + `db` 迁入本仓（QT-05 的 kit 腿）[P0] ⬜

> **拆分说明（2026-10-21）**：QT-05（S1 契约层与工具层落地）主单在 `tellurion-contracts:docs/TASKS.md`（同号 QT-05，六字段正文权威版，2026-10-21 自 core `tellurion:docs/TASKS.md` 迁入）；本单是其 **kit 腿**（主单 ③-② 的执行单），按仓归属拆出、判据以主单为准。两腿同期交付，**kit 先行落码**（contracts 依赖本仓 `decimal`/`errors`，主单 ② 实测）。

- **① 目标**：横切工具**物理上住进本仓**（12 个包 + 新增 `db`），零业务依赖；core 与 agent 经**版本依赖**引用本仓。
- **② 上下文**：主单 QT-05 的 ② 与形态裁决（ADR-D31：代码直接写进新仓，core/agent 用 `require <伪版本>`，全程不 `replace`；本地联调用 `go.work` 且不入库）。包清单（12 个横切包）源出 `tellurion:docs/design/qt-decoupling/target-module-topology.md` 包边界表（含依赖方向订正 banner：`kit ← contracts(含 domain) ← core`）。本仓现状：`github.com/ruoxizhnya/tellurion-kit`、`go 1.25`、public。`docs/SCOPE.md` 载体是 KT-01，不在本单。
- **③ 要求**：主单 ③-② 原文——`tellurion-kit`：迁 12 个包 + 新增 **`db`**（仅 connection/migrate/DSN 解析，**零表结构**，K 裁决落地）；另落地主单 ④/⑤ 中直接约束本仓的判据：一次只搬一个包、每搬一次 `go build ./...` + `go vet ./...`；本仓 import `tellurion/...` 必须**编译失败**；`kit` 业务词零命中（core 专用类型 `screen`/`sector`/`supply_chain`/`macro`/`mining` 不进本仓，主单 ③-⑨）。
- **④ 约束与非目标**：**不改任何业务逻辑**；`db` **零表结构**（K 裁决）；不搬 `fees` 的计算部分（主单 ④）；不在本单打 tag（S2 = QT-07）；不新增 `replace`。
- **⑤ 验收**（主单 ⑤ 中本仓腿的判据）：单独 checkout 本仓 `go build ./...` + `go test ./...` 绿；在本仓 import `tellurion/...` 编译失败；`kit` 业务词零命中；core/agent 侧 `go list -m all` 的 kit 行带真实伪版本且 `git grep -E '=> *\.\./tellurion-kit'` 零命中（core 腿执行，结果对拍）；**破坏验证**：在本仓任一包塞一个业务词 import ⇒ 护栏/构建红，还原绿。
- **⑥ 边界**：只做本仓包落地；契约类型是主单 QT-05（contracts 腿）；core/agent 依赖切换与 `check_contract_pin.py` 是主单 ③-③/④/⑤/⑦（core 侧）；tag / 版本策略 / 跨仓护栏属 QT-07 / QT-09 / QT-19；`docs/SCOPE.md` 属 KT-01。

---

## 施工序（候选单，开工前必须展开成六字段）

> 这张表是**施工顺序**，不是任务台账：每条在执行前要展开成六字段单（目标/上下文/要求/约束/验收/边界），并写明冻结契约的路径与节号。
> 顺序有依赖，**不要跳级**：跳序的代价是返工（例：X1a 的协议没冻结就写 aggregator，等于把协议写死在一个服务里）。

| 阶段 | 内容 | 上游单号 | 前置 |
|---|---|---|---|
| S1 | 迁 12 个横切包 + 新增 `db`（connection/migrate/DSN，零表结构）（正文见「一、本仓任务」QT-05-B） | `QT-05-B` | W0 全绿 |
| S1 | `kit` 内业务词零命中（验收项） | `QT-05-B` | 同上 |
| S2 | 打 `v0.1.0`；下游换 tag | `QT-07` | S1 全绿 |

## 二、跨仓指针（不在本仓，勿双写）

| 单号 | 主题 | 在哪 |
|---|---|---|
| QT-05 | S1 主单（契约层腿 + core 侧依赖切换），本仓腿拆为 QT-05-B | `tellurion-contracts:docs/TASKS.md`（2026-10-21 自 core 迁入；core `tellurion:docs/TASKS.md` 留同号指针） |
| QT-07 | S2：打 `v0.1.0` 并把下游伪版本换成 tag | `tellurion:docs/TASKS.md` |
| QT-17 | agent 数据自持会用到本仓的 `db` | `tellurion-agent:docs/TASKS.md`（2026-10-21 自 core 迁入；core 留同号指针） |

---

## 维护本文件

- **新增单**：写全六字段；六字段里必须显式引用**冻结契约的路径与节号**（契约仓的包路径、设计文档的节号），不许写"参考相关实现"。
- **派工前先证伪前提**：读码核对工单写的前提；不符就停下回报 + 给拆法，不硬做、不扩范围。
- **跨仓引用写 `repo:path`**（如 `tellurion:docs/TASKS.md`），不做可点击相对链接。
- **顺序纪律**：同一波次内撞同一批文件的单按单号串行，一单一 commit。
