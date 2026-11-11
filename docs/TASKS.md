---
status: active
last-verified: 2026-11-12
verified-by: 2026-11-12 KIT-DEP-SECURITY-BUMP 收口（govulncheck EXIT=0 / 六门禁绿 / 破坏验证真做真还原 / 发版 `v0.4.0`）；此前 2026-11-11 全面审查 W7 台账立单（A5 授权）——新立 3 条六字段单（`KIT-DEP-SECURITY-BUMP` / `KIT-SURFACE-FREEZE` / `KIT-REQUESTID-STDLIB`），证据坐标引 ODR-069/072（`tellurion-docs:docs/archive/odr/`）；`KIT-DEP-SECURITY-BUMP` 是 `tellurion:X-DEP-BUMP-CONSUMERS` 的前置
---

# TASKS — tellurion-kit 未完成项

> **归档说明（2026-11-11）**：本文件已瘦身为「只装未做的事」。2026-11-10 及以前已收口（✅）的任务原文、施工记录、验收输出与破坏验证，已逐字移入 [`docs/archive/tasks-archive-2026-11.md`](archive/tasks-archive-2026-11.md)；该文件顶部有「检索索引」表（ID / 标题 / 收口日期 / commit / 证据链接）。检索：`grep -n "<任务ID或commit>" docs/archive/tasks-archive-2026-11.md`。
>
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `⏸ 挂起（外部依赖或已裁决暂不实施；**不计入开放项**，重开须另立单并满足该条「重开条件」）` / `✅ 完成（附日期，随即移入归档）`

## 一、未做 / 挂起 / 外部依赖

### KIT-DEP-SECURITY-BUMP · 底座依赖安全 bump（pgx GO-2026-5004 / x/net / x/text）[P0] ✅ 2026-11-12

- **① 目标**：把 kit 的 pgx / x/net / x/text 升到安全版并发版，让依赖图底座先闭合漏洞，7 消费仓再跟升（`tellurion:X-DEP-BUMP-CONSUMERS`）。
- **② 上下文**：审查发现 W5-F1（ODR-072）+ 裁决 D1=b（pgx 一步收敛 v5.11.0）。GO-2026-5004（pgx SQLi）fixed v5.9.2，kit `db` 现 pin `pgx v5.7.3`；x/net×4 fixed v0.60.0、x/text×2 需升。kit 是底座，bump 后消费仓跟升才真生效。contracts 是唯一 govulncheck EXIT=0 的仓（零 pgx）。
- **③ 要求**：① pgx 升 v5.11.0（D1=b 一步到位，不逐档）；② x/net→v0.60.0、x/text 升安全版；③ `govulncheck ./...` EXIT=0；④ 发版三件套（CHANGELOG + 兼容矩阵 + tag 推远端）；⑤ AGENTS §4 依赖表同批。
- **④ 约束与非目标**：不手改版本号、不加 replace；不在本单顺手升无关依赖；pgx 大版本收敛的 breaking 面如实登记（底座一次打到 contracts/core/agent）。
- **⑤ 验收**：`govulncheck ./...` EXIT=0（E1）；`go build/vet/test ./... -count=1` 绿；`check_release_matrix` 绿；破坏验证 = pin 回旧版 ⇒ govulncheck 红。
- **⑥ 边界**：kit 依赖与发版；消费仓跟升归 `tellurion:X-DEP-BUMP-CONSUMERS`。
- **owner**：待定
- **重开条件**：无外部前置（A5 已授权立单；发版推远端由主任做）。
- **收口记录（2026-11-12）**：
  - 落地：pgx `v5.7.3 → v5.11.0`、x/net `v0.55.0 → v0.60.0`、x/text `v0.38.0 → v0.42.0`；MVS 连带 x/crypto `v0.53.0 → v0.57.0`、x/sys `v0.46.0 → v0.48.0`、x/sync `v0.21.0 → v0.23.0`。
  - **`go` 指令 1.25.0 → 1.26.11**（x/net v0.60.0 要求 ≥1.26.0；stdlib 29 条 advisory 要到 1.26.11 才齐，故钉补丁版）。版本判定 **MINOR → `v0.4.0`**（依据：工具链底线抬升是消费方可见变化，不是 PATCH 的「无行为修正」；导出面无改动）。
  - 验收（本机真实输出）：`govulncheck ./...` EXIT=0；`go build ./...` / `go vet ./...` / `go test ./... -count=1` 全绿；`go run ./tools/check_scope` 与 `check_release_matrix` 绿。
  - 破坏验证：x/text 降回 v0.38.0（x/net 被 MVS 压到 v0.56.0）⇒ govulncheck EXIT=3（6 条符号级，含 GO-2026-6629 / GO-2026-6617 / GO-2026-5970）；`go.mod` / `go.sum` 逐字节还原（SHA256 一致）后复绿 EXIT=0。单降 pgx 只到 package 级、默认模式 EXIT=0（kit 不调用其漏洞符号），故不作为破坏验证证据。
  - 残留（如实登记）：x/net v0.60.0 两条 module 级（fixed v0.62.0）、x/crypto 一条（GO-2026-5932），govulncheck 判「不调用」；本单版本由裁决指定为 x/net v0.60.0，不擅自再升。
  - 发版：tag `v0.4.0`（annotated，`main` 上）已推远端；CHANGELOG `v0.4.0` 节 + README 兼容矩阵 `v0.4` 行 + AGENTS §1/§4 同批。

### KIT-SURFACE-FREEZE · kit 缺导出面冻结护栏（对照 contracts 的 check_surface_freeze）[P1] ⬜

- **① 目标**：给 kit 一个导出面冻结护栏，让「新增/改名/删除导出符号」可机器判定，对齐 contracts 仓已有的 `check_api_surface` + `check_surface_freeze`。
- **② 上下文**：审查发现 W2-F1（ODR-069）。kit `docs/` 无 `API.md`/`api-surface.txt`；`tools/` 只有 `check_scope` + `check_release_matrix`，无导出面冻结。kit 是底座，breaking 面比别仓宽（README「判定补充」：移动/删除/改名包记 breaking），却无机器护栏钉住导出面。
- **③ 要求**：① 建 `docs/api-surface.txt` 基线 + `tools/check_surface_freeze`（对照 contracts 形状）；② 新增导出符号须先登记再更新基线；③ CI 增该护栏步骤；④ AGENTS §6 第四份真相源登记。
- **④ 约束与非目标**：不改现有导出符号（只建护栏）；不照搬 contracts 的业务词表（kit 零业务概念）；护栏破坏验证须真做。
- **⑤ 验收**：`check_surface_freeze` 绿；破坏验证 = 临时加一个未登记导出符号 ⇒ 红，还原复绿。
- **⑥ 边界**：kit 导出面护栏；contracts 侧已有护栏不动。
- **owner**：待定

### KIT-REQUESTID-STDLIB · RequestID 中间件只有 gin 版，6 服务仓 stdlib mux 装不上 [P1] ⬜

- **① 目标**：给 kit 一个 **stdlib `net/http` 兼容**的 RequestID 中间件，让 6 个用 stdlib mux 的服务仓能真装上请求 ID 传播。
- **② 上下文**：审查发现 W5-F3（ODR-072）。kit 现 RequestID 中间件是 gin 版，6 服务仓用 stdlib mux 装不上；core `observability/middleware.go:58` 零安装；agent 出站 HTTPTransport 有、入站 RequestIDMiddleware 无。ADR-027:1314-1316 / ADR-029:126-127/:448 / ADR-017:151 的「请求 ID 全链路」断言被证伪（订正归 `tellurion-docs:DOC-ADR-REQUESTID-CORRECTION`）。
- **③ 要求**：① kit 增 stdlib 版 RequestID 中间件（`func(http.Handler) http.Handler`）；② 与 gin 版共用同一 ID 生成/传播口径；③ 发版（只增不改，MINOR）；④ 6 服务仓安装归各仓 + `X-CI-RACE-STATICCHECK` 同批核对。
- **④ 约束与非目标**：不改 gin 版语义；不引入新依赖（用 stdlib + 现有 id 包）；不强制 core/agent 同批安装（各仓自持）。
- **⑤ 验收**：stdlib 中间件单测（入站生成/出站传播/已有 ID 透传）；破坏验证 = 删透传 ⇒ 用例红；`check_scope` 绿（不撞业务词）。
- **⑥ 边界**：kit stdlib 中间件与发版；各仓安装归各仓；ADR 断言订正归 docs 仓。
- **owner**：待定

**其余无未做 / 挂起 / 外部依赖项**（截至 2026-11-11 台账重整 + 全面审查 W7 立单）。已收口单的正文与施工记录全部在归档文件里；新增工作一律另立单。

## 二、跨仓指针（不在本仓，勿双写）

| 事项 | 去哪 |
|---|---|
| 契约形状与版本纪律 | `tellurion-contracts:docs/TASKS.md` |
| core 侧装配与鉴权消费 | `tellurion:docs/TASKS.md` |

## 维护本文件

- **新增单**：写全六字段；六字段里必须显式引用**冻结契约的路径与节号**（契约仓的包路径、设计文档的节号），不许写"参考相关实现"。
- **派工前先证伪前提**：读码核对工单写的前提；不符就停下回报 + 给拆法，不硬做、不扩范围。
- **跨仓引用写 `repo:path`**（如 `tellurion:docs/TASKS.md`），不做可点击相对链接。
- **顺序纪律**：同一波次内撞同一批文件的单按单号串行，一单一 commit。
