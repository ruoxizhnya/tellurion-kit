# TASKS 归档 · 2026-10-25

> 规则见 [`README.md`](README.md)。本文件保存当日结案的任务正文与结案记录，**旧代正文不改写**（订正只在文件头加说明）。

---

## KT-08 · `auth` 的 net/http 适配 + 启动闸（裁-15 的机制腿）[P1] ✅ 完成（2026-10-25）

- **① 目标**：给 `auth` 补一层**框架无关**的 net/http 适配与启动闸，供 risk / analytics / data 三个 net/http 服务复用；各服务的「端点 × 角色」策略表仍留各自仓。
- **② 上下文**（2026-10-25 实测）：`auth/middleware.go` 的 `Middleware()` / `RequireRole()` 是 `gin.HandlerFunc`，`authz.Policy` 也依赖 `gin.Context`；而三个服务全是 net/http（risk `127.0.0.1:8086`、analytics `127.0.0.1:8087`、data `:8081`）。目前可直接复用的只有 `TokenService.Parse`。core 的 `bootstrap.DecideAuthStartup` 语义正确但不在本仓。
- **③ 要求**：① 新增框架无关适配（建议 `auth/httpadapter`）：从 `Authorization: Bearer` 解析 → `Parse` → 拒 `kind != access` → 角色 any-of 判定，返回可映射为 401/403 的哨兵；② `DecideStartup(secret, allowInsecure, bindHost, exposure)` 与 core **同语义**：无 secret + 回环 = 开放，无 secret + 非回环 = 拒绝启动；③ 提供 `Policy`（route key → Public / AnyAuthenticated / Roles），**未知变更路由 fail-closed 403**；④ 打 `v0.1.1` 并在 CHANGELOG 记 MINOR。
- **④ 约束与非目标**：不搬用户表 / audit / 角色词表（角色仍是调用方的不透明字符串）；新包不 import gin；不改既有 `auth` / `authz` 导出面与其行为；不替服务决定 policy 内容。
- **⑤ 验收**：表驱动：无 secret+回环放行、无 secret+非回环拒绝启动、无 token 401、坏 token 401、refresh 当 access 401、角色不足 403 / 足够放行；policy 未声明变更路由 403；`go build ./...` / `go vet ./...` / `go test ./... -count=1` / `go run ./tools/check_scope` 全 0；**破坏验证**：把「未知变更路由」改成放行 ⇒ 对应用例红，逐字节还原后绿。
- **⑥ 边界**：只做机制与适配；三条服务腿各自的 policy 表与接线不在本单。

### 结案记录（2026-10-25）

**落地形状**：新增子包 [`auth/httpadapter`](../../auth/httpadapter)（`httpadapter.go` / `policy.go` / `httpadapter_test.go`），
导出面 **只增不改**：

| 符号 | 语义 |
|---|---|
| `ErrUnauthenticated` / `ErrForbidden` | 401 / 403 哨兵，调用方用 `errors.Is` 映射 |
| `Authenticate(svc, r)` | `Authorization: Bearer` → `auth.Parse` → 拒 `kind != access` → 校验 subject；失败包装 `ErrUnauthenticated` |
| `Authorize(claims, rule)` | Public / AnyAuthenticated / Roles 三档；角色 any-of；不足 `ErrForbidden` |
| `Guard(svc, policy)` | net/http 中间件；成功后 `ClaimsFromContext(r.Context())` 可读 claims；鉴权关闭时 no-op |
| `DecideStartup(secret, allowInsecure, bindHost, exposure)` / `IsLoopbackHost` / `InsecureExposureLoopbackPublished` / `Startup` | 启动闸，与 core `bootstrap.DecideAuthStartup` 同语义 |
| `Policy` / `Rule` / `NewPolicy` / `MustNewPolicy` / `Rule` / `Entries` | route key → 三档规则；构造期校验「恰好一档」；未声明变更路由运行时 403 |
| `IsMutating` / `RouteKey` / `ClaimsFromContext` | 方法分流、route key 解析（`http.Request.Pattern`，带方法前缀归一）、claims 读取 |

- **route key 来源**：`http.Request.Pattern`（net/http 在调用 handler 前从 `ServeMux` 注册模式填入），
  因此 `Guard` 按路由挂载：`mux.Handle(pattern, httpadapter.Guard(svc, policy)(h))`。拿不到 pattern 时 fail-closed（变更方法 403、读方法 401）。
- **gin 边界（如实登记）**：本包**源码 import 列表**无 gin（`go list -f '{{.Imports}}' ./auth/httpadapter` =
  `context encoding/json errors fmt kit/auth net net/http strings`）。但 `auth` 包同时托管 gin 中间件（`middleware.go`），
  Go 的包粒度使**编译依赖闭包仍含 gin**（`go list -deps ./auth/httpadapter` 命中 `gin-gonic/gin`）。
  彻底摘掉需把 `auth` 拆成 token 核心 + gin 适配，移动导出路径属 breaking，不在本单；口径写进 ADR 0003。

**表驱动用例**（`auth/httpadapter/httpadapter_test.go`，`go test ./auth/httpadapter/ -count=1`）：

- `TestDecideStartup` **9 例**：secret 开鉴权 / secret 压过 insecure 开关 / 无 secret 且未 opt-in 拒绝 / 纯空白 secret 拒绝 /
  无 secret + `127.0.0.1` 放行 / `localhost` 放行 / `0.0.0.0` 拒绝 / 容器 `loopback-published` 放行 / 拼错豁免拒绝。
- `TestGuard` **14 例**：无 token 401、坏 token 401、refresh 当 access 401、角色不足 403、两个在列角色放行、
  模板路由 `DELETE /api/widgets/{id}` 命中、any-authenticated（有/无 token）、Public（无 token / 坏 token 都放行）、
  **未声明变更路由 403（admin 也拦）**、未声明读路由仍要 token（无 token 401 / 有 token 200）。
- 其余：`TestGuard_UnknownMutatingRouteNamesTheRoute`（403 body 点名 route key）、`TestGuard_AuthDisabledIsOpenAccess`（nil / 无 secret 服务均放行）、
  `TestGuard_AttachesClaims`、`TestAuthenticate_Sentinels`、`TestAuthorize` 6 例、`TestNewPolicy_RejectsMalformedRules` 5 例、
  `TestMustNewPolicy_PanicsOnMalformed`、`TestPolicy_ReturnsCopies`、`TestRouteKey`、`TestIsMutating`、`TestIsLoopbackHost` 8 例。

**破坏验证（真退出码 + 逐字节还原）**：

1. 突变：`policy.go` 的 `if !declared && IsMutating(r.Method) {` 改为 `if false && !declared && IsMutating(r.Method) {`（等价于「未声明变更路由放行」）。
2. `go test ./auth/httpadapter/ -count=1 -run 'TestGuard$|TestGuard_UnknownMutatingRouteNamesTheRoute'` → **exit=1**：
   - `TestGuard/undeclared_mutating_route_fails_closed_for_admin`：expected 403，actual 200；
   - `TestGuard_UnknownMutatingRouteNamesTheRoute`：expected 403，actual 200。
3. 还原同一行后 `policy.go` SHA256 = `2B32336410A5C7276BADA204B9C552FE0E5F8B2456F435F74887445D0B4DCA00`（与突变前一致），
   同一条命令 → **exit=0**。

**全仓验收（真退出码，2026-10-25）**：

- `gofmt -l .` → 无输出（clean）。
- `go build ./...` → **exit=0**。
- `go vet ./...` → **exit=0**。
- `go test ./... -count=1` → **exit=0**（`auth/httpadapter` 0.114s；`httpclient` 15.166s 为既有退避睡眠）。
- `go run ./tools/check_scope` → **exit=0**：`16 个库包已扫 / 2 个 main 包跳过｜SCOPE.md 登记 13 包`，业务词零命中、无越界 import。
- `go run ./tools/check_release_matrix` → **exit=0**：`规则表 3 行｜矩阵序列 1｜CHANGELOG 版本节 2（条目 26）｜tag 1`（打 `v0.1.1` 前；打 tag 后为 2）。

**发版账**：`CHANGELOG.md` 新增 `## [v0.1.1] - 2026-10-25`（新增 `auth/httpadapter` 记 MINOR；同时收编 v0.1.0 之后的两条 PATCH：KT-06 `decimal.Div` 定标、KT-07 `logging.Init` 级别生效）；
`README.md` 兼容矩阵 `v0.1` 行更新到 `v0.1.1`；tag `v0.1.1`（annotated）打在结案 commit 上。

**未做 / 限制（如实登记）**：

- `-race` 未在本机执行（沿用 QT-05-B / KT-03 的同一限制：本机无 C 编译器）；CI 在 ubuntu 上带 race。
- 三条服务腿（risk / analytics / data）的 `Policy` 表与接线不在本单；本仓只交付机制。
- `Guard` 按路由挂载，漏挂的端点不受保护；服务侧需用 `Policy.Entries()` 与路由表双向对拍（backstop 不能替代守卫测试）。