---
status: active
last-verified: 2026-10-20
verified-by: 功能边界固化（三份契约：db / api / dataflow）
---

# Dataflow — tellurion-kit

> 目的：把「数据从哪来、在谁手里变形、到哪去」写死。**任何跨仓数据交换都必须能在本图上找到一条边**；找不到 = 这段代码在偷偷建第二条路。

## 1. 主链路

```
（无内部依赖：kit 是依赖图的根）

   third-party libs
        ↓
      kit  ──►  contracts  ──►  core / agent / 各服务
        └──────────────────────────►  （各服务直接用 kit 的 logging/db/http）
```

## 2. 关键时序

1. 服务启动：`db.Open(dsn, ...)` → 连接池
2. 服务自己跑 `migrations/*.sql`（kit 只执行，不提供内容）
3. 运行期：`logging` 记日志、`observability` 暴露指标、`httpclient` 出网
4. 时钟与随机一律从 `clock` 注入（保证可测试与可复现）

## 3. 失败模式（每一处 fail-closed 都要能指名）

| 失败点 | 症状 | 兜底 |
|---|---|---|
| kit 里混进业务 | 所有仓都被迫理解某个域的概念 | 业务词扫描（QT-05 验收） |
| kit 里放表结构 | 出现第二个写者 | code review：`db` 包禁 `CREATE TABLE` |
| 各自实现日志/HTTP | 观测与错误响应形状分裂 | 统一走 `logging` / `httpserver` |

## 4. 不变量（可被测试钉住的）

1. `go build ./...` 在无其他仓时绿。
2. `grep -riE "backtest|strategy|order|factor|market"` 在 `kit` 内零命中（除注释解释"为什么不放这"）。
3. `db` 包不含 DDL。

## 5. 与契约文档的关系

- 表结构 → `docs/db-schema.md`
- 接口形状 → `docs/api-schema.md`
- 本文件只讲**流向与顺序**，不重复字段定义。
