---
status: active
last-verified: 2026-10-25
verified-by: KT-08 / KT-09 / KT-10 落地（§8 分期表补 S2.5 / S2.6 / S3 三行；§2 包清单与内部边未变）
---

# 架构 — tellurion-kit

## 1. 定位

横切工具层：**零业务概念**。它是依赖图的底座，被所有仓 `require`，正因为如此，它一旦混进业务概念，就会变成新的耦合中心。

## 2. 分层

```
（依赖图底座：不 import 任何仓外内部模块）

  logging · errors · id · decimal · statistics
  httpclient · httpserver · auth/authz · observability
  clock · sandbox(staticcheck/runner) · db(connection/migrate/DSN)
```

十三包 QT-05-B 已全部落地（2026-10-22）；`auth/httpadapter` 子包由 KT-08 增加（2026-10-25）。包间**只有五条**同仓内部边，实测（`grep -rn "ruoxizhnya/tellurion-kit/" --include="*.go"`，排除测试与 `tools/`）：

```
httpclient → logging · errors      httpserver → errors      authz → auth      auth/httpadapter → auth
```
（`httpclient → errors` 由 KT-04 增加，理由见 [ADR 0002](adr/0002-httpclient-4xx-作为错误.md)；`auth/httpadapter → auth` 由 KT-08 增加，理由见 [ADR 0003](adr/0003-auth-net-http-适配子包.md)。）

再加一条就是新耦合，要在 ADR 里留账（判据与出处见 [`SCOPE.md`](SCOPE.md) §3 的「kit 内允许依赖」列与 §4）。

## 3. 组件职责

| 包 | 边界 | 反例 |
|---|---|---|
| `logging` | 日志形状与级别 | 不放业务埋点 |
| `errors` | 包装与分类 | 不定义业务错误码语义（那是 contracts） |
| `decimal` | 定点运算 | 不定义费率模型（那是 contracts 的 `fees`） |
| `db` | 连接 / 迁移执行 / DSN | 不含表结构 |
| `sandbox` | 进程隔离与静态检查 | 不加载策略（那是 core 的 WASM host） |

> 这张表是**定位级简表**（5 行示意）。逐包「收什么 / 不收什么」的完整清单、业务词表、允许的 import 面，
> 以及「core 今天的形状与清单有哪些出入」的实测偏差登记，都在 [`SCOPE.md`](SCOPE.md)——那才是 `tools/check_scope` 核对的登记表。
> 两处有出入时以 `SCOPE.md` 与 CI 为准。

## 4. 三份契约（功能的边界就在这三份文件里）

| 维度 | 文件 | 回答什么问题 |
|---|---|---|
| 数据形状 | [`db-schema.md`](db-schema.md) | 本仓拥有哪些表？谁能写？ |
| 接口形状 | [`api-schema.md`](api-schema.md) | 对外暴露什么？错误码与幂等怎么做？ |
| 流向与顺序 | [`dataflow.md`](dataflow.md) | 数据怎么流？哪里 fail-closed？ |

## 5. 依赖方向

```
third-party libs
        ↓
       kit
        ↑
contracts / core / agent / 各服务（都可以直接 require kit）
```

## 6. 关键决策

- [`adr/0001-零业务概念.md`](adr/0001-零业务概念.md)
- [`adr/0002-httpclient-4xx-作为错误.md`](adr/0002-httpclient-4xx-作为错误.md)
- [`adr/0003-auth-net-http-适配子包.md`](adr/0003-auth-net-http-适配子包.md)

## 7. 非目标（明确不做）

- 不做业务抽象（没有"通用 Repository"这种东西）
- 不做 DI / 框架
- 不放任何表结构或迁移 SQL
- 不替代标准库能干净解决的事

## 8. 分期与状态

| 阶段 | 内容 | 状态 |
|---|---|---|
| S1（QT-05） | 迁 12 个横切包 + 新增 `db` | ✅ 本仓腿 **QT-05-B** 完成（2026-10-22）：13 包全部有实现；contracts 腿与 core/agent 依赖切换仍在主单 QT-05 |
| S2（QT-07） | 打 `v0.1.0` | ✅ 已发版（2026-10-23） |
| S2.5（KT-08） | `auth` 的 net/http 适配 + 启动闸，打 `v0.1.1` | ✅ 本仓完成（2026-10-25） |
| S2.6（KT-09） | `auth.IssueServiceToken` 服务账号 access token，打 `v0.1.2` | ✅ 本仓完成（2026-10-25） |
| S3（KT-10） | 二期切包前置：`statistics.LogReturns` / `AnnualizeVolatility`、`httpserver.IsLoopbackRemoteAddr`，打 `v0.1.3` | ✅ 本仓完成（2026-10-25）；`LoopbackOnly` gin 中间件留 core |
| QT-17 | agent 用 `kit/db` 自建 `ai.*` | ⬜ |
