# TASKS 归档 · 2026-10-23

> 规则见 [`README.md`](README.md)。本文件保存当日结案的任务正文与结案记录，**旧代正文不改写**（订正只在文件头加说明）。

---

## KT-03 · `decimal.Div` 的舍入口径与源码注释不一致（裁决口径 + 让注释与实现同一句话）[P1] ✅ 完成（2026-10-23）

- **① 目标**：`decimal` 的除法舍入行为有**一个**成文口径，且源码注释、实现、`docs/api-schema.md` 三处说的是同一句话。
- **② 上下文**：`decimal/decimal.go:385` 注释写「Round-half-to-even (banker's rounding)」，`:400` 又写「Simple round-half-away-from-zero for determinism」，
  而 `:406` 的判据是 `2*absRem > absB` —— 恰好一半（`2*absRem == absB`）时 `shouldRound=false`，即**向零截断**。三句话互相矛盾，实现是第三种。
  代码逐字搬自 core `pkg/decimal/decimal.go`（同一份矛盾在那边也在），QT-05-B 按 ④「不改任何业务逻辑」没动它。
  `README.md`「版本规则」把「精度与舍入口径」明确列为 **MAJOR 语义**，所以这不是注释洁癖。
- **③ 要求**：① 先取证：写一条**钉住当前行为**的表驱动测试（含恰好一半的正负两侧、奇偶除数、`resultScale` 大于/小于被除数 scale），
  让"今天到底是什么口径"变成可执行的断言；② 裁决目标口径（banker's / half-away-from-zero / half-toward-zero 三选一，写清为什么——
  钱的场景通常要 half-even 或 half-up，向零截断会系统性少算）；③ 若裁决与当前实现不同，**分两步**：先改注释与文档到"当前真实行为"，
  再单独一次改实现并按 README 记 breaking（0.x 期间 bump MINOR + CHANGELOG 写「破坏性：旧 → 新」）；④ 同步 core 侧同一份代码的处置（改 / 不改 / 等切依赖后自然收敛），在 `tellurion:docs/TASKS.md` 留同号指针。
- **④ 约束与非目标**：**不顺手改 `Mul` / `Round` / `FromString` 的口径**（各自单独裁决）；不引入 rounding-mode 配置项（那是新导出面，要单独一张单）；
  不在本单里做 core 的切换（属 QT-05 主单）。
- **⑤ 验收**：`go test ./decimal/... -count=1` 绿，且新增测试**先红后绿**（钉当前行为的那条在实现变更前就绿；变更口径的那条在变更前红、变更后绿——两次真退出码都留证）；
  `grep -n "banker" decimal/decimal.go` 的结果与裁决一致（要么消失，要么实现真的是 banker's）；`docs/api-schema.md` §1 的 `decimal` 行不再挂 ⚠️；
  **破坏验证**：把判据从 `2*absRem > absB` 改成 `>=` → 钉当前行为的测试必须红（证明测试真的在测口径，不是恒真）。
- **⑥ 边界**：只做 `decimal` 的除法舍入口径；`httpclient` 的重试缺陷是 KT-04；core 侧切换是 QT-05 主单。

### 结案记录（2026-10-23）

**裁决：half away from zero（四舍五入）**。三条理由：

1. **同包一致性**：`Round`（`decimal.go` 的 `remainder >= pow/2`，测试 `3.145 → 3.15`）已经是 half-away-from-zero。`Div` 若走 banker's 或向零截断，
   同一个包内「先除再舍」与「先舍再除」就会给出不同答案，而没有任何成文解释。
2. **金额方向**：向零截断在正数上系统性少算、在负数上系统性多算，钱的场景不可接受。
3. **代价最小**：banker's 虽无偏，但要让全包一致就得连带改 `Round`（另一张单、另一笔 breaking）；half-away 只动 `Div` 一处。

**两步走（两条真退出码）**：

- **① 取证 + 文档对齐现状**（commit `56a8326`）：新增 `decimal/div_rounding_test.go`，表驱动钉住当时的口径（17 条恰好一半的行 + 4 条 `resultScale` 小于被除数 scale 的行 + 其余不变行），
  `go test ./decimal/... -count=1` **EXIT=0**；`decimal.go` 的 `Div` 注释与 `docs/api-schema.md` §1 改成与实现同一句话，并删掉自相矛盾的 banker's 注释与空分支死代码（行为零变化）。
- **② 变更口径**：先把表中恰好一半的期望值与 `2.5 / 4 @ scale 0` 改成目标值，`go test ./decimal/... -run TestDiv_Rounding -count=1` → **EXIT=1**（tie 行与 `2.5_div_4_at_0` 全红，实际值都是旧口径）；
  再把判据从 `2*absRem > absB` 改为 `absRem >= absB-absRem`（等价于 `2*absRem >= absB`，但不会在 `absB > 2^62` 时溢出），`go test ./decimal/... -count=1` → **EXIT=0**。

**破坏验证**：把 `>=` 还原成 `>` → `go test ./decimal/... -run TestDiv_Rounding -count=1` **EXIT=1**（17 条 tie 行全部红）；还原后 **EXIT=0**。证明这张表真在测口径，不是恒真。

**其余验收**：

- `grep -n "banker" decimal/decimal.go` → **0 命中**（旧注释已删，实现真的是 half-away-from-zero）。
- `docs/api-schema.md` §1 的 `decimal` 行不再挂 ⚠️，写的是「余数超过或恰好等于半个单位都向远离零进位，与 `Round` 同口径」。
- 全仓 `go build ./...` / `go vet ./...` / `go test ./... -count=1` / `go run ./tools/check_scope` / `go run ./tools/check_release_matrix` 全部 **EXIT=0**。
- 本机 `CGO_ENABLED=0` 且无 C 编译器，`-race` 跑不了（与 QT-05-B 结案记录同一条限制）；CI 在 ubuntu 上带 race，本条如实登记为未在本机执行。
- CHANGELOG 的 `[Unreleased]` 下新增 breaking 条目，逐条写明「旧 → 新」与使用方影响面（本仓与 core 今天零生产调用点）。

**取证时新发现（拆单不修）**：`resultScale` 小于被除数 scale 时先截断被除数再除（`2.5 / 1` 取整得 `2`），以及 `shift > 0` 溢出时静默少乘 10 的幂（错量级）。
两条都不属「恰好一半怎么走」，按「不顺手修未列出语义」的规则拆为 **KT-06**（仍在 `docs/TASKS.md`）。

**core 侧处置**：**不改** core `pkg/decimal` 的同一份代码——该包在 core 仍标「待采用」（零生产调用），且 core 切本仓依赖属 QT-05 主单；
处置结论与同号指针落在 `tellurion:docs/TASKS.md` §八「跨仓任务索引」。
