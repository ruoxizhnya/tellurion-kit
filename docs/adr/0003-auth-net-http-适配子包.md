# 0003 · auth 的 net/http 适配走独立子包

- **状态**：已裁决（来源：KT-08 / 裁-15，2026-10-25）
- **场景**：`auth.Middleware` / `auth.RequireRole` 与 `authz.Middleware` 都是 `gin.HandlerFunc`，而 risk / analytics / data 三个服务用的是标准库 `net/http`（risk / analytics 已用 Go 1.22+ `ServeMux` 的 `"POST /path/{id}"` 模式注册路由）。直接复用只有 `auth.TokenService.Parse`，中间件与启动闸得各服务各写一份，语义必然漂移。
- **决策**：
  1. 新增子包 `auth/httpadapter`（**不 import gin**），提供 `Authenticate`（`Authorization: Bearer` → `Parse` → 拒 `kind != access` → 校验 subject）、`Authorize`（Public / AnyAuthenticated / Roles 三档，角色 any-of）、`Guard`（net/http 中间件，成功时把 claims 挂进 `r.Context()`）、`DecideStartup`（与 core `bootstrap.DecideAuthStartup` 同语义）与 `Policy`（route key → Rule）。
  2. 失败面用哨兵：`ErrUnauthenticated` → 401，`ErrForbidden` → 403，调用方用 `errors.Is` 映射。
  3. route key 取 `http.Request.Pattern`（net/http 从 `ServeMux` 注册模式填入），因此 `Guard` 按路由挂载；拿不到 pattern 时 fail-closed（变更方法 403、读方法 401）。
  4. 未在 `Policy` 声明的**变更方法**路由一律 403 并点名 route key —— 与 gin 侧 `authz.Middleware` 同一条 runtime backstop。
- **理由**：把「token 怎么验、角色怎么比、没有密钥能不能启动」留在 kit，把「哪个端点要哪个角色」留在各服务（SCOPE §5 #4）。独立子包让 net/http 服务的调用面只有标准库签名（`*http.Request` / `http.Handler`），源码层不 import gin，也不用为了鉴权在服务里写第二套 token 校验。它复用 `auth` 的 token 机制而不是复制一份，因此新增一条明确的内部边 `auth/httpadapter → auth`（父包是 stdlib + jwt，无回边风险）。
- **代价 / 后果**：
  - **依赖闭包仍含 gin**：`auth` 包同时托管 gin 中间件（`middleware.go`），Go 的包粒度使 `auth/httpadapter → auth` 的编译依赖闭包包含 gin。源码层（本包 import 列表）不 import gin；要彻底摘掉需把 `auth` 拆成 token 核心 + gin 适配，移动导出路径属 breaking，不在本单。
  - `Guard` 必须逐路由挂载（`mux.Handle(pattern, Guard(svc, policy)(h))`），漏挂的端点不受保护。缓解手段与 gin 侧相同：服务侧守卫测试用 `Policy.Entries()` 与路由表做双向对拍。
  - 包装整个 mux（而非逐路由）时 `r.Pattern` 为空，所有请求 fail-closed（403/401）而不是静默放行；这是有意的失败方向。
  - 读方法未在策略表声明时仍要求有效 access token（Public 读端点须显式登记）。这与 gin 侧「auth 管全部非公开路径、authz 只管变更方法」的分工一致，只是合并进了一个中间件。
- **替代关系**：不替代 0001；与 0002 无重叠。