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

---

## KT-04 · `httpclient` 重试发空 body + 4xx 静默当成功（两条都已实测取证）[P1] ✅ 完成（2026-10-23）

- **① 目标**：`httpclient.Client.Do` 的重试**不丢请求体**，且调用方不会把一个 4xx 当成成功。
- **② 上下文**：实测取证（2026-10-22，临时 module `replace` 到本仓跑真 `httptest` 服务，退出码 0）：
  `maxRetries=2` 的 POST 第一次服务端收到 `len=17`，**第二次收到 `len=0`** —— `client.go:57-64` 把 body 做成一个 `*bytes.Reader` 后在重试循环里复用，
  首次尝试已把它读到 EOF，`http.NewRequestWithContext` 不会自动 rewind（`GetBody` 只用于重定向）；
  同一次实测里 400 响应回来是 `err=<nil>` + `StatusCode:400`（`client.go:135` 对非 5xx/429 一律当成功返回）。
  重试触发面是「传输错误 / 读 body 失败 / 5xx / 429」且**不看方法**，所以丢 body 会打在 POST 上。
  代码逐字搬自 core `pkg/httpclient/client.go`，QT-05-B 未改。现状已登记在 `docs/api-schema.md` §3。
  ⚠️ **顺带实测的代价**：退避是硬编码的 `time.After(1<<(attempt-1) * time.Second)`，不可注入，所以 6 条重试测试真的睡——
  `go test ./httpclient/... -count=1` 实测 **≈14s**（两次取样 14.156s / 14.140s；4 条各 3s + 2 条各 1s），是本仓测试耗时的绝大部分。
  改重试语义时若不让退避可注入（clock 或 interval 参数），这份代价会照旧留在每次 CI 上。
- **③ 要求**：① 先写**复现测试**（`httptest` 数收到的 body 长度序列 + 断言 4xx 的 `err`/`StatusCode` 语义），改代码前必须红；
  ② 修 body 复用：每次尝试重新构造 reader（保留已 marshal 的 `[]byte`，`bytes.NewReader(data)` 每轮新建）；
  ③ 裁决 4xx 的返回语义（返回 `*AppError`（`ErrCodeInvalidInput` 之类）还是保持 `err=nil` 只在文档里钉死），
  **两者选一即可，但必须与 `docs/api-schema.md` §3 同一句话**；若选前者，那是 breaking（调用方的 `if err != nil` 分支会变），按 README 规则记账；
  ④ 裁决「非幂等方法要不要重试」：要么加一个「只对幂等方法重试」的开关，要么在文档里写明由调用方用 `maxRetries=0` 关掉——不要留成隐含行为；
  ⑤ 同步 core 侧处置，在 `tellurion:docs/TASKS.md` 留同号指针。
- **④ 约束与非目标**：**不引熔断器**（`api-schema.md` §5 已登记为"今天没有，要加是 MINOR 且先要有使用方"）；不改 `Request` / `Response` 的字段形状（那是导出面 breaking，单独裁决）；
  不动退避算法（`1<<(attempt-1)` 秒）——要改另立单。
- **⑤ 验收**：`go test ./httpclient/... -count=1` 绿；复现测试**先红后绿**（真退出码留证）；
  **破坏验证**：把「每轮新建 reader」改回「循环外复用同一个 reader」→ 复现测试必须红；
  `docs/api-schema.md` §3 的 ⚠️ 段落与最终语义一致（不再有"已登记缺陷"的口径，或明确写清保留哪一半）。
- **⑥ 边界**：只做 `httpclient`；`decimal` 是 KT-03；core 侧切换是 QT-05 主单。

### 结案记录（2026-10-23）

**裁决（③）：4xx 返回 `nil` + `*errors.AppError`（不保留 `err == nil` 的旧语义）**。理由：KT-04 ① 的目标写的是「调用方不会把一个 4xx 当成成功」，
只有 fail-closed 才真正满足；`err != nil` 必须是唯一的「没成功」门。错误码按状态映射：400 / 422 / 未映射 4xx → `INVALID_INPUT`，
401 / 403 → `PERMISSION`，404 → `NOT_FOUND`，408 → `TIMEOUT`，409 → `CONFLICT`；`message` 为 `unexpected HTTP status <code>[: <body 前 200 字>]`，`op` 为 `"<METHOD> <Path>"`。
429 仍走重试，重试耗尽后维持「max retries exceeded」错误。该语义已与 `docs/api-schema.md` §3 写成同一句话；breaking 记账在 CHANGELOG `[Unreleased]`。
由此新增 `kit` 内部边 `httpclient → errors`，按 `SCOPE.md` §4 的要求留了 [ADR 0002](../adr/0002-httpclient-4xx-作为错误.md)，并同步更新 `SCOPE.md` §3 / §4 / §5 #9 与 `ARCHITECTURE.md` 的边清单。

**裁决（④）：重试仍不看方法，文档钉死**。不加「只对幂等方法重试」的开关（那是新导出面，且会改变默认行为）；有副作用的调用用 `New(baseURL, timeout, 0)` 关重试，
口径写在 `docs/api-schema.md` §3。这样没有隐含行为，也不引入本单范围外的新配置面。

**先红后绿（真退出码）**：

- **红**（改代码前）：新增 `httpclient/retry_body_test.go`（`TestDo_RetriesWithBody_PreservesBodyOnEveryAttempt` + `TestDo_4xxReturnsAppError` 表驱动 + `TestDo_2xxIsStillSuccess`），
  并把 `TestDo_DoesNotRetryOn404` 的期望从「`err == nil`」改为「有错误」。
  `go test ./httpclient/ -count=1 -run "TestDo_(RetriesWithBody_PreservesBodyOnEveryAttempt|4xxReturnsAppError|DoesNotRetryOn404)$"` → **EXIT=1**，
  失败原文：`attempt 2 body = "", want "{\"n\":1}"`、`expected an error for 404`、以及 400/401/403/404/408/409/422/418 八个子用例的 `expected an error ... got nil`。
- **绿**（改代码后同一条命令）→ **EXIT=0**；随后 `go test ./httpclient/ -count=1` → **EXIT=0**（两次取样 15.107s / 15.234s），`go vet ./httpclient/` → **EXIT=0**，`gofmt -l httpclient` 无输出。
  自审后另补 `TestDo_4xxTruncatesUpstreamBody`（250 个汉字的上游 body → 错误里只保留 200 字 + `...`，且 `utf8.ValidString` 为真），钉住截断口径与「按 rune 切、不劈开 UTF-8」。

**破坏验证**：把「每轮新建 reader」改回「循环外复用同一个 reader」→ `go test ./httpclient/ -count=1 -run TestDo_RetriesWithBody_PreservesBodyOnEveryAttempt` **EXIT=1**，
原文 `attempt 2 body = "", want "{\"n\":1}"`；还原后 **EXIT=0**。证明该测试真在测重发，不是恒真。

**core 侧处置**：**不改** core `pkg/httpclient` 的同一份代码（core 的调用点仍用自己的副本，本单不改其行为）；core 切到本仓依赖属 QT-05 主单，
届时各调用点都已先判 `err != nil`，4xx 语义只会变严格。处置结论与同号指针落在 `tellurion:docs/TASKS.md` §八「跨仓任务索引」。

**顺带登记（不修）**：退避仍是硬编码 `1<<(attempt-1)` 秒、不可注入，`go test ./httpclient/... -count=1` 的耗时结构不变（本单 ④ 约束明确「不动退避算法」）。

