---
status: active
last-verified: 2026-10-22
verified-by: KT-01 落地（§3 简表标注为示意，逐包真相在 SCOPE.md）
---

# 架构 — tellurion-kit

## 1. 定位

横切工具层：**零业务概念**。它是依赖图的底座，被所有仓 `require`，正因为如此，它一旦混进业务概念，就会变成新的耦合中心。

## 2. 分层

```
（底座，无内部依赖）

  logging · errors · id · decimal · statistics
  httpclient · httpserver · auth/authz · observability
  clock · sandbox(staticcheck/runner) · db(connection/migrate/DSN)
```

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

## 7. 非目标（明确不做）

- 不做业务抽象（没有"通用 Repository"这种东西）
- 不做 DI / 框架
- 不放任何表结构或迁移 SQL
- 不替代标准库能干净解决的事

## 8. 分期与状态

| 阶段 | 内容 | 状态 |
|---|---|---|
| S1（QT-05） | 迁 12 个横切包 + 新增 `db` | ⬜ |
| S2（QT-07） | 打 `v0.1.0` | ⬜ |
| QT-17 | agent 用 `kit/db` 自建 `ai.*` | ⬜ |
