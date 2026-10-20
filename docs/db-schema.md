---
status: active
last-verified: 2026-10-20
verified-by: 功能边界固化（三份契约：db / api / dataflow）
---

# DB Schema — tellurion-kit

> **所有权声明**：本仓**不持有任何业务表**。`db` 包只做连接 / 迁移执行 / DSN 解析，**零表结构** —— 表结构属于各业务仓（`market` / `strategy` / `risk` / `analysis` / `run` / `ai`）。
> 上级依据：`tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md` §1 —— **DDL 只有一个 owner**，同一个 PG 实例、每服务一个 schema，将来拆库是"按 schema 搬走"而不是重写。
> 本文件是本仓表结构的**唯一真相源**：代码里的建表语句与它不一致时，**改代码**；要改表结构，先改本文件（结构级变更还要立 ADR）。

## 1. 归属表

> **无表**。
>
> `kit/db` 提供的是"连上去并跑迁移"的能力，**不含任何 DDL**。曾经的诱惑是把公共表（如 `symbols`）放进 kit 省一次重复 —— 那会让所有仓都能写它，等于取消数据归属（ADR-D19 §1 的三条纪律之一：每张表只归一个服务写）。

| kit 能力 | 它碰什么 | 它不碰什么 |
|---|---|---|
| `db` | 连接池、DSN 解析、跑 `migrations/` | 建表、改表、业务查询 |

## 2. 只读输入（不属于本仓，禁止写入）

不适用（工具层不访问业务数据；调用方把自己的连接传进来）。

## 3. 红线

1. **零业务概念**：`backtest` / `strategy` / `order` / `factor` / `market` 这类词出现即违规（QT-05 验收项之一）。
2. **零表结构**：`db` 里出现 `CREATE TABLE` 就是边界破了。
3. **不反向依赖**：kit 是依赖图底座，不许 import contracts 或任何应用仓。
4. 不做"框架"：只提供能力，不规定调用方的分层。
5. 新增第三方依赖前先在 `AGENTS.md` 第 4 节登记。

## 4. 迁移策略

- 迁移脚本落在 `migrations/`，命名 `<序号>_<动作>.sql`（全仓一个递增序列，不复用号）。
- **只前滚**：已合并的迁移文件不改；修错靠新增一条。
- 每条迁移必须可重入（`IF NOT EXISTS` / `IF EXISTS`），且能在空库上一次跑完。
- 跨 schema 取数只走**只读视图或 API**；出现 `join` 别人的表就是边界腐化的第一信号。
- 表结构变更的同步义务：本文件 + `docs/ARCHITECTURE.md` 的「数据归属」节必须同一次提交里更新。

## 5. 待定

- `db` 的迁移执行语义（谁持有序号、并发跑怎么办）随 QT-05 落地时定稿。
